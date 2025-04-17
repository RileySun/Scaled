package main

import (
	"log"
	"time"
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
	
	//Create Exchange
	err = sender.channel.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	sender.onError(err, "Exchange Declare Error")
	
	return sender
}

//Actions
func(s *Sender) Send(message string) {
	//Create Timeout Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := s.channel.PublishWithContext(ctx, "logs", "", false, false, amqp.Publishing{ContentType:"text/plain", Body:[]byte(message),})
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