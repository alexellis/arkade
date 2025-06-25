// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

type Tool = get.Tool

func TestCheckForTool(t *testing.T) {
	tests := []struct {
		name       string
		appName    string
		tools      []Tool
		want       string
		expectFail bool
	}{
		{
			name:    "Tool exists but expected is wrong",
			appName: "kubectl",
			tools: []Tool{
				{Name: "kubectl"},
				{Name: "helm"},
				{Name: "k9s"},
			},
			want:       "this test should fail",
			expectFail: true,
		},
		{
			name:    "Tool exists",
			appName: "kubectl",
			tools: []Tool{
				{Name: "kubectl"},
				{Name: "helm"},
				{Name: "k9s"},
			},
			want:       "no such app. kubectl is available as a tool, run \"arkade get kubectl\" to get it",
			expectFail: false,
		},
		{
			name:    "Tool does not exist",
			appName: "randomtool",
			tools: []Tool{
				{Name: "kubectl"},
				{Name: "helm"},
			},
			want:       "no such app: randomtool, run \"arkade install --help\" for a list of apps",
			expectFail: false,
		},
		{
			name:    "Case-insensitive match",
			appName: "KUBECTL",
			tools: []Tool{
				{Name: "kubectl"},
				{Name: "helm"},
			},
			want:       "no such app. KUBECTL is available as a tool, run \"arkade get KUBECTL\" to get it",
			expectFail: false,
		},
		{
			name:       "Empty tool list",
			appName:    "kubectl",
			tools:      []Tool{},
			want:       "no such app: kubectl, run \"arkade install --help\" for a list of apps",
			expectFail: false,
		},
		{
			name:    "Empty app name",
			appName: "",
			tools: []Tool{
				{Name: "kubectl"},
			},
			want:       "no such app: , run \"arkade install --help\" for a list of apps",
			expectFail: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := checkForTool(tc.appName, tc.tools)
			if got != tc.want && !tc.expectFail {
				t.Errorf("%s\nwant: %s\n got: %s",
					tc.name, tc.want, got)
			}
		})
	}
}

func TestRenderTable(t *testing.T) {

	var buf bytes.Buffer

	appMap := map[string]ArkadeApp{
		"argocd": {
			Name: "argocd",
			Installer: func() *cobra.Command {
				return &cobra.Command{Short: "Install Argo CD"}
			},
		},
		"cert-manager": {
			Name: "cert-manager",
			Installer: func() *cobra.Command {
				return &cobra.Command{Short: "Install Cert Manager"}
			},
		},
	}
	expected := `|     TOOL     |     DESCRIPTION      |
|--------------|----------------------|
| argocd       | Install Argo CD      |
| cert-manager | Install Cert Manager |

There are 2 apps that you can install on your cluster.
`
	renderTable(&buf, appMap)
	actual := buf.String()

	want := strings.TrimSpace(expected)
	got := strings.TrimSpace(actual)

	if actual != expected {
		t.Errorf("Output did not match expected.\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}
