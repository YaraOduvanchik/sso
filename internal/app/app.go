package app

import (
	"log/slog"
	"sso/internal/app/grpcapp"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	logger *slog.Logger,
	port int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировать хранилище

	// TODO: инициализировать auth-сервер

	grpcApp := grpcapp.New(logger, port)

	return &App{
		GRPCServer: grpcApp,
	}
}
