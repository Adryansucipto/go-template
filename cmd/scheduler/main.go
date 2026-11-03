package main

import (
	"context"
	"go-template/internal/controller"
	"go-template/pkg/config"
	pkgDi "go-template/pkg/di"
	"go-template/pkg/resource/model"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

func onRunScheduler(container *dig.Container) error {
	return container.Invoke(func(
		ctrl controller.Controller,
		cfg config.Config,
		db model.Database,
		kafkaRes model.Kafka,
		logger *zap.Logger,
	) error {
		defer func() {
			_ = logger.Sync()
		}()

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		logger.Info("starting standalone Cron Scheduler", zap.String("env", cfg.AppConfig.Env))

		// Start Cron Scheduler only (configured via configs/*.yaml)
		ctrl.Scheduler.Start(ctx)

		<-ctx.Done()
		logger.Info("shutting down standalone scheduler gracefully...")

		if err := kafkaRes.Close(); err != nil {
			logger.Error("failed closing kafka writer", zap.Error(err))
		}

		if sqlDB, err := db.PostgreSql.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Error("failed closing database connection", zap.Error(err))
			}
		}

		logger.Info("scheduler exited cleanly")
		return nil
	})
}

func main() {
	container, err := pkgDi.Container()
	if err != nil {
		panic(err)
	}
	if err := onRunScheduler(container); err != nil {
		panic(err)
	}
}

