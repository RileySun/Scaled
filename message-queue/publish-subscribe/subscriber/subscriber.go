package main

import (
	"log"

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
	
	//Create Exchange
	err = reciever.channel.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	reciever.onError(err, "Exchange Declare Error")
	
	//Set Fair Dispatch
	err = reciever.channel.Qos(1, 0, false)
	reciever.onError(err, "QoS Setup Error")
	
	//Create Queue
	reciever.queue, err = reciever.channel.QueueDeclare("", false, false, true, false, nil)
	reciever.onError(err, "Queue Open Error")
	
	//Bind Queue to Exchange
	err = reciever.channel.QueueBind(reciever.queue.Name, "", "logs", false, nil)
	reciever.onError(err, "Queue Bind Error")
	
	return reciever
}

//Actions
func (r *Reciever) Start() {
	msgs, err := r.channel.Consume(r.queue.Name, "", false, false, false, false, nil)
	r.onError(err, "Consume Register Error")
	
	var forever chan struct{}
	
	go func() {
		for d := range msgs {
			publishedMessage(d.Body)
			d.Ack(false)
		}
	}()
	
	log.Printf("* Message Reciever Running *")
	
	<-forever
}

func publishedMessage(message []byte) {
	log.Println(string(message[:]))
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