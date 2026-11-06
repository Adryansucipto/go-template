package grpc

import (
	"context"
	"go-template/internal/controller/grpc/user"
	"go-template/pkg/config"
	userv1 "go-template/proto/user/v1"
	"strings"
	"time"

	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Controller struct {
	dig.In

	UserServer user.Server
	Config     config.Config
	Logger     *zap.Logger
}

func (c *Controller) NewServer() *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			c.loggingUnaryInterceptor(),
			c.apiKeyUnaryInterceptor(),
		),
	)

	userv1.RegisterUserServiceServer(srv, &c.UserServer)
	return srv
}

func (c *Controller) loggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		c.Logger.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("latency", time.Since(start)),
			zap.Error(err),
		)
		return resp, err
	}
}

func (c *Controller) apiKeyUnaryInterceptor() grpc.UnaryServerInterceptor {
	headerKey := strings.ToLower(c.Config.APIKey.HeaderName())

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if !c.Config.APIKey.Enabled {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get(headerKey)
		if len(values) == 0 {
			values = md.Get("x-api-key")
		}

		var incomingKey string
		if len(values) > 0 {
			incomingKey = values[0]
		}

		if _, valid := c.Config.APIKey.Validate(incomingKey); !valid {
			return nil, status.Error(codes.Unauthenticated, "invalid or missing x-api-key in metadata")
		}

		return handler(ctx, req)
	}
}
