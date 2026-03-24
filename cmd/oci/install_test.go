// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func TestResolveShortcutImage(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantImage         string
		wantAnonymousAuth bool
	}{
		{
			name:              "vmmeter shortcut uses anonymous auth",
			input:             "vmmeter",
			wantImage:         "ghcr.io/openfaasltd/vmmeter",
			wantAnonymousAuth: true,
		},
		{
			name:              "slicer shortcut uses anonymous auth",
			input:             "slicer",
			wantImage:         "ghcr.io/openfaasltd/slicer",
			wantAnonymousAuth: true,
		},
		{
			name:              "superterm shortcut uses anonymous auth",
			input:             "superterm",
			wantImage:         "ghcr.io/openfaasltd/superterm",
			wantAnonymousAuth: true,
		},
		{
			name:              "k3sup-pro shortcut uses anonymous auth",
			input:             "k3sup-pro",
			wantImage:         "ghcr.io/openfaasltd/k3sup-pro",
			wantAnonymousAuth: true,
		},
		{
			name:              "fully qualified image is unchanged",
			input:             "ghcr.io/openfaasltd/custom-tool",
			wantImage:         "ghcr.io/openfaasltd/custom-tool",
			wantAnonymousAuth: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotImage, gotAnonymousAuth := resolveShortcutImage(tc.input)
			if gotImage != tc.wantImage {
				t.Fatalf("want image %q, got %q", tc.wantImage, gotImage)
			}
			if gotAnonymousAuth != tc.wantAnonymousAuth {
				t.Fatalf("want anonymous auth %v, got %v", tc.wantAnonymousAuth, gotAnonymousAuth)
			}
		})
	}
}

func TestBuildPullOptions(t *testing.T) {
	platform := &v1.Platform{Architecture: "amd64", OS: "linux"}

	t.Run("default uses keychain auth", func(t *testing.T) {
		opts := crane.GetOptions(buildPullOptions(platform, false)...)
		if opts.Platform != platform {
			t.Fatalf("want platform %p, got %p", platform, opts.Platform)
		}

		authOptionName := runtime.FuncForPC(reflect.ValueOf(opts.Remote[0]).Pointer()).Name()
		if !strings.Contains(authOptionName, "WithAuthFromKeychain") {
			t.Fatalf("want keychain auth option, got %q", authOptionName)
		}
	})

	t.Run("anonymous auth overrides keychain", func(t *testing.T) {
		opts := crane.GetOptions(buildPullOptions(platform, true)...)
		if opts.Platform != platform {
			t.Fatalf("want platform %p, got %p", platform, opts.Platform)
		}

		authOptionName := runtime.FuncForPC(reflect.ValueOf(opts.Remote[0]).Pointer()).Name()
		if !strings.Contains(authOptionName, "WithAuth") || strings.Contains(authOptionName, "WithAuthFromKeychain") {
			t.Fatalf("want anonymous auth option, got %q", authOptionName)
		}
	})
}
