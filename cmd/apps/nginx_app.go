// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallNginx() *cobra.Command {
	var nginx = &cobra.Command{
		Use:     "ingress-nginx",
		Aliases: []string{"nginx-ingress"}, // backward compatibility
		Short:   "Install ingress-nginx",
		Long: `Install ingress-nginx. This app can be installed with Host networking for
cases where an external LB is not available. please see the --host-mode
flag and the ingress-nginx docs for more info`,
		Example:      `  arkade install ingress-nginx --namespace default`,
		SilenceUsage: true,
	}

	nginx.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	nginx.Flags().Bool("update-repo", true, "Update the helm repo")
	nginx.Flags().Bool("host-mode", false, "If we should install ingress-nginx in host mode.")
	nginx.Flags().Bool("default-ingress", true, "Is this the default ingressClass for the cluster?")
	nginx.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	nginx.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		wait, _ := command.Flags().GetBool("wait")

		namespace, _ := command.Flags().GetString("namespace")

		updateRepo, _ := nginx.Flags().GetBool("update-repo")

		overrides := map[string]string{}

		hostMode, flagErr := command.Flags().GetBool("host-mode")
		if flagErr != nil {
			return flagErr
		}
		if hostMode {
			fmt.Println("Running in host networking mode")
			overrides["controller.hostNetwork"] = "true"
			overrides["controller.hostPort.enabled"] = "true"
			overrides["controller.service.type"] = "NodePort"
			overrides["dnsPolicy"] = "ClusterFirstWithHostNet"
			overrides["controller.kind"] = "DaemonSet"
		}

		defaultIngress, flagErr := command.Flags().GetBool("default-ingress")
		if flagErr != nil {
			return flagErr
		}
		if defaultIngress {
			overrides["controller.ingressClassResource.default"] = "true"
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		nginxOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("ingress-nginx/ingress-nginx").
			WithHelmURL("https://kubernetes.github.io/ingress-nginx").
			WithHelmUpdateRepo(updateRepo).
			WithOverrides(overrides).
			WithWait(wait).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(nginxOptions)

		if err != nil {
			return err
		}

		fmt.Println(nginxIngressInstallMsg)

		return nil
	}

	return nginx
}

const NginxIngressInfoMsg = `# If you're using a local environment such as "minikube" or "KinD",
# then try the inlets operator with "arkade install inlets-operator"

# If you're using a managed Kubernetes service, then you'll find
# your LoadBalancer's IP under "EXTERNAL-IP" via:

kubectl get svc ingress-nginx-controller

# Find out more at:
# https://github.com/kubernetes/ingress-nginx/tree/master/charts/ingress-nginx`

const nginxIngressInstallMsg = `=======================================================================
= ingress-nginx has been installed.                                   =
=======================================================================` +
	"\n\n" + NginxIngressInfoMsg + "\n\n" + pkg.SupportMessageShort
