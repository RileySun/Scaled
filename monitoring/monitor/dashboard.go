package main

import(
	"log"
	"strconv"
	"net/http"
	"html/template"
	
	"github.com/julienschmidt/httprouter"
)

//Structs
type Dashboard struct {
	
}

type Item struct {
	ID, Name, Address, Speed, Memory, Uptime, Health string
}

//Create
func NewDashboard() *Dashboard {
	dash := &Dashboard{
	
	}
	
	return dash
}

//Actions
func (d *Dashboard) getItems() []*Item {
	var items []*Item
	
	//TODO get other servers health
	for i:=0; i<5; i++ {
		//Name, Address, Speed, Memory, Uptime, Health
		item := &Item{
			ID:"S-"+strconv.Itoa(i),
			Name:"Server "+ strconv.Itoa(i),
			Address:"localhost:8080",
			Speed:"100ms", Memory:"10%", 
			Uptime:"1d 20h 10m", Health:"OK",
		}
		items = append(items, item)
	}
	
	return items
}

//Routes
func (d *Dashboard) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl, parseErr := template.ParseFS(HTMLFiles, "html/Dashboard.html")
	if parseErr != nil {
		log.Println("Dashboard Template Parse: ", parseErr)
	}
	
	items := d.getItems()
	
	//Get Status Data
	templateData := struct {
    	Name string
    	Items []*Item
	}{
		"Micro", items,
	}
	
	tmpl.Execute(w, templateData)
}

func (d *Dashboard) Restart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Validate Server ID
	rawID := ps.ByName("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}
	
	//DEBUG - Mock Server Restart
	log.Println("Server '", id, "' restarted")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (d *Dashboard) Shutdown(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//Validate Server ID
	rawID := ps.ByName("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}
	
	//DEBUG - Mock Server Shutdown
	log.Println("Server '", id, "' shutdown")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}