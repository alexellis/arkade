// Package license reads a license file and
// returns the license key.
package license

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ReadLicense    Read the first line of the text file in the filelocation
// FileLocation   {param}   String file location
func ReadLicense(fileLocation string) (string, error) {
	filePointer, err := os.Open(fileLocation)

	if err != nil {
		return "", err
	}

	defer filePointer.Close()
	return readContents(filePointer)
}

//keyTextLimit   Checks if the size of the text is in the range or not.
//keyLen   bool   true if the size is in the limit
func isLicenseKeyLenValid(keyLen int64) (string, error) {

	if keyLen > 1000 {
		return "", fmt.Errorf("Invalid key of length %d", keyLen)
	}

	return "", nil
}

// readContents Reads the content of the for the
// pointer to open file.
// It returns the license token or any error encountered.
func readContents(fileHndle io.Reader) (string, error) {
	var licenseKey string = ""
	var err error

	sc := bufio.NewScanner(fileHndle)
	for sc.Scan() {
		licenseKey += sc.Text()
	}

	if err := sc.Err(); err != nil {
		return "", err
	}

	return licenseKey, err
}
