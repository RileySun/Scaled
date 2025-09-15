package circuit

import(
	"sync"
	"time"
	"errors"
	"context"
	
	"github.com/chen3feng/atomiccounter"
)

type Circuit struct {
	FailLimit int64
	
	state string
	mutex sync.Mutex
	lastAttempt time.Time
	failures atomiccounter.Int64
	
	//External
	OnBreak func()
}

func NewCircuit(maxFailLimit int64) *Circuit {
	circuit := &Circuit{
		FailLimit:maxFailLimit,
		state:"Working",
	}
	
	return circuit
}

func (c *Circuit) Execute(ctx context.Context, action interface{}, err error) (interface{}, error) {
	if c.failures.Read() > c.FailLimit {
		ctx.Done()
		ctx.Err()
		
		c.mutex.Lock()
		c.state = "Broken"
		c.mutex.Unlock()
		
		if c.OnBreak != nil {
			c.OnBreak()
		}
		
		return nil, Errors.New("Circuit breaker flipped shutting down")
	}
	
	result, err := req()
	
	if err != nil {
		c.mutex.Lock()
		c.failures.Inc()
		c.LastAttempt = time.Now()
		c.mutex.Unlock()
	}
	
	return result, err
}

func (c *Circuit) State() string {
	c.mutex.Lock()
	state := c.state
	c.mutex.Unlock()
	return state
}

func (c *Circuit) Failures() int64 {
	c.mutex.Lock()
	failures := c.failures
	c.mutex.Unlock()
	return failures
}