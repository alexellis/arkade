// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"errors"
	"fmt"
	"log"
	"os"

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
		Use:   "osm",
		Short: "Install osm",
		Long: `Install Open Service Mesh (OSM) - a lightweight, extensible, cloud native 
service mesh created by Microsoft Azure.`,
		Example:      `  arkade install osm`,
		SilenceUsage: true,
	}

	osm.Flags().String("mesh", "osm", "Give a specific mesh name override")
	osm.Flags().StringP("namespace", "n", "osm-system", "Give a specific mesh namespace override")

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
		ns, _ := osm.Flags().GetString("namespace")
		meshName, _ := osm.Flags().GetString("mesh")

		_, err = osmCli("check", "--pre-install", "--namespace="+ns)
		if err != nil {
			return err
		}

		res, err := osmCli("install", "--mesh-name="+meshName)
		if err != nil {
			return err
		}
		if res.ExitCode != 0 {
			return fmt.Errorf("exit code %d, error: %s", res.ExitCode, res.Stderr)
		}

		fmt.Println(`=======================================================================
= OSM has been installed.                                             =
=======================================================================
` +
			OSMInfoMsg + pkg.ThanksForUsing)
		return nil
	}

	return osm
}

func downloadOSM(userPath, arch, clientOS string) error {

	tools := get.MakeTools()
	var tool *get.Tool
	for _, t := range tools {
		if t.Name == "osm" {
			tool = &t
			break
		}
	}
	if tool == nil {
		return fmt.Errorf("unable to find tool definition")
	}

	if _, err := os.Stat(fmt.Sprintf("%s", env.LocalBinary(tool.Name, ""))); errors.Is(err, os.ErrNotExist) {

		outPath, finalName, err := get.Download(tool, arch, clientOS, tool.Version, get.DownloadArkadeDir)
		if err != nil {
			return err
		}

		fmt.Println("Downloaded to: ", outPath, finalName)
	} else {
		fmt.Printf("%s already exists, skipping download.\n", tool.Name)
	}
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
# $HOME/.bin/arkade/osm

# Find out more at:
# https://github.com/openservicemesh/osm

# Docs are live at:
# https://openservicemesh.io

# Walk-through a demo at:
# https://github.com/openservicemesh/osm/blob/main/docs/example/README.md

# To use the OSM CLI set this path:

export PATH=$PATH:` + getExportPath() + `
osm --help
`

var osmInstallMsg = `=======================================================================
= Open Service Mesh (OSM) has been installed.                                         =
=======================================================================` +
	"\n\n" + OSMInfoMsg + "\n\n" + pkg.ThanksForUsing
