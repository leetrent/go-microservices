package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	logSnippet := "\n[broker-service][event][NewEventEmitter] =>"

	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		log.Printf("%s (ERROR-emitter.setup): %s", logSnippet, err.Error())
		return Emitter{}, err
	}
	log.Printf("%s (SUCCESS-emitter.setup)", logSnippet)

	return emitter, nil
}

func (e *Emitter) setup() error {
	logSnippet := "\n[broker-service][event][setup] =>"

	channel, err := e.connection.Channel()
	if err != nil {
		log.Printf("%s (ERROR-emitter.connection.Channel): %s", logSnippet, err.Error())
		return err
	}
	defer channel.Close()
	log.Printf("%s (SUCCESS-emitter.connection.Channel)", logSnippet)

	return nil
}

func (e *Emitter) Push(event string, severity string) error {
	logSnippet := "\n[broker-service][event][Push] =>"

	channel, err := e.connection.Channel()
	if err != nil {
		log.Printf("%s (ERROR-emitter.connection.Channel): %s", logSnippet, err.Error())
		return err
	}
	defer channel.Close()
	log.Printf("%s (SUCCESS-emitter.connection.Channel)", logSnippet)

	err = channel.Publish(
		"logs_topic", // exchange
		severity,     // key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		}, // message
	)

	if err != nil {
		log.Printf("%s (ERROR-channel.Publish): %s", logSnippet, err.Error())
		return err
	}
	log.Printf("%s (SUCCESS-channel.Publish)", logSnippet)

	return nil
}
