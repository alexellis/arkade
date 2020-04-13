// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package helm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

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
		DownloadHelm(userPath, clientArch, clientOS, subdir, helm3)

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

	dest := path.Join(path.Join(userPath, "bin"), subdir)
	os.MkdirAll(dest, 0700)

	defer res.Body.Close()
	r := ioutil.NopCloser(res.Body)
	untarErr := Untar(r, dest)
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
