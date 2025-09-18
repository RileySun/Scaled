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
		failures:atomiccounter.MakeInt64(),
	}
	
	return circuit
}

func (c *Circuit) Execute(ctx context.Context, action func() (any, error)) (interface{}, error) {
	if c.failures.Read() >= c.FailLimit {
		ctx.Done()
		ctx.Err()
		
		c.mutex.Lock()
		c.state = "Broken"
		c.mutex.Unlock()
		
		if c.OnBreak != nil {
			c.OnBreak()
		}
		
		return nil, errors.New("Circuit breaker flipped shutting down")
	}
	
	result, err := action()
	
	if err != nil {
		c.mutex.Lock()
		c.failures.Inc()
		c.lastAttempt = time.Now()
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
	failures := c.failures.Read()
	c.mutex.Unlock()
	return failures
}