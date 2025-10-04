package main

import (
	"context"
	"os"
	"time"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/ekristen/go-project-template/pkg/common"

	_ "github.com/ekristen/go-project-template/pkg/commands/example"
	_ "github.com/ekristen/go-project-template/pkg/commands/server"
)

func main() {
	var exitCode int

	func() {
		defer func() {
			// Shutdown telemetry on exit
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := common.Shutdown(shutdownCtx); err != nil {
				log.Warn().Err(err).Msg("failed to shutdown telemetry")
			}

			if r := recover(); r != nil {
				// log panics using zerolog and set exit code
				log.Error().Interface("panic", r).Msg("panic recovered")
				exitCode = 1
			}
		}()

		app := &cli.Command{
			Name:    common.AppVersion.Name,
			Usage:   common.AppVersion.Name,
			Version: common.AppVersion.Summary,
			Authors: []any{
				"Erik Kristensen <erik@erikkristensen.com>",
			},
			Commands: common.GetCommands(),
			CommandNotFound: func(ctx context.Context, command *cli.Command, s string) {
				log.Error().Str("command", s).Msg("command not found")
			},
			EnableShellCompletion: true,
			Before:                common.Before,
			Flags:                 common.Flags(),
		}

		ctx := signals.SetupSignalContext()
		if err := app.Run(ctx, os.Args); err != nil {
			log.Error().Err(err).Msg("fatal error")
			exitCode = 1
		}
	}()

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
