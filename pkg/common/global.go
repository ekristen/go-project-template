package common

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Flags() []cli.Flag {
	globalFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log Level",
			Aliases: []string{"l"},
			EnvVars: []string{"LOGLEVEL"},
			Value:   "info",
		},
	}

	return globalFlags
}

func Before(c *cli.Context) error {
	config := zap.NewProductionConfig()

	// Handle color settings
	if c.Bool("log-disable-color") {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Handle timestamp settings
	if c.Bool("log-full-timestamp") {
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Handle caller settings
	config.DisableCaller = !c.Bool("log-caller")

	// Set log level
	switch c.String("log-level") {
	case "trace", "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	zap.ReplaceGlobals(logger)
	return nil
}
