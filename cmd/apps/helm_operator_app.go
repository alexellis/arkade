package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallHelmOperator() *cobra.Command {
	var operator = &cobra.Command{
		Use:          "helm-operator",
		Short:        "Install the helm operator",
		Long:         `Install the helm operator`,
		Example:      `arkade install helm-operator --namespace default`,
		SilenceUsage: true,
	}

	operator.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	operator.Flags().Bool("update-repo", true, "Update the helm repo")
	operator.Flags().Bool("helm3", true, "Use helm3, if set to false uses helm2")

	operator.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()
		wait, _ := command.Flags().GetBool("wait")

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		updateRepo, _ := operator.Flags().GetBool("update-repo")

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)
		helm3, _ := command.Flags().GetBool("helm3")

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}
		namespace, _ := command.Flags().GetString("namespace")
		if namespace != "default" {
			return fmt.Errorf(`to override the "default", install via tiller`)
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		err = addHelmRepo("fluxcd", "https://charts.fluxcd.io", helm3)
		if err != nil {
			return err
		}

		if updateRepo {
			err = updateHelmRepos(helm3)
			if err != nil {
				return err
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")
		err = fetchChart(chartPath, "fluxcd/helm-operator", "0.7.0", helm3)

		if err != nil {
			return err
		}

		overrides := map[string]string{}
		// the helm chart always installs the crds
		//overrides["createCRD"] = "true"
		overrides["helm.versions"] = "v3 "

		arch := getNodeArchitecture()

		fmt.Printf("Node architecture: %q\n", arch)

		if arch != "amd64" {
			return fmt.Errorf("This chart does not support %s", arch)
		}

		fmt.Println("Chart path: ", chartPath)

		ns := "default"

		log.Printf("Applying CRD: CRDs will be applied by the helm chart\n")

		if helm3 {
			outputPath := path.Join(chartPath, "helm-operator")

			err := helm3Upgrade(outputPath, "fluxcd/helm-operator", ns,
				"values.yaml",
				"",
				overrides,
				wait,
			)

			if err != nil {
				return err
			}
		} else {
			outputPath := path.Join(chartPath, "helm-operator/rendered")

			err = templateChart(chartPath,
				"helm-operator",
				ns,
				outputPath,
				"values.yaml",
				overrides)

			if err != nil {
				return err
			}

			err = kubectl("apply", "-R", "-f", outputPath)

			if err != nil {
				return err
			}
		}

		fmt.Println(helmOperatorInstallMsg)
		return nil
	}

	return operator
}

const helmOperatorInstallMsg = `# The helm-operator has been configured

For example you can install kubernetes nginx-ingress this way:

kubectl apply -f - <<EOF
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  name: nginx-ingress
  namespace: kube-system
spec:
  releaseName: nginx-ingress
  targetNamespace: kube-system
  timeout: 300
  resetValues: false
  wait: false
  forceUpgrade: false
  chart:
    repository: https://kubernetes-charts.storage.googleapis.com
    name: nginx-ingress
    version: 1.30.0
EOF

After some time you will see the helm-release:

kubectl get helmreleases

or 

helm list -A

# Find out more at:
# https://docs.fluxcd.io/projects/helm-operator/en/latest/references/helmrelease-custom-resource.html`
