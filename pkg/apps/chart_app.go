package apps

import (
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/commands"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallChart(options *types.InstallerOptions) error {
	if options.Namespace != "default" {
		if err := commands.CreateNamespace(options.Namespace); err != nil {
			return err
		}
	}

	if err := config.SetKubeconfig(options.KubeconfigPath); err != nil {
		return nil, err
	}

	err := helm.AddHelmRepo(options.Helm.Repo.Name, options.Helm.Repo.URL, options.Helm.UpdateRepo)
	if err != nil {
		return err
	}

	if err := helm.FetchChart(options.Helm.Repo.Name, options.Helm.Repo.Version); err != nil {
		return err
	}

	// Preinstall commands
	if len(options.PreChartCommands) > 0 {
		for _, preCmd := range options.PreChartCommands {
			if err := preCmd(); err != nil {
				return err
			}
		}
	}

	if err := helm.Helm3Upgrade(options.Helm.Repo.Name, options.Namespace,
		options.Helm.ValuesFile,
		options.Helm.Repo.Version,
		options.Helm.Overrides,
		options.Helm.Wait); err != nil {
		return err
	}

	// Postinstall commands
	if len(options.PostChartCommands) == 0 {
		return nil
	}

	for _, postCmd := range options.PostChartCommands {
		err := postCmd()
		if err != nil {
			return err
		}
	}

	return nil
}
