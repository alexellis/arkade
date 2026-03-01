package gha

import (
	"github.com/spf13/cobra"
)

func MakeGHA() *cobra.Command {

	command := &cobra.Command{
		Use:          "gha",
		Short:        "GitHub Actions utilities",
		Long:         `Utilities for GitHub Actions workflows.`,
		Example:      `  arkade gha upgrade --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeUpgrade())

	return command
}
