package get

func MakeTools() []Tool {
	tools := []Tool{
		{
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
{{- end -}}`,
		},
		//https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/darwin/amd64/kubectl
		{
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

https://storage.googleapis.com/kubernetes-release/release/{{.Version}}/bin/{{$os}}/{{$arch}}/kubectl{{$ext}}`,
		},
		{
			Owner:       "ahmetb",
			Repo:        "kubectx",
			Name:        "kubectx",
			Version:     "v0.9.0",
			URLTemplate: `https://github.com/ahmetb/kubectx/releases/download/{{.Version}}/kubectx`,
			// Author recommends to keep using Bash version in this release https://github.com/ahmetb/kubectx/releases/tag/v0.9.0
			NoExtension: true,
		},
		{
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
		},
		{
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
		},
		{
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
		},
		{
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
		},
	}
	return tools
}
