package pooling

import(
	
)

type Pool struct {
	clients chan *Client
}

func NewPool(poolSize int) *Pool {
	pool := &Pool{
		clients:make(chan *Client, poolSize),
	}
	
	for i:=0; i<poolSize; i++ {
		pool.clients <- NewClient()
	}
	
	return pool
}

func (p *Pool) Acquire() *Client {
	return <-p.clients
}

func (p *Pool) Release(client *Client) {
	p.clients <- client
}