// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import "github.com/spf13/cobra"

func MakeInstallGitLab() * cobra.Command {
	var gitlabApp = &cobra.Command{
		Use: "gitlab",
		Short: "Install GitLab",
		Long: "Install GitLab, a complete DevOps platform",
		Example: "arkade install gitlab",
		SilenceUsage: true,
	}

	gitlabApp.Flags().StringP("namespace", "n", "default", "The namespace to install GitLab (default: default)")
	gitlabApp.Flags().Bool("persistence", false, "Use a Persistent Volume to store data")
	gitlabApp.Flags().Int("persistence-size", 10, "Set size of persistent storage in Gi (default: 10Gi)")
	gitlabApp.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set global.hosts.domain)")
	gitlabApp.Flags().Bool("update-repo", true, "Update the helm repo")



	return gitlabApp
}