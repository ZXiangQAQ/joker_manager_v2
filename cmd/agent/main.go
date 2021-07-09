package main

import (
	"flag"
	"fmt"
	"joker_manager_v2/modules/agent"
	"os"
)

var (
	v = flag.Bool("version", false, "show version information")
	ip = flag.String("ip", "0.0.0.0", "ipaddress")
	port = flag.Int("port", 23333, "port")
	username = flag.String("conul username", "", "consul username")
	password = flag.String("password", "", "consul BasicAuth password")
)

func init() {
	flag.Parse()
	if *v {
		fmt.Println("version")
		os.Exit(0)
	}
}
func main() {
	srvCfg := agent.Config{
		Ip: *ip,
		Port: *port,
	}
	srv := agent.NewServer(srvCfg)
	srv.Start()
}
