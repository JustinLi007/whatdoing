package pubsub

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/JustinLi007/whatdoing/libs/go/configs"
	"github.com/JustinLi007/whatdoing/libs/go/utils"
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

func NewPublisher(c *configs.Config) Publisher {
	connStr := c.Get("DB_URL")
	if connStr == "" {
		utils.RequireNoError(errors.New("invalid db url"), "error")
	}

	db, err := database.NewDb(connStr)
	utils.RequireNoError(err, "error: publisher failed to connect to database")

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
			p.publishJSON("whatdoing", "anime.create", []byte("test"))
			timer.Reset(p.interval)
		case <-ctx.Done():
			return
		}
	}
}

func (p *publisher) publishJSON(exchange, key string, msg []byte) error {
	return p.ch.PublishWithContext(
		context.Background(), // ctx
		exchange,             // exchange
		key,                  // key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)
}

func (p *publisher) connect() {
	conn, err := amqp.Dial(p.url)
	utils.RequireNoError(err, "Error: publisher failed to establish a connection")
	p.conn = conn

	ch, err := conn.Channel()
	utils.RequireNoError(err, "Error: publisher failed to create a channel")
	p.ch = ch
}

func (p *publisher) declareExchanges() {
	err := p.declareExchange("whatdoing", "topic", true, nil)
	utils.RequireNoError(err, "Error: publisher failed to declare an exchange")
}

func (p *publisher) declareExchange(name, kind string, durable bool, args amqp.Table) error {
	return p.ch.ExchangeDeclare(
		name,     // name
		kind,     // kind
		durable,  // durable
		!durable, // auto delete
		false,    // internal
		false,    // no wait
		args,     // args
	)
}
