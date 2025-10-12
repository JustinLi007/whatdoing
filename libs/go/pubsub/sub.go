package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler[T any] func(T) AckType
type MessageUnmarshaller[T any] func([]byte) (T, error)

func SubscribeJSON[T any](
	ch *amqp.Channel,
	exchange,
	queueName,
	key string,
	queueType QueueType,
	table amqp.Table,
	handler MessageHandler[T],
) error {
	unmarshaller := func(data []byte) (T, error) {
		var temp T
		if err := json.Unmarshal(data, &temp); err != nil {
			return temp, err
		}
		return temp, nil
	}

	return subscribe(
		ch,
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
	ch *amqp.Channel,
	exchange,
	queueName,
	key string,
	queueType QueueType,
	args amqp.Table,
	handler MessageHandler[T],
	unmarshaller MessageUnmarshaller[T],
) error {
	queue, err := QueueDeclareAndBind(ch, exchange, queueName, key, queueType, args)
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
