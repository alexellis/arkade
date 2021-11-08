// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/k8s"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallInletsTcpClient() *cobra.Command {
	var inletsProClient = &cobra.Command{
		Use:          "inlets-tcp-client",
		Short:        "Install inlets PRO TCP client",
		Long:         `Install an inlets PRO TCP client to any Kubernetes cluster`,
		Example:      `  arkade install inlets-tcp-client`,
		SilenceUsage: true,
	}

	inletsProClient.Flags().StringP("namespace", "n", "default", "The namespace used for installation")
	inletsProClient.Flags().Bool("update-repo", true, "Update the helm repo")

	inletsProClient.Flags().String("url", "", "URL for remote server's control-plane, i.e. wss://127.0.0.1:8123")
	inletsProClient.Flags().Bool("auto-tls", true, "Toggle use of automated TLS, fetching CA from the server on start-up. Disable when providing your own TLS termination on the server")
	inletsProClient.Flags().String("upstream", "localhost", "Forward traffic from the server here, give a hostname or IP address")
	inletsProClient.Flags().IntSlice("ports", []int{}, "Publish a TCP port on the server")

	inletsProClient.Flags().String("license", "", "License JWT or Gumroad token")
	inletsProClient.Flags().String("license-file", "", "Path to license JWT file or Gumroad token")

	inletsProClient.Flags().String("token", "", "Authentication token")
	inletsProClient.Flags().String("token-file", "", "Read the authentication token from a file")

	inletsProClient.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	inletsProClient.PreRunE = func(command *cobra.Command, args []string) error {
		tokenFile, err := command.Flags().GetString("token-file")
		if err != nil {
			return fmt.Errorf("error with --token-file usage: %s", err)
		}
		tokenString, err := command.Flags().GetString("token")
		if err != nil {
			return fmt.Errorf("error with --token usage: %s", err)
		}

		if tokenString == "" && tokenFile == "" {
			return fmt.Errorf("either --token or --token-file is required")
		}

		licenseFile, err := command.Flags().GetString("license-file")
		if err != nil {
			return fmt.Errorf("error with --license-file usage: %s", err)
		}
		licenseString, err := command.Flags().GetString("license")
		if err != nil {
			return fmt.Errorf("error with --license usage: %s", err)
		}

		if licenseString == "" && licenseFile == "" {
			return fmt.Errorf("either --license or --license-file is required")
		}

		ports, err := command.Flags().GetIntSlice("ports")
		if err != nil {
			return fmt.Errorf("error with --ports usage: %s", err)
		}

		if len(ports) == 0 {
			return fmt.Errorf("you must specify at least one remote TCP port with --ports")
		}

		return nil
	}

	inletsProClient.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		appOpts := types.DefaultInstallOptions()

		updateRepo, _ := inletsProClient.Flags().GetBool("update-repo")
		namespace, _ := inletsProClient.Flags().GetString("namespace")
		overrides := map[string]string{}

		url, _ := command.Flags().GetString("url")
		autoTLS, _ := command.Flags().GetBool("auto-tls")
		upstream, _ := command.Flags().GetString("upstream")
		ports, _ := command.Flags().GetIntSlice("ports")

		overrides["url"] = url
		overrides["autoTLS"] = strings.ToLower(strconv.FormatBool(autoTLS))
		overrides["upstream"] = upstream
		overrides["ports"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ports)), ","), "[]")
		overrides["tokenSecretName"] = "inlets-pro-token"

		tokenFile, _ := command.Flags().GetString("token-file")
		tokenString, _ := command.Flags().GetString("token")

		if len(tokenFile) > 0 {
			secretData := []types.SecretsData{
				{Type: types.FromFileSecret, Key: "token", Value: tokenFile},
			}

			tokenSecret := types.NewGenericSecret("inlets-pro-token", namespace, secretData)
			appOpts.WithSecret(tokenSecret)
		} else {
			secretData := []types.SecretsData{
				{Type: types.StringLiteralSecret, Key: "token", Value: tokenString},
			}

			tokenSecret := types.NewGenericSecret("inlets-pro-token", namespace, secretData)
			appOpts.WithSecret(tokenSecret)
		}

		licenseFile, _ := command.Flags().GetString("license-file")
		licenseString, _ := command.Flags().GetString("license")

		if len(licenseFile) > 0 {
			secretData := []types.SecretsData{
				{Type: types.FromFileSecret, Key: "license", Value: licenseFile},
			}

			licenseSecret := types.NewGenericSecret("inlets-license", namespace, secretData)
			appOpts.WithSecret(licenseSecret)
		} else {
			secretData := []types.SecretsData{
				{Type: types.StringLiteralSecret, Key: "license", Value: licenseString},
			}

			licenseSecret := types.NewGenericSecret("inlets-license", namespace, secretData)
			appOpts.WithSecret(licenseSecret)
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := mergeFlags(overrides, customFlags); err != nil {
			return err
		}

		options := appOpts.
			WithNamespace(namespace).
			WithHelmRepo("inlets-pro/inlets-pro-client").
			WithHelmURL("https://inlets.github.io/inlets-pro/charts").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(options)
		if err != nil {
			return err
		}

		fmt.Println(inletsTcpClientInstallMsg)
		return nil
	}

	return inletsProClient
}

const InletsTcpClientInfoMsg = `# Find out more at:
# https://inlets.dev`

const inletsTcpClientInstallMsg = `=======================================================================
= inlets PRO TCP client has been installed.                           =
=======================================================================` +
	"\n\n" + InletsTcpClientInfoMsg + "\n\n" + pkg.ThanksForUsing
