# arkade - Your one-stop CLI for Kubernetes

arkade provides a simple Golang CLI with strongly-typed flags to install charts and apps to your cluster in one command.

You can also download CLIs like `kubectl`, `kind`, `kubectx` and `helm` faster than you can type "apt-get/brew update".

<img src="docs/arkade-logo-sm.png" alt="arkade logo" width="150" height="150">

[![Build Status](https://travis-ci.com/alexellis/arkade.svg?branch=master)](https://travis-ci.com/alexellis/arkade)
[![GoDoc](https://godoc.org/github.com/alexellis/arkade?status.svg)](https://godoc.org/github.com/alexellis/arkade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexellis/arkade)](https://goreportcard.com/report/github.com/alexellis/arkade)
![GitHub All Releases](https://img.shields.io/github/downloads/alexellis/arkade/total)

Gone are the days of contending with dozens of README files just to get the right version of [helm](https://helm.sh) and to install a chart with sane defaults.

## Get arkade

```bash
# Note: you can also run without `sudo` and move the binary yourself
curl -sLS https://dl.get-arkade.dev | sudo sh

arkade --help
ark --help  # a handy alias

# Windows users with Git Bash
curl -sLS https://dl.get-arkade.dev | sh
```

> Windows users: arkade requires bash to be available, therefore Windows users can [install Git Bash](https://git-scm.com/downloads).

An alias of `ark` is created at installation time, so you can also run `ark install APP`

<img src="https://cdn.shopify.com/s/files/1/0340/9739/7895/products/mockup-ea71d220_1024x1024@2x.jpg?v=1589366177" width="250">

[You can buy an arkade t-shirt in the OpenFaaS Ltd store](https://store.openfaas.com/products/arkade-t-shirt)

## Usage

Here's a few examples of apps you can install, for a complete list run: `arkade install --help`.

* `arkade install` - install an app
* `arkade info` - the post-install screen for an app
* `arkade get` - install a CLI tool such as `kubectl` or `faas-cli`
* `arkade update` - print instructions to update arkade itself

### Install a CLI tool

arkade downloads the correct version of a CLI for your OS and CPU.

With automatic detection of: Windows / MacOS / Linux / Intel / ARM.

```bash
arkade get faas-cli
arkade get helm
arkade get inletsctl
arkade get k3d
arkade get k3sup
arkade get kind
arkade get kubectl
arkade get kubectx
arkade get kubeseal
```

> This is a time saver compared to searching for download pages every time you need a tool.

Think of `arkade get TOOL` as a doing for CLIs, what `arkade install` does for helm.

### Create a Kubernetes cluster

If you have Docker installed, then you can install Kubernetes using KinD in a matter of moments:

```bash
arkade get kubectl
arkade get kind

kind create cluster
```

You can also download k3d [k3s](https://github.com/rancher/k3s) in the same way with `arkade get k3d`.

### Install an app

No need to worry about whether you're installing to Intel or ARM architecture, the correct values will be set for you automatically.

```bash
arkade install openfaas --gateways 2 --load-balancer false
```

Remember how awkward it was last time you installed the [Kubernetes dashboard](https://github.com/kubernetes/dashboard)? And how you could never remember the command to get the token to log in?

```bash
arkade install kubernetes-dashboard
```

Forgot your token? `arkade info kubernetes-dashboard`

Prefer [Portainer](https://www.portainer.io)? Just run: `arkade install portainer`

### Uninstall an app

Run `arkade uninstall` or `arkade delete` for more information on how to remove applications from a Kubernetes cluster.

#### Reduce the repetition

[Normally up to a dozen commands](https://cert-manager.io/docs/installation/kubernetes/) (including finding and downloading helm), now just one. No searching for the correct CRD to apply, no trying to install helm, no trying to find the correct helm repo to add:

```bash
arkade install cert-manager
```

Other common tools:

```bash
arkade install ingress-nginx

arkade install metrics-server
```

#### Bye-bye values.yaml, hello flags

We use strongly typed Go CLI flags, so that you can run `--help` instead of trawling through countless Helm chart README files to find the correct `--set` combination for what you want.

```bash
arkade install ingress-nginx --help

Install ingress-nginx. This app can be installed with Host networking for
cases where an external LB is not available. please see the --host-mode
flag and the ingress-nginx docs for more info

Usage:
  arkade install ingress-nginx [flags]

Aliases:
  ingress-nginx, nginx-ingress

Examples:
  arkade install ingress-nginx --namespace default

Flags:
      --helm3              Use helm3, if set to false uses helm2 (default true)
  -h, --help               help for ingress-nginx
      --host-mode          If we should install ingress-nginx in host mode.
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
arkade install ingress-nginx
arkade install cert-manager
arkade install docker-registry
arkade install docker-registry-ingress \
  --email web@example.com \
  --domain reg.example.com
```

#### Get OpenFaaS with TLS

The same for OpenFaaS would look like this:

```bash
arkade install ingress-nginx
arkade install cert-manager
arkade install openfaas
arkade install openfaas-ingress \
  --email web@example.com \
  --domain reg.example.com
```

#### Get a public IP for a private cluster and your IngressController

And if you're running on a private cloud, on-premises or on your laptop, you can simply add the [inlets-operator](https://github.com/inlets/inlets-operator/) using [inlets PRO](https://docs.inlets.dev/) to get a secure TCP tunnel and a public IP address.

```bash
arkade install inlets-operator \
  --access-token $HOME/digitalocean-token \
  --region lon1 \
  --license $(cat $HOME/license.txt)
```

### Tutorials & community blog posts

#### Video review from Rancher Labs

* [Tool of the Day with Adrian from Rancher Labs](https://youtu.be/IWtEtfpqoRg?t=1425)

#### Watch a video walk-through by Alex Ellis

[![](http://img.youtube.com/vi/8wU9s_mua8M/hqdefault.jpg)](https://www.youtube.com/watch?v=8wU9s_mua8M)

#### Featured tutorials

* [arkade by example — Kubernetes apps, the easy way 😎](https://medium.com/@alexellisuk/kubernetes-apps-the-easy-way-f06d9e5cad3c) - Alex Ellis
* [Walk-through — install Kubernetes to your Raspberry Pi in 15 minutes](https://medium.com/@alexellisuk/walk-through-install-kubernetes-to-your-raspberry-pi-in-15-minutes-84a8492dc95a)
* [Get a TLS-enabled Docker registry in 5 minutes](https://blog.alexellis.io/get-a-tls-enabled-docker-registry-in-5-minutes/) - Alex Ellis
* [Get TLS for OpenFaaS the easy way with arkade](https://blog.alexellis.io/tls-the-easy-way-with-openfaas-and-k3sup/) - Alex Ellis

#### Community posts

* [A bit of Istio before tea-time](https://blog.alexellis.io/a-bit-of-istio-before-tea-time/) - Alex Ellis
* [Kubernetes: Automatic Let's Encrypt Certificates for Services with arkade](https://medium.com/@admantium/kubernetes-automatic-lets-encrypt-certificates-for-services-2a5f4aa7f886)
* [Introducing Arkade - The Kubernetes app installer](https://blog.heyal.co.uk/introducing-arkade/) - Alistair Hey
* [Portainer for kubernetes in less than 60 seconds!!](https://www.portainer.io/2020/04/portainer-for-kubernetes-in-less-than-60-seconds/) - by Saiyam Pathak
* [Video walk-through with DJ Adams - Pi & Kubernetes with k3s, k3sup, arkade and OpenFaaS](https://www.youtube.com/watch?v=ZiR3QEfBivk)
* [Coffee chat: Easy way to install Kubernetes Apps - arkade (ark)](https://sachcode.com/tech/coffee-chat-easy-way-install-kubernetes-apps/) by Sachin Jha
* [Arkade & OpenFaaS: serverless on the spot](https://zero2datadog.readthedocs.io/en/latest/faas.html) by Blaise Pabon

#### Explore the apps

You can view the various apps available with `arkade install / --help`, more are available when you run the command yourself.

```bash
arkade install --help

  argocd                  Install argocd
  cert-manager            Install cert-manager
  chart                   Install the specified helm chart
  cron-connector          Install cron-connector for OpenFaaS
  crossplane              Install Crossplane
  docker-registry         Install a Docker registry
  docker-registry-ingress Install registry ingress with TLS
  grafana                 Install grafana
  info                    Find info about a Kubernetes app
  ingress-nginx           Install ingress-nginx
  inlets-operator         Install inlets-operator
  istio                   Install istio
  jenkins                 Install jenkins
  kafka-connector         Install kafka-connector for OpenFaaS
  kube-state-metrics      Install kube-state-metrics
  kubernetes-dashboard    Install kubernetes-dashboard
  linkerd                 Install linkerd
  loki                    Install Loki for monitoring and tracing
  metrics-server          Install metrics-server
  minio                   Install minio
  mongodb                 Install mongodb
  nats-connector          Install OpenFaaS connector for NATS
  nfs-client-provisioner  Install nfs client provisioner
  openfaas                Install openfaas
  openfaas-ingress        Install openfaas ingress with TLS
  openfaas-loki           Install Loki-OpenFaaS and Configure Loki logs provider for OpenFaaS
  portainer               Install portainer to visualise and manage containers
  postgresql              Install postgresql
  redis                   Install redis
  tekton                  Install Tekton pipelines and dashboard
  traefik2                Install traefik2
```

## Community & contributing

### Suggesting a new app

To suggest a new app, please check past issues and [raise an issue for it](https://github.com/alexellis/arkade). Think also whether your app suggestion would be a good candidate for a Sponsored App.

### Sponsored apps

You can now propose your project or product as a Sponsored App. Sponsored Apps work just like any other app that we've curated, however they will have a note next to them in the app description `(sponsored)` and a link to your chosen site upon installation. An app sponsorship can be purchased for a minimum of 12 months and includes free development and support of the app for arkade.

When your sponsorship expires your app can be renewed at that time, or it will disappear automatically based upon the end-date.

Email [sales@openfaas.com](mailto:sales@openfaas.com) to find out more.

### How does `arkade` compare to `helm`?

In the same way that [brew](https://brew.sh) uses git and Makefiles to compile applications for your Mac, `arkade` uses upstream [helm](https://helm.sh) charts and `kubectl` to install applications to your Kubernetes cluster. arkade exposes strongly-typed flags for the various popular options for helm charts, and enables easier discovery through `arkade install --help` and `arkade install APP --help`.

### Is arkade suitable for production use?

If you consider helm suitable, and `kubectl` then yes, arkade by definition uses those tools and the upstream artifacts of OSS projects.

Do you want to run arkade in a CI or CD pipeline? Go ahead.

### What is in scope for `arkade get`?

Generally speaking, tools that are used with the various arkade apps or with Kubernetes are in scope. If you want to propose a tool, raise a GitHub issue.

What about package management? `arkade get` provides a faster alternative to package managers like `apt` and `brew`, you're free to use either or both at the same time.

### Automatic download of tools

When required, tools, CLIs, and the helm binaries are downloaded and extracted to `$HOME/.arkade`.

If installing a tool which uses helm3, arkade will check for a cached version and use that, otherwise it will download it on demand.

Did you accidentally run arkade as root? **Running as root is not required**, and will mean your KUBECONFIG environment variable will be ignored. You can revert this using [the notes on release 0.1.18](https://github.com/alexellis/arkade/releases/tag/0.1.8).

### Improving the code or fixing an issue

Before contributing code, please see the [CONTRIBUTING guide](https://github.com/alexellis/inlets/blob/master/CONTRIBUTING.md). Note that arkade uses the same guide as [inlets.dev](https://inlets.dev/).

Both Issues and PRs have their own templates. Please fill out the whole template.

All commits must be signed-off as part of the [Developer Certificate of Origin (DCO)](https://developercertificate.org)

### Join us on Slack

Join `#contributors` at [https://slack.openfaas.io](https://slack.openfaas.io)

### License

MIT
