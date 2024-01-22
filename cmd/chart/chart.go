// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// chart contains packages for Linux servers and workstations
package chart

import (
	"github.com/spf13/cobra"
)

func MakeChart() *cobra.Command {

	command := &cobra.Command{
		Use:     "chart",
		Short:   "Chart utilities",
		Long:    `Utilities for Helm charts.`,
		Aliases: []string{"c"},
		Example: `  arkade chart verify --help
`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeVerify())
	command.AddCommand(MakeUpgrade())
	command.AddCommand(MakeBump())

	return command
}
