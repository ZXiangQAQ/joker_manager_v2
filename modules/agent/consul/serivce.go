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

//func (s *serviceSet) broadcast(ss []*ServiceInstance) {
//	s.services.Store(ss)
//	s.lock.RLock()
//	defer s.lock.RUnlock()
//	for k := range s.watcher {
//		select {
//		case k.event <- struct{}{}:
//		default:
//		}
//	}
//}
