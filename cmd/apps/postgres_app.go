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

func MakeInstallPostgresql() *cobra.Command {
	var postgresql = &cobra.Command{
		Use:          "postgresql",
		Short:        "Install postgresql",
		Long:         `Install postgresql`,
		Example:      `  arkade install postgresql`,
		SilenceUsage: true,
	}

	postgresql.Flags().Bool("update-repo", true, "Update the helm repo")
	postgresql.Flags().String("namespace", "default", "Kubernetes namespace for the application")

	postgresql.Flags().Bool("persistence", false, "Enable persistence")

	postgresql.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	postgresql.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		updateRepo, _ := postgresql.Flags().GetBool("update-repo")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		ns, _ := postgresql.Flags().GetString("namespace")

		if ns != "default" {
			return fmt.Errorf("please use the helm chart if you'd like to change the namespace to %s", ns)
		}

		persistence, _ := postgresql.Flags().GetBool("persistence")

		overrides := map[string]string{}

		overrides["persistence.enabled"] = strings.ToLower(strconv.FormatBool(persistence))

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		postgresqlAppOptions := types.DefaultInstallOptions().
			WithNamespace(ns).
			WithHelmRepo("bitnami/postgresql").
			WithHelmURL("https://charts.bitnami.com/bitnami").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(postgresqlAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(postgresqlInstallMsg)
		return nil
	}

	return postgresql
}

const PostgresqlInfoMsg = `PostgreSQL can be accessed via port 5432 on the following DNS name from within your cluster:

	postgresql.default.svc.cluster.local - Read/Write connection

To get the password for "postgres" run:

    export POSTGRES_PASSWORD=$(kubectl get secret --namespace default postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode)

To connect to your database run the following command:

    kubectl run postgresql-client --rm --tty -i --restart='Never' --namespace default --image docker.io/bitnami/postgresql:11.6.0-debian-9-r0 --env="PGPASSWORD=$POSTGRES_PASSWORD" --command -- psql --host postgresql -U postgres -d postgres -p 5432

To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace default svc/postgresql 5432:5432 &
	PGPASSWORD="$POSTGRES_PASSWORD" psql --host 127.0.0.1 -U postgres -d postgres -p 5432

# Find out more at: https://github.com/bitnami/charts/tree/main/bitnami/postgresql`

const postgresqlInstallMsg = `=======================================================================
= PostgreSQL has been installed.                                      =
=======================================================================` +
	"\n\n" + PostgresqlInfoMsg + "\n\n" + pkg.SupportMessageShort
