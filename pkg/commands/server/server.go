package server

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/ekristen/go-project-template/pkg/server"
)

func Execute(c *cli.Context) error {
	return server.RunServer(c.Context, &server.Options{
		Port: c.Int("port"),
		Log:  zap.L().With(zap.String("component", "server")),
	})
}

func init() {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   4242,
		},
	}

	cmd := &cli.Command{
		Name:        "server",
		Usage:       "run the api server",
		Description: "this is a restful base api-server with automatic openapi spec generation",
		Before:      common.Before,
		Flags:       append(common.Flags(), flags...),
		Action:      Execute,
	}

	common.RegisterCommand(cmd)
}