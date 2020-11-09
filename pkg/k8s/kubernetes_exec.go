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
	secretData := flattenSecretData(secret.KeyValues)

	args := []string{"-n", secret.Namespace, "create", "secret", secret.Type, secret.Name}
	args = append(args, secretData...)

	res, secretErr := KubectlTask(args...)

	if secretErr != nil {
		return secretErr
	}
	if res.ExitCode != 0 {
		fmt.Printf("[Warning] unable to create secret %s, may already exist: %s", "basic-auth", res.Stderr)
	}

	return nil
}

func flattenSecretData(data map[string]string) []string {
	var output []string

	for key, value := range data {
		output = append(output, fmt.Sprintf("--from-literal=%s=%s", key, value))
	}

	return output
}
