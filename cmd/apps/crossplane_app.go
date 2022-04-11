// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/helm"
	"github.com/spf13/cobra"
)

func MakeInstallCrossplane() *cobra.Command {
	var crossplane = &cobra.Command{
		Use:   "crossplane",
		Short: "Install Crossplane",
		Long: `Install Crossplane to deploy managed services across cloud providers and
schedule workloads to any Kubernetes cluster`,
		Example:      `  arkade install crossplane`,
		SilenceUsage: true,
	}

	crossplane.Flags().StringP("namespace", "n", "crossplane-system", "The namespace used for installation")
	crossplane.Flags().Bool("update-repo", true, "Update the helm repo")
	crossplane.Flags().StringArray("set", []string{},
		"Use custom flags or override existing flags \n(example --set persistence.enabled=true)")

	crossplane.RunE = func(command *cobra.Command, args []string) error {
		wait, _ := command.Flags().GetBool("wait")
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		namespace, _ := command.Flags().GetString("namespace")

		if namespace != "crossplane-system" {
			return fmt.Errorf(`to override the namespace, install crossplane via helm manually`)
		}

		arch := k8s.GetNodeArchitecture()
		if !strings.Contains(arch, "64") {
			return fmt.Errorf(`crossplane is currently only supported on 64-bit architectures`)
		}
		fmt.Printf("Node architecture: %q\n", arch)

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}

		overrideValues, err := crossplane.Flags().GetStringArray("set")
		if err != nil {
			return err
		}
		values := make(map[string]string)
		if err := mergeFlags(values, overrideValues); err != nil {
			return err
		}

		clientArch, clientOS := env.GetClientArch()

		fmt.Printf("Client: %q, %q\n", clientArch, clientOS)

		log.Printf("User dir established as: %s\n", userPath)

		os.Setenv("HELM_HOME", path.Join(userPath, ".helm"))

		_, err = helm.TryDownloadHelm(userPath, clientArch, clientOS)
		if err != nil {
			return err
		}

		updateRepo, _ := crossplane.Flags().GetBool("update-repo")
		err = helm.AddHelmRepo("crossplane-alpha", "https://charts.crossplane.io/alpha", updateRepo)
		if err != nil {
			return err
		}

		err = helm.FetchChart("crossplane-alpha/crossplane", defaultVersion)
		if err != nil {
			return err
		}

		_, nsErr := k8s.KubectlTask("create", "namespace", "crossplane-system")
		if nsErr != nil && !strings.Contains(nsErr.Error(), "AlreadyExists") {
			return nsErr
		}

		err = helm.Helm3Upgrade("crossplane-alpha/crossplane",
			namespace, "values.yaml", "", values, wait)
		if err != nil {
			return err
		}

		fmt.Println(crossplaneInstallMsg)
		return nil
	}

	return crossplane
}

const CrossplaneInfoMsg = `# Get started by installing a stack for your favorite provider:
* provider-gcp: https://crossplane.io/docs/master/install-crossplane.html#gcp-provider
* provider-aws: https://crossplane.io/docs/master/install-crossplane.html#aws-provider
* provider-azure: https://crossplane.io/docs/master/install-crossplane.html#azure-provider

Learn more about Crossplane: https://crossplaneio.github.io/docs/`

const crossplaneInstallMsg = `=======================================================================
= Crossplane has been installed.                                      =
=======================================================================` +
	"\n\n" + CrossplaneInfoMsg + "\n\n" + pkg.SupportMessageShort
