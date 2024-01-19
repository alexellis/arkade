// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

var linkerdVersion = "stable-2.13.0"

func MakeInstallLinkerd() *cobra.Command {
	var linkerd = &cobra.Command{
		Use:          "linkerd",
		Short:        "Install linkerd",
		Long:         `Install linkerd the light-weight service mesh for Kubernetes.`,
		Example:      `  arkade install linkerd`,
		SilenceUsage: true,
	}

	linkerd.Flags().StringP("version", "v", linkerdVersion, "Specify a version of linkerd")

	linkerd.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		version, err := command.Flags().GetString("version")
		if err != nil {
			return err
		}
		if command.Flags().Changed("version") {
			linkerdVersion = version
		}

		if len(version) == 0 {
			return fmt.Errorf("you must give a value for the --version flag")
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		userPath, err := getUserPath()
		if err != nil {
			return err
		}

		arch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q\n", clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		err = downloadLinkerd(userPath, arch, clientOS, linkerdVersion)
		if err != nil {
			return err
		}
		fmt.Println("Running linkerd2 check, this may take a few moments.")

		_, err = linkerdCli("check", "--pre")
		if err != nil {
			return err
		}

		res, err := linkerdCli("install")
		if err != nil {
			return err
		}
		file, err := os.CreateTemp("", "linkerd")
		if err != nil {
			return err
		}

		w := bufio.NewWriter(file)
		_, err = w.WriteString(res.Stdout)
		if err != nil {
			return err
		}
		w.Flush()

		err = k8s.Kubectl("apply", "-R", "-f", file.Name())
		if err != nil {
			return err
		}

		defer os.Remove(file.Name())

		_, err = linkerdCli("check")
		if err != nil {
			return err
		}

		fmt.Println(`=======================================================================
= Linkerd has been installed.                                        =
=======================================================================

# Find out more at:
# https://linkerd.io

# To use the linkerd2 CLI set this path:

export PATH=$PATH:` + path.Join(userPath, "bin/") + `
linkerd2 --help

` + pkg.SupportMessageShort)
		return nil
	}

	return linkerd
}

// func getLinkerdURL(os, version string) string {
// 	osSuffix := strings.ToLower(os)
// 	return fmt.Sprintf("https://github.com/linkerd/linkerd2/releases/download/%s/linkerd2-cli-%s-%s", version, version, osSuffix)
// }

func downloadLinkerd(userPath, arch, clientOS, version string) error {

	tools := get.MakeTools()
	var tool *get.Tool
	for _, t := range tools {
		if t.Name == "linkerd2" {
			tool = &t
			break
		}
	}

	if tool == nil {
		return fmt.Errorf("unable to find tool definition")
	}

	if _, err := os.Stat(env.LocalBinary(tool.Name, "")); errors.Is(err, os.ErrNotExist) {

		var (
			progress bool
			quiet    bool
		)

		defaultMovePath := ""
		outPath, finalName, err := get.Download(tool, arch, clientOS, version, defaultMovePath, progress, quiet)
		if err != nil {
			return err
		}

		fmt.Println("Downloaded to: ", outPath, finalName)
	} else {
		fmt.Printf("%s already exists, skipping download.\n", tool.Name)
	}

	return nil
}

func linkerdCli(parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     env.LocalBinary("linkerd2", ""),
		Args:        parts,
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute(context.Background())

	if err != nil {
		return res, err
	}

	if res.ExitCode != 0 {
		return res, fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}

	return res, nil
}

func getUserPath() (string, error) {
	userPath, err := config.InitUserDir()
	return userPath, err
}

func getExportPath() string {
	userPath := config.GetUserDir()
	return path.Join(userPath, "bin/")
}

var LinkerdInfoMsg = `# Find out more at:
# https://linkerd.io

# To use the linkerd2 CLI set this path:

export PATH=$PATH:` + getExportPath() + `
linkerd2 --help`
