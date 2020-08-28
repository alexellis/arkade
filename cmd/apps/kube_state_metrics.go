// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
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

	kubeStateMetrics.RunE = func(command *cobra.Command, args []string) error {
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

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		clientArch, clientOS := env.GetClientArch()
		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)
		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		err = helm.UpdateHelmRepos(true)
		if err != nil {
			return err
		}

		chartPath := path.Join(os.TempDir(), "charts")
		err = helm.FetchChart("stable/kube-state-metrics", defaultVersion)

		if err != nil {
			return err
		}

		setMap := map[string]string{}
		setVals, _ := kubeStateMetrics.Flags().GetStringArray("set")

		for _, setV := range setVals {
			var k string
			var v string

			if index := strings.Index(setV, "="); index > -1 {
				k = setV[:index]
				v = setV[index+1:]
				setMap[k] = v
			}
		}

		fmt.Println("Chart path: ", chartPath)

		err = helm.Helm3Upgrade("stable/kube-state-metrics", namespace,
			"values.yaml",
			defaultVersion,
			setMap,
			wait)

		if err != nil {
			return err
		}

		fmt.Println(`=======================================================================
=             kube-state-metrics has been installed.                  =
=======================================================================

# Port-forward
kubectl port-forward -n ` + namespace + ` service/kube-state-metrics 9000:8080 &

# Then access via:
http://localhost:9000/metrics
` + KubeStateMetricsInfoMsg + `
` + pkg.ThanksForUsing)

		return nil
	}

	return kubeStateMetrics
}

const KubeStateMetricsInfoMsg = `
# Find out more at:
# https://github.com/kubernetes/kube-state-metrics
`
