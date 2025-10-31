package example

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/ekristen/go-project-template/pkg/common"
)

func Execute(_ context.Context, _ *cli.Command) error {
	logrus.Info("example called")
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
