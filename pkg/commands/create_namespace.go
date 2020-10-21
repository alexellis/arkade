package commands

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/k8s"
)

func CreateNamespace(namespace string) error {
	getNs, err := k8s.KubectlTask("get", "namespace", namespace)
	if err != nil {
		return err
	}

	if getNs.ExitCode == 0 {
		fmt.Println(fmt.Sprintf("[Info] namespace exists: %s", namespace))
		return nil
	}

	nsRes, err := k8s.KubectlTask("create", "namespace", namespace)
	if err != nil {
		return err
	}

	if nsRes.ExitCode != 0 {
		fmt.Println(fmt.Sprintf("[Error] unable to create namespace %s: %s", namespace, nsRes.Stderr))
	}
	fmt.Println(fmt.Sprintf("[Info] namespace created: %s", namespace))
	return nil
}
