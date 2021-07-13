package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"sync"
)

type Registrar interface {
	// Register 注册实例
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister 反注册实例
	Deregister(ctx context.Context, service *ServiceInstance) error
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

type Config struct {
	*api.Config
}

type Registry struct {
	cfg *Config
	cli *Client

	registry map[string]*serviceSet
	lock     sync.RWMutex // protect following
}

func New(apiClient *api.Client) *Registry {
	return &Registry{
		cli: NewClient(apiClient),
	}
}

// Register register service
func (r *Registry) Register(ctx context.Context, svc *ServiceInstance) error {
	return r.cli.Register(ctx, svc)
}

// Deregister deregister service
func (r *Registry) Deregister(ctx context.Context, svc *ServiceInstance) error {
	return r.cli.Deregister(ctx, svc.ID)
}

// GetService return service by name
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
