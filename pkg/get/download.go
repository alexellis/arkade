package get

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alexellis/arkade/pkg/archive"
	"github.com/alexellis/arkade/pkg/config"
	"github.com/alexellis/arkade/pkg/env"
)

const (
	DownloadTempDir   = iota
	DownloadArkadeDir = iota
)

func Download(tool *Tool, arch, operatingSystem, version string, downloadMode int) (string, string, error) {

	downloadURL, err := GetDownloadURL(tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version)
	if err != nil {
		return "", "", err
	}

	fmt.Println(downloadURL)

	res, err := http.DefaultClient.Get(downloadURL)
	if err != nil {
		return "", "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("incorrect status for downloading tool: %d", res.StatusCode)
	}

	_, fileName := path.Split(downloadURL)
	tmp := os.TempDir()

	outFilePath := path.Join(tmp, fileName)

	if tool.IsArchive() {
		outFilePathDir := filepath.Dir(outFilePath)
		if len(tool.BinaryTemplate) > 0 {
			fileName, err = GetBinaryName(tool, strings.ToLower(operatingSystem), strings.ToLower(arch), version)
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
		r := ioutil.NopCloser(res.Body)
		if strings.HasSuffix(downloadURL, "tar.gz") || strings.HasSuffix(downloadURL, "tgz") {
			untarErr := archive.Untar(r, outFilePathDir)
			if untarErr != nil {
				return "", "", untarErr
			}
		} else if strings.HasSuffix(downloadURL, "zip") {
			buff := bytes.NewBuffer([]byte{})
			size, err := io.Copy(buff, res.Body)
			if err != nil {
				return "", "", err
			}

			reader := bytes.NewReader(buff.Bytes())

			unzipErr := archive.Unzip(reader, size, outFilePathDir)
			if unzipErr != nil {
				return "", "", unzipErr
			}
		}

	} else {
		out, err := os.Create(outFilePath)
		if err != nil {
			return "", "", err
		}
		defer out.Close()

		if _, err = io.Copy(out, res.Body); err != nil {
			return "", "", err
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
