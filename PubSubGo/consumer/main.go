package main

import (
	"NewsFeedApplication/PubSubGo/models"
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	// fmt.Println("heloo world")
	conn, err := amqp.Dial(models.Config.AMQPConnectionURL)
	handleError(err, "Can't Connect to AMQP")
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Couldn't configure QoS")

	messageChannel, err := amqpChannel.Consume(queue.Name, "", false, false, false, false, nil)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {

		log.Printf("Consumer ready, PID: %d", os.Getpid())

		for d := range messageChannel {
			log.Printf("Recieved a message: %s", d.Body)

			addTask := &models.AddTask{}

			err := json.Unmarshal(d.Body, addTask)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			log.Printf("Result of %d + %d is : %d", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	<-stopChan

}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
