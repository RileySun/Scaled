package main

import(

)

func main() {	
	//Reciever
	reciever := NewReciever()
	defer reciever.Close()
	reciever.Start()
}