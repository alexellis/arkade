package cmd

import (
	"fmt"
	"strings"

	"github.com/alexellis/arkade/cmd/bundle"
	"github.com/spf13/cobra"
)

func MakeBundle() *cobra.Command {
	var command = &cobra.Command{
		Use:          "bundle",
		Short:        "Bundle installs a curated list of apps",
		Long:         `Bundle installs a curated list of apps. Currently only helm3 apps are supported.`,
		Example:      `  arkade bundle minicloud`,
		SilenceUsage: false,
	}

	//command.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			fmt.Printf("You can install: %s\n%s\n\n", strings.TrimRight("\n - "+strings.Join(getBundle(), "\n - "), "\n - "),
				`Run arkade bundle NAME --help to see configuration options.`)
			return nil
		}

		return nil
	}

	command.AddCommand(bundle.MakeBundleMinicloud())

	return command
}

func getBundle() []string {
	return []string{"minicloud"}
}
