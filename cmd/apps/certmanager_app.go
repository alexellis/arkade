// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"golang.org/x/mod/semver"
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

func MakeInstallCertManager() *cobra.Command {
	var certManager = &cobra.Command{
		Use:          "cert-manager",
		Short:        "Install cert-manager",
		Long:         "Install cert-manager for TLS certificates management",
		Example:      "arkade install cert-manager",
		SilenceUsage: true,
	}

	certManager.Flags().StringP("namespace", "n", "cert-manager", "The namespace to install cert-manager")
	certManager.Flags().StringP("version", "v", "v0.15.2", "The version of cert-manager to install, has to be >=0.15.0")
	certManager.Flags().Bool("update-repo", true, "Update the helm repo")
	certManager.Flags().Bool("helm3", true, "Use helm3, if set to false uses helm2")

	certManager.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)
		helm3, _ := command.Flags().GetBool("helm3")

		if helm3 {
			fmt.Println("Using helm3")
		}
		namespace, _ := command.Flags().GetString("namespace")
		version, _ := command.Flags().GetString("version")

		if !semver.IsValid(version) {
			return fmt.Errorf("%q is not a valid semver version", version)
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		updateRepo, _ := certManager.Flags().GetBool("update-repo")
		err = helm.AddHelmRepo("jetstack", "https://charts.jetstack.io", updateRepo, helm3)
		if err != nil {
			return err
		}

		nsRes, nsErr := k8s.KubectlTask("create", "namespace", namespace)
		if nsErr != nil {
			return nsErr
		}

		if nsRes.ExitCode != 0 {
			fmt.Printf("[Warning] unable to create namespace %s, may already exist: %s", namespace, nsRes.Stderr)
		}

		chartPath := path.Join(os.TempDir(), "charts")

		err = helm.FetchChart("jetstack/cert-manager", version, helm3)
		if err != nil {
			return err
		}

		overrides := map[string]string{}

		// if <0.15 install CRDs using kubectl else use Helm
		if semver.Compare(version, "v0.15.0") < 0 {
			log.Printf("Applying CRD\n")
			crdsURL := fmt.Sprintf("https://raw.githubusercontent.com/jetstack/cert-manager/release-%s/deploy/manifests/00-crds.yaml", strings.Replace(semver.MajorMinor(version), "v", "", -1))
			res, err := k8s.KubectlTask("apply", "--validate=false", "-f",
				crdsURL)
			if err != nil {
				return err
			}

			if res.ExitCode > 0 {
				return fmt.Errorf("error applying CRD from: %s, error: %s", crdsURL, res.Stderr)
			}
		} else {
			overrides["installCRDs"] = "true"
		}

		outputPath := path.Join(chartPath, "cert-manager/rendered")
		if helm3 {
			err := helm.Helm3Upgrade("jetstack/cert-manager", namespace,
				"values.yaml",
				version,
				overrides,
				wait)

			if err != nil {
				return err
			}
		} else {
			err = helm.TemplateChart(chartPath, "cert-manager", namespace, outputPath, "values.yaml", nil)
			if err != nil {
				return err
			}

			applyRes, applyErr := k8s.KubectlTask("apply", "-R", "-f", outputPath)
			if applyErr != nil {
				return applyErr
			}

			if applyRes.ExitCode > 0 {
				return fmt.Errorf("error applying templated YAML files, error: %s", applyRes.Stderr)
			}
		}

		fmt.Println(certManagerInstallMsg)

		return nil
	}

	return certManager
}

const CertManagerInfoMsg = `# Get started with cert-manager here:
# https://docs.cert-manager.io/en/latest/tutorials/acme/http-validation.html`

const certManagerInstallMsg = `=======================================================================
= cert-manager  has been installed.                                   =
=======================================================================` +
	"\n\n" + CertManagerInfoMsg + "\n\n" + pkg.ThanksForUsing
