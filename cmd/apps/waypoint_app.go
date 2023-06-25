// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallWaypoint() *cobra.Command {
	var waypoint = &cobra.Command{
		Use:          "waypoint",
		Short:        "Install Waypoint",
		Long:         `Install Waypoint to any Kubernetes cluster`,
		Example:      `  arkade install waypoint`,
		SilenceUsage: true,
	}

	waypoint.Flags().StringP("namespace", "n", "waypoint", "The namespace used for installation")
	waypoint.Flags().Bool("update-repo", true, "Update the helm repo")
	waypoint.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	waypoint.RunE = func(command *cobra.Command, args []string) error {
		updateRepo, _ := waypoint.Flags().GetBool("update-repo")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		namespace, _ := waypoint.Flags().GetString("namespace")

		overrides := map[string]string{}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		waypointOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("hashicorp/waypoint").
			WithHelmURL("https://helm.releases.hashicorp.com").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(waypointOptions)
		if err != nil {
			return err
		}

		fmt.Println(waypointInstallMsg)
		return nil
	}

	return waypoint
}

const WaypointInfoMsg = `# Find out more at:
# https://learn.hashicorp.com/collections/waypoint/get-started-kubernetes`

const waypointInstallMsg = `=======================================================================
= Waypoint has been installed.                                        =
=======================================================================` +
	"\n\n" + WaypointInfoMsg + "\n\n" + pkg.SupportMessageShort
