// Copyright (c) arkade author(s) 2022. All rights reserved.
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
	metricsServer.Flags().StringP("tag", "t", "v0.6.3", "The tag or version of the metrics-server to install")
	metricsServer.Flags().Bool("update-repo", true, "Update the helm repo")

	metricsServer.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := command.Flags().GetString("namespace")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		updateRepo, _ := metricsServer.Flags().GetBool("update-repo")

		overrides := map[string]string{}
		overrides["args"] = `{--kubelet-insecure-tls,--kubelet-preferred-address-types=InternalIP\,ExternalIP\,Hostname}`

		tag, _ := command.Flags().GetString("tag")

		overrides["image.tag"] = tag

		customFlags, _ := command.Flags().GetStringArray("set")
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		nfsProvisionerOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("metrics-server/metrics-server").
			WithHelmURL("https://kubernetes-sigs.github.io/metrics-server").
			WithHelmUpdateRepo(updateRepo).
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
# https://artifacthub.io/packages/helm/metrics-server/metrics-server
`
