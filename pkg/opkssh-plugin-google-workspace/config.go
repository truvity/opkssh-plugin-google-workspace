package opksshplugingoogleworkspace

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"log/slog"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Google ConfigGoogle `json:"google"          yaml:"google"`
		Policy Policy       `json:"policy"          yaml:"policy"`
		Cache  *ConfigCache `json:"cache,omitempty" yaml:"cache,omitempty"`
	}
)

func LoadConfig(
	ctx context.Context,
	logger *slog.Logger,
	pathConfig string,
	pathCache string,
	cacheDuration time.Duration,
) (*Config, error) {
	// get absolute path to config file
	abs, err := filepath.Abs(pathConfig)
	if err != nil {
		const message = "failed to get absolute path"
		logger.ErrorContext(ctx,
			message,
			slog.String("path", pathConfig),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s path %s %w",
			message,
			pathConfig,
			err,
		)
		return nil, err
	}
	pathConfig = abs

	// read config file
	logger.DebugContext(ctx, "read config file", slog.String("path", pathConfig))
	data, err := os.ReadFile(pathConfig)
	if err != nil {
		const message = "failed to read config file"
		logger.ErrorContext(ctx,
			message,
			slog.String("path", pathConfig),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s path %s %w",
			message,
			pathConfig,
			err,
		)
		return nil, err
	}

	// parse config file
	logger.DebugContext(ctx, "parse config file", slog.String("path", pathConfig))
	var result Config
	if err = yaml.Unmarshal(data, &result); err != nil {
		const message = "failed to parse config file"
		logger.ErrorContext(ctx,
			message,
			slog.String("path", pathConfig),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s path %s %w",
			message,
			pathConfig,
			err,
		)
		return nil, err
	}

	for principal := range result.Policy {
		policy := result.Policy[principal]
		sort.Strings(policy.User)
		sort.Strings(policy.Group)
	}

	logger.DebugContext(ctx, "load config file completed")

	// get absolute path to servire account key file
	serviceAccountKeyPath := result.Google.ServiceAccount.KeyFile
	if !filepath.IsAbs(serviceAccountKeyPath) {
		serviceAccountKeyPath = filepath.Join(
			filepath.Dir(pathConfig),
			serviceAccountKeyPath,
		)
	}

	// load service account key file
	logger.DebugContext(ctx, "read service account file",
		slog.String("path", serviceAccountKeyPath),
	)
	data, err = os.ReadFile(serviceAccountKeyPath)
	if err != nil {
		const message = "failed to read service account file"
		logger.ErrorContext(ctx, message,
			slog.String("path", serviceAccountKeyPath),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s path %s %w",
			message,
			pathConfig,
			err,
		)
		return nil, err
	}

	// parse service account key file
	logger.DebugContext(ctx, "parse service account file",
		slog.String("path", serviceAccountKeyPath),
	)
	result.Google.ServiceAccount.Key, err = google.JWTConfigFromJSON(data, admin.AdminDirectoryGroupMemberReadonlyScope)
	if err != nil {
		const message = "failed to parse Google Service Account key file"
		logger.ErrorContext(ctx, message,
			slog.String("path", serviceAccountKeyPath),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s path %s %w",
			message,
			pathConfig,
			err,
		)
		return nil, err
	}

	logger.DebugContext(ctx, "load service account file completed")

	// defaults
	if result.Cache == nil {
		result.Cache = &ConfigCache{}
	}
	if result.Cache.Path == nil {
		result.Cache.Path = &pathCache
	}
	if result.Cache.Duration == nil {
		result.Cache.Duration = &cacheDuration
	}

	return &result, nil
}
