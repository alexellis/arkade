// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/spf13/cobra"
)

func MakeInstallKubeStateMetrics() *cobra.Command {
	var kubeStateMetrics = &cobra.Command{
		Use:          "kube-state-metrics",
		Short:        "Install kube-state-metrics",
		Long:         `Install kube-state-metrics to generate and expose cluster-level metrics.`,
		Example:      `  arkade install kube-state-metrics --namespace default  --set replicas=2`,
		SilenceUsage: true,
	}

	kubeStateMetrics.Flags().StringP("namespace", "n", "kube-system", "The namespace used for installation")
	kubeStateMetrics.Flags().StringArray("set", []string{}, "Set individual values in the helm chart")
	kubeStateMetrics.Flags().Bool("update-repo", true, "Update the helm repo")

	kubeStateMetrics.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		namespace, _ := command.Flags().GetString("namespace")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		clientArch, clientOS := env.GetClientArch()
		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)

		updateRepo, _ := kubeStateMetrics.Flags().GetBool("update-repo")

		overrides := map[string]string{}
		setVals, err := kubeStateMetrics.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		if err := config.MergeFlags(overrides, setVals); err != nil {
			return err
		}

		kubeStateMetricsOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("prometheus-community/kube-state-metrics").
			WithHelmURL("https://prometheus-community.github.io/helm-charts").
			WithHelmUpdateRepo(updateRepo).
			WithOverrides(overrides).
			WithWait(wait).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(kubeStateMetricsOptions)
		if err != nil {
			return err
		}

		println(`=======================================================================
=             kube-state-metrics has been installed.                  =
=======================================================================

# Port-forward
kubectl port-forward -n ` + namespace + ` service/kube-state-metrics 9000:8080 &

# Then access via:
http://localhost:9000/metrics
` + KubeStateMetricsInfoMsg + `
` + pkg.SupportMessageShort)

		return nil
	}

	return kubeStateMetrics
}

const KubeStateMetricsInfoMsg = `
# Find out more at:
# https://github.com/kubernetes/kube-state-metrics
`
