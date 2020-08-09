package get

func MakeTools() []Tool {
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
			Version:     "v0.9.0",
			URLTemplate: `https://github.com/ahmetb/kubectx/releases/download/{{.Version}}/kubectx`,
			// Author recommends to keep using Bash version in this release https://github.com/ahmetb/kubectx/releases/tag/v0.9.0
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
			Repo:    "linkerd",
			Name:    "linkerd2",
			Version: "stable-2.8.1",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
{{.Name}}-cli-{{.Version}}-windows.exe
{{- else if eq .OS "darwin" -}}
{{.Name}}-cli-{{.Version}}-darwin
{{- else if eq .OS "linux" -}}
{{.Name}}-cli-{{.Version}}-linux
{{- end -}}
`,
		})

	return tools
}
