package main

import (
	"context"
	"errors"
	"fmt"
	"go-template/internal/controller"
	"go-template/pkg/config"
	pkgDi "go-template/pkg/di"
	"go-template/pkg/resource/model"
	"go-template/pkg/validator"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type (
	ContainerCall func() (*dig.Container, error)
	Invoke        func(container *dig.Container) error
	InvokeError   func(container *dig.Container, err error)
)

func Run(containerCall ContainerCall, invoke Invoke, onError InvokeError) error {
	container, err := containerCall()
	if err != nil {
		return err
	}

	if err := invoke(container); err != nil {
		onError(container, err)
	}
	return nil
}

func onRun(container *dig.Container) error {
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

		e := echo.New()
		e.HideBanner = true
		e.Validator = validator.New()

		ctrl.HTTP.Routes(e)

		port := os.Getenv("PORT")
		if port == "" {
			port = cfg.AppConfig.Port
		}
		if port == "" {
			port = "9001"
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			logger.Info("starting HTTP API server", zap.String("port", port), zap.String("env", cfg.AppConfig.Env))
			if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal("shutting down the server", zap.Error(err))
			}
		}()

		<-ctx.Done()
		logger.Info("shutting down HTTP API server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(shutdownCtx); err != nil {
			logger.Error("server forced to shutdown", zap.Error(err))
		}

		if err := kafkaRes.Close(); err != nil {
			logger.Error("failed closing kafka writer", zap.Error(err))
		}

		if sqlDB, err := db.PostgreSql.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Error("failed closing database connection", zap.Error(err))
			}
		}

		logger.Info("HTTP API server exited cleanly")
		return nil
	})
}

func OnError(container *dig.Container, err error) {
	if container != nil && err != nil {
		panic(err)
	}
}

func main() {
	if err := Run(pkgDi.Container, onRun, OnError); err != nil {
		panic(err)
	}
}
