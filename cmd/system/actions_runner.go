// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallActionsRunner() *cobra.Command {

	command := &cobra.Command{
		Use:   "actions-runner",
		Short: "Install GitHub Actions Runner",
		Long:  `Install GitHub Actions Runner for self-hosted CI.`,
		Example: `  arkade system install actions-runner
  arkade system install actions-runner --version 2.290.1`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "$HOME/actions-runner/", "Installation path, where a actions-runner subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {

		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		fmt.Printf("Installing Actions Runner to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

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
			dlArch = "arm"
		}

		if version == "" {
			v, err := get.FindGitHubRelease("actions", "runner")
			if err != nil {
				return err
			}
			version = v
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		filename := fmt.Sprintf("actions-runner-linux-%s-%s.tar.gz", dlArch, strings.TrimPrefix(version, "v"))
		dlURL := fmt.Sprintf(githubDownloadTemplate, "actions", "runner", version, filename)
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

		fmt.Printf("Unpacking Actions Runner to: %s\n", path.Join(installPath, "actions-runner"))

		if err := archive.UntarNested(f, installPath); err != nil {
			return err
		}

		return nil
	}

	return command
}
