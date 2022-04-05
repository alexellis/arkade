// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// system contains a suite of Sponsored Apps for arkade
package system

import (
	"github.com/spf13/cobra"
)

func MakeSystem() *cobra.Command {

	command := &cobra.Command{
		Use:     "system",
		Short:   "System apps",
		Long:    `Apps for systems.`,
		Aliases: []string{"s"},
		Example: `  arkade system install [APP]
  arkade s i [APP]`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeInstall())

	return command
}
