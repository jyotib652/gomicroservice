package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type
		true,         // durable => is it durable?
		false,        // autodeleted => do you delete it when you are done with it?
		false,        // internal => is this an exchange just used internally? We are going to use this exchange between microservices. [noLocal means internal]
		false,        // no-wait => it's not important, so no need to worry about this.
		nil,          // arguments => we're not going to have any specific argument for this exchange.
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name of the queue => we're not giving it any name
		false, // durable => is it durable? no, just get rid off it when you're done with it
		false, // autodeleted => do you delete it when you are done with it? No.
		true,  // exclusive => yes, this is exclusive channel for my current operations. Don't share it around.
		false, // no-wait => it's not important, so no need to worry about this.
		nil,   // arguments => we're not going to have any specific argument for this exchange.
	)
}
