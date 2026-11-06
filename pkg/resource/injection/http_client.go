package injection

import (
	pkgConfig "go-template/pkg/config"
	"go-template/pkg/resource/model"
	"net/http"
	"time"
)

func NewHTTPClient(cfg pkgConfig.Config) (model.HTTPClient, error) {
	timeoutSec := cfg.Client.TimeoutSeconds
	if timeoutSec <= 0 {
		timeoutSec = 10
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Duration(timeoutSec) * time.Second,
		Transport: transport,
	}

	return model.HTTPClient{
		DefaultClient: client,
		Config:        cfg.Client,
	}, nil
}
