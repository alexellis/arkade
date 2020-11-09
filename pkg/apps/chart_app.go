package apps

import (
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallChart(options *types.InstallerOptions) (*types.InstallerOutput, error) {
	result := &types.InstallerOutput{}

	if err := config.SetKubeconfig(options.KubeconfigPath); err != nil {
		return nil, err
	}

	for _, secret := range options.Secrets {
		if err := k8s.CreateSecret(secret); err != nil {
			return nil, err
		}
	}

	err := helm.AddHelmRepo(options.Helm.Repo.Name, options.Helm.Repo.URL, options.Helm.UpdateRepo)
	if err != nil {
		return result, err
	}

	if err := helm.FetchChart(options.Helm.Repo.Name, options.Helm.Repo.Version); err != nil {
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
