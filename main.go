// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"os"

	"github.com/alexellis/arkade/cmd"
	"github.com/spf13/cobra"
)

func main() {

	cmdVersion := cmd.MakeVersion()
	cmdInstall := cmd.MakeInstall()
	cmdInfo := cmd.MakeInfo()

	printarkadeASCIIArt := cmd.PrintArkadeASCIIArt

	var rootCmd = &cobra.Command{
		Use: "arkade",
		Run: func(cmd *cobra.Command, args []string) {
			printarkadeASCIIArt()
			cmd.Help()
		},
	}

	rootCmd.AddCommand(cmdInstall)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmd.MakeUpdate())
	rootCmd.AddCommand(cmd.MakeGet())
	rootCmd.AddCommand(cmd.MakeUninstall())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
