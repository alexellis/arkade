// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"strings"

	"github.com/alexellis/arkade/pkg/config"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
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
	certManager.Flags().StringP("version", "v", "v1.5.4", "The version of cert-manager to install, has to be >= v1.0.0")
	certManager.Flags().Bool("update-repo", true, "Update the helm repo")
	certManager.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set key=value)")
	certManager.Flags().StringArray("dns-server", []string{}, "Use custom flags or override existing flags \n(example --dns-servers key=value)")

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
		overrides["installCRDs"] = "true"

		dnsServers, err := command.Flags().GetStringArray("dns-server")
		if err != nil {
			return err
		}

		if len(dnsServers) > 0 {

			for _, v := range dnsServers {
				if !strings.Contains(v, ":") {
					return fmt.Errorf("dns-server need a a specific port i.e. 8.8.8.8:53")
				}
			}

			st := ""
			if len(dnsServers) > 1 {
				st = `"` + strings.Join(dnsServers, ",") + `"`
			} else {
				st = dnsServers[0]
			}

			overrides["dns01RecursiveNameservers"] = st
			overrides["dns01RecursiveNameserversOnly"] = "true"
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
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
	"\n\n" + CertManagerInfoMsg + "\n\n" + pkg.SupportMessageShort
