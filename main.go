package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

const SERVER = "amqp://root:admin@127.0.0.1:5672/"
const EXCHANGE_NAME = "doro_exchange"
const EXCHANGE_TYPE = "direct"
const ROUTING_KEY = "#"

type bind struct {
	queue       string
	routingKeys []string
	callback    string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func setupQueues(channel *amqp.Channel, bindings []bind) {
	for _, bind := range bindings {
		_, err := channel.QueueDeclare(bind.queue, true, false, false, false, nil)
		failOnError(err, "queue.declare")

		for _, routingKey := range bind.routingKeys {
			err = channel.QueueBind(bind.queue, routingKey, EXCHANGE_NAME, false, nil)
			failOnError(err, fmt.Sprintf("queue.bind to routiangKey <%s>", routingKey))
		}
	}
}

func main() {
	//bindingConfig := []bind{
	//	{"reporting-queue-1", []string{"scheduled.event", ""}, ""},
	//	{"reporting-queue-1", []string{"rescheduled.event"}, ""},
	//}

	connection, err := amqp.Dial(SERVER)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()
	//setupQueues(channel, bindingConfig)

	queue, err := channel.QueueInspect("doro-testing")
	if err != nil {
		log.Printf("Queue doesn't exist")
		os.Exit(0)
	}
	log.Printf("Hay %d mensajes en la cola", queue.Messages)

	messages, err := channel.Consume(
		"doro-testing",
		"no-name-consumer",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
		os.Exit(1)
	}

	go func() {
		for msg := range messages {
			// msg.Ack(false)
			log.Printf("Message consumed: %s", msg.Body)
		}
	}()

	//queue, err := channel.QueueDeclare(
	//	"retaildash-scrapy", // name
	//	false,               // durable
	//	false,               // delete when unused
	//	false,               // exclusive
	//	false,               // no-wait (wait time for processing)
	//	nil,                 // arguments
	//)
	//	if err != nil {
	//		log.Printf("Failed to declare a queue")
	//	}

	/*
		messages, err := channel.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")
		for d := range messages {
			log.Printf("Received a message: %s", d.Body) //any kind of further processing code
		}
	*/
}
