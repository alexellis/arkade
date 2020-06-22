// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallLoki() *cobra.Command {
	var lokiApp = &cobra.Command{
		Use:          "loki",
		Short:        "Install Loki for monitoring and tracing",
		Long:         "Install Loki, part of the Grafana products for Logging and Tracing",
		Example:      "arkade install loki",
		SilenceUsage: true,
	}

	lokiApp.Flags().StringP("namespace", "n", "default", "The namespace to install loki (default: default")
	lokiApp.Flags().Bool("update-repo", true, "Update the helm repo")
	lokiApp.Flags().Bool("persistence", false, "Use a 10Gi Persistent Volume to store data")
	lokiApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set grafana.enabled=true)")
	lokiApp.Flags().Bool("grafana", false, "Install Grafana alongside Loki (default: false)")

	lokiApp.RunE = func(command *cobra.Command, args []string) error {
		helm3 := true

		namespace, _ := lokiApp.Flags().GetString("namespace")
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

		persistence, _ := lokiApp.Flags().GetBool("persistence")
		installGrafana, _ := lokiApp.Flags().GetBool("grafana")

		overrides := map[string]string{}

		if installGrafana {
			overrides["grafana.enabled"] = "true"
		}
		if persistence {
			overrides["loki.persistence.enabled"] = "true"
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		lokiOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("loki/loki-stack").
			WithHelmURL("https://grafana.github.io/loki/charts").
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

		println(lokiInstallMsg)
		return nil
	}

	return lokiApp
}

const LokiInfoMsg = `# Get started with loki here:
# https://github.com/grafana/loki/blob/master/docs/README.md

# See how to integrate loki with Grafana here
# https://github.com/grafana/loki/blob/master/docs/getting-started/grafana.md

# Check loki's logs with:

kubectl logs svc/loki-stack

kubectl logs svc/loki-stack-headless


# If you installed with Grafana you can access the dashboard with the username "admin" and password shown below
 # To get password
 kubectl get secret loki-stack-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
 
 # Forward traffic to your localhost
 kubectl port-forward service/loki-stack-grafana 3000:80

`

const lokiInstallMsg = `=======================================================================
= loki has been installed.                                   =
=======================================================================` +
	"\n\n" + LokiInfoMsg + "\n\n" + pkg.ThanksForUsing
