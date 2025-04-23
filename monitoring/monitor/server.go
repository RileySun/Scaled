package main

import(
	"io"
	"log"
	"time"
	"errors"
	"strconv"
	"net/http"
	"encoding/json"
	
	_ "embed"
	
	"github.com/RileySun/Scaled/utils"
)

//go:embed addresses.json
var addresses []byte

//Struct
type Server struct {
	ID, Name, Address, Speed, Memory, Uptime, Health string
}

//Create
func LoadServers(updateType string) []*Server {
	var servers []*Server
	
	//Get Addresses
	var addressList []string
	err := json.Unmarshal(addresses, &addressList)
	if err != nil {
		log.Println(err)
		return servers
	}
	
	for i, a := range addressList {
		server := &Server{
			ID:"S-" + strconv.Itoa(i),
			Name:"Error",
			Address:a,
			Speed:"Error",
			Memory:"Error",
			Uptime:"Error",
			Health:"Error",
		}
		
		servers = append(servers, server)
	}
	
	for _, s := range servers {
		s.Update(updateType)
	}
	
	return servers
}

//Actions
func (s *Server) Update(updateType string) error {
	switch updateType {
		case "microservice":
			return s.UpdateMicro()
		case "container":
			return s.UpdateContainer()
		case "basic":
			return s.UpdateBasic()
		default:
			return errors.New("Invalid Update Type (microservice, container, or basic only)")
	}
}

func (s *Server) UpdateMicro() error {
	return nil
}

func (s *Server) UpdateContainer() error {
	return nil
}

func (s *Server) UpdateBasic() error {	
	startTime := time.Now()
	resp, err := http.Get("http://" + s.Address + "/health")
	if err != nil {
		return err
	}
	endTime := int(time.Since(startTime).Milliseconds())
	timeStr := strconv.Itoa(endTime) + "ms"
	
	body, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return bodyErr
	}
	defer resp.Body.Close()
	
	var data *utils.ServerData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	
	//Check Health
	health := "OK"
		//Response Speed
	if endTime > 400 && endTime < 700 {
		health = "Slow"
	} else if endTime > 700 {
		health = "Error"
	}
		//Memory
	mem, floatErr := strconv.ParseFloat(data.Memory, 32)
	if floatErr != nil {
		health = "Error"
	}
	if mem > 50 && mem < 100 {
		health = "Slow"
	} else if mem > 100 {
		health = "Error"
	}
	
	//Set Props
	s.Name = data.Name
	s.Uptime = data.Uptime
	s.Speed = timeStr
	s.Memory = data.Memory + "MB"
	s.Health = health
	
	return nil
}