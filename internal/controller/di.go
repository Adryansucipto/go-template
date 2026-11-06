package controller

import (
	"go-template/internal/controller/grpc"
	"go-template/internal/controller/http"

	"go.uber.org/dig"
)

// Controller strictly handles Inbound API delivery (REST HTTP & gRPC Server).
type Controller struct {
	dig.In
	HTTP http.Controller
	GRPC grpc.Controller
}
