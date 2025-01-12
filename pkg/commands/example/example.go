package example

import (
	"github.com/urfave/cli/v2"

	"go.uber.org/zap"

	"github.com/ekristen/go-project-template/pkg/common"
)

func Execute(c *cli.Context) error {
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
