// Copyright (c) arkade author(s) 2026. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/docker/cli/cli/config"
	"github.com/google/go-containerregistry/pkg/authn"
)

func mask(s string) string {
	if s == "" {
		return ""
	}

	return strings.Repeat("*", len(s))
}

func TestOciLoginStoresCredentials(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("DOCKER_CONFIG", tempDir)

	opts := &ociLoginOptions{
		serverAddress: "ghcr.io",
		user:          "alex",
		password:      "hunter2",
	}

	if err := ociLogin(opts); err != nil {
		t.Fatalf("ociLogin error: %s", err)
	}

	cf, err := config.Load(tempDir)
	if err != nil {
		t.Fatalf("loading written config: %s", err)
	}

	authCfg, ok := cf.AuthConfigs["ghcr.io"]
	if !ok {
		t.Fatal("want ghcr.io entry in auths")
	}
	if authCfg.Username != "alex" || authCfg.Password != "hunter2" {
		t.Fatalf("want credentials for alex, got username %q", authCfg.Username)
	}

	wantUser, wantPass := "alex", "hunter2"
	if authCfg.Username != wantUser || authCfg.Password != wantPass {
		t.Fatalf("want credentials for %s, got username %q password %q",
			wantUser, authCfg.Username, mask(authCfg.Password))
	}

	if _, err := os.Stat(path.Join(tempDir, "config.json")); err != nil {
		t.Fatalf("want config.json written: %s", err)
	}
}

func TestOciLoginDefaultRegistryKey(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("DOCKER_CONFIG", tempDir)

	opts := &ociLoginOptions{
		serverAddress: "index.docker.io",
		user:          "alex",
		password:      "hunter2",
	}

	if err := ociLogin(opts); err != nil {
		t.Fatalf("ociLogin error: %s", err)
	}

	cf, err := config.Load(tempDir)
	if err != nil {
		t.Fatalf("loading written config: %s", err)
	}

	if _, ok := cf.AuthConfigs[authn.DefaultAuthKey]; !ok {
		t.Fatalf("want default registry key %q in auths, got keys from config", authn.DefaultAuthKey)
	}
}

func TestOciLoginRequiresCredentials(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("DOCKER_CONFIG", tempDir)

	err := ociLogin(&ociLoginOptions{serverAddress: "ghcr.io"})
	if err == nil {
		t.Fatal("want error when username and password are empty")
	}
}
