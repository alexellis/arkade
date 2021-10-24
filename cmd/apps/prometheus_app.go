// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallPrometheus() *cobra.Command {
	var kubePrometheusApp = &cobra.Command{
		Use:          "prometheus",
		Short:        "Install Prometheus for monitoring",
		Long:         "Install Prometheus, provides Kubernetes native deployment and management of Prometheus and related monitoring components.",
		Example:      "arkade install prometheus",
		SilenceUsage: true,
	}

	kubePrometheusApp.Flags().StringP("namespace", "n", "default", "The namespace to install prometheus (default: default")
	kubePrometheusApp.Flags().Bool("update-repo", true, "Update the helm repo")
	kubePrometheusApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set grafana.enabled=true)")
	kubePrometheusApp.Flags().Bool("alertmanager", true, "Install AlertManager (default: true)")
	kubePrometheusApp.Flags().Bool("node-exporter", true, "Install Node Exporter (default: true)")
	kubePrometheusApp.Flags().Bool("kube-state-metrics", true, "Install Kube State Metrics (default: true)")
	kubePrometheusApp.Flags().Bool("pushgateway", true, "Install Push Gateway (default: true)")
	kubePrometheusApp.Flags().Bool("prometheus", true, "Install Prometheus instance (default: true)")

	kubePrometheusApp.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		log.Println(kubeConfigPath)
		namespace, _ := kubePrometheusApp.Flags().GetString("namespace")
		installAlertManager, _ := kubePrometheusApp.Flags().GetBool("alertmanager")
		installNodeExporter, _ := kubePrometheusApp.Flags().GetBool("node-exporter")
		installKubeStateMetrics, _ := kubePrometheusApp.Flags().GetBool("kube-state-metrics")
		installPushGateway, _ := kubePrometheusApp.Flags().GetBool("pushgateway")
		installPrometheus, _ := kubePrometheusApp.Flags().GetBool("prometheus")

		overrides := map[string]string{}

		if !installAlertManager {
			overrides["alertmanager.enabled"] = "false"
		}

		if !installNodeExporter {
			overrides["nodeExporter.enabled"] = "false"
		}

		if !installPushGateway {
			overrides["pushgateway.enabled"] = "false"
		}

		if !installKubeStateMetrics {
			overrides["kubeStateMetrics.enabled"] = "false"
		}

		if !installPrometheus {
			overrides["server.enabled"] = "false"
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kubePromStackOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("prometheus-community/prometheus").
			WithHelmURL("https://prometheus-community.github.io/helm-charts").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(kubePromStackOptions)
		if err != nil {
			return err
		}

		println(PrometheusInstallMsg)
		return nil
	}

	return kubePrometheusApp
}

const PrometheusInfoMsg = `# Get started with Prometheus here:
# https://github.com/prometheus-community/helm-charts/blob/main/charts/prometheus/README.md

 # Forward traffic to your localhost for prometheus
 kubectl port-forward service/prometheus-server 8080:80

`

const PrometheusInstallMsg = `=======================================================================
= prometheus has been installed.                                      =
=======================================================================` +
	"\n\n" + PrometheusInfoMsg + "\n\n" + pkg.ThanksForUsing
