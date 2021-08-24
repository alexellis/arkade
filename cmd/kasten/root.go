// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// kasten contains a suite of Sponsored Apps for arkade
package kasten

import (
	"github.com/spf13/cobra"
)

func MakeK10() *cobra.Command {

	command := &cobra.Command{
		Use:     "kasten",
		Short:   "Sponsored Apps for kasten",
		Long:    `Sponsored apps for kasten.io. Kasten K10 by Veeam is purpose-built for Kubernetes backup and restore`,
		Aliases: []string{"k10"},
		Example: `  arkade kasten install [APP]
  arkade kasten info [APP]`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeInfo())
	command.AddCommand(MakeInstall())

	return command
}

// Kasten K10 by Veeam is purpose-built for Kubernetes backup and restore. Kasten K10 uses an application-aware approach, is easy-to-use, secure and free.
