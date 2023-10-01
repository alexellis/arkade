// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
)

func MakeInstallHarbor() *cobra.Command {
	var harbor = &cobra.Command{
		Use:   "harbor",
		Short: "Install harbor",
		Long:  `Install harbor`,
		Example: `  arkade install harbor
  # with ingress and custom domain
  arkade install harbor --ingress=true --domain=example.com`,
		SilenceUsage: true,
	}

	harbor.Flags().StringP("namespace", "n", "harbor", "The namespace to install the chart")
	harbor.Flags().Bool("update-repo", true, "Update the helm repo")
	harbor.Flags().Bool("ingress", false, "Enable ingress")
	harbor.Flags().Bool("trivy", false, "Enable trivy")
	harbor.Flags().Bool("notary", false, "Enable notary")
	harbor.Flags().String("domain", "", "Set ingress domain")
	harbor.Flags().Bool("persistence", false, "Enable harbor server persistence")
	harbor.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image.tag=1.11.2)")

	harbor.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("update-repo")
		if err != nil {
			return err
		}

		ingressEnabled, err := cmd.Flags().GetBool("ingress")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("trivy")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("notary")
		if err != nil {
			return err
		}

		ingressDomain, err := cmd.Flags().GetString("domain")
		if err != nil {
			return err
		}

		if !ingressEnabled && ingressDomain != "" {
			return fmt.Errorf("--domain option should be used only with --ingress=true")
		}

		_, err = cmd.Flags().GetBool("persistence")
		if err != nil {
			return err
		}

		return nil
	}

	harbor.RunE = func(cmd *cobra.Command, args []string) error {
		kubeConfigPath, _ := cmd.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		namespace, _ := cmd.Flags().GetString("namespace")
		updateRepo, _ := cmd.Flags().GetBool("update-repo")
		ingressEnabled, _ := cmd.Flags().GetBool("ingress")
		ingressDomain, _ := cmd.Flags().GetString("domain")
		persistence, _ := cmd.Flags().GetBool("persistence")
		trivyEnabled, _ := cmd.Flags().GetBool("trivy")
		notaryEnabled, _ := cmd.Flags().GetBool("notary")
		customFlags, _ := cmd.Flags().GetStringArray("set")

		overrides := map[string]string{
			"expose.type":                "clusterIP",
			"externalURL":                "harbor",
			"trivy.enabled":              "false",
			"notary.enabled":             "false",
			"persistence.enabled":        "false",
			"expose.tls.auto.commonName": "harbor",
			"chartmuseum.enabled":        "false", // disable chartmuseum per https://github.com/goharbor/harbor/discussions/15057
		}

		if ingressDomain != "" {
			overrides["expose.ingress.hosts.core"] = fmt.Sprintf("harbor.%s", ingressDomain)
			overrides["externalURL"] = fmt.Sprintf("https://harbor.%s", ingressDomain)
		}

		if ingressEnabled {
			overrides["expose.type"] = "ingress"
		}

		if persistence {
			overrides["persistence.enabled"] = "true"
		}

		if trivyEnabled {
			overrides["trivy.enabled"] = "true"
		}

		if notaryEnabled {
			overrides["notary.enabled"] = "true"
			overrides["expose.ingress.hosts.notary"] = fmt.Sprintf("notary.%s", ingressDomain)
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		harborOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("harbor/harbor").
			WithHelmURL("https://helm.goharbor.io").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(harborOptions)
		if err != nil {
			return err
		}

		fmt.Println(HarborInstallMsg)

		return nil
	}

	return harbor
}

const HarborInfoMsg = `Please wait for several minutes for Harbor deployment to complete.
Then you should be able to visit the Harbor portal at https://harbor.<your-domain>/ or via kubectl port-forward:
kubectl port-forward -n harbor svc/harbor 8081:443

# use https when port-forwarding, e.g.
https://localhost:8081

To login, get the admin password by running: 
kubectl get secrets -n harbor harbor-core -o jsonpath="{.data.HARBOR_ADMIN_PASSWORD}" | base64 -d

For more details, please visit https://github.com/goharbor/harbor
`

const HarborInstallMsg = `=======================================================================
=                     Harbor has been installed.                       =
=======================================================================` +
	"\n\n" + HarborInfoMsg + "\n\n" + pkg.SupportMessageShort
