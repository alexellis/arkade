package types

import "github.com/alexellis/arkade/pkg/config"

type InstallerOptions struct {
	Namespace       string
	CreateNamespace bool
	KubeconfigPath  string
	NodeArch        string
	Helm            *HelmConfig
	Verbose         bool
	Secrets         []K8sSecret
}

type K8sSecret struct {
	Type       string
	Name       string
	SecretData []SecretsData
	Namespace  string
}

type SecretsData struct {
	Type  string // file or literal
	Key   string
	Value string
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

func (o *InstallerOptions) WithWait(wait bool) *InstallerOptions {
	o.Helm.Wait = wait
	return o
}

func (o *InstallerOptions) WithHelmRepo(s string) *InstallerOptions {
	o.Helm.Repo.Name = s
	return o
}

func (o *InstallerOptions) WithHelmRepoVersion(s string) *InstallerOptions {
	o.Helm.Repo.Version = s
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

func (o *InstallerOptions) WithValuesFile(filename string) *InstallerOptions {
	o.Helm.ValuesFile = filename
	return o
}

func (o *InstallerOptions) WithSecret(secret K8sSecret) *InstallerOptions {
	o.Secrets = append(o.Secrets, secret)
	return o
}

func (o *InstallerOptions) WithInstallNamespace(b bool) *InstallerOptions {
	o.CreateNamespace = b
	return o
}

func DefaultInstallOptions() *InstallerOptions {
	return &InstallerOptions{
		Namespace:       "default",
		KubeconfigPath:  config.GetDefaultKubeconfig(),
		NodeArch:        "x86_64",
		CreateNamespace: false,
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

func NewGenericSecret(name, namespace string, secretData []SecretsData) K8sSecret {
	return K8sSecret{
		Type:       KubernetesGenericSecret,
		Name:       name,
		Namespace:  namespace,
		SecretData: secretData,
	}
}

const KubernetesGenericSecret = "generic"
const StringLiteralSecret = "string-literal"
const FromFileSecret = "from-file"
