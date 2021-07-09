package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"net/url"
	"strconv"
)

type Client struct {
	cli *api.Client
}

func NewClient(cli *api.Client) *Client {
	return &Client{cli: cli}
}

func (c *Client) Register(ctx context.Context, svc *ServiceInstance) error {
	addresses := make(map[string]api.ServiceAddress)
	var addr string
	var port uint64

	for _, endpoint := range svc.Endpoints {
		raw, err := url.Parse(endpoint)
		if err != nil {
			return err
		}
		addr = raw.Hostname()
		port, _ = strconv.ParseUint(raw.Port(), 10, 16)
		addresses[raw.Scheme] = api.ServiceAddress{Address: endpoint, Port: int(port)}
	}

	asr := &api.AgentServiceRegistration{
		ID:              svc.ID,
		Name:            svc.Name,
		Tags:            svc.Tags,
		Port:            int(port),
		Address:         addr,
		TaggedAddresses: addresses,
		Meta:            svc.Metadata,
		Checks: []*api.AgentServiceCheck{
			{
				TCP:                            fmt.Sprintf("%s:%d", addr, port),
				Interval:                       "10s",
				DeregisterCriticalServiceAfter: "180s",
			},
		},
	}

	ch := make(chan error)
	go func() {
		err := c.cli.Agent().ServiceRegister(asr)
		ch <- err
	}()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ch:
	}

	return err
}

func (c *Client) Deregister(ctx context.Context, serviceID string) error {
	ch := make(chan error)
	go func() {
		err := c.cli.Agent().ServiceDeregister(serviceID)
		ch <- err
	}()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-ch:
	}

	return err
}
