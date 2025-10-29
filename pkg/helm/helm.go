// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package helm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	execute "github.com/alexellis/go-execute/v2"
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

		verify := false
		defaultMovePath := ""
		outPath, finalName, err := get.Download(tool,
			clientArch,
			clientOS,
			tool.Version,
			defaultMovePath,
			progress,
			quiet,
			verify)
		if err != nil {
			return err
		}

		fmt.Println("Downloaded to: ", outPath, finalName)
	} else {
		fmt.Printf("%s already exists, skipping download.\n", tool.Name)
	}

	return nil
}

func UpdateHelmRepos(helm3 bool) error {
	subdir := ""

	task := execute.ExecTask{
		Command:     env.LocalBinary("helm", subdir),
		Args:        []string{"repo", "update"},
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute(context.Background())

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
		Command:     env.LocalBinary("helm", subdir),
		Args:        []string{"repo", "add", name, url, "--force-update"},
		Env:         os.Environ(),
		StreamStdio: true,
	}
	res, err := task.Execute(context.Background())
	if err != nil {
		return err
	}

	println(res.Stderr)

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	if update {
		task := execute.ExecTask{
			Command:     env.LocalBinary("helm", subdir),
			Args:        []string{"repo", "update"},
			Env:         os.Environ(),
			StreamStdio: true,
		}
		res, err := task.Execute(context.Background())

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

	subdir := ""

	// First remove any existing folder
	os.RemoveAll(chartsPath)

	if err := os.MkdirAll(chartsPath, 0700); err != nil {
		return err
	}

	task := execute.ExecTask{
		Command:     env.LocalBinary("helm", subdir),
		Args:        []string{"fetch", chart, "--untar=true", "--untardir", chartsPath},
		Env:         os.Environ(),
		StreamStdio: true,
	}

	if len(version) > 0 {
		task.Args = append(task.Args, "--version", version)
	}

	res, err := task.Execute(context.Background())

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}
	return nil
}

func IsOCI(chart string) bool {
	return strings.HasPrefix(chart, "oci://")
}

// processValuesFiles handles URL-based values files by ensuring they are passed through unchanged
// since Helm supports URLs in --values flags
func processValuesFiles(valuesFiles []string) ([]string, error) {
	// No processing needed - Helm handles URLs directly
	return valuesFiles, nil
}

// isURL checks if a string is a URL
func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func Helm3Upgrade(chart, namespace string, valuesFile []string, version string, overrides map[string]string, wait bool) error {

	chartName := chart
	if index := strings.Index(chartName, "/"); index > -1 {
		chartName = chartName[index+1:]
	}

	basePath := path.Join(os.TempDir(), "charts", chartName)

	// Process values files to download any URLs
	processedValuesFiles, err := processValuesFiles(valuesFile)
	if err != nil {
		return err
	}

	args := []string{"upgrade", "--install", chartName, chart, "--namespace", namespace}
	if len(version) > 0 {
		args = append(args, "--version", version)
	}

	if wait {
		args = append(args, "--wait")
	}

	for _, valueFile := range processedValuesFiles {
		args = append(args, "--values")
		if isURL(valueFile) {
			// Helm supports URLs directly, pass through unchanged
			args = append(args, valueFile)
		} else if !strings.HasPrefix(valueFile, "/") {
			args = append(args, path.Join(basePath, valueFile))
		} else {
			args = append(args, valueFile)
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
	res, err := task.Execute(context.Background())

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

func Helm3OCIUpgrade(chart, namespace string, valuesFiles []string, version string, overrides map[string]string, wait bool) error {

	if !IsOCI(chart) {
		return fmt.Errorf("chart %s is not an OCI chart URL", chart)
	}

	index := strings.LastIndex(chart, "/")
	if index < 0 {
		return fmt.Errorf("chart %s is not an OCI chart URL", chart)
	}
	name := chart[index+1:]

	// Process values files to download any URLs
	processedValuesFiles, err := processValuesFiles(valuesFiles)
	if err != nil {
		return err
	}

	args := []string{"upgrade", "--install", name, chart, "--namespace", namespace}
	if len(version) > 0 {
		args = append(args, "--version", version)
	}

	if wait {
		args = append(args, "--wait")
	}

	for _, valueFile := range processedValuesFiles {
		args = append(args, "--values")
		args = append(args, valueFile)
	}

	for k, v := range overrides {
		args = append(args, "--set")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	task := execute.ExecTask{
		Command:     env.LocalBinary("helm", ""),
		Args:        args,
		Env:         os.Environ(),
		StreamStdio: true,
	}

	fmt.Printf("Command: %s %s\n", task.Command, task.Args)
	res, err := task.Execute(context.TODO())

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
