package main

import (
	"context"
	"fmt"
	"go-template/internal/controller"
	"go-template/pkg/config"
	pkgDi "go-template/pkg/di"
	"go-template/pkg/resource/model"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func onRunGRPC(container *dig.Container) error {
	return container.Invoke(func(
		ctrl controller.Controller,
		cfg config.Config,
		db model.Database,
		kafkaRes model.Kafka,
		grpcClient *model.GRPCClient,
		logger *zap.Logger,
	) error {
		defer func() {
			_ = logger.Sync()
		}()

		port := os.Getenv("GRPC_PORT")
		if port == "" {
			port = cfg.GRPC.Port
		}
		if port == "" {
			port = "9002"
		}

		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			return fmt.Errorf("failed to listen on gRPC port %s: %w", port, err)
		}

		grpcServer := ctrl.GRPC.NewServer()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			logger.Info("starting standalone gRPC Server", zap.String("port", port), zap.String("env", cfg.AppConfig.Env))
			if err := grpcServer.Serve(lis); err != nil {
				logger.Fatal("gRPC server stopped unexpectedly", zap.Error(err))
			}
		}()

		<-ctx.Done()
		logger.Info("shutting down standalone gRPC Server gracefully...")

		grpcServer.GracefulStop()

		if err := grpcClient.Close(); err != nil {
			logger.Error("failed closing outbound gRPC client connections", zap.Error(err))
		}

		if err := kafkaRes.Close(); err != nil {
			logger.Error("failed closing kafka writer", zap.Error(err))
		}

		if sqlDB, err := db.PostgreSql.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Error("failed closing database connection", zap.Error(err))
			}
		}

		logger.Info("gRPC Server exited cleanly")
		return nil
	})
}

func main() {
	container, err := pkgDi.Container()
	if err != nil {
		panic(err)
	}
	if err := onRunGRPC(container); err != nil {
		panic(err)
	}
}

