package injection

import (
	pkgConfig "go-template/pkg/config"
	"go-template/pkg/resource/model"
)

func NewGRPCClient(cfg pkgConfig.Config) (*model.GRPCClient, error) {
	return model.NewGRPCClientManager(cfg.Client), nil
}
