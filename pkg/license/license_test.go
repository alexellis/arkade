package license

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestReadFileContent(t *testing.T) {
	expectedLicenseKey := "Valid license key.\nOther content. \nOther content."
	expectedLicenseKeyReader := strings.NewReader(expectedLicenseKey)

	t.Run("Should return the same content", func(t *testing.T) {
		actualLicenceKey, err := readContents(expectedLicenseKeyReader)

		if err != nil && actualLicenceKey == expectedLicenseKey {
			t.Errorf("readContents() error = %v", err)
		}
	})
}

func TestReadLicense(t *testing.T) {
	// Add temporary file.
	licenseFile := "/tmp/test-file.txt"
	licenceFileContent := []byte("This is the key. Must ignore other lines.\nMust be ignored.\nMust be ignored.\n")

	err := ioutil.WriteFile(licenseFile, licenceFileContent, 0400)
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Run("Should return valid license key", func(t *testing.T) {
		actualLicenseKey, err := ReadLicense(licenseFile)

		if err != nil && bytes.Compare(licenceFileContent, []byte(actualLicenseKey)) != 0 {
			t.Errorf("readContents() error = %v", err)
		}
	})

	//Remove temporary file.
	err = os.Remove(licenseFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}
