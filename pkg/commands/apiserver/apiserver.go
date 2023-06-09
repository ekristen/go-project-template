package apiserver

import (
	"github.com/ekristen/go-project-template/pkg/apiserver"
	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/urfave/cli/v2"
)

func Execute(c *cli.Context) error {
	return apiserver.RunServer(c.Context, &apiserver.Options{
		Port: c.Int("port"),
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
		Name:        "api-server",
		Usage:       "api-server",
		Description: "api-server",
		Before:      common.Before,
		Flags:       append(common.Flags(), flags...),
		Action:      Execute,
	}

	common.RegisterCommand(cmd)
}
