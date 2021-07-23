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
			Version:     "1.6",
			Description: "jq is a lightweight and flexible command-line JSON processor",
			URLTemplate: `{{$arch := "arm"}}

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

https://github.com/stedolan/jq/releases/download/jq-{{.Version}}/jq-{{$os}}{{$arch}}{{$ext}}`,
		})

	// https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/darwin/amd64/kubectl
	tools = append(tools,
		Tool{
			Owner:       "kubernetes",
			Repo:        "kubernetes",
			Name:        "kubectl",
			Version:     "v1.20.0",
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
			Owner:       "ahmetb",
			Repo:        "kubectx",
			Name:        "kubectx",
			Version:     "v0.9.1",
			Description: "Faster way to switch between clusters.",
			URLTemplate: `https://github.com/ahmetb/kubectx/releases/download/{{.Version}}/kubectx`,
			NoExtension: true,
		})

	tools = append(tools,
		Tool{
			Owner:       "ahmetb",
			Repo:        "kubectx",
			Name:        "kubens",
			Version:     "v0.9.1",
			Description: "Switch between Kubernetes namespaces smoothly.",
			URLTemplate: `https://github.com/ahmetb/kubectx/releases/download/{{.Version}}/kubens`,
			NoExtension: true,
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
			Version:     "0.8.2",
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
			Version:     "v0.7.0",
			Description: "Open Service Mesh uniformly manages, secures, and gets out-of-the-box observability features.",
			URLTemplate: `
	{{$osStr := ""}}
	{{ if HasPrefix .OS "ming" -}}
	{{$osStr = "windows"}}
	{{- else if eq .OS "linux" -}}
	{{$osStr = "linux"}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "darwin"}}
	{{- end -}}
	https://github.com/openservicemesh/osm/releases/download/{{.Version}}/osm-{{.Version}}-{{$osStr}}-amd64.tar.gz`,
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
			Version:     "v4.0.0",
			Description: "Customization of kubernetes YAML configurations",
			URLTemplate: `
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
	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/kustomize%2F{{.Version}}/{{.Name}}_{{.Version}}_{{$osStr}}.tar.gz`,
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
			Version:     "v0.24.10",
			Description: "Provides a terminal UI to interact with your Kubernetes clusters.",
			URLTemplate: `
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
https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{.Version}}_{{$osStr}}_{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "derailed",
			Repo:        "popeye",
			Name:        "popeye",
			Version:     "v0.9.0",
			Description: "Scans live Kubernetes cluster and reports potential issues with deployed resources and configurations.",
			URLTemplate: `
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
	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{$osStr}}_{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:          "civo",
			Repo:           "cli",
			Name:           "civo",
			Version:        "0.7.11",
			Description:    "CLI for interacting with your Civo resources.",
			BinaryTemplate: `civo`,
			URLTemplate: `

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
https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}-{{.Version}}-{{$osStr}}-{{$archStr}}.{{$extStr}}`,
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
			Owner:          "cli",
			Repo:           "cli",
			Name:           "gh",
			Version:        "1.13.1",
			Description:    "GitHub’s official command line tool.",
			BinaryTemplate: `gh`,
			URLTemplate: `

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
https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/gh_{{.Version}}_{{$osStr}}_{{$archStr}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "buildpacks",
			Repo:        "pack",
			Name:        "pack",
			Description: "Build apps using Cloud Native Buildpacks.",
			Version:     "0.14.2",
			URLTemplate: `

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

https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/pack-v{{.Version}}-{{$osStr}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:          "docker",
			Repo:           "buildx",
			Name:           "buildx",
			Version:        "0.4.2",
			Description:    "Docker CLI plugin for extended build capabilities with BuildKit.",
			BinaryTemplate: `buildx`,
			URLTemplate: `

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
		https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}-v{{.Version}}.{{$osStr}}-{{$archStr}}{{$extStr}}`,
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
			Version:     "v1.18.2",
			Description: "Production Grade K8s Installation, Upgrades, and Management.",
			URLTemplate: `
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

	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{$osStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes-sigs",
			Repo:        "krew",
			Name:        "krew",
			Version:     "v0.4.0",
			Description: "Package manager for kubectl plugins.",
			URLTemplate: `
		https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}.tar.gz`,
			BinaryTemplate: `
		{{$osStr := ""}}
		{{- if eq .OS "linux" -}}
		{{- if eq .Arch "x86_64" -}}
		{{$osStr = "linux_amd64"}}
		{{- else if eq .Arch "armv7l" -}}
		{{$osStr = "linux_arm"}}
		{{- end -}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "darwin_amd64"}}
		{{ else if HasPrefix .OS "ming" -}}
		{{$osStr ="windows_amd64.exe"}}
		{{- end -}}
		{{.Name}}-{{$osStr}}
	`,
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
			Version:     "1.13.0",
			Description: "Multi pod and container log tailing for Kubernetes.",
			URLTemplate: `{{$arch := "arm"}}

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
https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "boz",
			Repo:        "kail",
			Name:        "kail",
			Version:     "0.15.0",
			Description: "Kubernetes log viewer.",
			URLTemplate: `{{$arch := "arm"}}

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
	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.tar.gz`,
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
			URLTemplate: `
{{$arch := "arm"}}
{{$os := .OS}}

{{- if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- else if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- else if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- end -}}

https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "gohugoio",
			Repo:        "hugo",
			Name:        "hugo",
			Version:     "0.79.0",
			Description: "Static HTML and CSS website generator.",
			URLTemplate: `
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
	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}_{{.Version}}_{{$osStr}}-{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "docker",
			Repo:        "compose",
			Name:        "docker-compose",
			Version:     "1.29.1",
			Description: "Define and run multi-container applications with Docker.",
			URLTemplate: `
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
https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{$osStr}}-{{.Arch}}{{$ext}}`,
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
			Version:     "0.0.23",
			Description: "Utility to interact with and manage NATS.",
			URLTemplate: `{{$arch := .Arch}}
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
			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.Version}}-{{$osStr}}-{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "argoproj",
			Repo:        "argo-cd",
			Name:        "argocd",
			Version:     "v1.8.6",
			Description: "Declarative, GitOps continuous delivery tool for Kubernetes.",
			URLTemplate: `
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

			https://github.com/argoproj/argo-cd/releases/download/{{.Version}}/argocd-{{$osStr}}-{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "containerd",
			Repo:        "nerdctl",
			Name:        "nerdctl",
			Version:     "v0.7.2",
			Description: "Docker-compatible CLI for containerd, with support for Compose",
			URLTemplate: `
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

https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.VersionNumber}}-{{.OS}}-{{$file}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "istio",
			Repo:        "istio",
			Name:        "istioctl",
			Version:     "1.10.0",
			Description: "Service Mesh to establish a programmable, application-aware network using the Envoy service proxy.",
			URLTemplate: `
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

				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.VersionNumber}}-{{$versionString}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "tektoncd",
			Repo:        "cli",
			Name:        "tkn",
			Version:     "0.17.2",
			Description: "A CLI for interacting with Tekton.",
			URLTemplate: `
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

				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/tkn_{{.Version}}_{{$osString}}_{{$arch}}{{$ext}}`,
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
			Version:     `v0.1.0-alpha.12`,
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
			Version:     "0.17.2",
			Description: "Vulnerability Scanner for Containers and other Artifacts, Suitable for CI.",
			URLTemplate: `
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

				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/trivy_{{.Version}}_{{$osString}}-{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "fluxcd",
			Repo:        "flux2",
			Name:        "flux",
			Version:     "0.13.4",
			Description: "Continuous Delivery solution for Kubernetes powered by GitOps Toolkit.",
			URLTemplate: `
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

				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}_{{.Version}}_{{$osString}}_{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "FairwindsOps",
			Repo:        "polaris",
			Name:        "polaris",
			Version:     "3.2.1",
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
				{{.Name}}_{{.Version}}_{{$osString}}_{{$arch}}{{$ext}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "influxdata",
			Repo:        "influxdb",
			Name:        "influx",
			Version:     "2.0.7",
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
			Version:     "0.2.1",
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
			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/v{{.Version}}/{{.Name}}-{{$osString}}-{{$arch}}{{$ext}}
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
			Version:     "2.3.2",
			Description: "Find outdated or deprecated Helm charts running in your cluster.",
			URLTemplate: `
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
				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{.Version}}_{{$osString}}_{{$arch}}{{$ext}}
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
			Version:     "0.3.0",
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
			Owner:       "packethost",
			Repo:        "packet-cli",
			Name:        "packet",
			Version:     "0.5.0",
			Description: "The Equinix Metal CLI allows interaction with Equinix Metal platform.",
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
			{{ if eq .Arch "armv7l" -}}
			{{$archStr = "armv7"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$archStr = "arm64"}}
			{{- end -}}
			{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
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

	return tools
}
