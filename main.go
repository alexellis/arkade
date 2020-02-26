package main

import (
	"os"

	"github.com/alexellis/bazaar/cmd"
	"github.com/spf13/cobra"
)

func main() {

	cmdVersion := cmd.MakeVersion()
	cmdInstall := cmd.MakeInstall()
	cmdInfo := cmd.MakeInfo()

	printbazaarASCIIArt := cmd.PrintBazaarASCIIArt

	var rootCmd = &cobra.Command{
		Use: "bazaar",
		Run: func(cmd *cobra.Command, args []string) {
			printbazaarASCIIArt()
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
