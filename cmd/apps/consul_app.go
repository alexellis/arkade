// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/alexellis/arkade/pkg/config"
	"strconv"
	"strings"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"

	"github.com/alexellis/arkade/pkg/k8s"

	"github.com/alexellis/arkade/pkg"
	"github.com/spf13/cobra"
)

func MakeInstallConsul() *cobra.Command {
	var consul = &cobra.Command{
		Use:          "consul-connect",
		Short:        "Install Consul Service Mesh",
		Long:         `Install Consul Service Mesh to any Kubernetes cluster`,
		Example:      `  arkade install consul-connect`,
		SilenceUsage: true,
	}

	consul.Flags().StringP("namespace", "n", "consul-system", "The namespace used for installation")
	consul.Flags().Bool("update-repo", true, "Update the helm repo")
	consul.Flags().StringP("datacenter", "d", "dc1", "The name of the datacenter that the agents should register as")
	consul.Flags().Bool("enable-connect-injector", true, "If true, all the resources necessary for the Connect injector process to run will be installed")
	consul.Flags().Bool("enable-tls-encryption", true, "If true, TLS encryption across the cluster to verify authenticity of the Consul servers and clients is enabled")
	consul.Flags().Bool("enable-gossip-encryption", true, "If true, Consul's gossip encryption is enabled")
	consul.Flags().String("gossip-encryption-key", "", "The gossip encryption key; when empty, a new, random key is generated")
	consul.Flags().Bool("manage-system-acls", true, "If true, the ACL tokens and policies for all Consul and consul-k8s components will automatically be managed")
	consul.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image=org/repo:tag)")

	consul.RunE = func(command *cobra.Command, args []string) error {
		updateRepo, _ := consul.Flags().GetBool("update-repo")

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		namespace, _ := consul.Flags().GetString("namespace")

		overrides := map[string]string{}
		overrides["global.name"] = "consul"

		datacenter, _ := command.Flags().GetString("datacenter")
		overrides["global.datacenter"] = datacenter

		connectInjectorEnabled, _ := command.Flags().GetBool("enable-connect-injector")
		overrides["connectInject.enabled"] = strings.ToLower(strconv.FormatBool(connectInjectorEnabled))

		tlsEnabled, _ := command.Flags().GetBool("enable-tls-encryption")
		overrides["global.tls.enabled"] = strings.ToLower(strconv.FormatBool(tlsEnabled))

		manageSystemACLs, _ := command.Flags().GetBool("manage-system-acls")
		overrides["global.acls.manageSystemACLs"] = strings.ToLower(strconv.FormatBool(manageSystemACLs))

		gossipEncryptionEnabled, _ := command.Flags().GetBool("enable-gossip-encryption")
		gossipEncryptionKey, _ := command.Flags().GetString("gossip-encryption-key")

		var err error
		if gossipEncryptionEnabled && gossipEncryptionKey == "" {
			gossipEncryptionKey, err = generateGossipEncryptionKey()
			if err != nil {
				return err
			}
		}

		if gossipEncryptionEnabled {
			res, err := k8s.KubectlTask("create", "secret", "generic",
				"consul-gossip-encryption-key",
				"--namespace="+namespace,
				"--from-literal", "key="+gossipEncryptionKey)
			if err != nil {
				return err
			} else if len(res.Stderr) > 0 && strings.Contains(res.Stderr, "AlreadyExists") {
				fmt.Println("[Warning] secret consul-gossip-encryption-key already exists and will be used.")
			} else if len(res.Stderr) > 0 {
				return fmt.Errorf("error from kubectl\n%q", res.Stderr)
			}

			overrides["global.gossipEncryption.secretName"] = "consul-gossip-encryption-key"
			overrides["global.gossipEncryption.secretKey"] = "key"
		}

		customFlags, _ := command.Flags().GetStringArray("set")

		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")

		consulOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("hashicorp/consul").
			WithHelmURL("https://helm.releases.hashicorp.com").
			WithOverrides(overrides).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(consulOptions)
		if err != nil {
			return err
		}

		fmt.Println(consulInstallMsg)
		return nil
	}

	return consul
}

func generateGossipEncryptionKey() (string, error) {
	key := make([]byte, 32)
	n, err := rand.Reader.Read(key)
	if err != nil {
		return "", err
	}
	if n != 32 {
		return "", fmt.Errorf("couldn't read enough entropy, generate more entropy")
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

const ConsulInfoMsg = `# Find out more at:
# https://www.consul.io/docs/k8s`

const consulInstallMsg = `=======================================================================
= Consul has been installed.                                          =
=======================================================================` +
	"\n\n" + ConsulInfoMsg + "\n\n" + pkg.SupportMessageShort
