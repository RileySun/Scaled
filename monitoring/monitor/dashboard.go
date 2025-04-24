package main

import(
	"log"
	"strings"
	"strconv"
	"net/http"
	"html/template"
	
	"github.com/julienschmidt/httprouter"
)

//Struct
type Dashboard struct {
	Servers []*Server
	UpdateType string
}

//Create
func NewDashboard(updateType string) *Dashboard {
	dash := &Dashboard{
		UpdateType:updateType,
		Servers:LoadServers(updateType),
	}
	
	return dash
}

//Actions
func (d *Dashboard) getServers() []*Server {
	var servers []*Server
	
	//TODO get other servers health
	for i:=0; i<5; i++ {
		//Name, Address, Speed, Memory, Uptime, Health
		server := &Server{
			ID:"S-"+strconv.Itoa(i),
			Name:"Server "+ strconv.Itoa(i),
			Address:"localhost:8080",
			Speed:"100ms", Memory:"10%", 
			Uptime:"1d 20h 10m", Health:"OK",
		}
		servers = append(servers, server)
	}
	
	return servers
}

//Routes
func (d *Dashboard) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/Dashboard.html")
	if parseErr != nil {
		log.Println("Dashboard Template Parse: ", parseErr)
	}
	
	for _, s := range d.Servers {
		s.Update(d.UpdateType)
	}
	
	//Get Status Data
	templateData := struct {
    	Name string
    	Servers []*Server
	}{
		strings.Title(d.UpdateType), d.Servers,
	}
	
	tmpl.Execute(w, templateData)
}

func (d *Dashboard) Export(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Validate Server ID
	var server *Server
	valid := false
	id := ps.ByName("id")
	for _, s := range d.Servers {
		if s.ID == id {
			server = s
			valid = true
		}
	}
	if !valid || server == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}
	
	server.Update(d.UpdateType)
	
	if server.Health == "Error" {
		w.WriteHeader(http.StatusInternalServerError)
	}
	
	json := server.Export()
	if json == nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	
	//DEBUG - Mock Server Restart
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename=ServerExport.json")
	w.Write(json)
}

func (d *Dashboard) Shutdown(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Validate Server ID
	var server *Server
	valid := false
	id := ps.ByName("id")
	for _, s := range d.Servers {
		if s.ID == id {
			server = s
			valid = true
		}
	}
	if !valid || server == nil  {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}
	
	server.Shutdown()
	
	//DEBUG - Mock Server Shutdown
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}