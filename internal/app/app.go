package app

import (
	"log/slog"
	"time"
	grpcapp "usekit-auth/internal/app/grpc"
)

type App struct {
	GrpcServer *grpcapp.AppGrpc
}

func New(
	logger *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировать хранилище (storage)

	// TODO: инициализировать сервисный слой auth

	grpcApp := grpcapp.New(logger, grpcPort)

	return &App{
		GrpcServer: grpcApp,
	}
}
