// Copyright (c) arkade author(s) 2022. All rights reserved.
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

func MakeInstallKubernetesDashboard() *cobra.Command {
	var kubeDashboard = &cobra.Command{
		Use:          "kubernetes-dashboard",
		Short:        "Install kubernetes-dashboard",
		Long:         `Install kubernetes-dashboard`,
		Example:      `  arkade install kubernetes-dashboard`,
		SilenceUsage: true,
	}

	kubeDashboard.Flags().StringP("namespace", "n", "kubernetes-dashboard", "The namespace to install the chart")
	kubeDashboard.Flags().StringArray("set", []string{}, "Use custom flags or override existing flags \n(example --set image.tag=v2.5.0)")

	kubeDashboard.PreRunE = func(cmd *cobra.Command, args []string) error {
		_, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetStringArray("set")
		if err != nil {
			return err
		}

		return nil
	}

	kubeDashboard.RunE = func(cmd *cobra.Command, args []string) error {
		kubeConfigPath, _ := cmd.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		customFlags, _ := cmd.Flags().GetStringArray("set")
		namespace, _ := cmd.Flags().GetString("namespace")

		overrides := map[string]string{}
		if err := config.MergeFlags(overrides, customFlags); err != nil {
			return err
		}

		k8sDashboardOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("kubernetes-dashboard/kubernetes-dashboard").
			WithHelmURL("https://kubernetes.github.io/dashboard/").
			WithOverrides(overrides).
			WithKubeconfigPath(kubeConfigPath)

		_, err := apps.MakeInstallChart(k8sDashboardOptions)
		if err != nil {
			return err
		}

		fmt.Println(KubernetesDashboardInstallMsg)

		return nil
	}

	return kubeDashboard
}

const KubernetesDashboardInfoMsg = `# To create the Service Account and the ClusterRoleBinding
# @See https://github.com/kubernetes/dashboard/blob/master/docs/user/access-control/creating-sample-user.md#creating-sample-user

cat <<EOF | kubectl apply -f -
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: kubernetes-dashboard
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: kubernetes-dashboard
---
EOF

# To forward the dashboard to your local machine 
kubectl -n kubernetes-dashboard port-forward services/kubernetes-dashboard 8443:443

# To get your Token for logging in

## K8s v1.24 or above
### Generate token without storing it
kubectl -n kubernetes-dashboard create token admin-user

### Generate a token and store it in a secret
kubectl -n kubernetes-dashboard create token admin-user --bound-object-kind Secret --bound-object-name admin-user-token

## K8s v1.23 or below
kubectl -n kubernetes-dashboard describe secret $(kubectl -n kubernetes-dashboard get secret | grep admin-user-token | awk '{print $1}')

# Once Proxying you can navigate to the below
https://127.0.0.1:8443/#/login`

const KubernetesDashboardInstallMsg = `=======================================================================
= Kubernetes Dashboard has been installed.                            =
=======================================================================` +
	"\n\n" + KubernetesDashboardInfoMsg + "\n\n" + pkg.SupportMessageShort
