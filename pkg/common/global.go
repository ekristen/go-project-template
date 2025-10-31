package common

import (
	"context"
	"os"

	zerologhook "github.com/ekristen/go-telemetry/hooks/zerolog/v2"
	"github.com/ekristen/go-telemetry/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var telemetryInstance *telemetry.Telemetry

func Flags() []cli.Flag {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log Level",
			Aliases: []string{"l"},
			Sources: cli.EnvVars("LOG_LEVEL"),
			Value:   "info",
		},
		&cli.BoolFlag{
			Name:    "log-caller",
			Usage:   "log the caller (aka line number and file)",
			Sources: cli.EnvVars("LOG_CALLER"),
			Value:   true,
		},
		&cli.StringFlag{
			Name:    "log-format",
			Usage:   "the log format to use, defaults to auto, options are auto, json, console",
			Sources: cli.EnvVars("LOG_FORMAT"),
			Value:   "auto",
		},
	}

	return globalFlags
}

func Before(ctx context.Context, c *cli.Command) (context.Context, error) {
	// Parse log level
	logLevel := c.String("log-level")
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return ctx, err
	}

	// Configure global zerolog level
	zerolog.SetGlobalLevel(level)

	// Set up console writer based on format preference and terminal detection
	if c.String("log-format") == "json" || (!term.IsTerminal(int(os.Stdout.Fd())) && c.String("log-format") == "auto") {
		// Use JSON format for non-TTY or when explicitly requested
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		// Use console format with colors for TTY
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02T15:04:05Z07:00",
		}
		log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	}

	// Configure caller information
	if c.Bool("log-caller") {
		log.Logger = log.Logger.With().Caller().Logger()
	}

	// Initialize telemetry (this will set up OTEL providers)
	opts := &telemetry.Options{
		ServiceName:    AppVersion.Name,
		ServiceVersion: AppVersion.Summary,
		BatchExport:    true, // False by default, true batches for production
	}

	telemetryInstance, err = telemetry.New(ctx, opts)
	if err != nil {
		log.Warn().Err(err).Msg("failed to initialize telemetry")
	}

	// Attach OTel hook to the logger if telemetry was initialized successfully
	if telemetryInstance != nil && telemetryInstance.LoggerProvider() != nil {
		hook := zerologhook.New(
			telemetryInstance.ServiceName(),
			telemetryInstance.ServiceVersion(),
			telemetryInstance.LoggerProvider(),
		)
		if hook != nil {
			log.Logger = log.Logger.Hook(hook)
		}
	}

	return ctx, nil
}

// GetTelemetry returns the global telemetry instance
func GetTelemetry() *telemetry.Telemetry {
	return telemetryInstance
}

// Shutdown gracefully shuts down the telemetry instance
func Shutdown(ctx context.Context) error {
	if telemetryInstance != nil {
		return telemetryInstance.Shutdown(ctx)
	}
	return nil
}
