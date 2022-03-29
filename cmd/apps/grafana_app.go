// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"
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

	grafana.Flags().StringP("namespace", "n", "grafana", "The namespace to install grafana")
	grafana.Flags().Bool("update-repo", true, "Update the helm repo")
	grafana.Flags().Bool("persistence", false, "Make grafana persistent")
	grafana.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	grafana.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetBool("wait")
		if err != nil {
			return fmt.Errorf("error with --wait usage: %s", err)
		}

		_, err = command.Flags().GetBool("persistence")
		if err != nil {
			return fmt.Errorf("error with --persistence usage: %s", err)
		}

		_, err = command.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("error with --namespace usage: %s", err)
		}

		_, err = command.Flags().GetBool("update-repo")
		if err != nil {
			return fmt.Errorf("error with --update-repo usage: %s", err)
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		return nil
	}

	grafana.RunE = func(command *cobra.Command, args []string) error {

		const chartVersion = "6.24.1"

		// Get all flags
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		wait, _ := command.Flags().GetBool("wait")
		persistence, _ := command.Flags().GetBool("persistence")
		namespace, _ := command.Flags().GetString("namespace")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		customFlags, _ := command.Flags().GetStringArray("set")

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

		// set custom flags
		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		grafanaAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("grafana/grafana").
			WithHelmURL("https://grafana.github.io/helm-charts/").
			WithHelmRepoVersion(chartVersion).
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err := apps.MakeInstallChart(grafanaAppOptions)
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

# Access via http://127.0.0.1:3000

# Optionally, enable persistence if required:

arkade install grafana --persistence
`

var grafanaInstallMsg = `=======================================================================
=                      grafana has been installed                     =
=======================================================================` +
	"\n\n" + GrafanaInfoMsg + "\n\n" + pkg.ThanksForUsing
