package zapfx

import (
	"context"
	"fmt"

	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewZapParams defines the dependencies of the zapfx module.
type NewZapParams struct {
	fx.In

	Config    config.Provider
	Lifecycle fx.Lifecycle
}

// NewZapResult defines the objects that the zapfx module provides.
type NewZapResult struct {
	fx.Out

	Level  zap.AtomicLevel
	Logger *zap.Logger
}

// NewZap exports functionality similar to Module, but allows the caller to wrap
// or modify NewZapResult. Most users should use Module instead.
func NewZap(p NewZapParams) (NewZapResult, error) {
	var (
		c   = zap.NewProductionConfig()
		raw = p.Config.Get(ConfigurationKey)
	)
	if err := raw.Populate(&c); err != nil {
		return NewZapResult{}, fmt.Errorf("failed to load logging config: %v", err)
	}

	logger, err := c.Build()
	if err != nil {
		return NewZapResult{}, fmt.Errorf("failed to create zap logger: %v", err)
	}

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			logger.Sync()
			return nil
		},
	})

	return NewZapResult{
		Level:  c.Level,
		Logger: logger,
	}, nil
}
