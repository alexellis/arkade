package types

import "github.com/alexellis/arkade/pkg/config"

type InstallerOptions struct {
	Namespace      string
	KubeconfigPath string
	NodeArch       string
	Helm           *HelmConfig
	Verbose        bool
}

type HelmConfig struct {
	Repo       *HelmRepo
	Helm3      bool
	HelmPath   string
	Overrides  map[string]string
	UpdateRepo bool
	Wait       bool
	ValuesFile string
}

type HelmRepo struct {
	Name    string
	URL     string
	Version string
}

type InstallerOutput struct {
}

func (o *InstallerOptions) WithKubeconfigPath(path string) *InstallerOptions {
	o.KubeconfigPath = path
	return o
}

func (o *InstallerOptions) WithNamespace(namespace string) *InstallerOptions {
	o.Namespace = namespace
	return o
}

func (o *InstallerOptions) WithHelmPath(helmPath string) *InstallerOptions {
	o.Helm.HelmPath = helmPath
	return o
}

func (o *InstallerOptions) WithWait(wait bool) *InstallerOptions {
	o.Helm.Wait = wait
	return o
}

func (o *InstallerOptions) WithHelmRepo(s string) *InstallerOptions {
	o.Helm.Repo.Name = s
	return o
}

func (o *InstallerOptions) WithHelmURL(s string) *InstallerOptions {
	o.Helm.Repo.URL = s
	return o
}

func (o *InstallerOptions) WithHelmUpdateRepo(update bool) *InstallerOptions {
	o.Helm.UpdateRepo = update
	return o
}

func (o *InstallerOptions) WithOverrides(overrides map[string]string) *InstallerOptions {
	o.Helm.Overrides = overrides
	return o
}

func DefaultInstallOptions() *InstallerOptions {
	return &InstallerOptions{
		Namespace:      "default",
		KubeconfigPath: config.GetDefaultKubeconfig(),
		NodeArch:       "x86_64",
		Helm: &HelmConfig{
			Repo: &HelmRepo{
				Version: "",
			},
			ValuesFile: "values.yaml",
			Helm3:      true,
			Wait:       false,
		},
		Verbose: false,
	}
}
