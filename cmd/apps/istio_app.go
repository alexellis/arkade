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

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/arkade/pkg/k8s"
	execute "github.com/alexellis/go-execute/v2"
	"github.com/spf13/cobra"
)

const (
	// A default Istio version, get the latest from:
	// https://github.com/istio/istio/releases/latest
	istioVer = "1.16.1"
)

func MakeInstallIstio() *cobra.Command {
	var istio = &cobra.Command{
		Use:          "istio",
		Short:        "Install istio",
		Long:         `Install istio`,
		Example:      `  arkade install istio --loadbalancer`,
		SilenceUsage: true,
	}

	istio.Flags().StringP("version", "v", istioVer, "Specify a version of Istio")
	istio.Flags().String("namespace", "default", "Namespace for the app")
	istio.Flags().String("istio-namespace", "istio-system", "Namespace for the app")
	istio.Flags().String("profile", "default", "Set istio profile")
	istio.Flags().String("cpu", "100m", "Allocate CPU resource")
	istio.Flags().String("memory", "100Mi", "Allocate Memory resource")

	istio.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set prometheus.enabled=false)")

	istio.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("version")
		if err != nil {
			return fmt.Errorf("error with --version usage: %s", err)
		}

		_, err = command.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("error with --namespace usage: %s", err)
		}

		_, err = command.Flags().GetString("istio-namespace")
		if err != nil {
			return fmt.Errorf("error with --istio-namespace usage: %s", err)
		}

		_, err = command.Flags().GetString("profile")
		if err != nil {
			return fmt.Errorf("error with --profile usage: %s", err)
		}

		_, err = command.Flags().GetString("cpu")
		if err != nil {
			return fmt.Errorf("error with --cpu usage: %s", err)
		}

		_, err = command.Flags().GetString("memory")
		if err != nil {
			return fmt.Errorf("error with --memory usage: %s", err)
		}

		_, err = command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		return nil
	}

	istio.RunE = func(command *cobra.Command, args []string) error {
		version, _ := command.Flags().GetString("version")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := command.Flags().GetString("namespace")
		istioNamespace, _ := command.Flags().GetString("istio-namespace")
		profile, _ := command.Flags().GetString("profile")
		cpu, _ := command.Flags().GetString("cpu")
		memory, _ := command.Flags().GetString("memory")
		setOverrides, _ := command.Flags().GetStringArray("set")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if suffix := getValuesSuffix(arch); suffix == "-armhf" {
			return fmt.Errorf(`istio is currently not supported on armhf architectures`)
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		arch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q\n", clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		err = downloadIstio(userPath, arch, clientOS, version)
		if err != nil {
			return err
		}

		fmt.Println("Running istio check, this may take a few moments.")
		defaultFlags := []string{"--namespace", namespace, "--istioNamespace", istioNamespace}

		if len(kubeConfigPath) > 0 {
			defaultFlags = append(defaultFlags, "--kubeconfig", kubeConfigPath)
		}

		preCheckFlags := mergeFlagsSlices([]string{"experimental", "precheck"}, defaultFlags)
		_, err = istioCli(preCheckFlags...)
		if err != nil {
			return err
		}

		// set installation flags
		installFlags := mergeFlagsSlices([]string{"install", "--skip-confirmation",
			"--set", fmt.Sprintf("values.pilot.resources.requests.cpu=%s", cpu),
			"--set", fmt.Sprintf("values.pilot.resources.requests.memory=%s", memory),
			"--set", fmt.Sprintf("profile=%s", profile)}, defaultFlags, fmtSetFlags(setOverrides))

		res, err := istioCli(installFlags...)
		if err != nil {
			return err
		}
		file, err := os.CreateTemp("", "istio")
		if err != nil {
			return err
		}

		w := bufio.NewWriter(file)
		_, err = w.WriteString(res.Stdout)
		if err != nil {
			return err
		}
		w.Flush()

		defer os.Remove(file.Name())

		verifyFlags := mergeFlagsSlices([]string{"verify-install"}, defaultFlags)
		_, err = istioCli(verifyFlags...)
		if err != nil {
			return err
		}

		fmt.Println(istioPostInstallMsg)
		return nil
	}

	return istio
}

const IstioInfoMsg = `# Find out more at:
# https://github.com/istio/`

const istioPostInstallMsg = `=======================================================================
= Istio has been installed.                                        =
=======================================================================` +
	"\n\n" + IstioInfoMsg + "\n\n" + pkg.SupportMessageShort

func downloadIstio(userPath, arch, clientOS, version string) error {

	tools := get.MakeTools()
	var tool *get.Tool
	for _, t := range tools {
		if t.Name == "istioctl" {
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

func istioCli(parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     env.LocalBinary("istioctl", ""),
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

func fmtSetFlags(setOverrides []string) []string {
	fmtOverrides := []string{}
	for _, setOverride := range setOverrides {
		fmtOverrides = append(fmtOverrides, "--set", setOverride)
	}
	return fmtOverrides
}

func mergeFlagsSlices(args ...[]string) []string {
	mergedSlice := make([]string, 0)
	for _, oneSlice := range args {
		mergedSlice = append(mergedSlice, oneSlice...)
	}

	return mergedSlice
}
