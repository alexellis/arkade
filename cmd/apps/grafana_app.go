package apps

import (
	"fmt"
	"log"
	"os"
	"path"

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

	grafana.Flags().StringP("namespace", "n", "grafana", "The namespace to install grafana")
	grafana.Flags().Bool("update-repo", true, "Update the helm repo")
	grafana.Flags().Bool("persistence", false, "Make grafana persistent")

	grafana.RunE = func(command *cobra.Command, args []string) error {

		const chartVersion = "5.0.4"

		// Get all flags
		wait, _ := command.Flags().GetBool("wait")
		persistence, _ := command.Flags().GetBool("persistence")
		namespace, _ := command.Flags().GetString("namespace")

		// initialize client env
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		log.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, true)
		if err != nil {
			return err
		}

		// Update chart repositories
		updateRepo, _ := grafana.Flags().GetBool("update-repo")

		if updateRepo {
			err = updateHelmRepos(true)
			if err != nil {
				return err
			}
		}

		// create the namespace
		nsRes, nsErr := kubectlTask("create", "namespace", namespace)
		if nsErr != nil {
			return nsErr
		}

		// ignore errors
		if nsRes.ExitCode != 0 {
			log.Printf("[Warning] unable to create namespace %s, may already exist: %s", namespace, nsRes.Stderr)
		}

		// download the chart
		chartPath := path.Join(os.TempDir(), "charts")
		err = fetchChart(chartPath, "stable/grafana", chartVersion, true)
		if err != nil {
			return err
		}

		outputPath := path.Join(chartPath, "grafana")

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
		err = helm3Upgrade(outputPath, "stable/grafana", namespace,
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

const grafanaInstallMsg = `=======================================================================
=                      grafana has been installed                     =
=======================================================================

# Get the admin password:

  kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

# Expose the service via port-forward:

  kubectl --namespace grafana port-forward service/grafana 3000:80

# Enable persistence:

  arkade install grafana --persistence

` +
	"\n\n" + pkg.ThanksForUsing
