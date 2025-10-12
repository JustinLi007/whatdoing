package pubsub

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/JustinLi007/whatdoing/libs/go/config"
	"github.com/JustinLi007/whatdoing/libs/go/pubsub"
	"github.com/JustinLi007/whatdoing/libs/go/util"
	"github.com/JustinLi007/whatdoing/services/progress/internal/database"

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

func NewSubscriber(c *config.Config) Subscriber {
	connStr := c.Get("DB_URL")
	if connStr == "" {
		util.RequireNoError(errors.New("invalid db url"), "error")
	}

	db, err := database.NewDb(connStr)
	util.RequireNoError(err, "error: publisher failed to connect to database")

	animeProgressService := database.NewServiceProgress(db)
	return newSubscriber(animeProgressService)
}

func newSubscriber(animeProgressService database.ServiceAnimeProgress) Subscriber {
	if subscriberInstance != nil {
		return subscriberInstance
	}

	newSubscriber := &subscriber{
		conn:                 nil,
		url:                  "amqp://guest:guest@whatdoing-msg-broker:5672/",
		interval:             time.Second * 5,
		animeProgressService: animeProgressService,
	}

	newSubscriber.connect()
	newSubscriber.declareExchanges()

	subscriberInstance = newSubscriber
	return subscriberInstance
}

func handlerFoo(e *database.AnimeProgress) pubsub.MessageHandler[*database.AnimeProgress] {
	return func(t *database.AnimeProgress) pubsub.AckType {
		log.Println("test message handler...")
		return pubsub.NACK_REQUEUE
	}
}

func (s *subscriber) Start(ctx context.Context) {
	// TODO:
	pubsub.SubscribeJSON(
		s.ch,
		"whatdoing",
		"",
		"anime.*",
		pubsub.QUEUE_TYPE_TRANSIENT,
		nil,
		handlerFoo(nil),
	)
}

func (s *subscriber) connect() {
	conn, err := amqp.Dial(s.url)
	util.RequireNoError(err, "error: subscriber failed to establish a connection")
	ch, err := conn.Channel()
	util.RequireNoError(err, "error: subscriber failed to create a channel")
	s.conn = conn
	s.ch = ch
}

func (s *subscriber) declareExchanges() {
	err := pubsub.ExchangeDeclare(
		s.ch,
		"whatdoing",
		"topic",
		pubsub.EXCHANGE_TYPE_DURABLE,
		nil,
	)
	util.RequireNoError(err, "error: subscriber failed to declare an exchange")
}
