# arkade - get Kubernetes apps, the easy way

Gone are the days of contending with dozens of README files just to get the right version of helm and to install a chart with sane defaults. arkade (ark for short) provides a clean CLI with strongly-typed flags to install charts and apps to your cluster in one command.

[![Build
Status](https://travis-ci.com/alexellis/arkade.svg?branch=master)](https://travis-ci.com/alexellis/arkade)
[![GoDoc](https://godoc.org/github.com/alexellis/arkade?status.svg)](https://godoc.org/github.com/alexellis/arkade) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub All Releases](https://img.shields.io/github/downloads/alexellis/arkade/total)

## What about helm and `k3sup`?

In the same way that brew uses git and Makefiles to compile applications for your Mac, `arkade` uses upstream helm charts and kubectl to install applications to your Kubernetes cluster.

On k3sup vs. arkade: The codebase in this project is derived from `k3sup`. k3sup (ketchup) was developed to automate building of k3s clusters over SSH, then gained the powerful feature to install apps in a single command. The presence of the word "k3s" in the name of the application confused many people, this is why arkade has come to exist.

And yes, of course it works with k3s and where possible, apps are available for ARM.

## Get arkade

```bash
curl -sLS https://dl.get-arkade.dev | sh
sudo install arkade /usr/local/bin/

arkade --help
```

An alias of `ark` is created at installation time.

## Usage

Here's a few examples of apps you can install, for a complete list run: `[ark]ade install --help`.

```bash
[ark]ade install openfaas --gateways 2 --load-balancer false

[ark]ade install cert-manager

[ark]ade install nginx-ingress

[ark]ade install inlets-operator --access-token $HOME/digitalocean --region lon1
```

Here's how you can get a self-hosted Docker registry with TLS and authentication in just 5 commands on an empty cluster:

```bash
arkade install nginx-ingress
arkade install cert-manager
arkade install docker-registry
arkade install docker-registry-ingress \
  --email web@example.com \
  --domain reg.example.com
```

The same for OpenFaaS would look like this:

```bash
arkade install nginx-ingress
arkade install cert-manager
arkade install openfaas
arkade install openfaas-ingress \
  --email web@example.com \
  --domain reg.example.com
```

And if you're running on a private cloud, on-premises or on your laptop, you can simply add the inlets-operator using inlets-pro to get a secure TCP tunnel and a public IP address.

```bash
[ark]ade install inlets-operator \
  --access-token $HOME/digitalocean \
  --region lon1 \
  --license $(cat $HOME/license.txt)
```

## Contributing

### Suggesting a new app

To suggest a new app, please check past issues and [raise an issue for it](https://github.com/alexellis/arkade).

### Improving the code or fixing an issue

Before contributing code, please see the [CONTRIBUTING guide](https://github.com/alexellis/inlets/blob/master/CONTRIBUTING.md). Note that arkade uses the same guide as [inlets.dev](https://inlets.dev/).

Both Issues and PRs have their own templates. Please fill out the whole template.

All commits must be signed-off as part of the [Developer Certificate of Origin (DCO)](https://developercertificate.org)

### License

MIT

