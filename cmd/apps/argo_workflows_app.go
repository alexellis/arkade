// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"strings"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallArgoWorkflows() *cobra.Command {
	var command = &cobra.Command{
		Use:          "argo-workflows",
		Short:        "Install argo-workflows",
		Long:         `Install argo-workflows`,
		Example:      `  arkade install argo-workflows`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "stable", "The version of argo-workflows to install")

	command.RunE = func(command *cobra.Command, args []string) error {
		version, _ := command.Flags().GetString("version")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		url := "https://raw.githubusercontent.com/argoproj/argo-workflows/" + version + "/manifests/install.yaml"
		
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		_, err := k8s.KubectlTask("create", "ns",
			"argo")
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				return err
			}
		}

		_, err = k8s.KubectlTask("apply", "-f", url, "-n", "argo")
		if err != nil {
			return err
		}

		fmt.Println(ArgoWorkflowsInfoMsgInstallMsg)

		return nil
	}

	return command
}

const ArgoWorkflowsInfoMsg = `
# Install the argo CLI
arkade get argo-workflows

# Install argo to your kubernetes cluster
arkade install argo-workflows

# Install a specific version of argo
arkade install argo-workflows --version=v2.11.8

# Port-forward the Argo API server
kubectl -n argo port-forward deployment/argo-server 2746:2746

# Open the UI:
http://localhost:2746

# Get started with Argo Workflows at
# https://argoproj.github.io/argo-workflows/quick-start/`

const ArgoWorkflowsInfoMsgInstallMsg = `=======================================================================
= Argo has been installed                                           =
=======================================================================` +
	"\n\n" + ArgoWorkflowsInfoMsg + "\n\n" + pkg.ThanksForUsing
