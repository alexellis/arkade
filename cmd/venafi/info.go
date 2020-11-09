// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package venafi

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeInfo() *cobra.Command {

	command := &cobra.Command{
		Use:          "info",
		Short:        "Info for an app",
		Long:         `Info for an app`,
		Example:      `  arkade venafi info [APP]`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("give an app as an argument")
		}

		info := "None found."
		if args[0] == "cloud-issuer" {
			info = CloudIssuerInfo
		}

		fmt.Printf("Info for your app: %s\n\n%s\n\n", args[0], info)

		return nil
	}

	return command
}
