// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package system

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

func MakeInstallCaddyServer() *cobra.Command {
	command := &cobra.Command{
		Use:   "caddy",
		Short: "Install Caddy Server",
		Long:  `Install Caddy Server which is an extensible server platform that uses TLS by default`,
		Example: `  arkade system install caddy
  arkade system install caddy --version <version>`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "/usr/bin", "Installation path, where a caddy subfolder will be created")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		installPath, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		progress, _ := cmd.Flags().GetBool("progress")

		arch, osVer := env.GetClientArch()

		if cmd.Flags().Changed("os") {
			osVer, _ = cmd.Flags().GetString("os")
		}
		if cmd.Flags().Changed("arch") {
			arch, _ = cmd.Flags().GetString("arch")
		}

		if strings.ToLower(osVer) != "linux" {
			return fmt.Errorf("this app only supports Linux")
		}

		tools := get.MakeTools()
		var tool *get.Tool
		for _, t := range tools {
			if t.Name == "caddy" {
				tool = &t
				break
			}
		}

		if tool == nil {
			return fmt.Errorf("unable to find caddy definition")
		}

		fmt.Printf("Installing Caddy Server to %s\n", installPath)

		installPath = strings.ReplaceAll(installPath, "$HOME", os.Getenv("HOME"))

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		if version == "" {
			v, err := get.FindGitHubRelease("caddyserver", "caddy")
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

		svcTmpPath, err := get.DownloadFileP("https://raw.githubusercontent.com/caddyserver/dist/master/init/caddy.service", false)
		if err != nil {
			return err
		}
		fmt.Printf("Downloaded caddy.service file to %s\n", svcTmpPath)
		defer os.Remove(svcTmpPath)

		caddySystemFile := "/etc/systemd/system/caddy.service"
		if _, err = get.CopyFile(svcTmpPath, caddySystemFile); err != nil {
			return err
		}
		fmt.Printf("Copied caddy.service file to %s\n", caddySystemFile)

		caddyHomeDir := "/var/lib/caddy"
		caddyConfDir := "/etc/caddy"
		caddyUser := "caddy"
		if err = createCaddyConf(caddyHomeDir, caddyConfDir); err != nil {
			return err
		}
		fmt.Printf("Created caddy home %s and Conf %s directory\n", caddyHomeDir, caddyConfDir)

		if _, err = user.Lookup(caddyUser); errors.Is(err, user.UnknownUserError(caddyUser)) {
			if _, err = executeShellCmd(context.Background(), "useradd", "--system", "--home", caddyHomeDir, "--shell", "/bin/false", caddyUser); err != nil {
				return err
			}
			fmt.Printf("User created for caddy server.\n")
		}

		if _, err = executeShellCmd(context.Background(), "chown", "--recursive", "caddy:caddy", "/var/lib/caddy"); err != nil {
			return err
		}

		if _, err = executeShellCmd(context.Background(), "chown", "--recursive", "caddy:caddy", "/etc/caddy"); err != nil {
			return err
		}

		if _, err = executeShellCmd(context.Background(), "systemctl", "enable", "caddy"); err != nil {
			return err
		}

		if _, err = executeShellCmd(context.Background(), "systemctl", "daemon-reload"); err != nil {
			return err
		}

		return nil
	}

	return command
}

func createCaddyConf(caddyHomeDir, caddyConfDir string) error {
	os.MkdirAll(caddyHomeDir, 0755)
	os.MkdirAll(caddyConfDir, 0755)

	caddyConfFilePath := fmt.Sprintf("%s/Caddyfile", caddyConfDir)
	caddyFile, err := os.Create(caddyConfFilePath)
	if err != nil {
		return err
	}
	defer caddyFile.Close()
	return nil
}

func executeShellCmd(ctx context.Context, cmd string, parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     cmd,
		Args:        parts,
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute(ctx)

	if err != nil {
		return res, err
	}

	if res.ExitCode != 0 {
		return res, fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}

	return res, nil
}
