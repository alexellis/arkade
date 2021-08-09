package get

import (
	"fmt"
	"io"
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

func Download(tool *Tool, arch, operatingSystem, version string, downloadMode int, displayProgress bool) (string, string, error) {

	downloadURL, err := GetDownloadURL(tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version)
	if err != nil {
		return "", "", err
	}

	fmt.Println(downloadURL)
	outFilePath, err := downloadFile(downloadURL, displayProgress)
	if err != nil {
		return "", "", err
	}

	if isArchive, err := tool.IsArchive(); isArchive {
		if err != nil {
			return "", "", err
		}

		archiveFile, err := os.Open(outFilePath)
		if err != nil {
			return "", "", err
		}

		outFilePathDir := filepath.Dir(outFilePath)
		if len(tool.BinaryTemplate) > 0 {
			fileName, err := GetBinaryName(tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version)
			if err != nil {
				return "", "", err
			}
			outFilePath = path.Join(outFilePathDir, fileName)
		} else {
			outFilePath = path.Join(outFilePathDir, tool.Name)
		}

		if strings.Contains(strings.ToLower(operatingSystem), "mingw") && tool.NoExtension == false {
			outFilePath += ".exe"
		}

		if strings.HasSuffix(downloadURL, "tar.gz") || strings.HasSuffix(downloadURL, "tgz") {
			untarErr := archive.Untar(archiveFile, outFilePathDir)
			if untarErr != nil {
				return "", "", untarErr
			}
		} else if strings.HasSuffix(downloadURL, "zip") {
			fInfo, err := archiveFile.Stat()
			if err != nil {
				return "", "", err
			}

			fmt.Println("name", fInfo.Name(), "size: ", fInfo.Size())

			unzipErr := archive.Unzip(archiveFile, fInfo.Size(), outFilePathDir)
			if unzipErr != nil {
				return "", "", unzipErr
			}
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

		_, err = copyFile(path.Join(outFilePath), localPath)
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
