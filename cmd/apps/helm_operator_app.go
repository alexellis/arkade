// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func MakeInstallHelmOperator() *cobra.Command {
	var helmOperator = &cobra.Command{
		Use:          "helm-operator",
		Short:        "Install helm-operator",
		Long:         "Install helm-operator",
		Example:      "arkade install helm-operator --namespace default",
		SilenceUsage: true,
	}

	helmOperator.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	helmOperator.Flags().Bool("update-repo", true, "Update the helm repo")
	helmOperator.Flags().Bool("helm3", true, "Use helm3, if set to false uses helm2")
	helmOperator.Flags().StringP("version", "v", "v1.2.0", "The version of helm-operator to install")

	helmOperator.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")

		namespace, _ := command.Flags().GetString("namespace")
		version, _ := command.Flags().GetString("version")
		helm3, _ := command.Flags().GetBool("helm3")
		updateRepo, _ := helmOperator.Flags().GetBool("update-repo")

		if !semver.IsValid(version) {
			return fmt.Errorf("%q is not a valid semver version", version)
		}

		if namespace != "default" {
			return fmt.Errorf(`to override the "default", install via tiller`)
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		helmRepo := "fluxcd/helm-operator"
		if !helm3 {
			helmRepo = "helm-operator/rendered"
		}

		err = helm.FetchChart(helmRepo, defaultVersion)
		if err != nil {
			return err
		}

		overrides := map[string]string{}
		overrides["helm.versions"] = "v3"

		arch := k8s.GetNodeArchitecture()
		if arch != "amd64" {
			return fmt.Errorf("This chart does not support %s", arch)
		}

		installOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo(helmRepo).
			WithHelmURL("https://charts.fluxcd.io").
			WithOverrides(overrides).
			WithWait(wait).
			WithHelmUpdateRepo(updateRepo)

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
			installOptions.WithKubeconfigPath(kubeConfigPath)
		}

		_, err = apps.MakeInstallChart(installOptions)
		if err != nil {
			return err
		}

		fmt.Println(HelmOperatorInfoMsg)

		return nil
	}

	return helmOperator
}

const HelmOperatorInfoMsg = `# The helm-operator has been configured
For example you can install kubernetes nginx-ingress this way:
This conversation was marked as resolved by aidun
kubectl apply -f - <<EOF
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: nginx-ingress
  namespace: kube-system
spec:
  releaseName: nginx-ingress
  targetNamespace: kube-system
  timeout: 300
  resetValues: false
  wait: false
  forceUpgrade: false
  chart:
    repository: https://kubernetes-charts.storage.googleapis.com
    name: nginx-ingress
    version: 1.30.0
EOF
After some time you will see the helm-release:
kubectl get helmreleases
or 
helm list -A
# Find out more at:
# https://docs.fluxcd.io/projects/helm-operator/en/latest/references/helmrelease-custom-resource.html`
