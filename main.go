package main

import (
	"context"
	"os"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"

	"github.com/ekristen/go-project-template/pkg/common"

	_ "github.com/ekristen/go-project-template/pkg/commands/example"
	_ "github.com/ekristen/go-project-template/pkg/commands/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// log panics using zap and force exit
			zap.L().Error("panic recovered", zap.Any("panic", r))
			os.Exit(1)
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
			zap.L().Error("command not found", zap.String("command", s))
		},
		EnableShellCompletion: true,
		Before:                common.Before,
		Flags:                 common.Flags(),
	}

	ctx := signals.SetupSignalContext()
	if err := app.Run(ctx, os.Args); err != nil {
		zap.L().Fatal("fatal error", zap.Error(err))
	}
}
