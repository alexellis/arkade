// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/alexellis/go-execute/v2"
	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

func MakeUpdate() *cobra.Command {
	var command = &cobra.Command{
		Use:   "update",
		Short: "Replace the running binary with an updated version",
		Long: `The latest release version of arkade will be downloaded from GitHub.

If that text is not found by running "arkade version", then the release
will be downloaded, along with its .sha256 checksum file.

If the checksum matches the downloaded file then the running binary will be
replaced with the new binary.

This command can be run as often as you require, won't download the same 
version twice.`,
		Example:       `  arkade update`,
		Aliases:       []string{"u"},
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	command.Flags().Bool("verify", true, "Verify the checksum of the downloaded binary")
	command.Flags().Bool("force", false, "Force a download of the latest binary, even if up to date, the --verify flag still applies")

	command.RunE = func(cmd *cobra.Command, args []string) error {

		verifyDigest, _ := cmd.Flags().GetBool("verify")
		forceDownload, _ := cmd.Flags().GetBool("force")

		name := "arkade"
		toolList := get.MakeTools()
		var tool *get.Tool
		for _, t := range toolList {
			if t.Name == name {
				tool = &t
				break
			}
		}

		release, err := get.FindGitHubRelease("alexellis", name)
		if err != nil {
			return err
		}

		executable, err := os.Executable()
		if err != nil {
			return err
		}

		task := execute.ExecTask{
			Command: executable,
			Args:    []string{"version"},
		}

		res, err := task.Execute(context.TODO())
		if err != nil {
			return err
		}

		fmt.Printf("Latest release: %s\n", release)

		if !forceDownload && strings.Contains(res.Stdout, release) {
			fmt.Println("You are already using the latest version of arkade.")

			fmt.Println("\n\n", aec.Bold.Apply(pkg.SupportMessageShort))

			return nil
		}

		arch, operatingSystem := env.GetClientArch()
		arch = strings.ToLower(arch)
		operatingSystem = strings.ToLower(operatingSystem)

		if arch == "x86_64" {
			arch = "amd64"
		}

		downloadUrl, err := get.GetDownloadURL(tool, operatingSystem, arch, release, false)
		if err != nil {
			return err
		}

		newBinary, err := get.DownloadFileP(downloadUrl, true)
		if err != nil {
			return err
		}

		if verifyDigest {
			digest, err := downloadDigest(downloadUrl + ".sha256")
			if err != nil {
				return err
			}

			if err := compareSHA(digest, newBinary); err != nil {
				return fmt.Errorf("checksum failed for %s, error: %w", newBinary, err)
			}

			fmt.Printf("Checksum verified..OK.\n")
		}

		if err := replaceExec(executable, newBinary); err != nil {
			return err
		}

		fmt.Printf("Replaced: %s..OK.", executable)

		fmt.Println("\n\n", aec.Bold.Apply(pkg.SupportMessageShort))

		return nil
	}
	return command
}

func downloadDigest(uri string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	var body []byte
	if res.Body != nil {
		defer res.Body.Close()
		body, _ = io.ReadAll(res.Body)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d, body: %s", res.StatusCode, string(body))
	}

	return string(body), nil
}

// Copy the new binary to the same directory as the current binary before calling os.Rename to prevent an
// 'invalid cross-device link' error because the source and destination are not on the same file system.
func replaceExec(currentExec, newBinary string) error {
	targetDir := filepath.Dir(currentExec)
	filename := filepath.Base(currentExec)
	newExec := filepath.Join(targetDir, fmt.Sprintf(".%s.new", filename))

	// Copy the contents of newbinary to a new executable file
	sf, err := os.Open(newBinary)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.OpenFile(newExec, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer df.Close()

	if _, err := io.Copy(df, sf); err != nil {
		return err
	}

	// Replace the current executable file with the new executable file
	if err := os.Rename(newExec, currentExec); err != nil {
		return err
	}

	return nil
}

// compareSHA returns a nil error if the local digest matches the remote digest
func compareSHA(remoteDigest, localFile string) error {

	// GitHub format may sometimes include the binary name and a space, i.e.
	// "9dcfd1611440aa15333980b860220bcd55ca1d6875692facc458caf7eb1cd042  bin/arkade-darwin-arm64"
	if strings.Contains(remoteDigest, " ") {
		t, _, _ := strings.Cut(remoteDigest, " ")
		remoteDigest = t
	}

	localDigest, err := getSHA256Checksum(localFile)
	if err != nil {
		return err
	}

	if remoteDigest != localDigest {
		return fmt.Errorf("checksum mismatch, want: %s, but got: %s", remoteDigest, localDigest)
	}

	return nil
}

func getSHA256Checksum(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(f)), nil
}
