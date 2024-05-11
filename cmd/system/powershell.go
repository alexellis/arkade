// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/alexellis/arkade/pkg/archive"
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

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "x64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		} else if arch == "armv7" || arch == "armv7l" {
			dlArch = "arm32"
		}

		if version == "" {
			v, err := get.FindGitHubRelease("PowerShell", "PowerShell")
			if err != nil {
				return err
			}
			version = v
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		semVer := semver.MustParse(version)
		majorVersion := semVer.Major()
		// semVer := strings.TrimPrefix(version, "v")

		// majorVerDemlimiter := strings.Index(semVer, ".")
		// majorVersion := semVer[:majorVerDemlimiter]

		installPath = fmt.Sprintf("%s/%d", installPath, majorVersion)

		fmt.Printf("Installing Powershell to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		filename := fmt.Sprintf("powershell-%s-linux-%s.tar.gz", semVer, dlArch)
		dlURL := fmt.Sprintf(githubDownloadTemplate, "PowerShell", "PowerShell", version, filename)

		fmt.Printf("Downloading from: %s\n", dlURL)

		progress, _ := cmd.Flags().GetBool("progress")
		outPath, err := get.DownloadFileP(dlURL, progress)
		if err != nil {
			return err
		}
		defer os.Remove(outPath)

		fmt.Printf("Downloaded to: %s\n", outPath)

		f, err := os.OpenFile(outPath, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Printf("Unpacking Powershell to: %s\n", installPath)

		if err := archive.Untar(f, installPath, true, true); err != nil {
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
