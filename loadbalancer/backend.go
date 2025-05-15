package loadbalancer

import(
	"sync"
	"context"
	"net/url"
	"net/http"
	"net/http/httputil"
	
	"github.com/julienschmidt/httprouter"
	"github.com/RileySun/Scaled/utils"
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

func (b *backend) SetUrl(url *url.URL) {
	b.mux.Lock()
	b.url = url
	b.mux.Unlock()
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

func (b *backend) Serve(ctx context.Context) {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write("Hello!")
	})
	utils.StartHTTPServer(ctx, b.url, router)
} //Mock Service