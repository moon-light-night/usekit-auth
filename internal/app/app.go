package app

// main app

import (
	"log/slog"
	"time"
	grpcapp "usekit-auth/internal/app/grpc"
	"usekit-auth/internal/services/auth"
	"usekit-auth/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	// TODO: инициализировать сервисный слой auth
	authService := auth.New(logger, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(logger, authService, grpcPort)

	return &App{
		GrpcServer: grpcApp,
	}
}
