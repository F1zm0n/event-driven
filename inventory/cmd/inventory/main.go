package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/F1zm0n/event-driven/inventory/internal/repository"
	transport "github.com/F1zm0n/event-driven/inventory/internal/transport/grpc"
	"github.com/F1zm0n/event-driven/inventory/internal/usecases"
)

func main() {
	err := godotenv.Load("/home/F1zm0/Files/projects/gov1/event-driven/inventory/.env")
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := sqlx.MustOpen("postgres", os.Getenv("DSN"))

	var (
		repo    = repository.NewInventoryRepository(logger, db)
		service = usecases.NewInventoryUsecases(repo)
		app     = transport.New(logger, service, os.Getenv("PORT"))
	)
	defer db.Close()

	go app.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	logger.Info("stopping application", slog.String("signal", sign.String()))

	app.Stop()

	logger.Info("application stopped")
}
