package consul

import (
	"sync"
	"sync/atomic"
)

type serviceSet struct {
	serviceName string
	//watcher     map[*watcher]struct{}
	services *atomic.Value
	lock     sync.RWMutex
}
