// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"os"

	"github.com/alexellis/arkade/cmd"
	"github.com/alexellis/arkade/cmd/kasten"
	"github.com/alexellis/arkade/cmd/venafi"
	"github.com/spf13/cobra"
)

func main() {
	printarkadeASCIIArt := cmd.PrintArkadeASCIIArt

	var rootCmd = &cobra.Command{
		Use: "arkade",
		Run: func(cmd *cobra.Command, args []string) {
			printarkadeASCIIArt()
			cmd.Help()
		},
	}

	rootCmd.AddCommand(cmd.MakeInstall())
	rootCmd.AddCommand(cmd.MakeVersion())
	rootCmd.AddCommand(cmd.MakeInfo())
	rootCmd.AddCommand(cmd.MakeUpdate())
	rootCmd.AddCommand(cmd.MakeGet())
	rootCmd.AddCommand(cmd.MakeUninstall())
	rootCmd.AddCommand(cmd.MakeShellCompletion())

	rootCmd.AddCommand(venafi.MakeVenafi())
	rootCmd.AddCommand(kasten.MakeK10())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
