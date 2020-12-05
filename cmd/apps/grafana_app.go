// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/k8s"
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

		const chartVersion = "5.0.4"

		// Get all flags
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		wait, _ := command.Flags().GetBool("wait")
		persistence, _ := command.Flags().GetBool("persistence")
		namespace, _ := command.Flags().GetString("namespace")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}
		// initialize client env
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		log.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		// create the namespace
		nsRes, nsErr := k8s.KubectlTask("create", "namespace", namespace)
		if nsErr != nil {
			return nsErr
		}

		// ignore errors
		if nsRes.ExitCode != 0 {
			log.Printf("[Warning] unable to create namespace %s, may already exist: %s", namespace, nsRes.Stderr)
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

		// set custom flags
		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		grafanaAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("grafana/grafana").
			WithHelmURL("https://grafana.github.io/helm-charts/").
			WithHelmRepoVersion(chartVersion).
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(grafanaAppOptions)
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
