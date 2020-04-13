// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
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

	jenkins.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()
		wait, _ := command.Flags().GetBool("wait")

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}
		updateRepo, _ := jenkins.Flags().GetBool("update-repo")

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		arch := getNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		ns, _ := jenkins.Flags().GetString("namespace")

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, true)
		if err != nil {
			return err
		}

		err = addHelmRepo("stable", "https://kubernetes-charts.storage.googleapis.com", true)
		if err != nil {
			return err
		}

		if updateRepo {
			err = updateHelmRepos(true)
			if err != nil {
				return err
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")
		err = fetchChart(chartPath, "stable/jenkins", defaultVersion, true)

		if err != nil {
			return err
		}

		persistence, _ := jenkins.Flags().GetBool("persistence")
		overrides := map[string]string{}

		overrides["persistence.enabled"] = strings.ToLower(strconv.FormatBool(persistence))

		customFlags, err := jenkins.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		outputPath := path.Join(chartPath, "jenkins")

		err = helm3Upgrade(outputPath, "stable/jenkins", ns,
			"values.yaml",
			defaultVersion,
			overrides,
			wait)

		if err != nil {
			return err
		}

		fmt.Println(jenkinsInstallMsg)
		return nil
	}

	return jenkins
}

var JenkinsInfoMsg = `# Jenkins can take several minutes to install, check its status with:
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
	"\n\n" + JenkinsInfoMsg + "\n\n" + pkg.ThanksForUsing
