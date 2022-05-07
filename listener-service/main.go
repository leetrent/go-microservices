package main

import (
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	logSnippet := "\n[listener-service][main][main] =>"

	// Connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Printf("%s (ERROR - connect()): %s", logSnippet, err.Error())
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Printf("%s (INFO - Connected to RabbitMQ...)", logSnippet)

	// Start listening for messages

	// Create message consumer

	// Watch the queue and consume events
}

func connect() (*amqp.Connection, error) {
	logSnippet := "\n[listener-service][main][connect] =>"

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// Don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			log.Printf("%s (WARNING - RabbitMQ not yet ready...)", logSnippet)
		} else {
			log.Printf("%s (SUCCESS - Connected to RabbitMQ...)", logSnippet)
			connection = c
			break
		}

		if counts > 5 {
			log.Printf("%s (ERROR - Unable to connect to RabbitMQ): %s", logSnippet, err.Error())
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("%s (INFO - Backing off ...)", logSnippet)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
