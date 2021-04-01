// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallOPAGateKeeper() *cobra.Command {
	var opaGatekeeper = &cobra.Command{
		Use:   "opa-gatekeeper",
		Short: "Install Open Policy Agent (OPA) Gatekeeper",
		Long: `Install Open Policy Agent's Gatekeeper which brings Policy and Governance
for the Kubernetes API.`,
		Aliases: []string{"gatekeeper"},
		Example: `  arkade install opa-gatekeeper
  arkade install opa-gatekeeper --set helmKey=value
  arkade install gatekeeper`,
		SilenceUsage: true,
	}

	opaGatekeeper.Flags().Bool("update-repo", true, "Update the helm repo")
	// gatekeeper chart default namespace is fixed in helm chart
	opaGatekeeper.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set auditInterval=60)")

	opaGatekeeper.PreRunE = func(command *cobra.Command, args []string) error {
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

	opaGatekeeper.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		wait, _ := command.Flags().GetBool("wait")
		updateRepo, _ := command.Flags().GetBool("update-repo")
		customFlags, _ := command.Flags().GetStringArray("set")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		overrides := map[string]string{}
		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		opaGatekeeperAppOptions := types.DefaultInstallOptions().
			WithHelmRepo("gatekeeper/gatekeeper").
			WithHelmURL("https://open-policy-agent.github.io/gatekeeper/charts").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath).
			WithWait(wait)

		_, err := apps.MakeInstallChart(opaGatekeeperAppOptions)
		if err != nil {
			return err
		}

		fmt.Println(opaGatekeeperInstallMsg)
		return nil
	}

	return opaGatekeeper
}

var OPAGatekeeperInfoMsg = `
Open Policy Agent (OPA) Gatekeeper

# Find out more at:
https://github.com/open-policy-agent/gatekeeper

# Learn about how to use gatekeeper
https://github.com/open-policy-agent/gatekeeper#how-to-use-gatekeeper

# Uninstall Gatekeeper
https://github.com/open-policy-agent/gatekeeper#how-to-use-gatekeeper
`

var opaGatekeeperInstallMsg = `=======================================================================
= Open Policy Agent Gatekeeper has been installed.                                           =
=======================================================================` +
	"\n\n" + OPAGatekeeperInfoMsg + "\n\n" + pkg.ThanksForUsing
