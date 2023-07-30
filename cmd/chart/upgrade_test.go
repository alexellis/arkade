package chart

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

type test struct {
	title   string
	current string
	tags    []string
	want    string
}

type mockLister struct {
	imageTags map[string][]string
}

func (ml mockLister) listTags(image string) ([]string, error) {
	tags, ok := ml.imageTags[image]
	if !ok {
		return nil, fmt.Errorf("Image %s not found", image)
	}

	return tags, nil
}

func newMockLister(t *testing.T, tests []test) *mockLister {
	tagsMap := make(map[string][]string)
	lister := mockLister{imageTags: tagsMap}

	if len(tests) == 0 {
		t.Fatal("Tests cannot be empty, please provide test cases")
	}

	for _, tt := range tests {
		key, _ := splitImageName(tt.current)
		// Ensuring that a test entry with the same key doesn't already exist to prevent overwriting
		if _, ok := lister.imageTags[key]; !ok {
			lister.imageTags[key] = tt.tags
		} else {
			t.Fatalf("Duplicate image name found which might break tests: %s", tt.current)
		}
	}

	return &lister
}

func Test_ChartUpgrade(t *testing.T) {
	tests := []test{
		{
			title:   "Promote from a rootless build to the latest rootless build",
			current: "moby/buildkit:v0.11.5-rootless",
			want:    "moby/buildkit:v0.11.6-rootless",
			tags:    []string{"0.10.6-rootless", "v0.11.5", "0.11.6", "0.11.6-rootless", "0.12.0-rc1", "0.12.0-rc1-rootless"},
		},
		{
			title:   "Promote from a rootless build to the latest rootless build",
			current: "moby1/buildkit:v0.10.0-rootless",
			want:    "moby1/buildkit:v0.11.6-rootless",
			tags:    []string{"0.10.0", "0.10.6", "0.10.6-rootless", "0.11.6", "0.11.6-rootless", "0.12.0-rc1", "0.12.0-rc1-rootless"},
		},
		{
			title:   "Promote from non-rootless build to the latest stable build",
			current: "moby2/buildkit:v0.8.0",
			want:    "moby2/buildkit:v0.11.6",
			tags:    []string{"v0.8.0", "0.10.6-rootless", "0.11.6", "0.11.6-rootless", "0.12.0-rc1", "0.12.0-rc1-rootless"},
		},
		{
			title:   "Promote from stable release to the latest stable release",
			current: "prom/prometheus:v2.42.0",
			want:    "prom/prometheus:v2.44.0",
			tags:    []string{"2.41.5", "2.42.0", "2.43.0", "2.44.0", "2.45.0-rc.0", "2.45.0-rc.1"},
		},
		{
			title:   "Promote from stable release to the latest stable release",
			current: "prom1/prometheus:v2.2.0",
			want:    "prom1/prometheus:v2.44.0",
			tags:    []string{"v2.1.0", "v2.2.0", "2.41.5", "2.42.0", "2.43.0", "2.44.0", "2.45.0-rc.0", "2.45.0-rc.1"},
		},
		{
			title:   "Promote from release condidate to the latest stable release",
			current: "bitnami/postgresql:v2.42.0-rc.2",
			want:    "bitnami/postgresql:v2.44.0",
			tags:    []string{"v2.42.0-rc.5", "2.41.5", "2.42.0", "2.43.0", "2.44.0", "2.45.0-rc.0", "2.45.0-rc.1"},
		},
		{
			title:   "Promote from prerelease version to the latest prerelease version",
			current: "bitnami/rabbitmq:v2.22.0-rc1",
			want:    "bitnami/rabbitmq:v2.22.0-rc3",
			tags:    []string{"v2.21.0", "v2.22.0-rc1", "v2.22.0-rc2", "v2.22.0-rc3"},
		},
		{
			title:   "Remain on current version when there isn't a later version",
			current: "openfaas/openfaas:v2.23.3",
			want:    "openfaas/openfaas:v2.23.3",
			tags:    []string{"v2.23.3", "v2.21.0", "v2.15.6"},
		},
	}

	testFile, err := os.CreateTemp(os.TempDir(), "arkade_*.yml")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(testFile.Name())

	testData := make(map[string]map[string]string)
	for i, t := range tests {
		title := fmt.Sprintf("test%d", i)
		testData[title] = map[string]string{"image": t.current}
	}

	yamlBytes, err := yaml.Marshal(testData)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := testFile.Write(yamlBytes); err != nil {
		t.Fatal(err)
	}

	lister := newMockLister(t, tests)

	cmd := MakeUpgrade(lister)
	cmd.SetArgs([]string{"--write", "--file", testFile.Name()})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	yamlData, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	var final map[string]map[string]string
	err = yaml.Unmarshal(yamlData, &final)
	if err != nil {
		t.Fatal(err)
	}

	for i, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			title := fmt.Sprintf("test%d", i)
			got := final[title]["image"]
			if got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
