package main

import(
	"context"
	"net/http"
	
	"github.com/julienschmidt/httprouter"
	
	"github.com/RileySun/Scaled/utils"
)

func main() {
	dash := NewDashboard("basic")
	
	//Create Router
	router := httprouter.New()
	router.GET("/", dash.Handle)
	router.GET("/Export/:id", dash.Export)
	router.GET("/shutdown/:id", dash.Shutdown)
	router.ServeFiles("/static/*filepath", http.Dir("html/static"))
		
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	utils.StartHTTPServer(ctx, "8080", router)
}