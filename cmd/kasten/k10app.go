// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

// kasten contains a suite of Sponsored Apps for arkade
package kasten

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallK10() *cobra.Command {
	var k10cmd = &cobra.Command{
		Use:   "k10",
		Short: "Install K10",
		Long: `Kasten K10 by Veeam is purpose-built for Kubernetes backup and restore.

Note: K10 performs best if your cluster supports a CSI driver, see the following command:
  kubectl get storageclasses
`,
		Example: `  arkade install k10
  arkade install k10 --help
  arkade install k10 \
    --set eula.accept=true \
    --set clusterName=my-k10 \
    --set prometheus.server.enabled=false

See also: all helm chart options:
https://docs.kasten.io/latest/install/advanced.html#complete-list-of-k10-helm-options`,
		SilenceUsage: true,
	}

	k10cmd.Flags().StringP("namespace", "n", "kasten-io", "The namespace used for installation")
	k10cmd.Flags().Bool("update-repo", true, "Update the helm repo")
	k10cmd.Flags().Bool("eula", false, "Accept the EULA")
	k10cmd.Flags().Bool("token-auth", false, "Change to token mode for authentication")

	k10cmd.Flags().Bool("prometheus", true, "Enable Prometheus server")
	k10cmd.Flags().Bool("grafana", true, "Enable Grafana server")

	k10cmd.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	k10cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if _, err := cmd.Flags().GetBool("eula"); err != nil {
			if err != nil {
				return fmt.Errorf("error with \"eula\" flag %w", err)
			}
		}
		if _, err := cmd.Flags().GetBool("prometheus"); err != nil {
			if err != nil {
				return fmt.Errorf("error with \"prometheus\" flag %w", err)
			}
		}
		if _, err := cmd.Flags().GetBool("grafana"); err != nil {
			if err != nil {
				return fmt.Errorf("error with \"grafana\" flag %w", err)
			}
		}
		if _, err := cmd.Flags().GetBool("token-auth"); err != nil {
			if err != nil {
				return fmt.Errorf("error with \"token-auth\" flag %w", err)
			}
		}
		return nil
	}

	k10cmd.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}
		eula, _ := command.Flags().GetBool("eula")
		tokenAuth, _ := command.Flags().GetBool("token-auth")
		grafana, _ := command.Flags().GetBool("grafana")
		prometheus, _ := command.Flags().GetBool("prometheus")

		wait, _ := command.Flags().GetBool("wait")
		namespace, _ := command.Flags().GetString("namespace")
		overrides := map[string]string{
			"prometheus.server.enabled": strconv.FormatBool(prometheus),
			"eula.accept":               strconv.FormatBool(eula),
			"auth.tokenAuth.enabled":    strconv.FormatBool(tokenAuth),
			"grafana.enabled":           strconv.FormatBool(grafana),
		}
		if command.Flags().Changed("cluster-name") {
			v, _ := command.Flags().GetString("cluster-name")
			overrides["clusterName"] = v
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		k10cmdOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kasten/k10").
			WithHelmURL("https://charts.kasten.io/").
			WithOverrides(overrides).
			WithWait(wait).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(k10cmdOptions)

		if err != nil {
			return err
		}

		fmt.Println(k10InstallCmd)

		return nil
	}

	return k10cmd
}

const k10InfoMsg = `# The K10 app has been installed

# You may also need to install pre-requisites and configure a 
# CSI driver for your cluster.

https://docs.kasten.io/latest/install/storage.html

# The app may take a few moments to come up, run the following to
# wait for it:

kubectl rollout status -n kasten-io deploy/frontend-svc

# Then access the dashboard via:

kubectl -n kasten-io port-forward service/gateway 8080:8000

http://127.0.0.1:8080/k10/#/

# Find out your next steps here:

https://docs.kasten.io/latest/install/install.html`

const k10InstallCmd = `=======================================================================
= k10 has been installed.                                   =
=======================================================================` +
	"\n\n" + k10InfoMsg + "\n\n"

func mergeFlags(existingMap map[string]string, setOverrides []string) error {
	for _, setOverride := range setOverrides {
		flag := strings.Split(setOverride, "=")
		if len(flag) != 2 {
			return fmt.Errorf("incorrect format for custom flag `%s`", setOverride)
		}
		existingMap[flag[0]] = flag[1]
	}
	return nil
}
