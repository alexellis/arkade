// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/archive"
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

		dlArch := arch
		if arch == "x86_64" {
			dlArch = "amd64"
		} else if arch == "aarch64" {
			dlArch = "arm64"
		} else if arch == "armv7" || arch == "armv7l" {
			dlArch = "armv6l"
		}

		if len(version) == 0 {
			v, err := getGoVersion()
			if err != nil {
				return err
			}

			version = v
		} else if !strings.HasPrefix(version, "go") {
			version = "go" + version
		}

		fmt.Printf("Installing version: %s for: %s\n", version, dlArch)

		dlURL := fmt.Sprintf("https://go.dev/dl/%s.%s-%s.tar.gz", version, strings.ToLower(osVer), dlArch)
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

		fmt.Printf("Unpacking Go to: %s\n", path.Join(installPath, "go"))

		if err := archive.UntarNested(f, installPath, true, false); err != nil {
			return err
		}

		fmt.Printf("\nexport PATH=$PATH:%s:$HOME/go/bin\n"+
			"export GOPATH=$HOME/go/\n", path.Join(installPath, "go", "bin"))

		return nil
	}

	return command
}

func getGoVersion() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://go.dev/VERSION?m=text", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body == nil {
		return "", fmt.Errorf("unexpected empty body")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	content := strings.TrimSpace(string(body))
	version, _, ok := strings.Cut(content, "\n")
	if !ok {
		return "", fmt.Errorf("format unexpected: %q", content)
	}

	return version, nil
}
