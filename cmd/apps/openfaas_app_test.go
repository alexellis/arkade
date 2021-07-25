// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"reflect"
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

func Test_mergeFlags(t *testing.T) {
	existingMap := map[string]string{}

	want := map[string]string{
		"httpRetryCodes":    `429\,502\,500\,504\,408`,
		"queueWorker.image": "alexellis2/pro-queue-worker-demo:0.1.1",
	}

	mergeFlags(existingMap, []string{"queueWorker.image=alexellis2/pro-queue-worker-demo:0.1.1", "httpRetryCodes=429,502,500,504,408"})

	eq := reflect.DeepEqual(existingMap, want)

	if !eq {
		t.Errorf("want: %s, got: %s", want, existingMap)
	}

}
