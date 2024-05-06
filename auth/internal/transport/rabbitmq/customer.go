package rabbitmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/F1zm0n/event-driven/auth/internal/dto"
	"github.com/F1zm0n/event-driven/auth/internal/usecases"
)

type RabbitCustomer struct {
	q       amqp.Queue
	ch      *amqp.Channel
	usecase usecases.CustomerUsecases
	forever chan struct{}
	logger  *slog.Logger
}

func NewRabbitCustomer(
	q amqp.Queue,
	ch *amqp.Channel,
	usecase usecases.CustomerUsecases,
	logger *slog.Logger,
) RabbitCustomer {
	return RabbitCustomer{
		q:       q,
		ch:      ch,
		usecase: usecase,
		forever: make(chan struct{}),
		logger:  logger,
	}
}

func (r RabbitCustomer) ConsumeQueue(ctx context.Context) {
	msgs, err := r.ch.ConsumeWithContext(
		ctx,
		r.q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range msgs {
			r.logger.Info(
				"receiver message from queue",
			)
			if v, ok := msg.Headers["operation"]; ok && v.(string) == os.Getenv("REGISTER_HEADER") {
				r.logger.Info(
					"receiver message from queue",
					slog.String("header", os.Getenv("REGISTER_HEADER")),
				)
				r.register(ctx, msg.Body)
			}
		}
	}()
	r.logger.Info("waiting for queue messages")
	<-r.forever
}

func (r RabbitCustomer) register(ctx context.Context, body []byte) {
	var customerDto dto.CustomerDto
	err := json.Unmarshal(body, &customerDto)
	if err != nil {
		return
	}

	if err = r.usecase.Register(ctx, customerDto); err != nil {
		return
	}
}
