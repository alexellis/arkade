package apps_test

import (
	"testing"

	"github.com/alexellis/arkade/pkg/app"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/google/go-cmp/cmp"
)

func TestMakeMetricsServer(t *testing.T) {

	want := &app.HelmApp{
		ChartName:       "metrics-server",
		ChartRepository: "https://kubernetes-charts.storage.googleapis.com",
		ChartVersion:    "2.10.0",
		InfoMessage:     apps.MetricsInfoMsg,
		Name:            "metrics-server",
		Namespace:       "kube-system",
		Values: map[string]string{
			"args": `{--kubelet-insecure-tls,--kubelet-preferred-address-types=InternalIP\,ExternalIP\,Hostname}`,
		},
	}

	got := apps.MakeAppMetricsServer()

	if !cmp.Equal(want, got) {
		t.Errorf("want: %q, got: %q", want, got)
	}

}
