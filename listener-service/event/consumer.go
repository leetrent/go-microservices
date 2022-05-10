package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConsumer(pConn *amqp.Connection) (Consumer, error) {
	logSnippet := "\n[listener-service][event][consumer][NewConsumer] =>"

	consumer := Consumer{
		conn: pConn,
	}

	err := consumer.setup()
	if err != nil {
		log.Printf("%s (ERROR-consumer.setup): %s", logSnippet, err.Error())
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	logSnippet := "\n[listener-service][event][consumer][setup] =>"

	channel, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("%s (ERROR-consumer.conn.Channel): %s", logSnippet, err.Error())
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	logSnippet := "\n[listener-service][event][consumer][Listen] =>"

	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Printf("%s (ERROR-consumer.conn.Channel): %s", logSnippet, err.Error())
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		log.Printf("%s (ERROR-declareRandomQueue): %s", logSnippet, err.Error())
		return err
	}

	for index, key := range topics {
		err = ch.QueueBind(
			q.Name,       // name
			key,          // key
			"logs_topic", // exchange
			false,        // noWait
			nil,          // args
		)
		if err != nil {
			log.Printf("%s (ERROR-ch.QueueBind / index=%d): %s", logSnippet, index, err.Error())
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		log.Printf("%s (ERROR-ch.Consume): %s", logSnippet, err.Error())
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			err = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	log.Printf("%s ( INFO - Waiting for message [Exchange, Queue] [logs_topic=%s] )", logSnippet, q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	logSnippet := "\n[listener-service][event][consumer][handlePayload] =>"

	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Printf("%s ( ERROR - logEvent[case:'auth', 'event'] ): %s", logSnippet, err.Error())
		}
	case "auth":
		// TODO: handle authenication processing using RabbitMQ
	case "mail":
		// TODO: handle sending mail using RabbitMQ
	default:
		err := logEvent(payload)
		if err != nil {
			log.Printf("%s ( ERROR - logEvent[case:'default'] ): %s", logSnippet, err.Error())
		}
	}
}

func logEvent(entry Payload) error {
	logSnippet := "\n[listener-service][event][consumer][logEvent] =>"

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		log.Printf("%s (ERROR - json.MarshalIndent): %s", logSnippet, err.Error())
		return err
	}

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("%s (ERROR - http.NewRequest): %s", logSnippet, err.Error())
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("%s (ERROR - client.Do): %s", logSnippet, err.Error())
		return err
	}
	defer response.Body.Close()

	log.Printf("\n%s (response.StatusCode): %d", logSnippet, response.StatusCode)

	if response.StatusCode != http.StatusAccepted {
		log.Printf("%s (ERROR - response.StatusCode): %d", logSnippet, response.StatusCode)
		return errors.New("Invalid HTTP Status Code: " + string(response.StatusCode))
	}

	return nil
}
