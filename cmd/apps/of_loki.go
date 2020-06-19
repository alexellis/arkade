// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallOpenFaaSLoki() *cobra.Command {
	var OpenFaaSlokiApp = &cobra.Command{
		Use:          "openfaas-loki",
		Short:        "Install Loki-OpenFaaS and Configure Loki logs provider for OpenFaaS",
		Long:         "Install Loki-OpenFaaS and Configure Loki logs provider for OpenFaaS",
		Example:      "arkade install openfaas-loki",
		SilenceUsage: true,
	}

	OpenFaaSlokiApp.Flags().StringP("namespace", "n", "default", "The namespace to install loki (default: default")
	OpenFaaSlokiApp.Flags().Bool("update-repo", true, "Update the helm repo")
	OpenFaaSlokiApp.Flags().String("openfaas-namespace", "openfaas", "set the namespace that OpenFaaS is installed into")
	OpenFaaSlokiApp.Flags().String("loki-url", "http://loki-stack.default:3100", "set the loki url (default http://loki-stack.default:3100)")
	OpenFaaSlokiApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set grafana.enabled=true)")

	OpenFaaSlokiApp.RunE = func(command *cobra.Command, args []string) error {
		helm3 := true

		namespace, _ := OpenFaaSlokiApp.Flags().GetString("namespace")
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		log.Printf("Client: %s, %s\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		if err := os.Setenv("HELM_HOME", path.Join(userPath, ".helm")); err != nil {
			return err
		}

		openfaasNamespace, _ := OpenFaaSlokiApp.Flags().GetString("openfaas-namespace")
		lokiURL, _ := OpenFaaSlokiApp.Flags().GetString("loki-url")

		overrides := map[string]string{}
		overrides["lokiURL"] = lokiURL

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		lokiOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("lucas/openfaas-loki").
			WithHelmURL("https://lucasroesler.com/openfaas-loki").
			WithOverrides(overrides)

		if command.Flags().Changed("kubeconfig") {
			kubeconfigPath, _ := command.Flags().GetString("kubeconfig")
			lokiOptions.WithKubeconfigPath(kubeconfigPath)
		}

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(lokiOptions)
		if err != nil {
			return err
		}

		// Post install config of openfaas-loki
		k8s.Kubectl("-n", openfaasNamespace,
			"set", "env", "deployment/gateway",
			"-c", "gateway",
			"-e", fmt.Sprintf("logs_provider_url=http://openfaas-loki.%s:9191/", namespace))

		println(lokiOFInstallMsg)
		return nil
	}

	return OpenFaaSlokiApp
}

const LokiOFInfoMsg = `# Get started with openfaas-loki here:

# If you are authenticated with your openfaas gateway
faas-cli logs

# If you installed loki with grafana, with 'arkade install loki --grafana'
# You can use the grafana dashboard to see the OpenFaaS Logs, you can see 
# how to get your grafana password with 'arkade info loki'

# We have automatically configured OpenFaaS to use the Loki logs URL, you can set 'gateway.logsProviderURL'
# When installing openfaas with Helm or use the '--log-provider-url' flag in arkade.
# The url is in the format 'http://loki-stack.namespace:3100/' (where namespace is the installed namespace for loki-stack)
`

const lokiOFInstallMsg = `=======================================================================
= OpenFaaS loki has been installed.                                   =
=======================================================================` +
	"\n\n" + LokiOFInfoMsg + "\n\n" + pkg.ThanksForUsing
