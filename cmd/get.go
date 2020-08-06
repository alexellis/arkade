// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

// MakeGet creates the Get command to download software
func MakeGet() *cobra.Command {
	tools := get.MakeTools()

	var command = &cobra.Command{
		Use:   "get",
		Short: `The get command downloads a tool`,
		Long: `The get command downloads a CLI or application from the specific tool's 
releases or downloads page. The tool is usually downloaded in binary format 
and provides a fast and easy alternative to a package manager.`,
		Example: `  arkade get faas-cli
  arkade get helm
  arkade get kind
  arkade get kubectx`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(arkadeGet)
			return nil
		}
		var tool *get.Tool

		if len(args) == 1 {
			for _, t := range tools {
				if t.Name == args[0] {
					tool = &t
					break
				}
			}
		}
		if tool == nil {
			return fmt.Errorf("cannot get tool: %s", args[0])
		}

		fmt.Printf("Downloading %s\n", tool.Name)

		arch, operatingSystem := env.GetClientArch()
		version := ""

		outFilePath, finalName, err := get.Download(tool, arch, operatingSystem, version, get.DownloadTempDir)

		if err != nil {
			return err
		}

		fmt.Printf(`Tool written to: %s

Run the following to copy to install the tool:

chmod +x %s
sudo install -m 755 %s /usr/local/bin/%s
`, outFilePath, outFilePath, outFilePath, finalName)

		return err
	}

	return command
}

const arkadeGet = `Use "arkade get TOOL" to download a tool or application:

- faas-cli
- helm
- inletsctl
- k3d
- k3sup
- kind
- kubectl
- kubectx
- kubeseal
  `
