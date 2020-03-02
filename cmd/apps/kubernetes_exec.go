package apps

import (
	"github.com/alexellis/arkade/pkg/kubernetes"
	execute "github.com/alexellis/go-execute/pkg/v1"

	"github.com/alexellis/arkade/pkg/helm"
)

// stay here for legacy
func fetchChart(path, chart, version string, helm3 bool) error {
	return helm.FetchChart(path, chart, version, helm3)
}

// stay here for legacy
func getNodeArchitecture() string {
	return kubernetes.GetNodeArchitecture()
}

// stay here for legacy
func helm3Upgrade(basePath, chart, namespace, values, version string, overrides map[string]string, wait bool) error {
	return helm.Helm3Upgrade(basePath, chart, namespace, values, version, overrides, wait)
}

// stay here for legacy
func templateChart(basePath, chart, namespace, outputPath, values string, overrides map[string]string) error {
	return helm.TemplateChart(basePath, chart, namespace, outputPath, values, overrides)
}

// stay here for legacy
func addHelmRepo(name, url string, helm3 bool) error {
	return helm.AddHelmRepo(name, url, helm3)
}

// stay here for legacy
func updateHelmRepos(helm3 bool) error {
	return helm.UpdateHelmRepos(helm3)
}

// stay here for legacy
func kubectlTask(parts ...string) (execute.ExecResult, error) {
	return kubernetes.KubectlTask(parts...)
}

// stay here for legacy
func kubectl(parts ...string) error {
	return kubernetes.Kubectl(parts...)
}

// stay here for legacy
func getDefaultKubeconfig() string {
	return kubernetes.GetDefaultKubeconfig()
}
