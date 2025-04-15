package main

import (
	"log"
	"time"
	"bytes"

	amqp "github.com/rabbitmq/amqp091-go"
)

//Struct
type Reciever struct {
	connection *amqp.Connection
	channel *amqp.Channel
	queue amqp.Queue
}

//Create
func NewReciever() *Reciever {
	reciever := &Reciever{}
	
	//Error Re-Use
	var err error

	//Create Queue Connection
	reciever.connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	reciever.onError(err, "Message Queue Connection Error")
	
	//Create Channel
	reciever.channel, err = reciever.connection.Channel()
	reciever.onError(err, "Channel Open Error")
	
	//Create Queue
	reciever.queue, err = reciever.channel.QueueDeclare("Main", false, false, false, false, nil)
	reciever.onError(err, "Queue Open Error")
	
	return reciever
}

//Actions
func (r *Reciever) Start() {
	msgs, err := r.channel.Consume(r.queue.Name, "", true, false, false, false, nil)
	r.onError(err, "Consume Register Error")
	
	var forever chan struct{}
	
	go func() {
		for d := range msgs {
			distributedWork(d.Body)
		}
	}()
	
	log.Printf("* Message Reciever Running *")
	
	<-forever
}

func distributedWork(message string) {
	//
}

//Utils
func (r *Reciever) Close() {
	r.connection.Close()
	r.channel.Close()
}

func (r *Reciever) onError(err error, message string) {
	if err != nil {
		log.Panicf("%s: %s", message, err)
	}
}