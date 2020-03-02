package bundle

import (
	"log"

	"github.com/alexellis/arkade/pkg/app"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/spf13/cobra"
)

func MakeBundleMinicloud() *cobra.Command {
	var minicloud = &cobra.Command{
		Use:          "minicloud",
		Short:        "Installs minicloud",
		Long:         `Installs minicloud`,
		Example:      `  arkade bundle minicloud --namespace minicloud`,
		SilenceUsage: true,
	}

	minicloud.Flags().StringP("namespace", "n", "minicloud", "The namespace used for installation")

	minicloud.RunE = func(command *cobra.Command, args []string) error {

		namespace, _ := command.Flags().GetString("namespace")

		appList := []*app.HelmApp{}

		metricsserverApp := apps.MakeAppMetricsServer()
		metricsserverApp.Namespace = namespace

		appList = append(appList, metricsserverApp)

		for _, app := range appList {

			// check if app is already installed app.Verify()
			err := app.Install()

			if err != nil {
				log.Println(err)
			}
		}

		return nil
	}

	return minicloud
}
