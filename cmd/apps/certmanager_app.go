// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
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
	certManager.Flags().StringP("version", "v", "v1.0.4", "The version of cert-manager to install, has to be >=0.15.0")
	certManager.Flags().Bool("update-repo", true, "Update the helm repo")
	certManager.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set key=value)")

	certManager.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		wait, _ := command.Flags().GetBool("wait")

		namespace, _ := command.Flags().GetString("namespace")
		version, _ := command.Flags().GetString("version")

		if !semver.IsValid(version) {
			return fmt.Errorf("%q is not a valid semver version", version)
		}

		updateRepo, _ := certManager.Flags().GetBool("update-repo")

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

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return err
		}
		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		certmanagerOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("jetstack/cert-manager").
			WithHelmURL("https://charts.jetstack.io").
			WithOverrides(overrides).
			WithWait(wait).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(certmanagerOptions)
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
