package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/spf13/cobra"
)

func MakeInstallMetricsServer() *cobra.Command {
	var metricsServer = &cobra.Command{
		Use:          "metrics-server",
		Short:        "Install metrics-server",
		Long:         `Install metrics-server to provide metrics on nodes and Pods in your cluster.`,
		Example:      `  arkade install metrics-server --namespace kube-system --helm3`,
		SilenceUsage: true,
	}

	metricsServer.Flags().StringP("namespace", "n", "kube-system", "The namespace used for installation")

	metricsServer.RunE = func(command *cobra.Command, args []string) error {

		namespace, _ := command.Flags().GetString("namespace")

		app := apps.MakeAppMetricsServer()
		app.Namespace = namespace

		err := app.Install()

		if err != nil {
			return err
		}

		fmt.Print(app.GetInfoMessage())

		return nil
	}

	return metricsServer
}

const MetricsInfoMsg = `=======================================================================
= metrics-server has been installed.                                  =
=======================================================================

# It can take a few minutes for the metrics-server to collect data
# from the cluster. Try these commands and wait a few moments if
# no data is showing.

# Check pod usage

kubectl top pod

# Check node usage

kubectl top node


# Find out more at:
# https://github.com/helm/charts/tree/master/stable/metrics-server

` + pkg.ThanksForUsing
