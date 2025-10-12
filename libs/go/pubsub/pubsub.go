package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ExchangeType int
type QueueType int
type AckType int

const (
	EXCHANGE_TYPE_DURABLE ExchangeType = iota
	EXCHANGE_TYPE_TRANSIENT
)

const (
	QUEUE_TYPE_DURABLE QueueType = iota
	QUEUE_TYPE_TRANSIENT
)

const (
	ACK AckType = iota
	NACK_REQUEUE
	NACK
)

func ExchangeDeclare(
	ch *amqp.Channel,
	name,
	kind string,
	exchangeType ExchangeType,
	args amqp.Table,
) error {
	return ch.ExchangeDeclare(
		name,
		kind,
		exchangeType == EXCHANGE_TYPE_DURABLE,
		exchangeType != EXCHANGE_TYPE_DURABLE,
		false,
		false,
		args,
	)
}

func QueueDeclareAndBind(
	ch *amqp.Channel,
	exchange,
	queueName,
	key string,
	queueType QueueType,
	args amqp.Table,
) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueName,
		queueType == QUEUE_TYPE_DURABLE,
		queueType != QUEUE_TYPE_DURABLE,
		queueType != QUEUE_TYPE_DURABLE,
		false,
		args,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	if err := ch.QueueBind(queue.Name, key, exchange, false, nil); err != nil {
		return amqp.Queue{}, err
	}

	return queue, nil
}
