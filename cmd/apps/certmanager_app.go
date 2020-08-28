// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/mod/semver"

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

	certManager.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

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

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		updateRepo, _ := certManager.Flags().GetBool("update-repo")
		err = helm.AddHelmRepo("jetstack", "https://charts.jetstack.io", updateRepo)
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

		err = helm.FetchChart("jetstack/cert-manager", version)
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

		err = helm.Helm3Upgrade("jetstack/cert-manager", namespace,
			"values.yaml",
			version,
			overrides,
			wait)

		if err != nil {
			return err
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
