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

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
