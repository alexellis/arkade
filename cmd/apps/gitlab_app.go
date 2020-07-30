// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/spf13/cobra"
	"log"
)

func MakeInstallGitLab() * cobra.Command {
	var gitlabApp = &cobra.Command{
		Use: "gitlab",
		Short: "Install GitLab",
		Long: "Install GitLab, a complete DevOps platform",
		Example: "arkade install gitlab",
		SilenceUsage: true,
	}

	gitlabApp.Flags().StringP("namespace", "n", "default", "The namespace to install GitLab (default: default)")
	gitlabApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set global.hosts.domain)")
	gitlabApp.Flags().Bool("update-repo", true, "Update the helm repo")

	gitlabApp.Flags().Bool("ce", false, "Install the Community Edition of GitLab (default: false)")

	// Extras.
	gitlabApp.Flags().Bool("pgsql", true, "Install PostgreSQL alongside GitLab (default: true)")
	gitlabApp.Flags().Bool("redis", true, "Install Redis alongside GitLab (default: true)")
	gitlabApp.Flags().Bool("minio", true, "Install MinIO alongside GitLab (default: true)")

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

		overrides := map[string]string {}

		//ceEdition, _ := cmd.Flags().GetBool("ce")
		installPgsql, _ := cmd.Flags().GetBool("pgsql")
		installRedis, _ := cmd.Flags().GetBool("redis")
		installMinio, _ := cmd.Flags().GetBool("minio")

		if !installPgsql {
			overrides["postgresql.install"] = "false"
		}

		if !installRedis {
			overrides["redis.install"] = "false"
		}

		if !installMinio {
			overrides["minio.enabled"] = "false"
		}


		return nil
	}


	return gitlabApp
}