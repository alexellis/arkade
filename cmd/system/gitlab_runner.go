// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallGitLabRunner() *cobra.Command {
	command := &cobra.Command{
		Use:     "gitlab-runner",
		Short:   "Install GitLab Runner",
		Long:    `Install GitLab Runner for self-hosted CI.`,
		Aliases: []string{"gitlab"},
		Example: `  # Install GitLab Runner to the default path
  arkade system install gitlab-runner

  # Install a specific version
  arkade system install gitlab-runner --version <version>

  # Install on macOS
  arkade system install gitlab-runner --os darwin --arch arm64`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "$HOME/gitlab-runner", "Installation path, where gitlab-runner binary file is downloaded")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. x86_64 or arm64")
	command.Flags().String("os", "", "Operating system i.e. linux or darwin, leave blank to detect the client OS")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		fmt.Printf("Installing GitLab Runner to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		arch, osVer := env.GetClientArch()

		if cmd.Flags().Changed("os") {
			osVer, _ = cmd.Flags().GetString("os")
		}

		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		dlOS := strings.ToLower(osVer)
		if dlOS != "linux" && dlOS != "darwin" {
			return fmt.Errorf("unsupported operating system: %q, use linux or darwin (macOS)", osVer)
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "amd64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		} else if arch == "armv7" || arch == "armv7l" {
			dlArch = "arm"
		}

		if version == "" {
			version = "latest"
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s / %s\n", version, dlOS, dlArch)

		dlURL := fmt.Sprintf("https://gitlab-runner-downloads.s3.amazonaws.com/%s/binaries/gitlab-runner-%s-%s", version, dlOS, dlArch)

		fmt.Printf("Downloading from: %s\n", dlURL)

		progress, _ := cmd.Flags().GetBool("progress")
		outPath, err := get.DownloadFileP(dlURL, progress)
		if err != nil {
			return err
		}
		defer os.Remove(outPath)

		fmt.Printf("Downloaded to: %s\n", outPath)

		if _, err := get.CopyFile(outPath, installPath); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Binary file copied to: %s\n", installPath)

		return nil
	}

	return command
}
