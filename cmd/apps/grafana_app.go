// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/commands"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallGrafana() *cobra.Command {
	var grafana = &cobra.Command{
		Use:          "grafana",
		Short:        "Install grafana",
		Long:         "Install grafana for creating dashboards",
		Example:      "arkade install grafana",
		SilenceUsage: true,
	}

	grafana.Flags().Bool("persistence", false, "Make grafana persistent")

	grafana.RunE = func(command *cobra.Command, args []string) error {

		const chartVersion = "5.0.4"

		// Get all flags
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}
		wait, _ := command.Flags().GetBool("wait")
		persistence, _ := command.Flags().GetBool("persistence")

		namespace, err := commands.GetNamespace(command.Flags(), "grafana")
		if err != nil {
			return err
		}
		if err := commands.CreateNamespace(namespace); err != nil {
			return err
		}

		// initialize client env
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		updateRepo, _ := grafana.Flags().GetBool("update-repo")
		err = helm.AddHelmRepo("stable", "https://kubernetes-charts.storage.googleapis.com", updateRepo)
		if err != nil {
			return err
		}

		// create the namespace

		// download the chart
		err = helm.FetchChart("stable/grafana", chartVersion)
		if err != nil {
			return err
		}

		// define the values to override
		// due the missing arm support. datasource and dashboard sidecars are not possible
		overrides := map[string]string{
			"sidecar.datasources.enabled": "false",
			"sidecar.dashboards.enabled":  "false",
		}

		if persistence {
			overrides["persistence.enabled"] = "true"
			overrides["persistence.size"] = "2Gi"
		}

		// install the chart
		err = helm.Helm3Upgrade("stable/grafana", namespace,
			"values.yaml",
			chartVersion,
			overrides,
			wait)

		if err != nil {
			return err
		}

		fmt.Println(grafanaInstallMsg)

		return nil
	}

	return grafana
}

const GrafanaInfoMsg = `
# Get the admin password:

  kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

# Expose the service via port-forward:

  kubectl --namespace grafana port-forward service/grafana 3000:80

# Enable persistence:

  arkade install grafana --persistence

`

var grafanaInstallMsg = `=======================================================================
=                      grafana has been installed                     =
=======================================================================` +
	"\n\n" + GrafanaInfoMsg + "\n\n" + pkg.ThanksForUsing
