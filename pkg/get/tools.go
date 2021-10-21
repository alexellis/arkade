package get

import (
	"strings"
)

func (t Tools) Len() int { return len(t) }

func (t Tools) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t Tools) Less(i, j int) bool {
	var ti = t[i]
	var tj = t[j]
	var tiNameLower = strings.ToLower(ti.Name)
	var tjNameLower = strings.ToLower(tj.Name)
	if tiNameLower == tjNameLower {
		return ti.Name < tj.Name
	}
	return tiNameLower < tjNameLower
}

type Tools []Tool

func MakeTools() Tools {
	tools := []Tool{}

	tools = append(tools,
		Tool{
			Owner:       "openfaas",
			Repo:        "faas-cli",
			Name:        "faas-cli",
			Description: "Official CLI for OpenFaaS.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}.exe
{{- else if eq .OS "darwin" -}}
{{.Name}}-darwin
{{- else if eq .Arch "armv6l" -}}
{{.Name}}-armhf
{{- else if eq .Arch "armv7l" -}}
{{.Name}}-armhf
{{- else if eq .Arch "aarch64" -}}
{{.Name}}-arm64
{{- else -}}
{{.Name}}
{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "helm",
			Repo:        "helm",
			Name:        "helm",
			Version:     "v3.5.2",
			Description: "The Kubernetes Package Manager: Think of it like apt/yum/homebrew for Kubernetes.",
			URLTemplate: `{{$arch := "amd64"}}

{{- if eq .Arch "armv7l" -}}
{{$arch = "arm"}}
{{- end -}}

{{- if eq .OS "linux" -}}
	{{- if eq .Arch "aarch64" -}}
	{{$arch = "arm64"}}
	{{- end -}}
{{- end -}}

{{$os := .OS}}
{{$ext := "tar.gz"}}

{{ if HasPrefix .OS "ming" -}}
{{$os = "windows"}}
{{$ext = "zip"}}
{{- end -}}

https://get.helm.sh/helm-{{.Version}}-{{$os}}-{{$arch}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "roboll",
			Repo:        "helmfile",
			Name:        "helmfile",
			Description: "Deploy Kubernetes Helm Charts",
			BinaryTemplate: `{{$arch := "386"}}
	{{- if eq .Arch "x86_64" -}}
	{{$arch = "amd64"}}
	{{- end -}}

	{{$os := .OS}}
	{{$ext := ""}}

	{{ if HasPrefix .OS "ming" -}}
	{{$os = "windows"}}
	{{$ext = ".exe"}}
	{{- end -}}

helmfile_{{$os}}_{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "stedolan",
			Repo:        "jq",
			Name:        "jq",
			Description: "jq is a lightweight and flexible command-line JSON processor",
			BinaryTemplate: `{{$arch := "arm"}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "64"}}
{{- else -}}
{{$arch = "32"}}
{{- end -}}

{{$ext := ""}}
{{$os := .OS}}

{{ if HasPrefix .OS "ming" -}}
{{$ext = ".exe"}}
{{$os = "win"}}
{{- else if eq .OS "darwin" -}}
{{$os = "osx-amd"}}
{{- end -}}

{{.Version}}/jq-{{$os}}{{$arch}}{{$ext}}`,
		})

	// https://storage.googleapis.com/kubernetes-release/release/v1.22.2/bin/darwin/amd64/kubectl
	tools = append(tools,
		Tool{
			Owner:       "kubernetes",
			Repo:        "kubernetes",
			Name:        "kubectl",
			Version:     "v1.22.2",
			Description: "Run commands against Kubernetes clusters",
			URLTemplate: `{{$arch := "arm"}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- end -}}

{{- if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- end -}}

{{$ext := ""}}
{{$os := .OS}}

{{ if HasPrefix .OS "ming" -}}
{{$ext = ".exe"}}
{{$os = "windows"}}
{{- end -}}

https://storage.googleapis.com/kubernetes-release/release/{{.Version}}/bin/{{$os}}/{{$arch}}/kubectl{{$ext}}`})

	tools = append(tools,
		Tool{
			Owner:          "ahmetb",
			Repo:           "kubectx",
			Name:           "kubectx",
			Version:        "v0.9.1",
			Description:    "Faster way to switch between clusters.",
			BinaryTemplate: `kubectx`,
			NoExtension:    true,
		})

	tools = append(tools,
		Tool{
			Owner:          "ahmetb",
			Repo:           "kubectx",
			Name:           "kubens",
			Version:        "v0.9.1",
			Description:    "Switch between Kubernetes namespaces smoothly.",
			BinaryTemplate: `kubens`,
			NoExtension:    true,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes-sigs",
			Repo:        "kind",
			Name:        "kind",
			Description: "Run local Kubernetes clusters using Docker container nodes.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}-windows-amd64
{{- else if eq .OS "darwin" -}}
{{.Name}}-darwin-amd64
{{- else if eq .Arch "armv6l" -}}
{{.Name}}-linux-arm
{{- else if eq .Arch "armv7l" -}}
{{.Name}}-linux-arm
{{- else if eq .Arch "aarch64" -}}
{{.Name}}-linux-arm64
{{- else -}}
{{.Name}}-linux-amd64
{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "k3s-io",
			Repo:        "k3s",
			Name:        "k3s",
			Description: "Lightweight Kubernetes",
			BinaryTemplate: `
{{- if eq .OS "darwin" -}}
{{.Name}}-darwin
{{ else if eq .Arch "aarch64" -}}
{{.Name}}-arm64
{{- else -}}
{{.Name}}
{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "cnrancher",
			Repo:        "autok3s",
			Name:        "autok3s",
			Description: "Run Rancher Lab's lightweight Kubernetes distribution k3s everywhere.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
	{{.Name}}_windows_amd64.exe
	{{- else if eq .OS "darwin" -}}
	{{.Name}}_darwin_amd64
	{{- else if eq .Arch "armv6l" -}}
	{{.Name}}_linux_arm
	{{- else if eq .Arch "armv7l" -}}
	{{.Name}}_linux_arm
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}_linux_arm64
	{{- else -}}
	{{.Name}}_linux_amd64
	{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "loft-sh",
			Repo:        "devspace",
			Name:        "devspace",
			Description: "Automate your deployment workflow with DevSpace and develop software directly inside Kubernetes.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
	{{.Name}}-windows-amd64.exe
	{{- else if and (eq .OS "darwin") (eq .Arch "aarch64") -}}
	{{.Name}}-darwin-arm64
	{{- else if and (eq .OS "darwin") -}}
	{{.Name}}-darwin-amd64
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}-linux-arm64
	{{- else -}}
	{{.Name}}-linux-amd64
	{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "tilt-dev",
			Repo:        "tilt",
			Name:        "tilt",
			Description: "A multi-service dev environment for teams on Kubernetes.",
			BinaryTemplate: `{{$version:=slice .Version 1}}
	{{ if HasPrefix .OS "ming" -}}
	{{.Name}}.{{$version}}.windows.x86_64.zip
	{{- else if and (eq .OS "darwin") (eq .Arch "aarch64") -}}
	{{.Name}}.{{$version}}.mac.arm64_ALPHA.tar.gz
	{{- else if eq .OS "darwin" -}}
	{{.Name}}.{{$version}}.mac.x86_64.tar.gz
	{{- else if eq .Arch "armv6l" -}}
	{{.Name}}.{{$version}}.linux.arm_ALPHA.tar.gz
	{{- else if eq .Arch "armv7l" -}}
	{{.Name}}.{{$version}}.linux.arm_ALPHA.tar.gz
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}.{{$version}}.linux.arm64_ALPHA.tar.gz
	{{- else -}}
	{{.Name}}.{{$version}}.linux.x86_64.tar.gz
	{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "rancher",
			Repo:        "k3d",
			Name:        "k3d",
			Description: "Helper to run Rancher Lab's k3s in Docker.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
	{{.Name}}-windows-amd64
	{{- else if eq .OS "darwin" -}}
	{{.Name}}-darwin-amd64
	{{- else if eq .Arch "armv6l" -}}
	{{.Name}}-linux-arm
	{{- else if eq .Arch "armv7l" -}}
	{{.Name}}-linux-arm
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}-linux-arm64
	{{- else -}}
	{{.Name}}-linux-amd64
	{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "alexellis",
			Repo:        "k3sup",
			Name:        "k3sup",
			Description: "Bootstrap Kubernetes with k3s over SSH < 1 min.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
	{{.Name}}.exe
	{{- else if eq .OS "darwin" -}}
	{{.Name}}-darwin
	{{- else if eq .Arch "armv6l" -}}
	{{.Name}}-armhf
	{{- else if eq .Arch "armv7l" -}}
	{{.Name}}-armhf
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}-arm64
	{{- else -}}
	{{.Name}}
	{{- end -}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "alexellis",
			Repo:        "arkade",
			Name:        "arkade",
			Description: "Portable marketplace for downloading your favourite devops CLIs and installing helm charts, with a single command.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
		{{.Name}}.exe
		{{- else if eq .OS "darwin" -}}
		{{.Name}}-darwin
		{{- else if eq .Arch "armv6l" -}}
		{{.Name}}-armhf
		{{- else if eq .Arch "armv7l" -}}
		{{.Name}}-armhf
		{{- else if eq .Arch "aarch64" -}}
		{{.Name}}-arm64
		{{- else -}}
		{{.Name}}
		{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "bitnami-labs",
			Repo:        "sealed-secrets",
			Name:        "kubeseal",
			Version:     "v0.14.1",
			Description: "A Kubernetes controller and tool for one-way encrypted Secrets",
			URLTemplate: `{{$arch := "arm"}}
{{- if eq .Arch "armv7l" -}}
https://github.com/bitnami-labs/sealed-secrets/releases/download/{{.Version}}/kubeseal-{{$arch}}
{{- end -}}

{{- if eq .Arch "arm64" -}}
https://github.com/bitnami-labs/sealed-secrets/releases/download/{{.Version}}/kubeseal-{{.Arch}}
{{- end -}}

{{- if HasPrefix .OS "ming" -}}
https://github.com/bitnami-labs/sealed-secrets/releases/download/{{.Version}}/kubeseal.exe
{{- end -}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- end -}}

{{- if and ( not ( or ( eq $arch "arm") ( eq $arch "arm64")) ) ( or ( eq .OS "darwin" ) ( eq .OS "linux" )) -}}
https://github.com/bitnami-labs/sealed-secrets/releases/download/{{.Version}}/kubeseal-{{.OS}}-{{$arch}}
{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "inlets",
			Repo:        "inletsctl",
			Name:        "inletsctl",
			Description: "Automates the task of creating an exit-server (tunnel server) on public cloud infrastructure.",
			URLTemplate: `
{{$fileName := ""}}
{{- if eq .Arch "armv6l" -}}
{{$fileName = "inletsctl-armhf.tgz"}}
{{- else if eq .Arch "armv7l" }}
{{$fileName = "inletsctl-armhf.tgz"}}
{{- else if eq .Arch "arm64" -}}
{{$fileName = "inletsctl-arm64.tgz"}}
{{ else if HasPrefix .OS "ming" -}}
{{$fileName = "inletsctl.exe.tgz"}}
{{- else if eq .OS "linux" -}}
{{$fileName = "inletsctl.tgz"}}
{{- else if eq .OS "darwin" -}}
{{$fileName = "inletsctl-darwin.tgz"}}
{{- end -}}
https://github.com/inlets/inletsctl/releases/download/{{.Version}}/{{$fileName}}`,
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}.exe
{{- else if eq .OS "darwin" -}}
{{.Name}}-darwin
{{- else if eq .OS "linux" -}}
{{.Name}}
{{- else if eq .Arch "armv6l" -}}
{{.Name}}-armhf
{{- else if eq .Arch "armv7l" -}}
{{.Name}}-armhf
{{- else if eq .Arch "aarch64" -}}
{{.Name}}-arm64
{{- end -}}`,
		},
		Tool{
			Name:        "osm",
			Repo:        "osm",
			Owner:       "openservicemesh",
			Description: "Open Service Mesh uniformly manages, secures, and gets out-of-the-box observability features.",
			BinaryTemplate: `
	{{$osStr := ""}}
	{{ if HasPrefix .OS "ming" -}}
	{{$osStr = "windows"}}
	{{- else if eq .OS "linux" -}}
	{{$osStr = "linux"}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "darwin"}}
	{{- end -}}
	{{.Version}}/osm-{{.Version}}-{{$osStr}}-amd64.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "linkerd",
			Repo:        "linkerd2",
			Name:        "linkerd2",
			Version:     "stable-2.9.1",
			Description: "Ultralight, security-first service mesh for Kubernetes.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}-cli-{{.Version}}-windows.exe
{{- else if eq .OS "darwin" -}}
{{.Name}}-cli-{{.Version}}-darwin
{{- else if eq .Arch "x86_64" -}}
{{.Name}}-cli-{{.Version}}-linux-amd64
{{- else if eq .Arch "aarch64" -}}
{{.Name}}-cli-{{.Version}}-linux-arm64
{{- end -}}
`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes-sigs",
			Repo:        "kubebuilder",
			Name:        "kubebuilder",
			NoExtension: true,
			Version:     "3.1.0",
			Description: "Framework for building Kubernetes APIs using custom resource definitions (CRDs).",
			URLTemplate: `{{$arch := "arm64"}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- end -}}

			{{$osStr := ""}}
			{{- if eq .OS "linux" -}}
			{{$osStr = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "darwin"}}
			{{- end -}}
			https://github.com/kubernetes-sigs/kubebuilder/releases/download/v{{.Version}}/kubebuilder_{{$osStr}}_{{$arch}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes-sigs",
			Repo:        "kustomize",
			Name:        "kustomize",
			Description: "Customization of kubernetes YAML configurations",
			BinaryTemplate: `
	{{$osStr := ""}}
	{{- if eq .OS "linux" -}}
	{{- if eq .Arch "x86_64" -}}
	{{$osStr = "linux_amd64"}}
	{{- else if eq .Arch "aarch64" -}}
  {{$osStr = "linux_arm64"}}
	{{- end -}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "darwin_amd64"}}
	{{- end -}}
	kustomize%2F{{.Version}}/{{.Name}}_{{.Version}}_{{$osStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "digitalocean",
			Repo:        "doctl",
			Name:        "doctl",
			Version:     "1.56.0",
			Description: "Official command line interface for the DigitalOcean API.",
			URLTemplate: `
		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "darwin"}}
		{{- end -}}

		{{$archStr := ""}}
		{{- if eq .Arch "x86_64" -}}
		{{$archStr = "amd64"}}
		{{- end -}}

		{{$archiveStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$archiveStr = "zip"}}
		{{- else -}}
		{{$archiveStr = "tar.gz"}}
		{{- end -}}

		https://github.com/digitalocean/doctl/releases/download/v{{.Version}}/doctl-{{.Version}}-{{$osStr}}-{{$archStr}}.{{$archiveStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "derailed",
			Repo:        "k9s",
			Name:        "k9s",
			Description: "Provides a terminal UI to interact with your Kubernetes clusters.",
			BinaryTemplate: `
		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "Windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "Linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "Darwin"}}
		{{- end -}}

		{{$archStr := .Arch}}
		{{- if eq .Arch "armv7l" -}}
		{{$archStr = "arm"}}
		{{- else if eq .Arch "aarch64" -}}
		{{$archStr = "arm64"}}
		{{- end -}}

		{{.Version}}/{{.Name}}_{{$osStr}}_{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "derailed",
			Repo:        "popeye",
			Name:        "popeye",
			Description: "Scans live Kubernetes cluster and reports potential issues with deployed resources and configurations.",
			BinaryTemplate: `
			{{$osStr := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$osStr = "Windows"}}
			{{- else if eq .OS "linux" -}}
			{{$osStr = "Linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "Darwin"}}
			{{- end -}}

			{{$archStr := .Arch}}
			{{- if eq .Arch "armv7l" -}}
			{{$archStr = "arm"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$archStr = "arm64"}}
			{{- end -}}

			{{.Version}}/{{.Name}}_{{$osStr}}_{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "civo",
			Repo:        "cli",
			Name:        "civo",
			Description: "CLI for interacting with your Civo resources.",
			// BinaryTemplate: `civo`,
			BinaryTemplate: `

		{{$extStr := "tar.gz"}}
		{{ if HasPrefix .OS "ming" -}}
		{{$extStr = "zip"}}
		{{- end -}}

		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "darwin"}}
		{{- end -}}

		{{$archStr := .Arch}}
		{{- if eq .Arch "armv7l" -}}
		{{$archStr = "arm"}}
		{{- else if eq .Arch "x86_64" -}}
		{{$archStr = "amd64"}}
		{{- end -}}

		{{.Version}}/{{.Name}}-{{.VersionNumber}}-{{$osStr}}-{{$archStr}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "hashicorp",
			Repo:        "terraform",
			Name:        "terraform",
			Version:     "1.0.0",
			Description: "Infrastructure as Code for major cloud providers.",
			URLTemplate: `{{$arch := .Arch}}

			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "arm"}}
			{{- end -}}

			{{$os := .OS}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{- end -}}

			https://releases.hashicorp.com/{{.Name}}/{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.zip`})

	tools = append(tools,
		Tool{
			Owner:       "hashicorp",
			Repo:        "vagrant",
			Name:        "vagrant",
			Version:     "2.2.14",
			Description: "Tool for building and distributing development environments.",
			URLTemplate: `{{$arch := .Arch}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- else if eq .Arch "armv7l" -}}
{{$arch = "arm"}}
{{- end -}}

{{$os := .OS}}
{{ if HasPrefix .OS "ming" -}}
{{$os = "windows"}}
{{- end -}}

https://releases.hashicorp.com/{{.Name}}/{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.zip`})

	tools = append(tools,
		Tool{
			Owner:       "hashicorp",
			Repo:        "packer",
			Name:        "packer",
			Version:     "1.6.5",
			Description: "Build identical machine images for multiple platforms from a single source configuration.",
			URLTemplate: `{{$arch := .Arch}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- else if eq .Arch "armv7l" -}}
{{$arch = "arm"}}
{{- end -}}

{{$os := .OS}}
{{ if HasPrefix .OS "ming" -}}
{{$os = "windows"}}
{{- end -}}

https://releases.hashicorp.com/{{.Name}}/{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.zip`})

	tools = append(tools,
		Tool{
			Owner:       "cli",
			Repo:        "cli",
			Name:        "gh",
			Description: "GitHub’s official command line tool.",
			BinaryTemplate: `

	{{$extStr := "tar.gz"}}
	{{ if HasPrefix .OS "ming" -}}
	{{$extStr = "zip"}}
	{{- end -}}

	{{$osStr := ""}}
	{{ if HasPrefix .OS "ming" -}}
	{{$osStr = "windows"}}
	{{- else if eq .OS "linux" -}}
	{{$osStr = "linux"}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "macOS"}}
	{{- end -}}

	{{$archStr := .Arch}}
	{{- if eq .Arch "aarch64" -}}
	{{$archStr = "arm64"}}
	{{- else if eq .Arch "x86_64" -}}
	{{$archStr = "amd64"}}
	{{- end -}}

	{{.Version}}/gh_{{.VersionNumber}}_{{$osStr}}_{{$archStr}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "buildpacks",
			Repo:        "pack",
			Name:        "pack",
			Description: "Build apps using Cloud Native Buildpacks.",
			BinaryTemplate: `

	{{$osStr := ""}}
	{{ if HasPrefix .OS "ming" -}}
	{{$osStr = "windows"}}
	{{- else if eq .OS "linux" -}}
	{{$osStr = "linux"}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "macos"}}
	{{- end -}}

	{{$extStr := "tgz"}}
	{{ if HasPrefix .OS "ming" -}}
	{{$extStr = "zip"}}
	{{- end -}}

	{{.Version}}/pack-{{.Version}}-{{$osStr}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "docker",
			Repo:        "buildx",
			Name:        "buildx",
			Description: "Docker CLI plugin for extended build capabilities with BuildKit.",
			BinaryTemplate: `
				{{$extStr := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$extStr = ".exe"}}
				{{- end -}}

				{{$osStr := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$osStr = "windows"}}
				{{- else if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$osStr = "darwin"}}
				{{- end -}}

				{{$archStr := .Arch}}
				{{- if eq .Arch "armv6l" -}}
				{{$archStr = "arm-v6"}}
				{{- else if eq .Arch "armv7l" -}}
				{{$archStr = "arm-v7"}}
				{{- else if eq .Arch "x86_64" -}}
				{{$archStr = "amd64"}}
				{{- end -}}

				{{.Version}}/{{.Name}}-{{.Version}}.{{$osStr}}-{{$archStr}}{{$extStr}}`,
		})

	// See issue: https://github.com/rakyll/hey/issues/229
	// 	tools = append(tools,
	// 		Tool{
	// 			Owner:   "rakyll",
	// 			Repo:    "hey",
	// 			Name:    "hey",
	// 			Version: "v0.1.2",
	// 			URLTemplate: `{{$arch := "amd64"}}
	// {{$ext := ""}}
	// {{$os := .OS}}

	// {{ if HasPrefix .OS "ming" -}}
	// {{$os = "windows"}}
	// {{$ext = ".exe"}}
	// {{$ext := ""}}
	// {{- end -}}

	// 	https://storage.googleapis.com/jblabs/dist/hey_{{$os}}_{{$.Version}}{{$ext}}`})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes",
			Repo:        "kops",
			Name:        "kops",
			Description: "Production Grade K8s Installation, Upgrades, and Management.",
			BinaryTemplate: `
	{{$osStr := ""}}
	{{- if eq .OS "linux" -}}
	{{- if eq .Arch "x86_64" -}}
	{{$osStr = "linux-amd64"}}
	{{- else if eq .Arch "aarch64" -}}
	{{$osStr = "linux-arm64"}}
	{{- end -}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "darwin-amd64"}}
	{{ else if HasPrefix .OS "ming" -}}
	{{$osStr ="windows-amd64"}}
	{{- end -}}

	{{.Version}}/{{.Name}}-{{$osStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes-sigs",
			Repo:        "krew",
			Name:        "krew",
			Description: "Package manager for kubectl plugins.",
			URLTemplate: `
			{{$osStr := ""}}
			{{- if eq .OS "linux" -}}
			{{- if eq .Arch "x86_64" -}}
			{{$osStr = "linux_amd64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$osStr = "linux_arm"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$osStr = "linux_arm64"}}
			{{- end -}}
			{{- else if eq .OS "darwin" -}}
			{{-  if eq .Arch "aarch64" -}}
			{{$osStr = "darwin_arm64"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$osStr = "darwin_amd64"}}
			{{- end -}}
			{{ else if HasPrefix .OS "ming" -}}
			{{-  if eq .Arch "aarch64" -}}
			{{$osStr = "darwin_arm64"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$osStr ="windows_amd64"}}
			{{- end -}}
			{{- end -}}
			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{$osStr}}.tar.gz`,
			BinaryTemplate: `
			{{$osStr := ""}}
			{{- if eq .OS "linux" -}}
			{{- if eq .Arch "x86_64" -}}
			{{$osStr = "linux_amd64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$osStr = "linux_arm"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$osStr = "linux_arm64"}}
			{{- end -}}
			{{- else if eq .OS "darwin" -}}
			{{-  if eq .Arch "aarch64" -}}
			{{$osStr = "darwin_arm64"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$osStr = "darwin_amd64"}}
			{{- end -}}
			{{ else if HasPrefix .OS "ming" -}}
			{{-  if eq .Arch "aarch64" -}}
			{{$osStr = "darwin_arm64"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$osStr ="windows_amd64"}}
			{{- end -}}
			{{- end -}}
			{{.Name}}-{{$osStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes",
			Repo:        "minikube",
			Name:        "minikube",
			Description: "Runs the latest stable release of Kubernetes, with support for standard Kubernetes features.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
	{{.Name}}-windows-amd64.exe
	{{- else if eq .OS "darwin" -}}
	{{.Name}}-darwin-amd64
	{{- else if eq .Arch "armv6l" -}}
	{{.Name}}-linux-arm
	{{- else if eq .Arch "armv7l" -}}
	{{.Name}}-linux-arm
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}-linux-arm64
	{{- else -}}
	{{.Name}}-linux-amd64
	{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "stern",
			Repo:        "stern",
			Name:        "stern",
			Description: "Multi pod and container log tailing for Kubernetes.",
			BinaryTemplate: `{{$arch := "arm"}}

{{- if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- else if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- end -}}

{{$os := .OS}}
{{$ext := "tar.gz"}}

{{ if HasPrefix .OS "ming" -}}
{{$os = "windows"}}
{{- end -}}

{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "boz",
			Repo:        "kail",
			Name:        "kail",
			Description: "Kubernetes log viewer.",
			BinaryTemplate: `{{$arch := "arm"}}

	{{- if eq .Arch "aarch64" -}}
	{{$arch = "arm64"}}
	{{- else if eq .Arch "x86_64" -}}
	{{$arch = "amd64"}}
	{{- end -}}

	{{$os := .OS}}
	{{$ext := "tar.gz"}}
	{{ if HasPrefix .OS "ming" -}}
	{{$os = "windows"}}
	{{- end -}}

	{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "mikefarah",
			Repo:        "yq",
			Name:        "yq",
			Description: "Portable command-line YAML processor.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
	{{.Name}}_windows_amd64.exe
	{{- else if eq .OS "darwin" -}}
	{{.Name}}_darwin_amd64
	{{- else if eq .Arch "armv6l" -}}
	{{.Name}}_linux_arm
	{{- else if eq .Arch "armv7l" -}}
	{{.Name}}_linux_arm
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}_linux_arm64
	{{- else -}}
	{{.Name}}_linux_amd64
	{{- end -}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "aquasecurity",
			Repo:        "kube-bench",
			Name:        "kube-bench",
			Version:     "0.4.0",
			Description: "Checks whether Kubernetes is deployed securely by running the checks documented in the CIS Kubernetes Benchmark.",
			BinaryTemplate: `
{{$arch := "arm"}}
{{$os := .OS}}

{{- if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- else if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- else if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- end -}}

{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "gohugoio",
			Repo:        "hugo",
			Name:        "hugo",
			Description: "Static HTML and CSS website generator.",
			BinaryTemplate: `
			{{$osStr := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$osStr = "Windows"}}
			{{- else if eq .OS "linux" -}}
			{{$osStr = "Linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "macOS"}}
			{{- end -}}

			{{$archStr := "64bit"}}
			{{- if eq .Arch "armv7l" -}}
			{{$archStr = "ARM"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$archStr = "ARM64"}}
			{{- end -}}

			{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$osStr}}-{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "docker",
			Repo:        "compose",
			Name:        "docker-compose",
			Version:     "1.29.1",
			Description: "Define and run multi-container applications with Docker.",
			BinaryTemplate: `
{{$osStr := ""}}
{{ if HasPrefix .OS "ming" -}}
{{$osStr = "Windows"}}
{{- else if eq .OS "linux" -}}
{{$osStr = "Linux"}}
{{- else if eq .OS "darwin" -}}
{{$osStr = "Darwin"}}
{{- end -}}
{{$ext := ""}}
{{ if HasPrefix .OS "ming" -}}
{{$ext = "exe"}}
{{- end -}}

{{.Version}}/{{.Name}}-{{$osStr}}-{{.Arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "open-policy-agent",
			Repo:        "opa",
			Name:        "opa",
			Description: "General-purpose policy engine that enables unified, context-aware policy enforcement across the entire stack.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
			{{.Name}}_windows_amd64.exe
			{{- else if eq .OS "darwin" -}}
			{{.Name}}_darwin_amd64
			{{- else -}}
			{{.Name}}_linux_amd64
			{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "minio",
			Repo:        "mc",
			Name:        "mc",
			Description: "MinIO Client is a replacement for ls, cp, mkdir, diff and rsync commands for filesystems and object storage.",
			URLTemplate: `{{$arch := .Arch}}
			{{ if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "armv6l" -}}
			{{$arch = "arm"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "arm"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- end -}}
			{{$osStr := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$osStr = "windows"}}
			{{- else if eq .OS "linux" -}}
			{{$osStr = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "darwin"}}
			{{- end -}}
			{{$ext := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$ext = ".exe"}}
			{{- end -}}
			https://dl.min.io/client/{{.Repo}}/release/{{$osStr}}-{{$arch}}/{{.Name}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "nats-io",
			Repo:        "natscli",
			Name:        "nats",
			Description: "Utility to interact with and manage NATS.",
			BinaryTemplate: `{{$arch := .Arch}}
			{{ if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "armv6l" -}}
			{{$arch = "arm6"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "arm7"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- end -}}

			{{$osStr := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$osStr = "windows"}}
			{{- else if eq .OS "linux" -}}
			{{$osStr = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "darwin"}}
			{{- end -}}

			{{.Version}}/{{.Name}}-{{.Version}}-{{$osStr}}-{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "argoproj",
			Repo:        "argo-cd",
			Name:        "argocd",
			Description: "Declarative, GitOps continuous delivery tool for Kubernetes.",
			BinaryTemplate: `
			{{$arch := ""}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- end -}}

			{{$osStr := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$osStr = "windows"}}
			{{- else if eq .OS "linux" -}}
			{{$osStr = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "darwin"}}
			{{- end -}}

			{{$ext := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$ext = ".exe"}}
			{{- end -}}

			{{.Version}}/argocd-{{$osStr}}-{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "containerd",
			Repo:        "nerdctl",
			Name:        "nerdctl",
			Description: "Docker-compatible CLI for containerd, with support for Compose",
			BinaryTemplate: `
{{ $file := "" }}
{{- if eq .OS "linux" -}}
	{{- if eq .Arch "armv6l" -}}
		{{ $file = "arm-v7.tar.gz" }}
	{{- else if eq .Arch "armv7l" -}}
		{{ $file = "arm-v7.tar.gz" }}
	{{- else if eq .Arch "aarch64" -}}
		{{ $file = "arm64.tar.gz" }}
	{{- else -}}
		{{ $file = "amd64.tar.gz" }}
	{{- end -}}
{{- end -}}

{{.Version}}/{{.Name}}-{{.VersionNumber}}-{{.OS}}-{{$file}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "istio",
			Repo:        "istio",
			Name:        "istioctl",
			Description: "Service Mesh to establish a programmable, application-aware network using the Envoy service proxy.",
			BinaryTemplate: `
				{{$arch := .Arch}}
				{{ if eq .Arch "x86_64" -}}
				{{$arch = "amd64"}}
				{{- else if eq .Arch "arm" -}}
				{{$arch = "armv7"}}
				{{- else if eq .Arch "armv6l" -}}
				{{$arch = "armv7"}}
				{{- else if eq .Arch "armv7l" -}}
				{{$arch = "armv7"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{$versionString:=(printf "%s-%s" .OS $arch)}}
				{{ if HasPrefix .OS "ming" -}}
				{{$versionString = "win"}}
				{{- else if eq .OS "darwin" -}}
				{{$versionString = "osx"}}
				{{- end -}}

				{{$ext := ".tar.gz"}}
				{{ if HasPrefix .OS "ming" -}}
				{{$ext = ".zip"}}
				{{- end -}}

				{{.Version}}/{{.Name}}-{{.VersionNumber}}-{{$versionString}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "tektoncd",
			Repo:        "cli",
			Name:        "tkn",
			Description: "A CLI for interacting with Tekton.",
			BinaryTemplate: `
				{{$arch := .Arch}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "x86_64"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{$osString:= .OS}}
				{{ if HasPrefix .OS "ming" -}}
				{{$osString = "Windows"}}
				{{- else if eq .OS "darwin" -}}
				{{$osString = "Darwin"}}
				{{- else if eq .OS "linux" -}}
				{{$osString = "Linux"}}
				{{- end -}}

				{{$ext := ".tar.gz"}}
				{{ if HasPrefix .OS "ming" -}}
				{{$ext = ".zip"}}
				{{- end -}}

				{{.Version}}/tkn_{{.VersionNumber}}_{{$osString}}_{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "inlets",
			Repo:        "inlets-pro",
			Name:        "inlets-pro",
			Description: "Cloud Native Tunnel for HTTP and TCP traffic.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}.exe
{{- else if eq .OS "darwin" -}}
{{.Name}}-darwin
{{- else if eq .Arch "armv6l" -}}
{{.Name}}-armhf
{{- else if eq .Arch "armv7l" -}}
{{.Name}}-armhf
{{- else if eq .Arch "aarch64" -}}
{{.Name}}-arm64
{{- else -}}
{{.Name}}
{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "rancher",
			Repo:        "kim",
			Name:        "kim",
			Version:     "v0.1.0-beta.4",
			Description: "Build container images inside of Kubernetes. (Experimental)",
			BinaryTemplate: `
			{{ $ext := "" }}
			{{ $osStr := "linux" }}
			{{ if HasPrefix .OS "ming" -}}
			{{	$osStr = "windows" }}
			{{ $ext = ".exe" }}
			{{- else if eq .OS "darwin" -}}
			{{  $osStr = "darwin" }}
			{{- end -}}

			{{ $archStr := "amd64" }}

			{{- if eq .Arch "armv6l" -}}
			{{ $archStr = "arm" }}
			{{- else if eq .Arch "armv7l" -}}
			{{ $archStr = "arm" }}
			{{- else if eq .Arch "aarch64" -}}
			{{ $archStr = "arm64" }}
			{{- end -}}

			{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "aquasecurity",
			Repo:        "trivy",
			Name:        "trivy",
			Description: "Vulnerability Scanner for Containers and other Artifacts, Suitable for CI.",
			BinaryTemplate: `
				{{$arch := .Arch}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "64bit"}}
				{{- else if eq .Arch "armv7l" -}}
				{{$arch = "ARM"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$arch = "ARM64"}}
				{{- end -}}

				{{$osString:= .OS}}
				{{ if HasPrefix .OS "darwin" -}}
				{{$osString = "macOS"}}
				{{- else if eq .OS "linux" -}}
				{{$osString = "Linux"}}
				{{- end -}}

				{{$ext := ".tar.gz"}}

				{{.Version}}/trivy_{{.VersionNumber}}_{{$osString}}-{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "fluxcd",
			Repo:        "flux2",
			Name:        "flux",
			Description: "Continuous Delivery solution for Kubernetes powered by GitOps Toolkit.",
			BinaryTemplate: `
				{{$arch := .Arch}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "amd64"}}
				{{- else if eq .Arch "armv7l" -}}
				{{$arch = "arm"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{$osString := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$osString = "windows"}}
				{{- else if eq .OS "linux" -}}
				{{$osString = "linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$osString = "darwin"}}
				{{- end -}}

				{{$ext := ".tar.gz"}}
				{{ if HasPrefix .OS "ming" -}}
				{{$ext = ".zip"}}
				{{- end -}}

				{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$osString}}_{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "FairwindsOps",
			Repo:        "polaris",
			Name:        "polaris",
			Description: "Run checks to ensure Kubernetes pods and controllers are configured using best practices.",
			BinaryTemplate: `
				{{$arch := "amd64"}}
				{{if eq .Arch "armv7l" -}}
				{{$arch = "armv7"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{$osString:= .OS}}
				{{ if HasPrefix .OS "darwin" -}}
				{{$osString = "darwin"}}
				{{- else if eq .OS "linux" -}}
				{{$osString = "linux"}}
				{{- end -}}
				{{$ext := ".tar.gz"}}
				{{.Name}}_{{$osString}}_{{$arch}}{{$ext}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "influxdata",
			Repo:        "influxdb",
			Name:        "influx",
			Version:     "2.0.8",
			Description: "InfluxDB’s command line interface (influx) is an interactive shell for the HTTP API.",
			URLTemplate: `{{$arch := .Arch}}
		{{ if eq .Arch "x86_64" -}}
		{{$arch = "amd64"}}
		{{- else if eq .Arch "aarch64" -}}
		{{$arch = "arm64"}}
		{{- end -}}

		{{$ext := ".tar.gz"}}
		{{ if HasPrefix .OS "ming" -}}
		{{$ext = ".zip"}}
		{{- end -}}

				https://dl.{{.Owner}}.com/{{.Repo}}/releases/{{.Repo}}2-client-{{.Version}}-{{.OS}}-{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "argoproj-labs",
			Repo:        "argocd-autopilot",
			Name:        "argocd-autopilot",
			Description: "An opinionated way of installing Argo-CD and managing GitOps repositories.",
			URLTemplate: `
{{$arch := ""}}
{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- else if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- end -}}

{{$osString := ""}}
{{ if HasPrefix .OS "ming" -}}
{{$osString = "windows"}}
{{- else if eq .OS "linux" -}}
{{$osString = "linux"}}
{{- else if eq .OS "darwin" -}}
{{$osString = "darwin"}}
{{- end -}}
{{$ext := ".tar.gz"}}
https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{$osString}}-{{$arch}}{{$ext}}
`,
			BinaryTemplate: `
{{$arch := ""}}
{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- else if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- end -}}

{{$osString := ""}}
{{ if HasPrefix .OS "ming" -}}
{{$osString = "windows"}}
{{- else if eq .OS "linux" -}}
{{$osString = "linux"}}
{{- else if eq .OS "darwin" -}}
{{$osString = "darwin"}}
{{- end -}}
{{.Name}}-{{$osString}}-{{$arch}}
`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "FairwindsOps",
			Repo:        "nova",
			Name:        "nova",
			Description: "Find outdated or deprecated Helm charts running in your cluster.",
			BinaryTemplate: `
				{{$arch := "amd64"}}
				{{if eq .Arch "armv7l" -}}
				{{$arch = "armv7"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{$osString:= .OS}}
				{{ if HasPrefix .OS "darwin" -}}
				{{$osString = "darwin"}}
				{{- else if eq .OS "linux" -}}
				{{$osString = "linux"}}
				{{- end -}}
				{{$ext := ".tar.gz"}}

				{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$osString}}_{{$arch}}{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "johanhaleby",
			Repo:        "kubetail",
			Name:        "kubetail",
			Version:     "1.6.13",
			Description: "Bash script to tail Kubernetes logs from multiple pods at the same time.",
			URLTemplate: `https://raw.githubusercontent.com/{{.Owner}}/{{.Repo}}/{{.Version}}/{{.Name}}`,
		})

	tools = append(tools,
		Tool{

			Owner:       "squat",
			Repo:        "kilo",
			Name:        "kgctl",
			Description: "A CLI to manage Kilo, a multi-cloud network overlay built on WireGuard and designed for Kubernetes.",
			BinaryTemplate: `
{{$os := .OS}}
{{ if HasPrefix .OS "ming" -}}
{{$os = "windows"}}
{{- end -}}
{{$arch := "amd64"}}
{{ if eq .Arch "armv7l" -}}
{{$arch = "arm"}}
{{- else if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- end -}}
{{.Name}}-{{$os}}-{{$arch}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "getporter",
			Repo:        "porter",
			Name:        "porter",
			Version:     "v0.38.4",
			Description: "With Porter you can package your application artifact, tools, etc. as a bundle that can distribute and install.",
			BinaryTemplate: `
			{{ $ext := "" }}
			{{ $osStr := "linux" }}
			{{ if HasPrefix .OS "ming" -}}
			{{	$osStr = "windows" }}
			{{ $ext = ".exe" }}
			{{- else if eq .OS "darwin" -}}
			{{  $osStr = "darwin" }}
			{{- end -}}

			{{ $archStr := "amd64" }}
			{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "k0sproject",
			Repo:        "k0s",
			Name:        "k0s",
			Description: "Zero Friction Kubernetes",
			BinaryTemplate: `
			{{$arch := ""}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- end -}}
			{{.Name}}-{{.Version}}-{{$arch}}
			`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "k0sproject",
			Repo:        "k0sctl",
			Name:        "k0sctl",
			Description: "A bootstrapping and management tool for k0s clusters",
			BinaryTemplate: `{{$arch := "x64"}}
	{{- if eq .Arch "aarch64" -}}
	{{$arch = "arm64"}}
	{{- end -}}

	{{$os := .OS}}
	{{$ext := ""}}

	{{ if HasPrefix .OS "ming" -}}
	{{$os = "win"}}
	{{$ext = ".exe"}}
	{{- end -}}

	{{.Name}}-{{$os}}-{{$arch}}{{$ext}}`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "equinix",
			Repo:        "metal-cli",
			Name:        "metal",
			Version:     "0.6.0-alpha2",
			Description: "Official Equinix Metal CLI",
			BinaryTemplate: `{{ $ext := "" }}
				{{ $osStr := "linux" }}
				{{ if HasPrefix .OS "ming" -}}
				{{	$osStr = "windows" }}
				{{ $ext = ".exe" }}
				{{- else if eq .OS "darwin" -}}
				{{  $osStr = "darwin" }}
				{{- end -}}

				{{ $archStr := "amd64" }}
				{{- if eq .Arch "armv6l" -}}
				{{ $archStr = "armv6" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $archStr = "armv7" }}
				{{- else if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- end -}}
				{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "kanisterio",
			Repo:        "kanister",
			Name:        "kanctl",
			Description: "Framework for application-level data management on Kubernetes.",
			URLTemplate: `
{{ $osStr := "linux" }}
{{- if eq .OS "darwin" -}}
{{ $osStr = "darwin" }}
{{- end -}}

https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Repo}}_{{$.Version}}_{{$osStr}}_amd64.tar.gz`,
			BinaryTemplate: `{{.Name}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kastenhq",
			Repo:        "kubestr",
			Name:        "kubestr",
			Description: "Kubestr discovers, validates and evaluates your Kubernetes storage options.",

			URLTemplate: `
{{ $ext := "tar.gz" }}
{{ $osStr := "linux" }}

{{- if eq .OS "darwin" -}}
{{ $osStr = "darwin" }}
{{- else if HasPrefix .OS "ming" -}}
{{ $osStr = "windows" }}
{{ $ext = ".zip" }}
{{- end -}}

https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.Version}}-{{$osStr}}-amd64.{{$ext}}`,
			BinaryTemplate: `{{.Name}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kastenhq",
			Repo:        "external-tools",
			Name:        "k10multicluster",
			Description: "Multi-cluster support for K10.",

			BinaryTemplate: `
	{{ $osStr := "linux" }}
	{{ $archStr := "amd64" }}

	{{- if eq .Arch "aarch64" -}}
	{{ $archStr = "arm64" }}
	{{- end -}}

	{{- if eq .OS "darwin" -}}
	{{ $osStr = "macOS" }}
	{{- end -}}

	{{.Name}}_{{.Version}}_{{$osStr}}_{{$archStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kastenhq",
			Repo:        "external-tools",
			Name:        "k10tools",
			Description: "Tools for evaluating and debugging K10.",

			BinaryTemplate: `
	{{ $osStr := "linux" }}
	{{ $archStr := "amd64" }}

	{{- if eq .Arch "aarch64" -}}
	{{ $archStr = "arm64" }}
	{{- end -}}

	{{- if eq .OS "darwin" -}}
	{{ $osStr = "macOS" }}
	{{- end -}}

	{{.Name}}_{{.Version}}_{{$osStr}}_{{$archStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "sigstore",
			Repo:        "cosign",
			Name:        "cosign",
			Description: "Container Signing, Verification and Storage in an OCI registry.",
			BinaryTemplate: `{{ $ext := "" }}
				{{ $osStr := "linux" }}
				{{ if HasPrefix .OS "ming" -}}
				{{ $osStr = "windows" }}
				{{ $ext = ".exe" }}
				{{- else if eq .OS "darwin" -}}
				{{  $osStr = "darwin" }}
				{{- end -}}

				{{ $archStr := "amd64" }}

				{{- if eq .Arch "armv6l" -}}
				{{ $archStr = "arm" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $archStr = "arm" }}
				{{- else if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- end -}}

				{{.Version}}/{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		},
	)

	tools = append(tools,
		Tool{
			Owner: "sigstore",
			Repo:  "rekor",
			Name:  "rekor-cli",

			Description: "Secure Supply Chain - Transparency Log",
			BinaryTemplate: `{{ $ext := "" }}
			{{ $osStr := "linux" }}
			{{ if HasPrefix .OS "ming" -}}
			{{ $osStr = "windows" }}
			{{ $ext = ".exe" }}
			{{- else if eq .OS "darwin" -}}
			{{  $osStr = "darwin" }}
			{{- end -}}

			{{ $archStr := "amd64" }}

			{{- if eq .Arch "armv6l" -}}
			{{ $archStr = "arm" }}
			{{- else if eq .Arch "armv7l" -}}
			{{ $archStr = "arm" }}
			{{- else if eq .Arch "aarch64" -}}
			{{ $archStr = "arm64" }}
			{{- end -}}

			{{.Version}}/{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "aquasecurity",
			Repo:        "tfsec",
			Name:        "tfsec",
			Version:     "v0.57.1",
			Description: "Security scanner for your Terraform code",
			BinaryTemplate: `{{ $ext := "" }}
				{{ $osStr := "linux" }}
				{{ if HasPrefix .OS "ming" -}}
				{{	$osStr = "windows" }}
				{{ $ext = ".exe" }}
				{{- else if eq .OS "darwin" -}}
				{{  $osStr = "darwin" }}
				{{- end -}}

				{{ $archStr := "amd64" }}
				{{- if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- end -}}
				{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "wagoodman",
			Repo:        "dive",
			Name:        "dive",
			Version:     "0.10.0",
			Description: "A tool for exploring each layer in a docker image",
			URLTemplate: `{{$osStr := ""}}
			{{- if HasPrefix .OS "ming" -}}
			{{$osStr = "windows"}}
			{{- else if eq .OS "linux" -}}
			{{$osStr = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "darwin"}}
			{{- end -}}

			{{$archiveStr := ""}}
			{{- if HasPrefix .OS "ming" -}}
			{{$archiveStr = ".zip"}}
			{{- else -}}
			{{$archiveStr = ".tar.gz"}}
			{{- end -}}

			https://github.com/{{.Owner}}/{{.Name}}/releases/download/v{{.Version}}/{{.Name}}_{{.Version}}_{{$osStr}}_amd64{{$archiveStr}}`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "goreleaser",
			Repo:        "goreleaser",
			Name:        "goreleaser",
			Description: "Deliver Go binaries as fast and easily as possible",
			BinaryTemplate: `
		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "Windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "Linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "Darwin"}}
		{{- end -}}

		{{$archStr := ""}}
		{{- if eq .Arch "x86_64" -}}
		{{$archStr = "x86_64"}}
		{{- else if eq .Arch "aarch64" -}}
        {{$archStr = "arm64"}}
		{{- end -}}

		{{$archiveStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$archiveStr = "zip"}}
		{{- else -}}
		{{$archiveStr = "tar.gz"}}
		{{- end -}}

		{{.Name}}_{{$osStr}}_{{$archStr}}.{{$archiveStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "armosec",
			Repo:        "kubescape",
			Name:        "kubescape",
			Description: "kubescape is the first tool for testing if Kubernetes is deployed securely as defined in Kubernetes Hardening Guidance by to NSA and CISA",
			BinaryTemplate: `
		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "ubuntu"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "macos"}}
		{{- end -}}


		{{.Name}}-{{$osStr}}-latest`,
		})

	return tools
}
