package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// This will be the codes that pushes events to the queue(rabbitmq) [emit an event to rabbitmq]
type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	return declareExchange(channel)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	err = channel.Publish(
		"logs_topic",
		severity, // log.INFO, log.WARNING, log.ERROR (3 severities)
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event), // text message
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
