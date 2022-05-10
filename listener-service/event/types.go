package event

import amqp "github.com/rabbitmq/amqp091-go"

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
