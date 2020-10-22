package commands

import (
	"fmt"

	"github.com/spf13/pflag"

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

func GetNamespace(flags *pflag.FlagSet, defaultNs string) (string, error) {
	namespace, err := flags.GetString("namespace")
	if err != nil {
		return namespace, err
	}

	if len(namespace) == 0 {
		namespace = defaultNs
	}

	return namespace, nil
}
