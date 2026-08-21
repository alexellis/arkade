// Copyright (c) arkade author(s) 2024. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/spf13/cobra"
)

func MakeOciPublish() *cobra.Command {
	command := &cobra.Command{
		Use:     "publish",
		Aliases: []string{"push", "export"},
		Short:   "Publish a directory or set of tarballs as an OCI image",
		Long: `Bundle and publish the contents of a directory, or a set of pre-built
tarballs, as an OCI image. This is the inverse of "arkade oci install":
instead of pulling an image and extracting its layers, you push a layer
built from your files and tag it on a registry.

Without --bundle, the positional directory is tarred into a single layer
and pushed. With one or more --bundle flags, each pre-built tarball is
attached as a layer carrying its platform (os/arch, i.e. linux/amd64) and
the result is published as a multi-arch index, so "arkade oci install"
can pull the correct platform for a given architecture.

Credentials come from your local Docker keychain (the same login that
docker login writes). On GitHub Actions, that keychain is populated by
the docker/login-action step for ghcr.io, so this command uses the
credentials already available in your workflow. Anonymous registries such
as ttl.sh need no login at all.`,
		Example: `  # Publish a single directory as one layer to a single tag
  arkade oci publish ./dist -t ghcr.io/me/app:0.1.0

  # Multiple tags, one -t per identifier, like docker build
  arkade oci publish ./dist \
    -t ghcr.io/me/app:0.1.0 \
    -t ghcr.io/me/app:latest

  # Omitting the tag falls back to :latest
  arkade oci publish ./dist -t ghcr.io/me/app

  # Multi-arch index from pre-built tarballs, one per platform
  arkade oci publish \
    -t ghcr.io/me/app:0.1.0 \
    --bundle linux/amd64=./app-amd64.tgz \
    --bundle linux/arm64=./app-arm64.tgz

  # Anonymous registry for testing; the :1h tag expires after an hour,
  # so only use that tag and not :latest, which would persist
  arkade oci publish ./dist -t ttl.sh/me/app:1h`,
		SilenceUsage: true,
	}

	command.Flags().StringArrayP("tag", "t", []string{}, "Image identifier (format: [registry/]repository[:tag]), repeatable, i.e. -t ghcr.io/me/app:0.1.0 -t ghcr.io/me/app:latest")
	command.Flags().StringArray("bundle", []string{}, "Pre-built tarball with a platform key, i.e. linux/amd64=./app-amd64.tgz (repeatable)")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		tags, _ := cmd.Flags().GetStringArray("tag")
		bundles, _ := cmd.Flags().GetStringArray("bundle")

		if len(tags) == 0 {
			return fmt.Errorf("please provide at least one image via -t/--tag, i.e. -t ghcr.io/me/app:0.1.0")
		}
		if len(bundles) > 0 && len(args) > 0 {
			return fmt.Errorf("provide either a source directory or --bundle, not both")
		}
		if len(bundles) == 0 && len(args) == 0 {
			return fmt.Errorf("please provide a source directory or at least one --bundle")
		}
		if len(args) > 1 {
			return fmt.Errorf("expected a single source directory, got %d arguments", len(args))
		}

		var publish func(ref string) error
		if len(bundles) > 0 {
			idx, err := publishBundles(bundles)
			if err != nil {
				return err
			}
			publish = func(ref string) error {
				r, err := name.ParseReference(ref)
				if err != nil {
					return err
				}
				return remote.WriteIndex(r, idx, remote.WithAuthFromKeychain(authn.DefaultKeychain))
			}
		} else {
			img, cleanup, err := buildImageFromDir(args[0])
			if err != nil {
				return err
			}
			if cleanup != nil {
				defer cleanup()
			}
			publish = func(ref string) error {
				r, err := name.ParseReference(ref)
				if err != nil {
					return err
				}
				return crane.Push(img, r.Name())
			}
		}

		for _, t := range tags {
			repo, embedded, hasEmbedded, err := splitImage(t)
			if err != nil {
				return err
			}
			if !hasEmbedded {
				embedded = "latest"
			}
			ref := repo + ":" + embedded
			if err := publish(ref); err != nil {
				return fmt.Errorf("publishing %s: %w", ref, err)
			}
			fmt.Printf("Published %s\n", ref)
		}
		return nil
	}

	return command
}

// splitImage separates a repository name from an optional embedded tag,
// e.g. ghcr.io/me/app:0.1.0 -> (ghcr.io/me/app, 0.1.0, true). A bare repo
// with no tag is reported as not having an explicit tag so the caller can
// default to :latest. Registry hosts with ports are handled by
// go-containerregistry's name parser.
func splitImage(image string) (repo, tag string, hasTag bool, err error) {
	ref, err := name.ParseReference(image)
	if err != nil {
		return "", "", false, fmt.Errorf("parsing image %s: %w", image, err)
	}
	repo = ref.Context().Name()
	if t, ok := ref.(name.Tag); ok {
		lastSlash := strings.LastIndex(image, "/")
		if strings.Contains(image[lastSlash+1:], ":") {
			tag = t.TagStr()
			hasTag = true
		}
	}
	return repo, tag, hasTag, nil
}

func buildImageFromDir(dir string) (v1.Image, func(), error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("source %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, nil, fmt.Errorf("source %s is not a directory", dir)
	}

	tmp, err := os.CreateTemp("", "arkade-oci-publish-*.tar")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() { os.Remove(tmp.Name()) }

	if err := tarDir(dir, tmp); err != nil {
		tmp.Close()
		cleanup()
		return nil, nil, err
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return nil, nil, err
	}

	layer, err := tarball.LayerFromFile(tmp.Name())
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("tarring %s: %w", dir, err)
	}
	img, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	cfg, err := img.ConfigFile()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	cfg = cfg.DeepCopy()
	arch, osName := env.GetClientArch()
	downloadArch, downloadOS := getDownloadArch(arch, osName)
	if downloadOS == "microsoft windows" {
		downloadOS = "windows"
	}
	cfg.OS = downloadOS
	cfg.Architecture = downloadArch
	img, err = mutate.ConfigFile(img, cfg)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return img, cleanup, nil
}

func tarDir(dir string, w io.Writer) error {
	tw := tar.NewWriter(w)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		var linkTarget string
		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}
		hdr, err := tar.FileInfoHeader(info, linkTarget)
		if err != nil {
			return err
		}
		hdr.Name = filepath.ToSlash(rel)
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, f); err != nil {
				f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return tw.Close()
}

func buildImageFromBundle(path, platform string) (v1.Image, error) {
	layer, err := tarball.LayerFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	img, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		return nil, err
	}
	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}
	cfg = cfg.DeepCopy()
	if idx := strings.Index(platform, "/"); idx > 0 {
		cfg.OS = platform[:idx]
		cfg.Architecture = platform[idx+1:]
	}
	return mutate.ConfigFile(img, cfg)
}

func publishBundles(bundles []string) (v1.ImageIndex, error) {
	var idx v1.ImageIndex = empty.Index
	for _, b := range bundles {
		parts := strings.SplitN(b, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid --bundle %q, want os/arch=file", b)
		}
		platform := strings.TrimSpace(parts[0])
		path := strings.TrimSpace(parts[1])
		platformParts := strings.Split(platform, "/")
		if len(platformParts) != 2 || platformParts[0] == "" || platformParts[1] == "" {
			return nil, fmt.Errorf("invalid platform %q in --bundle %q, want os/arch i.e. linux/amd64", platform, b)
		}
		img, err := buildImageFromBundle(path, platform)
		if err != nil {
			return nil, err
		}
		cfg, _ := img.ConfigFile()
		addendum := mutate.IndexAddendum{
			Add: img,
			Descriptor: v1.Descriptor{
				Platform: &v1.Platform{OS: cfg.OS, Architecture: cfg.Architecture},
			},
		}
		idx = mutate.AppendManifests(idx, addendum)
	}
	return idx, nil
}
