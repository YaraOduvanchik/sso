﻿package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"sso/internal/grpc/authgrpc"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	logger := a.logger.With(slog.String("op", op), slog.Int("port", a.port))

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("starting gRPC server", slog.String("address", listen.Addr().String()))

	if err := a.gRPCServer.Serve(listen); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.logger.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.String("address", fmt.Sprintf(":%d", a.port)))

	a.gRPCServer.GracefulStop()
}
