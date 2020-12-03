// Copyright (c) arkade author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_buildYAML_Issuer(t *testing.T) {
	templBytes, _ := buildIssuerYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", "openfaas-gateway", false, false, "openfaas")
	var want = `
apiVersion: cert-manager.io/v1
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
    - selector: {}
      http01:
        ingress:
          class: traefik`

	got := string(templBytes)
	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}

func Test_buildYAML_IssuerTakesEmailOverride(t *testing.T) {
	templBytes, _ := buildIssuerYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", "openfaas-gateway", false, false, "openfaas")
	var want = `
apiVersion: cert-manager.io/v1
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
    - selector: {}
      http01:
        ingress:
          class: traefik`

	got := string(templBytes)
	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}

func Test_buildIssuerYAMLStaging(t *testing.T) {
	templBytes, _ := buildIssuerYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", "openfaas-gateway", true, false, "openfaas")
	var want = `
apiVersion: cert-manager.io/v1
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
    - selector: {}
      http01:
        ingress:
          class: traefik`

	got := string(templBytes)
	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}

func Test_buildIngressYAMLStaging(t *testing.T) {
	templBytes, _ := buildOpenfaasIngressYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", "openfaas-gateway", true, false, "", "openfaas")
	var want = `
apiVersion: extensions/v1beta1 
kind: Ingress
metadata:
  name: openfaas-gateway
  namespace: openfaas
  annotations:
    cert-manager.io/issuer: letsencrypt-staging
    kubernetes.io/ingress.class: traefik
    cert-manager.io/common-name: openfaas.subdomain.example.com
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
`

	got := string(templBytes)
	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}

func Test_buildIngress_WithCustomIssuername(t *testing.T) {
	templBytes, _ := buildOpenfaasIngressYAML("openfaas.example.com", "openfaas@subdomain.example.com", "traefik", "openfaas-gateway", true, false, "venafi-tpp", "openfaas")
	var want = `
apiVersion: extensions/v1beta1 
kind: Ingress
metadata:
  name: openfaas-gateway
  namespace: openfaas
  annotations:
    cert-manager.io/issuer: venafi-tpp
    kubernetes.io/ingress.class: traefik
    cert-manager.io/common-name: openfaas.example.com
spec:
  rules:
  - host: openfaas.example.com
    http:
      paths:
      - backend:
          serviceName: gateway
          servicePort: 8080
        path: /
  tls:
  - hosts:
    - openfaas.example.com
    secretName: openfaas-gateway
`

	got := string(templBytes)
	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}

func Test_buildYAMLClusterIssuer_HasNoNamespace(t *testing.T) {
	templBytes, _ := buildIssuerYAML("openfaas.subdomain.example.com", "openfaas@subdomain.example.com", "traefik", "openfaas-gateway", true, true, "openfaas")
	var want = `
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    email: openfaas@subdomain.example.com
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: example-issuer-account-key
    solvers:
    - selector: {}
      http01:
        ingress:
          class: traefik`

	got := string(templBytes)
	if want != got {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}
func Test_writeTempFile_writes_to_tmp(t *testing.T) {
	var want = "some input string"
	tmpLocation, _ := writeTempFile([]byte(want), "tmp_file_name.yaml")

	got, _ := ioutil.ReadFile(tmpLocation)
	if string(got) != want {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}

func Test_createTempDirectory_creates(t *testing.T) {
	var want = filepath.Join(os.TempDir(), ".arkade")

	got, _ := createTempDirectory(".arkade")

	if got != want {
		t.Errorf("want:\n%q\ngot:\n%q\n", want, got)
	}
}
