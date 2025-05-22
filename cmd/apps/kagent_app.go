// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/arkade/pkg/k8s"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

const (
	// A default Kagent version, get the latest from:
	// https://github.com/kagent-dev/kagent/releases/latest
	kagentVer = "0.3.2"
)

func MakeInstallKagent() *cobra.Command {
	var kagent = &cobra.Command{
		Use:          "kagent",
		Short:        "Install kagent",
		Long:         `Install kagent`,
		Example:      `  arkade install kagent`,
		SilenceUsage: true,
	}

	kagent.Flags().StringP("version", "v", kagentVer, "Specify a version of kagent")

	kagent.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("version")
		if err != nil {
			return fmt.Errorf("error with --version usage: %s", err)
		}

		_, err = command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}
		return nil
	}

	kagent.RunE = func(command *cobra.Command, args []string) error {
		version, _ := command.Flags().GetString("version")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if suffix := getValuesSuffix(arch); suffix == "-armhf" {
			return fmt.Errorf(`kagent is currently not supported on armhf architectures`)
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		arch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q\n", clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		err = downloadKagent(userPath, arch, clientOS, version)
		if err != nil {
			return err
		}

		fmt.Println("Running kagent check, this may take a few moments.")
		defaultFlags := []string{""}

		if len(kubeConfigPath) > 0 {
			defaultFlags = append(defaultFlags, "--kubeconfig", kubeConfigPath)
		}

		fmt.Println(kagentPostInstallMsg)
		return nil
	}

	return kagent
}

const KagentInfoMsg = `# Find out more at:
# https://kagent.dev
`
const kagentPostInstallMsg = `=======================================================================
= Kagent has been installed.                                        =
=======================================================================` +
	"\n\n" + KagentInfoMsg + "\n\n" + pkg.SupportMessageShort

func downloadKagent(userPath, arch, clientOS, version string) error {

	tools := get.MakeTools()
	var tool *get.Tool
	for _, t := range tools {
		if t.Name == "kagent" {
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

func kagentCli(parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     env.LocalBinary("kagentctl", ""),
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
