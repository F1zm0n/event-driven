package transport

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/F1zm0n/event-driven/auth/internal/usecases"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(log *slog.Logger, service usecases.CustomerUsecases, port string) App {
	srv := grpc.NewServer()
	RegisterServer(srv, service)

	return App{
		logger:     log,
		gRPCServer: srv,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "auth#transport#grpc#prepare.Run"
	l := a.logger.With(slog.String("op", op), slog.String("port", a.port))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return err
	}
	l.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	const op = "auth#transport#grpc#prepare.Stop"
	l := a.logger.With(slog.String("op", op))

	l.Info("stopping gRPC server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
