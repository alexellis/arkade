package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallMongoDB() *cobra.Command {
	var command = &cobra.Command{
		Use:          "mongodb",
		Short:        "Install mongodb",
		Long:         `Install mongodb`,
		Example:      `  arkade install mongodb`,
		SilenceUsage: true,
	}
	command.Flags().String("namespace", "default", "Namespace for the app")

	command.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set=mongodbUsername=admin)")
	command.Flags().Bool("persistence", false, "Create and bind a persistent volume, not recommended for development")

	command.RunE = func(command *cobra.Command, args []string) error {

		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath := getDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		namespace, _ := command.Flags().GetString("namespace")

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		helm3 := true

		persistence, _ := command.Flags().GetBool("persistence")

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, helm3)
		if err != nil {
			return err
		}

		err = addHelmRepo("stable", "https://kubernetes-charts.storage.googleapis.com/", helm3)
		if err != nil {
			return fmt.Errorf("unable to add repo %s", err)
		}

		updateRepo, _ := command.Flags().GetBool("update-repo")

		if updateRepo {
			err = updateHelmRepos(helm3)
			if err != nil {
				return fmt.Errorf("unable to update repos %s", err)
			}
		}

		chartPath := path.Join(os.TempDir(), "charts")

		err = fetchChart(chartPath, "stable/mongodb", defaultVersion, helm3)

		if err != nil {
			return fmt.Errorf("unable fetch chart %s", err)
		}

		overrides := map[string]string{}

		overrides["persistence.enabled"] = strconv.FormatBool(persistence)

		outputPath := path.Join(chartPath, "mongodb")

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		err = helm3Upgrade(outputPath, "stable/mongodb",
			namespace, "values.yaml", defaultVersion, overrides, wait)
		if err != nil {
			return fmt.Errorf("unable to mongodb chart with helm %s", err)
		}
		fmt.Println(mongoDBPostInstallMsg)
		return nil
	}
	return command
}

const mongoDBPostInstallMsg = `=======================================================================
=                  MongoDB has been installed.                        =
=======================================================================` +
	"\n\n" + pkg.ThanksForUsing

var MongoDBInfoMsg = `
# MongoDB can be accessed via port 27017 on the following DNS name from within your cluster:

mongodb.{{namespace}}.svc.cluster.local

# To get the root password run:

export MONGODB_ROOT_PASSWORD=$(kubectl get secret --namespace {{namespace}} mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)

# To connect to your database run the following command:

kubectl run --namespace {{namespace}} mongodb-client --rm --tty -i --restart='Never' --image bitnami/mongodb --command -- mongo admin --host mongodb --authenticationDatabase admin -u root -p $MONGODB_ROOT_PASSWORD

# To connect to your database from outside the cluster execute the following commands:

kubectl port-forward --namespace {{namespace}} svc/mongodb 27017:27017 &
mongo --host 127.0.0.1 --authenticationDatabase admin -p $MONGODB_ROOT_PASSWORD

# More on GitHub : https://github.com/helm/charts/tree/master/stable/mongodb`
