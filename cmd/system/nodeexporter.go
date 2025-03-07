// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

func MakeInstallNodeExporter() *cobra.Command {
	command := &cobra.Command{
		Use:   "node_exporter",
		Short: "Install Node Exporter",
		Long: `Install Node Exporter which is a Prometheus exporter for hardware and OS 
metrics exposed by a server/container such as CPU/RAM/Disk/Network.`,
		RunE:         installNodeExporterE,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version for Go, or leave blank for pinned version")
	command.Flags().String("path", "/usr/local/", "Installation path, where a go subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().Bool("systemd", false, "Install and start a systemd service")

	return command
}

func installNodeExporterE(cmd *cobra.Command, args []string) error {
	installPath, _ := cmd.Flags().GetString("path")
	version, _ := cmd.Flags().GetString("version")
	progress, _ := cmd.Flags().GetBool("progress")
	systemd, _ := cmd.Flags().GetBool("systemd")

	arch, osVer := env.GetClientArch()

	if cmd.Flags().Changed("os") {
		osVer, _ = cmd.Flags().GetString("os")
	}
	if cmd.Flags().Changed("arch") {
		arch, _ = cmd.Flags().GetString("arch")
	}

	if strings.ToLower(osVer) != "linux" && strings.ToLower(osVer) != "darwin" {
		return fmt.Errorf("this app only supports Linux and Darwin")
	}

	tools := get.MakeTools()
	var tool *get.Tool
	for _, t := range tools {
		if t.Name == "node_exporter" {
			tool = &t
			break
		}
	}

	if tool == nil {
		return fmt.Errorf("unable to find node_exporter definition")
	}

	fmt.Printf("Installing node_exporter Server to %s\n", installPath)

	installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

	if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
	}

	if version == "" {
		v, err := get.FindGitHubRelease("prometheus", "node_exporter")
		if err != nil {
			return err
		}
		version = v
	} else if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	outFilePath, _, err := get.Download(tool, arch, osVer, version, installPath, progress, !progress)
	if err != nil {
		return err
	}
	if err = os.Chmod(outFilePath, readWriteExecuteEveryone); err != nil {
		return err
	}

	if systemd && strings.ToLower(osVer) != "linux" {
		return fmt.Errorf("systemd is only supported on Linux")
	}

	if systemd {
		systemdUnit := generateUnit(outFilePath)
		unitName := "/etc/systemd/system/node_exporter.service"
		if err := os.WriteFile(unitName, []byte(systemdUnit), readWriteExecuteEveryone); err != nil {
			return err
		}

		fmt.Printf("Wrote: %s\n", unitName)

		if _, err = executeShellCmd(context.Background(), "systemctl", "daemon-reload"); err != nil {
			return err
		}

		if _, err = executeShellCmd(context.Background(), "systemctl", "enable", "node_exporter", "--now"); err != nil {
			return err
		}

		fmt.Printf(`Started service: node_exporter

Check status with: sudo journalctl -u node_exporter -f

View metrics at: http://127.0.0.1:9100/metrics

`)

	}
	return nil
}

func generateUnit(outFilePath string) string {
	return fmt.Sprintf(`[Unit]
Description=Node Exporter
After=network.target

[Service]
ExecStart=%s

[Install]
WantedBy=multi-user.target
`, outFilePath)
}
