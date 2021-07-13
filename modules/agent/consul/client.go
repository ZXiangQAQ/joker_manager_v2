package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"net/url"
	"strconv"
	"time"
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

	svc.Metadata["version"] = svc.Version

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
				Interval:                       "5s",
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

func (c *Client) Service(ctx context.Context, service string, tag string, index uint64, passingOnly bool) ([]*ServiceInstance, uint64, error) {
	opts := &api.QueryOptions{
		WaitIndex: index,
		WaitTime:  time.Second * 55,
	}
	opts = opts.WithContext(ctx)
	entries, meta, err := c.cli.Health().Service(service, tag, passingOnly, opts)
	if err != nil {
		return nil, 0, err
	}
	var services []*ServiceInstance
	for _, entry := range entries {
		var version string
		var endpoints []string
		for _, addr := range entry.Service.TaggedAddresses {
			endpoints = append(endpoints, addr.Address)
		}
		services = append(services, &ServiceInstance{
			ID:        entry.Service.ID,
			Name:      entry.Service.Service,
			Metadata:  entry.Service.Meta,
			Tags:      entry.Service.Tags,
			Version:   version,
			Endpoints: endpoints,
		})
	}
	return services, meta.LastIndex, nil
}
