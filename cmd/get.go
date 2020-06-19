// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexellis/arkade/pkg/env"
	"github.com/alexellis/arkade/pkg/get"
	"github.com/spf13/cobra"
)

// TODO: Untar logic is copied from helm.go. Need to refactor this later on.

// Untar reads the gzip-compressed tar file from r and writes it into dir.
func Untar(r io.Reader, dir string) error {
	return untar(r, dir)
}

func untar(r io.Reader, dir string) (err error) {
	t0 := time.Now()
	nFiles := 0
	madeDir := map[string]bool{}
	defer func() {
		td := time.Since(t0)
		if err == nil {
			log.Printf("extracted tarball into %s: %d files, %d dirs (%v)", dir, nFiles, len(madeDir), td)
		} else {
			log.Printf("error extracting tarball into %s after %d files, %d dirs, %v: %v", dir, nFiles, len(madeDir), td, err)
		}
	}()
	zr, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("requires gzip-compressed body: %v", err)
	}
	tr := tar.NewReader(zr)
	loggedChtimesError := false
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("tar reading error: %v", err)
			return fmt.Errorf("tar error: %v", err)
		}
		if !validRelPath(f.Name) {
			return fmt.Errorf("tar contained invalid name error %q", f.Name)
		}
		baseFile := filepath.Base(f.Name)
		abs := path.Join(dir, baseFile)
		fmt.Println(abs, f.Name)

		fi := f.FileInfo()
		mode := fi.Mode()
		switch {
		case mode.IsDir():

			break

		case mode.IsRegular():

			wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}
			n, err := io.Copy(wf, tr)
			if closeErr := wf.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
			if err != nil {
				return fmt.Errorf("error writing to %s: %v", abs, err)
			}
			if n != f.Size {
				return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, abs, f.Size)
			}
			modTime := f.ModTime
			if modTime.After(t0) {
				// Clamp modtimes at system time. See
				// golang.org/issue/19062 when clock on
				// buildlet was behind the gitmirror server
				// doing the git-archive.
				modTime = t0
			}
			if !modTime.IsZero() {
				if err := os.Chtimes(abs, modTime, modTime); err != nil && !loggedChtimesError {
					// benign error. Gerrit doesn't even set the
					// modtime in these, and we don't end up relying
					// on it anywhere (the gomote push command relies
					// on digests only), so this is a little pointless
					// for now.
					log.Printf("error changing modtime: %v (further Chtimes errors suppressed)", err)
					loggedChtimesError = true // once is enough
				}
			}
			nFiles++
		default:
		}
	}
	return nil
}

func validRelativeDir(dir string) bool {
	if strings.Contains(dir, `\`) || path.IsAbs(dir) {
		return false
	}
	dir = path.Clean(dir)
	if strings.HasPrefix(dir, "../") || strings.HasSuffix(dir, "/..") || dir == ".." {
		return false
	}
	return true
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

func MakeGet() *cobra.Command {
	tools := get.MakeTools()

	var command = &cobra.Command{
		Use:   "get",
		Short: "Get a release of a tool or application and install it on your local computer.",
		Example: `  arkade get kubectl
  arkade get openfaas
  arkade get kubectx`,
		SilenceUsage: true,
	}

	command.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println(arkadeGet)
			return nil
		}
		var tool *get.Tool

		if len(args) == 1 {
			for _, t := range tools {
				if t.Name == args[0] {
					tool = &t
					break
				}
			}
		}
		if tool == nil {
			return fmt.Errorf("cannot get tool: %s", args[0])
		}

		fmt.Printf("Downloading %s\n", tool.Name)

		arch, operatingSystem := env.GetClientArch()
		version := ""

		downloadURL, err := get.GetDownloadURL(tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version)
		if err != nil {
			return err
		}

		fmt.Println(downloadURL)

		res, err := http.DefaultClient.Get(downloadURL)
		if err != nil {
			return err
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("incorrect status for downloading tool: %d", res.StatusCode)
		}

		_, fileName := path.Split(downloadURL)
		tmp := os.TempDir()

		outFilePath := path.Join(tmp, fileName)

		if tool.NeedDecompression == true {
			outFilePathDir := filepath.Dir(outFilePath)
			outFilePath = path.Join(outFilePathDir, tool.Name)
			r := ioutil.NopCloser(res.Body)
			untarErr := Untar(r, outFilePathDir)
			if untarErr != nil {
				return untarErr
			}
		} else {
			out, err := os.Create(outFilePath)
			if err != nil {
				return err
			}
			defer out.Close()

			if _, err = io.Copy(out, res.Body); err != nil {
				return err
			}
		}

		finalName := tool.Name
		if strings.Contains(strings.ToLower(operatingSystem), "mingw") && tool.NoExtension == false {
			finalName = finalName + ".exe"
		}

		fmt.Printf(`Tool written to: %s

Run the following to copy to install the tool:

chmod +x %s
sudo install -m 755 %s /usr/local/bin/%s
`, outFilePath, outFilePath, outFilePath, finalName)

		return err
	}

	return command
}

const arkadeGet = `Use "arkade get TOOL" to download a tool or application:

  - kubectl
  - faas-cli
  - kubectx
  `
