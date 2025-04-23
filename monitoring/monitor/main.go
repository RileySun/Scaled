package main

import(
	"os"
	"os/signal"
	"syscall"
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
	router.GET("/restart/:id", dash.Restart)
	router.GET("/shutdown/:id", dash.Shutdown)
	router.ServeFiles("/static/*filepath", http.Dir("html/static"))
		
	//Listen with server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	srv := utils.StartHTTPServer(router, "8080")
	<-done
	
	//Context for shutting down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}