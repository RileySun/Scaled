package loadbalancer

import(
	"sync"
	"errors"
)

//Pool Interface
type ServerPool interface {
	Get() []Backend
	NextValidPeer() (Backend, error)
	Add(backend)
	PoolSize() int
}

//Round Robin
type roundRobinPool struct {
	backends []backend
	mux sync.RWMutex
	current int
	
	ServerPool
}

func (s *roundRobinPool) Get() []backend {
	return s.backends
}

func (s *roundRobinPool) Rotate() backend {
	s.mux.Lock()
	s.current = (s.current + 1) % s.PoolSize()
	s.mux.Unlock()
	return s.backends[s.current]
}

func (s *roundRobinPool) NextValidPeer() (backend, error) {
	for i := 0; i < s.PoolSize(); i++ {
		next := s.Rotate()
		if next.IsAlive() {
			return next, nil
		}
	}
	return backend{}, errors.New("No valid peers")
}

func (s *roundRobinPool) Add(newBackend backend) {
	s.mux.Lock()
	s.backends = append(s.backends, newBackend)
	s.mux.Unlock()
}

func (s *roundRobinPool) PoolSize() int {
	return len(s.backends)
}

//Least Connections
type lcPool struct {
	backends []backend
	mux	sync.RWMutex
	
	ServerPool
}

func (s *lcPool) Get() []backend {
	return s.backends
}

func (s *lcPool) GetNextValidPeer() backend {
	var least backend
	
	//Clever trick where you make sure at least one is alive before bothering comparing connection numbers
	for _, b := range s.backends {
		if b.IsAlive() {
			least = b
			break
		}
	}
	
	//Now check which has the least amount of connections
	for _, b := range s.backends {
		//Skip non alive backends
		if !b.IsAlive() {
			continue
		}
		
		//Compare Active Connections
		if least.GetConnections() > b.GetConnections() {
			least = b
		}
	}
	
	return least
}

func (s *lcPool) Add(newBackend backend) {
	s.mux.Lock()
	s.backends = append(s.backends, newBackend)
	s.mux.Unlock()
}

func (s *lcPool) PoolSize() int {
	return len(s.backends)
}