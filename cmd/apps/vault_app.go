// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallVault() *cobra.Command {
	var vault = &cobra.Command{
		Use:   "vault",
		Short: "Install vault",
		Long:  `Install vault`,
		Example: `  arkade install vault
  # with ingress and custom domain
  arkade install vault --ingress=true --domain=vault.example.com`,
		SilenceUsage: true,
	}

	vault.Flags().StringP("namespace", "n", "vault", "The namespace to install the chart")
	vault.Flags().Bool("update-repo", true, "Update the helm repo")
	vault.Flags().Bool("ingress", false, "Enable ingress")
	vault.Flags().String("domain", "", "Set ingress domain")
	vault.Flags().Bool("injector", false, "Enable sidecar injector")
	vault.Flags().Bool("persistence", false, "Enable vault server persistence")
	vault.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image.tag=1.11.2)")

	vault.PreRunE = func(cmd *cobra.Command, args []string) error {
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

		ingressDomain, err := cmd.Flags().GetString("domain")
		if err != nil {
			return err
		}

		if !ingressEnabled && ingressDomain != "" {
			return fmt.Errorf("--domain option should be used only with --ingress=true")
		}

		_, err = cmd.Flags().GetBool("injector")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetBool("persistence")
		if err != nil {
			return err
		}

		return nil
	}

	vault.RunE = func(cmd *cobra.Command, args []string) error {
		kubeConfigPath, _ := cmd.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		namespace, _ := cmd.Flags().GetString("namespace")
		updateRepo, _ := cmd.Flags().GetBool("update-repo")
		ingressEnabled, _ := cmd.Flags().GetBool("ingress")
		ingressDomain, _ := cmd.Flags().GetString("domain")
		injectorEnabled, _ := cmd.Flags().GetBool("injector")
		persistence, _ := cmd.Flags().GetBool("persistence")
		customFlags, _ := cmd.Flags().GetStringArray("set")

		overrides := map[string]string{
			"injector.enabled":           "false",
			"server.dataStorage.enabled": "false",
		}

		if ingressDomain != "" {
			overrides["server.ingress.hosts[0].host"] = ingressDomain
		}

		if ingressEnabled {
			overrides["server.ingress.enabled"] = "true"
		}
		if injectorEnabled {
			overrides["injector.enabled"] = "true"
		}
		if persistence {
			overrides["server.dataStorage.enabled"] = "true"
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		vaultOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("hashicorp/vault").
			WithHelmURL("https://helm.releases.hashicorp.com").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(vaultOptions)
		if err != nil {
			return err
		}

		fmt.Println(vaultInstallMsg)

		return nil
	}

	return vault
}

const VaultInfoMsg = `# To use Vault, you need to initialize it:
kubectl exec --stdin=true --tty=true -n vault vault-0 -- vault operator init

# Then save Unseal Keys and root token and execute command below 3 times with each key (e.g. Unseal Key 1, 2, 3)
kubectl exec --stdin=true --tty=true -n vault vault-0 -- vault operator unseal # ... Unseal Key 1

# Sealed should be false after that
Key             Value
---             -----
Seal Type       shamir
Initialized     true
Sealed          false
Total Shares    5
Threshold       3
Version         1.11.2
Build Date      2022-07-29T09:48:47Z
Storage Type    file
Cluster Name    vault-cluster-dad545cc
Cluster ID      eacce4ce-7954-62a7-78ff-be5e3f4ee2f7
HA Enabled      false

# Get started with Vault at
# https://www.vaultproject.io/docs/platform/k8s

# You can install the "vault" CLI via:

arkade get vault`

const vaultInstallMsg = `=======================================================================
=                     Vault has been installed.                       =
=======================================================================` +
	"\n\n" + VaultInfoMsg + "\n\n" + pkg.SupportMessageShort
