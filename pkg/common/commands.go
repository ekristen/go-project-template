package common

import (
	"github.com/urfave/cli/v2"

	"go.uber.org/zap"
)

var logger = zap.S()
var commands = make(map[string][]*cli.Command, 0)

// RegisterSubcommand allows you to register a command under a group
func RegisterSubcommand(group string, command *cli.Command) {
	logger.Debug("Registering", command.Name, "command...")
	commands[group] = append(commands[group], command)
}

// GetSubcommands retrieves all commands assigned to a group
func GetSubcommands(group string) []*cli.Command {
	return commands[group]
}

// RegisterCommand -- allows you to register a command under the main group
func RegisterCommand(command *cli.Command) {
	logger.Debug("Registering", command.Name, "command...")
	commands["_main_"] = append(commands["_main_"], command)
}

// GetCommands -- retrieves all commands assigned to the main group
func GetCommands() []*cli.Command {
	return commands["_main_"]
}
