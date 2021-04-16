// Copyright (c) arkade author(s) 2021. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_ProviderHetzner(t *testing.T) {
	cmdInlets := MakeInstallInletsOperator()

	// missing required flags
	err := cmdInlets.ParseFlags([]string{"--provider", "hetzner"})
	if err != nil {
		t.Errorf("expected no error when parsing flags, but got: %s", err.Error())
	}

	errWant := "invalid region set for provider hetzner. Valid regions are: fsn1, nbg1, hel1"
	_, err = getInletsOperatorOverrides(cmdInlets)
	if err == nil {
		t.Errorf("expected error when missing a flag")
	}
	if err.Error() != errWant {
		t.Errorf("expected error: %s, got: %s", errWant, err)
	}

	// required flags in place
	err = cmdInlets.ParseFlags([]string{"--provider", "hetzner", "--region", "fsn1"})
	if err != nil {
		t.Errorf("expected no error when parsing flags, but got: %s", err.Error())
	}

	overridesWant := map[string]string{"provider": "hetzner", "region": "fsn1"}
	overrides, err := getInletsOperatorOverrides(cmdInlets)
	if err != nil {
		t.Errorf("expected no error when parsing flags, but got: %s", err.Error())
	}
	if !reflect.DeepEqual(overrides, overridesWant) {
		t.Errorf("expected %v, but got: %v", overridesWant, overrides)
	}
}

func Test_ValidatePreRun(t *testing.T) {
	err := validatePreRun("my_honk_token", "", "")
	if err != nil {
		t.Errorf("expected no error when validating the PreRun, but got: %s", err.Error())
	}

	tokenFile, err := ioutil.TempFile("", "honk-token-file")
	if err != nil {
		t.Errorf("expected no error when creating the temporary test file, but got: %s", err.Error())
	}
	defer os.Remove(tokenFile.Name())

	err = validatePreRun("", tokenFile.Name(), "")
	if err != nil {
		t.Errorf("expected no error when validating the PreRun, but got: %s", err.Error())
	}

	secretFile, err := ioutil.TempFile("", "honk-secret-file")
	if err != nil {
		t.Errorf("expected no error when creating the temporary test file, but got: %s", err.Error())
	}
	defer os.Remove(secretFile.Name())

	err = validatePreRun("", tokenFile.Name(), secretFile.Name())
	if err != nil {
		t.Errorf("expected no error when validating the PreRun, but got: %s", err.Error())
	}

	// Invalid TokenFileName
	errWant := "failed to check the token file invalid-token-file: stat invalid-token-file: no such file or directory"
	err = validatePreRun("", "invalid-token-file", "")
	if err == nil {
		t.Errorf("expected error when the token file does not exist")
	}
	if err.Error() != errWant {
		t.Errorf("expected error: %s, got: %s", errWant, err)
	}

	errWant = "failed to check the secret file invalid-secret-file: stat invalid-secret-file: no such file or directory"
	err = validatePreRun("my_honk_token", "", "invalid-secret-file")
	if err == nil {
		t.Errorf("expected error when the secret file does not exist")
	}
	if err.Error() != errWant {
		t.Errorf("expected error: %s, got: %s", errWant, err)
	}

	errWant = "--token-file or --token is a required field for your cloud API token or service account JSON"
	err = validatePreRun("", "", "")
	if err == nil {
		t.Errorf("expected error when the secret file does not exist")
	}
	if err.Error() != errWant {
		t.Errorf("expected error: %s, got: %s", errWant, err)
	}
}
