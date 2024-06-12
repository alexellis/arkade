// Copyright (c) arkade author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package apps

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_buildYAML_Issuer(t *testing.T) {
	templBytes, _ := buildIssuerYAML("openfaas@subdomain.example.com", "traefik", false, false, "openfaas")
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
		t.Errorf("want:\n%q\n\ngot:\n%q\n", want, got)
	}
}

func Test_buildYAML_IssuerTakesEmailOverride(t *testing.T) {
	templBytes, _ := buildIssuerYAML("openfaas@subdomain.example.com", "traefik", false, false, "openfaas")
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
	templBytes, _ := buildIssuerYAML("openfaas@subdomain.example.com", "traefik", true, false, "openfaas")
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

func Test_buildIngressYAML(t *testing.T) {
	cases := []struct {
		name          string
		domain        string
		email         string
		ingressClass  string
		ingressName   string
		staging       bool
		clusterIssuer bool
		issuerName    string
		namespace     string
		hasNetworking bool
		want          string
	}{
		{
			name:          "build staging extensions/v1",
			domain:        "openfaas.subdomain.example.com",
			email:         "openfaas@subdomain.example.com",
			ingressClass:  "traefik",
			ingressName:   "openfaas-gateway",
			staging:       true,
			clusterIssuer: false,
			issuerName:    "",
			namespace:     "openfaas",
			hasNetworking: false,
			want: `
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
      `,
		},
		{
			name:          "build staging networking/v1",
			domain:        "openfaas.subdomain.example.com",
			email:         "openfaas@subdomain.example.com",
			ingressClass:  "traefik",
			ingressName:   "openfaas-gateway",
			staging:       true,
			clusterIssuer: false,
			issuerName:    "",
			namespace:     "openfaas",
			hasNetworking: true,
			want: `
apiVersion: networking.k8s.io/v1
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
      - path: /
        pathType: ImplementationSpecific
        backend:
          service:
            name: gateway
            port:
              number: 8080
  tls:
  - hosts:
    - openfaas.subdomain.example.com
    secretName: openfaas-gateway
      `,
		},
		{
			name:          "build with custom issuer",
			domain:        "openfaas.example.com",
			email:         "openfaas@example.com",
			ingressClass:  "traefik",
			ingressName:   "openfaas-gateway",
			staging:       true,
			clusterIssuer: false,
			issuerName:    "venafi-tpp",
			namespace:     "openfaas",
			hasNetworking: false,
			want: `
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
  `,
		},
		{
			name:          "build custom issuer networking/v1",
			domain:        "openfaas.subdomain.example.com",
			email:         "openfaas@subdomain.example.com",
			ingressClass:  "traefik",
			ingressName:   "openfaas-gateway",
			staging:       true,
			clusterIssuer: false,
			issuerName:    "awesome-issuer",
			namespace:     "openfaas",
			hasNetworking: true,
			want: `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: openfaas-gateway
  namespace: openfaas
  annotations:
    cert-manager.io/issuer: awesome-issuer
    kubernetes.io/ingress.class: traefik
    cert-manager.io/common-name: openfaas.subdomain.example.com
spec:
  rules:
  - host: openfaas.subdomain.example.com
    http:
      paths:
      - path: /
        pathType: ImplementationSpecific
        backend:
          service:
            name: gateway
            port:
              number: 8080
  tls:
  - hosts:
    - openfaas.subdomain.example.com
    secretName: openfaas-gateway
      `,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := buildOpenfaasIngressYAML(
				tc.domain,
				tc.email,
				tc.ingressClass,
				tc.ingressName,
				tc.staging,
				tc.clusterIssuer,
				tc.issuerName,
				tc.namespace,
				tc.hasNetworking,
			)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got := string(result)
			if strings.TrimSpace(tc.want) != strings.TrimSpace(got) {
				t.Errorf("want:\n%q\ngot:\n%q\n", tc.want, got)
			}

		})
	}

}

func Test_buildYAMLClusterIssuer_HasNoNamespace(t *testing.T) {
	templBytes, _ := buildIssuerYAML("openfaas@subdomain.example.com", "traefik", true, true, "openfaas")
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

	got, _ := os.ReadFile(tmpLocation)
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
