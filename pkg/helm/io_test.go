package helm

import (
	"os"
	"reflect"
	"testing"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    ValuesMap
		err     bool
	}{
		{
			name: "valid yaml file",
			content: `
replicaCount: 1
image: demo-operator
nginx:
  image: nginx
  pullPolicy: IfNotPresent
imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""`,
			want: ValuesMap{
				"replicaCount": 1,
				"image":        "demo-operator",
				"nginx": ValuesMap{
					"image":      "nginx",
					"pullPolicy": "IfNotPresent",
				},
				"imagePullSecrets": []interface{}{},
				"nameOverride":     "",
				"fullnameOverride": "",
			},
		},
		{
			name:    "empty yaml file",
			content: ``,
			want:    ValuesMap{},
		},
		{
			name: "invalid yaml file",
			content: `
			replicaCount: 1
			image:
			repository: nginx`,
			err: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "app-helm-*.yaml")
			if err != nil {
				t.Fatalf("failed to create yaml file for load test: %v", err)
			}
			// defer os.Remove(file.Name())

			_, err = file.WriteString(tc.content)
			if err != nil {
				t.Fatalf("failed to write yaml file for load test: %v", err)
			}

			err = file.Close()
			if err != nil {
				t.Fatalf("failed to close yaml file for load test: %v", err)
			}

			got, err := Load(file.Name())
			if !tc.err && err != nil {
				t.Fatalf("failed to parse yaml file full path for load test: %v", err)
			}

			if !tc.err && !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("fwant: %q\n but got: %q", tc.want, got)
			}
		})
	}
}

func Test_ReplaceValuesInHelmValuesFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		values      map[string]string
		want        string
	}{
		{
			name: "replace successfully",
			fileContent: `
              replicaCount: "1"
              nginx:
                image: NGINX_IMG
                pullPolicy: IfNotPresent
              imagePullSecrets: []
              nameOverride: ""
              fullnameOverride: ""`,
			values: map[string]string{
				"NGINX_IMG": "nginx",
			},
			want: `
              replicaCount: "1"
              nginx:
                image: nginx
                pullPolicy: IfNotPresent
              imagePullSecrets: []
              nameOverride: ""
              fullnameOverride: ""`,
		},
		{
			name: "replace values not found",
			fileContent: `
              replicaCount: "1"
              nginx:
                image: NGINX_IMG
                pullPolicy: IfNotPresent
              imagePullSecrets: []
              nameOverride: ""
              fullnameOverride: ""`,
			values: map[string]string{
				"TEST_IMG": "test",
			},
			want: `
              replicaCount: "1"
              nginx:
                image: NGINX_IMG
                pullPolicy: IfNotPresent
              imagePullSecrets: []
              nameOverride: ""
              fullnameOverride: ""`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "app-helm-*.yaml")
			if err != nil {
				t.Fatalf("failed to create yaml file for load test: %v", err)
			}
			defer os.Remove(file.Name())

			_, err = file.WriteString(tc.fileContent)
			if err != nil {
				t.Fatalf("failed to write yaml file for load test: %v", err)
			}

			err = file.Close()
			if err != nil {
				t.Fatalf("failed to close yaml file for load test: %v", err)
			}

			got, err := ReplaceValuesInHelmValuesFile(tc.values, file.Name())
			if err != nil {
				t.Fatalf("failed to replace values in yaml file: %v", err)
			}

			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("fwant: %q\n but got: %q", tc.want, got)
			}
		})
	}
}

func Test_FilterImagesUptoDepth(t *testing.T) {
	tests := []struct {
		name   string
		values ValuesMap
		depth  int
		want   map[string]string
	}{
		{
			name: "fetch images upto level 1",
			values: ValuesMap{
				"nginx": ValuesMap{
					"image":      "nginx",
					"pullPolicy": "IfNotPresent",
				},
			},
			depth: 1,
			want: map[string]string{
				"nginx": "nginx",
			},
		},
		{
			name: "fetch images upto level 2",
			values: ValuesMap{
				"deployment1": ValuesMap{
					"container1": ValuesMap{
						"image":      "nginx",
						"pullPolicy": "IfNotPresent",
					},
					"container2": ValuesMap{
						"image":      "deamonjob",
						"pullPolicy": "IfNotPresent",
					},
				},
			},
			depth: 2,
			want: map[string]string{
				"nginx":     "nginx",
				"deamonjob": "deamonjob",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FilterImagesUptoDepth(tc.values, tc.depth)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("fwant: %q\n but got: %q", tc.want, got)
			}
		})
	}
}
