// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallNginx() *cobra.Command {
	var nginx = &cobra.Command{
		Use:     "ingress-nginx",
		Aliases: []string{"nginx-ingress"},
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
	nginx.Flags().Bool("helm3", true, "Use helm3, if set to false uses helm2")

	nginx.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()
		wait, _ := command.Flags().GetBool("wait")

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		updateRepo, _ := nginx.Flags().GetBool("update-repo")

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)
		helm3, _ := command.Flags().GetBool("helm3")

		if helm3 {
			fmt.Println("Using helm3")
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}
		namespace, _ := command.Flags().GetString("namespace")

		if namespace != "default" {
			return fmt.Errorf(`to override the "default", install via tiller`)
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		err = addHelmRepo("ingress-nginx", "https://kubernetes.github.io/ingress-nginx", helm3)
		if err != nil {
			return err
		}

		if updateRepo {
			err = updateHelmRepos(helm3)
			if err != nil {
				return err
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")
		err = fetchChart(chartPath, "ingress-nginx/ingress-nginx", defaultVersion, helm3)

		if err != nil {
			return err
		}

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
		fmt.Println("Chart path: ", chartPath)

		ns := "default"

		if helm3 {
			outputPath := path.Join(chartPath, "ingress-nginx")

			err := helm3Upgrade(outputPath, "ingress-nginx/ingress-nginx", ns,
				"values.yaml",
				defaultVersion,
				overrides,
				wait)

			if err != nil {
				return err
			}
		} else {
			outputPath := path.Join(chartPath, "ingress-nginx/rendered")

			err = templateChart(chartPath,
				"ingress-nginx",
				ns,
				outputPath,
				"values.yaml",
				overrides)

			if err != nil {
				return err
			}

			err = kubectl("apply", "-R", "-f", outputPath)

			if err != nil {
				return err
			}
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
	"\n\n" + NginxIngressInfoMsg + "\n\n" + pkg.ThanksForUsing
