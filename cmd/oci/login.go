// Copyright (c) arkade author(s) 2026. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/types"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"
)

type ociLoginOptions struct {
	serverAddress string
	user          string
	password      string
	passwordStdin bool
}

func MakeOciLogin() *cobra.Command {
	opts := &ociLoginOptions{}

	command := &cobra.Command{
		Use:   "login REGISTRY",
		Short: "Log in to a container registry",
		Long: `Log in to a container registry, storing credentials in your Docker
config file (~/.docker/config.json), honouring any credential helper
configured there such as the macOS keychain or Docker Desktop.

The saved credentials are picked up automatically by "arkade oci install"
and "arkade oci publish" for private registries.

This is a narrow equivalent of "crane auth login". For interactive
login prompts, identity tokens, or credential exports, install crane
with "arkade get crane" or use "docker login" instead.`,
		Example: `  # Log in to reg.example.com
  arkade oci login reg.example.com -u user -p secret

  # Read the password from stdin
  cat ~/ghcr-token | arkade oci login ghcr.io -u user --password-stdin`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := name.NewRegistry(args[0])
			if err != nil {
				return fmt.Errorf("invalid registry %q: %w", args[0], err)
			}

			opts.serverAddress = reg.Name()

			return ociLogin(opts)
		},
	}

	command.Flags().StringVarP(&opts.user, "username", "u", "", "Username")
	command.Flags().StringVarP(&opts.password, "password", "p", "", "Password")
	command.Flags().BoolVarP(&opts.passwordStdin, "password-stdin", "", false, "Take the password from stdin")

	return command
}

func ociLogin(opts *ociLoginOptions) error {
	if opts.passwordStdin {
		contents, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		opts.password = strings.TrimSuffix(string(contents), "\n")
		opts.password = strings.TrimSuffix(opts.password, "\r")
	}

	if opts.user == "" && opts.password == "" {
		return errors.New("username and password required")
	}

	cf, err := config.Load(os.Getenv("DOCKER_CONFIG"))
	if err != nil {
		return err
	}

	serverAddress := opts.serverAddress
	creds := cf.GetCredentialsStore(serverAddress)
	if serverAddress == name.DefaultRegistry {
		serverAddress = authn.DefaultAuthKey
	}

	if err := creds.Store(types.AuthConfig{
		ServerAddress: serverAddress,
		Username:      opts.user,
		Password:      opts.password,
	}); err != nil {
		return err
	}

	if err := cf.Save(); err != nil {
		return err
	}

	fmt.Printf("Logged in via %s\n", cf.Filename)

	return nil
}
