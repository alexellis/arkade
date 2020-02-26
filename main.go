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
	cmdUpdate := cmd.MakeUpdate()

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
	rootCmd.AddCommand(cmdUpdate)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
