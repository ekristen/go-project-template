package example

import (
	"github.com/ekristen/go-project-template/pkg/common"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Execute(c *cli.Context) error {
	logrus.Info("example called")
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
