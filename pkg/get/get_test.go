package get

import "testing"

const faasCLIVersion = "0.12.4"
const arch64bit = "x86_64"

func Test_DownloadDarwin(t *testing.T) {
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
	want := "https://github.com/ahmetb/kubectx/releases/download/v0.9.0/kubectx"
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

	type test struct {
		os      string
		arch    string
		version string
		url     string
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
			arch:    "arm",
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
			t.Fatalf("want: %s, got: %s", tc.url, got)
		}
	}
}
