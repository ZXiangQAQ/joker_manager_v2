package agent

import (
	"fmt"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/codes"
	//"google.golang.org/grpc/health/grpc_health_v1"
	//"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Ip        string
	Port      int
	ConsulCfg string
}

type Server struct {
	server *grpc.Server
	config Config
	mu *sync.RWMutex
}

func NewServer(c Config) *Server {
	return &Server{
		config: c,
		mu: &sync.RWMutex{},
	}
}

func (s *Server) Start() {
	config := s.config
	addr := fmt.Sprintf("%s:%d", config.Ip, config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	log.Printf("server listen on tcp: %v", addr)

	s.server = grpc.NewServer()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Printf("%s\n", err)
		}
	}()

	sig := <-ch
	log.Printf("%s", sig)
	s.close()
}

func (s *Server) close() {
	s.server.GracefulStop()
}