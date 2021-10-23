// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallPrometheusOperator() *cobra.Command {
	var kubePromStackApp = &cobra.Command{
		Use:          "kube-prometheus-stack",
		Short:        "Install Kube Prometheus Stack for monitoring",
		Long:         "Install Kube Prometheus Stack, provides Kubernetes native deployment and management of Prometheus and related monitoring components.",
		Example:      "arkade install kube-prometheus-stack",
		SilenceUsage: true,
	}

	kubePromStackApp.Flags().StringP("namespace", "n", "default", "The namespace to install prometheus-operator (default: default")
	kubePromStackApp.Flags().Bool("update-repo", true, "Update the helm repo")
	kubePromStackApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set  prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false)")
	kubePromStackApp.Flags().Bool("alertmanager", true, "Install AlertManager (default: true)")
	kubePromStackApp.Flags().Bool("grafana", false, "Install Grafana (default: false)")
	kubePromStackApp.Flags().Bool("node-exporter", true, "Install Node Exporter (default: true)")
	kubePromStackApp.Flags().Bool("kube-state-metrics", true, "Install Kube State Metrics (default: true)")
	kubePromStackApp.Flags().Bool("prometheus-operator", true, "Install Prometheus Operator (default: true)")
	kubePromStackApp.Flags().Bool("prometheus", true, "Install Prometheus instance (default: true)")

	kubePromStackApp.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("error with --namespace usage: %s", err)
		}

		_, err = command.Flags().GetBool("alertmanager")
		if err != nil {
			return fmt.Errorf("error with --alertmanager usage: %s", err)
		}

		_, err = command.Flags().GetBool("grafana")
		if err != nil {
			return fmt.Errorf("error with --grafana usage: %s", err)
		}

		_, err = command.Flags().GetBool("node-exporter")
		if err != nil {
			return fmt.Errorf("error with --node-exporter usage: %s", err)
		}

		_, err = command.Flags().GetBool("kube-state-metrics")
		if err != nil {
			return fmt.Errorf("error with --kube-state-metrics usage: %s", err)
		}

		_, err = command.Flags().GetBool("prometheus-operator")
		if err != nil {
			return fmt.Errorf("error with --prometheus-operator usage: %s", err)
		}

		_, err = command.Flags().GetBool("prometheus")
		if err != nil {
			return fmt.Errorf("error with --prometheus usage: %s", err)
		}

		_, err = command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		return nil
	}

	kubePromStackApp.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		namespace, _ := kubePromStackApp.Flags().GetString("namespace")
		installAlertManager, _ := kubePromStackApp.Flags().GetBool("alertmanager")
		installGrafana, _ := kubePromStackApp.Flags().GetBool("grafana")
		installNodeExporter, _ := kubePromStackApp.Flags().GetBool("node-exporter")
		installKubeStateMetrics, _ := kubePromStackApp.Flags().GetBool("kube-state-metrics")
		installprometheusOperator, _ := kubePromStackApp.Flags().GetBool("prometheus-operator")
		installPrometheus, _ := kubePromStackApp.Flags().GetBool("prometheus")

		overrides := map[string]string{}

		if !installAlertManager {
			overrides["alertmanager.enabled"] = "false"
		}

		if !installGrafana {
			overrides["grafana.enabled"] = "false"
		}

		if !installNodeExporter {
			overrides["nodeExporter.enabled"] = "false"
		}

		if !installprometheusOperator {
			overrides["prometheusOperator.enabled"] = "false"
		}

		if !installKubeStateMetrics {
			overrides["kubeStateMetrics.enabled"] = "false"
		}

		if !installPrometheus {
			overrides["prometheus.enabled"] = "false"
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kubePromStackOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("prometheus-community/kube-prometheus-stack").
			WithHelmURL("https://prometheus-community.github.io/helm-charts").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(kubePromStackOptions)
		if err != nil {
			return err
		}

		println(kubePromInstallMsg)
		return nil
	}

	return kubePromStackApp
}

const KubePromInfoMsg = `# Get started with kube-prometheus-stack here:
# https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-prometheus-stack/README.md

Visit https://github.com/prometheus-operator/kube-prometheus for instructions on how to create & configure Alertmanager and Prometheus instances using the Operator.

 # Forward traffic to your localhost for prometheus
 kubectl port-forward service/prometheus-operated 9090:9090

 # Forward traffic to your localhost for alertmanager (if installed)
 kubectl port-forward service/alertmanager-operated 9093:9093

`

const kubePromInstallMsg = `=======================================================================
= kube-prometheus-stack has been installed.                           =
=======================================================================` +
	"\n\n" + KubePromInfoMsg + "\n\n" + pkg.SupportMessageShort
