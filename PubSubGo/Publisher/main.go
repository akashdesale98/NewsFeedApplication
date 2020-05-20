package main

import (
	"NewsFeedApplication/PubSubGo/models"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial(models.Config.AMQPConnectionURL)
	handleError(err, "Can't Connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	for {

		rand.Seed(time.Now().UnixNano())

		addTask := models.AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
		body, err := json.Marshal(addTask)

		if err != nil {
			handleError(err, "Error Encoding JSON")
		}

		err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

		if err != nil {
			log.Fatalf("Error Publishing message: %s", err)
		}
	}

	// log.Printf("AddTask: %d+%d", addTask.Number1, addTask.Number2)
}
