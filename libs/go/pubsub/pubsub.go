package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueType int
type AckType int

const (
	QUEUE_DURABLE QueueType = iota
	QUEUE_TRANSIENT
)

const (
	ACK = iota
	NACK_REQUEUE
	NACK
)

func declareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType QueueType,
	table amqp.Table,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueType == QUEUE_DURABLE,
		queueType != QUEUE_DURABLE,
		queueType != QUEUE_DURABLE,
		false,
		table,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	if err := ch.QueueBind(queue.Name, key, exchange, false, nil); err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}
