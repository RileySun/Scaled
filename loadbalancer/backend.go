package loadbalancer

import(
	"sync"
	"net/url"
	"net/http"
	"net/http/httputil"
) 

//Interface
type Backend interface {
	SetAlive(bool)
	IsAlive() bool
	GetUrl() *url.URL
	GetConnections() int //Active connections
	GetWeight() int //Weighted round robin
	Serve(http.ResponseWriter, *http.Request)
}

//Struct
type backend struct {
	rProxy *httputil.ReverseProxy //reverse proxy
	mux sync.RWMutex
	url *url.URL
	connections, weight int
	alive bool
	
	Backend
} //Struct orders matter!

//Methods
func (b *backend) SetAlive(newAlive bool) {
	b.mux.Lock()
	b.alive = newAlive
	b.mux.Unlock()
}

func (b *backend) IsAlive() bool {
	return b.alive
}

func (b *backend) GetUrl() *url.URL {
	return b.url
}

func (b *backend) GetConnections() int {
	return b.connections
}

func (b *backend) GetWeight() int {
	return b.weight
}

func (b *backend) Serve(http.ResponseWriter, *http.Request) {
	//Use ur imagination for the actual server code, or see any other project about go servers
}