package exchange

import (
	"example.com/binance/internal/config"
	"github.com/adshao/go-binance/v2"
)

// Service provides access to Binance REST API operations
type Service struct {
	client *binance.Client
	cfg    *config.Config
}

// NewService creates a new Binance exchange service using the provided configuration
func NewService(cfg *config.Config) *Service {
	if cfg == nil {
		cfg = config.Load()
	}
	return &Service{
		client: binance.NewClient(cfg.APIKey, cfg.SecretKey),
		cfg:    cfg,
	}
}

// DefaultService returns a Service instance initialized from the environment
func DefaultService() *Service {
	return NewService(config.Load())
}
