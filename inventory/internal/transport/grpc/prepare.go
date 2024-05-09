package transport

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/F1zm0n/event-driven/inventory/internal/usecases"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(log *slog.Logger, service usecases.InventoryUsecases, port string) App {
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
	const op = "inventory#transport#grpc#prepare.Run"
	l := a.logger.With(slog.String("op", op))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		l.Info("error opening tcp connection", slog.String("error", err.Error()))
		return err
	}

	l.Info("starting gRPC server", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		l.Error("error starting gRPC server", slog.String("error", err.Error()))
		return err
	}

	l.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	return nil
}

func (a *App) Stop() {
	const op = "inventory#transport#grpc#prepare.Run"
	l := a.logger.With("op", op)
	l.Info("stopping grpc server")

	a.gRPCServer.GracefulStop()

	l.Info("stopped grpc server")
}
