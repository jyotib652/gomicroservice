package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer will be used, to receive events from queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// creating an instance of "Consumer"
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// open up a channel and declare an exchange (these are specific to amqp)
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// Payload will be used to push events to (queue) rabbitmq
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen listens to the queue for specific topics
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	// 'q' for queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Now, we want to do this forever; we want to consume everything that is coming from rabbitmq until we exit this application
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			// here, we're firing off another go-routine within a go-routine. We're doing this just to make things as fast as possible
			go handlePayload(payload)
		}
	}()

	fmt.Printf("waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	// since, "forever" is a channel, so the below line will keep it going for ever.
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event": // "log" or "event"
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}

	case "auth":
		// authenticate

	// you can have as many cases as you want here, but naturally you'll have to write the logic
	// to connect to a given microservice

	default:
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func logEvent(entry Payload) error {
	// create some json we'll send to the logger microservice
	jsonData, _ := json.Marshal(entry)

	// call the service
	// Now, to call the auth service, we need to build the equivalent
	// of http request
	// NewRequest params => method, uri(this uri would the name of the service in the docker-compose.yml file/
	//	endpoint, which we want to hit.)
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	// create a http client
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
