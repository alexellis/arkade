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

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
	"github.com/cheggaaa/pb/v3"
)

const (
	DownloadTempDir   = iota
	DownloadArkadeDir = iota
)

func Download(tool *Tool, arch, operatingSystem, version string, downloadMode int, displayProgress, quiet bool) (string, string, error) {

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

	if isArchive, err := tool.IsArchive(quiet); isArchive {
		if err != nil {
			return "", "", err
		}

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
	if strings.Contains(strings.ToLower(operatingSystem), "mingw") && tool.NoExtension == false {
		finalName = finalName + ".exe"
	}

	if downloadMode == DownloadArkadeDir {
		_, err := config.InitUserDir()
		if err != nil {
			return "", "", err
		}

		localPath := env.LocalBinary(finalName, "")

		if !quiet {
			log.Printf("Copying %s to %s\n", outFilePath, localPath)
		}
		_, err = copyFile(outFilePath, localPath)
		if err != nil {
			return "", "", err
		}

		outFilePath = localPath
	}

	return outFilePath, finalName, nil
}

func downloadFile(downloadURL string, displayProgress bool) (string, error) {
	res, err := http.DefaultClient.Get(downloadURL)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("incorrect status for downloading tool: %d", res.StatusCode)
	}

	_, fileName := path.Split(downloadURL)
	tmp := os.TempDir()
	outFilePath := path.Join(tmp, fileName)
	wrappedReader := withProgressBar(res.Body, int(res.ContentLength), displayProgress)
	out, err := os.Create(outFilePath)
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

func copyFile(src, dst string) (int64, error) {
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

	userReadWriteExecute := 0700
	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(userReadWriteExecute))
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
		if err := archive.Untar(archiveFile, outFilePathDir, forceQuiet); err != nil {
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
