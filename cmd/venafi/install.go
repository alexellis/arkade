// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package venafi

import "github.com/spf13/cobra"

func MakeInstall() *cobra.Command {

	command := &cobra.Command{
		Use:     "install",
		Short:   "Install Sponsored Apps for Venafi",
		Long:    `Install Sponsored Apps for Venafi`,
		Aliases: []string{"i"},
		Example: `  arkade venafi install [APP]
  arkade venafi install --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}
	return command
}
