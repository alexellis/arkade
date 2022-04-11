// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
)

func MakeInstallGitLab() *cobra.Command {
	var gitlabApp = &cobra.Command{
		Use:          "gitlab",
		Short:        "Install GitLab",
		Long:         "Install GitLab, the complete DevOps platform",
		Example:      "arkade install gitlab",
		SilenceUsage: true,
	}

	gitlabApp.Flags().StringP("namespace", "n", "default", "The namespace to install GitLab")
	gitlabApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set global.hosts.domain)")
	gitlabApp.Flags().Bool("update-repo", true, "Update the helm repo")
	// Maybe these should be arguments as they are required?
	gitlabApp.Flags().StringP("domain", "d", "", "Domain name that will be used for all publicly exposed services (required)")
	gitlabApp.Flags().StringP("external-ip", "i", "", "Static IP to assign to NGINX Ingress Controller (required)")
	// EE is the default type to install.
	gitlabApp.Flags().Bool("ce", false, "Install the Community Edition of GitLab")
	// The following dependencies are optional, as external instances of those can be used.
	// This does though require that some values are changed in the helm values via the `--set` flag.
	gitlabApp.Flags().Bool("no-pgsql", false, "Do not install PostgreSQL alongside GitLab")
	gitlabApp.Flags().Bool("no-redis", false, "Do not install Redis alongside GitLab")
	gitlabApp.Flags().Bool("no-minio", false, "Do not install MinIO alongside GitLab")

	_ = gitlabApp.MarkFlagRequired("domain")
	_ = gitlabApp.MarkFlagRequired("external-ip")

	gitlabApp.RunE = func(cmd *cobra.Command, args []string) error {
		namespace, _ := cmd.Flags().GetString("namespace")
		kubeConfigPath, _ := cmd.Flags().GetString("kubeconfig")

		overrides := map[string]string{}
		overrides["global.hosts.domain"], _ = cmd.Flags().GetString("domain")
		overrides["global.hosts.externalIP"], _ = cmd.Flags().GetString("external-ip")

		ceEdition, _ := cmd.Flags().GetBool("ce")
		noInstallPgsql, _ := cmd.Flags().GetBool("no-pgsql")
		noInstallRedis, _ := cmd.Flags().GetBool("no-redis")
		noInstallMinio, _ := cmd.Flags().GetBool("no-minio")

		if ceEdition {
			overrides["global.edition"] = "ce"
		}

		if noInstallPgsql {
			overrides["global.postgresql.install"] = "false"
		}

		if noInstallRedis {
			overrides["global.redis.install"] = "false"
		}

		if noInstallMinio {
			overrides["global.minio.enabled"] = "false"
		}

		customFlags, _ := cmd.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		options := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("gitlab/gitlab").
			WithHelmURL("https://charts.gitlab.io").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(options)
		if err != nil {
			return err
		}

		println(gitlabInstallMsg)
		return nil
	}

	return gitlabApp
}

const GitlabInfoMsg = `# For full configuration information, make sure to check out
# https://docs.gitlab.com/charts/charts/globals.html

# To access your GitLab installation, visit the external URI after all the pods are running.
# You can get the initial root administration password with the following command:

kubectl get secret gitlab-gitlab-initial-root-password -o jsonpath='{.data.password}' | base64 --decode ; echo

`

const gitlabInstallMsg = `=======================================================================
= GitLab has been installed.                                          =
=======================================================================` +
	"\n\n" + GitlabInfoMsg + "\n\n" + pkg.SupportMessageShort
