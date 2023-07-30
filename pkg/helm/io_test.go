package helm

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func Test_ReplaceValuesInHelmValuesFile(t *testing.T) {
	tests := []struct {
		current string
		update  string
	}{
		{current: "ghcr.io/openfaas/faas-netes:0.1.0", update: "ghcr.io/openfaas/faas-netes:0.2.0"},
		{current: "ghcr.io/openfaas/faas-netes:0.1.0-rc", update: "ghcr.io/openfaas/faas-netes:0.2.0"},
		{current: "prom/prometheus:v2.43.0", update: "prom/prometheus:v2.45.0"},
		{current: "prom/prometheus:v2.43.0-rc.0", update: "prom/prometheus:v2.45.0"},
		{current: "ghcr.io/openfaas/faas-netes:0.1.0-rc.1", update: "ghcr.io/openfaas/faas-netes:0.2.0"},
		{current: "ghcr.io/openfaas/faas-netes:0.1.0-rc.1-rc.2", update: "ghcr.io/openfaas/faas-netes:0.2.0"},
	}

	testFile, err := os.CreateTemp(os.TempDir(), "arkade_*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testFile.Name())

	yamlData := make(map[string]map[string]string)
	for i, t := range tests {
		title := fmt.Sprintf("test%d", i)
		yamlData[title] = map[string]string{"image": t.current}
	}

	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := testFile.Write(yamlBytes); err != nil {
		t.Fatal(err)
	}

	input := make(map[string]string)
	for _, t := range tests {
		input[t.current] = t.update
	}

	// Test is run multiple times to verify consistent and correct output independent of the
	// input map's iteration order, addressing potential non-deterministic behaviour
	count := 10
	for i := 0; i < count; i++ {
		out, err := ReplaceValuesInHelmValuesFile(input, testFile.Name())
		if err != nil {
			t.Fatal(err)
		}

		var results map[string]map[string]string
		if err := yaml.Unmarshal([]byte(out), &results); err != nil {
			t.Fatal(err)
		}

		for i, tc := range tests {
			title := fmt.Sprintf("test%d", i)
			got := results[title]["image"]
			if got != tc.update {
				t.Fatalf("want: %s, got: %s", tc.update, got)
			}
		}
	}
}
