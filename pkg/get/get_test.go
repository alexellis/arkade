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
const baseURL string = "https://github.com/%s/%s/releases/download/%s/%s"

type test struct {
	os      string
	arch    string
	version string
	binary  string
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
	v := "v1.20.0"

	const overrideBaseURL string = "https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s"

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
			version: v,
			binary:  "darwin/amd64/kubectl"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "linux/amd64/kubectl"},
		{os: "linux",
			arch:    archARM64,
			version: v,
			binary:  "linux/arm64/kubectl"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(overrideBaseURL, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
	v := "v0.12.4"

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
			version: v,
			binary:  "kubeseal.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "kubeseal-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "kubeseal-darwin-amd64"},
		{os: "linux",
			arch:    "armv7l",
			version: v,
			binary:  "kubeseal-arm"},
		{os: "linux",
			arch:    "arm64",
			version: v,
			binary:  "kubeseal-arm64"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadKind(t *testing.T) {
	tools := MakeTools()
	name := "kind"
	v := "v0.8.1"

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
			version: v,
			binary:  "kind-windows-amd64"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "kind-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "kind-darwin-amd64"},
		{os: "linux",
			arch:    "armv7l",
			version: v,
			binary:  "kind-linux-arm"},
		{os: "linux",
			arch:    "aarch64",
			version: v,
			binary:  "kind-linux-arm64"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadK3d(t *testing.T) {
	tools := MakeTools()
	name := "k3d"
	v := "v3.0.0"

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
			version: v,
			binary:  "k3d-windows-amd64"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "k3d-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "k3d-darwin-amd64"},
		{os: "linux",
			arch:    "armv7l",
			version: v,
			binary:  "k3d-linux-arm"},
		{os: "linux",
			arch:    "aarch64",
			version: v,
			binary:  "k3d-linux-arm64"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadK3sup(t *testing.T) {
	tools := MakeTools()
	name := "k3sup"
	v := "0.9.2"

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
			version: v,
			binary:  "k3sup.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "k3sup"},
		{os: "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "k3sup-darwin"},
		{os: "linux",
			arch:    "armv7l",
			version: v,
			binary:  "k3sup-armhf"},
		{os: "linux",
			arch:    "aarch64",
			version: v,
			binary:  "k3sup-arm64"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadInletsctl(t *testing.T) {
	tools := MakeTools()
	name := "inletsctl"
	v := "0.5.4"

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
			version: v,
			binary:  "inletsctl.exe.tgz"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "inletsctl.tgz"},
		{os: "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "inletsctl-darwin.tgz"},
		{os: "linux",
			arch:    "armv6l",
			version: v,
			binary:  "inletsctl-armhf.tgz"},
		{os: "linux",
			arch:    "arm64",
			version: v,
			binary:  "inletsctl-arm64.tgz"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadKubebuilder(t *testing.T) {
	tools := MakeTools()
	name := "kubebuilder"
	v := "3.1.0"

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
			version: v,
			binary:  "kubebuilder_darwin_amd64"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "kubebuilder_linux_amd64"},
		{os: "linux",
			arch:    "arm64",
			version: v,
			binary:  "kubebuilder_linux_arm64"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  fmt.Sprintf("kustomize_%s_linux_amd64.tar.gz", ver),
		},
		{os: "darwin",
			arch:    arch64bit,
			version: ver,
			binary:  fmt.Sprintf("kustomize_%s_darwin_amd64.tar.gz", ver),
		},
		{os: "linux",
			arch:    archARM64,
			version: ver,
			binary:  fmt.Sprintf("kustomize_%s_linux_arm64.tar.gz", ver),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "kustomize%2F"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadDigitalOcean(t *testing.T) {
	tools := MakeTools()
	name := "doctl"

	tool := getTool(name, tools)

	const toolVersion = "1.46.0"

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("doctl-%s-windows-amd64.zip", toolVersion)},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("doctl-%s-linux-amd64.tar.gz", toolVersion)},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("doctl-%s-darwin-amd64.tar.gz", toolVersion)},
		// this asserts that we can build a URL for ARM processors, but no asset exists and will yield a 404
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			binary:  fmt.Sprintf("doctl-%s-linux-.tar.gz", toolVersion)},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  fmt.Sprintf("k9s_%s_Windows_x86_64.tar.gz", toolVersion),
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("k9s_%s_Linux_x86_64.tar.gz", toolVersion),
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("k9s_%s_Darwin_x86_64.tar.gz", toolVersion),
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			binary:  fmt.Sprintf("k9s_%s_Linux_arm.tar.gz", toolVersion),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  fmt.Sprintf("civo-%s-windows-amd64.zip", toolVersion),
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("civo-%s-linux-amd64.tar.gz", toolVersion),
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("civo-%s-darwin-amd64.tar.gz", toolVersion),
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			binary:  fmt.Sprintf("civo-%s-linux-arm.tar.gz", toolVersion),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadTerraform(t *testing.T) {
	tools := MakeTools()
	name := "terraform"
	const toolVersion = "1.0.0"
	const overrideBaseURL string = "https://releases.hashicorp.com/terraform/%s/%s"

	tool := getTool(name, tools)

	tests := []test{
		{os: "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("terraform_%s_windows_amd64.zip", toolVersion),
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("terraform_%s_linux_amd64.zip", toolVersion),
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("terraform_%s_darwin_amd64.zip", toolVersion),
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			binary:  fmt.Sprintf("terraform_%s_linux_arm.zip", toolVersion),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(overrideBaseURL, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  fmt.Sprintf("gh_%s_windows_amd64.zip", toolVersion),
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("gh_%s_linux_amd64.tar.gz", toolVersion),
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("gh_%s_macOS_amd64.tar.gz", toolVersion),
		},
		{os: "linux",
			arch:    archARM64,
			version: toolVersion,
			binary:  fmt.Sprintf("gh_%s_linux_arm64.tar.gz", toolVersion),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  fmt.Sprintf("pack-v%s-windows.zip", toolVersion),
		},
		{os: "darwin",
			version: toolVersion,
			binary:  fmt.Sprintf("pack-v%s-macos.tgz", toolVersion),
		},
		{os: "linux",
			version: toolVersion,
			binary:  fmt.Sprintf("pack-v%s-linux.tgz", toolVersion),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  fmt.Sprintf("buildx-v%s.windows-amd64.exe", toolVersion),
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("buildx-v%s.linux-amd64", toolVersion),
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			binary:  fmt.Sprintf("buildx-v%s.darwin-amd64", toolVersion),
		},
		{os: "linux",
			arch:    archARM7,
			version: toolVersion,
			binary:  fmt.Sprintf("buildx-v%s.linux-arm-v7", toolVersion),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  "helmfile_windows_amd64.exe",
		},
		{os: "linux",
			arch:    arch64bit,
			version: toolVersion,
			binary:  "helmfile_linux_amd64",
		},
		{os: "darwin",
			arch:    arch64bit,
			version: toolVersion,
			binary:  "helmfile_darwin_amd64",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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
			binary:  "opa_windows_amd64.exe",
		},
		{os: "linux",
			version: toolVersion,
			binary:  "opa_linux_amd64",
		},
		{os: "darwin",
			version: toolVersion,
			binary:  "opa_darwin_amd64",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
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

	const overrideBaseURL string = "https://dl.min.io/client/mc/release/%s"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:     "ming",
			arch:   "amd64",
			binary: "windows-amd64/mc.exe",
		},
		{
			os:     "linux",
			arch:   "amd64",
			binary: "linux-amd64/mc",
		},
		{
			os:     "linux",
			arch:   "arm",
			binary: "linux-arm/mc",
		},
		{
			os:     "linux",
			arch:   "armv6l",
			binary: "linux-arm/mc",
		},
		{
			os:     "linux",
			arch:   "armv7l",
			binary: "linux-arm/mc",
		},
		{
			os:     "linux",
			arch:   archARM64,
			binary: "linux-arm64/mc",
		},
		{
			os:     "darwin",
			arch:   "amd64",
			binary: "darwin-amd64/mc",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(t *testing.T) {
			expectedURL := fmt.Sprintf(overrideBaseURL, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, "")
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadNats(t *testing.T) {
	tools := MakeTools()
	name := "nats"
	v := "0.0.21"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    "amd64",
			version: v,
			binary:  fmt.Sprintf("nats-%s-windows-amd64.zip", v),
		},
		{
			os:      "linux",
			arch:    "amd64",
			version: v,
			binary:  fmt.Sprintf("nats-%s-linux-amd64.zip", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("nats-%s-linux-arm64.zip", v),
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: v,
			binary:  fmt.Sprintf("nats-%s-linux-arm6.zip", v),
		},
		{
			os:      "linux",
			arch:    "armv7l",
			version: v,
			binary:  fmt.Sprintf("nats-%s-linux-arm7.zip", v),
		},
		{
			os:      "darwin",
			arch:    "amd64",
			version: v,
			binary:  fmt.Sprintf("nats-%s-darwin-amd64.zip", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadLinkerd(t *testing.T) {
	tools := MakeTools()
	name := "linkerd2"
	v := "stable-2.9.1"

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
			version: v,
			binary:  "linkerd2-cli-stable-2.9.1-windows.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "linkerd2-cli-stable-2.9.1-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "linkerd2-cli-stable-2.9.1-darwin"},
		{os: "linux",
			arch:    archARM64,
			version: v,
			binary:  "linkerd2-cli-stable-2.9.1-linux-arm64"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadArgocd(t *testing.T) {
	tools := MakeTools()
	name := "argocd"
	v := "v1.8.6"

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
			version: v,
			binary:  "argocd-windows-amd64.exe"},
		{os: "linux",
			arch:    arch64bit,
			version: v,
			binary:  "argocd-linux-amd64"},
		{os: "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "argocd-darwin-amd64"},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadNerdctl(t *testing.T) {
	tools := MakeTools()
	name := "nerdctl"
	v := "0.7.2"

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
			version: v,
			binary:  fmt.Sprintf("nerdctl-%s-linux-amd64.tar.gz", v),
		},
		{os: "linux",
			arch:    archARM7,
			version: v,
			binary:  fmt.Sprintf("nerdctl-%s-linux-arm-v7.tar.gz", v),
		},
		{os: "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("nerdctl-%s-linux-arm64.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadIstioCtl(t *testing.T) {
	tools := MakeTools()
	name := "istioctl"
	v := "1.9.1"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    "amd64",
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-win.zip", v),
		},
		{
			os:      "linux",
			arch:    "x86_64",
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-linux-amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    "amd64",
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-linux-amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    "arm",
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-linux-armv7.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-linux-armv7.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-linux-armv7.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-linux-arm64.tar.gz", v),
		},
		{
			os:      "darwin",
			arch:    "amd64",
			version: v,
			binary:  fmt.Sprintf("istioctl-%s-osx.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(t *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadTektonCli(t *testing.T) {
	tools := MakeTools()
	name := "tkn"
	v := "0.17.2"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("tkn_%s_Windows_x86_64.zip", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("tkn_%s_Linux_x86_64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("tkn_%s_Linux_arm64.tar.gz", v),
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("tkn_%s_Darwin_x86_64.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadInfluxCli(t *testing.T) {
	tools := MakeTools()
	name := "influx"
	v := "2.0.7"

	const overrideBaseURL string = "https://dl.influxdata.com/influxdb/releases/%s"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "windows",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("influxdb2-client-%s-windows-amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("influxdb2-client-%s-linux-amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("influxdb2-client-%s-linux-arm64.tar.gz", v),
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("influxdb2-client-%s-darwin-amd64.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(overrideBaseURL, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadInletsProCli(t *testing.T) {
	tools := MakeTools()
	name := "inlets-pro"
	v := "0.8.3"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: v,
			binary:  "inlets-pro.exe",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  "inlets-pro",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  "inlets-pro-arm64",
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: v,
			binary:  "inlets-pro-armhf",
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: v,
			binary:  "inlets-pro-armhf",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "inlets-pro-darwin",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadKim(t *testing.T) {
	tools := MakeTools()
	name := "kim"
	v := "v0.1.0-alpha.12"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "ming",
			arch:    arch64bit,
			version: v,
			binary:  "kim-windows-amd64.exe",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  "kim-linux-amd64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  "kim-linux-arm64",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "kim-darwin-amd64",
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("Download for: %s %s %s", tc.os, tc.arch, tc.version), func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				r.Fatal(err)
			}
			if got != expectedURL {
				r.Errorf("\nwant: %s\ngot:  %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadTrivyCli(t *testing.T) {
	tools := MakeTools()
	name := "trivy"
	v := "0.17.2"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("trivy_%s_Linux-64bit.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: v,
			binary:  fmt.Sprintf("trivy_%s_Linux-ARM.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("trivy_%s_Linux-ARM64.tar.gz", v),
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("trivy_%s_macOS-64bit.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadFluxCli(t *testing.T) {
	tools := MakeTools()
	name := "flux"
	v := "0.13.4"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("flux_%s_linux_amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: v,
			binary:  fmt.Sprintf("flux_%s_linux_arm.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("flux_%s_linux_arm64.tar.gz", v),
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("flux_%s_darwin_amd64.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadPolarisCli(t *testing.T) {
	tools := MakeTools()
	name := "polaris"
	v := "3.2.1"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("polaris_%s_darwin_amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("polaris_%s_linux_amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("polaris_%s_linux_arm64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: v,
			binary:  fmt.Sprintf("polaris_%s_linux_armv7.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadHelm(t *testing.T) {
	tools := MakeTools()
	name := "helm"
	v := "3.5.4"

	const overrideBaseURL string = "https://get.helm.sh/%s"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("helm-%s-linux-amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: v,
			binary:  fmt.Sprintf("helm-%s-linux-arm.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("helm-%s-linux-arm64.tar.gz", v),
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("helm-%s-darwin-amd64.tar.gz", v),
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("helm-%s-darwin-amd64.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(overrideBaseURL, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				r.Fatal(err)
			}
			if got != expectedURL {
				r.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadArgoCDAutopilotCli(t *testing.T) {
	tools := MakeTools()
	name := "argocd-autopilot"
	v := "0.2.13"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  "argocd-autopilot-linux-amd64.tar.gz",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  "argocd-autopilot-linux-arm64.tar.gz",
		},
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  "argocd-autopilot-darwin-amd64.tar.gz",
		},
	}

	for _, tc := range tests {
		expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

		got, err := tool.GetURL(tc.os, tc.arch, tc.version)
		if err != nil {
			t.Fatal(err)
		}
		if got != expectedURL {
			t.Errorf("want: %s, got: %s", expectedURL, got)
		}
	}
}

func Test_DownloadNovaCli(t *testing.T) {
	tools := MakeTools()
	name := "nova"
	v := "2.3.2"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("nova_%s_darwin_amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("nova_%s_linux_amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: v,
			binary:  fmt.Sprintf("nova_%s_linux_arm64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: v,
			binary:  fmt.Sprintf("nova_%s_linux_armv7.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadKubetailCli(t *testing.T) {
	tools := MakeTools()
	name := "kubetail"

	const overrideBaseURL string = "https://raw.githubusercontent.com/johanhaleby/kubetail/%s/%s"

	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: "1.6.13",
			binary:  "kubetail",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "1.6.13",
			binary:  "kubetail",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "1.6.13",
			binary:  "kubetail",
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "1.6.13",
			binary:  "kubetail",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(overrideBaseURL, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
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
			binary:  "kgctl-darwin-amd64",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "0.3.0",
			binary:  "kgctl-darwin-arm64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.3.0",
			binary:  "kgctl-linux-amd64",
		},
		{
			os:      "linux",
			arch:    archARM7,
			version: "0.3.0",
			binary:  "kgctl-linux-arm",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "0.3.0",
			binary:  "kgctl-linux-arm64",
		},
		{
			os:      "mingw64_nt-10.0-18362",
			arch:    arch64bit,
			version: "0.3.0",
			binary:  "kgctl-windows-amd64",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
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
			binary:  "metal-darwin-amd64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.6.0-alpha2",
			binary:  "metal-linux-amd64",
		},
		{
			os:      "linux",
			arch:    "aarch64",
			version: "0.6.0-alpha2",
			binary:  "metal-linux-arm64",
		},
		{
			os:      "linux",
			arch:    "armv7l",
			version: "0.6.0-alpha2",
			binary:  "metal-linux-armv7",
		},
		{
			os:      "linux",
			arch:    "armv6l",
			version: "0.6.0-alpha2",
			binary:  "metal-linux-armv6",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "0.6.0-alpha2",
			binary:  "metal-windows-amd64.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
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
			binary:  "porter-darwin-amd64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.38.4",
			binary:  "porter-linux-amd64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.38.4",
			binary:  "porter-windows-amd64.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadJq(t *testing.T) {
	tools := MakeTools()
	name := "jq"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: tool.Version,
			binary:  "jq-osx-amd64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: tool.Version,
			binary:  "jq-linux64",
		},
		{
			os:      "linux",
			arch:    arch32bit,
			version: tool.Version,
			binary:  "jq-linux32",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: tool.Version,
			binary:  "jq-win64.exe",
		},
		{
			os:      "ming",
			arch:    arch32bit,
			version: tool.Version,
			binary:  "jq-win32.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "jq-"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tool.Version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
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
			binary:  "cosign-darwin-amd64",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "1.0.0",
			binary:  "cosign-darwin-arm64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "1.0.0",
			binary:  "cosign-linux-amd64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "1.0.0",
			binary:  "cosign-windows-amd64.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want:\n%s\ngot:\n%s", expectedURL, got)
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
			binary:  fmt.Sprintf("kanister_%s_darwin_amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("kanister_%s_linux_amd64.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadKubestr(t *testing.T) {
	tools := MakeTools()
	name := "kubestr"
	v := "v0.4.17"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("kubestr-%s-darwin-amd64.tar.gz", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("kubestr-%s-linux-amd64.tar.gz", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
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
			binary:  fmt.Sprintf("k10multicluster_%s_macOS_amd64", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("k10multicluster_%s_linux_amd64", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}

func Test_DownloadK10tools(t *testing.T) {
	tools := MakeTools()
	name := "k10tools"
	v := "4.0.6"
	tool := getTool(name, tools)

	tests := []test{
		{
			os:      "darwin",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("k10tools_%s_macOS_amd64", v),
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: v,
			binary:  fmt.Sprintf("k10tools_%s_linux_amd64", v),
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want:\n%s\ngot:\n%s", expectedURL, got)
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
			version: "0.3.0",
			binary:  "rekor-cli-darwin-amd64",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "0.3.0",
			binary:  "rekor-cli-darwin-arm64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "0.3.0",
			binary:  "rekor-cli-linux-amd64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "0.3.0",
			binary:  "rekor-cli-windows-amd64.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {
			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, "v"+tc.version, tc.binary)

			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
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
			binary:  "tfsec-darwin-amd64",
		},
		{
			os:      "darwin",
			arch:    archARM64,
			version: "v0.57.1",
			binary:  "tfsec-darwin-arm64",
		},
		{
			os:      "linux",
			arch:    arch64bit,
			version: "v0.57.1",
			binary:  "tfsec-linux-amd64",
		},
		{
			os:      "linux",
			arch:    archARM64,
			version: "v0.57.1",
			binary:  "tfsec-linux-arm64",
		},
		{
			os:      "ming",
			arch:    arch64bit,
			version: "v0.57.1",
			binary:  "tfsec-windows-amd64.exe",
		},
	}

	for _, tc := range tests {
		t.Run(tc.os+" "+tc.arch+" "+tc.version, func(r *testing.T) {

			expectedURL := fmt.Sprintf(baseURL, tool.Owner, tool.Repo, tc.version, tc.binary)
			got, err := tool.GetURL(tc.os, tc.arch, tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if got != expectedURL {
				t.Errorf("want: %s, got: %s", expectedURL, got)
			}
		})
	}
}
