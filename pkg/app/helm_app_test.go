package app_test

import (
	"testing"

	"github.com/alexellis/arkade/pkg/app"
	"github.com/google/go-cmp/cmp"
)

func TestHelmAppChainingOptions(t *testing.T) {

	values := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	want := &app.HelmApp{
		ChartName:       "test-name",
		ChartRepository: "http://example.com/chart/",
		ChartVersion:    "1.2.3",
		CrdURL:          "http://example2.com/crd",
		InfoMessage:     "cool message",
		Name:            "test-name-app",
		Namespace:       "kube-system",
		Values:          values,
	}

	got := &app.HelmApp{}

	got.SetChartName("test-name").
		SetChartRepository("http://example.com/chart/").
		SetChartVersion("1.2.3").
		SetCrdUrl("http://example2.com/crd").
		SetInfoMessage("cool message").
		SetName("test-name-app").
		SetNamespace("kube-system").
		AddValue("key1", "value1").
		AddValue("key2", "value2")

	if !cmp.Equal(want, got) {
		t.Errorf("want: %q, got: %q", want, got)
	}

}
