package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"usekit-auth/internal/app"
	"usekit-auth/internal/config"
)

const (
	envDev  = "development"
	envProd = "production"
)

// главная функция, которая собирает и запускает приложение, по сути, точка входа
func main() {
	// TODO: инициализировать объект конфига
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: инициализировать логгер
	logger := setupLogger(cfg.Env)
	fmt.Println(logger)
	logger.Info("starting application", slog.String("env", cfg.Env))

	// TODO: инициализировать приложение (app)
	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	// TODO: запустить grpc-сервер приложения
	go application.GrpcServer.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	// слушаются события ОС, при вызове которых в канал stop запишется ...
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// читаем из канала stop(блокирующая операция, пока выполняется горутина MustRun)
	// и пишем сигнал в переменную sign, который пришел от ОС в канал stop
	sign := <-stop
	logger.Info("stopping", slog.String("signal", sign.String()))

	// завершаем работу приложения(работающие в это время процессы выполнятся до конца)
	application.GrpcServer.Stop()
	logger.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
