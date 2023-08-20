package get

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/Masterminds/semver"
)

var faasCLIVersionConstraint, _ = semver.NewConstraint(">= 0.13.2")

const arch32bit = "i686"
const arch64bit = "x86_64"
const archARM7 = "armv7l"
const archARM64 = "aarch64"
const archDarwinARM64 = "arm64"

type test struct {
	os      string
	arch    string
	version string
	url     string
	// Optional fields
	binary string
}

func getTool(name string, tools []Tool) *Tool {
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}
	return tool
}

func Test_MakeSureNoDuplicates(t *testing.T) {
	count := map[string]int{}
	tools := MakeTools()
	dupes := []string{}

	for _, tool := range tools {
		count[tool.Name]++

		if count[tool.Name] > 1 {
			dupes = append(dupes, tool.Name)
		}
	}
	if len(dupes) > 0 {
		t.Fatalf("Duplicate tools found which will break get-arkade GitHub Action: %v", dupes)
	}
}

func Test_MakeSureToolsAreSorted(t *testing.T) {
	got := Tools{
		{
			Owner: "roboll",
			Repo:  "helmfile",
			Name:  "helmfile",
		},
		{
			Owner: "kubernetes",
			Repo:  "kubernetes",
			Name:  "kubectl",
		},
		{
			Owner: "digitalocean",
			Repo:  "doctl",
			Name:  "doctl",
		},
	}

	sort.Sort(got)

	want := Tools{
		{
			Owner: "digitalocean",
			Repo:  "doctl",
			Name:  "doctl",
		},
		{
			Owner: "roboll",
			Repo:  "helmfile",
			Name:  "helmfile",
		},
		{
			Owner: "kubernetes",
			Repo:  "kubernetes",
			Name:  "kubectl",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func Test_PostInstallationMsg(t *testing.T) {

	testCases := []struct {
		defaultDownloadDir string
		localToolsStore    []ToolLocal
		want               string
	}{
		{
			defaultDownloadDir: "",
			localToolsStore: []ToolLocal{
				{Name: "yq",
					Path: "/home/user/.arkade/bin/yq",
				},
				{
					Name: "jq",
					Path: "/home/user/.arkade/bin/jq",
				}},
			want: `# Add arkade binary directory to your PATH variable
export PATH=$PATH:$HOME/.arkade/bin/

# Test the binary:
/home/user/.arkade/bin/yq
/home/user/.arkade/bin/jq

# Or install with:
sudo mv /home/user/.arkade/bin/yq /usr/local/bin/
sudo mv /home/user/.arkade/bin/jq /usr/local/bin/`,
		},
		{
			defaultDownloadDir: "/tmp/bin/",
			localToolsStore: []ToolLocal{
				{Name: "yq",
					Path: "/tmp/bin/yq_linux_amd64",
				},
				{
					Name: "jq",
					Path: "/tmp/bin/jq-linux64",
				}},
			want: `Run the following to copy to install the tool:
sudo install -m 755 /tmp/bin/yq_linux_amd64 /usr/local/bin/yq
sudo install -m 755 /tmp/bin/jq-linux64 /usr/local/bin/jq`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.localToolsStore[0].Name, func(t *testing.T) {
			defaultDownloadDir := tt.defaultDownloadDir
			msg, _ := PostInstallationMsg(defaultDownloadDir, tt.localToolsStore)

			got := string(msg)

			if got != tt.want {
				t.Errorf("got\n%s\n\nwant\n%s", got, tt.want)
			}
		})
	}
}

func Test_GetDownloadURLs(t *testing.T) {
	tools := MakeTools()
	kubectlVersion := "v1.26.2"

	tests := []struct {
		name    string
		url     string
		version string
		os      string
		arch    string
	}{
		{
			name:    "kubectl",
			url:     "https://storage.googleapis.com/kubernetes-release/release/" + kubectlVersion + "/bin/linux/amd64/kubectl",
			version: kubectlVersion,
			os:      "linux",
			arch:    "x86_64",
		},
		{
			name:    "kubectl",
			url:     "https://storage.googleapis.com/kubernetes-release/release/" + kubectlVersion + "/bin/darwin/amd64/kubectl",
			version: kubectlVersion,
			os:      "darwin",
			arch:    "x86_64",
		},
		{
			name:    "kubectl",
			url:     "https://storage.googleapis.com/kubernetes-release/release/" + kubectlVersion + "/bin/linux/arm64/kubectl",
			version: kubectlVersion,
			os:      "linux",
			arch:    "aarch64",
		},
		{
			name:    "kubectl",
			url:     "https://storage.googleapis.com/kubernetes-release/release/" + kubectlVersion + "/bin/darwin/arm64/kubectl",
			version: kubectlVersion,
			os:      "darwin",
			arch:    archDarwinARM64,
		},
		{
			name:    "kubectl",
			url:     "https://storage.googleapis.com/kubernetes-release/release/" + kubectlVersion + "/bin/linux/amd64/kubectl",
			version: kubectlVersion,
			os:      "linux",
			arch:    "x86_64",
		},
		{
			name:    "faas-cli",
			url:     "https://github.com/openfaas/faas-cli/releases/download/0.13.14/faas-cli",
			version: "0.13.14",
			os:      "linux",
			arch:    "x86_64",
		},
		{
			name:    "terraform",
			url:     "https://releases.hashicorp.com/terraform/1.3.9/terraform_1.3.9_linux_amd64.zip",
			version: "1.3.9",
			os:      "linux",
			arch:    "x86_64",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			toolList, err := GetDownloadURLs(tools, []string{tc.name}, tc.version)
			if err != nil {
				t.Fatal(err)
			}

			tool := toolList[0]
			got, err := tool.GetURL(tc.os, tc.arch, tool.Version, false)
			if err != nil {
				t.Fatal(err)
			}

			if got != tc.url {
				t.Fatalf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadArkade(t *testing.T) {
	tools := MakeTools()
	name := "arkade"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.8.28",
			url:     "https://github.com/alexellis/arkade/releases/download/0.8.28/arkade.exe"},
		{os: "darwin",
			arch:    arch64bit,
			version: "0.8.28",
			url:     "https://github.com/alexellis/arkade/releases/download/0.8.28/arkade-darwin"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "0.8.28",
			url:     "https://github.com/alexellis/arkade/releases/download/0.8.28/arkade-darwin-arm64"},
		{os: "linux",
			arch:    arch64bit,
			version: "0.8.28",
			url:     "https://github.com/alexellis/arkade/releases/download/0.8.28/arkade"},
		{os: "linux",
			arch:    "armv6l",
			version: "0.8.28",
			url:     "https://github.com/alexellis/arkade/releases/download/0.8.28/arkade-armhf"},
		{os: "linux",
			arch:    "armv7l",
			version: "0.8.28",
			url:     "https://github.com/alexellis/arkade/releases/download/0.8.28/arkade-armhf"},
		{os: "linux",
			arch:    archARM64,
			version: "0.8.28",
			url:     "https://github.com/alexellis/arkade/releases/download/0.8.28/arkade-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_Download_RunJob(t *testing.T) {
	tools := MakeTools()
	name := "run-job"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/alexellis/run-job/releases/download/0.0.1/run-job.exe"},
		{os: "darwin",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/alexellis/run-job/releases/download/0.0.1/run-job-darwin"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "0.0.1",
			url:     "https://github.com/alexellis/run-job/releases/download/0.0.1/run-job-darwin-arm64"},
		{os: "linux",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/alexellis/run-job/releases/download/0.0.1/run-job"},
		{os: "linux",
			arch:    "armv6l",
			version: "0.0.1",
			url:     "https://github.com/alexellis/run-job/releases/download/0.0.1/run-job-armhf"},
		{os: "linux",
			arch:    "armv7l",
			version: "0.0.1",
			url:     "https://github.com/alexellis/run-job/releases/download/0.0.1/run-job-armhf"},
		{os: "linux",
			arch:    archARM64,
			version: "0.0.1",
			url:     "https://github.com/alexellis/run-job/releases/download/0.0.1/run-job-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_Download_ActuatedCLI(t *testing.T) {
	tools := MakeTools()
	name := "actuated-cli"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/self-actuated/actuated-cli/releases/download/0.0.1/actuated-cli.exe"},
		{os: "darwin",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/self-actuated/actuated-cli/releases/download/0.0.1/actuated-cli-darwin"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "0.0.1",
			url:     "https://github.com/self-actuated/actuated-cli/releases/download/0.0.1/actuated-cli-darwin-arm64"},
		{os: "linux",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/self-actuated/actuated-cli/releases/download/0.0.1/actuated-cli"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_Download_mixctl(t *testing.T) {
	tools := MakeTools()
	name := "mixctl"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/inlets/mixctl/releases/download/0.0.1/mixctl.exe"},
		{os: "darwin",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/inlets/mixctl/releases/download/0.0.1/mixctl-darwin"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "0.0.1",
			url:     "https://github.com/inlets/mixctl/releases/download/0.0.1/mixctl-darwin-arm64"},
		{os: "linux",
			arch:    arch64bit,
			version: "0.0.1",
			url:     "https://github.com/inlets/mixctl/releases/download/0.0.1/mixctl"},
		{os: "linux",
			arch:    "armv6l",
			version: "0.0.1",
			url:     "https://github.com/inlets/mixctl/releases/download/0.0.1/mixctl-armhf"},
		{os: "linux",
			arch:    "armv7l",
			version: "0.0.1",
			url:     "https://github.com/inlets/mixctl/releases/download/0.0.1/mixctl-armhf"},
		{os: "linux",
			arch:    archARM64,
			version: "0.0.1",
			url:     "https://github.com/inlets/mixctl/releases/download/0.0.1/mixctl-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadKubectl(t *testing.T) {
	tools := MakeTools()
	name := "kubectl"

	tool := getTool(name, tools)

	tests := []test{
		{os: "darwin",
			arch:    arch64bit,
			version: "v1.20.0",
			url:     "https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/darwin/amd64/kubectl"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "v1.20.0",
			url:     "https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/darwin/arm64/kubectl"},
		{os: "linux",
			arch:    arch64bit,
			version: "v1.20.0",
			url:     "https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kubectl"},
		{os: "linux",
			arch:    archARM64,
			version: "v1.20.0",
			url:     "https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/arm64/kubectl"},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadKubectx(t *testing.T) {
	tools := MakeTools()
	name := "kubectx"
	tool := getTool(name, tools)

	got, err := tool.GetURL("linux", arch64bit, "v0.9.4", false)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/ahmetb/kubectx/releases/download/v0.9.4/kubectx"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadKubens(t *testing.T) {
	tools := MakeTools()
	name := "kubens"
	tool := getTool(name, tools)

	got, err := tool.GetURL("linux", arch64bit, tool.Version, false)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/ahmetb/kubectx/releases/download/v0.9.5/kubens"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadKubeseal(t *testing.T) {
	tools := MakeTools()
	name := "kubeseal"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v0.17.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.4/kubeseal-0.17.4-windows-amd64.tar.gz"},
		{os: "linux",
			arch:    arch64bit,
			version: "v0.17.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.4/kubeseal-0.17.4-linux-amd64.tar.gz"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "v0.17.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.4/kubeseal-0.17.4-darwin-arm64.tar.gz"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v0.17.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.4/kubeseal-0.17.4-darwin-amd64.tar.gz"},
		{os: "linux",
			arch:    archARM7,
			version: "v0.17.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.4/kubeseal-0.17.4-linux-arm.tar.gz"},
		{os: "linux",
			arch:    archARM64,
			version: "v0.17.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.17.4/kubeseal-0.17.4-linux-arm64.tar.gz"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("for %s/%s, want: %q, but got: %q", tc.os, tc.arch, tc.url, got)
		}
	}
}

func Test_DownloadKind(t *testing.T) {
	tools := MakeTools()
	name := "kind"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v0.11.0",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.11.0/kind-windows-amd64"},
		{os: "linux",
			arch:    arch64bit,
			version: "v0.11.0",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.11.0/kind-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v0.11.0",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.11.0/kind-darwin-amd64"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "v0.11.0",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.11.0/kind-darwin-arm64"},
		{os: "linux",
			arch:    archARM7,
			version: "v0.11.0",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.11.0/kind-linux-arm"},
		{os: "linux",
			arch:    "aarch64",
			version: "v0.11.0",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.11.0/kind-linux-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadK3d(t *testing.T) {
	tools := MakeTools()
	name := "k3d"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v3.0.0",
			url:     "https://github.com/k3d-io/k3d/releases/download/v3.0.0/k3d-windows-amd64.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: "v3.0.0",
			url:     "https://github.com/k3d-io/k3d/releases/download/v3.0.0/k3d-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v3.0.0",
			url:     "https://github.com/k3d-io/k3d/releases/download/v3.0.0/k3d-darwin-amd64"},
		{os: "linux",
			arch:    archARM7,
			version: "v3.0.0",
			url:     "https://github.com/k3d-io/k3d/releases/download/v3.0.0/k3d-linux-arm"},
		{os: "linux",
			arch:    "aarch64",
			version: "v3.0.0",
			url:     "https://github.com/k3d-io/k3d/releases/download/v3.0.0/k3d-linux-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadK3s(t *testing.T) {
	tools := MakeTools()
	name := "k3s"

	tool := getTool(name, tools)

	tests := []test{
		{os: "linux",
			arch:    arch64bit,
			version: "v1.21.4+k3s1",
			url:     "https://github.com/k3s-io/k3s/releases/download/v1.21.4+k3s1/k3s"},
		{os: "linux",
			arch:    "aarch64",
			version: "v1.21.4+k3s1",
			url:     "https://github.com/k3s-io/k3s/releases/download/v1.21.4+k3s1/k3s-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadK0s(t *testing.T) {
	tools := MakeTools()
	name := "k0s"

	tool := getTool(name, tools)

	tests := []test{
		{os: "linux",
			arch:    arch64bit,
			version: "v1.27.4+k0s.0",
			url:     "https://github.com/k0sproject/k0s/releases/download/v1.27.4+k0s.0/k0s-v1.27.4+k0s.0-amd64"},
		{os: "linux",
			arch:    "aarch64",
			version: "v1.27.4+k0s.0",
			url:     "https://github.com/k0sproject/k0s/releases/download/v1.27.4+k0s.0/k0s-v1.27.4+k0s.0-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want:\n%q, got:\n%q", tc.url, got)
		}
	}
}

func Test_DownloadDevspace(t *testing.T) {
	tools := MakeTools()
	name := "devspace"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v5.15.0",
			url:     "https://github.com/devspace-sh/devspace/releases/download/v5.15.0/devspace-windows-amd64.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: "v5.15.0",
			url:     "https://github.com/devspace-sh/devspace/releases/download/v5.15.0/devspace-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v5.15.0",
			url:     "https://github.com/devspace-sh/devspace/releases/download/v5.15.0/devspace-darwin-amd64"},
		{os: "darwin",
			arch:    "aarch64",
			version: "v5.15.0",
			url:     "https://github.com/devspace-sh/devspace/releases/download/v5.15.0/devspace-darwin-arm64"},
		{os: "linux",
			arch:    "aarch64",
			version: "v5.15.0",
			url:     "https://github.com/devspace-sh/devspace/releases/download/v5.15.0/devspace-linux-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadTilt(t *testing.T) {
	tools := MakeTools()
	name := "tilt"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v0.22.8",
			url:     "https://github.com/tilt-dev/tilt/releases/download/v0.22.8/tilt.0.22.8.windows.x86_64.zip"},
		{os: "linux",
			arch:    arch64bit,
			version: "v0.22.8",
			url:     "https://github.com/tilt-dev/tilt/releases/download/v0.22.8/tilt.0.22.8.linux.x86_64.tar.gz"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v0.22.8",
			url:     "https://github.com/tilt-dev/tilt/releases/download/v0.22.8/tilt.0.22.8.mac.x86_64.tar.gz"},
		{os: "darwin",
			arch:    "aarch64",
			version: "v0.22.8",
			url:     "https://github.com/tilt-dev/tilt/releases/download/v0.22.8/tilt.0.22.8.mac.arm64_ALPHA.tar.gz"},
		{os: "linux",
			arch:    "aarch64",
			version: "v0.22.8",
			url:     "https://github.com/tilt-dev/tilt/releases/download/v0.22.8/tilt.0.22.8.linux.arm64_ALPHA.tar.gz"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadK3sup(t *testing.T) {
	tools := MakeTools()
	name := "k3sup"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.9.2",
			url:     "https://github.com/alexellis/k3sup/releases/download/0.9.2/k3sup.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: "0.9.2",
			url:     "https://github.com/alexellis/k3sup/releases/download/0.9.2/k3sup"},
		{os: "darwin",
			arch:    arch64bit,
			version: "0.9.2",
			url:     "https://github.com/alexellis/k3sup/releases/download/0.9.2/k3sup-darwin"},
		{os: "linux",
			arch:    archARM7,
			version: "0.9.2",
			url:     "https://github.com/alexellis/k3sup/releases/download/0.9.2/k3sup-armhf"},
		{os: "linux",
			arch:    "aarch64",
			version: "0.9.2",
			url:     "https://github.com/alexellis/k3sup/releases/download/0.9.2/k3sup-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}
func Test_DownloadAutok3s(t *testing.T) {
	tools := MakeTools()
	name := "autok3s"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v0.4.4",
			url:     "https://github.com/cnrancher/autok3s/releases/download/v0.4.4/autok3s_windows_amd64.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: "v0.4.4",
			url:     "https://github.com/cnrancher/autok3s/releases/download/v0.4.4/autok3s_linux_amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v0.4.4",
			url:     "https://github.com/cnrancher/autok3s/releases/download/v0.4.4/autok3s_darwin_amd64"},
		{os: "linux",
			arch:    archARM7,
			version: "v0.4.4",
			url:     "https://github.com/cnrancher/autok3s/releases/download/v0.4.4/autok3s_linux_arm"},
		{os: "linux",
			arch:    "aarch64",
			version: "v0.4.4",
			url:     "https://github.com/cnrancher/autok3s/releases/download/v0.4.4/autok3s_linux_arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadInletsctl(t *testing.T) {
	tools := MakeTools()
	name := "inletsctl"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.8.16",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.8.16/inletsctl.exe.tgz",
			binary:  "inletsctl"},
		{os: "darwin",
			arch:    arch64bit,
			version: "0.8.16",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.8.16/inletsctl-darwin.tgz",
			binary:  "inletsctl-darwin"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: "0.8.16",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.8.16/inletsctl-darwin-arm64.tgz",
			binary:  "inletsctl-darwin-arm64"},
		{os: "linux",
			arch:    arch64bit,
			version: "0.8.16",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.8.16/inletsctl.tgz",
			binary:  "inletsctl"},
		{os: "linux",
			arch:    "armv6l",
			version: "0.8.16",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.8.16/inletsctl-armhf.tgz",
			binary:  "inletsctl-armhf"},
		{os: "linux",
			arch:    "armv7l",
			version: "0.8.16",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.8.16/inletsctl-armhf.tgz",
			binary:  "inletsctl-armhf"},
		{os: "linux",
			arch:    archARM64,
			version: "0.8.16",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.8.16/inletsctl-arm64.tgz",
			binary:  "inletsctl-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("for %s/%s, want: %q, but got: %q", tc.os, tc.arch, tc.url, got)
		}
		binary, err := GetBinaryName(tool, tc.os, tc.arch, tc.version)
		if err != nil {
			t.Fatal(err)
		}
		if binary != tc.binary {
			t.Errorf("for %s/%s, want: %q, but got: %q", tc.os, tc.arch, tc.binary, binary)
		}
	}
}

func Test_DownloadKubebuilder(t *testing.T) {
	tools := MakeTools()
	name := "kubebuilder"

	tool := getTool(name, tools)

	tests := []test{
		{os: "darwin",
			arch:    arch64bit,
			version: "3.1.0",
			url:     "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v3.1.0/kubebuilder_darwin_amd64"},
		{os: "linux",
			arch:    arch64bit,
			version: "3.1.0",
			url:     "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v3.1.0/kubebuilder_linux_amd64"},
		{os: "linux",
			arch:    "arm64",
			version: "3.1.0",
			url:     "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v3.1.0/kubebuilder_linux_arm64"},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadKustomize(t *testing.T) {
	tools := MakeTools()
	name := "kustomize"

	tool := getTool(name, tools)

	ver := "v5.0.3"

	tests := []test{
		{os: "linux",
			arch:    arch64bit,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv5.0.3/kustomize_v5.0.3_linux_amd64.tar.gz",
		},
		{os: "darwin",
			arch:    arch64bit,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv5.0.3/kustomize_v5.0.3_darwin_amd64.tar.gz",
		},
		{os: "linux",
			arch:    archARM64,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv5.0.3/kustomize_v5.0.3_linux_arm64.tar.gz",
		},
		{os: "mingw64_nt-10.0-18362",

			arch:    arch64bit,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv5.0.3/kustomize_v5.0.3_windows_amd64.tar.gz",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadCrane(t *testing.T) {
	tools := MakeTools()
	name := "crane"

	const toolVersion = "v0.11.0"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/google/go-containerregistry/releases/download/v0.11.0/go-containerregistry_Windows_x86_64.tar.gz"},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/google/go-containerregistry/releases/download/v0.11.0/go-containerregistry_Linux_x86_64.tar.gz"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     "https://github.com/google/go-containerregistry/releases/download/v0.11.0/go-containerregistry_Darwin_arm64.tar.gz"},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/google/go-containerregistry/releases/download/v0.11.0/go-containerregistry_Darwin_x86_64.tar.gz"},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     "https://github.com/google/go-containerregistry/releases/download/v0.11.0/go-containerregistry_Linux_arm64.tar.gz"},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("for %s/%s, want: %q, but got: %q", tc.os, tc.arch, tc.url, got)
		}
	}
}

func Test_DownloadDigitalOcean(t *testing.T) {
	tools := MakeTools()
	name := "doctl"

	tool := getTool(name, tools)

	const toolVersion = "1.46.0"
	const urlTemplate = "https://github.com/digitalocean/doctl/releases/download/v1.46.0/doctl-1.46.0-%s-%s.%s"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     fmt.Sprintf(urlTemplate, "windows", "amd64", "zip")},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     fmt.Sprintf(urlTemplate, "linux", "amd64", "tar.gz")},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     fmt.Sprintf(urlTemplate, "darwin", "amd64", "tar.gz")},
		// this asserts that we can build a URL for ARM processors, but no asset exists and will yield a 404
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     fmt.Sprintf(urlTemplate, "linux", "", "tar.gz")},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadEKSCTL(t *testing.T) {
	tools := MakeTools()
	name := "eksctl"

	tool := getTool(name, tools)

	const toolVersion = "v0.79.0"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/eksctl-io/eksctl/releases/download/v0.79.0/eksctl_Windows_amd64.zip"},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/eksctl-io/eksctl/releases/download/v0.79.0/eksctl_Linux_amd64.tar.gz"},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     "https://github.com/eksctl-io/eksctl/releases/download/v0.79.0/eksctl_Linux_arm64.tar.gz"},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     "https://github.com/eksctl-io/eksctl/releases/download/v0.79.0/eksctl_Darwin_arm64.tar.gz"},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/eksctl-io/eksctl/releases/download/v0.79.0/eksctl_Darwin_amd64.tar.gz"},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     "https://github.com/eksctl-io/eksctl/releases/download/v0.79.0/eksctl_Linux_armv7.tar.gz"},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadEKSCTLANYWHERE(t *testing.T) {
	tools := MakeTools()
	name := "eksctl-anywhere"

	tool := getTool(name, tools)

	const toolVersion = "v0.12.1"

	tests := []test{
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/aws/eks-anywhere/releases/download/v0.12.1/eksctl-anywhere-v0.12.1-linux-amd64.tar.gz"},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/aws/eks-anywhere/releases/download/v0.12.1/eksctl-anywhere-v0.12.1-darwin-amd64.tar.gz"},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadK9s(t *testing.T) {
	tools := MakeTools()
	name := "k9s"

	tool := getTool(name, tools)

	const toolVersion = "v0.24.10"

	tests := []test{
		{os: "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_Windows_amd64.tar.gz`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_Linux_amd64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_Darwin_amd64.tar.gz`,
		},
		{os: "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_Darwin_arm64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_Linux_arm64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_Linux_arm.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadEtcd(t *testing.T) {
	tools := MakeTools()
	name := "etcd"

	tool := getTool(name, tools)

	const toolVersion = "v3.5.9"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-linux-amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-darwin-amd64.zip`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-darwin-arm64.zip`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-linux-arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/etcd-io/etcd/releases/download/v3.5.9/etcd-v3.5.9-windows-amd64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadCivo(t *testing.T) {
	tools := MakeTools()
	name := "civo"

	tool := getTool(name, tools)

	const toolVersion = "v0.7.11"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.7.11/civo-0.7.11-windows-amd64.zip`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.7.11/civo-0.7.11-linux-amd64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.7.11/civo-0.7.11-darwin-amd64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.7.11/civo-0.7.11-linux-arm.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadWaypoint(t *testing.T) {
	tools := MakeTools()
	name := "waypoint"

	tool := getTool(name, tools)

	const toolVersion = "0.6.1"

	tests := []test{
		{os: "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/waypoint/0.6.1/waypoint_0.6.1_windows_amd64.zip`,
		},
		{os: "darwin",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/waypoint/0.6.1/waypoint_0.6.1_darwin_arm64.zip`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/waypoint/0.6.1/waypoint_0.6.1_darwin_amd64.zip`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/waypoint/0.6.1/waypoint_0.6.1_linux_arm.zip`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/waypoint/0.6.1/waypoint_0.6.1_linux_amd64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadTerraform(t *testing.T) {
	tools := MakeTools()
	name := "terraform"

	tool := getTool(name, tools)

	const toolVersion = "1.3.9"

	tests := []test{
		{
			os:      "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/1.3.9/terraform_1.3.9_windows_amd64.zip`,
		},
		{
			url:     "https://releases.hashicorp.com/terraform/1.3.9/terraform_1.3.9_linux_amd64.zip",
			version: toolVersion,
			os:      "linux",
			arch:    arch64bit,
		},
		{
			url:     "https://releases.hashicorp.com/terraform/1.3.9/terraform_1.3.9_linux_arm.zip",
			version: toolVersion,
			os:      "linux",
			arch:    archARM7,
		},
		{
			url:     "https://releases.hashicorp.com/terraform/1.3.9/terraform_1.3.9_linux_arm64.zip",
			version: toolVersion,
			os:      "linux",
			arch:    archARM64,
		},
		{
			url:     "https://releases.hashicorp.com/terraform/1.3.9/terraform_1.3.9_darwin_arm64.zip",
			version: toolVersion,
			os:      "darwin",
			arch:    archDarwinARM64,
		},
		{
			url:     "https://releases.hashicorp.com/terraform/1.3.9/terraform_1.3.9_darwin_amd64.zip",
			version: toolVersion,
			os:      "darwin",
			arch:    arch64bit,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadPacker(t *testing.T) {
	tools := MakeTools()
	name := "packer"

	tool := getTool(name, tools)

	const toolVersion = "1.6.5"

	tests := []test{
		{
			os:      "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/packer/1.6.5/packer_1.6.5_windows_amd64.zip`,
		},
		{
			url:     "https://releases.hashicorp.com/packer/1.6.5/packer_1.6.5_linux_amd64.zip",
			version: toolVersion,
			os:      "linux",
			arch:    arch64bit,
		},
		{
			url:     "https://releases.hashicorp.com/packer/1.6.5/packer_1.6.5_linux_arm.zip",
			version: toolVersion,
			os:      "linux",
			arch:    archARM7,
		},
		{
			url:     "https://releases.hashicorp.com/packer/1.6.5/packer_1.6.5_linux_arm64.zip",
			version: toolVersion,
			os:      "linux",
			arch:    archARM64,
		},
		{
			url:     "https://releases.hashicorp.com/packer/1.6.5/packer_1.6.5_darwin_arm64.zip",
			version: toolVersion,
			os:      "darwin",
			arch:    archARM64,
		},
		{
			url:     "https://releases.hashicorp.com/packer/1.6.5/packer_1.6.5_darwin_amd64.zip",
			version: toolVersion,
			os:      "darwin",
			arch:    arch64bit,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadGH(t *testing.T) {
	tools := MakeTools()
	name := "gh"

	tool := getTool(name, tools)

	const toolVersion = "v1.6.1"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.6.1/gh_1.6.1_windows_amd64.zip`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.6.1/gh_1.6.1_linux_amd64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.6.1/gh_1.6.1_macOS_amd64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.6.1/gh_1.6.1_linux_arm64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadPack(t *testing.T) {
	tools := MakeTools()
	name := "pack"

	tool := getTool(name, tools)

	const toolVersion = "v0.14.2"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			version: toolVersion,
			url:     `https://github.com/buildpacks/pack/releases/download/v0.14.2/pack-v0.14.2-windows.zip`,
		},
		{os: "darwin",
			version: toolVersion,
			url:     `https://github.com/buildpacks/pack/releases/download/v0.14.2/pack-v0.14.2-macos.tgz`,
		},
		{os: "linux",
			version: toolVersion,
			url:     `https://github.com/buildpacks/pack/releases/download/v0.14.2/pack-v0.14.2-linux.tgz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, "", tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadBuildx(t *testing.T) {
	tools := MakeTools()
	name := "buildx"

	tool := getTool(name, tools)

	const toolVersion = "v0.8.2"

	tests := []test{
		{
			os:      "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.8.2/buildx-v0.8.2.windows-amd64.exe`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.8.2/buildx-v0.8.2.darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.8.2/buildx-v0.8.2.darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.8.2/buildx-v0.8.2.linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.8.2/buildx-v0.8.2.linux-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.8.2/buildx-v0.8.2.linux-arm-v7`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadDockerCompose(t *testing.T) {
	tools := MakeTools()
	name := "docker-compose"

	tool := getTool(name, tools)

	const toolVersion = "v2.3.4"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/compose/releases/download/v2.3.4/docker-compose-windows-x86_64.exe`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/compose/releases/download/v2.3.4/docker-compose-linux-x86_64`,
		},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/docker/compose/releases/download/v2.3.4/docker-compose-linux-aarch64`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/compose/releases/download/v2.3.4/docker-compose-darwin-x86_64`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/docker/compose/releases/download/v2.3.4/docker-compose-linux-armv7`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadHelmfile(t *testing.T) {
	tools := MakeTools()
	name := "helmfile"

	tool := getTool(name, tools)

	const toolVersion = "v0.145.4"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/helmfile/helmfile/releases/download/v0.145.4/helmfile_0.145.4_windows_amd64.tar.gz`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/helmfile/helmfile/releases/download/v0.145.4/helmfile_0.145.4_linux_amd64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/helmfile/helmfile/releases/download/v0.145.4/helmfile_0.145.4_darwin_amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadOpa(t *testing.T) {
	tools := MakeTools()
	name := "opa"

	tool := getTool(name, tools)

	const toolVersion = "v0.24.0"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			version: toolVersion,
			url:     `https://github.com/open-policy-agent/opa/releases/download/v0.24.0/opa_windows_amd64.exe`,
		},
		{os: "linux",
			version: toolVersion,
			url:     `https://github.com/open-policy-agent/opa/releases/download/v0.24.0/opa_linux_amd64`,
		},
		{os: "darwin",
			version: toolVersion,
			url:     `https://github.com/open-policy-agent/opa/releases/download/v0.24.0/opa_darwin_amd64`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, "", tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_getBinaryURL_SlashInDownloadPath(t *testing.T) {
	got := getBinaryURL("roboll", "helmfile", "0.134.0", "v0.134.0/helmfile_0.134.0_darwin_amd64")
	want := "https://github.com/roboll/helmfile/releases/download/v0.134.0/helmfile_0.134.0_darwin_amd64"
	if got != want {
		t.Fatalf("want %s, but got: %s", want, got)
	}
}

func Test_getBinaryURL_NoSlashInDownloadPath(t *testing.T) {
	got := getBinaryURL("openfaas", "faas-cli", "0.19.0", "faas-cli_darwin")
	want := "https://github.com/openfaas/faas-cli/releases/download/0.19.0/faas-cli_darwin"
	if got != want {
		t.Fatalf("want %s, but got: %s", want, got)
	}
}

func Test_DownloadMinikube(t *testing.T) {
	tools := MakeTools()
	name := "minikube"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    "amd64",
			version: "v1.25.2",
			url:     `https://github.com/kubernetes/minikube/releases/download/v1.25.2/minikube-windows-amd64.exe`,
		},
		{
			os:      "linux",
			arch:    "amd64",
			version: "v1.25.2",
			url:     `https://github.com/kubernetes/minikube/releases/download/v1.25.2/minikube-linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v1.25.2",
			url:     `https://github.com/kubernetes/minikube/releases/download/v1.25.2/minikube-linux-arm64`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: "v1.25.2",
			url:     `https://github.com/kubernetes/minikube/releases/download/v1.25.2/minikube-linux-armv6`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "v1.25.2",
			url:     `https://github.com/kubernetes/minikube/releases/download/v1.25.2/minikube-linux-arm`,
		},
		{
			os:      "darwin",
			arch:    "amd64",
			version: "v1.25.2",
			url:     `https://github.com/kubernetes/minikube/releases/download/v1.25.2/minikube-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: "v1.25.2",
			url:     `https://github.com/kubernetes/minikube/releases/download/v1.25.2/minikube-darwin-arm64`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadMinio(t *testing.T) {
	tools := MakeTools()
	name := "mc"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:   "ming",
			arch: "amd64",
			url:  `https://dl.min.io/client/mc/release/windows-amd64/mc.exe`,
		},
		{
			os:   "linux",
			arch: "amd64",
			url:  `https://dl.min.io/client/mc/release/linux-amd64/mc`,
		},
		{
			os:   "linux",
			arch: "arm",
			url:  `https://dl.min.io/client/mc/release/linux-arm/mc`,
		},
		{
			os:   "linux",
			arch: "armv6l",
			url:  `https://dl.min.io/client/mc/release/linux-arm/mc`,
		},
		{
			os:   "linux",
			arch: archARM7,
			url:  `https://dl.min.io/client/mc/release/linux-arm/mc`,
		},
		{
			os:   "linux",
			arch: archARM64,
			url:  `https://dl.min.io/client/mc/release/linux-arm64/mc`,
		},
		{
			os:   "darwin",
			arch: "amd64",
			url:  `https://dl.min.io/client/mc/release/darwin-amd64/mc`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(t *testing.T) {

			got, err := tool.GetURL(tc.os, tc.arch, "", false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadNats(t *testing.T) {
	tools := MakeTools()
	name := "nats"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    "amd64",
			version: "v0.0.28",
			url:     `https://github.com/nats-io/natscli/releases/download/v0.0.28/nats-0.0.28-windows-amd64.zip`,
		},
		{
			os:      "linux",
			arch:    "amd64",
			version: "v0.0.28",
			url:     `https://github.com/nats-io/natscli/releases/download/v0.0.28/nats-0.0.28-linux-amd64.zip`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.0.28",
			url:     `https://github.com/nats-io/natscli/releases/download/v0.0.28/nats-0.0.28-linux-arm64.zip`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: "v0.0.28",
			url:     `https://github.com/nats-io/natscli/releases/download/v0.0.28/nats-0.0.28-linux-arm6.zip`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "v0.0.28",
			url:     `https://github.com/nats-io/natscli/releases/download/v0.0.28/nats-0.0.28-linux-arm7.zip`,
		},
		{
			os:      "darwin",
			arch:    "amd64",
			version: "v0.0.28",
			url:     `https://github.com/nats-io/natscli/releases/download/v0.0.28/nats-0.0.28-darwin-amd64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadLinkerd(t *testing.T) {
	tools := MakeTools()
	name := "linkerd2"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "stable-2.9.1",
			url:     "https://github.com/linkerd/linkerd2/releases/download/stable-2.9.1/linkerd2-cli-stable-2.9.1-windows.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: "stable-2.9.1",
			url:     "https://github.com/linkerd/linkerd2/releases/download/stable-2.9.1/linkerd2-cli-stable-2.9.1-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "stable-2.9.1",
			url:     "https://github.com/linkerd/linkerd2/releases/download/stable-2.9.1/linkerd2-cli-stable-2.9.1-darwin"},
		{os: "linux",
			arch:    archARM64,
			version: "stable-2.9.1",
			url:     "https://github.com/linkerd/linkerd2/releases/download/stable-2.9.1/linkerd2-cli-stable-2.9.1-linux-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadArgocd(t *testing.T) {
	tools := MakeTools()
	name := "argocd"
	version := "v2.4.14"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/argoproj/argo-cd/releases/download/v2.4.14/argocd-windows-amd64.exe",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/argoproj/argo-cd/releases/download/v2.4.14/argocd-linux-amd64",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/argoproj/argo-cd/releases/download/v2.4.14/argocd-darwin-amd64",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/argoproj/argo-cd/releases/download/v2.4.14/argocd-darwin-arm64",
		},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadNerdctl(t *testing.T) {
	tools := MakeTools()
	name := "nerdctl"

	tool := getTool(name, tools)

	tests := []test{
		{os: "linux",
			arch:    arch64bit,
			version: "v0.7.2",
			url:     "https://github.com/containerd/nerdctl/releases/download/v0.7.2/nerdctl-0.7.2-linux-amd64.tar.gz",
		},
		{os: "linux",
			arch:    archARM7,
			version: "v0.7.2",
			url:     "https://github.com/containerd/nerdctl/releases/download/v0.7.2/nerdctl-0.7.2-linux-arm-v7.tar.gz",
		},
		{os: "linux",
			arch:    archARM64,
			version: "v0.7.2",
			url:     "https://github.com/containerd/nerdctl/releases/download/v0.7.2/nerdctl-0.7.2-linux-arm64.tar.gz",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("for %s/%s, want: %q, but got: %q", tc.os, tc.arch, tc.url, got)
		}
	}
}

func Test_DownloadIstioCtl(t *testing.T) {
	tools := MakeTools()
	name := "istioctl"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    "amd64",
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-win.zip`,
		},
		{
			os:      "linux",
			arch:    "x86_64",
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    "amd64",
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    "arm",
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-linux-armv7.tar.gz`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-linux-armv7.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-linux-armv7.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    "amd64",
			version: "1.9.1",
			url:     `https://github.com/istio/istio/releases/download/1.9.1/istioctl-1.9.1-osx.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(t *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadTektonCli(t *testing.T) {
	tools := MakeTools()
	name := "tkn"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Windows_x86_64.zip`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Linux_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Darwin_x86_64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadInfluxCli(t *testing.T) {
	tools := MakeTools()
	name := "influx"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "windows",
			arch:    arch64bit,
			version: "2.0.7",
			url:     `https://dl.influxdata.com/influxdb/releases/influxdb2-client-2.0.7-windows-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "2.0.7",
			url:     `https://dl.influxdata.com/influxdb/releases/influxdb2-client-2.0.7-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "2.0.7",
			url:     `https://dl.influxdata.com/influxdb/releases/influxdb2-client-2.0.7-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "2.0.7",
			url:     `https://dl.influxdata.com/influxdb/releases/influxdb2-client-2.0.7-darwin-amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}

}

func Test_DownloadInletsProCli(t *testing.T) {
	tools := MakeTools()
	name := "inlets-pro"
	const version = "0.9.1"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.9.1/inlets-pro.exe`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.9.1/inlets-pro`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.9.1/inlets-pro-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.9.1/inlets-pro-armhf`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: version,
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.9.1/inlets-pro-armhf`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.9.1/inlets-pro-darwin`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.9.1/inlets-pro-darwin-arm64`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}

}

func Test_DownloadFaaSCLI(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	const version = "0.16.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/openfaas/faas-cli/releases/download/0.16.0/faas-cli.exe`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/openfaas/faas-cli/releases/download/0.16.0/faas-cli`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/openfaas/faas-cli/releases/download/0.16.0/faas-cli-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/openfaas/faas-cli/releases/download/0.16.0/faas-cli-armhf`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: version,
			url:     `https://github.com/openfaas/faas-cli/releases/download/0.16.0/faas-cli-armhf`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/openfaas/faas-cli/releases/download/0.16.0/faas-cli-darwin`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/openfaas/faas-cli/releases/download/0.16.0/faas-cli-darwin-arm64`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}

}

func Test_DownloadKim(t *testing.T) {
	tools := MakeTools()
	name := "kim"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.1.0-alpha.12",
			url:     `https://github.com/rancher/kim/releases/download/v0.1.0-alpha.12/kim-windows-amd64.exe`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.1.0-alpha.12",
			url:     `https://github.com/rancher/kim/releases/download/v0.1.0-alpha.12/kim-linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.1.0-alpha.12",
			url:     `https://github.com/rancher/kim/releases/download/v0.1.0-alpha.12/kim-linux-arm64`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.1.0-alpha.12",
			url:     `https://github.com/rancher/kim/releases/download/v0.1.0-alpha.12/kim-darwin-amd64`,
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Download for: %s %s %s", tc.os, tc.arch, tc.version), func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				r.Fatal(err)
			}
			if got != tc.url {
				r.Errorf("\nwant: %s\ngot:  %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadTrivyCli(t *testing.T) {
	tools := MakeTools()
	name := "trivy"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_Linux-64bit.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "v0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_Linux-ARM.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_Linux-ARM64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_macOS-64bit.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}

}

func Test_DownloadFluxCli(t *testing.T) {
	tools := MakeTools()
	name := "flux"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "v0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_linux_arm.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_darwin_amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadPolarisCli(t *testing.T) {
	tools := MakeTools()
	name := "polaris"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/v3.2.1/polaris_darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/v3.2.1/polaris_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/v3.2.1/polaris_linux_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "v3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/v3.2.1/polaris_linux_armv7.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {

			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadHelm(t *testing.T) {
	tools := MakeTools()
	name := "helm"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: "3.5.4",
			url:     `https://get.helm.sh/helm-3.5.4-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "3.5.4",
			url:     `https://get.helm.sh/helm-3.5.4-linux-arm.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "3.5.4",
			url:     `https://get.helm.sh/helm-3.5.4-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "3.5.4",
			url:     `https://get.helm.sh/helm-3.5.4-darwin-amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "3.5.4",
			url:     `https://get.helm.sh/helm-3.5.4-darwin-arm64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				r.Fatal(err)
			}
			if got != tc.url {
				r.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadArgoCDAutopilotCli(t *testing.T) {
	tools := MakeTools()
	name := "argocd-autopilot"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.2.13",
			url:     `https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.2.13/argocd-autopilot-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.2.13",
			url:     `https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.2.13/argocd-autopilot-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.2.13",
			url:     `https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.2.13/argocd-autopilot-darwin-amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}

}

func Test_DownloadNovaCli(t *testing.T) {
	tools := MakeTools()
	name := "nova"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "2.3.2",
			url:     `https://github.com/FairwindsOps/nova/releases/download/2.3.2/nova_2.3.2_darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "2.3.2",
			url:     `https://github.com/FairwindsOps/nova/releases/download/2.3.2/nova_2.3.2_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "2.3.2",
			url:     `https://github.com/FairwindsOps/nova/releases/download/2.3.2/nova_2.3.2_linux_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "2.3.2",
			url:     `https://github.com/FairwindsOps/nova/releases/download/2.3.2/nova_2.3.2_linux_armv7.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {

			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadKubetailCli(t *testing.T) {
	tools := MakeTools()
	name := "kubetail"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "1.6.13",
			url:     `https://raw.githubusercontent.com/johanhaleby/kubetail/1.6.13/kubetail`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "1.6.13",
			url:     `https://raw.githubusercontent.com/johanhaleby/kubetail/1.6.13/kubetail`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "1.6.13",
			url:     `https://raw.githubusercontent.com/johanhaleby/kubetail/1.6.13/kubetail`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "1.6.13",
			url:     `https://raw.githubusercontent.com/johanhaleby/kubetail/1.6.13/kubetail`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadKgctl(t *testing.T) {
	tools := MakeTools()
	name := "kgctl"
	tool := getTool(name, tools)
	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "0.3.0",
			url:     `https://github.com/squat/kilo/releases/download/0.3.0/kgctl-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "0.3.0",
			url:     `https://github.com/squat/kilo/releases/download/0.3.0/kgctl-darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.3.0",
			url:     `https://github.com/squat/kilo/releases/download/0.3.0/kgctl-linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "0.3.0",
			url:     `https://github.com/squat/kilo/releases/download/0.3.0/kgctl-linux-arm`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.3.0",
			url:     `https://github.com/squat/kilo/releases/download/0.3.0/kgctl-linux-arm64`,
		},
		{
			os:      "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.3.0",
			url:     `https://github.com/squat/kilo/releases/download/0.3.0/kgctl-windows-amd64`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadEquinixMetalCli(t *testing.T) {
	tools := MakeTools()
	name := "metal"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "0.6.0-alpha2",
			url:     `https://github.com/equinix/metal-cli/releases/download/0.6.0-alpha2/metal-darwin-amd64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.6.0-alpha2",
			url:     `https://github.com/equinix/metal-cli/releases/download/0.6.0-alpha2/metal-linux-amd64`,
		},
		{
			os:      "linux",
			arch:    "aarch64",
			version: "0.6.0-alpha2",
			url:     `https://github.com/equinix/metal-cli/releases/download/0.6.0-alpha2/metal-linux-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "0.6.0-alpha2",
			url:     `https://github.com/equinix/metal-cli/releases/download/0.6.0-alpha2/metal-linux-armv7`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: "0.6.0-alpha2",
			url:     `https://github.com/equinix/metal-cli/releases/download/0.6.0-alpha2/metal-linux-armv6`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "0.6.0-alpha2",
			url:     `https://github.com/equinix/metal-cli/releases/download/0.6.0-alpha2/metal-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadPorterCli(t *testing.T) {
	tools := MakeTools()
	name := "porter"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.38.4",
			url:     `https://github.com/getporter/porter/releases/download/v0.38.4/porter-darwin-amd64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.38.4",
			url:     `https://github.com/getporter/porter/releases/download/v0.38.4/porter-linux-amd64`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.38.4",
			url:     `https://github.com/getporter/porter/releases/download/v0.38.4/porter-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadJq(t *testing.T) {
	tools := MakeTools()
	name := "jq"
	tool := getTool(name, tools)
	tool.Version = "jq-1.7"
	prefix := "https://github.com/" + tool.Owner + "/" + tool.Repo + "/releases/download/" + tool.Version + "/"

	tests := []test{
		{
			os:   "darwin",
			arch: arch64bit,
			url:  prefix + "jq-osx-amd64",
		},
		{
			os:   "darwin",
			arch: archDarwinARM64,
			url:  prefix + "jq-osx-amd64",
		},
		{
			os:   "linux",
			arch: arch64bit,
			url:  prefix + "jq-linux64",
		},
		{
			os:   "linux",
			arch: arch32bit,
			url:  prefix + "jq-linux32",
		},
		{
			os:   "ming",
			arch: arch64bit,
			url:  prefix + "jq-win64.exe",
		},
		{
			os:   "ming",
			arch: arch32bit,
			url:  prefix + "jq-win32.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tool.Version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadOperatorSDK(t *testing.T) {
	tools := MakeTools()
	tool := getTool("operator-sdk", tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v1.13.1",
			url:     "https://github.com/operator-framework/operator-sdk/releases/download/v1.13.1/operator-sdk_darwin_amd64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v1.13.1",
			url:     "https://github.com/operator-framework/operator-sdk/releases/download/v1.13.1/operator-sdk_linux_amd64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v1.13.1",
			url:     "https://github.com/operator-framework/operator-sdk/releases/download/v1.13.1/operator-sdk_linux_arm64",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadCosignCli(t *testing.T) {
	tools := MakeTools()
	name := "cosign"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-linux-amd64`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want:\n%s\ngot:\n%s", tc.url, got)
			}
		})
	}
}

func Test_DownloadKanister(t *testing.T) {
	tools := MakeTools()
	name := "kanctl"
	v := "0.63.0"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kanisterio/kanister/releases/download/0.63.0/kanister_0.63.0_darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kanisterio/kanister/releases/download/0.63.0/kanister_0.63.0_linux_amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadKubestr(t *testing.T) {
	tools := MakeTools()
	name := "kubestr"
	v := "v0.4.31"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kastenhq/kubestr/releases/download/v0.4.31/kubestr_0.4.31_MacOS_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: v,
			url:     `https://github.com/kastenhq/kubestr/releases/download/v0.4.31/kubestr_0.4.31_MacOS_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kastenhq/kubestr/releases/download/v0.4.31/kubestr_0.4.31_Linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			url:     `https://github.com/kastenhq/kubestr/releases/download/v0.4.31/kubestr_0.4.31_Linux_arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kastenhq/kubestr/releases/download/v0.4.31/kubestr_0.4.31_Windows_amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadK10multicluster(t *testing.T) {
	tools := MakeTools()
	name := "k10multicluster"
	v := "4.0.6"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kastenhq/external-tools/releases/download/4.0.6/k10multicluster_4.0.6_macOS_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kastenhq/external-tools/releases/download/4.0.6/k10multicluster_4.0.6_linux_amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadK10tools(t *testing.T) {
	tools := MakeTools()
	name := "k10tools"
	v := "4.5.9"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kastenhq/external-tools/releases/download/4.5.9/k10tools_4.5.9_macOS_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: v,
			url:     `https://github.com/kastenhq/external-tools/releases/download/4.5.9/k10tools_4.5.9_macOS_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			url:     `https://github.com/kastenhq/external-tools/releases/download/4.5.9/k10tools_4.5.9_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			url:     `https://github.com/kastenhq/external-tools/releases/download/4.5.9/k10tools_4.5.9_linux_arm64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want:\n%s\ngot:\n%s", tc.url, got)
			}
		})
	}
}

func Test_DownloadRekorCli(t *testing.T) {
	tools := MakeTools()
	name := "rekor-cli"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.3.0",
			url:     `https://github.com/sigstore/rekor/releases/download/v0.3.0/rekor-cli-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v0.3.0",
			url:     `https://github.com/sigstore/rekor/releases/download/v0.3.0/rekor-cli-darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.3.0",
			url:     `https://github.com/sigstore/rekor/releases/download/v0.3.0/rekor-cli-linux-amd64`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.3.0",
			url:     `https://github.com/sigstore/rekor/releases/download/v0.3.0/rekor-cli-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadTFSecCli(t *testing.T) {
	tools := MakeTools()
	name := "tfsec"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.57.1",
			url:     `https://github.com/aquasecurity/tfsec/releases/download/v0.57.1/tfsec-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v0.57.1",
			url:     `https://github.com/aquasecurity/tfsec/releases/download/v0.57.1/tfsec-darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.57.1",
			url:     `https://github.com/aquasecurity/tfsec/releases/download/v0.57.1/tfsec-linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.57.1",
			url:     `https://github.com/aquasecurity/tfsec/releases/download/v0.57.1/tfsec-linux-arm64`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.57.1",
			url:     `https://github.com/aquasecurity/tfsec/releases/download/v0.57.1/tfsec-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadDive(t *testing.T) {
	tools := MakeTools()
	name := "dive"
	version := "0.10.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/wagoodman/dive/releases/download/v0.10.0/dive_0.10.0_darwin_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/wagoodman/dive/releases/download/v0.10.0/dive_0.10.0_darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/wagoodman/dive/releases/download/v0.10.0/dive_0.10.0_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/wagoodman/dive/releases/download/v0.10.0/dive_0.10.0_linux_amd64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/wagoodman/dive/releases/download/v0.10.0/dive_0.10.0_windows_amd64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadGoReleaserCli(t *testing.T) {
	tools := MakeTools()
	name := "goreleaser"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.177.0",
			url:     `https://github.com/goreleaser/goreleaser/releases/download/v0.177.0/goreleaser_Darwin_x86_64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v0.177.0",
			url:     `https://github.com/goreleaser/goreleaser/releases/download/v0.177.0/goreleaser_Darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.177.0",
			url:     `https://github.com/goreleaser/goreleaser/releases/download/v0.177.0/goreleaser_Linux_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.177.0",
			url:     `https://github.com/goreleaser/goreleaser/releases/download/v0.177.0/goreleaser_Linux_arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.177.0",
			url:     `https://github.com/goreleaser/goreleaser/releases/download/v0.177.0/goreleaser_Windows_x86_64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadKubescape(t *testing.T) {
	tools := MakeTools()
	name := "kubescape"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			version: "v1.0.69",
			url:     `https://github.com/kubescape/kubescape/releases/download/v1.0.69/kubescape-macos-latest`,
		},
		{
			os:      "linux",
			version: "v1.0.69",
			url:     `https://github.com/kubescape/kubescape/releases/download/v1.0.69/kubescape-ubuntu-latest`,
		},
		{
			os:      "ming",
			version: "v1.0.69",
			url:     `https://github.com/kubescape/kubescape/releases/download/v1.0.69/kubescape-windows-latest`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadKrew(t *testing.T) {
	tools := MakeTools()
	name := "krew"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.4.3",
			url:     `https://github.com/kubernetes-sigs/krew/releases/download/v0.4.3/krew-darwin_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: "v0.4.3",
			url:     `https://github.com/kubernetes-sigs/krew/releases/download/v0.4.3/krew-darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.4.3",
			url:     `https://github.com/kubernetes-sigs/krew/releases/download/v0.4.3/krew-linux_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.4.3",
			url:     `https://github.com/kubernetes-sigs/krew/releases/download/v0.4.3/krew-linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "v0.4.3",
			url:     `https://github.com/kubernetes-sigs/krew/releases/download/v0.4.3/krew-linux_arm.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.4.3",
			url:     `https://github.com/kubernetes-sigs/krew/releases/download/v0.4.3/krew-windows_amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s for %s on %s, got: %s", tc.url, tc.os, tc.arch, got)
			}
		})
	}

}

func Test_DownloadKubeBench(t *testing.T) {
	tools := MakeTools()
	name := "kube-bench"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.6.5",
			url:     "https://github.com/aquasecurity/kube-bench/releases/download/v0.6.5/kube-bench_0.6.5_linux_amd64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.6.5",
			url:     "https://github.com/aquasecurity/kube-bench/releases/download/v0.6.5/kube-bench_0.6.5_darwin_amd64.tar.gz",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadClusterctl(t *testing.T) {
	tools := MakeTools()
	name := "clusterctl"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v1.0.0",
			url:     "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.0.0/clusterctl-linux-amd64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v1.0.0",
			url:     "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.0.0/clusterctl-linux-arm64",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v1.0.0",
			url:     "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.0.0/clusterctl-darwin-amd64",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v1.0.0",
			url:     "https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.0.0/clusterctl-darwin-arm64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v1.0.0",
			url:     `https://github.com/kubernetes-sigs/cluster-api/releases/download/v1.0.0/clusterctl-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadvCluster(t *testing.T) {
	tools := MakeTools()
	name := "vcluster"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v0.4.5",
			url:     `https://github.com/loft-sh/vcluster/releases/download/v0.4.5/vcluster-windows-amd64.exe`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.4.5",
			url:     "https://github.com/loft-sh/vcluster/releases/download/v0.4.5/vcluster-linux-amd64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.4.5",
			url:     "https://github.com/loft-sh/vcluster/releases/download/v0.4.5/vcluster-linux-arm64",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.4.5",
			url:     "https://github.com/loft-sh/vcluster/releases/download/v0.4.5/vcluster-darwin-amd64",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v0.4.5",
			url:     "https://github.com/loft-sh/vcluster/releases/download/v0.4.5/vcluster-darwin-arm64",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadHostcl(t *testing.T) {
	tools := MakeTools()
	name := "hostctl"
	version := "v1.1.3"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/guumaster/hostctl/releases/download/v1.1.3/hostctl_1.1.3_linux_64-bit.tar.gz",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/guumaster/hostctl/releases/download/v1.1.3/hostctl_1.1.3_macOS_64-bit.tar.gz",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/guumaster/hostctl/releases/download/v1.1.3/hostctl_1.1.3_macOS_arm64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    "arm64",
			version: version,
			url:     "https://github.com/guumaster/hostctl/releases/download/v1.1.3/hostctl_1.1.3_macOS_arm64.tar.gz",
		},
		{
			os:      "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/guumaster/hostctl/releases/download/v1.1.3/hostctl_1.1.3_windows_64-bit.zip",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/guumaster/hostctl/releases/download/v1.1.3/hostctl_1.1.3_linux_arm64.tar.gz",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadKubecm(t *testing.T) {
	tools := MakeTools()
	name := "kubecm"
	version := "v0.16.2"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/sunny0826/kubecm/releases/download/v0.16.2/kubecm_v0.16.2_Darwin_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/sunny0826/kubecm/releases/download/v0.16.2/kubecm_v0.16.2_Linux_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/sunny0826/kubecm/releases/download/v0.16.2/kubecm_v0.16.2_Linux_arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/sunny0826/kubecm/releases/download/v0.16.2/kubecm_v0.16.2_Windows_x86_64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadMkcert(t *testing.T) {
	tools := MakeTools()
	name := "mkcert"
	version := "v1.4.2"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/FiloSottile/mkcert/releases/download/v1.4.2/mkcert-v1.4.2-darwin-amd64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/FiloSottile/mkcert/releases/download/v1.4.2/mkcert-v1.4.2-linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/FiloSottile/mkcert/releases/download/v1.4.2/mkcert-v1.4.2-linux-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/FiloSottile/mkcert/releases/download/v1.4.2/mkcert-v1.4.2-linux-arm`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/FiloSottile/mkcert/releases/download/v1.4.2/mkcert-v1.4.2-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadSOPS(t *testing.T) {
	tools := MakeTools()
	name := "sops"
	version := "v3.7.2"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/getsops/sops/releases/download/v3.7.2/sops-v3.7.2.linux.amd64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/getsops/sops/releases/download/v3.7.2/sops-v3.7.2.linux.arm64",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/getsops/sops/releases/download/v3.7.2/sops-v3.7.2.darwin.amd64",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/getsops/sops/releases/download/v3.7.2/sops-v3.7.2.darwin.arm64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/getsops/sops/releases/download/v3.7.2/sops-v3.7.2.exe",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadDagger(t *testing.T) {
	tools := MakeTools()
	name := "dagger"

	tool := getTool(name, tools)

	version := "v0.2.4"
	tests := []test{
		{
			os:      "darwin",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/dagger/dagger/releases/download/v0.2.4/dagger_v0.2.4_darwin_arm64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/dagger/dagger/releases/download/v0.2.4/dagger_v0.2.4_linux_arm64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/dagger/dagger/releases/download/v0.2.4/dagger_v0.2.4_darwin_amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/dagger/dagger/releases/download/v0.2.4/dagger_v0.2.4_linux_amd64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/dagger/dagger/releases/download/v0.2.4/dagger_v0.2.4_windows_amd64.zip",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want:\n%s\ngot:\n%s", tc.url, got)
		}
	}
}

func Test_DownloadOhMyPosh(t *testing.T) {
	tools := MakeTools()
	name := "oh-my-posh"

	tool := getTool(name, tools)

	const toolVersion = "v7.55.2"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/jandedobbeleer/oh-my-posh/releases/download/v7.55.2/posh-windows-amd64.exe`,
		},
		{os: "mingw64_nt-10.0-18362",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/jandedobbeleer/oh-my-posh/releases/download/v7.55.2/posh-windows-arm64.exe`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/jandedobbeleer/oh-my-posh/releases/download/v7.55.2/posh-linux-amd64`,
		},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/jandedobbeleer/oh-my-posh/releases/download/v7.55.2/posh-linux-arm64`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/jandedobbeleer/oh-my-posh/releases/download/v7.55.2/posh-darwin-amd64`,
		},
		{os: "darwin",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/jandedobbeleer/oh-my-posh/releases/download/v7.55.2/posh-darwin-arm64`,
		}}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadKumactl(t *testing.T) {
	tools := MakeTools()
	name := "kumactl"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "1.4.1",
			url:     "https://download.konghq.com/mesh-alpine/kuma-1.4.1-darwin-amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "1.4.1",
			url:     "https://download.konghq.com/mesh-alpine/kuma-1.4.1-ubuntu-amd64.tar.gz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadHey(t *testing.T) {
	tools := MakeTools()
	name := "hey"
	version := "v0.0.1-rc1"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/alexellis/hey/releases/download/v0.0.1-rc1/hey-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/alexellis/hey/releases/download/v0.0.1-rc1/hey-darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/alexellis/hey/releases/download/v0.0.1-rc1/hey`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/alexellis/hey/releases/download/v0.0.1-rc1/hey-linux-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/alexellis/hey/releases/download/v0.0.1-rc1/hey-linux-armv7`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/alexellis/hey/releases/download/v0.0.1-rc1/hey.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadCaddy(t *testing.T) {
	tools := MakeTools()
	name := "caddy"
	version := "v2.5.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/caddyserver/caddy/releases/download/v2.5.0/caddy_2.5.0_mac_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/caddyserver/caddy/releases/download/v2.5.0/caddy_2.5.0_mac_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/caddyserver/caddy/releases/download/v2.5.0/caddy_2.5.0_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/caddyserver/caddy/releases/download/v2.5.0/caddy_2.5.0_linux_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/caddyserver/caddy/releases/download/v2.5.0/caddy_2.5.0_linux_armv7.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/caddyserver/caddy/releases/download/v2.5.0/caddy_2.5.0_windows_amd64.zip`,
		},
		{
			os:      "ming",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/caddyserver/caddy/releases/download/v2.5.0/caddy_2.5.0_windows_arm64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadNatsServer(t *testing.T) {
	tools := MakeTools()
	name := "nats-server"
	version := "v2.7.4"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-darwin-amd64.zip`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-darwin-arm64.zip`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-linux-amd64.zip`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-linux-arm64.zip`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-linux-arm7.zip`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/nats-io/nats-server/releases/download/v2.7.4/nats-server-v2.7.4-windows-amd64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadCilium(t *testing.T) {
	tools := MakeTools()
	name := "cilium"
	version := "v0.11.9"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/cilium/cilium-cli/releases/download/v0.11.9/cilium-darwin-amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/cilium/cilium-cli/releases/download/v0.11.9/cilium-darwin-arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/cilium/cilium-cli/releases/download/v0.11.9/cilium-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/cilium/cilium-cli/releases/download/v0.11.9/cilium-linux-arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/cilium/cilium-cli/releases/download/v0.11.9/cilium-windows-amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadTerragrunt(t *testing.T) {
	tools := MakeTools()
	name := "terragrunt"
	version := "v0.37.1"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/gruntwork-io/terragrunt/releases/download/v0.37.1/terragrunt_darwin_amd64`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/gruntwork-io/terragrunt/releases/download/v0.37.1/terragrunt_darwin_arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/gruntwork-io/terragrunt/releases/download/v0.37.1/terragrunt_linux_amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/gruntwork-io/terragrunt/releases/download/v0.37.1/terragrunt_linux_arm64`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/gruntwork-io/terragrunt/releases/download/v0.37.1/terragrunt_windows_amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadFzf(t *testing.T) {
	tools := MakeTools()
	name := "fzf"
	version := "0.30.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/junegunn/fzf/releases/download/0.30.0/fzf-0.30.0-darwin_amd64.zip`,
		},
		{
			os:      "darwin",
			arch:    "arm64",
			version: version,
			url:     `https://github.com/junegunn/fzf/releases/download/0.30.0/fzf-0.30.0-darwin_arm64.zip`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/junegunn/fzf/releases/download/0.30.0/fzf-0.30.0-linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/junegunn/fzf/releases/download/0.30.0/fzf-0.30.0-linux_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/junegunn/fzf/releases/download/0.30.0/fzf-0.30.0-linux_armv7.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/junegunn/fzf/releases/download/0.30.0/fzf-0.30.0-windows_amd64.zip`,
		},
		{
			os:      "ming",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/junegunn/fzf/releases/download/0.30.0/fzf-0.30.0-windows_arm64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadHubble(t *testing.T) {
	tools := MakeTools()
	name := "hubble"
	version := "v0.10.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/cilium/hubble/releases/download/v0.10.0/hubble-darwin-amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/cilium/hubble/releases/download/v0.10.0/hubble-darwin-arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/cilium/hubble/releases/download/v0.10.0/hubble-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/cilium/hubble/releases/download/v0.10.0/hubble-linux-arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/cilium/hubble/releases/download/v0.10.0/hubble-windows-amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadGomplate(t *testing.T) {
	tools := MakeTools()
	name := "gomplate"
	version := "v3.11.1"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/hairyhenderson/gomplate/releases/download/v3.11.1/gomplate_darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/hairyhenderson/gomplate/releases/download/v3.11.1/gomplate_darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/hairyhenderson/gomplate/releases/download/v3.11.1/gomplate_linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/hairyhenderson/gomplate/releases/download/v3.11.1/gomplate_linux-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/hairyhenderson/gomplate/releases/download/v3.11.1/gomplate_linux-armv7`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/hairyhenderson/gomplate/releases/download/v3.11.1/gomplate_windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadJust(t *testing.T) {
	tools := MakeTools()
	name := "just"
	version := "1.3.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/casey/just/releases/download/1.3.0/just-1.3.0-x86_64-apple-darwin.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/casey/just/releases/download/1.3.0/just-1.3.0-aarch64-apple-darwin.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/casey/just/releases/download/1.3.0/just-1.3.0-x86_64-unknown-linux-musl.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/casey/just/releases/download/1.3.0/just-1.3.0-aarch64-unknown-linux-musl.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/casey/just/releases/download/1.3.0/just-1.3.0-armv7-unknown-linux-musleabihf.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/casey/just/releases/download/1.3.0/just-1.3.0-x86_64-pc-windows-msvc.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadPromtool(t *testing.T) {
	tools := MakeTools()
	name := "promtool"
	version := "v2.37.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.darwin-amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    "arm64",
			version: version,
			url:     `https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.darwin-arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.linux-arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.linux-armv7.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.windows-amd64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/prometheus/prometheus/releases/download/v2.37.0/prometheus-2.37.0.windows-arm64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadTalosctl(t *testing.T) {
	tools := MakeTools()
	name := "talosctl"
	version := "v1.1.2"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/siderolabs/talos/releases/download/v1.1.2/talosctl-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/siderolabs/talos/releases/download/v1.1.2/talosctl-darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/siderolabs/talos/releases/download/v1.1.2/talosctl-linux-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/siderolabs/talos/releases/download/v1.1.2/talosctl-linux-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/siderolabs/talos/releases/download/v1.1.2/talosctl-linux-armv7`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/siderolabs/talos/releases/download/v1.1.2/talosctl-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadTerrascan(t *testing.T) {
	tools := MakeTools()
	name := "terrascan"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v1.11.0",
			url:     `https://github.com/tenable/terrascan/releases/download/v1.11.0/terrascan_1.11.0_Darwin_x86_64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v1.11.0",
			url:     `https://github.com/tenable/terrascan/releases/download/v1.11.0/terrascan_1.11.0_Darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v1.11.0",
			url:     `https://github.com/tenable/terrascan/releases/download/v1.11.0/terrascan_1.11.0_Linux_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v1.11.0",
			url:     `https://github.com/tenable/terrascan/releases/download/v1.11.0/terrascan_1.11.0_Linux_arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v1.11.0",
			url:     `https://github.com/tenable/terrascan/releases/download/v1.11.0/terrascan_1.11.0_Windows_x86_64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadGolangciLint(t *testing.T) {
	tools := MakeTools()
	name := "golangci-lint"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v1.48.0",
			url:     `https://github.com/golangci/golangci-lint/releases/download/v1.48.0/golangci-lint-1.48.0-darwin-amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: "v1.48.0",
			url:     `https://github.com/golangci/golangci-lint/releases/download/v1.48.0/golangci-lint-1.48.0-darwin-arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v1.48.0",
			url:     `https://github.com/golangci/golangci-lint/releases/download/v1.48.0/golangci-lint-1.48.0-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v1.48.0",
			url:     `https://github.com/golangci/golangci-lint/releases/download/v1.48.0/golangci-lint-1.48.0-linux-arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v1.48.0",
			url:     `https://github.com/golangci/golangci-lint/releases/download/v1.48.0/golangci-lint-1.48.0-windows-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "v1.48.0",
			url:     `https://github.com/golangci/golangci-lint/releases/download/v1.48.0/golangci-lint-1.48.0-linux-armv7.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadBun(t *testing.T) {
	tools := MakeTools()
	name := "bun"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.1.8",
			url:     `https://github.com/oven-sh/bun/releases/download/v0.1.8/bun-darwin-x64.zip`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: "v0.1.8",
			url:     `https://github.com/oven-sh/bun/releases/download/v0.1.8/bun-darwin-aarch64.zip`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.1.8",
			url:     `https://github.com/oven-sh/bun/releases/download/v0.1.8/bun-linux-x64.zip`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.1.8",
			url:     `https://github.com/oven-sh/bun/releases/download/v0.1.8/bun-linux-aarch64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadLazygit(t *testing.T) {
	tools := MakeTools()
	name := "lazygit"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "v0.35",
			url:     `https://github.com/jesseduffield/lazygit/releases/download/v0.35/lazygit_0.35_Darwin_x86_64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: "v0.35",
			url:     `https://github.com/jesseduffield/lazygit/releases/download/v0.35/lazygit_0.35_Darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.35",
			url:     `https://github.com/jesseduffield/lazygit/releases/download/v0.35/lazygit_0.35_Linux_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.35",
			url:     `https://github.com/jesseduffield/lazygit/releases/download/v0.35/lazygit_0.35_Linux_arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.35",
			url:     `https://github.com/jesseduffield/lazygit/releases/download/v0.35/lazygit_0.35_Windows_x86_64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}

func Test_DownloadRpk(t *testing.T) {
	tools := MakeTools()
	name := "rpk"
	version := "v22.1.7"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/redpanda-data/redpanda/releases/download/v22.1.7/rpk-darwin-amd64.zip`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/redpanda-data/redpanda/releases/download/v22.1.7/rpk-darwin-arm64.zip`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/redpanda-data/redpanda/releases/download/v22.1.7/rpk-linux-amd64.zip`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/redpanda-data/redpanda/releases/download/v22.1.7/rpk-linux-arm64.zip`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadVault(t *testing.T) {
	tools := MakeTools()
	name := "vault"

	tool := getTool(name, tools)

	const toolVersion = "1.11.2"

	tests := []test{
		{
			url:     `https://releases.hashicorp.com/vault/1.11.2/vault_1.11.2_windows_amd64.zip`,
			version: toolVersion,
			os:      "ming",
			arch:    arch64bit,
		},
		{
			url:     "https://releases.hashicorp.com/vault/1.11.2/vault_1.11.2_linux_amd64.zip",
			version: toolVersion,
			os:      "linux",
			arch:    arch64bit,
		},
		{
			url:     "https://releases.hashicorp.com/vault/1.11.2/vault_1.11.2_linux_arm.zip",
			version: toolVersion,
			os:      "linux",
			arch:    archARM7,
		},
		{
			url:     "https://releases.hashicorp.com/vault/1.11.2/vault_1.11.2_linux_arm64.zip",
			version: toolVersion,
			os:      "linux",
			arch:    archARM64,
		},
		{
			url:     "https://releases.hashicorp.com/vault/1.11.2/vault_1.11.2_darwin_arm64.zip",
			version: toolVersion,
			os:      "darwin",
			arch:    archDarwinARM64,
		},
		{
			url:     "https://releases.hashicorp.com/vault/1.11.2/vault_1.11.2_darwin_amd64.zip",
			version: toolVersion,
			os:      "darwin",
			arch:    arch64bit,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadCr(t *testing.T) {
	tools := MakeTools()
	name := "cr"
	version := "v1.4.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/helm/chart-releaser/releases/download/v1.4.0/chart-releaser_1.4.0_darwin_amd64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/helm/chart-releaser/releases/download/v1.4.0/chart-releaser_1.4.0_darwin_arm64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/helm/chart-releaser/releases/download/v1.4.0/chart-releaser_1.4.0_linux_amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/helm/chart-releaser/releases/download/v1.4.0/chart-releaser_1.4.0_linux_arm64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     "https://github.com/helm/chart-releaser/releases/download/v1.4.0/chart-releaser_1.4.0_linux_armv7.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/helm/chart-releaser/releases/download/v1.4.0/chart-releaser_1.4.0_windows_amd64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadHadolint(t *testing.T) {
	tools := MakeTools()
	name := "hadolint"
	version := "v2.10.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/hadolint/hadolint/releases/download/v2.10.0/hadolint-Darwin-x86_64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/hadolint/hadolint/releases/download/v2.10.0/hadolint-Linux-x86_64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/hadolint/hadolint/releases/download/v2.10.0/hadolint-Linux-arm64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/hadolint/hadolint/releases/download/v2.10.0/hadolint-Windows-x86_64.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadButane(t *testing.T) {
	tools := MakeTools()
	name := "butane"
	version := "v0.15.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/coreos/butane/releases/download/v0.15.0/butane-x86_64-apple-darwin",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/coreos/butane/releases/download/v0.15.0/butane-x86_64-unknown-linux-gnu",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/coreos/butane/releases/download/v0.15.0/butane-aarch64-unknown-linux-gnu",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/coreos/butane/releases/download/v0.15.0/butane-x86_64-pc-windows-gnu.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadFlyctl(t *testing.T) {
	tools := MakeTools()
	name := "flyctl"
	version := "v0.0.388"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/superfly/flyctl/releases/download/v0.0.388/flyctl_0.0.388_macOS_x86_64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/superfly/flyctl/releases/download/v0.0.388/flyctl_0.0.388_Linux_x86_64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/superfly/flyctl/releases/download/v0.0.388/flyctl_0.0.388_Linux_arm64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/superfly/flyctl/releases/download/v0.0.388/flyctl_0.0.388_Windows_x86_64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadKubeconform(t *testing.T) {
	tools := MakeTools()
	name := "kubeconform"
	version := "v0.4.14"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/yannh/kubeconform/releases/download/v0.4.14/kubeconform-darwin-amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/yannh/kubeconform/releases/download/v0.4.14/kubeconform-linux-amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/yannh/kubeconform/releases/download/v0.4.14/kubeconform-linux-arm64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/yannh/kubeconform/releases/download/v0.4.14/kubeconform-windows-amd64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadConftest(t *testing.T) {
	tools := MakeTools()
	name := "conftest"
	version := "v0.34.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/open-policy-agent/conftest/releases/download/v0.34.0/conftest_0.34.0_Darwin_x86_64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/open-policy-agent/conftest/releases/download/v0.34.0/conftest_0.34.0_Linux_x86_64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/open-policy-agent/conftest/releases/download/v0.34.0/conftest_0.34.0_Linux_arm64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/open-policy-agent/conftest/releases/download/v0.34.0/conftest_0.34.0_Windows_x86_64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadKubeval(t *testing.T) {
	tools := MakeTools()
	name := "kubeval"
	version := "v0.16.1"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/instrumenta/kubeval/releases/download/v0.16.1/kubeval-darwin-amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/instrumenta/kubeval/releases/download/v0.16.1/kubeval-linux-amd64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/instrumenta/kubeval/releases/download/v0.16.1/kubeval-Windows-amd64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadViddy(t *testing.T) {
	tools := MakeTools()
	name := "viddy"
	version := "v0.3.6"

	tool := getTool(name, tools)
	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/sachaos/viddy/releases/download/v0.3.6/viddy_Darwin_x86_64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/sachaos/viddy/releases/download/v0.3.6/viddy_Darwin_arm64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/sachaos/viddy/releases/download/v0.3.6/viddy_Linux_x86_64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/sachaos/viddy/releases/download/v0.3.6/viddy_Linux_arm64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/sachaos/viddy/releases/download/v0.3.6/viddy_Windows_x86_64.tar.gz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadTctl(t *testing.T) {
	tools := MakeTools()
	name := "tctl"
	version := "v1.16.3"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/temporalio/tctl/releases/download/v1.16.3/tctl_1.16.3_darwin_amd64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/temporalio/tctl/releases/download/v1.16.3/tctl_1.16.3_darwin_arm64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/temporalio/tctl/releases/download/v1.16.3/tctl_1.16.3_linux_amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/temporalio/tctl/releases/download/v1.16.3/tctl_1.16.3_linux_arm64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/temporalio/tctl/releases/download/v1.16.3/tctl_1.16.3_windows_amd64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadFirectl(t *testing.T) {
	tools := MakeTools()
	name := "firectl"
	version := "v0.2.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/firecracker-microvm/firectl/releases/download/v0.2.0/firectl-v0.2.0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_GrafanaAgent(t *testing.T) {
	tools := MakeTools()
	name := "grafana-agent"
	version := "v0.31.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/grafana/agent/releases/download/v0.31.0/grafana-agent-linux-amd64.zip",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/grafana/agent/releases/download/v0.31.0/grafana-agent-linux-arm64.zip",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/grafana/agent/releases/download/v0.31.0/grafana-agent-darwin-amd64.zip",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/grafana/agent/releases/download/v0.31.0/grafana-agent-darwin-arm64.zip",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/grafana/agent/releases/download/v0.31.0/grafana-agent-windows-amd64.exe.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			fmt.Println(tc)
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_ScalewayCli(t *testing.T) {
	tools := MakeTools()
	name := "scaleway-cli"
	version := "v2.7.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/scaleway/scaleway-cli/releases/download/v2.7.0/scaleway-cli_2.7.0_linux_amd64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/scaleway/scaleway-cli/releases/download/v2.7.0/scaleway-cli_2.7.0_linux_arm64",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/scaleway/scaleway-cli/releases/download/v2.7.0/scaleway-cli_2.7.0_darwin_amd64",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/scaleway/scaleway-cli/releases/download/v2.7.0/scaleway-cli_2.7.0_darwin_arm64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/scaleway/scaleway-cli/releases/download/v2.7.0/scaleway-cli_2.7.0_windows_amd64.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			fmt.Println(tc)
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadSyft(t *testing.T) {
	tools := MakeTools()
	name := "syft"
	version := "v0.68.1"

	tool := getTool(name, tools)
	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/anchore/syft/releases/download/v0.68.1/syft_0.68.1_darwin_amd64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/anchore/syft/releases/download/v0.68.1/syft_0.68.1_darwin_arm64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/anchore/syft/releases/download/v0.68.1/syft_0.68.1_linux_amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/anchore/syft/releases/download/v0.68.1/syft_0.68.1_linux_arm64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/anchore/syft/releases/download/v0.68.1/syft_0.68.1_windows_amd64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadGrype(t *testing.T) {
	tools := MakeTools()
	name := "grype"
	version := "v0.55.0"

	tool := getTool(name, tools)
	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/anchore/grype/releases/download/v0.55.0/grype_0.55.0_darwin_amd64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     "https://github.com/anchore/grype/releases/download/v0.55.0/grype_0.55.0_darwin_arm64.tar.gz",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/anchore/grype/releases/download/v0.55.0/grype_0.55.0_linux_amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     "https://github.com/anchore/grype/releases/download/v0.55.0/grype_0.55.0_linux_arm64.tar.gz",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     "https://github.com/anchore/grype/releases/download/v0.55.0/grype_0.55.0_windows_amd64.zip",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloaClusterawsadm(t *testing.T) {
	tools := MakeTools()
	name := "clusterawsadm"

	tool := getTool(name, tools)

	const toolVersion = "v2.0.2"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/v2.0.2/clusterawsadm-linux-amd64`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/v2.0.2/clusterawsadm-darwin-amd64`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/v2.0.2/clusterawsadm-linux-arm64`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/v2.0.2/clusterawsadm-darwin-arm64`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases/download/v2.0.2/clusterawsadm-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloaCroc(t *testing.T) {
	tools := MakeTools()
	name := "croc"

	tool := getTool(name, tools)

	const toolVersion = "v9.6.2"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/schollz/croc/releases/download/v9.6.2/croc_9.6.2_Linux-64bit.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/schollz/croc/releases/download/v9.6.2/croc_9.6.2_macOS-64bit.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/schollz/croc/releases/download/v9.6.2/croc_9.6.2_Linux-ARM64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/schollz/croc/releases/download/v9.6.2/croc_9.6.2_macOS-ARM64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/schollz/croc/releases/download/v9.6.2/croc_9.6.2_Linux-ARM.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/schollz/croc/releases/download/v9.6.2/croc_9.6.2_Windows-64bit.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

// func Test_DownloadKubectlcnpg(t *testing.T) {
// 	tools := MakeTools()
// 	name := "kubectl-cnpg"

// 	tool := getTool(name, tools)

// 	const toolVersion = "v1.19.0"

// 	tests := []test{
// 		{
// 			os:      "linux",
// 			arch:    arch64bit,
// 			version: toolVersion,
// 			url:     `https://github.com/cloudnative-pg/cloudnative-pg/releases/download/v1.19.0/kubectl-cnpg_1.19.0_linux_x86_64.tar.gz`,
// 		},
// 		{
// 			os:      "darwin",
// 			arch:    arch64bit,
// 			version: toolVersion,
// 			url:     `https://github.com/cloudnative-pg/cloudnative-pg/releases/download/v1.19.0/kubectl-cnpg_1.19.0_darwin_x86_64.tar.gz`,
// 		},
// 		{
// 			os:      "linux",
// 			arch:    archARM64,
// 			version: toolVersion,
// 			url:     `https://github.com/cloudnative-pg/cloudnative-pg/releases/download/v1.19.0/kubectl-cnpg_1.19.0_linux_arm64.tar.gz`,
// 		},
// 		{
// 			os:      "darwin",
// 			arch:    archDarwinARM64,
// 			version: toolVersion,
// 			url:     `https://github.com/cloudnative-pg/cloudnative-pg/releases/download/v1.19.0/kubectl-cnpg_1.19.0_darwin_arm64.tar.gz`,
// 		},
// 		{
// 			os:      "linux",
// 			arch:    archARM7,
// 			version: toolVersion,
// 			url:     `https://github.com/cloudnative-pg/cloudnative-pg/releases/download/v1.19.0/kubectl-cnpg_1.19.0_linux_armv7.tar.gz`,
// 		},
// 		{
// 			os:      "ming",
// 			arch:    arch64bit,
// 			version: toolVersion,
// 			url:     `https://github.com/cloudnative-pg/cloudnative-pg/releases/download/v1.19.0/kubectl-cnpg_1.19.0_windows_x86_64.tar.gz`,
// 		},
// 	}

// 	for _, tc := range tests {
// 		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if got != tc.url {
// 			t.Errorf("want: %s, got: %s", tc.url, got)
// 		}
// 	}
// }

func Test_DownloadFstail(t *testing.T) {
	tools := MakeTools()
	name := "fstail"
	const version = "0.1.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/alexellis/fstail/releases/download/0.1.0/fstail.exe`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/alexellis/fstail/releases/download/0.1.0/fstail`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/alexellis/fstail/releases/download/0.1.0/fstail-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/alexellis/fstail/releases/download/0.1.0/fstail-armhf`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: version,
			url:     `https://github.com/alexellis/fstail/releases/download/0.1.0/fstail-armhf`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/alexellis/fstail/releases/download/0.1.0/fstail-darwin`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/alexellis/fstail/releases/download/0.1.0/fstail-darwin-arm64`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}

}

func Test_DownloadYt(t *testing.T) {
	tools := MakeTools()
	name := "yt-dlp"
	const version = "0.1.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/yt-dlp/yt-dlp/releases/download/0.1.0/yt-dlp_linux`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/yt-dlp/yt-dlp/releases/download/0.1.0/yt-dlp_linux_armv7l`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/yt-dlp/yt-dlp/releases/download/0.1.0/yt-dlp_linux_aarch64`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/yt-dlp/yt-dlp/releases/download/0.1.0/yt-dlp_macos`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/yt-dlp/yt-dlp/releases/download/0.1.0/yt-dlp_x86.exe`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadActuatedCLI(t *testing.T) {
	tools := MakeTools()
	name := "actions-usage"
	const version = "0.1.0"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/self-actuated/actions-usage/releases/download/0.1.0/actions-usage.exe`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/self-actuated/actions-usage/releases/download/0.1.0/actions-usage`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: version,
			url:     `https://github.com/self-actuated/actions-usage/releases/download/0.1.0/actions-usage-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: version,
			url:     `https://github.com/self-actuated/actions-usage/releases/download/0.1.0/actions-usage-armhf`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: version,
			url:     `https://github.com/self-actuated/actions-usage/releases/download/0.1.0/actions-usage-armhf`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: version,
			url:     `https://github.com/self-actuated/actions-usage/releases/download/0.1.0/actions-usage-darwin`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: version,
			url:     `https://github.com/self-actuated/actions-usage/releases/download/0.1.0/actions-usage-darwin-arm64`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadCmctl(t *testing.T) {
	tools := MakeTools()
	name := "cmctl"

	tool := getTool(name, tools)

	const toolVersion = "v1.11.0"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cmctl-linux-amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cmctl-darwin-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cmctl-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cmctl-darwin-arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cmctl-linux-arm.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cmctl-windows-amd64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadTimoni(t *testing.T) {
	tools := MakeTools()
	name := "timoni"

	tool := getTool(name, tools)

	const toolVersion = "v0.3.0"

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/stefanprodan/timoni/releases/download/v0.3.0/timoni_0.3.0_windows_amd64.zip`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/stefanprodan/timoni/releases/download/v0.3.0/timoni_0.3.0_linux_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/stefanprodan/timoni/releases/download/v0.3.0/timoni_0.3.0_darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/stefanprodan/timoni/releases/download/v0.3.0/timoni_0.3.0_linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/stefanprodan/timoni/releases/download/v0.3.0/timoni_0.3.0_darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/stefanprodan/timoni/releases/download/v0.3.0/timoni_0.3.0_linux_armv7l.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadSeaweedFS(t *testing.T) {
	tools := MakeTools()
	name := "seaweedfs"

	tool := getTool(name, tools)

	const toolVersion = "3.45"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/seaweedfs/seaweedfs/releases/download/3.45/linux_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/seaweedfs/seaweedfs/releases/download/3.45/darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/seaweedfs/seaweedfs/releases/download/3.45/linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/seaweedfs/seaweedfs/releases/download/3.45/darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/seaweedfs/seaweedfs/releases/download/3.45/linux_arm.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadKyverno(t *testing.T) {
	tools := MakeTools()
	name := "kyverno"

	tool := getTool(name, tools)

	const toolVersion = "v1.9.2"

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/kyverno/kyverno/releases/download/v1.9.2/kyverno-cli_v1.9.2_darwin_x86_64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/kyverno/kyverno/releases/download/v1.9.2/kyverno-cli_v1.9.2_darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/kyverno/kyverno/releases/download/v1.9.2/kyverno-cli_v1.9.2_linux_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/kyverno/kyverno/releases/download/v1.9.2/kyverno-cli_v1.9.2_linux_arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/kyverno/kyverno/releases/download/v1.9.2/kyverno-cli_v1.9.2_windows_x86_64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadBuildKit(t *testing.T) {
	tools := MakeTools()
	name := "replicated"

	tool := getTool(name, tools)

	const toolVersion = "v0.45.0"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/replicatedhq/replicated/releases/download/v0.45.0/replicated_0.45.0_linux_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/replicatedhq/replicated/releases/download/v0.45.0/replicated_0.45.0_darwin_all.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/replicatedhq/replicated/releases/download/v0.45.0/replicated_0.45.0_darwin_all.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}
}

func Test_DownloadKtop(t *testing.T) {
	tools := MakeTools()
	name := "ktop"

	tool := getTool(name, tools)

	const toolVersion = "v0.3.5"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/vladimirvivien/ktop/releases/download/v0.3.5/ktop_v0.3.5_linux_amd64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/vladimirvivien/ktop/releases/download/v0.3.5/ktop_v0.3.5_darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/vladimirvivien/ktop/releases/download/v0.3.5/ktop_v0.3.5_linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/vladimirvivien/ktop/releases/download/v0.3.5/ktop_v0.3.5_darwin_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/vladimirvivien/ktop/releases/download/v0.3.5/ktop_v0.3.5_linux_armv7.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadKubeBurner(t *testing.T) {
	tools := MakeTools()
	name := "kube-burner"

	tool := getTool(name, tools)

	const toolVersion = "v1.6"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cloud-bulldozer/kube-burner/releases/download/v1.6/kube-burner-V1.6-Linux-x86_64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cloud-bulldozer/kube-burner/releases/download/v1.6/kube-burner-V1.6-Darwin-x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/cloud-bulldozer/kube-burner/releases/download/v1.6/kube-burner-V1.6-Linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/cloud-bulldozer/kube-burner/releases/download/v1.6/kube-burner-V1.6-Darwin-arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cloud-bulldozer/kube-burner/releases/download/v1.6/kube-burner-V1.6-Windows-x86_64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadOpenshiftInstall(t *testing.T) {
	tools := MakeTools()
	name := "openshift-install"

	tool := getTool(name, tools)

	const toolVersion = "4.13.1"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-install-linux.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-install-mac.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-install-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-install-mac-arm64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadOpenshiftCLI(t *testing.T) {
	tools := MakeTools()
	name := "oc"

	tool := getTool(name, tools)

	const toolVersion = "4.13.1"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-client-linux.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-client-mac.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-client-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-client-mac-arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.13.1/openshift-client-windows.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadAtuin(t *testing.T) {
	tools := MakeTools()
	name := "atuin"

	tool := getTool(name, tools)

	const toolVersion = "v15.0.0"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/atuinsh/atuin/releases/download/v15.0.0/atuin-v15.0.0-x86_64-unknown-linux-gnu.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/atuinsh/atuin/releases/download/v15.0.0/atuin-v15.0.0-x86_64-apple-darwin.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}

}

func Test_Copacetic(t *testing.T) {
	tools := MakeTools()
	name := "copa"

	tool := getTool(name, tools)

	const toolVersion = "v0.2.0"

	test := test{
		os:      "linux",
		arch:    arch64bit,
		version: toolVersion,
		url:     `https://github.com/project-copacetic/copacetic/releases/download/v0.2.0/copa_0.2.0_linux_amd64.tar.gz`,
	}

	t.Run(test.os+" "+test.arch+" "+test.version, func(r *testing.T) {
		got, err := tool.GetURL(test.os, test.arch, test.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.url {
			t.Errorf("want: %s, got: %s", test.url, got)
		}
	})

}

func Test_DownloadTask(t *testing.T) {
	tools := MakeTools()
	name := "task"

	tool := getTool(name, tools)

	const toolVersion = "v3.26.0"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/go-task/task/releases/download/v3.26.0/task_linux_amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     "https://github.com/go-task/task/releases/download/v3.26.0/task_linux_arm.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     "https://github.com/go-task/task/releases/download/v3.26.0/task_linux_arm64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     "https://github.com/go-task/task/releases/download/v3.26.0/task_darwin_arm64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://github.com/go-task/task/releases/download/v3.26.0/task_darwin_amd64.tar.gz",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_Download1Password(t *testing.T) {
	tools := MakeTools()
	name := "op"

	tool := getTool(name, tools)

	const toolVersion = "v2.17.0"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://cache.agilebits.com/dist/1P/op2/pkg/v2.17.0/op_linux_amd64_v2.17.0.zip",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "",
			url:     "https://cache.agilebits.com/dist/1P/op2/pkg/v2.17.0/op_linux_amd64_v2.17.0.zip",
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     "https://cache.agilebits.com/dist/1P/op2/pkg/v2.17.0/op_linux_arm_v2.17.0.zip",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     "https://cache.agilebits.com/dist/1P/op2/pkg/v2.17.0/op_linux_arm64_v2.17.0.zip",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     "https://cache.agilebits.com/dist/1P/op2/pkg/v2.17.0/op_windows_amd64_v2.17.0.zip",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_VHS(t *testing.T) {
	tools := MakeTools()
	name := "vhs"

	tool := getTool(name, tools)

	const toolVersion = "v0.5.0"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/charmbracelet/vhs/releases/download/v0.5.0/vhs_0.5.0_Linux_x86_64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/charmbracelet/vhs/releases/download/v0.5.0/vhs_Darwin_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/charmbracelet/vhs/releases/download/v0.5.0/vhs_0.5.0_Linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/charmbracelet/vhs/releases/download/v0.5.0/vhs_Darwin_arm64.tar.gz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/charmbracelet/vhs/releases/download/v0.5.0/vhs_0.5.0_Windows_x86_64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloaSkupper(t *testing.T) {
	tools := MakeTools()
	name := "skupper"

	tool := getTool(name, tools)

	const toolVersion = "1.4.2"

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/skupperproject/skupper/releases/download/1.4.2/skupper-cli-1.4.2-linux-amd64.tgz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/skupperproject/skupper/releases/download/1.4.2/skupper-cli-1.4.2-mac-amd64.tgz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/skupperproject/skupper/releases/download/1.4.2/skupper-cli-1.4.2-linux-arm64.tgz`,
		},
		{
			os:      "darwin",
			arch:    archDarwinARM64,
			version: toolVersion,
			url:     `https://github.com/skupperproject/skupper/releases/download/1.4.2/skupper-cli-1.4.2-mac-arm64.tgz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/skupperproject/skupper/releases/download/1.4.2/skupper-cli-1.4.2-linux-arm32.tgz`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/skupperproject/skupper/releases/download/1.4.2/skupper-cli-1.4.2-windows-amd64.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version, false)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
		}
	}
}
