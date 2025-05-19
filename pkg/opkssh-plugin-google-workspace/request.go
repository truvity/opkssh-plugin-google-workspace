package opksshplugingoogleworkspace

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
)

type (
	Request struct {
		Principal     string `env:"OPKSSH_PLUGIN_U"`
		Email         string `env:"OPKSSH_PLUGIN_EMAIL"`
		EmailVerified bool   `env:"OPKSSH_PLUGIN_EMAIL_VERIFIED"`
		ClientID      string `env:"OPKSSH_PLUGIN_AUD"`
	}
)

func LoadRequest(ctx context.Context, logger *slog.Logger, environ []string) (*Request, error) {
	if environ == nil {
		environ = os.Environ()
	}

	logger.DebugContext(ctx, "load plugin")

	var plugin Request
	err := env.ParseWithOptions(&plugin, env.Options{
		Environment: env.ToMap(environ),
	})
	if err != nil {
		const message = "failed to parse environment variables"
		logger.ErrorContext(ctx, message,
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s %w", message, err)
		return nil, err
	}

	logger.DebugContext(ctx, "request loaded")

	return &plugin, nil
}
