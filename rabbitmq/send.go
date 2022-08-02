package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"strconv"
)

var (
	exchangeName = "hello-exchange"
	queueName    = "hello-queue"
	routingKey   = "hello-routing-key"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://jeven:1@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//定义一个exchange
	err = ch.ExchangeDeclare(exchangeName, "fanout", false, false, false, false, nil)
	failOnError(err, "Failed to declare a exchange")

	//定义一个queue
	_, err = ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
	failOnError(err, "Failed to bind exchange and queue")

	var body = "Hello World!~~~~~~~~" + strconv.Itoa(rand.Int())
	err = ch.PublishWithContext(
		context.TODO(),
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
