// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package helm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	execute "github.com/alexellis/go-execute/pkg/v1"
)

func TryDownloadHelm(userPath, clientArch, clientOS string) (string, error) {
	helmVal := "helm"
	subdir := ""

	helmBinaryPath := path.Join(path.Join(userPath, "bin"), helmVal)
	if _, statErr := os.Stat(helmBinaryPath); statErr != nil {
		if err := DownloadHelm(userPath, clientArch, clientOS, subdir); err != nil {
			return "", err
		}

	}
	return helmBinaryPath, nil
}

func GetHelmURL(arch, os, version string) string {
	archSuffix := "amd64"
	osSuffix := strings.ToLower(os)

	if strings.HasPrefix(arch, "armv7") {
		archSuffix = "arm"
	} else if strings.HasPrefix(arch, "aarch64") {
		archSuffix = "arm64"
	}
	if strings.Contains(strings.ToLower(os), "mingw") {
		osSuffix = "windows"
	}

	return fmt.Sprintf("https://get.helm.sh/helm-%s-%s-%s.tar.gz", version, osSuffix, archSuffix)
}

func DownloadHelm(userPath, clientArch, clientOS, subdir string) error {
	tools := get.MakeTools()
	var tool *get.Tool
	for _, t := range tools {
		if t.Name == "helm" {
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
		outPath, finalName, err := get.Download(tool,
			clientArch,
			clientOS,
			tool.Version,
			defaultMovePath,
			progress,
			quiet)
		if err != nil {
			return err
		}

		fmt.Println("Downloaded to: ", outPath, finalName)
	} else {
		fmt.Printf("%s already exists, skipping download.\n", tool.Name)
	}

	return nil
}

func HelmInit() error {
	fmt.Printf("Running helm init.\n")
	subdir := ""

	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s", env.LocalBinary("helm", subdir)),
		Env:         os.Environ(),
		Args:        []string{"init", "--client-only"},
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func UpdateHelmRepos(helm3 bool) error {
	subdir := ""

	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s repo update", env.LocalBinary("helm", subdir)),
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func AddHelmRepo(name, url string, update bool) error {
	subdir := ""

	if index := strings.Index(name, "/"); index > -1 {
		name = name[:index]
	}

	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s repo add %s %s", env.LocalBinary("helm", subdir), name, url),
		Env:         os.Environ(),
		StreamStdio: true,
	}
	res, err := task.Execute()

	if err != nil {
		return err
	}

	println(res.Stderr)

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	if update {
		task := execute.ExecTask{
			Command:     fmt.Sprintf("%s repo update", env.LocalBinary("helm", subdir)),
			Env:         os.Environ(),
			StreamStdio: true,
		}
		res, err := task.Execute()

		if err != nil {
			return err
		}

		println(res.Stderr)

		if res.ExitCode != 0 {
			return fmt.Errorf("exit code %d", res.ExitCode)
		}
	}

	return nil
}

func FetchChart(chart, version string) error {
	chartsPath := path.Join(os.TempDir(), "charts")
	versionStr := ""

	if len(version) > 0 {
		// Issue in helm where adding a space to the command makes it think that it's another chart of " " we want to template,
		// So we add the space before version here rather than on the command
		versionStr = " --version " + version
	}
	subdir := ""

	// First remove any existing folder
	os.RemoveAll(chartsPath)

	mkErr := os.MkdirAll(chartsPath, 0700)

	if mkErr != nil {
		return mkErr
	}
	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s fetch %s --untar=true --untardir %s%s", env.LocalBinary("helm", subdir), chart, chartsPath, versionStr),
		Env:         os.Environ(),
		StreamStdio: true,
	}
	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func Helm3Upgrade(chart, namespace, values, version string, overrides map[string]string, wait bool) error {

	chartName := chart
	if index := strings.Index(chartName, "/"); index > -1 {
		chartName = chartName[index+1:]
	}

	basePath := path.Join(os.TempDir(), "charts", chartName)

	args := []string{"upgrade", "--install", chartName, chart, "--namespace", namespace}
	if len(version) > 0 {
		args = append(args, "--version", version)
	}

	if wait {
		args = append(args, "--wait")
	}

	fmt.Println("VALUES", values)
	if len(values) > 0 {
		args = append(args, "--values")
		if !strings.HasPrefix(values, "/") {
			args = append(args, path.Join(basePath, values))
		} else {
			args = append(args, values)
		}
	}

	for k, v := range overrides {
		args = append(args, "--set")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	task := execute.ExecTask{
		Command:     env.LocalBinary("helm", ""),
		Args:        args,
		Env:         os.Environ(),
		Cwd:         basePath,
		StreamStdio: true,
	}

	fmt.Printf("Command: %s %s\n", task.Command, task.Args)
	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}

	if len(res.Stderr) > 0 {
		log.Printf("stderr: %s\n", res.Stderr)
	}

	return nil
}
