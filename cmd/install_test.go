// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
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

func dummyCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: name + " description",
	}
}
func TestFindAppByName(t *testing.T) {

	apps := map[string]ArkadeApp{
		"openfaas": {
			Installer:   func() *cobra.Command { return dummyCommand("openfaas") },
			InfoMessage: "OpenFaaS info",
		},
		"grafana": {
			Installer:   func() *cobra.Command { return dummyCommand("grafana") },
			InfoMessage: "Grafana info",
		},
	}

	t.Run("find existing app", func(t *testing.T) {
		app := findAppByName("openfaas", apps)
		if app == nil {
			t.Fatal("\nwant: 'openfaas'\n got: nil")
		}
		if app.InfoMessage != "OpenFaaS info" {
			t.Errorf("\nwant: 'OpenFaaS info'\n  got: '%s'", app.InfoMessage)
		}
	})

	t.Run("find non-existing app", func(t *testing.T) {
		app := findAppByName("doesnotexist", apps)
		if app != nil {
			t.Errorf("\nwant: nil\n got: '%s'", app.InfoMessage)
		}
	})
}
