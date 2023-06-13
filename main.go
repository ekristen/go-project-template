package main

import (
	"os"
	"path"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/ekristen/go-project-template/pkg/common"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = common.AppVersion.Name
	app.Version = common.AppVersion.Summary
	app.Authors = []*cli.Author{
		{
			Name:  "Erik Kristensen",
			Email: "erik@erikkristensen",
		},
	}

	app.Before = common.Before
	app.Flags = common.Flags()

	app.Commands = common.GetCommands()
	app.CommandNotFound = func(context *cli.Context, command string) {
		logrus.Fatalf("Command %s not found.", command)
	}

	ctx := signals.SetupSignalContext()
	if err := app.RunContext(ctx, os.Args); err != nil {
		logrus.Fatal(err)
	}
}
