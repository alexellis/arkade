// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// kasten contains a suite of Sponsored Apps for arkade
package kasten

import "github.com/spf13/cobra"

func MakeInstall() *cobra.Command {

	command := &cobra.Command{
		Use:     "install",
		Short:   "Install Sponsored Apps for kasten",
		Long:    `Install Sponsored Apps for kasten`,
		Aliases: []string{"i"},
		Example: `  arkade kasten install [APP]
  arkade kasten install --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeInstallK10())
	command.AddCommand(MakeInstallK10Preflight())

	return command
}
