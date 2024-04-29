// Copyright (c) arkade author(s) 2022. All rights reserved.
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

func MakeInstallArgoCD() *cobra.Command {
	var command = &cobra.Command{
		Use:          "argocd",
		Short:        "Install argocd",
		Long:         `Install argocd`,
		Example:      `  arkade install argocd`,
		SilenceUsage: true,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		_, err := k8s.KubectlTask("create", "ns",
			"argocd")
		if err != nil {
			if !strings.Contains(err.Error(), "exists") {
				return err
			}
		}

		_, err = k8s.KubectlTask("apply", "-f",
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
# Install the "argocd" CLI:
arkade get argocd

# Port-forward the ArgoCD API server
kubectl port-forward svc/argocd-server -n argocd 8443:443 &

# Get the password
PASS=$(kubectl get secret argocd-initial-admin-secret \
  -n argocd \
  -o jsonpath="{.data.password}" | base64 -d)
echo $PASS

# Or log in:
argocd login --name local 127.0.0.1:8443 --insecure \
 --username admin \
 --password $PASS

# Open the UI:
https://127.0.0.1:8443

# Get started with ArgoCD at
# https://argoproj.github.io/argo-cd/#quick-start`

const ArgoCDInfoMsgInstallMsg = `=======================================================================
= ArgoCD has been installed                                           =
=======================================================================` +
	"\n\n" + ArgoCDInfoMsg + "\n\n" + pkg.SupportMessageShort
