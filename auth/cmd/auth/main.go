package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/F1zm0n/event-driven/auth/internal/repository"
	transport "github.com/F1zm0n/event-driven/auth/internal/transport/grpc"
	"github.com/F1zm0n/event-driven/auth/internal/transport/rabbitmq"
	"github.com/F1zm0n/event-driven/auth/internal/usecases"
)

func main() {
	tokenTTL := 24 * time.Hour

	err := godotenv.Load("/home/F1zm0/Files/projects/gov1/event-driven/auth/.env")
	if err != nil {
		log.Fatal(err)
	}
	var (
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

		db   = sqlx.MustOpen("postgres", os.Getenv("DSN"))
		repo = repository.NewCustomerRepository(logger, db)

		usecase = usecases.NewCustomerUsecases(repo, logger, tokenTTL)
		grpcSrv = transport.New(logger, usecase, os.Getenv("PORT"))
	)

	conn, err := amqp.Dial(os.Getenv("RABBIT_DSN"))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	queue := rabbitmq.MustMakeQueue(ch)
	broker := rabbitmq.NewRabbitCustomer(queue, ch, usecase, logger)

	rabbitmq.MustMakeExchangeBindings(ch, queue)

	go broker.ConsumeQueue(context.TODO())

	defer db.Close()
	defer conn.Close()

	go grpcSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	logger.Info("stopping application", slog.String("signal", sign.String()))

	grpcSrv.Stop()

	logger.Info("application stopped")
}
