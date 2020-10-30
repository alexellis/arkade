package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallNginxIncIngress() *cobra.Command {
	var command = &cobra.Command{
		Use:          "nginx-inc",
		Short:        "Install nginx-inc for OpenFaaS",
		Long:         `Install nginx-inc for OpenFaaS`,
		Example:      `arkade install nginx-inc`,
		SilenceUsage: true,
	}

	command.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	command.Flags().Bool("update-repo", true, "Update the helm repo")
	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set key=value)")
	command.Flags().String("controller-image", "", "Controller Image for NGINX Plus: (assuming you have pushed the Ingress controller image nginx-ingress to your private registry myregistry.example.com)")
	command.Flags().Bool("plus", false, "Install Nginx Plus")
	command.Flags().Int("replicas", 1, "The number of replicas of the Ingress controller deployment")
	command.Flags().Bool("prometheus", false, "Install Nginx Plus")
	command.Flags().Int64("prometheus-port", 9113, "Configures the prometheus port to scrape the metrics.")
	command.Flags().String("kind", "deployment", "The kind of the Ingress controller installation - deployment or daemonset")

	command.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")

		updateRepo, _ := command.Flags().GetBool("update-repo")

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

		nginxPlus, err := command.Flags().GetBool("plus")
		if err != nil {
			return fmt.Errorf("error with --plus usage: %s", err)
		}
		overrides["controller.nginxplus"] = strconv.FormatBool(nginxPlus)

		if command.Flags().Changed("controller-image") {
			controllerImage, err := command.Flags().GetString("controller-image")
			if err != nil {
				return fmt.Errorf("error with --controller-image usage: %s", err)
			}
			overrides["controller.image.repository"] = controllerImage
		}

		replicaCount, err := command.Flags().GetInt("replicas")
		if err != nil {
			return fmt.Errorf("error with --replicas usage: %s", err)
		}
		overrides["controller.replicaCount"] = fmt.Sprintf("%d", replicaCount)

		prometheusCreate, err := command.Flags().GetBool("prometheus")
		if err != nil {
			return fmt.Errorf("error with --prometheus usage: %s", err)
		}
		overrides["prometheus.create"] = strconv.FormatBool(prometheusCreate)
		if prometheusCreate {
			prometheusPort, err := command.Flags().GetInt64("prometheus-port")
			if err != nil {
				return fmt.Errorf("error with --prometheus-port usage: %s", err)
			}
			overrides["prometheus.port"] = fmt.Sprintf("%d", prometheusPort)
		}

		kindName, err := command.Flags().GetString("kind")
		if err != nil {
			return fmt.Errorf("error with --kind usage: %s", err)
		}
		overrides["controller.kind"] = kindName

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		nginxIncIngressAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("nginx-stable/nginx-ingress").
			WithHelmURL("https://helm.nginx.com/stable").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithWait(wait)

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
			nginxIncIngressAppOptions.WithKubeconfigPath(kubeConfigPath)
		}

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		_, err = apps.MakeInstallChart(nginxIncIngressAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(NginxIncIngressInstallMsg)

		return nil
	}

	return command
}

const NginxIncIngressInfoMsg = `# If you're using a local environment such as "minikube" or "KinD",
# then try the inlets operator with "arkade install inlets-operator"

# Find out more on the project homepage:
# https://github.com/nginxinc/kubernetes-ingress`

const NginxIncIngressInstallMsg = `=======================================================================
= nginx-inc has been installed.                                  =
=======================================================================` +
	"\n\n" + NginxIncIngressInfoMsg + "\n\n" + pkg.ThanksForUsing
