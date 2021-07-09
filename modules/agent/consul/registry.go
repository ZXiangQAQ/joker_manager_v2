package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"sync"
	"time"
)

type Config struct {
	*api.Config
}

type Registry struct {
	cfg      *Config
	cli      *Client
	registry map[string]*serviceSet

	lock sync.RWMutex
}

type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Version  string            `json:"Version"`
	Metadata map[string]string `json:"metadata"`
	Tags     []string          `json:"tags"`
	// Endpoints is endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoints []string `json:"Endpoints"`
}

func New(apiClient *api.Client) *Registry {
	return &Registry{
		cli: NewClient(apiClient),
	}
}

func (r *Registry) Register(ctx context.Context, svc *ServiceInstance) error {
	return r.cli.Register(ctx, svc)
}

func (r *Registry) Deregister(ctx context.Context, svc *ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

func (r *Registry) GetService(ctx context.Context, name string) (services []*ServiceInstance, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	set := r.registry[name]
	if set == nil {
		return nil, fmt.Errorf("service %s not resolved in registry", name)
	}
	ss, _ := set.services.Load().([]*ServiceInstance)
	if ss == nil {
		return nil, fmt.Errorf("service %s not found in registry", name)
	}
	for _, s := range ss {
		services = append(services, s)
	}
	return
}

func (r *Registry) resolve(ss *serviceSet) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	services, idx, err := r.cli.Service(ctx, ss.serviceName, 0, true)
	cancel()
	if err == nil && len(services) > 0 {
		ss.broadcast(services)
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		tmpService, tmpIdx, err := r.cli.Service(ctx, ss.serviceName, idx, true)
		cancel()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		if len(tmpService) != 0 && tmpIdx != idx {
			services = tmpService
			ss.broadcast(services)
		}
		idx = tmpIdx
	}
}
