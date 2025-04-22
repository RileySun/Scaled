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

var shutdown func()

func main() {
	shutdown = startServer()
}


func startServer() func() {
	//Create Router
	router := httprouter.New()
	router.GET("/:message", Handle)
		
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
	return func() {
		srv.Shutdown(ctx)
		cancel()
	}
}

func Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	message := ps.ByName("message")
	
	//w.WriteHeader(http.StatusInternalServerError)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}