package kubernetes

import "strings"

func GetNodeArchitecture() string {
	res, _ := KubectlTask("get", "nodes", `--output`, `jsonpath={range $.items[0]}{.status.nodeInfo.architecture}`)

	arch := strings.TrimSpace(string(res.Stdout))

	return arch
}
