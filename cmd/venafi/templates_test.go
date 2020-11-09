// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package venafi

import (
	"testing"
)

func TestCloudIssuerTemplate_NamespacedIssuer(t *testing.T) {
	manifest, err := templateManifest(cloudIssuerTemplate, struct {
		Name      string
		Namespace string
		Zone      string
		Kind      string
	}{
		Name:      "cloud",
		Namespace: "default",
		Zone:      "dev",
		Kind:      "Issuer",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Uncomment to capture formatted results
	// ioutil.WriteFile("/tmp/test.yaml", manifest, os.ModePerm)
	want := `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: cloud
  namespace: default
spec:
  venafi:
    zone: "dev"
    cloud:
      apiTokenSecretRef:
        name: cloud-secret
        key: apikey
`
	got := string(manifest)
	if got != want {
		t.Errorf(`want
%q
but got
%q`, want, got)
	}
}

func TestCloudIssuerTemplate_ClusterIssuer(t *testing.T) {
	manifest, err := templateManifest(cloudIssuerTemplate, struct {
		Name      string
		Namespace string
		Zone      string
		Kind      string
	}{
		Name:      "cloud",
		Namespace: "default",
		Zone:      "dev",
		Kind:      "ClusterIssuer",
	})

	if err != nil {
		t.Fatal(err)
	}

	// Uncomment to capture formatted results
	// ioutil.WriteFile("/tmp/test.yaml", manifest, os.ModePerm)
	want := `apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: cloud
spec:
  venafi:
    zone: "dev"
    cloud:
      apiTokenSecretRef:
        name: cloud-secret
        key: apikey
`
	got := string(manifest)
	if got != want {
		t.Errorf(`want
%q
but got
%q`, want, got)
	}
}
