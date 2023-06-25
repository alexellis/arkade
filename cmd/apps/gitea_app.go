// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"
	"strconv"
	"strings"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
)

func MakeInstallGitea() *cobra.Command {
	var gitea = &cobra.Command{
		Use:          "gitea",
		Short:        "Install gitea",
		Long:         `Install gitea`,
		Example:      `  arkade install gitea`,
		SilenceUsage: true,
	}

	gitea.Flags().Bool("update-repo", true, "Update the helm repo")
	gitea.Flags().String("namespace", "default", "Kubernetes namespace for the application")
	gitea.Flags().Bool("persistence", false, "Enable persistence")
	gitea.Flags().StringP("user", "u", "gitea_admin", "Username of admin user")
	gitea.Flags().StringP("password", "p", "", "Overide the default random admin-password if this is set")
	gitea.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	gitea.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		updateRepo, _ := gitea.Flags().GetBool("update-repo")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		ns, _ := gitea.Flags().GetString("namespace")

		persistence, _ := gitea.Flags().GetBool("persistence")
		overrides := map[string]string{}

		overrides["persistence.enabled"] = strings.ToLower(strconv.FormatBool(persistence))

		pass, _ := command.Flags().GetString("password")

		if len(pass) == 0 {
			var err error
			pass, err = password.Generate(25, 10, 0, false, true)
			if err != nil {
				return err
			}
		}

		overrides["gitea.admin.password"] = pass

		adminUsername, err := command.Flags().GetString("user")
		if err != nil {
			return err
		}
		overrides["gitea.admin.username"] = adminUsername

		// disabling password complexity check by default per NIST guidelines
		// as it has been disabled upstream just waiting on a stable release
		// user is free to enable password checks using --set flag
		overrides["gitea.config.security.PASSWORD_COMPLEXITY"] = "off"

		customFlags, err := gitea.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		giteaAppOptions := types.DefaultInstallOptions().
			WithNamespace(ns).
			WithHelmRepo("gitea-charts/gitea").
			WithHelmURL("https://dl.gitea.io/charts").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(giteaAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(giteaInstallMsg)
		return nil
	}

	return gitea
}

const GiteaInfoMsg = `# Forward the gateway to your machine
kubectl rollout status --namespace {{ namespace }} sts/gitea
kubectl --namespace {{ namespace }} port-forward svc/gitea-http 3000:3000

# Open up http://127.0.0.1:3000 to use your application

# Find out more at:
# https://gitea.com/gitea/helm-chart`

var giteaInstallMsg = `=======================================================================
=                    Gitea has been installed.                      =
=======================================================================` +
	"\n\n" + GiteaInfoMsg + "\n\n" + pkg.SupportMessageShort
