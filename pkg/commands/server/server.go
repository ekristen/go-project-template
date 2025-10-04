package server

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/ekristen/go-project-template/pkg/server"
)

func Execute(ctx context.Context, c *cli.Command) error {
	// Get the telemetry instance for tracing and metrics
	telemetryInstance := common.GetTelemetry()

	return server.Run(ctx, &server.Options{
		Port:      c.Int("port"),
		Telemetry: telemetryInstance,
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
		Description: "this is a restful base api server with automatic openapi spec generation",
		Flags:       flags,
		Action:      Execute,
	}

	common.RegisterCommand(cmd)
}
