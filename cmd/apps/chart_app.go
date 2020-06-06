// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

const defaultVersion = "" // If we don't set version then we get latest

func MakeInstallChart() *cobra.Command {
	var chartCmd = &cobra.Command{
		Use:   "chart",
		Short: "Install the specified helm chart",
		Long: `Install the specified helm chart without using tiller.
Note: You may need to install a CRD or run other additional steps
before using the generic helm chart installer command.`,
		Example: `  arkade install chart --repo-name stable/nginx-ingress \
     --set controller.service.type=NodePort
  arkade install chart --repo-name inlets/inlets-operator \
     --repo-url https://inlets.github.io/inlets-operator/`,
		SilenceUsage: true,
	}

	chartCmd.Flags().StringP("namespace", "n", "default", "The namespace to install the chart")
	chartCmd.Flags().String("repo", "", "The chart repo to install from")
	chartCmd.Flags().String("values-file", "", "Give the values.yaml file to use from the upstream chart repo")
	chartCmd.Flags().String("repo-name", "", "Chart name")
	chartCmd.Flags().String("repo-url", "", "Chart repo")

	chartCmd.Flags().StringArray("set", []string{}, "Set individual values in the helm chart")

	chartCmd.RunE = func(command *cobra.Command, args []string) error {
		chartRepoName, _ := command.Flags().GetString("repo-name")
		chartRepoURL, _ := command.Flags().GetString("repo-url")

		chartName := chartRepoName
		if index := strings.Index(chartRepoName, "/"); index > -1 {
			chartName = chartRepoName[index+1:]
		}

		chartPrefix := chartRepoName
		if index := strings.Index(chartRepoName, "/"); index > -1 {
			chartPrefix = chartRepoName[:index]
		}

		if len(chartRepoName) == 0 {
			return fmt.Errorf("--repo-name required")
		}

		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		namespace, _ := command.Flags().GetString("namespace")

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, false)
		if err != nil {
			return err
		}

		if len(chartRepoURL) > 0 {
			err = helm.AddHelmRepo(chartPrefix, chartRepoURL, true, false)
			if err != nil {
				return err
			}
		}

		res, kcErr := k8s.KubectlTask("get", "namespace", namespace)

		if kcErr != nil {
			return err
		}

		if res.ExitCode != 0 {
			err = k8s.Kubectl("create", "namespace", namespace)
			if err != nil {
				return err
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")

		err = helm.FetchChart(chartRepoName, defaultVersion, false)
		if err != nil {
			return err
		}

		outputPath := path.Join(chartPath, "chart/rendered")

		setMap := map[string]string{}
		setVals, _ := chartCmd.Flags().GetStringArray("set")

		for _, setV := range setVals {
			var k string
			var v string

			if index := strings.Index(setV, "="); index > -1 {
				k = setV[:index]
				v = setV[index+1:]
				setMap[k] = v
			}
		}

		err = helm.TemplateChart(chartPath, chartName, namespace, outputPath, "values.yaml", setMap)
		if err != nil {
			return err
		}

		err = k8s.Kubectl("apply", "--namespace", namespace, "-R", "-f", outputPath)
		if err != nil {
			return err
		}

		fmt.Println(
			`=======================================================================
chart ` + chartRepoName + ` installed.
=======================================================================
		
` + pkg.ThanksForUsing)

		return nil
	}

	return chartCmd
}
