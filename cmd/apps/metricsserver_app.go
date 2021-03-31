// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/spf13/cobra"
)

func MakeInstallMetricsServer() *cobra.Command {
	var metricsServer = &cobra.Command{
		Use:          "metrics-server",
		Short:        "Install metrics-server",
		Long:         `Install metrics-server to provide metrics on nodes and Pods in your cluster.`,
		Example:      `  arkade install metrics-server --namespace kube-system`,
		SilenceUsage: true,
	}

	metricsServer.Flags().StringP("namespace", "n", "kube-system", "The namespace used for installation")
	metricsServer.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	metricsServer.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := command.Flags().GetString("namespace")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

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
		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		nfsProvisionerOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("stable/metrics-server").
			WithHelmURL("https://charts.helm.sh/stable").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(nfsProvisionerOptions)
		if err != nil {
			return err
		}

		println(MetricsInfoMsg)
		return nil
	}

	return metricsServer
}

const MetricsInfoMsg = `

You have installed the metrics-server for Kubernetes:

# Check pod usage
kubectl top pod

# Check node usage
kubectl top node

# Find out more at:
# https://github.com/helm/charts/tree/master/stable/metrics-server`
