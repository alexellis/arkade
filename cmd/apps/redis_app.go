// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"
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
	redis.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	redis.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		namespace, _ := command.Flags().GetString("namespace")
		wait, _ := command.Flags().GetBool("wait")
		updateRepo, _ := command.Flags().GetBool("update-repo")

		// exit on arm
		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		if arch != IntelArch {
			return fmt.Errorf(OnlyIntelArch)
		}

		overrides := map[string]string{
			"serviceAccount.create": "true",
			"rbac.create":           "true",
		}

		customFlags, err := command.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		if err = config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		redisAppOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("bitnami-redis/redis").
			WithHelmURL("https://charts.bitnami.com/bitnami").
			WithOverrides(overrides).
			WithWait(wait).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(redisAppOptions)
		if err != nil {
			return err
		}

		println(redisInstallMsg)
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
