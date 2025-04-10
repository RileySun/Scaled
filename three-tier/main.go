package main

import(
	
)

func main() {
	data := NewDataService()
	logic := NewLogicService(data)
	handler := NewHandlerService(logic)
	
	handler.StartService()
}