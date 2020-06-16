// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/spf13/cobra"
)

func MakeInstallRedis() *cobra.Command {
	var redis = &cobra.Command{
		Use:          "redis",
		Short:        "Install redis",
		Long:         "Install redis",
		Example:      "arkade install redis",
		SilenceUsage: true,
	}

	redis.Flags().StringP("namespace", "n", "redis", "The namespace to install redis")
	redis.Flags().Bool("update-repo", true, "Update the helm repo")

	redis.RunE = func(command *cobra.Command, args []string) error {

		const chartVersion = "10.5.7"

		// Get all flags
		wait, _ := command.Flags().GetBool("wait")
		namespace, _ := command.Flags().GetString("namespace")
		updateRepo, _ := redis.Flags().GetBool("update-repo")

		// exit on arm
		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(`only Intel, i.e. PC architecture is supported for this app`)
		}

		// initialize client env
		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		log.Printf("Client: %s, %s\n", clientArch, clientOS)
		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS, true)
		if err != nil {
			return err
		}

		// Add helm repo
		err = helm.AddHelmRepo("bitnami-redis", "https://charts.bitnami.com/bitnami", true, true)
		if err != nil {
			return fmt.Errorf("unable to add repo %s", err)
		}

		// Update chart repositories
		if updateRepo {
			err = helm.UpdateHelmRepos(true)
			if err != nil {
				return err
			}
		}

		// create the namespace
		nsRes, nsErr := k8s.KubectlTask("create", "namespace", namespace)
		if nsErr != nil {
			return nsErr
		}

		// ignore errors
		if nsRes.ExitCode != 0 {
			log.Printf("[Warning] unable to create namespace %s, may already exist: %s", namespace, nsRes.Stderr)
		}

		// download the chart
		err = helm.FetchChart("bitnami-redis/redis", chartVersion, true)
		if err != nil {
			return err
		}

		// define the values to override
		// due the missing arm support. datasource and dashboard sidecars are not possible
		overrides := map[string]string{
			"serviceAccount.create": "true",
			"rbac.create":           "true",
		}

		// install the chart
		err = helm.Helm3Upgrade("bitnami-redis/redis", namespace,
			"values.yaml",
			chartVersion,
			overrides,
			wait)

		if err != nil {
			return err
		}

		fmt.Println(redisInstallMsg)

		return nil
	}

	return redis
}

const redisInstallMsg = `=======================================================================
=                       redis has been installed                      =
=======================================================================
` + RedisInfoMsg + pkg.ThanksForUsing

const RedisInfoMsg = `
# Redis can be accessed via port 6379 on the following DNS names from within your cluster:

# redis-master.redis.svc.cluster.local for read/write operations
# redis-slave.redis.svc.cluster.local for read-only operations


# To get your password run:

  export REDIS_PASSWORD=$(kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 --decode)

# To connect to your Redis server:

# 1. Run a Redis pod that you can use as a client:

  kubectl run --namespace redis redis-client --rm --tty -i --restart='Never' \
   --env REDIS_PASSWORD=$REDIS_PASSWORD \
   --image docker.io/bitnami/redis:5.0.7-debian-10-r48 -- bash

# 2. Connect using the Redis CLI:
  redis-cli -h redis-master -a $REDIS_PASSWORD
  redis-cli -h redis-slave -a $REDIS_PASSWORD

# To connect to your database from outside the cluster execute the following commands:

  kubectl port-forward --namespace redis svc/redis-master 6379:6379 &
  redis-cli -h 127.0.0.1 -p 6379 -a $REDIS_PASSWORD

`
