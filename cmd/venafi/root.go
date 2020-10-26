// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package venafi

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeVenafi() *cobra.Command {

	command := &cobra.Command{
		Use:   "venafi",
		Short: "Sponsored Apps - Venafi",
		Long: `Sponsored apps by Venafi.com. Venafi specialises in Machine Identity and is 
the custodian of cert-manager`,
		Example: `  arkade venafi install [APP]
  arkade venafi info [APP]
  arkade venafi --help
  arkade venafi install --help`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println("No apps available yet")

		return nil

	}

	return command
}
