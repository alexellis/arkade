// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

func PrintArkadeASCIIArt() {
	arkadeLogo := aec.BlueF.Apply(arkadeFigletStr)
	fmt.Print(arkadeLogo)
}

func MakeVersion() *cobra.Command {
	var jsonOut bool

	var command = &cobra.Command{
		Use:          "version",
		Short:        "Print the version",
		Example:      `  arkade version`,
		Aliases:      []string{"v"},
		SilenceUsage: false,
	}

	command.Flags().BoolVarP(&jsonOut, "json", "j", false, "Output version as JSON")

	command.Run = func(cmd *cobra.Command, args []string) {
		if jsonOut {
			out := map[string]string{
				"version":    pkg.BuildString(),
				"commit":     "n/a",
				"build_date": "n/a",
			}
			if len(pkg.GitCommit) > 0 {
				out["commit"] = pkg.GitCommit
			}
			if bd := pkg.BuildDateString(); bd != "" {
				out["build_date"] = bd
			}

			b, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(b))
			return
		}

		PrintArkadeASCIIArt()
		commit := "n/a"
		if len(pkg.GitCommit) > 0 {
			commit = pkg.GitCommit
		}
		bd := "n/a"
		if buildDate := pkg.BuildDateString(); buildDate != "" {
			bd = buildDate
		}
		fmt.Printf("  commit:  %s\n", commit)
		fmt.Printf("  version: %s\n", pkg.BuildString())
		fmt.Printf("  build date: %s\n", bd)

		fmt.Println("\n", aec.Bold.Apply(pkg.SupportMessageShort))
	}
	return command
}

const arkadeFigletStr = `            _             _      
  __ _ _ __| | ____ _  __| | ___ 
 / _` + "`" + ` | '__| |/ / _` + "`" + ` |/ _` + "`" + ` |/ _ \
| (_| | |  |   < (_| | (_| |  __/
 \__,_|_|  |_|\_\__,_|\__,_|\___|

Open Source Marketplace For Developer Tools

`
