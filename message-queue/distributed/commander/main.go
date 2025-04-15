package main

import(
	
)

func main() {
	//Sender
	sender := NewSender()
	defer sender.Close()
	sender.Send("Distribute Work 1")
	sender.Send("Distribute Work 2")
	sender.Send("Distribute Work 3")
	sender.Send("Distribute Work 4")
	sender.Send("Distribute Work 5")
}