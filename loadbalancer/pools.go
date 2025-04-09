package loadbalancer

import(
	"sync"
	
	"github.com/RileySun/Scaled/utils"
)

//Pool Interface
type ServerPool interface {
	Get() []Backend
	NextValidPeer() (Backend, error)
	Add(Backend)
	PoolSize() int
}

//Round Robin
type roundRobinPool struct {
	backends []Backend
	mux sync.RWMutex
	current int
	
	ServerPool
}

func (s *roundRobinPool) Get() []Backend {
	return s.backends
}

func (s *roundRobinPool) Rotate() Backend {
	s.mux.Lock()
	s.current = (s.current + 1) % s.PoolSize()
	s.mux.Unlock()
	return s.backends[s.current]
}

func (s *roundRobinPool) NextValidPeer() Backend {
	for i := 0; i < s.PoolSize(); i++ {
		next := s.Rotate()
		if next.IsAlive() {
			return next
		}
	}
	return nil
}

func (s *roundRobinPool) Add(newBackend Backend) {
	s.mux.Lock()
	s.backends = append(s.backends, newBackend)
	s.mux.Unlock()
}

func (s *roundRobinPool) PoolSize() int {
	return len(s.backends)
}

//Weighted Round Robin
type weightedRoundRobinPool struct {
	roundRobinPool
	ServerPool
}

func (s *weightedRoundRobinPool) TotalWeight() int {
	var weight int
	for _, b := range s.backends {
		weight += b.GetWeight()
	}
	return weight
}

func (s *weightedRoundRobinPool) CumulativeWeights() []int {
	cumulativeWeights := make([]int, len(s.backends)) //slightly faster ops
	cumulativeWeights[0] = s.backends[0].GetWeight()
	for i, b := range s.backends {
		cumulativeWeights[i] = cumulativeWeights[i - 1] + b.GetWeight();
	}
	return cumulativeWeights
}

func (s *weightedRoundRobinPool) Rotate() Backend {
	//Get Index based off weights
	var index int
	random := utils.GetRandomInt(0, s.TotalWeight());
	for i, w := range s.CumulativeWeights() {
		if random < w {
			index = i
			break
		}
	}

	//Lock mux, then change current index
	s.mux.Lock()
	s.current = index
	s.mux.Unlock()
	
	//return
	return s.backends[s.current]
}


//Least Connections
type lcPool struct {
	backends []Backend
	mux	sync.RWMutex
	
	ServerPool
}

func (s *lcPool) Get() []Backend {
	return s.backends
}

func (s *lcPool) GetNextValidPeer() Backend {
	var least Backend
	
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

func (s *lcPool) Add(newBackend Backend) {
	s.mux.Lock()
	s.backends = append(s.backends, newBackend)
	s.mux.Unlock()
}

func (s *lcPool) PoolSize() int {
	return len(s.backends)
}