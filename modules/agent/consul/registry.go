package consul

import "github.com/hashicorp/consul/api"

type Config struct {
	*api.Config
}

type Registry struct {
	cfg *Config
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
