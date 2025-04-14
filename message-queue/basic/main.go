package main

import(

)

func main() {
	//Sender
	sender := NewSender()
	defer sender.Close()
	sender.Send("Potato")
	
	//Reciever
	reciever := NewReciever()
	defer reciever.Close()
	reciever.Start()
}