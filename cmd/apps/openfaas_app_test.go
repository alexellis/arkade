// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"testing"
)

func Test_getValuesSuffix_arm64(t *testing.T) {
	want := "-arm64"
	got := getValuesSuffix("arm64")
	if want != got {
		t.Errorf("suffix, want: %s, got: %s", want, got)
	}
}

func Test_getValuesSuffix_aarch64(t *testing.T) {
	want := "-arm64"
	got := getValuesSuffix("aarch64")
	if want != got {
		t.Errorf("suffix, want: %s, got: %s", want, got)
	}
}
