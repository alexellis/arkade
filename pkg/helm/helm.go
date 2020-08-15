// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package helm

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	execute "github.com/alexellis/go-execute/pkg/v1"
)

const helmVersion = "v2.16.0"
const helm3Version = "v3.1.2"

func TryDownloadHelm(userPath, clientArch, clientOS string, helm3 bool) (string, error) {
	helmVal := "helm"
	subdir := ""
	if helm3 {
		helmVal = "helm3"
		subdir = "helm3"
	}

	helmBinaryPath := path.Join(path.Join(userPath, "bin"), helmVal)
	if _, statErr := os.Stat(helmBinaryPath); statErr != nil {
		if err := DownloadHelm(userPath, clientArch, clientOS, subdir, helm3); err != nil {
			return "", err
		}

		if !helm3 {
			err := HelmInit()
			if err != nil {
				return "", err
			}
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

func DownloadHelm(userPath, clientArch, clientOS, subdir string, helm3 bool) error {

	useHelmVersion := helmVersion
	if helm3 {
		useHelmVersion = helm3Version
	}

	helmURL := GetHelmURL(clientArch, clientOS, useHelmVersion)
	fmt.Println(helmURL)
	parsedURL, _ := url.Parse(helmURL)

	res, err := http.DefaultClient.Get(parsedURL.String())
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("download of %s gave status: %d", parsedURL.String(), res.StatusCode)
	}

	dest := path.Join(path.Join(userPath, "bin"), subdir)
	mkErr := os.MkdirAll(dest, 0700)
	if mkErr != nil {
		return mkErr
	}

	defer res.Body.Close()
	r := ioutil.NopCloser(res.Body)
	untarErr := archive.Untar(r, dest)
	if untarErr != nil {
		return untarErr
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

func AddHelmRepo(name, url string, update, helm3 bool) error {
	subdir := ""
	if helm3 {
		subdir = "helm3"
	}

	if index := strings.Index(name, "/"); index > -1 {
		name = name[:index]
	}

	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s repo add %s %s", env.LocalBinary("helm", subdir), name, url),
		Env:         os.Environ(),
		StreamStdio: true,
	}
	res, err := task.Execute()

	println(res.Stderr)

	if err != nil {
		return err
	}
	if res.ExitCode != 0 {
		return fmt.Errorf("exit code %d", res.ExitCode)
	}

	if update {
		return UpdateHelmRepos(helm3)
	}

	return nil
}

func UpdateHelmRepos(helm3 bool) error {
	subdir := ""
	if helm3 {
		subdir = "helm3"
	}
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

func FetchChart(chart, version string, helm3 bool) error {
	chartsPath := path.Join(os.TempDir(), "charts")
	versionStr := ""

	if len(version) > 0 {
		// Issue in helm where adding a space to the command makes it think that it's another chart of " " we want to template,
		// So we add the space before version here rather than on the command
		versionStr = " --version " + version
	}
	subdir := ""
	if helm3 {
		subdir = "helm3"
	}

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
		Command:     env.LocalBinary("helm", "helm3"),
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

func TemplateChart(basePath, chart, namespace, outputPath, values string, overrides map[string]string) error {

	rmErr := os.RemoveAll(outputPath)

	if rmErr != nil {
		log.Printf("Error cleaning up: %s, %s\n", outputPath, rmErr.Error())
	}

	mkErr := os.MkdirAll(outputPath, 0700)
	if mkErr != nil {
		return mkErr
	}

	overridesStr := ""
	for k, v := range overrides {
		overridesStr += fmt.Sprintf(" --set %s=%s", k, v)
	}

	chartRoot := path.Join(basePath, chart)

	valuesStr := ""
	if len(values) > 0 {
		valuesStr = "--values " + path.Join(chartRoot, values)
	}

	task := execute.ExecTask{
		Command: fmt.Sprintf("%s template %s --name %s --namespace %s --output-dir %s %s %s",
			env.LocalBinary("helm", ""), chart, chart, namespace, outputPath, valuesStr, overridesStr),
		Env:         os.Environ(),
		Cwd:         basePath,
		StreamStdio: true,
	}

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
