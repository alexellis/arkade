package apps

import (
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallChart(options *types.InstallerOptions) (*types.InstallerOutput, error) {
	result := &types.InstallerOutput{}
	err := helm.AddHelmRepo(options.Helm.Repo.Name, options.Helm.Repo.URL, options.Helm.UpdateRepo, options.Helm.Helm3)
	if err != nil {
		return result, err
	}

	if err := helm.FetchChart(options.Helm.Repo.Name, options.Helm.Repo.Version, options.Helm.Helm3); err != nil {
		return result, err
	}

	if err := helm.Helm3Upgrade(options.Helm.Repo.Name, options.Namespace,
		options.Helm.ValuesFile,
		options.Helm.Repo.Version,
		options.Helm.Overrides,
		options.Helm.Wait); err != nil {
		return result, err
	}

	return result, nil
}
