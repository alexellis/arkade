package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallChart(options *types.InstallerOptions) (*types.InstallerOutput, error) {
	result := &types.InstallerOutput{}

	if err := config.SetKubeconfig(options.KubeconfigPath); err != nil {
		return nil, err
	}

	if options.CreateNamespace {
		if err := k8s.CreateNamespace(options.Namespace); err != nil {
			return nil, err
		}
	}

	for _, secret := range options.Secrets {
		if err := k8s.CreateSecret(secret); err != nil {
			return nil, err
		}
	}

	userPath, err := config.InitUserDir()
	if err != nil {
		return nil, err
	}

	clientArch, clientOS := env.GetClientArch()

	fmt.Printf("Client: %s, %s\n", clientArch, clientOS)

	log.Printf("User dir established as: %s\n", userPath)

	os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

	_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
	if err != nil {
		return nil, err
	}

	installer := helm.Helm3OCIUpgrade
	name := options.Helm.Repo.URL
	if !helm.IsOCI(options.Helm.Repo.URL) {
		if helm.IsURL(options.Helm.Repo.URL) && strings.HasSuffix(options.Helm.Repo.URL, ".tgz") {
			// Chart given as a direct archive URL: fetch, untar and
			// install from the local chart dir. Helm3Upgrade resolves
			// <chart> as <tempdir>/charts/<name-after-slash>, and the
			// untarred dir drops the version suffix.
			if err := helm.FetchChart(options.Helm.Repo.URL, options.Helm.Repo.Version); err != nil {
				return result, err
			}
			installer = helm.Helm3Upgrade
			chartName := strings.TrimSuffix(path.Base(options.Helm.Repo.URL), ".tgz")
			if index := strings.LastIndex(chartName, "-"); index > -1 {
				chartName = chartName[:index]
			}
			name = chartName
		} else {
			if err = helm.AddHelmRepo(options.Helm.Repo.Name, options.Helm.Repo.URL, options.Helm.UpdateRepo); err != nil {
				return result, err
			}

			if err := helm.FetchChart(options.Helm.Repo.Name, options.Helm.Repo.Version); err != nil {
				return result, err
			}
			installer = helm.Helm3Upgrade
			name = options.Helm.Repo.Name
		}
	}

	if err := installer(
		name,
		options.Namespace,
		options.Helm.ValuesFiles,
		options.Helm.Repo.Version,
		options.Helm.Overrides,
		options.Helm.Wait); err != nil {
		return result, err
	}

	return result, nil
}
