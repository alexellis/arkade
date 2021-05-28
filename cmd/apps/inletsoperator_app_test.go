// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_ValidatePreRun(t *testing.T) {
	testcases := []struct {
		name              string
		expectedError     string
		token             string
		tokenFileName     string
		secretKeyFileName string
	}{
		{
			name:              "token is set",
			expectedError:     "",
			token:             "my_honk_token",
			tokenFileName:     "",
			secretKeyFileName: "",
		},
		{
			name:              "tokenFileName is set",
			expectedError:     "",
			token:             "",
			tokenFileName:     "honk-token-file",
			secretKeyFileName: "",
		},
		{
			name:              "tokenFileName and secretKeyFileName are set",
			expectedError:     "",
			token:             "",
			tokenFileName:     "honk-token-file",
			secretKeyFileName: "honk-secret-file",
		},
		{
			name:              "missing token",
			expectedError:     "--token-file or --token is a required field for your cloud API token or service account JSON",
			token:             "",
			tokenFileName:     "",
			secretKeyFileName: "",
		},
		{
			name:              "invalid token file",
			expectedError:     "failed to check the token file invalid-token-file: stat invalid-token-file: no such file or directory",
			token:             "",
			tokenFileName:     "invalid-token-file",
			secretKeyFileName: "",
		},
		{
			name:              "invalid secret file",
			expectedError:     "failed to check the secret file invalid-secret-file: stat invalid-secret-file: no such file or directory",
			token:             "honk-token",
			tokenFileName:     "",
			secretKeyFileName: "invalid-secret-file",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			tokenTempfileName := ""
			if strings.Contains(tc.tokenFileName, "invalid") {
				tokenTempfileName = tc.tokenFileName
			} else if tc.tokenFileName != "" {
				tokenTempfileName = createTempFile(t, tc.tokenFileName)
				defer os.Remove(tokenTempfileName)
			}

			secretTempfileName := ""
			if strings.Contains(tc.secretKeyFileName, "invalid") {
				secretTempfileName = tc.secretKeyFileName
			} else if tc.secretKeyFileName != "" {
				secretTempfileName = createTempFile(t, tc.secretKeyFileName)
				defer os.Remove(secretTempfileName)
			}

			err := validatePreRun(tc.token, tokenTempfileName, secretTempfileName)
			if tc.expectedError == "" {
				if err != nil {
					t.Errorf("expected no error when validating the PreRun, but got: %s", err.Error())
				}
			} else {
				if err.Error() != tc.expectedError {
					t.Errorf("expected error: %s, got: %s", tc.expectedError, err)
				}
			}
		})
	}
}

func createTempFile(t *testing.T, name string) string {
	tempFile, err := ioutil.TempFile("", name)
	if err != nil {
		t.Errorf("expected no error when creating the temporary test file, but got: %s", err.Error())
	}

	return tempFile.Name()
}
