// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import "github.com/spf13/cobra"

func MakeInstall() *cobra.Command {

	command := &cobra.Command{
		Use:     "install",
		Short:   "Install system apps",
		Long:    `Install system apps for Linux hosts`,
		Aliases: []string{"i"},
		Example: `  arkade system install [APP]
  arkade system install --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeInstallGo())
	command.AddCommand(MakeInstallFirecracker())
	command.AddCommand(MakeInstallPrometheus())

	return command
}
