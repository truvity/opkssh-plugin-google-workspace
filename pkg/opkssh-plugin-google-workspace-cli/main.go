package opksshplugingoogleworkspacecli

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	opksshplugingoogleworkspace "github.com/truvity/opkssh-plugin-google-workspace/pkg/opkssh-plugin-google-workspace"
	"github.com/urfave/cli/v3"
)

const (
	FlagConfig     = "config"
	FlagCache      = "cache"
	FlagLog        = "log"
	FlagExpiration = "expiration"
	FlagVerbose    = "verbose"
	FlagQuiet      = "quiet"
)

func Main() {
	if (func() error {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()
		var logCloser io.Closer
		defer func() {
			if logCloser != nil {
				_ = logCloser.Close()
			}
		}()
		var logger *slog.Logger
		app := cli.Command{
			Name:        "opkssh-plugin-google-workspace",
			Usage:       "provide necessary environment variables (usually done by opkssh)",
			Description: "Plugin for opkssh to enable group-based authorization for Google Workspace",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        FlagConfig,
					Usage:       "path to config",
					DefaultText: opksshplugingoogleworkspace.DefaultConfigPath,
					Value:       opksshplugingoogleworkspace.DefaultConfigPath,
				},
				&cli.StringFlag{
					Name:        FlagCache,
					Usage:       "path to cache",
					DefaultText: opksshplugingoogleworkspace.DefaultCachePath,
					Value:       opksshplugingoogleworkspace.DefaultCachePath,
				},
				&cli.StringFlag{
					Name:        FlagLog,
					Usage:       "path to log",
					DefaultText: opksshplugingoogleworkspace.DefaultLogPath,
					Value:       opksshplugingoogleworkspace.DefaultLogPath,
				},
				&cli.DurationFlag{
					Name:        FlagExpiration,
					Usage:       "cache expiration",
					DefaultText: opksshplugingoogleworkspace.DefaultCacheDuration.String(),
					Value:       opksshplugingoogleworkspace.DefaultCacheDuration,
				},
				&cli.BoolFlag{
					Name:        FlagVerbose,
					Aliases:     []string{"v"},
					Usage:       "verbose logging (debug)",
					DefaultText: "false",
					Value:       false,
				},
				&cli.BoolFlag{
					Name:        FlagQuiet,
					Aliases:     []string{"q"},
					Usage:       "quiet logging (errors only)",
					DefaultText: "false",
					Value:       false,
				},
			},
			Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
				verbose := c.Bool(FlagVerbose)
				quiet := c.Bool(FlagQuiet)
				if verbose && quiet {
					return ctx, fmt.Errorf("flag %s and %s are mutually exclusive", FlagVerbose, FlagQuiet)
				}
				var level slog.Level
				if verbose {
					level = slog.LevelDebug
				}
				if quiet {
					level = slog.LevelError
				}
				logPath := c.String(FlagLog)
				var logWriter io.Writer
				switch logPath {
				case "stdout":
					logWriter = os.Stdout
				case "stderr":
					logWriter = os.Stderr
				default:
					logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
					if err != nil {
						return ctx, err
					}
					logWriter = logFile
					logCloser = logFile
				}
				logger = slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{
					Level: level,
				}))
				return ctx, nil
			},
			Action: func(ctx context.Context, c *cli.Command) error {
				if logger == nil {
					panic(logger)
				}
				config, err := opksshplugingoogleworkspace.LoadConfig(ctx, logger,
					c.String(FlagConfig),
					c.String(FlagCache),
					c.Duration(FlagExpiration),
				)
				if err != nil {
					return err
				}

				fetcher := opksshplugingoogleworkspace.NewGooglFetcher(config.Google.ServiceAccount)
				cache := opksshplugingoogleworkspace.NewCacheFetcher(*config.Cache, config.Google.Workspace.CustomerID, fetcher)

				request, err := opksshplugingoogleworkspace.LoadRequest(ctx, logger, nil)
				if err != nil {
					return err
				}
				allow, err := opksshplugingoogleworkspace.Verify(ctx, logger, cache, config, request)
				if err != nil {
					return err
				}
				if allow {
					fmt.Println("allow")
				}
				return nil
			},
		}
		if err := app.Run(ctx, os.Args); err != nil {
			return err
		}
		return nil
	})() != nil {
		os.Exit(1)
	}
}
