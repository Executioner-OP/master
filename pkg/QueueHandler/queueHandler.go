package queueHandler

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
)

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Initialize RabbitMQ connection, channel, and declare a durable task queue
func Init(URI string) error {
	var err error
	conn, err = amqp.Dial(URI)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	// Declare a durable task queue
	q, err = ch.QueueDeclare(
		"task_queue", // name of the queue
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	fmt.Println("Connected to RabbitMQ")

	return nil
}

func AddToQueue(taskID string) {

	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "text/plain",
			Body:         []byte(taskID),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent task ID: %s", taskID)
}

func Cleanup() {
	// Close the channel and connection when no longer needed
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
