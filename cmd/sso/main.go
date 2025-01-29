package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/logger"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := logger.Setup(cfg.Env)

	log.Info("starting app", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.Storage, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	// TODO: инициализировать приложение (App)

	// TODO: запустить gRPC-сервер приложения

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping app", slog.String("signal", sign.String()))

	application.GRPCServer.Stop()

	log.Info("stopped app")
}
