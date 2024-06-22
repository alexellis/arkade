// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallGo() *cobra.Command {

	command := &cobra.Command{
		Use:   "go",
		Short: "Install Go",
		Long:  `Install Go programming language and SDK.`,
		Example: `  arkade system install go
  arkade system install go --version v1.18.1`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version for Go, or leave blank for pinned version")
	command.Flags().String("path", "/usr/local/", "Installation path, where a go subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {

		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		fmt.Printf("Installing Go to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		tools := get.MakeTools()
		var tool *get.Tool
		for _, t := range tools {
			if t.Name == "go" {
				tool = &t
				break
			}
		}

		if tool == nil {
			return fmt.Errorf("unable to find go definition")
		}

		progress, _ := cmd.Flags().GetBool("progress")
		tempPath, err := get.DownloadNested(tool, arch, osVer, version, installPath, progress, !progress)
		if err != nil {
			return err
		}

		err = get.MoveTo(tempPath, installPath)
		if err != nil {
			return err
		}

		fmt.Printf("Downloaded to: %sgo\n", installPath)

		fmt.Printf("\nexport PATH=$PATH:%s:$HOME/go/bin\n"+
			"export GOPATH=$HOME/go/\n", path.Join(installPath, "go", "bin"))

		return nil
	}

	return command
}
