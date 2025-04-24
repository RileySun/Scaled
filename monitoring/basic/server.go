package main

import(	
	"log"
	"time"
	"context"
	"net/http"
	"encoding/json"
	
	"github.com/julienschmidt/httprouter"
	
	"github.com/RileySun/Scaled/utils"
)

//Structs
type Server struct {
	name, port, status string
	startTime time.Time
	Shutdown func()
	
	debug int
}

//Create
func NewServer(name, newPort string) *Server {
	server := &Server {
		name:name,
		port:newPort,
		status:"offline",
		debug:0,
	}
	
	return server
}

//Actions
func (s *Server) Start() {
	//Start Uptime
	s.startTime = time.Now()
	
	//Create Router
	router := httprouter.New()
	router.GET("/main/:message", s.Handle)
	router.GET("/health", s.Health)
	router.GET("/shutdown", s.Close)
	
	//Get Context and Set Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	s.Shutdown = func() {
		log.Println("Shutting down server...")
		cancel()
		s.status = "offline"
	}
	
	//Start Server
	s.status = "online"
	log.Println("Server Started")
	utils.StartHTTPServer(ctx, s.port, router)
}

//Routes
func (s *Server) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	message := ps.ByName("message")
	
	//w.WriteHeader(http.StatusInternalServerError)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Get Memory
	memoryStr := utils.GetMemory()
	
	//Get Uptime
	rawUptime := time.Since(s.startTime)
	uptimeStr := utils.FormatUptime(rawUptime)
	
	//Coallate
	data := &utils.ServerData {
		Name:s.name,
		Memory:memoryStr,
		Uptime:uptimeStr,
	}
	
	//Marshall
	jsonStr, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error - Can not check health status"))
		return
	}
	
	//Send
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (s *Server) Close(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.Shutdown()
}