package common

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"golang.org/x/term"
)

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
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Adds color
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // Readable time format
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if c.String("log-format") == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // JSON encoder for non-TTY
	} else if term.IsTerminal(int(os.Stdout.Fd())) {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // Console encoder with colors
	}

	logLevel := c.String("log-level")
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return ctx, err
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.Lock(os.Stdout)), // Output to stdout
		level,                                    // Log level
	)

	options := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel)}
	if c.Bool("log-caller") {
		options = append(options, zap.AddCaller())
	}

	logger := zap.New(core, options...)
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)

	return ctx, nil
}
