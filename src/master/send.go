package main

import (
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:ello@206.189.131.249:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue for the tasks
	q, err := ch.QueueDeclare(
		"task_queue", // name of the queue
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Push some IDs into the queue as tasks
	for i := 1; i <= 10; i++ {
		body := strconv.Itoa(i) // Convert the ID to string
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key (queue name)
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // Make message persistent
				ContentType:  "text/plain",
				Body:         []byte(body),
			},
		)
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent ID %s", body)
		time.Sleep(1 * time.Second) // Add delay to simulate task distribution
	}
}
