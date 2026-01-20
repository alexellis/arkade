package get

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/alexellis/arkade/pkg"
	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/cheggaaa/pb/v3"
)

const (
	DownloadTempDir   = iota
	DownloadArkadeDir = iota
)

type ErrNotFound struct {
}

func (e *ErrNotFound) Error() string {
	return "server returned status: 404"
}

func Download(tool *Tool, arch, operatingSystem, version string, movePath string, displayProgress, quiet, verify bool) (string, string, error) {

	downloadURL, resolvedVersion, err := GetDownloadURL(tool,
		strings.ToLower(operatingSystem),
		strings.ToLower(arch),
		version, quiet)
	if err != nil {
		return "", "", err
	}

	if !quiet {
		fmt.Printf("Downloading: %s\n", downloadURL)
	}

	outFilePath, err := downloadFile(downloadURL, displayProgress)
	if err != nil {
		return "", "", err
	}

	if !quiet {
		fmt.Printf("%s written.\n", outFilePath)
	}

	if verify {
		if tool.VerifyStrategy == ClaudeShasumStrategy {
			st := time.Now()
			tmpl := template.New(tool.Name + "sha")
			tmpl = tmpl.Funcs(templateFuncs)
			t, err := tmpl.Parse(tool.VerifyTemplate)
			if err != nil {
				return "", "", err
			}

			var buf bytes.Buffer
			inputs := map[string]string{
				"Name":          tool.Name,
				"Owner":         tool.Owner,
				"Repo":          tool.Repo,
				"Version":       resolvedVersion,
				"VersionNumber": strings.TrimPrefix(resolvedVersion, "v"),
				"Arch":          arch,
				"OS":            operatingSystem,
			}

			if err = t.Execute(&buf, inputs); err != nil {
				return "", "", err
			}

			verifyURL := strings.TrimSpace(buf.String())
			log.Printf("Downloading SHA sum from: %s", verifyURL)
			shaSumManifest, err := fetchText(verifyURL)
			if err != nil {
				return "", "", err
			}
			// {
			//   "version": "2.1.12",
			//   "buildDate": "2026-01-17T15:42:38Z",
			//   "platforms": {
			//     "darwin-arm64": {
			//       "checksum": "40be59519a84bd35eb1111aa46f72aa6b3443866d3f6336252a198fdcaefbbe5",
			//       "size": 177846896
			//     },

			var manifest struct {
				Version   string `json:"version"`
				BuildDate string `json:"buildDate"`
				Platforms map[string]struct {
					Checksum string `json:"checksum"`
					Size     int64  `json:"size"`
				} `json:"platforms"`
			}
			if err := json.Unmarshal([]byte(shaSumManifest), &manifest); err != nil {
				return "", "", err
			}

			var archMappingForClaude = map[string]string{
				"amd64":   "amd64",
				"x86_64":  "x64",
				"arm64":   "arm64",
				"aarch64": "arm64",
			}

			platformKey := fmt.Sprintf("%s-%s", strings.ToLower(operatingSystem), archMappingForClaude[arch])

			platformInfo, found := manifest.Platforms[platformKey]
			if !found {
				return "", "", fmt.Errorf("no checksum info found for platform: %s", platformKey)
			}

			if err := verifySHA(platformInfo.Checksum, outFilePath); err != nil {
				return "", "", err
			} else {
				log.Printf("SHA sum verified in %s.", time.Since(st).Round(time.Millisecond))
			}

		} else if tool.VerifyStrategy == HashicorpShasumStrategy {
			st := time.Now()
			tmpl := template.New(tool.Name + "sha")
			tmpl = tmpl.Funcs(templateFuncs)
			t, err := tmpl.Parse(tool.VerifyTemplate)
			if err != nil {
				return "", "", err
			}

			var buf bytes.Buffer
			inputs := map[string]string{
				"Name":          tool.Name,
				"Owner":         tool.Owner,
				"Repo":          tool.Repo,
				"Version":       resolvedVersion,
				"VersionNumber": strings.TrimPrefix(resolvedVersion, "v"),
				"Arch":          arch,
				"OS":            operatingSystem,
			}

			if err = t.Execute(&buf, inputs); err != nil {
				return "", "", err
			}

			verifyURL := strings.TrimSpace(buf.String())
			log.Printf("Downloading SHA sum from: %s", verifyURL)
			shaSum, err := fetchText(verifyURL)
			if err != nil {
				return "", "", err
			}
			if err := verifySHA(shaSum, outFilePath); err != nil {
				return "", "", err
			} else {
				log.Printf("SHA sum verified in %s.", time.Since(st).Round(time.Millisecond))
			}
		}
	}

	if isArchiveStr(downloadURL) {

		outPath, err := decompress(tool, downloadURL, outFilePath, operatingSystem, arch, version, quiet)
		if err != nil {
			return "", "", err
		}

		outFilePath = outPath
		if !quiet {
			log.Printf("Extracted: %s", outFilePath)
		}
	}

	finalName := tool.Name
	if strings.Contains(strings.ToLower(operatingSystem), "mingw") && !tool.NoExtension {
		finalName = finalName + ".exe"
	}

	var localPath string

	if movePath == "" {
		_, err := config.InitUserDir()
		if err != nil {
			return "", "", err
		}

		localPath = env.LocalBinary(finalName, "")
	} else {
		localPath = filepath.Join(movePath, finalName)
	}

	if !quiet {
		log.Printf("Copying %s to %s\n", outFilePath, localPath)
	}

	if _, err = CopyFile(outFilePath, localPath); err != nil {
		return "", "", err
	}

	// Remove parent folder of the binary
	tempPath := filepath.Dir(outFilePath)
	if err := os.RemoveAll(tempPath); err != nil {
		log.Printf("Error removing temporary directory: %s", err)
	}

	outFilePath = localPath

	return outFilePath, finalName, nil
}

// DownloadFile downloads a file to a temporary directory
// and returns the path to the file and any error.
func DownloadFileP(downloadURL string, displayProgress bool) (string, error) {
	return downloadFile(downloadURL, displayProgress)
}

func downloadFile(downloadURL string, displayProgress bool) (string, error) {

	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", pkg.UserAgent())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode == http.StatusNotFound {
		return "", &ErrNotFound{}
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status: %d", res.StatusCode)
	}

	_, fileName := path.Split(downloadURL)
	tmp := os.TempDir()

	customTmp, err := os.MkdirTemp(tmp, "arkade-*")
	if err != nil {
		return "", err
	}

	outFilePath := path.Join(customTmp, fileName)
	wrappedReader := withProgressBar(res.Body, int(res.ContentLength), displayProgress)

	// Owner/Group read/write/execute
	// World - execute
	out, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		return "", err
	}

	defer out.Close()
	defer wrappedReader.Close()

	if _, err := io.Copy(out, wrappedReader); err != nil {
		return "", err
	}

	return outFilePath, nil
}

func CopyFile(src, dst string) (int64, error) {
	return CopyFileP(src, dst, 0700)
}

func CopyFileP(src, dst string, permMode int) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(permMode))
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func withProgressBar(r io.ReadCloser, length int, displayProgress bool) io.ReadCloser {
	if !displayProgress {
		return r
	}

	bar := pb.Simple.New(length).Start()
	return bar.NewProxyReader(r)
}

func decompress(tool *Tool, downloadURL, outFilePath, operatingSystem, arch, version string, quiet bool) (string, error) {

	archiveFile, err := os.Open(outFilePath)
	if err != nil {
		return "", err
	}

	outFilePathDir := filepath.Dir(outFilePath)
	if len(tool.BinaryTemplate) == 0 && len(tool.URLTemplate) > 0 {
		outFilePath = path.Join(outFilePathDir, tool.Name)
	} else if len(tool.BinaryTemplate) > 0 && len(tool.URLTemplate) == 0 &&
		(!strings.Contains(tool.BinaryTemplate, "tar.gz") &&
			!strings.Contains(tool.BinaryTemplate, "zip") &&
			!strings.Contains(tool.BinaryTemplate, "tgz")) {
		fileName, err := GetBinaryName(tool,
			strings.ToLower(operatingSystem),
			strings.ToLower(arch),
			version)
		if err != nil {
			return "", err
		}

		outFilePath = path.Join(outFilePathDir, fileName)
	} else if len(tool.BinaryTemplate) > 0 && len(tool.URLTemplate) > 0 {
		fileName, err := GetBinaryName(tool,
			strings.ToLower(operatingSystem),
			strings.ToLower(arch),
			version)
		if err != nil {
			return "", err
		}
		outFilePath = path.Join(outFilePathDir, fileName)

	} else {
		outFilePath = path.Join(outFilePathDir, tool.Name)
	}

	if strings.Contains(strings.ToLower(operatingSystem), "mingw") && tool.NoExtension == false {
		outFilePath += ".exe"
	}

	forceQuiet := true

	if strings.HasSuffix(downloadURL, "tar.gz") || strings.HasSuffix(downloadURL, "tgz") {
		if err := archive.Untar(archiveFile, outFilePathDir, true, forceQuiet); err != nil {
			return "", err
		}
	} else if strings.HasSuffix(downloadURL, "zip") {
		fInfo, err := archiveFile.Stat()
		if err != nil {
			return "", err
		}

		if !quiet {
			fmt.Printf("Name: %s, size: %d", fInfo.Name(), fInfo.Size())
		}

		if err := archive.Unzip(archiveFile, fInfo.Size(), outFilePathDir, forceQuiet); err != nil {
			return "", err
		}
	}

	return outFilePath, nil
}

func fetchText(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

func verifySHA(shaSum, outFilePath string) error {

	outFileBaseName := filepath.Base(outFilePath)

	lines := strings.Split(shaSum, "\n")

	for _, line := range lines {
		remoteHash, file, ok := strings.Cut(line, " ")
		if ok {

			if file == outFileBaseName {
				calculated, err := getSHA256Checksum(outFilePath)
				if err != nil {
					return err
				}
				if calculated != remoteHash {
					return fmt.Errorf("checksum mismatch, want: %s, but got: %s", remoteHash, calculated)
				}
			}
		}
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
