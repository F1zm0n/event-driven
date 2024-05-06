package rabbitmq

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func MustMakeQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		os.Getenv("QUEUE_NAME"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return q
}

func MustMakeExchangeBindings(ch *amqp.Channel, q amqp.Queue) {
	err := ch.ExchangeDeclare(
		os.Getenv("EXCHANGE_NAME"),
		"headers",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		q.Name,
		os.Getenv("REGISTER_ROUTING_KEY"),
		os.Getenv("EXCHANGE_NAME"),
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
}
