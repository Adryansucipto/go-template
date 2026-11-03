package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-template/pkg/config"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	DefaultClient *http.Client
	Config        config.OutboundClientConfig
}

type HTTPRequestOption struct {
	ServiceName string
	Method      string
	Path        string
	Headers     map[string]string
	Body        any
}

func (h *HTTPClient) DoJSON(ctx context.Context, opt HTTPRequestOption, responseTarget any) (int, error) {
	svcCfg, exists := h.Config.Services[opt.ServiceName]
	if !exists || strings.TrimSpace(svcCfg.BaseURL) == "" {
		return 0, fmt.Errorf("[HTTPClient] external service '%s' base_url is not configured", opt.ServiceName)
	}

	fullURL := strings.TrimRight(svcCfg.BaseURL, "/") + "/" + strings.TrimLeft(opt.Path, "/")

	var bodyReader io.Reader
	if opt.Body != nil {
		payloadBytes, err := json.Marshal(opt.Body)
		if err != nil {
			return 0, fmt.Errorf("[HTTPClient] failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(payloadBytes)
	}

	reqCtx := ctx
	var cancel context.CancelFunc
	if svcCfg.TimeoutSec > 0 {
		reqCtx, cancel = context.WithTimeout(ctx, time.Duration(svcCfg.TimeoutSec)*time.Second)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(reqCtx, opt.Method, fullURL, bodyReader)
	if err != nil {
		return 0, fmt.Errorf("[HTTPClient] failed creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Automatically attach target service's X-API-Key if configured in YAML
	if svcCfg.APIKey != "" {
		req.Header.Set("X-API-Key", svcCfg.APIKey)
	}

	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	resp, err := h.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("[HTTPClient] request to %s failed: %w", fullURL, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("[HTTPClient] failed reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("[HTTPClient] upstream %s returned status %d: %s", opt.ServiceName, resp.StatusCode, string(respBytes))
	}

	if responseTarget != nil && len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, responseTarget); err != nil {
			return resp.StatusCode, fmt.Errorf("[HTTPClient] failed unmarshaling response JSON: %w", err)
		}
	}

	return resp.StatusCode, nil
}

