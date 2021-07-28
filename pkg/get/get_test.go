package get

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
)

var faasCLIVersionConstraint, _ = semver.NewConstraint(">= 0.13.2")

const arch32bit = "i686"
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

func getFaaSCLIVersion(url string, expectedBinaryName string) *semver.Version {
	faasCLIURLVersionRegex := regexp.MustCompile(
		"https://github.com/openfaas/faas-cli/releases/download/" +
			semver.SemVerRegex + "/" + expectedBinaryName)
	result := faasCLIURLVersionRegex.FindStringSubmatch(url)
	version, _ := semver.NewVersion(strings.Join(result[1:], ""))
	return version
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

	gotURL, err := tool.GetURL("darwin", "", "")
	if err != nil {
		t.Fatal(err)
	}
	gotVersion := getFaaSCLIVersion(gotURL, "faas-cli-darwin")
	valid, msgs := faasCLIVersionConstraint.Validate(gotVersion)
	if !valid {
		t.Fatalf("%s failed version constraint: %v", gotURL, msgs)
	}
}

func Test_DownloadKubectl(t *testing.T) {
	tools := MakeTools()
	name := "kubectl"

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
			version: "v1.20.0",
			url:     "https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/darwin/amd64/kubectl"},
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
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

func Test_DownloadFaaSCLIArmhf(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	gotURL, err := tool.GetURL("Linux", "armv7l", "")
	if err != nil {
		t.Fatal(err)
	}
	gotVersion := getFaaSCLIVersion(gotURL, "faas-cli-armhf")
	valid, msgs := faasCLIVersionConstraint.Validate(gotVersion)
	if !valid {
		t.Fatalf("%s failed version constraint: %v", gotURL, msgs)
	}
}

func Test_DownloadFaaSCLIArm64(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	gotURL, err := tool.GetURL("Linux", "aarch64", "")
	if err != nil {
		t.Fatal(err)
	}
	gotVersion := getFaaSCLIVersion(gotURL, "faas-cli-arm64")
	valid, msgs := faasCLIVersionConstraint.Validate(gotVersion)
	if !valid {
		t.Fatalf("%s failed version constraint: %v", gotURL, msgs)
	}
}

func Test_DownloadFaaSCLIWindows(t *testing.T) {
	tools := MakeTools()
	name := "faas-cli"
	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

	gotURL, err := tool.GetURL("mingw64_nt-10.0-18362", arch64bit, "")
	if err != nil {
		t.Fatal(err)
	}
	gotVersion := getFaaSCLIVersion(gotURL, "faas-cli.exe")
	valid, msgs := faasCLIVersionConstraint.Validate(gotVersion)
	if !valid {
		t.Fatalf("%s failed version constraint: %v", gotURL, msgs)
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

	ver := "v3.8.8"

	tests := []test{
		{os: "linux",
			arch:    arch64bit,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.8.8/kustomize_v3.8.8_linux_amd64.tar.gz",
		},
		{os: "darwin",
			arch:    arch64bit,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.8.8/kustomize_v3.8.8_darwin_amd64.tar.gz",
		},
		{os: "linux",
			arch:    archARM64,
			version: ver,
			url:     "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.8.8/kustomize_v3.8.8_linux_arm64.tar.gz",
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

	const toolVersion = "v0.24.10"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_v0.24.10_Windows_x86_64.tar.gz`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_v0.24.10_Linux_x86_64.tar.gz`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_v0.24.10_Darwin_x86_64.tar.gz`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://github.com/derailed/k9s/releases/download/v0.24.10/k9s_v0.24.10_Linux_arm.tar.gz`,
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

	const toolVersion = "0.7.11"

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

	const toolVersion = "1.0.0"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/1.0.0/terraform_1.0.0_windows_amd64.zip`,
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/1.0.0/terraform_1.0.0_linux_amd64.zip`,
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/1.0.0/terraform_1.0.0_darwin_amd64.zip`,
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			url:     `https://releases.hashicorp.com/terraform/1.0.0/terraform_1.0.0_linux_arm.zip`,
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

	const toolVersion = "1.6.1"

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
		got, err := tool.GetURL(tc.os, "", tc.version)
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
			arch: "armv7l",
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

			got, err := tool.GetURL(tc.os, tc.arch, "")
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
			version: "0.0.21",
			url:     `https://github.com/nats-io/natscli/releases/download/0.0.21/nats-0.0.21-windows-amd64.zip`,
		},
		{
			os:      "linux",
			arch:    "amd64",
			version: "0.0.21",
			url:     `https://github.com/nats-io/natscli/releases/download/0.0.21/nats-0.0.21-linux-amd64.zip`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.0.21",
			url:     `https://github.com/nats-io/natscli/releases/download/0.0.21/nats-0.0.21-linux-arm64.zip`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: "0.0.21",
			url:     `https://github.com/nats-io/natscli/releases/download/0.0.21/nats-0.0.21-linux-arm6.zip`,
		},
		{
			os:      "linux",
			arch:    "armv7l",
			version: "0.0.21",
			url:     `https://github.com/nats-io/natscli/releases/download/0.0.21/nats-0.0.21-linux-arm7.zip`,
		},
		{
			os:      "darwin",
			arch:    "amd64",
			version: "0.0.21",
			url:     `https://github.com/nats-io/natscli/releases/download/0.0.21/nats-0.0.21-darwin-amd64.zip`,
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

func Test_DownloadLinkerd(t *testing.T) {
	tools := MakeTools()
	name := "linkerd2"

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
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			version: "v1.8.6",
			url:     "https://github.com/argoproj/argo-cd/releases/download/v1.8.6/argocd-windows-amd64.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: "v1.8.6",
			url:     "https://github.com/argoproj/argo-cd/releases/download/v1.8.6/argocd-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: "v1.8.6",
			url:     "https://github.com/argoproj/argo-cd/releases/download/v1.8.6/argocd-darwin-amd64"},
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

func Test_DownloadNerdctl(t *testing.T) {
	tools := MakeTools()
	name := "nerdctl"

	var tool *Tool
	for _, target := range tools {
		if name == target.Name {
			tool = &target
			break
		}
	}

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
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			version: "0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Windows_x86_64.zip`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Linux_x86_64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "0.17.2",
			url:     `https://github.com/tektoncd/cli/releases/download/v0.17.2/tkn_0.17.2_Darwin_x86_64.tar.gz`,
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
		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: "0.8.3",
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.8.3/inlets-pro.exe`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.8.3",
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.8.3/inlets-pro`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.8.3",
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.8.3/inlets-pro-arm64`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "0.8.3",
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.8.3/inlets-pro-armhf`,
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: "0.8.3",
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.8.3/inlets-pro-armhf`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "0.8.3",
			url:     `https://github.com/inlets/inlets-pro/releases/download/0.8.3/inlets-pro-darwin`,
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
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			version: "0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_Linux-64bit.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_Linux-ARM.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_Linux-ARM64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "0.17.2",
			url:     `https://github.com/aquasecurity/trivy/releases/download/v0.17.2/trivy_0.17.2_macOS-64bit.tar.gz`,
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

func Test_DownloadFluxCli(t *testing.T) {
	tools := MakeTools()
	name := "flux"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_linux_arm.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_linux_arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "0.13.4",
			url:     `https://github.com/fluxcd/flux2/releases/download/v0.13.4/flux_0.13.4_darwin_amd64.tar.gz`,
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
func Test_DownloadPolarisCli(t *testing.T) {
	tools := MakeTools()
	name := "polaris"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/3.2.1/polaris_3.2.1_darwin_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/3.2.1/polaris_3.2.1_linux_amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/3.2.1/polaris_3.2.1_linux_arm64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "3.2.1",
			url:     `https://github.com/FairwindsOps/polaris/releases/download/3.2.1/polaris_3.2.1_linux_armv7.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			url:     `https://get.helm.sh/helm-3.5.4-darwin-amd64.tar.gz`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			version: "0.2.1",
			url:     `https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.2.1/argocd-autopilot-linux-amd64.tar.gz`,
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.2.1",
			url:     `https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.2.1/argocd-autopilot-linux-arm64.tar.gz`,
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "0.2.1",
			url:     `https://github.com/argoproj-labs/argocd-autopilot/releases/download/v0.2.1/argocd-autopilot-darwin-amd64.tar.gz`,
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

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			arch:    "armv7l",
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
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
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
	prefix := "https://github.com/" + tool.Owner + "/" + tool.Repo + "/releases/download/jq-" + tool.Version + "/"

	tests := []test{
		{
			os:   "darwin",
			arch: arch64bit,
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
			got, err := tool.GetURL(tc.os, tc.arch, tool.Version)
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
			version: "1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-darwin-amd64`,
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-darwin-arm64`,
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-linux-amd64`,
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "1.0.0",
			url:     `https://github.com/sigstore/cosign/releases/download/v1.0.0/cosign-windows-amd64.exe`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.url {
				t.Errorf("want: %s, got: %s", tc.url, got)
			}
		})
	}

}
