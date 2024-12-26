package cmds

import (
	"github.com/spf13/cobra"
)

type CommandOptions struct {
	Repositories []string
}

// NewCommandOptions creates a new CommandOptions with the given repositories
func NewCommandOptions(repositories []string) *CommandOptions {
	return &CommandOptions{
		Repositories: repositories,
	}
}

// NewCommands creates all prompto commands with the given options
func NewCommands(options *CommandOptions) []*cobra.Command {
	return []*cobra.Command{
		NewGetCommand(options),
		NewListCommand(options),
		NewServeCommand(options),
		NewWhichCommand(options),
		NewEditCommand(options),
	}
}
