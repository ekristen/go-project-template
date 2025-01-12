package example

import (
	"github.com/swade1987/go-project-template/pkg/common"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var logger = zap.L()

func Execute(c *cli.Context) error {
	logger.Info("example called")
	return nil
}

func init() {
	cmd := &cli.Command{
		Name:        "example",
		Usage:       "example",
		Description: `example command for the go-project-template`,
		Before:      common.Before,
		Flags:       common.Flags(),
		Action:      Execute,
	}

	common.RegisterCommand(cmd)
}
