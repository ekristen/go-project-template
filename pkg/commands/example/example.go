package example

import (
	"context"

	"github.com/urfave/cli/v3"

	"go.uber.org/zap"

	"github.com/ekristen/go-project-template/pkg/common"
)

func Execute(_ context.Context, _ *cli.Command) error {
	zap.L().Info("example called")
	return nil
}

func init() {
	cmd := &cli.Command{
		Name:        "example",
		Usage:       "example cli command",
		Description: `example command for the go-project-template`,
		Before:      common.Before,
		Flags:       common.Flags(),
		Action:      Execute,
	}

	common.RegisterCommand(cmd)
}
