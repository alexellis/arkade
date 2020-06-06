// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package k8s

import (
	"fmt"
	"strings"

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
