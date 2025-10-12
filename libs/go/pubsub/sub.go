package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType QueueType,
	table amqp.Table,
	handler func(T) AckType,
) error {
	unmarshaller := func(data []byte) (T, error) {
		var temp T
		if err := json.Unmarshal(data, &temp); err != nil {
			return temp, err
		}
		return temp, nil
	}

	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		table,
		handler,
		unmarshaller,
	)
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType QueueType,
	table amqp.Table,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	ch, queue, err := declareAndBind(conn, exchange, queueName, key, queueType, table)
	if err != nil {
		return err
	}

	if err := ch.Qos(10, 0, false); err != nil {
		return err
	}

	delivery, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func(delivery <-chan amqp.Delivery) {
		for v := range delivery {
			temp, err := unmarshaller(v.Body)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			ack := handler(temp)
			switch ack {
			case ACK:
				if err := v.Ack(false); err != nil {
					log.Printf("error: %v", err)
				}
			case NACK_REQUEUE:
				if err := v.Nack(false, true); err != nil {
					log.Printf("error: %v", err)
				}
			case NACK:
				if err := v.Nack(false, false); err != nil {
					log.Printf("error: %v", err)
				}
			default:
				fmt.Println("unknown ack type")
			}
		}
	}(delivery)

	return nil
}
