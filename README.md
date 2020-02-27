# arkade - get Kubernetes apps, the easy way

arkade (ark for short) provides a clean CLI with strongly-typed flags to install charts and apps to your cluster in one command.

<img src="docs/arkade-logo-sm.png" alt="arkade logo" width="150" height="150">

Gone are the days of contending with dozens of README files just to get the right version of [helm](https://helm.sh) and to install a chart with sane defaults.

[![Build
Status](https://travis-ci.com/alexellis/arkade.svg?branch=master)](https://travis-ci.com/alexellis/arkade)
[![GoDoc](https://godoc.org/github.com/alexellis/arkade?status.svg)](https://godoc.org/github.com/alexellis/arkade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexellis/arkade)](https://goreportcard.com/report/github.com/alexellis/arkade)
![GitHub All Releases](https://img.shields.io/github/downloads/alexellis/arkade/total)

## Get arkade

```bash
# Note: you can also run without `sudo` and move the binary yourself
curl -sLS https://dl.get-arkade.dev | sudo sh

arkade --help
ark --help  # a handy alias
```

## Usage

Here's a few examples of apps you can install, for a complete list run: `arkade install --help`.

* `arkade install` - install an app
* `arkade update` - update arkade
* `arkade info` - the post-install screen for an app

### Install an app

No need to worry about whether you're installing to Intel or ARM architecture, the correct values will be set for you automatically.

```bash
arkade install openfaas --gateways 2 --load-balancer false
```

#### Reduce the repetition

[Normally up to a dozen commands](https://cert-manager.io/docs/installation/kubernetes/) (including finding and downloading helm), now just one. No searching for the correct CRD to apply, no trying to install helm, no trying to find the correct helm repo to add:

```bash
arkade install cert-manager
```

Other common tools:

```bash
arkade install nginx-ingress

arkade install metrics-server
```

#### Bye-bye values.yaml, hello flags

We use strongly typed Go CLI flags, so that you can run `--help` instead of trawling through countless Helm chart README files to find the correct `--set` combination for what you want.

```bash
arkade install nginx-ingress --help

Install nginx-ingress. This app can be installed with Host networking for 
cases where an external LB is not available. please see the --host-mode 
flag and the nginx-ingress docs for more info

Usage:
  arkade install nginx-ingress [flags]

Examples:
  arkade install nginx-ingress --namespace default

Flags:
      --helm3              Use helm3, if set to false uses helm2 (default true)
  -h, --help               help for nginx-ingress
      --host-mode          If we should install nginx-ingress in host mode.
  -n, --namespace string   The namespace used for installation (default "default")
      --update-repo        Update the helm repo (default true)
```

#### Override with `--set`

You can also set helm overrides, for apps which use helm via `--set`

```bash
ark install openfaas --set=faasIdler.dryRun=false
```

After installation, an info message will be printed with help for usage, you can get back to this at any time via:

```bash
arkade info <NAME>
```

#### Get a self-hosted TLS registry with authentication

Here's how you can get a self-hosted Docker registry with TLS and authentication in just 5 commands on an empty cluster:

```bash
arkade install nginx-ingress
arkade install cert-manager
arkade install docker-registry
arkade install docker-registry-ingress \
  --email web@example.com \
  --domain reg.example.com
```

#### Get OpenFaaS with TLS

The same for OpenFaaS would look like this:

```bash
arkade install nginx-ingress
arkade install cert-manager
arkade install openfaas
arkade install openfaas-ingress \
  --email web@example.com \
  --domain reg.example.com
```

#### Get a public IP for a private cluster and your IngressController

And if you're running on a private cloud, on-premises or on your laptop, you can simply add the [inlets-operator](https://github.com/inlets/inlets-operator/) using [inlets-pro](https://docs.inlets.dev/) to get a secure TCP tunnel and a public IP address.

```bash
arkade install inlets-operator \
  --access-token $HOME/digitalocean-token \
  --region lon1 \
  --license $(cat $HOME/license.txt)
```

#### Explore the apps

```bash
arkade install --help
ark --help

cert-manager            Install cert-manager
chart                   Install the specified helm chart
cron-connector          Install cron-connector for OpenFaaS
crossplane              Install Crossplane
docker-registry         Install a Docker registry
docker-registry-ingress Install registry ingress with TLS
info                    Find info about a Kubernetes app
inlets-operator         Install inlets-operator
istio                   Install istio
kafka-connector         Install kafka-connector for OpenFaaS
kubernetes-dashboard    Install kubernetes-dashboard
linkerd                 Install linkerd
metrics-server          Install metrics-server
minio                   Install minio
mongodb                 Install mongodb
nginx-ingress           Install nginx-ingress
openfaas                Install openfaas
openfaas-ingress        Install openfaas ingress with TLS
postgresql              Install postgresql
```

## Community & contributing

### What about helm and `k3sup`?

In the same way that brew uses git and Makefiles to compile applications for your Mac, `arkade` uses upstream [helm](https://helm.sh) charts and kubectl to install applications to your Kubernetes cluster.

On k3sup vs. arkade: The codebase in this project is derived from `k3sup`. [k3sup (ketchup)](https://k3sup.dev/) was developed to automate building of k3s clusters over SSH, then gained the powerful feature to install apps in a single command. The presence of the word "k3s" in the name of the application confused many people, this is why arkade has come to exist.

And yes, of course it works with k3s and where possible, apps are available for ARM.

### Tools and cached versions of helm

When required, tools, CLIs, and the helm binaries are downloaded and extracted to `$HOME/.arkade`.

If installing a tool which uses helm3, arkade will check for a cached version and use that, otherwise it will download it on demand.

### Suggesting a new app

To suggest a new app, please check past issues and [raise an issue for it](https://github.com/alexellis/arkade).

### Improving the code or fixing an issue

Before contributing code, please see the [CONTRIBUTING guide](https://github.com/alexellis/inlets/blob/master/CONTRIBUTING.md). Note that arkade uses the same guide as [inlets.dev](https://inlets.dev/).

Both Issues and PRs have their own templates. Please fill out the whole template.

All commits must be signed-off as part of the [Developer Certificate of Origin (DCO)](https://developercertificate.org)

### Join us on Slack

Join #arkade on [slack.openfaas.io](https://slack.openfaas.io)

### License

MIT

