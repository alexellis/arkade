package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallKongIngress() *cobra.Command {
	var command = &cobra.Command{
		Use:          "kong-ingress",
		Short:        "Install kong-ingress for OpenFaaS",
		Long:         `Install kong-ingress for OpenFaaS`,
		Example:      `arkade install kong-ingress`,
		SilenceUsage: true,
	}

	command.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	command.Flags().Bool("update-repo", true, "Update the helm repo")

	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set key=value)")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		updateRepo, _ := command.Flags().GetBool("update-repo")

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		namespace, _ := command.Flags().GetString("namespace")

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		overrides := map[string]string{}
		// always need to be set to  false
		overrides["ingressController.installCRDs"] = "false"

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		kongIngressAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("kong/kong").
			WithHelmURL("https://charts.konghq.com/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(kongIngressAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(kongIngressInstallMsg)

		return nil
	}

	return command
}

const KongIngressInfoMsg = `# If you're using a local environment such as "minikube" or "KinD",
# then try the inlets operator with "arkade install inlets-operator"

# Find out more on the project homepage:
# https://github.com/Kong/kubernetes-ingress-controller`

const kongIngressInstallMsg = `=======================================================================
= kong-ingress has been installed.                                  =
=======================================================================` +
	"\n\n" + KongIngressInfoMsg + "\n\n" + pkg.ThanksForUsing
