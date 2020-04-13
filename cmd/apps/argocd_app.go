// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"strings"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallArgoCD() *cobra.Command {
	var command = &cobra.Command{
		Use:          "argocd",
		Short:        "Install argocd",
		Long:         `Install argocd`,
		Example:      `  arkade install argocd`,
		SilenceUsage: true,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		_, err := kubectlTask("create", "ns",
			"argocd")
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				return err
			}
		}

		_, err = kubectlTask("apply", "-f",
			"https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml", "-n", "argocd")
		if err != nil {
			return err
		}

		fmt.Println(ArgoCDInfoMsgInstallMsg)

		return nil
	}

	return command
}

const ArgoCDInfoMsg = `
# Get the ArgoCD CLI

brew tap argoproj/tap
brew install argoproj/tap/argocd

# Or download via https://github.com/argoproj/argo-cd/releases/latest

# Username is "admin", get the password

kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2

# Port-forward

kubectl port-forward svc/argocd-server -n argocd 8081:443 &

http://localhost:8081

# Get started with ArgoCD at
# https://argoproj.github.io/argo-cd/#quick-start`

const ArgoCDInfoMsgInstallMsg = `=======================================================================
= ArgoCD has been installed                                           =
=======================================================================` +
	"\n\n" + ArgoCDInfoMsg + "\n\n" + pkg.ThanksForUsing
