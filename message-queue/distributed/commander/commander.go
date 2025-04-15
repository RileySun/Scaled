package main

import (
	"os"
	"log"
	"time"
	"strings"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

//Struct
type Sender struct {
	connection *amqp.Connection
	channel *amqp.Channel
	queue amqp.Queue
}

//Create
func NewSender() *Sender {
	//Create
	sender := &Sender{}
	
	//Error Re-Use
	var err error

	//Create Queue Connection
	sender.connection, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	sender.onError(err, "Message Queue Connection Error")
	
	//Create Channel
	sender.channel, err = sender.connection.Channel()
	sender.onError(err, "Channel Open Error")
	
	//Set Fair Dispatch
	err = reciever.channel.Qos(1, 0, false)
	sender.onError(err, "QoS Setup Error")
	
	//Create Queue
	sender.queue, err = sender.channel.QueueDeclare("Main", true, false, false, false, nil)
	sender.onError(err, "Queue Open Error")
	
	return sender
}

//Actions
func(s *Sender) Send(message string) {
	//Create Timeout Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := s.channel.PublishWithContext(ctx, "", s.queue.Name, false, false, amqp.Publishing{DeliveryMode:amqp.Persistent, ContentType:"text/plain", Body:[]byte(message),})
	s.onError(err, "Message Publish Error")
}

//Utils
func (s *Sender) Close() {
	s.connection.Close()
	s.channel.Close()
}

func (s *Sender) onError(err error, message string) {
	if err != nil {
		log.Panicf("%s: %s", message, err)
	}
}