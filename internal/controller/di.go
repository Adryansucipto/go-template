package controller

import (
	"go-template/internal/controller/event"
	"go-template/internal/controller/grpc"
	"go-template/internal/controller/http"
	"go-template/internal/controller/scheduler"

	"go.uber.org/dig"
)

type Controller struct {
	dig.In
	HTTP      http.Controller
	GRPC      grpc.Controller
	Event     event.Controller
	Scheduler scheduler.Controller
}
