package model

import (
	"context"
	"encoding/json"
	"fmt"
	"go-template/pkg/config"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
)

func init() {
	encoding.RegisterCodec(JSONCodec{})
}

// JSONCodec allows gRPC services to exchange Go structs seamlessly alongside standard Protobuf messages.
type JSONCodec struct{}

func (JSONCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (JSONCodec) Name() string {
	return "json"
}

type GRPCClient struct {
	Config config.OutboundClientConfig
	mu     sync.RWMutex
	conns  map[string]*grpc.ClientConn
}

func NewGRPCClientManager(cfg config.OutboundClientConfig) *GRPCClient {
	return &GRPCClient{
		Config: cfg,
		conns:  make(map[string]*grpc.ClientConn),
	}
}

// GetConnection returns a pooled *grpc.ClientConn for the named service in configs/*.yaml.
// It automatically attaches x-api-key in outgoing gRPC metadata if configured.
func (g *GRPCClient) GetConnection(serviceName string, extraOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	g.mu.RLock()
	if conn, ok := g.conns[serviceName]; ok {
		g.mu.RUnlock()
		return conn, nil
	}
	g.mu.RUnlock()

	g.mu.Lock()
	defer g.mu.Unlock()

	if conn, ok := g.conns[serviceName]; ok {
		return conn, nil
	}

	svcCfg, exists := g.Config.Services[serviceName]
	if !exists || strings.TrimSpace(svcCfg.GRPCAddr) == "" {
		return nil, fmt.Errorf("[GRPCClient] external service '%s' grpc_addr is not configured", serviceName)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(apiKeyUnaryClientInterceptor(svcCfg.APIKey)),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(JSONCodec{}.Name())),
	}
	opts = append(opts, extraOpts...)

	conn, err := grpc.NewClient(svcCfg.GRPCAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("[GRPCClient] failed connecting to %s (%s): %w", serviceName, svcCfg.GRPCAddr, err)
	}

	g.conns[serviceName] = conn
	return conn, nil
}

func (g *GRPCClient) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	var firstErr error
	for name, conn := range g.conns {
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed closing grpc conn %s: %w", name, err)
		}
		delete(g.conns, name)
	}
	return firstErr
}

func apiKeyUnaryClientInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if apiKey != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", apiKey)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
