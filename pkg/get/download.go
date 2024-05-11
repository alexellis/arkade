package get

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

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

func Download(tool *Tool, arch, operatingSystem, version string, movePath string, displayProgress, quiet bool) (string, string, error) {

	downloadURL, err := GetDownloadURL(tool,
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
