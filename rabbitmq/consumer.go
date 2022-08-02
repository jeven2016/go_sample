package main

import "log"

var (
	exchangeName_consumer = "hello-exchange"
	queueName_consumer    = "hello-queue"
	routingKey_consumer   = "hello-routing-key"
	url_consumer          = "amqp://admin:admin@localhost:5672/"
)

func handleError(err error) {
	if err != nil {
		log.Printf("error encountered: %v", err)
	}
}

func main() {
	connection, err := Dial(url_consumer)
	handleError(err)

	channel, err := connection.Channel()
	handleError(err)

	_, err = channel.Consume(queueName_consumer, "", true, false, false, false, nil)
	handleError(err)
}
