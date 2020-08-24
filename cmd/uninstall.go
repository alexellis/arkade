// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeUninstall() *cobra.Command {
	var command = &cobra.Command{
		Use:          "uninstall",
		Short:        "Uninstall apps installed with arkade",
		Long:         `Uninstall apps installed with arkade`,
		Example:      `  arkade uninstall`,
		Aliases:      []string{"delete"},
		SilenceUsage: false,
	}

	command.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")
	command.PersistentFlags().Bool("wait", false, "If we should wait for the resource to be ready before returning (helm3 only, default false)")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			fmt.Printf(
				`Apps installed to Kubernetes can rarely be uninstalled in a single command 
and often leave clusters in an inconsistent state. Kubernetes does not 
track resources created by applications in the same way that something like 
Windows 10 or MacOS would do, so  it is often much easier for you to 
create a new cluster, than to remove an application.

With a tool like kind, you just run:

kind delete cluster; kind create cluster

Most arkade apps are installed with helm, so you can simply use helm to 
remove them, but beware that each project generally provides very specific 
guidance on how to clear up all resources that it may have created 
at installation time.

Get helm:
arkade get helm

List charts:
helm list --all-namespaces

Delete a chart:
helm delete -n openfaas openfaas

Delete any namespaces it created:
kubectl delete namespace openfaas openfaas-fn

Where an app was installed via manifest files or an upstream CLI like 
with linkerd or OSM, then you can usually delete its namespace, 
however this does not always remove all resources such as a 
ClusterRole. See the app's code for more:

https://github.com/alexellis/arkade/tree/master/cmd/apps

You can seek out technical support from the OpenFaaS community 
at https://slack.openfaas.io/

`)
			return nil
		}

		return nil
	}

	return command
}
