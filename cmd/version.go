// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
)

func PrintArkadeASCIIArt() {
	arkadeLogo := aec.BlueF.Apply(arkadeFigletStr)
	fmt.Print(arkadeLogo)
}

func MakeVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "Print the version",
		Example:      `  arkade version`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		PrintArkadeASCIIArt()
		if len(Version) == 0 {
			fmt.Println("Version: dev")
		} else {
			fmt.Println("Version:", Version)
		}
		fmt.Println("Git Commit:", GitCommit)
	}
	return command
}

const arkadeFigletStr = `            _             _      
  __ _ _ __| | ____ _  __| | ___ 
 / _` + "`" + ` | '__| |/ / _` + "`" + ` |/ _` + "`" + ` |/ _ \
| (_| | |  |   < (_| | (_| |  __/
 \__,_|_|  |_|\_\__,_|\__,_|\___|

Get Kubernetes apps the easy way

`
