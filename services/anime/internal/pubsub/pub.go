package pubsub

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/JustinLi007/whatdoing/libs/go/config"
	"github.com/JustinLi007/whatdoing/libs/go/pubsub"
	"github.com/JustinLi007/whatdoing/libs/go/util"
	"github.com/JustinLi007/whatdoing/services/anime/internal/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DURABLE int = iota
	TRANSIENT
)

type Publisher interface {
	Start(ctx context.Context)
}

type publisher struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	url           string
	interval      time.Duration
	outboxService database.ServiceOutbox
}

var publisherIntance *publisher

func NewPublisher(c *config.Config) Publisher {
	connStr := c.Get("DB_URL")
	if connStr == "" {
		util.RequireNoError(errors.New("invalid db url"), "error")
	}

	db, err := database.NewDb(connStr)
	util.RequireNoError(err, "error: publisher failed to connect to database")

	outboxService := database.NewServiceOutbox(db)
	return newPublisher(outboxService)
}

func newPublisher(outboxService database.ServiceOutbox) Publisher {
	if publisherIntance != nil {
		return publisherIntance
	}
	newPublisher := &publisher{
		conn:          nil,
		ch:            nil,
		url:           "amqp://guest:guest@whatdoing-msg-broker:5672/",
		interval:      time.Second * 5,
		outboxService: outboxService,
	}
	publisherIntance = newPublisher
	publisherIntance.connect()
	publisherIntance.declareExchanges()

	return publisherIntance
}

func (p *publisher) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go p.start(ctx, wg)
	wg.Wait()
}

func (p *publisher) start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(p.interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			pubsub.PublishJSON(
				p.ch,
				"whatdoing",
				"anime.test",
				[]byte("test"),
			)
			timer.Reset(p.interval)
		case <-ctx.Done():
			return
		}
	}
}

func (p *publisher) connect() {
	conn, err := amqp.Dial(p.url)
	util.RequireNoError(err, "error: publisher failed to establish a connection")
	p.conn = conn

	ch, err := conn.Channel()
	util.RequireNoError(err, "error: publisher failed to create a channel")
	p.ch = ch
}

func (p *publisher) declareExchanges() {
	err := pubsub.ExchangeDeclare(
		p.ch,
		"whatdoing",
		"topic",
		pubsub.EXCHANGE_TYPE_DURABLE,
		nil,
	)
	util.RequireNoError(err, "error: publisher failed to declare an exchange")
}
