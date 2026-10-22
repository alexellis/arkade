// Copyright (c) arkade author(s) 2026. All rights reserved.
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

func MakeInstallTinyGo() *cobra.Command {

	command := &cobra.Command{
		Use:   "tinygo",
		Short: "Install TinyGo",
		Long:  `Install TinyGo, a Go compiler for microcontrollers, WebAssembly (WASM), and CLI tools.`,
		Example: `  arkade system install tinygo
  arkade system install tinygo --version 0.43.0`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version for TinyGo, or leave blank for latest")
	command.Flags().String("path", "/usr/local/", "Installation path, where a tinygo subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		fmt.Printf("Installing TinyGo to %s\n", installPath)

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		arch, osVer := env.GetClientArch()

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "amd64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		}

		if dlArch != "amd64" && dlArch != "arm64" {
			return fmt.Errorf("this app only supports amd64 and arm64 architectures, found: %s", arch)
		}

		if len(version) == 0 {
			v, err := get.FindGitHubRelease("tinygo-org", "tinygo")
			if err != nil {
				return err
			}

			version = v
		} else if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		versionNumber := strings.TrimPrefix(version, "v")
		downloadName := fmt.Sprintf("tinygo%s.linux-%s.tar.gz", versionNumber, dlArch)
		dlURL := fmt.Sprintf(githubDownloadTemplate, "tinygo-org", "tinygo", version, downloadName)
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

		fmt.Printf("Unpacking TinyGo to: %s\n", path.Join(installPath, "tinygo"))

		if err := spinWhile("Unpacking TinyGo", func() error {
			return archive.UntarNested(f, installPath, true, true, true, false)
		}); err != nil {
			return err
		}

		fmt.Printf(`TinyGo installed to: %s

export PATH=$PATH:%s
`, path.Join(installPath, "tinygo", "bin", "tinygo"), path.Join(installPath, "tinygo", "bin"))

		return nil
	}

	return command
}
