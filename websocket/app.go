package main

import(
	"log"
	"net/http"
	
	"github.com/godbus/dbus/v5"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

//Const
const APPID = "com.sunshine.scaled"
const APP_PATH = "/"

//Struct
type App struct {
	Bus *dbus.Conn
	Conns []*websocket.Conn
	
	done bool
	upgrader websocket.Upgrader
	
	Close func()
}

//Create
func NewApp() *App {
	return &App{
		done:false,
		upgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
	}
}

//Actions
func (a *App) Broadcast(message []byte) {
	for _, c := range a.Conns {
		c.WriteMessage(1, message)
	}
}

func (a *App) Listen() {
	//dBus Connection
	conn, err := dbus.ConnectSessionBus()
	
	if err != nil {
		log.Fatal("dBus connection error: ", err)
	}
	
	//Set dBus
	a.Bus = conn
	
	//Filter Connections
	err = conn.AddMatchSignal(dbus.WithMatchObjectPath(APP_PATH))
	if err != nil {
		log.Fatal("Connection filter error: ", err)
	}
	
	//Create signal channel
	c := make(chan *dbus.Signal, 10)
	
	//Attach signal channel
	conn.Signal(c)
	
	//Broadcast Message
	for v := range c {
		message := v.Body[0].([]byte)
		a.Broadcast(message)
	}
}

func (a *App) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Upgrade connection to websocket instance
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrader error: ", err)
		return
	}
	
	//Add to Connection slice
	a.Conns = append(a.Conns, conn)
	
	//Read Message and Emit
	for {
		if a.done {
			a.Shutdown()
			return
		}
		
		_, b, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket read error: ", err)
		}
		
		if len(b) > 0 {
			a.Bus.Emit(APP_PATH, APPID, b)
		}
	}
}

func (a *App) Shutdown() {
	for _, c := range a.Conns {
		 c.Close()
	}
	
	a.Close()
}