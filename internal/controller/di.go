package controller

import (
	"go-template/internal/controller/http"

	"go.uber.org/dig"
)

type Controller struct {
	dig.In
	HTTP http.Controller
}
