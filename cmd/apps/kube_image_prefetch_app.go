// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"

	"github.com/spf13/cobra"
)

func MakeInstallKubeImagePrefetch() *cobra.Command {
	var command = &cobra.Command{
		Use:          "kube-image-prefetch",
		Short:        "Install kube-image-prefetch",
		Long:         `Install kube-image-prefetch`,
		Example:      `  arkade install kube-image-prefetch`,
		SilenceUsage: true,
	}

	command.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := config.GetDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}

		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

		_, err := k8s.KubectlTask("apply", "-f",
			"https://raw.githubusercontent.com/AverageMarcus/kube-image-prefetch/master/manifest.yaml")
		if err != nil {
			return err
		}

		fmt.Println(`=======================================================================
= kube-image-prefetch has been installed.                             =
=======================================================================` +
			"\n\n" + KubeImagePrefetchInfoMsg + "\n\n" + pkg.ThanksForUsing)

		return nil
	}

	return command
}

const KubeImagePrefetchInfoMsg = `
Pre-pulls all images, on all nodes.

To ignore deployments from having their images prefetched add the following annotation: kube-image-prefetch/ignore: "true"

To specify specific containers within a deployment to ignore when prefetching add the following annotation to the relevant deployments: kube-image-prefetch/ignore-containers: "container-name"

# Find out more at
# https://github.com/AverageMarcus/kube-image-prefetch`
