package example

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/ekristen/go-project-template/pkg/common"
)

func Execute(_ context.Context, _ *cli.Command) error {
	log.Info().Msg("example called")
	return nil
}

func init() {
	cmd := &cli.Command{
		Name:        "example",
		Usage:       "example cli command",
		Description: `example command for the go-project-template`,
		Action:      Execute,
	}

	common.RegisterCommand(cmd)
}
