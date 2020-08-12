// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/alexellis/arkade/cmd/apps"
	"github.com/spf13/cobra"
)

func MakeInstall() *cobra.Command {
	var command = &cobra.Command{
		Use:   "install",
		Short: "Install Kubernetes apps from helm charts or YAML files",
		Long: `Install Kubernetes apps from helm charts or YAML files using the "install"
command. Helm 3 is used by default unless you pass --helm3=false, then helm 2
will be used to generate YAML files which are applied without tiller.

You can also find the post-install message for each app with the "info"
command.`,
		Example: `  arkade install
  arkade install openfaas --helm3 --gateways=2
  arkade install inlets-operator --token-file $HOME/do-token`,
		SilenceUsage: false,
	}

	command.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")
	command.PersistentFlags().Bool("wait", false, "If we should wait for the resource to be ready before returning (helm3 only, default false)")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			fmt.Printf(
				`To see a complete list of apps run:

  arkade install --help

And to see options for a specific app before installing, run:

  arkade install APP --help
`)
			return nil
		}

		return nil
	}

	for _, app := range apps.Apps {
		command.AddCommand(app.MakeInstall)
	}

	command.AddCommand(MakeInfo())

	return command
}
