// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/types"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
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

	// Change those to arguments as they are required?
	gitlabApp.Flags().StringP("domain", "d", "", "Domain name that will be used for all publicly exposed services (required)")
	gitlabApp.Flags().StringP("external-ip", "i", "", "Static IP to assign to NGINX Ingress Controller (required)")

	gitlabApp.Flags().Bool("ce", false, "Install the Community Edition of GitLab")

	// Extras.
	gitlabApp.Flags().Bool("no-pgsql", false, "Do not install PostgreSQL alongside GitLab")
	gitlabApp.Flags().Bool("no-redis", false, "Do not install Redis alongside GitLab")
	gitlabApp.Flags().Bool("no-minio", false, "Do not install MinIO alongside GitLab")

	_ = gitlabApp.MarkFlagRequired("domain")
	_ = gitlabApp.MarkFlagRequired("external-ip")

	gitlabApp.RunE = func(cmd *cobra.Command, args []string) error {
		helm3 := true

		namespace, _ := cmd.Flags().GetString("namespace")
		userPath, err := config.InitUserDir()

		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()
		log.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		overrides := map[string]string{}
		overrides["global.domain"], _ = cmd.Flags().GetString("domain")
		overrides["global.ip"], _ = cmd.Flags().GetString("external-ip")

		ceEdition, _ := cmd.Flags().GetBool("ce")
		noInstallPgsql, _ := cmd.Flags().GetBool("no-pgsql")
		noInstallRedis, _ := cmd.Flags().GetBool("no-redis")
		noInstallMinio, _ := cmd.Flags().GetBool("no-minio")

		if ceEdition {
			overrides["global.edition"] = "ce"
		}

		if !noInstallPgsql {
			overrides["global.postgresql.install"] = "false"
		}

		if !noInstallRedis {
			overrides["global.redis.install"] = "false"
		}

		if !noInstallMinio {
			overrides["global.minio.enabled"] = "false"
		}

		customFlags, _ := cmd.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		options := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("gitlab/gitlab").
			WithHelmURL("https://charts.gitlab.io").
			WithOverrides(overrides)

		if cmd.Flags().Changed("kubeconfig") {
			kubeconfigPath, _ := cmd.Flags().GetString("kubeconfig")
			options.WithKubeconfigPath(kubeconfigPath)
		}

		_ = os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(options)
		if err != nil {
			return err
		}

		println(gitlabInstallMsg)
		return nil
	}

	return gitlabApp
}

const gitlabInfoMessage = ``

const gitlabInstallMsg = `=======================================================================
= GitLab has been installed.                                          =
=======================================================================` +
	"\n\n" + gitlabInfoMessage + "\n\n" + pkg.ThanksForUsing
