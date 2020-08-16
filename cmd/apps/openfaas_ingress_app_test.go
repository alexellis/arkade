// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_buildYAML_SubsitutesDomainEmailAndIngress(t *testing.T) {
	templBytes, _ := buildYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", false, false)
	var want = `
apiVersion: extensions/v1beta1 
kind: Ingress
metadata:
  name: openfaas-gateway
  namespace: openfaas
  annotations:
    cert-manager.io/issuer: letsencrypt-prod
    kubernetes.io/ingress.class: traefik
spec:
  rules:
  - host: openfaas.subdomain.example.com
    http:
      paths:
      - backend:
          serviceName: gateway
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - openfaas.subdomain.example.com
    secretName: openfaas-gateway
---
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: letsencrypt-prod
  namespace: openfaas
spec:
  acme:
    email: openfaas@subdomain.example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: example-issuer-account-key
    solvers:
    - http01:
        ingress:
          class: traefik`

	got := string(templBytes)
	if want != got {
		t.Errorf("suffix, want: %q, got: %q", want, got)
	}
}

func Test_buildYAMLStaging(t *testing.T) {
	templBytes, _ := buildYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", true, false)
	var want = `
apiVersion: extensions/v1beta1 
kind: Ingress
metadata:
  name: openfaas-gateway
  namespace: openfaas
  annotations:
    cert-manager.io/issuer: letsencrypt-staging
    kubernetes.io/ingress.class: traefik
spec:
  rules:
  - host: openfaas.subdomain.example.com
    http:
      paths:
      - backend:
          serviceName: gateway
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - openfaas.subdomain.example.com
    secretName: openfaas-gateway
---
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: letsencrypt-staging
  namespace: openfaas
spec:
  acme:
    email: openfaas@subdomain.example.com
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: example-issuer-account-key
    solvers:
    - http01:
        ingress:
          class: traefik`

	got := string(templBytes)
	if want != got {
		t.Errorf("suffix, want: %q, got: %q", want, got)
	}
}

func Test_buildYAMLClusterIssuer(t *testing.T) {
	templBytes, _ := buildYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", false, true)
	var want = `
apiVersion: extensions/v1beta1 
kind: Ingress
metadata:
  name: openfaas-gateway
  namespace: openfaas
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    kubernetes.io/ingress.class: traefik
spec:
  rules:
  - host: openfaas.subdomain.example.com
    http:
      paths:
      - backend:
          serviceName: gateway
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - openfaas.subdomain.example.com
    secretName: openfaas-gateway
---
apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    email: openfaas@subdomain.example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: example-issuer-account-key
    solvers:
    - http01:
        ingress:
          class: traefik`

	got := string(templBytes)
	if want != got {
		t.Errorf("suffix, want: %q, got: %q", want, got)
	}
}

func Test_writeTempFile_writes_to_tmp(t *testing.T) {
	var want = "some input string"
	tmpLocation, _ := writeTempFile([]byte(want), "tmp_file_name.yaml")

	got, _ := ioutil.ReadFile(tmpLocation)
	if string(got) != want {
		t.Errorf("suffix, want: %q, got: %q", want, got)
	}
}

func Test_createTempDirectory_creates(t *testing.T) {
	var want = filepath.Join(os.TempDir(), ".arkade")

	got, _ := createTempDirectory(".arkade")

	if got != want {
		t.Errorf("suffix, want: %q, got: %q", want, got)
	}
}
