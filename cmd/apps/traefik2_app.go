package apps

import (
	"fmt"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

func installTraefik2(parts ...string) (execute.ExecResult, error) {

	task := execute.ExecTask{
		Command:     "helm",
		Args:        parts,
		StreamStdio: true,
	}
	res, err := task.Execute()
	if err != nil {
		return res, err
	}
	if res.ExitCode != 0 {
		return res, fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}
	return res, nil
}

func MakeInstallTraefik2() *cobra.Command {
	var traefik2 = &cobra.Command{
		Use:          "traefik2",
		Short:        "Install traefik2",
		Long:         "Install traefik2",
		Example:      `  arkade app install traefik2`,
		SilenceUsage: true,
	}

	var dashboard bool
	traefik2.Flags().StringP("namespace", "n", "kube-system", "The namespace used for installation")
	traefik2.Flags().Bool("update-repo", true, "Update the helm repo")
	traefik2.Flags().BoolVar(&dashboard, "dashboard", false, "Expose dashboard if you want access to dashboard from the browser")
	traefik2.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set key=value)")

	traefik2.RunE = func(command *cobra.Command, args []string) error {
		PrintTraefikASCIIArt()
		kubeConfigPath := getDefaultKubeconfig()
		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		updateRepo, _ := traefik2.Flags().GetBool("update-repo")
		namespace, _ := traefik2.Flags().GetString("namespace")
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		_, clientOS := env.GetClientArch()
		clientArch, clientOS := env.GetClientArch()
		fmt.Printf("Client: %q\n", clientOS)

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, false)
		if err != nil {
			return err
		}

		err = addHelmRepo("traefik", "https://containous.github.io/traefik-helm-chart", false)
		if err != nil {
			return fmt.Errorf("Unable to add repo %s", err)
		}

		if updateRepo {
			err = updateHelmRepos(false)
			if err != nil {
				return err
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")
		err = fetchChart(chartPath, "traefik/traefik", "", false)
		if err != nil {
			return fmt.Errorf("Unable fetch char %s", err)
		}

		overrides := map[string]string{}
		overrides["service.type"] = "NodePort"

		if dashboard {
			overrides["dashboard.ingressRoute"] = "true"
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		outputPath := path.Join(chartPath, "traefik/traefik/")
		err = templateChart(chartPath,
			"traefik",
			namespace,
			outputPath,
			"values.yaml",
			overrides)

		if err != nil {
			return err
		}

		err = kubectl("apply", "-R", "-f", outputPath, "-n", namespace)
		if err != nil {
			return err
		}
		fmt.Println(traefikInstallMsg)
		return nil
	}

	return traefik2
}

func PrintTraefikASCIIArt() {
	arkadeLogo := aec.BlueF.Apply(traefikstarted)
	fmt.Print(arkadeLogo)
}

const Traefik2InfoMsg = `# Get started with traefik v2 here:
https://docs.traefik.io/v2.0/
 
# Install traefik version 2 enabling dashboard access

	$ arkade install traefik2 --dashboard
`

const traefikstarted = `

_______              __ _ _           ___
|__   __|            / _(_) |         |__ \
   | |_ __ __ _  ___| |_ _| | __ __   __ ) |
   | | '__/ _´ |/ _ \  _| | |/ / \ \ / // /
   | | | | (_| |  __/ | | |   <   \ V // /_
   |_|_|  \__,_|\___|_| |_|_|\_\   \_/|____|

`
const traefikInstallMsg = `=======================================================================
=                  traefik v2 has been installed.                        =
=======================================================================
` +
	"\n" + pkg.ThanksForUsing +
	`
NOTES:

1. Traefik has been started. You can find out the port numbers being used by traefik by running:

          $ kubectl describe svc traefik --namespace kube-system

2. To Access to Dashboard you can run kubectl proxy or expose the service, for example:

          $ kubectl port-forward service/traefik 8181:80 -n kube-system 

3. Access through https://localhost:8181/dashboard/

`
