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

func MakeInstallKubernetesDashboard() *cobra.Command {
	var kubeDashboard = &cobra.Command{
		Use:          "kubernetes-dashboard",
		Short:        "Install kubernetes-dashboard",
		Long:         `Install kubernetes-dashboard`,
		Example:      `  arkade install kubernetes-dashboard`,
		SilenceUsage: true,
	}

	kubeDashboard.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath, _ := command.Flags().GetString("kubeconfig")
		if err := config.SetKubeconfig(kubeConfigPath); err != nil {
			return err
		}

		arch := k8s.GetNodeArchitecture()
		fmt.Printf("Node architecture: %q\n", arch)

		_, err := k8s.KubectlTask("apply", "-f",
			"https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-rc5/aio/deploy/recommended.yaml")
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

#To forward the dashboard to your local machine 
kubectl proxy

#To get your Token for logging in
kubectl -n kubernetes-dashboard describe secret $(kubectl -n kubernetes-dashboard get secret | grep admin-user-token | awk '{print $1}')

# Once Proxying you can navigate to the below
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/login`

const KubernetesDashboardInstallMsg = `=======================================================================
= Kubernetes Dashboard has been installed.                            =
=======================================================================` +
	"\n\n" + KubernetesDashboardInfoMsg + "\n\n" + pkg.ThanksForUsing
