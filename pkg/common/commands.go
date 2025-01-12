package common

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var commands = make(map[string][]*cli.Command, 0)
var logger = zap.L()

// RegisterSubcommand allows you to register a command under a group
func RegisterSubcommand(group string, command *cli.Command) {
	logger.Debug("Registering", zap.String("command", command.Name), zap.String("group", group))
	commands[group] = append(commands[group], command)
}

// GetSubcommands retrieves all commands assigned to a group
func GetSubcommands(group string) []*cli.Command {
	return commands[group]
}

// RegisterCommand -- allows you to register a command under the main group
func RegisterCommand(command *cli.Command) {
	logger.Debug("Registering", zap.String("command", command.Name), zap.String("group", "_main_"))
	commands["_main_"] = append(commands["_main_"], command)
}

// GetCommands -- retrieves all commands assigned to the main group
func GetCommands() []*cli.Command {
	logger.Info("retrieving commands",
		zap.Int("command_count", len(commands["_main_"])),
		zap.Any("registered_commands", commands),
	)
	return commands["_main_"]
}
