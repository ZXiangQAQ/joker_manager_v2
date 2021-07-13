package consul

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"testing"
	"time"
)

func tcpServer(t *testing.T, lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		fmt.Println("get tcp")
		conn.Close()
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
func TestRegister(t *testing.T) {
	addr := fmt.Sprintf("%s:8081", getLocalIP())
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Errorf("listen tcp %s failed!", addr)
		t.Fail()
	}
	defer lis.Close()

	go tcpServer(t, lis)
	time.Sleep(time.Millisecond * 100)

	cli, err := api.NewClient(&api.Config{Address: "consul.ihr360.com", Scheme: "https"})
	if err != nil {
		t.Fatalf("create consul client failed: %v", err)
	}

	r := New(cli)
	assert.Nil(t, err)
	version := strconv.FormatInt(time.Now().Unix(), 10)
	svc := &ServiceInstance{
		ID:        "test23334",
		Name:      "test-provider",
		Version:   version,
		Metadata:  map[string]string{"app": "joker_manager_agent", "brand": "lb"},
		Endpoints: []string{fmt.Sprintf("tcp://%s?isSecure=false", addr)},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = r.Deregister(ctx, svc)
	assert.Nil(t, err)
	err = r.Register(ctx, svc)
	assert.Nil(t, err)

	// 等待consul健康检查通过
	time.Sleep(time.Second * 6)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	services, _, err := r.cli.Service(ctx, "test-provider", "", 0, true)
	assert.Equal(t, 1, len(services))
	for _, service := range services {
		fmt.Printf("%v\n", service)
	}
}
