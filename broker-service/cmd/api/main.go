package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

func main() {
	logSnippet := "\n[broker-service][main][main] =>"

	// Connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Printf("%s (EXIT - connect()): %s", logSnippet, err.Error())
		os.Exit(1)
	}
	log.Printf("%s (SUCESS - connect())", logSnippet)

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	logSnippet := "\n[broker-service][main][connect] =>"

	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// Don't continue until rabbitmq is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
