package cmd

import "github.com/spf13/cobra"

type baseCmd struct {
	cmd *cobra.Command
}

type cmder interface {
	getCommand() *cobra.Command
}

func addCommands(command *cobra.Command, cmds ...cmder) {
	for _, cmd := range cmds {
		c := cmd.getCommand()
		if c == nil {
			continue
		}
		command.AddCommand(c)
	}
}
