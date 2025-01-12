package main

import (
	"os"
	"path"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/urfave/cli/v2"
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

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = common.AppVersion.Name
	app.Version = common.AppVersion.Summary
	app.Authors = []*cli.Author{
		{
			Name:  "Erik Kristensen",
			Email: "erik@erikkristensen.com",
		},
	}

	app.Before = common.Before
	app.Flags = common.Flags()

	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		zap.L().Fatal("command not found.", zap.String("command", command))
	}

	ctx := signals.SetupSignalContext()
	if err := app.RunContext(ctx, os.Args); err != nil {
		zap.L().Fatal("fatal error", zap.Error(err))
	}
}
