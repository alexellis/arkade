package apps

import (
	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/app"
)

// MakeAppMetricsServer returns the customized metrics-server chart
// - All static configuration from the maintainer is done.
func MakeAppMetricsServer() *app.HelmApp {

	app := &app.HelmApp{}

	app.SetName("metrics-server").
		SetChartName("metrics-server").
		SetChartRepository("https://kubernetes-charts.storage.googleapis.com").
		SetChartVersion("2.10.0").
		SetNamespace("kube-system").
		SetInfoMessage(MetricsInfoMsg).
		AddValue("args", `{--kubelet-insecure-tls,--kubelet-preferred-address-types=InternalIP\,ExternalIP\,Hostname}`)

	return app
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
