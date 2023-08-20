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
  {{- if eq .Arch "arm64" -}}
{{.Name}}-darwin-arm64
  {{- else -}}
{{.Name}}-darwin
  {{- end -}}
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
			Owner:           "helm",
			Repo:            "helm",
			Name:            "helm",
			VersionStrategy: "github",
			Description:     "The Kubernetes Package Manager: Think of it like apt/yum/homebrew for Kubernetes.",
			URLTemplate: `
						{{$os := .OS}}
						{{$arch := .Arch}}
						{{$ext := "tar.gz"}}

						{{- if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
							{{$arch = "arm64"}}
						{{- else if eq .Arch "x86_64" -}}
							{{ $arch = "amd64" }}
						{{- else if eq .Arch "armv7l" -}}
							{{ $arch = "arm" }}
						{{- end -}}

						{{ if HasPrefix .OS "ming" -}}
						{{$os = "windows"}}
						{{$ext = "zip"}}
						{{- end -}}

						https://get.helm.sh/helm-{{.Version}}-{{$os}}-{{$arch}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "helmfile",
			Repo:        "helmfile",
			Name:        "helmfile",
			Description: "Deploy Kubernetes Helm Charts",
			BinaryTemplate: `{{$arch := ""}}
						{{- if eq .Arch "x86_64" -}}
						{{$arch = "amd64"}}
						{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
						{{$arch = "arm64"}}
						{{- end -}}

						{{$os := .OS}}
						{{ if HasPrefix .OS "ming" -}}
						{{$os = "windows"}}
						{{- end -}}

					helmfile_{{.VersionNumber}}_{{$os}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "jqlang",
			Repo:        "jq",
			Name:        "jq",
			Description: "jq is a lightweight and flexible command-line JSON processor",
			BinaryTemplate: `{{$arch := "arm"}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "64"}}
{{- else if eq .Arch "arm64" -}}
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
			Version:     "v1.24.2",
			Description: "Run commands against Kubernetes clusters",
			URLTemplate: `{{$arch := "arm"}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- end -}}

{{- if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- end -}}

{{- if eq .Arch "arm64" -}}
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
			Description:    "Faster way to switch between clusters.",
			BinaryTemplate: `kubectx`,
			NoExtension:    true,
		})

	tools = append(tools,
		Tool{
			Owner:          "ahmetb",
			Repo:           "kubectx",
			Name:           "kubens",
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
	{{- if eq .Arch "aarch64" -}}
	{{.Name}}-darwin-arm64
	{{- else if eq .Arch "arm64" -}}
	{{.Name}}-darwin-arm64
	{{- else -}}
	{{.Name}}-darwin-amd64
	{{- end -}}
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
			Owner:       "etcd-io",
			Repo:        "etcd",
			Name:        "etcd",
			Description: "Distributed reliable key-value store for the most critical data of a distributed system.",
			BinaryTemplate: `
					{{$ext := "zip"}}
					{{- if eq .OS "linux" -}}
						{{$ext = "tar.gz"}}
					{{- end -}}

					{{$arch := .Arch}}
					{{ if (eq .Arch "x86_64") -}}
						{{$arch = "amd64"}}
					{{- else if eq .Arch "aarch64" -}}
						{{$arch = "arm64"}}
					{{- end -}}

					{{$osString:= .OS}}
					{{ if HasPrefix .OS "ming" -}}
						{{$osString = "windows"}}
					{{- end -}}

					{{.Name}}-{{.Version}}-{{$osString}}-{{$arch}}.{{$ext}}
				`,
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
			Owner:       "devspace-sh",
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
			Owner:       "k3d-io",
			Repo:        "k3d",
			Name:        "k3d",
			Description: "Helper to run Rancher Lab's k3s in Docker.",
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
				{{ if eq .Arch "arm64" -}}
				{{.Name}}-darwin-arm64
				{{- else -}}
				{{.Name}}-darwin
				{{- end -}}
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
			Repo:        "run-job",
			Name:        "run-job",
			Description: "Run a Kubernetes Job and get the logs when it's done.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
			{{.Name}}.exe
			{{- else if eq .OS "darwin" -}}
				{{ if eq .Arch "arm64" -}}
				{{.Name}}-darwin-arm64
				{{- else -}}
				{{.Name}}-darwin
				{{- end -}}
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
			Owner:       "inlets",
			Repo:        "mixctl",
			Name:        "mixctl",
			Description: "A tiny TCP load-balancer.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
			{{.Name}}.exe
			{{- else if eq .OS "darwin" -}}
				{{ if eq .Arch "arm64" -}}
				{{.Name}}-darwin-arm64
				{{- else -}}
				{{.Name}}-darwin
				{{- end -}}
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

	// The kubeseal repo has releases for both the Helm chart and the CLI
	// tool, and there's no way to filter, so the version has to be hard-coded.
	tools = append(tools,
		Tool{
			Owner:       "bitnami-labs",
			Repo:        "sealed-secrets",
			Name:        "kubeseal",
			Version:     "v0.19.5",
			Description: "A Kubernetes controller and tool for one-way encrypted Secrets",
			BinaryTemplate: `{{$arch := ""}}
		{{- if eq .Arch "aarch64" -}}
		{{$arch = "arm64"}}
                {{- else if eq .Arch "arm64" -}}
                {{$arch = "arm64"}}
		{{- else if eq .Arch "x86_64" -}}
		{{$arch = "amd64"}}
		{{- else if eq .Arch "armv7l" -}}
		{{$arch = "arm"}}
		{{- end -}}

		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "darwin"}}
		{{- end -}}

		{{.Version}}/{{.Name}}-{{.VersionNumber}}-{{$osStr}}-{{$arch}}.tar.gz`,
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
{{- else if eq .Arch "aarch64" -}}
{{$fileName = "inletsctl-arm64.tgz"}}
{{ else if HasPrefix .OS "ming" -}}
{{$fileName = "inletsctl.exe.tgz"}}
{{- else if eq .OS "linux" -}}
{{$fileName = "inletsctl.tgz"}}
{{- else if eq .OS "darwin" -}}
	{{- if eq .Arch "arm64" -}}
	{{$fileName = "inletsctl-darwin-arm64.tgz"}}
	{{- else }}
	{{$fileName = "inletsctl-darwin.tgz"}}
	{{- end -}}
{{- end -}}
https://github.com/inlets/inletsctl/releases/download/{{.Version}}/{{$fileName}}`,
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}
{{- else if eq .OS "darwin" -}}
	{{- if eq .Arch "arm64" -}}
	{{.Name}}-darwin-arm64
	{{- else if eq .Arch "x86_64" -}}
	{{.Name}}-darwin
	{{- end -}}
{{- else if eq .OS "linux" -}}
	{{- if eq .Arch "armv6l" -}}
	{{.Name}}-armhf
	{{- else if eq .Arch "armv7l" -}}
	{{.Name}}-armhf
	{{- else if eq .Arch "aarch64" -}}
	{{.Name}}-arm64
	{{- else -}}
	{{.Name}}
	{{- end -}}
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
			Version:     "v5.0.3",
			BinaryTemplate: `
	{{$osStr := ""}}
	{{$ext := "tar.gz"}}
	{{- if eq .OS "linux" -}}
	{{- if eq .Arch "x86_64" -}}
	{{$osStr = "linux_amd64"}}
	{{- else if eq .Arch "aarch64" -}}
  {{$osStr = "linux_arm64"}}
	{{- end -}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "darwin_amd64"}}
	{{- end -}}
	{{ if HasPrefix .OS "ming" -}}
	{{$osStr = "windows_amd64"}}
	{{- end -}}
	kustomize%2F{{.Version}}/{{.Name}}_{{.Version}}_{{$osStr}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "google",
			Repo:        "go-containerregistry",
			Name:        "crane",
			Description: "crane is a tool for interacting with remote images and registries",
			BinaryTemplate: `{{$arch := ""}}
			{{- if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
			{{- else if eq .Arch "arm64" -}}
				{{$arch = "arm64"}}
			{{- else if eq .Arch "x86_64" -}}
				{{$arch = "x86_64"}}
			{{- else if eq .Arch "armv7l" -}}
				{{$arch = "armv6"}}
			{{- end -}}
	
			{{$osStr := ""}}
			{{ if HasPrefix .OS "ming" -}}
				{{$osStr = "Windows"}}
			{{- else if eq .OS "linux" -}}
				{{$osStr = "Linux"}}
			{{- else if eq .OS "darwin" -}}
				{{$osStr = "Darwin"}}
			{{- end -}}
	
			{{.Version}}/go-containerregistry_{{$osStr}}_{{$arch}}.tar.gz`,
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
			Owner:       "eksctl-io",
			Repo:        "eksctl",
			Name:        "eksctl",
			Description: "Amazon EKS Kubernetes cluster management",
			BinaryTemplate: `
			{{$arch := ""}}
			{{$extStr := "tar.gz"}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "arm64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "armv7"}}
			{{- else if eq .Arch "armv6l" -}}
			{{$arch = "armv6"}}
			{{- end -}}

			{{$os := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "Windows"}}
			{{$extStr = "zip"}}
			{{- else if eq .OS "linux" -}}
			{{$os = "Linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "Darwin"}}
			{{- end -}}
			{{.Name}}_{{$os}}_{{$arch}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "aws",
			Repo:        "eks-anywhere",
			Name:        "eksctl-anywhere",
			Description: "Run Amazon EKS on your own infrastructure",
			BinaryTemplate: `
			{{$os := .OS}}
			{{$ext := "tar.gz"}}
			{{$arch := .Arch}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- end -}}

			{{.Name}}-{{.Version}}-{{$os}}-{{$arch}}.{{$ext}}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "derailed",
			Repo:        "k9s",
			Name:        "k9s",
			Description: "Provides a terminal UI to interact with your Kubernetes clusters.",
			BinaryTemplate: `
		{{$os := "" }}
		{{ if HasPrefix .OS "ming" -}}
		{{$os = "Windows"}}
		{{- else if eq .OS "linux" -}}
		{{$os = "Linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$os = "Darwin"}}
		{{- end -}}

		{{$arch := .Arch}}
		{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
		{{$arch = "arm64"}}
		{{- else if eq .Arch "x86_64" -}}
		{{ $arch = "amd64" }}
		{{- else if eq .Arch "armv7l" -}}
		{{$arch = "arm"}}
		{{- end -}}

		{{.Name}}_{{$os}}_{{$arch}}.tar.gz`,
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
			Version:     "1.3.9",
			Description: "Infrastructure as Code for major cloud providers.",
			URLTemplate: `
			{{$arch := ""}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "arm64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "arm"}}
			{{- end -}}

			{{$os := .OS}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{- end -}}

			https://releases.hashicorp.com/{{.Name}}/{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "gruntwork-io",
			Repo:        "terragrunt",
			Name:        "terragrunt",
			Description: "Terragrunt is a thin wrapper for Terraform that provides extra tools for working with multiple Terraform modules",
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
			{{- if eq .Arch "x86_64" -}}
			{{$archStr = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$archStr = "arm64"}}
			{{- end -}}

			{{.Version}}/{{.Name}}_{{$osStr}}_{{$archStr}}{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "hashicorp",
			Repo:        "vagrant",
			Name:        "vagrant",
			Version:     "2.2.19",
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
			Version:     "1.8.0",
			Description: "Build identical machine images for multiple platforms from a single source configuration.",
			URLTemplate: `
			{{$arch := ""}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
            {{$arch = "arm64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "arm"}}
			{{- end -}}

			{{$os := .OS}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{- end -}}

			https://releases.hashicorp.com/{{.Name}}/{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "hashicorp",
			Repo:        "waypoint",
			Name:        "waypoint",
			Version:     "0.8.1",
			Description: "Easy application deployment for Kubernetes and Amazon ECS",
			URLTemplate: `
			{{$arch := ""}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
            {{$arch = "arm64"}}
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
			    {{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
				{{$archStr = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{$archStr = "amd64"}}
				{{- end -}}

				{{.Version}}/{{.Name}}-{{.Version}}.{{$osStr}}-{{$archStr}}{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "alexellis",
			Repo:        "hey",
			Name:        "hey",
			Description: "Load testing tool",
			BinaryTemplate: `
			{{$osStr := ""}}
			{{- if eq .OS "linux" -}}
				{{- if eq .Arch "x86_64" -}}
					{{$osStr = ""}}
				{{- else if eq .Arch "aarch64" -}}
					{{$osStr = "-linux-arm64"}}
				{{- else if eq .Arch "armv7l" -}}
					{{$osStr = "-linux-armv7"}}
				{{- end -}}
			{{- else if eq .OS "darwin" -}}
				{{- if eq .Arch "x86_64" -}}
					{{$osStr = "-darwin-amd64"}}
				{{- else if eq .Arch "arm64" -}}
					{{$osStr = "-darwin-arm64"}}
				{{- end -}}
			{{ else if HasPrefix .OS "ming" -}}
				{{$osStr =".exe"}}
			{{- end -}}
			
			{{.Name}}{{$osStr}}`,
		})

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
			{{-  if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
			{{$osStr = "darwin_arm64"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$osStr = "darwin_amd64"}}
			{{- end -}}
			{{ else if HasPrefix .OS "ming" -}}
			{{- if eq .Arch "x86_64" -}}
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
			{{-  if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
			{{$osStr = "darwin_arm64"}}
			{{- else if eq .Arch "x86_64" -}}
			{{$osStr = "darwin_amd64"}}
			{{- end -}}
			{{ else if HasPrefix .OS "ming" -}}
			{{-  if eq .Arch "x86_64" -}}
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
			BinaryTemplate: `{{$arch := .Arch}}
			{{ if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "armv6l" -}}
			{{$arch = "armv6"}}
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

			{{.Version}}/{{.Name}}-{{$osStr}}-{{$arch}}{{$ext}}`,
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

	{{.Version}}/{{.Name}}_v{{.VersionNumber}}_{{$os}}_{{$arch}}.tar.gz`,
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
			Description: "Define and run multi-container applications with Docker.",
			BinaryTemplate: `
{{$arch := .Arch}}

{{$osStr := ""}}
{{ if HasPrefix .OS "ming" -}}
{{$osStr = "windows"}}
{{- else if eq .OS "linux" -}}
{{$osStr = "linux"}}
  {{- if eq .Arch "armv7l" -}}
  {{ $arch = "armv7"}}
  {{- end }}
{{- else if eq .OS "darwin" -}}
{{$osStr = "darwin"}}
{{- end -}}
{{$ext := ""}}
{{ if HasPrefix .OS "ming" -}}
{{$ext = ".exe"}}
{{- end -}}

{{.Version}}/{{.Name}}-{{$osStr}}-{{$arch}}{{$ext}}`,
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

			{{.Version}}/{{.Name}}-{{.VersionNumber}}-{{$osStr}}-{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "argoproj",
			Repo:        "argo-cd",
			Name:        "argocd",
			Description: "Declarative, GitOps continuous delivery tool for Kubernetes.",
			BinaryTemplate: `
			{{$arch := .Arch}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
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
			BinaryTemplate: `
			{{$arch := ""}}
			{{$os := ""}}
			{{$ext := ""}}

			{{- if eq .Arch "aarch64" -}}
            {{$arch = "-arm64"}}
			{{- else if eq .Arch "arm64" -}}
			{{$arch = "-arm64"}}
			{{- else if (or (eq .Arch "armv6l") (eq .Arch "armv7l")) -}}
			{{$arch = "-armhf"}}
			{{- end -}}

			{{ if eq .OS "darwin" -}}
			{{$os = "-darwin"}}
			{{ else if HasPrefix .OS "ming" -}}
			{{$ext = ".exe"}}
			{{- end -}}
			{{.Name}}{{$os}}{{$arch}}{{$ext}}`,
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
			Owner:       "stefanprodan",
			Repo:        "timoni",
			Name:        "timoni",
			Description: "A package manager for Kubernetes powered by CUE.",
			BinaryTemplate: `
				{{$arch := .Arch}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
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
	{{ $osStr := "Linux" }}
	{{ $arch := .Arch }}
	
	{{- if eq .Arch "x86_64" -}}
	{{$arch = "amd64"}}
	{{- end -}}

	{{- if eq .Arch "aarch64" -}}
	{{$arch = "arm64"}}
	{{- end -}}

	{{- if eq .OS "darwin" -}}
	{{ $osStr = "MacOS" }}
	{{- else if HasPrefix .OS "ming" -}}
	{{ $osStr = "Windows" }}
	{{- end -}}

	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$osStr}}_{{$arch}}.{{$ext}}`,
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

	{{.Name}}_{{.Version}}_{{$osStr}}_{{$archStr}}.tar.gz`,
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

	{{- if eq .Arch "arm64" -}}
	{{ $archStr = "arm64" }}
	{{- end -}}

	{{- if eq .OS "darwin" -}}
	{{ $osStr = "macOS" }}
	{{- end -}}

	{{.Name}}_{{.Version}}_{{$osStr}}_{{$archStr}}.tar.gz`,
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
			Owner:       "kubescape",
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
	tools = append(tools,
		Tool{
			Owner:       "operator-framework",
			Repo:        "operator-sdk",
			Name:        "operator-sdk",
			Description: "Operator SDK is a tool for scaffolding and generating code for building Kubernetes operators",
			BinaryTemplate: `{{$arch := "amd64"}}

{{if eq .Arch "aarch64" -}}
{{$arch = "arm64"}}
{{- else if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
{{- end -}}

{{$os := .OS}}

{{- if eq .OS "darwin" -}}
{{$os = "darwin"}}
{{- end -}}

{{.Version}}/{{.Name}}_{{$os}}_{{$arch}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes-sigs",
			Repo:        "cluster-api",
			Name:        "clusterctl",
			Description: "The clusterctl CLI tool handles the lifecycle of a Cluster API management cluster",
			BinaryTemplate: `{{ $ext := "" }}
			{{ $osStr := "linux" }}
			{{ if HasPrefix .OS "ming" -}}
			{{ $osStr = "windows" }}
			{{ $ext = ".exe" }}
			{{- else if eq .OS "darwin" -}}
			{{ $osStr = "darwin" }}
			{{- end -}}

			{{ $archStr := "amd64" }}
			{{- if eq .Arch "aarch64" -}}
			{{ $archStr = "arm64" }}
			{{- end -}}
			{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "loft-sh",
			Repo:        "vcluster",
			Name:        "vcluster",
			Description: "Create fully functional virtual Kubernetes clusters - Each vcluster runs inside a namespace of the underlying k8s cluster.",
			BinaryTemplate: `{{ $ext := "" }}
			{{ $osStr := "linux" }}
			{{ if HasPrefix .OS "ming" -}}
            {{ $osStr = "windows" }}
			{{ $ext = ".exe" }}
			{{- else if eq .OS "darwin" -}}
			{{  $osStr = "darwin" }}
			{{- end -}}

			{{ $archStr := "amd64" }}
            {{- if eq .Arch "x86_64" -}}
            {{ $archStr = "amd64" }}
			{{- else if eq .Arch "aarch64" -}}
			{{ $archStr = "arm64" }}
			{{- end -}}
			{{.Name}}-{{$osStr}}-{{$archStr}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "guumaster",
			Repo:        "hostctl",
			Name:        "hostctl",
			Description: "Dev tool to manage /etc/hosts like a pro!",
			BinaryTemplate: `
			{{ $osStr := "" }}
			{{ $archStr := "" }}
			{{ $extStr := ".tar.gz" }}

			{{ if HasPrefix .OS "ming" -}}
			{{ $osStr = "windows" }}
			{{ $extStr = ".zip" }}
			{{- else if eq .OS "linux" -}}
			{{ $osStr = "linux" }}
			{{- else if eq .OS "darwin" -}}
			{{ $osStr = "macOS" }}
			{{- end -}}

			{{- if eq .Arch "x86_64" -}}
			{{ $archStr = "64-bit" }}
			{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
			{{ $archStr = "arm64" }}
			{{- end -}}

			{{ .Name }}_{{ .VersionNumber }}_{{ $osStr }}_{{ $archStr }}{{ $extStr }}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "sunny0826",
			Repo:        "kubecm",
			Name:        "kubecm",
			Description: "Easier management of kubeconfig. ",
			BinaryTemplate: `
			{{ $osStr := "" }}
			{{ $archStr := "" }}
			{{ $extStr := ".tar.gz" }}

			{{ if HasPrefix .OS "ming" -}}
			{{ $osStr = "Windows" }}
			{{- else if eq .OS "linux" -}}
			{{ $osStr = "Linux" }}
			{{- else if eq .OS "darwin" -}}
			{{ $osStr = "Darwin" }}
			{{- end -}}

			{{- if eq .Arch "x86_64" -}}
			{{ $archStr = "x86_64" }}
			{{- else if eq .Arch "aarch64" -}}
			{{ $archStr = "arm64" }}
			{{- end -}}

			{{.Version}}/{{.Name}}_{{.Version}}_{{$osStr}}_{{$archStr}}{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "FiloSottile",
			Repo:        "mkcert",
			Name:        "mkcert",
			Description: "A simple zero-config tool to make locally trusted development certificates with any names you'd like.",
			BinaryTemplate: `
				{{ $osStr := "" }}
				{{ $archStr := "" }}
				{{$archiveStr := ""}}
				{{- if HasPrefix .OS "ming" -}}
				{{$archiveStr = ".exe"}}
				{{- else -}}
				{{$archiveStr = ""}}
				{{- end -}}
	
				{{ if HasPrefix .OS "ming" -}}
				{{ $osStr = "windows" }}
				{{- else if eq .OS "linux" -}}
				{{ $osStr = "linux" }}
				{{- else if eq .OS "darwin" -}}
				{{ $osStr = "darwin" }}
				{{- end -}}
	
				{{- if eq .Arch "x86_64" -}}
				{{ $archStr = "amd64" }}
				{{- else if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $archStr = "arm" }}
				{{- end -}}
	
				{{.Version}}/{{.Name}}-v{{.VersionNumber}}-{{$osStr}}-{{$archStr}}{{$archiveStr}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "getsops",
			Repo:        "sops",
			Name:        "sops",
			Description: "Simple and flexible tool for managing secrets",
			BinaryTemplate: `
			{{ $archStr := "" }}

			{{- if eq .Arch "x86_64" -}}
			{{ $archStr = "amd64" }}
			{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
			{{ $archStr = "arm64" }}
			{{- end -}}

			{{- if HasPrefix .OS "ming" -}}
			{{ .Name }}-{{ .Version }}.exe
			{{- else -}}
			{{ .Name }}-{{ .Version }}.{{ .OS }}.{{ $archStr }}
			{{- end -}}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "dagger",
			Repo:        "dagger",
			Name:        "dagger",
			Description: "A portable devkit for CI/CD pipelines.",
			URLTemplate: `
	{{ $ext := ".tar.gz"}}
	{{- if HasPrefix .OS "ming" -}}
	{{ $ext = ".zip"}}
	{{- end -}}

	{{ $os := .OS }}
	{{- if HasPrefix .OS "ming" -}}
	{{ $os = "windows" }}
	{{- end -}}

	{{ $arch := .Arch }}
	{{- if eq .Arch "x86_64" -}}
	{{ $arch = "amd64" }}
	{{- else if eq .Arch "aarch64" -}}
	{{ $arch = "arm64" }}
	{{- end -}}

	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/dagger_{{.Version}}_{{$os}}_{{$arch}}{{$ext}}`,
			BinaryTemplate: `
			{{ $name := "dagger" }}

			{{ $os := .OS }}
			{{- if HasPrefix .OS "ming" -}}
			{{ $name = "dagger.exe" }}
			{{- end -}}

			{{$name}}
	`})

	tools = append(tools,
		Tool{
			Owner:       "kumahq",
			Repo:        "kuma",
			Name:        "kumactl",
			Version:     "1.4.1",
			Description: "kumactl is a CLI to interact with Kuma and its data",
			URLTemplate: `
			{{$osStr := ""}}
			{{$archStr := ""}}
			{{- if HasPrefix .OS "linux" -}}
			{{$osStr = "ubuntu"}}
			{{- else if eq .OS "darwin" -}}
			{{$osStr = "darwin"}}
			{{- end -}}

			{{- if eq .Arch "x86_64" -}}
			{{$archStr = "amd64"}}
			{{- else -}}
			{{$archStr = .Arch}}
			{{- end -}}
			https://download.konghq.com/mesh-alpine/{{.Repo}}-{{.Version}}-{{$osStr}}-{{$archStr}}.tar.gz`,
			BinaryTemplate: `{{.Name}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "jandedobbeleer",
			Repo:        "oh-my-posh",
			Name:        "oh-my-posh",
			Description: "A prompt theme engine for any shell that can display kubernetes information.",
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

			posh-{{$osStr}}-{{$archStr}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "caddyserver",
			Repo:        "caddy",
			Name:        "caddy",
			Description: "Caddy is an extensible server platform that uses TLS by default",
			URLTemplate: `
			{{ $os := "linux" }}
			{{ $arch := "amd64" }}
			{{ $ext := "tar.gz" }}

			{{- if eq .Arch "aarch64" -}}
			{{ $arch = "arm64" }}
			{{- else if eq .Arch "arm64" -}}
			{{ $arch = "arm64" }}
			{{- else if eq .Arch "armv7l" -}}
			{{ $arch = "armv7" }}
			{{- else if eq .Arch "armv6l" -}}
			{{ $arch = "armv6" }}
			{{- end -}}

			{{ if HasPrefix .OS "ming" -}}
			{{ $os = "windows" }}
			{{ $ext = "zip" }}
			{{- else if eq .OS "darwin" -}}
			{{  $os = "mac" }}
			{{- end -}}

			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{ .Version }}/{{ .Name }}_{{ .VersionNumber }}_{{ $os }}_{{ $arch }}.{{ $ext }}
			`,
			BinaryTemplate: `
			{{ if HasPrefix .OS "ming" -}}
			{{ .Name }}.exe
			{{- else -}}
			{{ .Name }}
			{{- end -}}
			`,
		},
	)

	tools = append(tools,
		Tool{
			Owner:       "nats-io",
			Repo:        "nats-server",
			Name:        "nats-server",
			Description: "Cloud native message bus and queue server",
			BinaryTemplate: `
				{{ $archStr := "" }}
				{{ $osStr := "linux" }}

				{{ if HasPrefix .OS "ming" -}}
				{{ $osStr = "windows" }}
				{{- else if eq .OS "darwin" -}}
				{{ $osStr = "darwin" }}
				{{- end -}}

				{{- if eq .Arch "x86_64" -}}
				{{ $archStr = "amd64" }}
				{{- else if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- else if eq .Arch "arm64" -}}
				{{ $archStr = "arm64" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $archStr = "arm7" }}
				{{- end -}}
	
				{{ .Name }}-{{ .Version }}-{{ $osStr }}-{{ $archStr }}.zip
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "cilium",
			Repo:        "cilium-cli",
			Name:        "cilium",
			Description: "CLI to install, manage & troubleshoot Kubernetes clusters running Cilium.",
			URLTemplate: `
			{{$arch := ""}}
			{{$extStr := "tar.gz"}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "arm64" -}}
			{{$arch = "arm64"}}
			{{- end -}}

			{{$os := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{- else if eq .OS "linux" -}}
			{{$os = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "darwin"}}
			{{- end -}}
			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{$os}}-{{$arch}}.{{$extStr}}`,
			BinaryTemplate: `
			{{ if HasPrefix .OS "ming" -}}
			{{ .Name }}.exe
			{{- else -}}
			{{ .Name }}
			{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "junegunn",
			Repo:        "fzf",
			Name:        "fzf",
			Description: "General-purpose command-line fuzzy finder",
			BinaryTemplate: `
				{{ $osStr := "linux" }}
				{{ $ext := ".tar.gz" }}
				{{ if HasPrefix .OS "ming" -}}
				{{ $osStr = "windows" }}
				{{ $ext = ".zip" }}
				{{- else if eq .OS "darwin" -}}
				{{  $osStr = "darwin" }}
				{{ $ext = ".zip" }}
				{{- end -}}
				{{ $archStr := "amd64" }}
				{{- if eq .Arch "armv6l" -}}
				{{ $archStr = "armv6" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $archStr = "armv7" }}
				{{- else if eq .Arch "arm64" -}}
				{{ $archStr = "arm64" }}
				{{- else if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- end -}}
				{{.Name}}-{{.VersionNumber}}-{{$osStr}}_{{$archStr}}{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "cilium",
			Repo:        "hubble",
			Name:        "hubble",
			Description: "CLI for network, service & security observability for Kubernetes clusters running Cilium.",
			URLTemplate: `
			{{$arch := ""}}
			{{$extStr := "tar.gz"}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "arm64" -}}
			{{$arch = "arm64"}}
			{{- end -}}

			{{$os := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{- else if eq .OS "linux" -}}
			{{$os = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "darwin"}}
			{{- end -}}
			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{$os}}-{{$arch}}.{{$extStr}}`,
			BinaryTemplate: `
			{{ if HasPrefix .OS "ming" -}}
			{{ .Name }}.exe
			{{- else -}}
			{{ .Name }}
			{{- end -}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "hairyhenderson",
			Repo:        "gomplate",
			Name:        "gomplate",
			Description: "A flexible commandline tool for template rendering. Supports lots of local and remote datasources.",
			URLTemplate: `
				{{ $os := "linux" }}
				{{ $arch := "amd64" }}
				{{ $ext := "" }}
	
				{{- if eq .Arch "aarch64" -}}
				{{ $arch = "arm64" }}
				{{- else if eq .Arch "arm64" -}}
				{{ $arch = "arm64" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $arch = "armv7" }}
				{{- end -}}
	
				{{ if HasPrefix .OS "ming" -}}
				{{ $os = "windows" }}
				{{ $ext = ".exe" }}
				{{- else if eq .OS "darwin" -}}
				{{  $os = "darwin" }}
				{{- end -}}
	
				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{ .Version }}/{{ .Name }}_{{ $os }}-{{ $arch }}{{ $ext }}
				`,
			BinaryTemplate: `
				{{ if HasPrefix .OS "ming" -}}
				{{ .Name }}.exe
				{{- else -}}
				{{ .Name }}
				{{- end -}}
				`,
		})

	tools = append(tools,
		Tool{
			Name:        "just",
			Owner:       "casey",
			Repo:        "just",
			Description: "Just a command runner",
			URLTemplate: `
				{{ $os := "unknown-linux" }}
				{{ $arch := "x86_64" }}
				{{ $ext := "-musl.tar.gz" }}

				{{- if eq .Arch "aarch64" -}}
				{{ $arch = "aarch64" }}
				{{- else if eq .Arch "arm64" -}}
				{{ $arch = "aarch64" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $arch = "armv7" }}
				{{ $ext = "-musleabihf.tar.gz" }}
				{{- end -}}

				{{ if HasPrefix .OS "ming" -}}
				{{ $os = "pc-windows" }}
				{{ $ext = "-msvc.zip" }}
				{{- else if eq .OS "darwin" -}}
				{{  $os = "apple-darwin" }}
				{{ $ext = ".tar.gz" }}
				{{- end -}}
			https://github.com/{{ .Owner }}/{{ .Repo }}/releases/download/{{ .VersionNumber }}/{{ .Name }}-{{ .VersionNumber }}-{{ $arch }}-{{ $os }}{{ $ext }}
			`,
			BinaryTemplate: `
			{{- if HasPrefix .OS "ming" -}}
			{{ .Name }}.exe
			{{- else -}}
			{{ .Name }}
			{{- end -}}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "prometheus",
			Repo:        "prometheus",
			Name:        "promtool",
			Description: "Prometheus rule tester and debugging utility",
			URLTemplate: `
			{{$arch := ""}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "arm64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{ $arch = "armv7" }}
			{{- end -}}

			{{$os := ""}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{- else if eq .OS "linux" -}}
			{{$os = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "darwin"}}
			{{- end -}}
			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Repo}}-{{.VersionNumber}}.{{$os}}-{{$arch}}.tar.gz`,
			BinaryTemplate: `
			{{ if HasPrefix .OS "ming" -}}
			{{ .Name }}.exe
			{{- else -}}
			{{ .Name }}
			{{- end -}}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "siderolabs",
			Repo:        "talos",
			Name:        "talosctl",
			Description: "The command-line tool for managing Talos Linux OS.",
			URLTemplate: `
					{{ $os := "linux" }}
					{{ $arch := "amd64" }}
					{{ $ext := "" }}
					{{- if eq .Arch "aarch64" -}}
					{{ $arch = "arm64" }}
					{{- else if eq .Arch "arm64" -}}
					{{ $arch = "arm64" }}
					{{- else if eq .Arch "armv7l" -}}
					{{ $arch = "armv7" }}
					{{- end -}}
					{{ if HasPrefix .OS "ming" -}}
					{{ $os = "windows" }}
					{{ $ext = ".exe" }}
					{{- else if eq .OS "darwin" -}}
					{{  $os = "darwin" }}
					{{- end -}}
					https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{ .Version }}/{{ .Name }}-{{ $os }}-{{ $arch }}{{ $ext }}
						`,
			BinaryTemplate: `
					{{ if HasPrefix .OS "ming" -}}
					{{ .Name }}.exe
					{{- else -}}
					{{ .Name }}
					{{- end -}}
					`,
		})

	tools = append(tools,
		Tool{
			Owner:       "tenable",
			Repo:        "terrascan",
			Name:        "terrascan",
			Description: "Detect compliance and security violations across Infrastructure as Code.",
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
						{{- if eq .Arch "aarch64" -}}
						{{$archStr = "arm64"}}
						{{- else if eq .Arch "x86_64" -}}
						{{$archStr = "x86_64"}}
						{{- end -}}
						{{.Name}}_{{slice .Version 1}}_{{$osStr}}_{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "golangci",
			Repo:        "golangci-lint",
			Name:        "golangci-lint",
			Description: "Go linters aggregator.",
			BinaryTemplate: `
							{{$os := ""}}
							{{ if HasPrefix .OS "ming" -}}
							{{$os = "windows"}}
							{{- else if eq .OS "linux" -}}
							{{$os = "linux"}}
							{{- else if eq .OS "darwin" -}}
							{{$os = "darwin"}}
							{{- end -}}
							{{$arch := .Arch}}
							{{- if eq .Arch "x86_64" -}}
							{{$arch = "amd64"}}
							{{- else if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
							{{$arch = "arm64"}}
							{{- else if eq .Arch "armv7l" -}}
							{{$arch = "armv7"}}
							{{- end -}}
							{{.Name}}-{{.VersionNumber}}-{{$os}}-{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "oven-sh",
			Repo:        "bun",
			Name:        "bun",
			Description: "Bun is an incredibly fast JavaScript runtime, bundler, transpiler and package manager – all in one.",
			BinaryTemplate: `
							{{$arch := .Arch}}
							{{- if eq .Arch "x86_64" -}}
							{{$arch = "x64"}}
							{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
							{{$arch = "aarch64"}}
							{{- end -}}
							{{.Name}}-{{ .OS }}-{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "jesseduffield",
			Repo:        "lazygit",
			Name:        "lazygit",
			Description: "A simple terminal UI for git commands.",
			BinaryTemplate: `
								{{$os := ""}}
								{{$ext := "tar.gz" }}
								{{ if HasPrefix .OS "ming" -}}
								{{$os = "Windows"}}
								{{$ext = "zip" }}
								{{- else if eq .OS "linux" -}}
								{{$os = "Linux"}}
								{{- else if eq .OS "darwin" -}}
								{{$os = "Darwin"}}
								{{- end -}}

								{{$arch := .Arch}}
								{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
								{{$arch = "x86_64"}}
								{{- else if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
								{{$arch = "arm64"}}
								{{- end -}}
								{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "redpanda-data",
			Repo:        "redpanda",
			Name:        "rpk",
			Description: "Kafka compatible streaming platform for mission critical workloads.",
			BinaryTemplate: `
			{{$os := ""}}
			{{$arch := ""}}
			{{- if eq .OS "linux" -}}
			{{$os = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "darwin"}}
			{{- end -}}

			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
			{{$arch = "arm64"}}
			{{- end -}}

			{{.Name}}-{{$os}}-{{$arch}}.zip
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "hashicorp",
			Repo:        "vault",
			Name:        "vault",
			Version:     "1.11.2",
			Description: "A tool for secrets management, encryption as a service, and privileged access management.",
			URLTemplate: `
			{{$arch := ""}}
			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "arm64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "arm"}}
			{{- end -}}

			{{$os := .OS}}
			{{ if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{- end -}}

			https://releases.hashicorp.com/{{.Name}}/{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "helm",
			Repo:        "chart-releaser",
			Name:        "cr",
			Description: "Hosting Helm Charts via GitHub Pages and Releases",
			URLTemplate: `
			{{$os := ""}}
			{{$arch := ""}}
			{{$ext := ".tar.gz"}}
			{{- if eq .OS "linux" -}}
			{{$os = "linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "darwin"}}
			{{- else if HasPrefix .OS "ming" -}}
			{{$os = "windows"}}
			{{$ext = ".zip"}}
			{{- end -}}

			{{- if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
			{{$arch = "arm64"}}
			{{- else if eq .Arch "armv7l" -}}
			{{$arch = "armv7"}}
			{{- end -}}
			https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Repo}}_{{.VersionNumber}}_{{$os}}_{{$arch}}{{$ext}}
			`,
			BinaryTemplate: `
			{{ if HasPrefix .OS "ming" -}}
			{{ .Name }}.exe
			{{- else -}}
			{{ .Name }}
			{{- end -}}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "hadolint",
			Repo:        "hadolint",
			Name:        "hadolint",
			Description: "A smarter Dockerfile linter that helps you build best practice Docker images",
			BinaryTemplate: `
			{{$os := ""}}
			{{$arch := .Arch}}
			{{$ext := ""}}
			{{- if eq .OS "linux" -}}
			{{$os = "Linux"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "Darwin"}}
			{{- else if HasPrefix .OS "ming" -}}
			{{$os = "Windows"}}
			{{$ext = ".exe"}}
			{{- end -}}

			{{- if eq .Arch "aarch64" -}}
			{{$arch = "arm64"}}
			{{- end -}}
			{{.Name}}-{{$os}}-{{$arch}}{{$ext}}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "coreos",
			Repo:        "butane",
			Name:        "butane",
			Description: "Translates human readable Butane Configs into machine readable Ignition Configs",
			BinaryTemplate: `
			{{$os := ""}}
			{{$ext := ""}}
			{{$arch := .Arch}}
			{{- if eq .OS "linux" -}}
			{{$os = "unknown-linux-gnu"}}
			{{- else if eq .OS "darwin" -}}
			{{$os = "apple-darwin"}}
			{{- else if HasPrefix .OS "ming" -}}
			{{$os = "pc-windows-gnu"}}
			{{$ext = ".exe"}}
			{{- end -}}

			{{- if eq .Arch "arm64" -}}
			{{$arch = "aarch64"}}
			{{- end -}}
			{{.Name}}-{{$arch}}-{{$os}}{{$ext}}
			`,
		})

	tools = append(tools,
		Tool{
			Owner:       "superfly",
			Repo:        "flyctl",
			Name:        "flyctl",
			Description: "Command line tools for fly.io services",
			URLTemplate: `
				{{$os := ""}}
				{{$arch := .Arch}}
				{{$ext := ".tar.gz"}}

				{{- if eq .OS "linux" -}}
				{{$os = "Linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$os = "macOS"}}
				{{- else if HasPrefix .OS "ming" -}}
				{{$os = "Windows"}}
				{{$ext = ".zip"}}
				{{- end -}}
	
				{{- if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- end -}}

				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Repo}}_{{.VersionNumber}}_{{$os}}_{{$arch}}{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "yannh",
			Repo:        "kubeconform",
			Name:        "kubeconform",
			Description: "A FAST Kubernetes manifests validator, with support for Custom Resources",
			BinaryTemplate: `
				{{$os := .OS}}
				{{$arch := .Arch}}
				{{$ext := ".tar.gz"}}
				{{- if HasPrefix .OS "ming" -}}
				{{$os = "windows"}}
				{{$ext = ".zip"}}
				{{- end -}}
	
				{{$arch := .Arch}}
				{{- if eq .Arch "x86_64" -}}
				{{$arch = "amd64"}}
				{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{.Name}}-{{$os}}-{{$arch}}{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "open-policy-agent",
			Repo:        "conftest",
			Name:        "conftest",
			Description: "Write tests against structured configuration data using the Open Policy Agent Rego query language",
			BinaryTemplate: `
				{{$os := .OS}}
				{{$arch := .Arch}}
				{{$ext := "tar.gz"}}

				{{$arch := .Arch}}
				{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{ if HasPrefix .OS "ming" -}}
				{{$os = "Windows"}}
				{{$ext = "zip"}}
				{{- else if eq .OS "linux" -}}
				{{$os = "Linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$os = "Darwin"}}
				{{- end -}}

				{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "instrumenta",
			Repo:        "kubeval",
			Name:        "kubeval",
			Description: "Validate your Kubernetes configuration files, supports multiple Kubernetes versions",
			BinaryTemplate: `
				{{$os := .OS}}
				{{$arch := .Arch}}
				{{$ext := "tar.gz"}}

				{{$arch := .Arch}}
				{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
				{{$arch = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{ $arch = "amd64" }}
				{{- end -}}

				{{ if HasPrefix .OS "ming" -}}
				{{$os = "Windows"}}
				{{$ext = "zip"}}
				{{- end -}}

				{{.Name}}-{{$os}}-{{$arch}}.{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "sachaos",
			Repo:        "viddy",
			Name:        "viddy",
			Description: "A modern watch command. Time machine and pager etc.",
			BinaryTemplate: `
					{{$arch := .Arch}}
					{{$ext := "tar.gz"}}
	
					{{$arch := .Arch}}
					{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
					{{$arch = "arm64"}}
					{{- end -}}

					{{$osStr := ""}}
					{{ if HasPrefix .OS "ming" -}}
						{{$osStr = "Windows"}}
					{{- else if eq .OS "linux" -}}
						{{$osStr = "Linux"}}
					{{- else if eq .OS "darwin" -}}
						{{$osStr = "Darwin"}}
					{{- end -}}

					{{.Name}}_{{$osStr}}_{{$arch}}.{{$ext}}
					`,
		})

	tools = append(tools,
		Tool{
			Owner:       "temporalio",
			Repo:        "tctl",
			Name:        "tctl",
			Description: "Temporal CLI.",
			BinaryTemplate: `
						{{$os := .OS}}
						{{$arch := .Arch}}
						{{$ext := "tar.gz"}}

						{{$arch := .Arch}}
						{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
						{{$arch = "arm64"}}
						{{- else if eq .Arch "x86_64" -}}
						{{ $arch = "amd64" }}
						{{- end -}}

						{{ if HasPrefix .OS "ming" -}}
						{{$os = "windows"}}
						{{$ext = "zip"}}
						{{- end -}}

						{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.{{$ext}}
						`,
		})

	tools = append(tools,
		Tool{
			Owner:          "firecracker-microvm",
			Repo:           "firectl",
			Name:           "firectl",
			Description:    "Command-line tool that lets you run arbitrary Firecracker MicroVMs",
			BinaryTemplate: `{{.Name}}-{{.Version}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "grafana",
			Repo:        "agent",
			Name:        "grafana-agent",
			Description: "Grafana Agent is a telemetry collector for sending metrics, logs, and trace data to the opinionated Grafana observability stack.",
			URLTemplate: `
						{{$os := .OS}}
						{{$arch := .Arch}}
						{{$ext := ".zip"}}

						{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
						{{$arch = "arm64"}}
						{{- else if eq .Arch "x86_64" -}}
						{{ $arch = "amd64" }}
						{{- else if eq .Arch "armv6l" -}}
						{{ $arch = "armv6" }}
						{{- else if eq .Arch "armv7l" -}}
						{{ $arch = "armv7" }}
						{{- end -}}

						{{ if HasPrefix .OS "ming" -}}
						{{$os = "windows"}}
						{{$ext = ".exe.zip"}}
						{{- end -}}
						https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/grafana-agent-{{$os}}-{{$arch}}{{$ext}}
						`,
			BinaryTemplate: `
						{{$os := .OS}}
						{{$arch := .Arch}}
						{{$ext := ""}}

						{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
						{{$arch = "arm64"}}
						{{- else if eq .Arch "x86_64" -}}
						{{ $arch = "amd64" }}
						{{- else if eq .Arch "armv6l" -}}
						{{ $arch = "armv6" }}
						{{- else if eq .Arch "armv7l" -}}
						{{ $arch = "armv7" }}
						{{- end -}}

						{{ if HasPrefix .OS "ming" -}}
						{{$os = "windows"}}
						{{$ext = ".exe"}}
						{{- end -}}
						grafana-agent-{{$os}}-{{$arch}}{{$ext}}
						`,
		})

	tools = append(tools,
		Tool{
			Owner:       "scaleway",
			Repo:        "scaleway-cli",
			Name:        "scaleway-cli",
			Description: "Scaleway CLI is a tool to help you pilot your Scaleway infrastructure directly from your terminal.",
			BinaryTemplate: `
							{{$os := .OS}}
							{{$arch := .Arch}}
							{{$ext := ""}}
	
							{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
							{{$arch = "arm64"}}
							{{- else if eq .Arch "x86_64" -}}
							{{ $arch = "amd64" }}
							{{- end -}}
	
							{{ if HasPrefix .OS "ming" -}}
							{{$os = "windows"}}
							{{$ext = ".exe"}}
							{{- end -}}
							{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}{{$ext}}
							`,
		})

	tools = append(tools,
		Tool{
			Owner:       "anchore",
			Repo:        "syft",
			Name:        "syft",
			Description: "CLI tool and library for generating a Software Bill of Materials from container images and filesystems",
			BinaryTemplate: `
				{{$os := .OS}}
				{{$arch := .Arch}}
				{{$ext := "tar.gz"}}

				{{$arch := .Arch}}
				{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
				{{$arch = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{ $arch = "amd64" }}
				{{- end -}}

				{{ if HasPrefix .OS "ming" -}}
				{{$os = "windows"}}
				{{$ext = "zip"}}
				{{- end -}}

				{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "anchore",
			Repo:        "grype",
			Name:        "grype",
			Description: "A vulnerability scanner for container images and filesystems",
			BinaryTemplate: `
				{{$os := .OS}}
				{{$arch := .Arch}}
				{{$ext := "tar.gz"}}

				{{$arch := .Arch}}
				{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
				{{$arch = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{ $arch = "amd64" }}
				{{- end -}}

				{{ if HasPrefix .OS "ming" -}}
				{{$os = "windows"}}
				{{$ext = "zip"}}
				{{- end -}}

				{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kubernetes-sigs",
			Repo:        "cluster-api-provider-aws",
			Name:        "clusterawsadm",
			Description: "Kubernetes Cluster API Provider AWS Management Utility",
			BinaryTemplate: `
				{{$os := .OS}}
				{{$arch := .Arch}}
				{{$ext := ""}}

				{{- if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{ $arch = "amd64" }}
				{{- end -}}

				{{ if HasPrefix .OS "ming" -}}
				{{$os = "windows"}}
				{{$ext = ".exe"}}
				{{- end -}}


				clusterawsadm-{{$os}}-{{$arch}}{{$ext}}
				`,
		})

	tools = append(tools,
		Tool{
			Owner:       "schollz",
			Repo:        "croc",
			Name:        "croc",
			Description: "Easily and securely send things from one computer to another",
			BinaryTemplate: `
					{{$os := .OS}}
					{{$arch := .Arch}}
					{{$ext := "tar.gz"}}
	
					{{- if eq .OS "darwin" -}}
					{{$os = "macOS"}}
					{{- else if eq .OS "linux" -}}
					{{ $os = "Linux" }}
					{{- end -}}

					{{- if eq .Arch "aarch64" -}}
					{{$arch = "ARM64"}}
					{{- else if eq .Arch "arm64" -}}
					{{ $arch = "ARM64" }}
					{{- else if eq .Arch "x86_64" -}}
					{{ $arch = "64bit" }}
					{{- else if eq .Arch "armv7l" -}}
					{{ $arch = "ARM" }}
					{{- end -}}
	
					{{ if HasPrefix .OS "ming" -}}
					{{$os = "Windows"}}
					{{$ext = "zip"}}
					{{- end -}}
	
	
					croc_{{.VersionNumber}}_{{$os}}-{{$arch}}.{{$ext}}
					`,
		})

	// tools = append(tools,
	// 	Tool{
	// 		Owner:       "cloudnative-pg",
	// 		Repo:        "cloudnative-pg",
	// 		Name:        "kubectl-cnpg",
	// 		Description: "This plugin provides multiple commands to help you manage your CloudNativePG clusters.",
	// 		BinaryTemplate: `
	// 				{{ $os := .OS }}
	// 				{{ $arch := .Arch }}

	// 				{{- if eq .Arch "aarch64" -}}
	// 				{{ $arch = "arm64" }}
	// 				{{- else if eq .Arch "arm64" -}}
	// 				{{ $arch = "arm64" }}
	// 				{{- else if eq .Arch "armv7l" -}}
	// 				{{ $arch = "armv7" }}
	// 				{{- end -}}

	// 				{{ if HasPrefix .OS "ming" -}}
	// 				{{$os = "windows"}}
	// 				{{- end -}}

	// 				kubectl-cnpg_{{ .VersionNumber }}_{{ $os }}_{{ $arch }}.tar.gz
	// 				`,
	// 	})

	tools = append(tools,
		Tool{
			Owner:       "alexellis",
			Repo:        "fstail",
			Name:        "fstail",
			Description: "Tail modified files in a directory.",
			BinaryTemplate: `
				{{$arch := ""}}
				{{$os := ""}}
				{{$ext := ""}}
	
				{{- if eq .Arch "aarch64" -}}
				{{$arch = "-arm64"}}
				{{- else if eq .Arch "arm64" -}}
				{{$arch = "-arm64"}}
				{{- else if (or (eq .Arch "armv6l") (eq .Arch "armv7l")) -}}
				{{$arch = "-armhf"}}
				{{- end -}}
	
				{{ if eq .OS "darwin" -}}
				{{$os = "-darwin"}}
				{{ else if HasPrefix .OS "ming" -}}
				{{$ext = ".exe"}}
				{{- end -}}
				{{.Name}}{{$os}}{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "self-actuated",
			Repo:        "actions-usage",
			Name:        "actions-usage",
			Description: "Get usage insights from GitHub Actions.",
			BinaryTemplate: `
				{{$arch := ""}}
				{{$os := ""}}
				{{$ext := ""}}
	
				{{- if eq .Arch "aarch64" -}}
				{{$arch = "-arm64"}}
				{{- else if eq .Arch "arm64" -}}
				{{$arch = "-arm64"}}
				{{- else if (or (eq .Arch "armv6l") (eq .Arch "armv7l")) -}}
				{{$arch = "-armhf"}}
				{{- end -}}
	
				{{ if eq .OS "darwin" -}}
				{{$os = "-darwin"}}
				{{ else if HasPrefix .OS "ming" -}}
				{{$ext = ".exe"}}
				{{- end -}}
				{{.Name}}{{$os}}{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "self-actuated",
			Repo:        "actuated-cli",
			Name:        "actuated-cli",
			Description: "Official CLI for actuated.dev",
			BinaryTemplate: `
					{{$arch := ""}}
					{{$os := ""}}
					{{$ext := ""}}
		
					{{- if eq .Arch "aarch64" -}}
					{{$arch = "-arm64"}}
					{{- else if eq .Arch "arm64" -}}
					{{$arch = "-arm64"}}
					{{- else if (or (eq .Arch "armv6l") (eq .Arch "armv7l")) -}}
					{{$arch = "-armhf"}}
					{{- end -}}
		
					{{ if eq .OS "darwin" -}}
					{{$os = "-darwin"}}
					{{ else if HasPrefix .OS "ming" -}}
					{{$ext = ".exe"}}
					{{- end -}}
					{{.Name}}{{$os}}{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "cert-manager",
			Repo:        "cert-manager",
			Name:        "cmctl",
			Description: "cmctl is a CLI tool that helps you manage cert-manager and its resources inside your cluster.",
			BinaryTemplate: `
						{{$os := .OS}}
						{{$arch := "arm"}}
						{{$ext := "tar.gz"}}
						{{- if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
						{{$arch = "arm64"}}
						{{- else if eq .Arch "x86_64" -}}
						{{$arch = "amd64"}}
						{{- end -}}
						{{ if HasPrefix .OS "ming" -}}
						{{$os = "windows"}}
						{{$ext = "zip"}}
						{{- end -}}
						cmctl-{{$os}}-{{$arch}}.{{$ext}}
						`,
		})

	tools = append(tools,
		Tool{
			Owner:       "yt-dlp",
			Repo:        "yt-dlp",
			Name:        "yt-dlp",
			Description: "Fork of youtube-dl with additional features and fixes",
			BinaryTemplate: `
						{{$arch := ""}}
						{{$os := ""}}
						{{$ext := ""}}
			
						{{- if eq .OS "linux" -}}
							{{$os = "linux"}}
						{{- else if eq .OS "darwin" -}}
							{{$os = "macos"}}
						{{- end }}

						{{- if eq .Arch "aarch64" -}}
						{{$arch = "_aarch64"}}
						{{- else if (or (eq .Arch "armv6l") (eq .Arch "armv7l")) -}}
						{{$arch = "_armv7l"}}
						{{- end -}}
			
						{{ if HasPrefix .OS "ming" -}}
						{{$ext = ".exe"}}
						{{$arch = "x86"}}
						{{- end -}}
						{{.Name}}_{{$os}}{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "seaweedfs",
			Repo:        "seaweedfs",
			Name:        "seaweedfs",
			Description: "SeaweedFS is a fast distributed storage system for blobs, objects, files, and data lake, for billions of files!",
			URLTemplate: `
							{{$arch := ""}}
							{{$os := .OS}}
							{{$ext := ".tar.gz"}}
	
							{{- if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
							{{$arch = "arm64"}}
							{{- else if (or (eq .Arch "armv6l") (eq .Arch "armv7l")) -}}
							{{$arch = "arm"}}
							{{- else if eq .Arch "x86_64" -}}
							{{$arch = "amd64"}}
							{{- end -}}
				
							https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{$os}}_{{$arch}}{{$ext}}
							`,
			BinaryTemplate: `weed`,
		})

	tools = append(tools,
		Tool{
			Owner:       "kyverno",
			Repo:        "kyverno",
			Name:        "kyverno",
			Description: "CLI to apply and test Kyverno policies outside a cluster.",
			URLTemplate: `
				{{$arch := .Arch}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "x86_64"}}
				{{- else if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
				{{$arch = "arm64"}}
				{{- end -}}
	
				{{$os := .OS}}
				{{$extStr := "tar.gz"}}
				
				{{ if HasPrefix .OS "ming" -}}
				{{$os = "windows"}}
				{{$extStr = "zip"}}
				{{- end -}}
				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-cli_{{.Version}}_{{$os}}_{{$arch}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "replicatedhq",
			Repo:        "replicated",
			Name:        "replicated",
			Description: "CLI for interacting with the Replicated Vendor API",
			URLTemplate: `
				{{$arch := ""}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "amd64"}}
				{{- end -}}

				{{$osStr := ""}}
				{{- if eq .OS "darwin" -}}
				{{$osStr = "darwin"}}
				{{$arch = "all"}}
				{{- else if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				{{- end -}}

				{{$extStr := "tar.gz"}}

				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$osStr}}_{{$arch}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "vladimirvivien",
			Repo:        "ktop",
			Name:        "ktop",
			Description: "A top-like tool for your Kubernetes cluster.",
			URLTemplate: `
					{{$arch := .Arch}}
					{{ if eq .Arch "x86_64" -}}
					{{$arch = "amd64"}}
					{{- else if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
					{{$arch = "arm64"}}
					{{- else if eq .Arch "armv7l" -}}
					{{$arch = "armv7"}}
					{{- end -}}
		
					{{$os := .OS}}
					{{$ext := "tar.gz"}}
					
					https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{.Version}}_{{$os}}_{{$arch}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "cloud-bulldozer",
			Repo:        "kube-burner",
			Name:        "kube-burner",
			Description: "A tool aimed at stressing Kubernetes clusters by creating or deleting a high quantity of objects.",
			BinaryTemplate: `
 					{{$os := .OS}}
 					{{$arch := .Arch}}
 					{{$ext := "tar.gz"}}
 	
 					{{- if eq .OS "darwin" -}}
 						{{$os = "Darwin"}}
 					{{- else if eq .OS "linux" -}}
 						{{ $os = "Linux" }}
					{{- else if HasPrefix .OS "ming" -}}
						{{$os = "Windows"}}
						{{$ext = "zip"}}
 					{{- end -}}

 					{{- if eq .Arch "aarch64" -}}
 						{{$arch = "arm64"}}
 					{{- else if eq .Arch "arm64" -}}
 						{{ $arch = "arm64" }}
 					{{- else if eq .Arch "x86_64" -}}
 						{{ $arch = "x86_64" }}
 					{{- end -}}

 					{{.Name}}-V{{.VersionNumber}}-{{$os}}-{{$arch}}.{{$ext}}
 					`,
		})

	tools = append(tools,
		Tool{
			Owner:       "openshift",
			Repo:        "installer",
			Name:        "openshift-install",
			Description: "CLI to install an OpenShift 4.x cluster.",
			URLTemplate: `
						{{$os := .OS}}
						{{$arch := .Arch}}
						{{$ext := "tar.gz"}}
						{{$version := .VersionNumber}}
		 
						{{- if eq .OS "darwin" -}}
							{{$os = "mac"}}
						{{- end -}}
	
						{{- if eq .Arch "aarch64" -}}
							{{$arch = "-arm64"}}
						{{- else if eq .Arch "arm64" -}}
							{{ $arch = "-arm64" }}
						{{- else if eq .Arch "x86_64" -}}
							{{ $arch = "" }}
						{{- end -}}

						{{- if eq .VersionNumber "" -}}
							{{$version = "4.13.0"}}
						{{- end -}}
	
						https://mirror.openshift.com/pub/openshift-v4/clients/ocp/{{$version}}/{{.Name}}-{{$os}}{{$arch}}.tar.gz
						`,
		})

	tools = append(tools,
		Tool{
			Owner:       "openshift",
			Repo:        "oc",
			Name:        "oc",
			Description: "Client to use an OpenShift 4.x cluster.",
			URLTemplate: `
						{{$os := .OS}}
						{{$arch := .Arch}}
						{{$ext := "tar.gz"}}
						{{$version := .VersionNumber}}
			
						{{- if eq .OS "darwin" -}}
							{{$os = "mac"}}
						{{- else if HasPrefix .OS "ming" -}}
							{{$os = "windows"}}
							{{$ext = "zip"}}
						{{- end -}}
	
						{{- if eq .Arch "aarch64" -}}
							{{$arch = "-arm64"}}
						{{- else if eq .Arch "arm64" -}}
							{{ $arch = "-arm64" }}
						{{- else if eq .Arch "x86_64" -}}
							{{ $arch = "" }}
						{{- end -}}

						{{- if eq .VersionNumber "" -}}
							{{$version = "latest"}}
						{{- end -}}
	
						https://mirror.openshift.com/pub/openshift-v4/clients/ocp/{{$version}}/openshift-client-{{$os}}{{$arch}}.{{$ext}}
						`,
		})

	tools = append(tools,
		Tool{
			Owner:       "atuinsh",
			Repo:        "atuin",
			Name:        "atuin",
			Description: "Sync, search and backup shell history with Atuin.",
			URLTemplate: `
					{{$os := .OS}}
					{{$arch := .Arch}}
					{{$ext := "tar.gz"}}
					
					{{- if eq .OS "darwin" -}}
						{{$os = "apple-darwin"}}
					{{- else if eq .OS "linux" -}}
						{{$os = "unknown-linux-gnu"}}
					{{- end -}}
					
					{{- if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
						{{$arch = "x86_64"}}
					{{- end -}}
					
					https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}-{{.Version}}-{{$arch}}-{{$os}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "project-copacetic",
			Repo:        "copacetic",
			Name:        "copa",
			Description: "CLI for patching container images",
			URLTemplate: `
				{{$arch := ""}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "amd64"}}
				{{- end -}}

				{{$osStr := ""}}
				{{- if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				{{- end -}}

				{{$extStr := "tar.gz"}}

				https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$osStr}}_{{$arch}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "go-task",
			Repo:        "task",
			Name:        "task",
			Description: "A simple task runner and build tool",
			BinaryTemplate: `
					{{$os := .OS}}
					{{$arch := .Arch}}
					{{$ext := "tar.gz"}}

					{{- if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
						{{$arch = "arm64"}}
					{{- else if eq .Arch "x86_64" -}}
						{{ $arch = "amd64" }}
					{{- else if eq .Arch "armv7l" -}}
						{{ $arch = "arm" }}
					{{- end -}}

					{{ if HasPrefix .OS "ming" -}}
					{{$os = "windows"}}
					{{$ext = "zip"}}
					{{- end -}}

					{{.Name}}_{{$os}}_{{$arch}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "1password",
			Name:        "op",
			Description: "1Password CLI enables you to automate administrative tasks and securely provision secrets across development environments.",
			URLTemplate: `
				{{$os := .OS}}
				{{$arch := .Arch}}
				{{$version := .Version}}

				{{- if eq .Version "" -}}
					{{ $version = "v2.17.0" }}
				{{- end -}}

				{{- if eq .Arch "aarch64" -}}
					{{ $arch = "arm64" }}
				{{- else if eq .Arch "x86_64" -}}
					{{ $arch = "amd64" }}
				{{- else if eq .Arch "armv7l" -}}
					{{ $arch = "arm" }}
				{{- end -}}

				{{ if HasPrefix .OS "ming" -}}
				{{$os = "windows"}}
				{{- end -}}

				https://cache.agilebits.com/dist/1P/op2/pkg/{{$version}}/op_{{$os}}_{{$arch}}_{{$version}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "charmbracelet",
			Repo:        "vhs",
			Name:        "vhs",
			Description: "CLI for recording demos",
			URLTemplate: `
					{{$arch := .Arch}}
					{{ if (eq .Arch "x86_64") -}}
					{{$arch = "x86_64"}}
					{{- else if eq .Arch "aarch64" -}}
					{{$arch = "arm64"}}
					{{- end -}}

					{{$osStr := ""}}
					{{$extStr := "tar.gz"}}
					{{- if eq .OS "darwin" -}}
					{{$osStr = "Darwin"}}
					{{- else if eq .OS "linux" -}}
					{{$osStr = "Linux"}}
					{{- else if HasPrefix .OS "ming" -}}
					{{$osStr = "Windows"}}
					{{$extStr = "zip"}}
					{{- end -}}

					{{- if eq $osStr "Darwin"}}
					https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{$osStr}}_{{$arch}}.{{$extStr}}
					{{- else if or (eq $osStr "Windows") (eq $osStr "Linux") -}}
					https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_{{.VersionNumber}}_{{$osStr}}_{{$arch}}.{{$extStr}}
					{{- end -}}
					`,
		})

	tools = append(tools,
		Tool{
			Owner:       "skupperproject",
			Repo:        "skupper",
			Name:        "skupper",
			Description: "Skupper is an implementation of a Virtual Application Network, enabling rich hybrid cloud communication",
			BinaryTemplate: `
					{{$os := .OS}}
					{{$arch := .Arch}}
					{{$ext := "tgz"}}
	
					{{- if eq .OS "darwin" -}}
					{{$os = "mac"}}
					{{- else if eq .OS "linux" -}}
					{{ $os = "linux" }}
					{{- end -}}

					{{- if eq .Arch "aarch64" -}}
					{{$arch = "arm64"}}
					{{- else if eq .Arch "arm64" -}}
					{{ $arch = "arm64" }}
					{{- else if eq .Arch "x86_64" -}}
					{{ $arch = "amd64" }}
					{{- else if eq .Arch "armv7l" -}}
					{{ $arch = "arm32" }}
					{{- end -}}
	
					{{ if HasPrefix .OS "ming" -}}
					{{$os = "windows"}}
					{{$ext = "zip"}}
					{{- end -}}
	
	
					skupper-cli-{{.VersionNumber}}-{{$os}}-{{$arch}}.{{$ext}}
					`,
		})
	return tools
}
