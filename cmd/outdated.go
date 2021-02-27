// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func isExecutable(mode os.FileMode) bool {
	return mode&0100 != 0
}

func MakeOutdated() *cobra.Command {
	var command = &cobra.Command{
		Use:          "outdated",
		Short:        "Check the latest version of installed tools",
		Example:      `  arkade outdated`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		tools := get.MakeTools()
		var toolNames []string = make([]string, len(tools))
		for _, t := range tools {
			toolNames = append(toolNames, t.Name)
		}

		config.InitUserDir()
		binDir := path.Join(config.GetUserDir(), "bin/")

		files, err := ioutil.ReadDir(binDir)
		if err != nil {
			log.Fatal(err)
		}

		foundTools := make([]string, 0)
		for _, f := range files {
			_, found := find(toolNames, f.Name())
			if found && isExecutable(f.Mode()) {
				foundTools = append(foundTools, f.Name())
			}
		}

		for i := 0; i < len(foundTools); i++ {
			var tool *get.Tool
			for _, target := range tools {
				if foundTools[i] == target.Name {
					tool = &target
					break
				}
			}

			_, versionStr := tool.GetInstalledVersion()
			fmt.Printf("found: %s, version %s\n", foundTools[i], versionStr)
		}

		// TODO(tuananh): find the current version of installed tools

		// TODO(tuananh): Fetch the latest version of the installed tools
	}
	return command
}
