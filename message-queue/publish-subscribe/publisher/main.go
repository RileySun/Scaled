package main

import(
	
)

func main() {
	//Sender
	sender := NewSender()
	defer sender.Close()
	sender.Send("Publish this message to all subscribers")
}