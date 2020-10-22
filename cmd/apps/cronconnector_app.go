package apps

import (
	"fmt"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg/commands"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallCronConnector() *cobra.Command {
	var command = &cobra.Command{
		Use:          "cron-connector",
		Short:        "Install cron-connector for OpenFaaS",
		Long:         `Install cron-connector for OpenFaaS`,
		Example:      `  arkade install cron-connector`,
		SilenceUsage: true,
	}

	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set key=value)")

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		updateRepo, _ := command.Flags().GetBool("update-repo")

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		namespace, err := commands.GetNamespace(command.Flags(), "openfaas")

		clientArch, clientOS := env.GetClientArch()

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		overrides := map[string]string{}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		cronConnectorAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmPath(path.Join(userPath, ".helm")).
			WithHelmRepo("openfaas/cron-connector").
			WithHelmURL("https://openfaas.github.io/faas-netes/").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		err = apps.MakeInstallChart(cronConnectorAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(cronConnectorInstallMsg)

		return nil
	}

	return command
}

const CronConnectorInfoMsg = `# Example usage to trigger nodeinfo every 5 minutes:

faas-cli store deploy nodeinfo \
  --annotation schedule="*/5 * * * *" \
  --annotation topic=cron-function

# View the connector's logs:

kubectl logs deploy/cron-connector -n openfaas -f

# Find out more on the project homepage:

# https://github.com/openfaas-incubator/cron-connector/`

const cronConnectorInstallMsg = `=======================================================================
= cron-connector has been installed.                                  =
=======================================================================` +
	"\n\n" + CronConnectorInfoMsg + "\n\n" + pkg.ThanksForUsing
