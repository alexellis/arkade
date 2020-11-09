// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package k8s

import (
	"fmt"
	"strings"

	"github.com/alexellis/arkade/pkg/types"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func GetNodeArchitecture() string {
	res, _ := KubectlTask("get", "nodes", `--output`, `jsonpath={range $.items[0]}{.status.nodeInfo.architecture}`)

	arch := strings.TrimSpace(string(res.Stdout))

	return arch
}

func KubectlTask(parts ...string) (execute.ExecResult, error) {
	task := execute.ExecTask{
		Command:     "kubectl",
		Args:        parts,
		StreamStdio: false,
	}

	res, err := task.Execute()

	return res, err
}

func Kubectl(parts ...string) error {
	task := execute.ExecTask{
		Command:     "kubectl",
		Args:        parts,
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("kubectl exit code %d, stderr: %s",
			res.ExitCode,
			res.Stderr)
	}
	return nil
}

func CreateSecret(secret types.K8sSecret) error {
	secretData, err := flattenSecretData(secret.SecretData)
	if err != nil {
		return err
	}

	args := []string{"-n", secret.Namespace, "create", "secret", secret.Type, secret.Name}
	args = append(args, secretData...)

	res, secretErr := KubectlTask(args...)

	if secretErr != nil {
		return secretErr
	}
	if res.ExitCode != 0 {
		fmt.Printf("[Warning] unable to create secret %s, may already exist: %s", secret.Name, res.Stderr)
	}

	return nil
}

func flattenSecretData(data []types.SecretsData) ([]string, error) {
	var output []string

	for _, value := range data {
		switch value.Type {
		case types.StringLiteralSecret:
			output = append(output, fmt.Sprintf("--from-literal=%s=%s", value.Key, value.Value))

		case types.FromFileSecret:
			output = append(output, fmt.Sprintf("--from-file=%s=%s", value.Key, value.Value))
		default:

			return nil, fmt.Errorf("could not create secret value of type %s. Please use one of [%s, %s]", value.Type, types.StringLiteralSecret, types.FromFileSecret)

		}
	}

	return output, nil
}
