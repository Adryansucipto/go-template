package resource

import (
	pkgConfig "go-template/pkg/config"

	"go.uber.org/dig"
)

type Resource struct {
	dig.In

	Config pkgConfig.Config
}
