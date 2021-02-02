package get

import "strings"

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
			Owner: "openfaas",
			Repo:  "faas-cli",
			Name:  "faas-cli",
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
			Owner:   "helm",
			Repo:    "helm",
			Name:    "helm",
			Version: "v3.2.4",
			URLTemplate: `{{$arch := "arm"}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
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
			Owner: "roboll",
			Repo:  "helmfile",
			Name:  "helmfile",
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

	// https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/darwin/amd64/kubectl
	tools = append(tools,
		Tool{
			Owner:   "kubernetes",
			Repo:    "kubernetes",
			Name:    "kubectl",
			Version: "v1.18.0",
			URLTemplate: `{{$arch := "arm"}}

{{- if eq .Arch "x86_64" -}}
{{$arch = "amd64"}}
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
			URLTemplate: `https://github.com/ahmetb/kubectx/releases/download/{{.Version}}/kubectx`,
			NoExtension: true,
		})

	tools = append(tools,
		Tool{
			Owner:       "ahmetb",
			Repo:        "kubectx",
			Name:        "kubens",
			Version:     "v0.9.1",
			URLTemplate: `https://github.com/ahmetb/kubectx/releases/download/{{.Version}}/kubens`,
			NoExtension: true,
		})

	tools = append(tools,
		Tool{
			Owner: "kubernetes-sigs",
			Repo:  "kind",
			Name:  "kind",
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
			Owner: "rancher",
			Repo:  "k3d",
			Name:  "k3d",
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
			Owner: "alexellis",
			Repo:  "k3sup",
			Name:  "k3sup",
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
			Owner: "alexellis",
			Repo:  "arkade",
			Name:  "arkade",
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
			Owner:   "bitnami-labs",
			Repo:    "sealed-secrets",
			Name:    "kubeseal",
			Version: "v0.12.4",
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
			Owner:   "inlets",
			Repo:    "inletsctl",
			Name:    "inletsctl",
			Version: "0.5.4",
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
			Name:    "osm",
			Repo:    "osm",
			Owner:   "openservicemesh",
			Version: "v0.1.0",
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
			Owner:   "linkerd",
			Repo:    "linkerd2",
			Name:    "linkerd2",
			Version: "stable-2.9.0",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}-cli-{{.Version}}-windows.exe
{{- else if eq .OS "darwin" -}}
{{.Name}}-cli-{{.Version}}-darwin
{{- else if eq .OS "linux" -}}
{{.Name}}-cli-{{.Version}}-linux
{{- end -}}
`,
		})

	tools = append(tools,
		Tool{
			Owner:   "kubernetes-sigs",
			Repo:    "kubebuilder",
			Name:    "kubebuilder",
			Version: "2.3.1",
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
			https://github.com/kubernetes-sigs/kubebuilder/releases/download/v{{.Version}}/kubebuilder_{{.Version}}_{{$osStr}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:   "kubernetes-sigs",
			Repo:    "kustomize",
			Name:    "kustomize",
			Version: "kustomize/v3.8.1",
			URLTemplate: `
	{{$osStr := ""}}
	{{- if eq .OS "linux" -}}
	{{- if eq .Arch "x86_64" -}}
	{{$osStr = "linux_amd64"}}
	{{- end -}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "darwin_amd64"}}
	{{- end -}}
	https://github.com/{{.Owner}}/{{.Repo}}/releases/download/{{.Version}}/{{.Name}}_v3.8.1_{{$osStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:   "digitalocean",
			Repo:    "doctl",
			Name:    "doctl",
			Version: "1.46.0",
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
			Owner:   "derailed",
			Repo:    "k9s",
			Name:    "k9s",
			Version: "v0.21.7",
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
			Owner:   "derailed",
			Repo:    "popeye",
			Name:    "popeye",
			Version: "v0.9.0",
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
			Version:        "0.6.35",
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

	// https://releases.hashicorp.com/terraform/0.13.1/terraform_0.13.1_linux_amd64.zip
	tools = append(tools,
		Tool{
			Owner:   "hashicorp",
			Repo:    "terraform",
			Name:    "terraform",
			Version: "0.13.1",
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
			Owner:   "hashicorp",
			Repo:    "vagrant",
			Name:    "vagrant",
			Version: "2.2.14",
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
			Owner:   "hashicorp",
			Repo:    "packer",
			Name:    "packer",
			Version: "1.6.5",
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
			Version:        "1.3.0",
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
			Owner:   "buildpacks",
			Repo:    "pack",
			Name:    "pack",
			Version: "0.14.2",
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
			Owner:   "kubernetes",
			Repo:    "kops",
			Name:    "kops",
			Version: "v1.18.2",
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
			Owner:   "kubernetes-sigs",
			Repo:    "krew",
			Name:    "krew",
			Version: "v0.4.0",
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
			Owner: "kubernetes",
			Repo:  "minikube",
			Name:  "minikube",
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
			Owner:   "stern",
			Repo:    "stern",
			Name:    "stern",
			Version: "1.13.0",
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
			Owner:   "boz",
			Repo:    "kail",
			Name:    "kail",
			Version: "0.15.0",
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
			Owner: "mikefarah",
			Repo:  "yq",
			Name:  "yq",
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
			Owner:   "aquasecurity",
			Repo:    "kube-bench",
			Name:    "kube-bench",
			Version: "0.4.0",
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
			Owner:   "gohugoio",
			Repo:    "hugo",
			Name:    "hugo",
			Version: "0.79.0",
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
			Owner:   "docker",
			Repo:    "compose",
			Name:    "docker-compose",
			Version: "1.27.4",
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
			Owner: "open-policy-agent",
			Repo:  "opa",
			Name:  "opa",
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
			Owner: "minio",
			Repo:  "mc",
			Name:  "mc",
			URLTemplate: `{{$arch := .Arch}}
			{{ if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "armv6l" -}}
			{$arch = "arm"}}
			{{- else if eq .Arch "armv7l" -}}
			{$arch = "arm"}}
			{{- else if eq .Arch "aarch64" -}}
			{$arch = "arm64"}}
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
			Owner:   "nats-io",
			Repo:    "natscli",
			Name:    "nats",
			Version: "0.0.21",
			URLTemplate: `{{$arch := .Arch}}
			{{ if eq .Arch "x86_64" -}}
			{{$arch = "amd64"}}
			{{- else if eq .Arch "armv6l" -}}
			{$arch = "arm6"}}
			{{- else if eq .Arch "armv7l" -}}
			{$arch = "arm7"}}
			{{- else if eq .Arch "aarch64" -}}
			{$arch = "arm64"}}
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

	return tools
}
