package pubsub

import (
	"service-anime/internal/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type subscriber struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	url  string
}

var subscriberInstance *subscriber

func (s *subscriber) connect() {
	conn, err := amqp.Dial(s.url)
	utils.RequireNoError(err, "Error: publisher failed to establish a connection")
	s.conn = conn

	ch, err := conn.Channel()
	utils.RequireNoError(err, "Error: publisher failed to create a channel")
	s.ch = ch
}

func (s *subscriber) declareExchange(name, kind string, durable bool, args amqp.Table) error {
	return s.ch.ExchangeDeclare(
		name,     // name
		kind,     // kind
		durable,  // durable
		!durable, // auto delete
		false,    // internal
		false,    // no wait
		args,     // args
	)
}

func (s *subscriber) declareQueue(name string, durable bool, args amqp.Table) (amqp.Queue, error) {
	return s.ch.QueueDeclare(
		name,     // name
		durable,  // durable
		!durable, // auto delete
		!durable, // exclusive
		false,    // no wait
		args,     // args
	)
}

func (s *subscriber) bindQueue(name, key, exchange string, args amqp.Table) error {
	return s.ch.QueueBind(
		name,     // name
		key,      // key
		exchange, // exchange
		false,    // no wait
		args,     // args
	)
}
