package conncheck

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type RegistryConfig struct {
	Endpoint string
	Username string
	Password string
}

type RegistryProbe struct {
	config RegistryConfig
	client *http.Client
}

func NewRegistryProbe(cfg RegistryConfig) *RegistryProbe {
	return &RegistryProbe{
		config: cfg,
		client: &http.Client{},
	}
}

func (p *RegistryProbe) Test(ctx context.Context) error {
	if err := p.checkV2(ctx); err != nil {
		logrus.Error("registry v2 check failed: ", err)
		return fmt.Errorf("registry v2 check failed: %w", err)
	}
	logrus.Info("registry v2 check passed")

	if err := p.checkCatalog(ctx); err != nil {
		logrus.Error("registry catalog check failed: ", err)
		return fmt.Errorf("registry catalog check failed: %w", err)
	}
	logrus.Info("registry catalog check passed")

	return nil
}

func (p *RegistryProbe) checkV2(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v2/", p.config.Endpoint), nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	if p.config.Username != "" {
		req.SetBasicAuth(p.config.Username, p.config.Password)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized: invalid credentials")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (p *RegistryProbe) checkCatalog(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v2/_catalog", p.config.Endpoint), nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	if p.config.Username != "" {
		req.SetBasicAuth(p.config.Username, p.config.Password)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
