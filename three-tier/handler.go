package main

import(
	"os"
	"time"
	"context"
	"syscall"
	"net/http"
	"os/signal"
	
	"github.com/RileySun/Scaled/utils"
	
	"github.com/julienschmidt/httprouter"
)

//Struct
type HandlerService struct {
	router *httprouter.Router
	server *http.Server
	
	logicService *LogicService
}

//Create
func NewHandlerService(ls *LogicService) *HandlerService {
	//Create handler
	handler := &HandlerService{
		logicService:ls,
		router:httprouter.New(),
	}
	
	//set routes
	handler.router.GET("/user/status/get/:id",  handler.userGetStatus) 
	handler.router.GET("/user/status/set/:id/:newStatus", handler.userSetStatus)
	
	return handler
}

//Start Handler Server
func (s *HandlerService) StartService() {
	//Create server (can close gracefully with Shutdown())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s.server = utils.StartHTTPServer(s.router, "8080")
	<-done
	
	//Context for shutting down
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		//Graceful shutdown functions here
		cancel()
	}()
	if err := s.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

//Routes
func (s *HandlerService) userGetStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("id")
	status, err := s.logicService.GetStatus(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	
	//Write header and Respond
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}

func (s *HandlerService) userSetStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("id")
	newStatus := ps.ByName("newStatus")
	
	err := s.logicService.SetStatus(userID, newStatus)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	
	//Write header and Respond
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}