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

func MakeInstallCNI() *cobra.Command {

	command := &cobra.Command{
		Use:   "cni",
		Short: "Install CNI plugins",
		Long:  `Install CNI plugins for use with faasd, actuated, Kubernetes, etc.`,
		Example: `  arkade system install cni
  arkade system install cni --version v1.4.0`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "v1.4.0", "The version for CNI to install")
	command.Flags().StringP("path", "p", "/opt/cni/bin/", "Installation path, where a go subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		owner := "containernetworking"
		repo := "plugins"

		fmt.Printf("Installing CNI to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()
		if cmd.Flags().Changed("arch") {
			archFlag, _ := cmd.Flags().GetString("arch")
			arch = archFlag
		}

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "amd64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		} else if arch == "armv7" || arch == "armv7l" {
			dlArch = "arm"
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

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		filename := fmt.Sprintf("cni-plugins-linux-%s-%s.tgz", dlArch, version)
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

		tempUnpackPath, err := os.MkdirTemp(os.TempDir(), "cni-plugins*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempUnpackPath)

		fmt.Printf("Unpacking CNI Plugins to: %s\n", tempUnpackPath)
		if err := archive.Untar(f, tempUnpackPath, true, true); err != nil {
			return err
		}

		dirs, err := os.ReadDir(tempUnpackPath)
		if err != nil {
			return err
		}

		for _, dir := range dirs {
			if !dir.IsDir() {
				src := path.Join(tempUnpackPath, dir.Name())
				dst := path.Join(installPath, dir.Name())
				fmt.Printf("Copying %s to: %s\n", src, dst)

				if _, err := get.CopyFileP(src, dst, readWriteExecuteEveryone); err != nil {
					return fmt.Errorf("unable to copy %s to %s, error: %s", src, dst, err)
				}
			}
		}

		return nil
	}

	return command
}
