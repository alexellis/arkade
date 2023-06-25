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
	"github.com/spf13/cobra"
)

func MakeInstallJenkins() *cobra.Command {
	var jenkins = &cobra.Command{
		Use:          "jenkins",
		Short:        "Install jenkins",
		Long:         `Install jenkins`,
		Example:      `  arkade install jenkins`,
		SilenceUsage: true,
	}

	jenkins.Flags().Bool("update-repo", true, "Update the helm repo")
	jenkins.Flags().String("namespace", "default", "Kubernetes namespace for the application")
	jenkins.Flags().Bool("persistence", false, "Enable persistence")
	jenkins.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	jenkins.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetBool("wait")
		if err != nil {
			return fmt.Errorf("error with --wait usage: %s", err)
		}

		_, err = command.Flags().GetBool("persistence")
		if err != nil {
			return fmt.Errorf("error with --persistence usage: %s", err)
		}

		_, err = command.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("error with --namespace usage: %s", err)
		}

		_, err = command.Flags().GetBool("update-repo")
		if err != nil {
			return fmt.Errorf("error with --update-repo usage: %s", err)
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		return nil
	}

	jenkins.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		updateRepo, _ := jenkins.Flags().GetBool("update-repo")
		ns, _ := command.Flags().GetString("namespace")
		persistence, _ := command.Flags().GetBool("persistence")
		customFlags, _ := command.Flags().GetStringArray("set")

		overrides := map[string]string{}
		overrides["persistence.enabled"] = strings.ToLower(strconv.FormatBool(persistence))

		// set custom flags
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		jenkinsAppOptions := types.DefaultInstallOptions().
			WithNamespace(ns).
			WithHelmRepo("jenkins/jenkins").
			WithHelmURL("https://charts.jenkins.io/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err := apps.MakeInstallChart(jenkinsAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(jenkinsInstallMsg)
		return nil
	}

	return jenkins
}

const JenkinsInfoMsg = `# Jenkins can take several minutes to install, check its status with:
kubectl rollout status deploy/jenkins --timeout 10m

# Get the Jenkins credentials:
export USER=$(kubectl get secret jenkins \
	-o jsonpath="{.data.jenkins-admin-user}" | base64 --decode)
export PASS=$(kubectl get secret jenkins \
	-o jsonpath="{.data.jenkins-admin-password}" | base64 --decode)

echo "Credentials: $USER / $PASS"

# Port-forward the Jenkins service
kubectl port-forward svc/jenkins 8080:8080 &

# Open the Jenkins UI at:
echo http://127.0.0.1:8080
`

var jenkinsInstallMsg = `=======================================================================
=                    Jenkins has been installed.                      =
=======================================================================` +
	"\n\n" + JenkinsInfoMsg + "\n\n" + pkg.SupportMessageShort
