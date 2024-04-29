// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

func MakeInstallContainerd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "containerd",
		Short:        "Install containerd",
		Long:         "Install container container runtime.",
		Example:      `arkade system install containerd`,
		SilenceUsage: true,
	}

	cmd.Flags().StringP("version", "v", "", "Version of the containerd binary pack, leave blank for latest")
	cmd.Flags().String("path", "/usr/local/bin", "Install path, where the containerd binaries will installed")
	cmd.Flags().Bool("systemd", true, "Add and enable systemd service for containerd")
	cmd.Flags().Bool("progress", true, "Show download progress")
	cmd.Flags().String("arch", "", "CPU architecture i.e. amd64")

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("progress")
		if err != nil {
			return err
		}

		return nil
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")
		systemd, _ := cmd.Flags().GetBool("systemd")

		arch, osVer := env.GetClientArch()
		if cmd.Flags().Changed("arch") {
			archFlag, _ := cmd.Flags().GetString("arch")
			arch = archFlag
		}

		fmt.Printf("Installing containerd to %s\n", installPath)

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app currently only supports Linux")
		}

		if version == "" {
			latestVerison, err := get.FindGitHubRelease("containerd", "containerd")
			if err != nil {
				return err
			}
			version = latestVerison
		}

		downloadArch := ""

		if arch == "x86_64" {
			downloadArch = "amd64"
		} else if arch == "aarch64" {
			downloadArch = "arm64"
		} else {
			return fmt.Errorf("this app currently only supports arm64 and amd64 archs")
		}

		containerdTool := get.Tool{
			Name:    "containerd",
			Repo:    "containerd",
			Owner:   "containerd",
			Version: version,
			BinaryTemplate: `
			{{$archStr := .Arch}}
			{{- if eq .Arch "aarch64" -}}
			{{$archStr = "arm64"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$archStr = "amd64"}}
			{{- end -}}
			{{.Name}}-{{.VersionNumber}}-{{.OS}}-{{$archStr}}.tar.gz
			`,
		}

		url, err := containerdTool.GetURL(osVer, downloadArch, containerdTool.Version, !progress)
		if err != nil {
			return err
		}

		outPath, err := get.DownloadFileP(url, progress)
		if err != nil {
			return err
		}
		fmt.Printf("Downloaded to: %s\n", outPath)

		f, err := os.OpenFile(outPath, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		tempDirName := os.TempDir() + "/containerd"

		if err := archive.UntarNested(f, tempDirName, true, false); err != nil {
			return err
		}

		fmt.Printf("Copying containerd binaries to: %s\n", installPath)

		dir, err := os.ReadDir(tempDirName + "/bin")
		if err != nil {
			return err
		}

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		for _, entry := range dir {
			_, err := get.CopyFile(tempDirName+"/bin/"+entry.Name(), installPath+"/"+entry.Name())
			if err != nil {
				return err
			}
		}

		if systemd {
			systemdUnitName := "containerd.service"
			systemdUnitUrl := fmt.Sprintf("https://raw.githubusercontent.com/containerd/containerd/%s/%s", version, systemdUnitName)

			response, err := http.Get(systemdUnitUrl)
			if err != nil {
				return err
			}

			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("error fetching systemd unit file, status code: %d, body: %s", response.StatusCode, string(body))
			}

			body = bytes.ReplaceAll(body, []byte("/usr/local/bin/containerd"), []byte(installPath+"/containerd"))

			if err := createSystemdUnit(systemdUnitName, body); err != nil {
				return err
			}
		}

		return nil
	}

	return cmd
}

func createSystemdUnit(systemdUnitName string, content []byte) error {
	fmt.Printf("Creating systemd unit file\n")

	if err := os.WriteFile("/lib/systemd/system/"+systemdUnitName, content, os.FileMode(0700)); err != nil {
		return err
	}

	taskReload := execute.ExecTask{
		Command:     "systemctl",
		Args:        []string{"daemon-reload"},
		StreamStdio: false,
	}
	if _, err := taskReload.Execute(context.Background()); err != nil {
		return err
	}

	taskEnable := execute.ExecTask{
		Command:     "systemctl",
		Args:        []string{"enable", "--now", systemdUnitName},
		StreamStdio: false,
	}

	result, err := taskEnable.Execute(context.Background())
	if err != nil {
		return err
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("error enabling systemd service, stderr: %s", result.Stderr)
	}
	return nil
}
