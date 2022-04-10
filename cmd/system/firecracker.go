// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallFirecracker() *cobra.Command {

	command := &cobra.Command{
		Use:   "firecracker",
		Short: "Install Firecracker",
		Long:  `Install Firecracker and its Jailer.`,
		Example: `  arkade system install firecracker
  arkade system install firecracker --version v1.0.0`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", githubLatest, "The version for Firecracker to install")
	command.Flags().StringP("path", "p", "/usr/local/bin", "Installation path, where a go subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		owner := "firecracker-microvm"
		repo := "firecracker"

		fmt.Printf("Installing Firecracker to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		if arch != "x86_64" && arch != "aarch64" {
			return fmt.Errorf("this app only supports x86_64 and aarch64 and not %s", arch)
		}

		if version == githubLatest {
			v, err := get.FindGitHubRelease(owner, repo)
			if err != nil {
				return err
			}

			version = v
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, arch)

		filename := fmt.Sprintf("firecracker-%s-%s.tgz", version, arch)
		dlURL := fmt.Sprintf(githubDownloadTemplate, owner, repo, version, filename)

		fmt.Printf("Downloading from: %s\n", dlURL)
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

		tempUnpackPath, err := os.MkdirTemp(os.TempDir(), "firecracker*")
		if err != nil {
			return err
		}
		fmt.Printf("Unpacking Firecracker to: %s\n", tempUnpackPath)
		if err := archive.Untar(f, tempUnpackPath, true); err != nil {
			return err
		}

		fmt.Printf("Copying Firecracker binaries to: %s\n", installPath)
		filesToCopy := map[string]string{
			fmt.Sprintf("%s/firecracker-%s-%s", tempUnpackPath, version, arch): fmt.Sprintf("%s/firecracker", installPath),
			fmt.Sprintf("%s/jailer-%s-%s", tempUnpackPath, version, arch):      fmt.Sprintf("%s/jailer", installPath),
		}
		for src, dst := range filesToCopy {
			if _, err := get.CopyFileP(src, dst, readWriteExecuteEveryone); err != nil {
				return err
			}
		}

		return nil
	}

	return command
}
