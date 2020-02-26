# bazaar - get Kubernetes apps, the easy way

Gone are the days of contending with dozens of README files just to get the right version of helm and to install a chart with sane defaults. bazaar (baz for short) provides a clean CLI with strongly-typed flags to install charts and apps to your cluster in one command.

[![Build
Status](https://travis-ci.com/alexellis/bazaar.svg?branch=master)](https://travis-ci.com/alexellis/bazaar)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexellis/bazaar)](https://goreportcard.com/report/github.com/alexellis/bazaar) 
[![GoDoc](https://godoc.org/github.com/alexellis/bazaar?status.svg)](https://godoc.org/github.com/alexellis/bazaar) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub All Releases](https://img.shields.io/github/downloads/alexellis/bazaar/total)

## Usage

Here's a few examples of apps you can install, for a complete list run: `[baz]aar install --help`.

```bash
[baz]aar install openfaas --gateways 2 --load-balancer false

[baz]aar install cert-manager

[baz]aar install nginx-ingress

[baz]aar install inlets-operator --access-token $HOME/digitalocean --region lon1
```

An alias of `baz` is created at installation time.

Here's how you can get a self-hosted Docker registry with TLS and authentication in just 5 commands on an empty cluster:

```bash
bazaar install nginx-ingress
bazaar install cert-manager
bazaar install docker-registry
bazaar install docker-registry-ingress \
  --email web@example.com \
  --domain reg.example.com
```

The same for OpenFaaS would look like this:

```bash
bazaar install nginx-ingress
bazaar install cert-manager
bazaar install openfaas
bazaar install openfaas-ingress \
  --email web@example.com \
  --domain reg.example.com
```

And if you're running on a private cloud, on-premises or on your laptop, you can simply add the inlets-operator using inlets-pro to get a secure TCP tunnel and a public IP address.

```bash
[baz]aar install inlets-operator \
  --access-token $HOME/digitalocean \
  --region lon1 \
  --license $(cat $HOME/license.txt)
```

## Contributing

### Suggesting a new app

To suggest a new app, please check past issues and [raise an issue for it](https://github.com/alexellis/bazaar).

### Improving the code or fixing an issue

Before contributing code, please see the [CONTRIBUTING guide](https://github.com/alexellis/inlets/blob/master/CONTRIBUTING.md). Note that bazaar uses the same guide as [inlets.dev](https://inlets.dev/).

Both Issues and PRs have their own templates. Please fill out the whole template.

All commits must be signed-off as part of the [Developer Certificate of Origin (DCO)](https://developercertificate.org)

### k3sup vs. bazaar

The codebase in this project is derived from `k3sup`. k3sup (ketchup) was developed to automate building of k3s clusters over SSH, then gained the powerful feature to install apps in a single command. The presence of the word "k3s" in the name of the application confused many people, this is why bazaar has come to exist.

### License

MIT

