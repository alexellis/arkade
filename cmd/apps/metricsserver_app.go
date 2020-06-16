// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallMetricsServer() *cobra.Command {
	var metricsServer = &cobra.Command{
		Use:          "metrics-server",
		Short:        "Install metrics-server",
		Long:         `Install metrics-server to provide metrics on nodes and Pods in your cluster.`,
		Example:      `  arkade install metrics-server --namespace kube-system --helm3`,
		SilenceUsage: true,
	}

	metricsServer.Flags().StringP("namespace", "n", "kube-system", "The namespace used for installation")
	metricsServer.Flags().Bool("helm3", true, "Use helm3, if set to false uses helm2")

	metricsServer.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}
		namespace, _ := command.Flags().GetString("namespace")

		if namespace != "kube-system" {
			return fmt.Errorf(`to override the "kube-system", install via tiller`)
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		helm3, _ := command.Flags().GetBool("helm3")

		if helm3 {
			fmt.Println("Using helm3")
		}

		clientArch, clientOS := env.GetClientArch()
		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)
		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		err = helm.UpdateHelmRepos(helm3)
		if err != nil {
			return err
		}

		chartPath := path.Join(os.TempDir(), "charts")
		err = helm.FetchChart("stable/metrics-server", defaultVersion, helm3)

		if err != nil {
			return err
		}

		overrides := map[string]string{}
		overrides["args"] = `{--kubelet-insecure-tls,--kubelet-preferred-address-types=InternalIP\,ExternalIP\,Hostname}`
		switch arch {
		case "arm":
			overrides["image.repository"] = `gcr.io/google_containers/metrics-server-arm`
			break
		case "arm64", "aarch64":
			overrides["image.repository"] = `gcr.io/google_containers/metrics-server-arm64`
			break
		}

		fmt.Println("Chart path: ", chartPath)

		if helm3 {
			err := helm.Helm3Upgrade("stable/metrics-server", namespace,
				"values.yaml",
				defaultVersion,
				overrides,
				wait)

			if err != nil {
				return err
			}

		} else {
			outputPath := path.Join(chartPath, "metrics-server/rendered")

			err = helm.TemplateChart(chartPath,
				"metrics-server",
				namespace,
				outputPath,
				"values.yaml",
				overrides)

			if err != nil {
				return err
			}

			applyRes, applyErr := k8s.KubectlTask("apply", "-n", namespace, "-R", "-f", outputPath)
			if applyErr != nil {
				return applyErr
			}
			if applyRes.ExitCode > 0 {
				return fmt.Errorf("error applying templated YAML files, error: %s", applyRes.Stderr)
			}

		}

		fmt.Println(`=======================================================================
= metrics-server has been installed.                                  =
=======================================================================

# It can take a few minutes for the metrics-server to collect data
# from the cluster. Try these commands and wait a few moments if
# no data is showing.

` + MetricsInfoMsg + `

` + pkg.ThanksForUsing)

		return nil
	}

	return metricsServer
}

const MetricsInfoMsg = `# Check pod usage

kubectl top pod

# Check node usage

kubectl top node


# Find out more at:
# https://github.com/helm/charts/tree/master/stable/metrics-server`
