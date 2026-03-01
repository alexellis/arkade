package gha

import (
	"github.com/alexellis/gha-bump/pkg/ghabump"
	"github.com/spf13/cobra"
)

func MakeGHA() *cobra.Command {

	command := &cobra.Command{
		Use:          "gha",
		Short:        "GitHub Actions utilities",
		Long:         `Utilities for GitHub Actions workflows.`,
		Example:      `  arkade gha bump --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeBump())

	return command
}

func MakeBump() *cobra.Command {
	var command = &cobra.Command{
		Use:     "bump",
		Short:   "Upgrade actions in GitHub Actions workflow files to the latest major version",
		Aliases: []string{"u"},
		Long: `Upgrade actions in GitHub Actions workflow files to the latest major version.

Processes all workflow YAML files in .github/workflows/ or a single file.
Only bumps major versions (e.g. actions/checkout@v3 to actions/checkout@v4).
`,
		Example: `  # Upgrade all workflows in the current directory
  arkade gha bump

  # Upgrade a single workflow file
  arkade gha bump -f .github/workflows/build.yaml

  # Dry-run mode, don't write changes
  arkade gha bump --write=false`,
		SilenceUsage: true,
	}

	command.Flags().StringP("file", "f", ".", "Path to workflow file or directory")
	command.Flags().BoolP("verbose", "v", true, "Verbose output")
	command.Flags().BoolP("write", "w", true, "Write the updated values back to the file")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("file")
		verbose, _ := cmd.Flags().GetBool("verbose")
		write, _ := cmd.Flags().GetBool("write")

		return ghabump.Run(ghabump.RunOptions{
			Target:  target,
			Verbose: verbose,
			Write:   write,
		})
	}

	return command
}
