package apps

import (
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallRedis(options *types.InstallerOptions) (*types.InstallerOutput, error) {
	result := &types.InstallerOutput{}

	k8s.KubectlTask("create", "namespace", options.Namespace)

	err := helm.AddHelmRepo(options.Helm.Repo.Name, options.Helm.Repo.URL, options.Helm.UpdateRepo, options.Helm.Helm3)
	if err != nil {
		return result, err
	}

	if err := helm.FetchChart(options.Helm.Repo.Name, options.Helm.Repo.Version, options.Helm.Helm3); err != nil {
		return result, err
	}

	if err := helm.Helm3Upgrade(options.Helm.Repo.Name, options.Namespace,
		"values.yaml",
		options.Helm.Repo.Version,
		options.Helm.Overrides,
		options.Helm.Wait); err != nil {
		return result, err
	}

	return result, nil
}
