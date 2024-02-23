// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// system contains packages for Linux servers and workstations
package oci

import (
	"github.com/spf13/cobra"
)

func MakeOci() *cobra.Command {

	command := &cobra.Command{
		Use:     "oci",
		Aliases: []string{"o"},
		Short:   "oci apps",
		Long:    `Apps from OCI images.`,
		Example: `  arkade oci install [container image]
  arkade oci i [container image]`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	}

	command.AddCommand(MakeOciInstall())

	return command
}
