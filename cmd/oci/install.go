// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package oci

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/spf13/cobra"
)

func MakeOciInstall() *cobra.Command {
	command := &cobra.Command{
		Use:     "install IMAGE [PATH]",
		Aliases: []string{"i", "extract"},
		Short:   "Install the contents of an OCI image to a given path",
		Long: `Use this command to install binaries or packages distributed within an
OCI image.`,
		Example: `  # Install slicer to /usr/local/bin (default)
  # Files will be extracted to /usr/local/bin/slicer
  arkade oci install ghcr.io/openfaasltd/slicer

  # Install to current directory
  arkade oci install ghcr.io/openfaasltd/slicer .

  # Install to a custom path like /tmp/
  # Files will be extracted to /tmp/slicer
  arkade oci install ghcr.io/openfaasltd/slicer /tmp --version 0.1.0

  # Install slicer for arm64 as an architecture override, instead of using uname
  arkade oci install ghcr.io/openfaasltd/slicer --arch arm64

  # Use a shortcut for the image name (vmmeter, slicer, k3sup-pro)
  arkade oci install k3sup-pro
`,
		SilenceUsage: true,
	}

	command.Flags().StringP("version", "v", "latest", "The version or leave blank to determine the latest available version")
	command.Flags().String("path", "/usr/local/bin", "(deprecated: use positional argument) Installation path")
	command.Flags().Bool("progress", true, "Show download progress")
	command.Flags().String("arch", "", "CPU architecture i.e. amd64")
	command.Flags().String("os", "", "OS i.e. linux")

	command.Flags().BoolP("gzipped", "g", false, "Is this a gzipped tarball?")
	command.Flags().Bool("quiet", false, "Suppress progress output")
	command.Flags().Bool("symlink", false, "Write symlinks when unpacking OCI image, only use with trusted sources")

	// Hide the deprecated --path flag
	command.Flags().MarkHidden("path")

	command.PreRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetString("version")
		gzipped, _ := cmd.Flags().GetBool("gzipped")
		quiet, _ := cmd.Flags().GetBool("quiet")
		allowSymlinks, _ := cmd.Flags().GetBool("symlink")
		showProgress, _ := cmd.Flags().GetBool("progress")

		if len(args) < 1 {
			return fmt.Errorf("please provide an image name")
		}

		imageName := args[0]

		// Determine installation path
		// Priority: arg[1] > --path flag > default
		installPath := "/usr/local/bin"
		if len(args) >= 2 {
			installPath = args[1]
		} else if cmd.Flags().Changed("path") {
			installPath, _ = cmd.Flags().GetString("path")
		}

		imageName, forceAnonymousAuth := resolveShortcutImage(imageName)

		if !strings.Contains(imageName, ":") {
			imageName = imageName + ":" + version
		}

		st := time.Now()

		if err := os.MkdirAll(installPath, 0755); err != nil && !os.IsExist(err) {
			fmt.Printf("Error creating directory %s, error: %s\n", installPath, err.Error())
		}

		clientArch, clientOS := env.GetClientArch()

		if cmd.Flags().Changed("arch") {
			clientArch, _ = cmd.Flags().GetString("arch")
		}

		if cmd.Flags().Changed("os") {
			clientOS, _ = cmd.Flags().GetString("os")
		} else {
			if strings.Contains(clientOS, "microsoft windows") {
				clientOS = "windows"
			}
		}

		tempFile, err := os.CreateTemp(os.TempDir(), "arkade-oci-*")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}

		defer os.Remove(tempFile.Name())

		f, err := os.Create(tempFile.Name())
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", tempFile.Name(), err)
		}
		defer f.Close()

		downloadArch, downloadOS := getDownloadArch(clientArch, clientOS)
		platform := &v1.Platform{Architecture: downloadArch, OS: downloadOS}

		// ── progress state ───────────────────────────────────
		p := &imageProgress{
			imageName: imageName,
			platform:  fmt.Sprintf("%s/%s", downloadOS, downloadArch),
			status:    stResolving,
			started:   time.Now(),
		}

		// Counting transport: every response body the crane stack
		// reads is wrapped, giving us live network-byte counts.
		ct := &countingTransport{base: remote.DefaultTransport, n: &p.bytesRead}
		opts := buildPullOptions(platform, forceAnonymousAuth)
		opts = append(opts, crane.WithTransport(ct))

		tty := !quiet && isTTY()
		renderLive := showProgress && !quiet

		// In TTY mode output goes to stderr so stdout stays clean for piping.
		out := os.Stdout
		if tty {
			out = os.Stderr
		}

		if !quiet {
			fmt.Fprintf(out, "Installing %s to %s\n", imageName, installPath)
		}

		if tty && renderLive {
			fmt.Fprint(out, "\033[?1049h") // enter alternate screen
			fmt.Fprint(out, "\033[?25l")   // hide cursor
		}
		leaveAlt := func() {
			if tty && renderLive {
				fmt.Fprint(out, "\033[?25h")   // restore cursor
				fmt.Fprint(out, "\033[?1049l") // leave alternate screen
			}
		}

		// Restore terminal on signal.
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-signalChan
			leaveAlt()
			os.Exit(2)
		}()

		// Worker performs the crane operations. We send the eventual
		// error (or nil) plus a phase-change signal back to the main
		// goroutine so the renderer can flip from resolving → downloading
		// at the right moment.
		type phase struct {
			downloading bool
			done        bool
			err         error
		}
		ph := make(chan phase, 2)

		go func() {
			img, pullErr := crane.Pull(imageName, opts...)
			if pullErr != nil {
				ph <- phase{done: true, err: fmt.Errorf("pulling %s: %w", imageName, pullErr)}
				return
			}
			// Compute total bytes from the manifest (sum of compressed
			// layer sizes). The manifest fetch itself contributes a few
			// KB to bytesRead, so reset before downloading begins.
			if manifest, mErr := img.Manifest(); mErr == nil {
				var total int64
				for _, l := range manifest.Layers {
					total += l.Size
				}
				atomic.StoreInt64(&p.totalBytes, total)
			}
			atomic.StoreInt64(&p.bytesRead, 0)
			ph <- phase{downloading: true}

			if expErr := crane.Export(img, f); expErr != nil {
				ph <- phase{done: true, err: fmt.Errorf("exporting %s: %w", imageName, expErr)}
				return
			}
			ph <- phase{done: true}
		}()

		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		var plainPrev string
		render := func() {
			if !renderLive {
				return
			}
			if tty {
				renderTTY(out, p)
			} else {
				plainPrev = renderPlain(out, p, plainPrev)
			}
		}

		render() // initial frame (resolving)

		var workErr error
	loop:
		for {
			select {
			case ev := <-ph:
				if ev.downloading {
					p.status = stDownloading
					render()
					continue
				}
				if ev.done {
					workErr = ev.err
					break loop
				}
			case <-ticker.C:
				render()
			}
		}

		// finalize: extract if download succeeded.
		if workErr == nil {
			p.status = stExtracting
			render()

			tarFile, openErr := os.Open(tempFile.Name())
			if openErr != nil {
				workErr = fmt.Errorf("failed to open %s: %w", tempFile.Name(), openErr)
			} else {
				defer tarFile.Close()
				// When the alt-screen is active, suppress UntarNested's
				// per-file logging so it doesn't corrupt the live frame.
				untarQuiet := quiet || (tty && renderLive)
				if uErr := archive.UntarNested(tarFile, installPath, gzipped, untarQuiet, allowSymlinks); uErr != nil {
					workErr = fmt.Errorf("failed to untar %s: %w", tempFile.Name(), uErr)
				}
			}
		}

		if workErr != nil {
			p.status = stFailed
			p.err = workErr
		} else {
			p.status = stDone
		}
		p.elapsed = time.Since(p.started)

		if renderLive {
			if tty {
				renderTTY(out, p)
				leaveAlt()
				renderTTYFinal(out, p)
			} else {
				renderPlain(out, p, plainPrev)
			}
		} else if tty {
			leaveAlt()
		}

		if workErr != nil {
			return workErr
		}

		if !quiet {
			fmt.Fprintf(out, "Took %s\n", time.Since(st).Round(time.Millisecond))
		}
		return nil
	}

	return command
}

func resolveShortcutImage(imageName string) (string, bool) {
	switch imageName {
	case "vmmeter":
		return "ghcr.io/openfaasltd/vmmeter", true
	case "slicer":
		return "ghcr.io/openfaasltd/slicer", true
	case "superterm":
		return "ghcr.io/openfaasltd/superterm", true
	case "k3sup-pro":
		return "ghcr.io/openfaasltd/k3sup-pro", true
	default:
		return imageName, false
	}
}

func buildPullOptions(platform *v1.Platform, forceAnonymousAuth bool) []crane.Option {
	options := []crane.Option{
		crane.WithPlatform(platform),
	}

	if forceAnonymousAuth {
		options = append(options, crane.WithAuth(authn.Anonymous))
	}

	return options
}

func getDownloadArch(clientArch, clientOS string) (arch string, os string) {
	downloadArch := strings.ToLower(clientArch)
	downloadOS := strings.ToLower(clientOS)

	if downloadArch == "x86_64" {
		downloadArch = "amd64"
	} else if downloadArch == "aarch64" {
		downloadArch = "arm64"
	}

	return downloadArch, downloadOS
}
