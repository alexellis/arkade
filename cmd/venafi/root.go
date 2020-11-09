// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// venafi contains a suite of Sponsored Apps for arkade
package venafi

import (
	"github.com/spf13/cobra"
)

func MakeVenafi() *cobra.Command {

	command := &cobra.Command{
		Use:   "venafi",
		Short: "Sponsored Apps for Venafi",
		Long: `Sponsored apps for Venafi.com. Venafi specialises in Machine Identity and 
support for cert-manager.`,
		Aliases: []string{"v"},
		Example: `  arkade venafi install [APP]
  arkade venafi info [APP]`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	install := MakeInstall()
	install.AddCommand(MakeCloudIssuer())
	install.AddCommand(MakeTPPIssuer())
	command.AddCommand(install)
	command.AddCommand(MakeInfo())

	return command
}
