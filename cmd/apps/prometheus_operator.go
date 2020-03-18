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

func MakeInstallPrometheusOperator() *cobra.Command {
	var command = &cobra.Command{
		Use:          "prometheus-operator",
		Short:        "Install prometheus operator",
		Long:         "Install prometheus operator",
		SilenceUsage: true,
	}

	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set=defaultRules.create=false)")
	command.Flags().String("namespace", "default", "Namespace for the app")

	command.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath := getDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		namespace, _ := command.Flags().GetString("namespace")

		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		helm3 := true

		// persistence, _ := command.Flags().GetBool("persistence")

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		err = addHelmRepo("stable", "https://kubernetes-charts.storage.googleapis.com/", helm3)
		if err != nil {
			return fmt.Errorf("unable to add repo %s", err)
		}

		updateRepo, _ := command.Flags().GetBool("update-repo")

		if updateRepo {
			err = updateHelmRepos(helm3)
			if err != nil {
				return fmt.Errorf("unable to update repos %s", err)
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")

		err = fetchChart(chartPath, "stable/prometheus-operator", defaultVersion, helm3)

		if err != nil {
			return fmt.Errorf("unable fetch chart %s", err)
		}

		overrides := map[string]string{}
		overrides["prometheusOperator.createCustomResource"] = "false"

		outputPath := path.Join(chartPath, "prometheus-operator")

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		err = helm3Upgrade(outputPath, "stable/prometheus-operator",
			namespace, "values.yaml", defaultVersion, overrides, wait)
		if err != nil {
			return fmt.Errorf("unable to install prometheus-operator chart with helm %s", err)
		}
		fmt.Println(prometheusOperatorInstallMsg)
		return nil
	}

	return command
}

const prometheusOperatorInstallMsg = `=======================================================================
=                  Prometheus Operator has been installed.                        =
=======================================================================` +
	"\n\n" + pkg.ThanksForUsing

var PrometheusOperatorInfoMsg = `
# Grafana can be access via port-forwarding on port 80 from within your cluster using below command:

kubectl port-forward svc/prometheus-operator-grafana 8888:80

Grafana can now be accessed at localhost:8888


# To get the grafan password run

export GRAFANA_ADMIN_PASSWORD=$(kubectl get secret --namespace {{namespace}} prometheus-operator-grafana -o jsonpath="{.data.admin-password}" | base64 --decode)

# Prometheus UI can be access via port-forwarding on port 9090 from within your cluster using below command:

kubectl port-forward svc/prometheus-operated 8000:9090

Grafana can now be accessed at localhost:8000

# Alert Manager can be access via port-forwarding on port 9093 from within your cluster using below command:

kubectl port-forward svc/alertmanager-operated 8080:9093

Grafana can now be accessed at localhost:8080

# More on GitHub : https://github.com/helm/charts/tree/master/stable/prometheus-operator`
