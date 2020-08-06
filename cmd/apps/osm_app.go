// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
)

func MakeInstallOSM() *cobra.Command {
	var osm = &cobra.Command{
		Use:          "osm",
		Short:        "Install osm",
		Long:         `Install Open Service Mesh (OSM) - a lightweight, extensible, cloud native service mesh`,
		Example:      `  arkade install osm`,
		SilenceUsage: true,
	}

	osm.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)
		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		userPath, err := getUserPath()
		if err != nil {
			return err
		}

		arch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q\n", clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		err = downloadOSM(userPath, arch, clientOS)
		if err != nil {
			return err
		}
		fmt.Println("Running osm check, this may take a few moments.")

		_, err = osmCli("check", "--pre-flight")
		if err != nil {
			return err
		}

		res, err := osmCli("install")
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf("exit code %d, error: %s", res.ExitCode, res.Stderr)
		}

		_, err = osmCli("check")
		if err != nil {
			return err
		}

		fmt.Println(`=======================================================================
= OSM has been installed.                                        =
=======================================================================

# Get the osm-cli
curl -sL https://run.osm.io/install | sh

# Find out more at:
# https://osm.io

# To use the OSM CLI set this path:

export PATH=$PATH:` + path.Join(userPath, "bin/") + `
osm --help

` + pkg.ThanksForUsing)
		return nil
	}

	return osm
}

func downloadOSM(userPath, arch, clientOS string) error {
	t := &get.Tool{
		Name:    "osm",
		Repo:    "osm",
		Owner:   "openservicemesh",
		Version: "v0.1.0",
		URLTemplate: `
{{$osStr := ""}}
{{ if HasPrefix .OS "ming" -}}
{{$osStr = "windows"}}
{{- else if eq .OS "Linux" -}}
{{$osStr = "linux"}}
{{- else if eq .OS "Darwin" -}}
{{$osStr = "darwin"}}
{{- end -}}
https://github.com/openservicemesh/osm/releases/download/{{.Version}}/osm-{{.Version}}-{{$osStr}}-amd64.tar.gz`,
	}

	u, err := get.GetDownloadURL(t, clientOS, arch, t.Version)
	if err != nil {
		return err
	}

	fmt.Printf("Download URL: %s\n", u)

	return nil
}

func osmCli(parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     fmt.Sprintf("%s", env.LocalBinary("osm", "")),
		Args:        parts,
		Env:         os.Environ(),
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return res, err
	}

	if res.ExitCode != 0 {
		return res, fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}

	return res, nil
}

var OSMInfoMsg = `# The osm CLI is installed at:
# $HOME/.bin/arkade

# Find out more at:
# https://github.com/openservicemesh/osm

# To use the OSM CLI set this path:

export PATH=$PATH:` + getExportPath() + `
osm --help`

var osmInstallMsg = `=======================================================================
= Open Service Mesh (OSM) has been installed.                                         =
=======================================================================` +
	"\n\n" + OSMInfoMsg + "\n\n" + pkg.ThanksForUsing
