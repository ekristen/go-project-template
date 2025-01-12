package main

import (
	"os"
	"path"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/swade1987/go-project-template/pkg/common"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	_ "github.com/swade1987/go-project-template/pkg/commands/example"
)

func initializeLogger() *zap.Logger {
	// Create a default production logger until the Before function configures it
	logger, _ := zap.NewProduction()
	return logger
}

// syncLogger flushes the logger buffer
func syncLogger(logger *zap.Logger) {
	_ = logger.Sync() // ignore sync errors
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			// For zap, we'll check if it's a logger error
			if _, ok := r.(error); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	// Initialize default logger
	logger := initializeLogger()
	defer syncLogger(logger)

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = common.AppVersion.Name
	app.Version = common.AppVersion.Summary
	app.Authors = []*cli.Author{
		{
			Name:  "Steve Wade",
			Email: "steven@stevenwade.co.uk",
		},
	}

	app.Before = common.Before
	app.Flags = common.Flags()

	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logger.Fatal("command not found",
			zap.String("command", command),
		)
	}

	ctx := signals.SetupSignalContext()
	if err := app.RunContext(ctx, os.Args); err != nil {
		logger.Fatal("application error",
			zap.Error(err),
		)
	}
}
