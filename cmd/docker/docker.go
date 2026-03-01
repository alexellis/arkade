package docker

import (
	"github.com/spf13/cobra"
)

func MakeDocker() *cobra.Command {

	command := &cobra.Command{
		Use:          "docker",
		Short:        "Docker utilities",
		Long:         `Utilities for Dockerfiles.`,
		Aliases:      []string{"d"},
		Example:      `  arkade docker upgrade --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeUpgrade())

	return command
}
