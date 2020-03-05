package main

import (
	"log"

	"github.com/streadway/amqp"
)

// fail on err func
// Helper fn to check errs
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// establish connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to establish connection")

	// close on return
	defer conn.Close()

	// create connection channel
	channel, err := conn.Channel()
	failOnError(err, "Failed to create channel")
	defer channel.Close()

	// declare queue to listen messages from
	queue, err := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // autodel
		false,   // exclusive
		false,   // no-wait
		nil,
	)
	failOnError(err, "Failed to declare queue")

	// register as a consumer/listener
	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register as consumer")

	// go channel for receivig messages
	forever := make(chan bool)

	// read the mesages in a goroutine
	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
