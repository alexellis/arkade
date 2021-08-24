// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// k10 contains a suite of Sponsored Apps for arkade
package k10

import (
	"github.com/spf13/cobra"
)

func MakeK10() *cobra.Command {

	command := &cobra.Command{
		Use:   "k10",
		Short: "Sponsored Apps for K10",
		Long:  `Sponsored apps for kasten.io. Kasten K10 by Veeam is purpose-built for Kubernetes backup and restore`,
		Example: `  arkade k10 install [APP]
  arkade k10 info [APP]`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeInfo())

	return command
}

// Kasten K10 by Veeam is purpose-built for Kubernetes backup and restore. Kasten K10 uses an application-aware approach, is easy-to-use, secure and free.
