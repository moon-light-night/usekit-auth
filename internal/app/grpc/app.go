package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	authgrpc "usekit-auth/internal/grpc/auth"
)

type AppGrpc struct {
	logger     *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int) *AppGrpc {
	// создается grpc сервер
	grpcServer := grpc.NewServer()

	// регистрируется grpc сервер
	authgrpc.Register(grpcServer)

	return &AppGrpc{
		logger:     logger,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (app *AppGrpc) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *AppGrpc) Run() error {
	const op = "grpcapp.Run"

	log := app.logger.With(slog.String("op", op), slog.Int("port", app.port))

	// создается слушатель для tcp соединения, которое необходимо для работы grpc
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", listener.Addr().String()))

	// запускается сервер и указывается listener для обработки запросов
	if err := app.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (app *AppGrpc) Stop() {
	const op = "grpcapp.Stop"

	app.logger.With(slog.String("op", op)).Info("grpc server is stopping", slog.Int("port", app.port))

	// прекращается прием новых запросов, блокируется выполнение кода до момента обработки уже выполняемых запросов
	app.grpcServer.GracefulStop()
}
