// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallPowershell() *cobra.Command {
	command := &cobra.Command{
		Use:   "pwsh",
		Short: "Install Powershell",
		Long:  `Install Powershell cross-platform task automation solution.`,
		Example: `  arkade system install pwsh
  arkade system install pwsh --version <version>`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "/opt/microsoft/powershell", "Installation path, where a powershell subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		tools := get.MakeTools()
		var tool *get.Tool
		for _, t := range tools {
			if t.Name == "pwsh" {
				tool = &t
				break
			}
		}

		if tool == nil {
			return fmt.Errorf("unable to find powershell definition")
		}

		if version == "" {
			v, err := get.FindGitHubRelease(tool.Owner, tool.Repo)
			if err != nil {
				return err
			}
			version = v
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		semVer := semver.MustParse(version)
		majorVersion := semVer.Major()

		installPath = fmt.Sprintf("%s/%d", installPath, majorVersion)

		fmt.Printf("Installing Powershell to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		progress, _ := cmd.Flags().GetBool("progress")
		tempPath, err := get.DownloadNested(tool, arch, osVer, version, installPath, progress, !progress)
		if err != nil {
			return err
		}

		err = get.MoveTo(tempPath, installPath)
		if err != nil {
			return err
		}

		lnPath := "/usr/bin/pwsh"
		fmt.Printf("Creating Symbolic link to: %s\n", lnPath)
		pwshPath := fmt.Sprintf("%s/pwsh", installPath)
		os.Symlink(pwshPath, lnPath)
		return nil
	}

	return command
}
