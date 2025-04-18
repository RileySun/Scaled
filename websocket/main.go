package main

import(
	"os"
	"os/signal"
	"syscall"
	"context"
	
	"github.com/julienschmidt/httprouter"
	
	"github.com/RileySun/Scaled/utils"
)

func main() {
	//Create App
	app := NewApp()
	
	//Create Router
	router := httprouter.New()
	router.GET("/websocket", app.Handle)
	
	//Start App
	go app.Listen()
		
	//Listen with server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	srv := utils.StartHTTPServer(router, "8080")
	<-done
	
	//Context for shutting down
	ctx, cancel := context.WithCancel(context.Background())
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	
	//Bind server shutdown to app.Close()
	app.Close = func() {
		srv.Shutdown(ctx)
		cancel()
	}
}