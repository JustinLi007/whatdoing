package pubsub

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/JustinLi007/whatdoing/libs/go/configs"
	"github.com/JustinLi007/whatdoing/libs/go/utils"
	"github.com/JustinLi007/whatdoing/progress/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber interface {
	Start(ctx context.Context)
}

type subscriber struct {
	conn                 *amqp.Connection
	ch                   *amqp.Channel
	url                  string
	interval             time.Duration
	animeProgressService database.ServiceAnimeProgress
}

var subscriberInstance *subscriber

func NewSubscriber(c *configs.Config) Subscriber {
	connStr := c.Get("DB_URL")
	if connStr == "" {
		utils.RequireNoError(errors.New("invalid db url"), "error")
	}

	db, err := database.NewDb(connStr)
	utils.RequireNoError(err, "error: publisher failed to connect to database")

	animeProgressService := database.NewServiceProgress(db)
	return newSubscriber(animeProgressService)
}

func newSubscriber(animeProgressService database.ServiceAnimeProgress) Subscriber {
	if subscriberInstance != nil {
		return subscriberInstance
	}

	newSubscriber := &subscriber{
		conn:                 nil,
		ch:                   nil,
		url:                  "amqp://guest:guest@whatdoing-msg-broker:5672/",
		interval:             time.Second * 5,
		animeProgressService: animeProgressService,
	}

	newSubscriber.connect()
	newSubscriber.declareExchanges()

	subscriberInstance = newSubscriber
	return subscriberInstance
}

func (s *subscriber) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go s.start(ctx, wg)
	wg.Wait()
}

func (s *subscriber) start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(s.interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// TODO: consume from queue
			fmt.Println("consume message...")
			timer.Reset(s.interval)
		case <-ctx.Done():
			return
		}
	}
}

func (s *subscriber) connect() {
	conn, err := amqp.Dial(s.url)
	utils.RequireNoError(err, "Error: subscriber failed to establish a connection")
	s.conn = conn

	ch, err := conn.Channel()
	utils.RequireNoError(err, "Error: subscriber failed to create a channel")
	s.ch = ch
}

func (s *subscriber) declareExchanges() {
	err := s.declareExchange("whatdoing", "topic", true, nil)
	utils.RequireNoError(err, "error: subscriber failed to declare an exchange")
	s.declareAndBind("whatdoing", "anime.*")
}

func (s *subscriber) declareAndBind(exchange, key string) {
	q, err := s.declareQueue(
		"",
		false,
		nil,
	)
	utils.RequireNoError(err, "error: subscriber failed to declare queue")

	err = s.bindQueue(
		q.Name,
		key,
		exchange,
		nil,
	)
	utils.RequireNoError(err, "error: subscriber failed to bind queue to exchange")
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
