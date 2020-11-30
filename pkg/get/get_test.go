package get

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

const faasCLIVersion = "0.12.19"
const arch64bit = "x86_64"
const archARM7 = "armv7l"
const archARM64 = "aarch64"

type test struct {
	os      string
	arch    string
	version string
	url     string
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

func Test_DownloadFaaSCLIDarwin(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("darwin", "", "")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/openfaas/faas-cli/releases/download/" + faasCLIVersion + "/faas-cli-darwin"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadKubectlDarwin(t *testing.T) {
	tools := MakeTools()
	name := "kubectl"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("darwin", arch64bit, tool.Version)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/darwin/amd64/kubectl"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadKubectlLinux(t *testing.T) {
	tools := MakeTools()
	name := "kubectl"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("linux", arch64bit, tool.Version)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadKubectx(t *testing.T) {
	tools := MakeTools()
	name := "kubectx"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("linux", arch64bit, tool.Version)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/ahmetb/kubectx/releases/download/v0.9.1/kubectx"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadKubens(t *testing.T) {
	tools := MakeTools()
	name := "kubens"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("linux", arch64bit, tool.Version)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/ahmetb/kubectx/releases/download/v0.9.1/kubens"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadArmhf(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("Linux", "armv7l", "")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/openfaas/faas-cli/releases/download/" + faasCLIVersion + "/faas-cli-armhf"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadArm64(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("Linux", "aarch64", "")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/openfaas/faas-cli/releases/download/" + faasCLIVersion + "/faas-cli-arm64"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadWindows(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("mingw64_nt-10.0-18362", arch64bit, "")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://github.com/openfaas/faas-cli/releases/download/" + faasCLIVersion + "/faas-cli.exe"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadHelmDarwin(t *testing.T) {
	tools := MakeTools()
	name := "helm"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("darwin", arch64bit, tool.Version)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://get.helm.sh/helm-v3.2.4-darwin-amd64.tar.gz"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadHelmLinux(t *testing.T) {
	tools := MakeTools()
	name := "helm"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("linux", arch64bit, tool.Version)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadHelmWindows(t *testing.T) {
	tools := MakeTools()
	name := "helm"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	got, err := tool.GetURL("mingw64_nt-10.0-18362", arch64bit, tool.Version)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://get.helm.sh/helm-v3.2.4-windows-amd64.zip"
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func Test_DownloadKubeseal(t *testing.T) {
	tools := MakeTools()
	name := "kubeseal"

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v0.12.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/kubeseal.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: "v0.12.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/kubeseal-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v0.12.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/kubeseal-darwin-amd64"},
		{os: "linux",
			arch:    "armv7l",
			version: "v0.12.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/kubeseal-arm"},
		{os: "linux",
			arch:    "arm64",
			version: "v0.12.4",
			url:     "https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.12.4/kubeseal-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v0.8.1",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.8.1/kind-windows-amd64"},
		{os: "linux",
			arch:    arch64bit,
			version: "v0.8.1",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.8.1/kind-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v0.8.1",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.8.1/kind-darwin-amd64"},
		{os: "linux",
			arch:    "armv7l",
			version: "v0.8.1",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.8.1/kind-linux-arm"},
		{os: "linux",
			arch:    "aarch64",
			version: "v0.8.1",
			url:     "https://github.com/kubernetes-sigs/kind/releases/download/v0.8.1/kind-linux-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "v3.0.0",
			url:     "https://github.com/rancher/k3d/releases/download/v3.0.0/k3d-windows-amd64"},
		{os: "linux",
			arch:    arch64bit,
			version: "v3.0.0",
			url:     "https://github.com/rancher/k3d/releases/download/v3.0.0/k3d-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v3.0.0",
			url:     "https://github.com/rancher/k3d/releases/download/v3.0.0/k3d-darwin-amd64"},
		{os: "linux",
			arch:    "armv7l",
			version: "v3.0.0",
			url:     "https://github.com/rancher/k3d/releases/download/v3.0.0/k3d-linux-arm"},
		{os: "linux",
			arch:    "aarch64",
			version: "v3.0.0",
			url:     "https://github.com/rancher/k3d/releases/download/v3.0.0/k3d-linux-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

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
			arch:    "armv7l",
			version: "0.9.2",
			url:     "https://github.com/alexellis/k3sup/releases/download/0.9.2/k3sup-armhf"},
		{os: "linux",
			arch:    "aarch64",
			version: "0.9.2",
			url:     "https://github.com/alexellis/k3sup/releases/download/0.9.2/k3sup-arm64"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.5.4",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.5.4/inletsctl.exe.tgz"},
		{os: "linux",
			arch:    arch64bit,
			version: "0.5.4",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.5.4/inletsctl.tgz"},
		{os: "darwin",
			arch:    arch64bit,
			version: "0.5.4",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.5.4/inletsctl-darwin.tgz"},
		{os: "linux",
			arch:    "armv6l",
			version: "0.5.4",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.5.4/inletsctl-armhf.tgz"},
		{os: "linux",
			arch:    "arm64",
			version: "0.5.4",
			url:     "https://github.com/inlets/inletsctl/releases/download/0.5.4/inletsctl-arm64.tgz"},
	}
	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}

func Test_DownloadKubebuilder(t *testing.T) {
	tools := MakeTools()
	name := "kubebuilder"

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	tests := []test{
		{os: "darwin",
			arch:    arch64bit,
			version: "2.3.1",
			url:     "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v2.3.1/kubebuilder_2.3.1_darwin_amd64.tar.gz"},
		{os: "linux",
			arch:    arch64bit,
			version: "2.3.1",
			url:     "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v2.3.1/kubebuilder_2.3.1_linux_amd64.tar.gz"},
		{os: "linux",
			arch:    "arm64",
			version: "2.3.1",
			url:     "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v2.3.1/kubebuilder_2.3.1_linux_arm64.tar.gz"},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	ver := "kustomize/v3.8.1"

	tests := []test{
		{os: "linux",
			arch:    arch64bit,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v3.8.1/kustomize_v3.8.1_linux_amd64.tar.gz",
		},
		{os: "darwin",
			arch:    arch64bit,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v3.8.1/kustomize_v3.8.1_darwin_amd64.tar.gz",
		},
		{os: "linux",
			arch:    "arm64",
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v3.8.1/kustomize_v3.8.1_.tar.gz",
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.url {
			t.Errorf("want: %s, got: %s", tc.url, got)
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
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	const toolVersion = "v0.21.7"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.21.7/k9s_Windows_x86_64.tar.gz`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.21.7/k9s_Linux_x86_64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.21.7/k9s_Darwin_x86_64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.21.7/k9s_Linux_arm.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	const toolVersion = "0.6.27"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.6.27/civo-0.6.27-windows-amd64.zip`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.6.27/civo-0.6.27-linux-amd64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.6.27/civo-0.6.27-darwin-amd64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/civo/cli/releases/download/v0.6.27/civo-0.6.27-linux-arm.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	const toolVersion = "0.13.1"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/0.13.1/terraform_0.13.1_windows_amd64.zip`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/0.13.1/terraform_0.13.1_linux_amd64.zip`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/0.13.1/terraform_0.13.1_darwin_amd64.zip`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/0.13.1/terraform_0.13.1_linux_arm.zip`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	const toolVersion = "1.0.0"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.0.0/gh_1.0.0_windows_amd64.zip`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.0.0/gh_1.0.0_linux_amd64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.0.0/gh_1.0.0_macOS_amd64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			url:     `https://github.com/cli/cli/releases/download/v1.0.0/gh_1.0.0_linux_arm64.tar.gz`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	const toolVersion = "0.14.2"

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
		got, err := tool.GetURL(tc.os, "", tc.version)
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

	const toolVersion = "0.4.2"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.4.2/buildx-v0.4.2.windows-amd64.exe`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.4.2/buildx-v0.4.2.linux-amd64`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.4.2/buildx-v0.4.2.darwin-amd64`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/docker/buildx/releases/download/v0.4.2/buildx-v0.4.2.linux-arm-v7`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	const toolVersion = "v0.132.1"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/roboll/helmfile/releases/download/v0.132.1/helmfile_windows_amd64.exe`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/roboll/helmfile/releases/download/v0.132.1/helmfile_linux_amd64`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/roboll/helmfile/releases/download/v0.132.1/helmfile_darwin_amd64`,
		},
	}

	for _, tc := range tests {
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
