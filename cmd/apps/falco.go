// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"github.com/alexellis/arkade/pkg/config"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallFalco() *cobra.Command {
	var falco = &cobra.Command{
		Use:          "falco",
		Short:        "Install Falco",
		Long:         `Install Falco which brings container runtime security`,
		Example:      `arkade install falco --set helmKey=value`,
		SilenceUsage: true,
	}

	falco.Flags().Bool("update-repo", true, "Update the helm repo")
	falco.Flags().Bool("sidekick", false, "install falcosidekick dependency")
	falco.Flags().Bool("webui", false, "enable falcosidekick web UI")
	falco.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set ebpf.enabled=true)")

	falco.PreRunE = func(command *cobra.Command, args []string) error {
		_, err := command.Flags().GetString("kubeconfig")
		if err != nil {
			return fmt.Errorf("error with --kubeconfig usage: %s", err)
		}

		_, err = command.Flags().GetBool("wait")
		if err != nil {
			return fmt.Errorf("error with --wait usage: %s", err)
		}

		_, err = command.Flags().GetBool("update-repo")
		if err != nil {
			return fmt.Errorf("error with --update-repo usage: %s", err)
		}

		_, err = command.Flags().GetStringArray("set")
		if err != nil {
			return fmt.Errorf("error with --set usage: %s", err)
		}

		return nil
	}

	falco.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		wait, _ := command.Flags().GetBool("wait")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		customFlags, _ := command.Flags().GetStringArray("set")
		sidekickEnabled, _ := command.Flags().GetBool("sidekick")
		webUIEnabled, _ := command.Flags().GetBool("webui")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		overrides := map[string]string{}
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		if sidekickEnabled {
			overrides["falcosidekick.enabled"] = "true"
		}

		if webUIEnabled {
			overrides["falcosidekick.webui.enabled"] = "true"
		}

		falcoAppOptions := types.DefaultInstallOptions().
			WithHelmRepo("falcosecurity/falco").
			WithHelmURL("https://falcosecurity.github.io/charts").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err := apps.MakeInstallChart(falcoAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(falcoInstallMsg)
		return nil
	}

	return falco
}

var FalcoInfoMsg = `
Falco is a Cloud Native Runtime Security tool designed to detect anomalous activity in your applications. 
You can use Falco to monitor runtime security of your Kubernetes applications and internal components.

# Find out more at:
https://github.com/falcosecurity/falco

# Learn about how to use Falco
https://falco.org/docs/
`

var falcoInstallMsg = `=======================================================================
= Falco has been installed.                                           =
=======================================================================` +
	"\n\n" + FalcoInfoMsg + "\n\n" + pkg.SupportMessageShort
