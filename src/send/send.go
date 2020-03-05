package main

import (
	"log"

	"github.com/streadway/amqp"
)

// Helper fn to check errs
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// connecto to rabbitmq server and catch err
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	// close connection on fn return
	defer conn.Close()

	// create a connection channel catching err
	channel, err := conn.Channel()
	failOnError(err, "Failed to open connection channel")

	// close channel on return
	defer channel.Close()

	// Declare queue where we'll publish messages
	queue, err := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // autodel
		false,   // exclusive
		false,   // no-wait
		nil,     // args-table
	)
	failOnError(err, "Failed to decare queue")

	// Message body
	body := "Hello, World!"

	// publish
	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish message")
}
