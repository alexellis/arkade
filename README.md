# arkade - Open Source Marketplace For Developer Tools

arkade is how developers install the latest versions of their favourite CLI tools and Kubernetes apps.

With `arkade get`, you'll have `kubectl`, `kind`, `terraform`, and `jq` on your machine faster than you can type `apt-get install` or `brew update`.

<img src="docs/arkade-logo-sm.png" alt="arkade logo" width="150" height="150">

[![Sponsor this](https://img.shields.io/static/v1?label=Sponsor&message=%E2%9D%A4&logo=GitHub&link=https://github.com/sponsors/alexellis)](https://github.com/sponsors/alexellis) [![CI Build](https://github.com/alexellis/arkade/actions/workflows/build.yml/badge.svg)](https://github.com/alexellis/arkade/actions/workflows/build.yml)
[![URL Checker](https://github.com/alexellis/arkade/actions/workflows/e2e-url-checker.yml/badge.svg)](https://github.com/alexellis/arkade/actions/workflows/e2e-url-checker.yml)
[![GoDoc](https://godoc.org/github.com/alexellis/arkade?status.svg)](https://godoc.org/github.com/alexellis/arkade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Downloads](https://img.shields.io/github/downloads/alexellis/arkade/total)

With over 120 CLIs and 55 Kubernetes apps (charts, manifests, installers) available for Kubernetes, gone are the days of contending with dozens of README files just to set up a development stack with the usual suspects like ingress-nginx, Postgres, and cert-manager.

- [arkade - Open Source Marketplace For Developer Tools](#arkade---open-source-marketplace-for-developer-tools)
  - [Support arkade 👋](#support-arkade-)
  - [Should you try arkade?](#should-you-try-arkade)
  - [Getting arkade](#getting-arkade)
  - [Usage overview](#usage-overview)
  - [Download CLI tools with arkade](#download-cli-tools-with-arkade)
  - [Install System Packages](#install-system-packages)
  - [Install packages from OCI images](#install-packages-from-oci-images)
  - [Install CLIs during CI with GitHub Actions](#install-clis-during-ci-with-github-actions)
  - [Bump Helm chart versions](#bump-helm-chart-versions)
  - [Verify and upgrade images in Helm charts](#verify-and-upgrade-images-in-helm-charts)
    - [Upgrade images within a Helm chart](#upgrade-images-within-a-helm-chart)
  - [Verify images within a helm chart](#verify-images-within-a-helm-chart)
  - [Installing apps with arkade](#installing-apps-with-arkade)
    - [Create a Kubernetes cluster](#create-a-kubernetes-cluster)
    - [Install a Kubernetes app](#install-a-kubernetes-app)
    - [Uninstall an app](#uninstall-an-app)
    - [Reduce the repetition](#reduce-the-repetition)
    - [Say goodbye to values.yaml and hello to flags](#say-goodbye-to-valuesyaml-and-hello-to-flags)
    - [Override with `--set`](#override-with---set)
    - [Compounding apps](#compounding-apps)
      - [Get a self-hosted TLS registry with authentication](#get-a-self-hosted-tls-registry-with-authentication)
      - [Get a public IP for a private cluster and your IngressController](#get-a-public-ip-for-a-private-cluster-and-your-ingresscontroller)
    - [Explore the apps](#explore-the-apps)
  - [Community \& contributing](#community--contributing)
    - [Tutorials \& community blog posts](#tutorials--community-blog-posts)
      - [Watch a video walk-through by Alex Ellis](#watch-a-video-walk-through-by-alex-ellis)
      - [Featured tutorials](#featured-tutorials)
      - [Official blog posts](#official-blog-posts)
      - [Community posts](#community-posts)
    - [Suggest a new app](#suggest-a-new-app)
  - [Sponsored apps](#sponsored-apps)
  - [FAQ](#faq)
    - [How does `arkade` compare to `helm`?](#how-does-arkade-compare-to-helm)
    - [Is arkade suitable for production use?](#is-arkade-suitable-for-production-use)
    - [What is in scope for `arkade get`?](#what-is-in-scope-for-arkade-get)
    - [Automatic download of tools](#automatic-download-of-tools)
    - [Improving the code or fixing an issue](#improving-the-code-or-fixing-an-issue)
    - [Join us on Slack](#join-us-on-slack)
    - [License](#license)
  - [Catalog of apps and CLIs](#catalog-of-apps-and-clis)
    - [Catalog of Apps](#catalog-of-apps)
    - [Catalog of CLIs](#catalog-of-clis)

## Support arkade 👋

Arkade is built to save you time so you can focus and get productive quickly.

<a href="https://github.com/sponsors/alexellis/">
<img alt="Sponsor this project" src="https://github.com/alexellis/alexellis/blob/master/sponsor-today.png" width="90%">
</a>

You can support Alex's work on arkade [via GitHub Sponsors](https://github.com/sponsors/alexellis/).

Or get a copy of his eBook on Go so you can learn how to build tools like k3sup, arkade, and OpenFaaS for yourself:

<a href="https://openfaas.gumroad.com/l/everyday-golang">
<img src="https://public-files.gumroad.com/7j27fj7c5xqxm3f9lyxj1pg8oa1w" alt="Buy Everyday Go" width="50%"></a>

## Should you try arkade?

> I was setting up a new dev environment yesterday. Kind, helm, kustomize, kubectl, all this stuff. My take is - arkade is highly underappreciated.
> I'd spend an hour in the past to install such tools. With arkade it was under ten minutes.
>
> [Ivan Velichko](https://twitter.com/iximiuz/status/1422605221226860548?s=20), SRE @ Booking.com

> Before arkade whenever I used to spin up an instance, I used to go to multiple sites and download the binary. Arkade is one of my favourite tools.
> 
> [Kumar Anurag](https://kubesimplify.com/arkade) - Cloud Native Enthusiast

> It's hard to use K8s without Arkade these days.
> My team at @lftechnology absolutely loves it.
>
> [@Yankexe](https://twitter.com/yankexe/status/1305427718050250754?s=20)

> arkade is really a great tool to install CLI tools, and system packages, check this blog on how to get started with arkade it's a time saver.
> 
> [Kiran Satya Raj](https://twitter.com/jksrtwt/status/1556592117627047936?s=20&t=g0gnSP98jg3ZwU7sQqUrLw)

> This is real magic get #kubernetes up and going in a second; then launch #openfaas a free better than lambda solution that uses docker images.
>
> [Greg](https://twitter.com/cactusanddove) runs Fullstack JS and is a JavaScript developer

> for getting the basics installed, nothing beats arkade
> it can install commonly used cli tools like kubectl locally for you, as well as common k8s pkgs like ingress-nginx or portainer
>
> [@arghzero](https://twitter.com/ArghZero/status/1346097288851070983?s=20)

> I finally got around to installing Arkade, super simple!
> quicker to install this than the argocli standalone commands, but there are lots of handy little tools in there.
> also, the neat little part about arkade, not only does it make it easy to install a ton of different apps and CLIs you can also get the info on them as well pretty quickly.
> 
> [Michael Cade @ Kasten](https://twitter.com/MichaelCade1/status/1390403831167700995?s=20)

> You've to install the latest and greatest tools for your daily @kubernetesio tasks? No problem, check out #arkade the open source #kubernetes marketplace 👍
>
> [Thorsten Hans](https://twitter.com/ThorstenHans/status/1457982292597608449?s=20) - Cloud Native consultant

> If you want to install quickly a new tool in your dev env or in your k8s cluster you can use the Arkade (https://github.com/alexellis/arkade) easy and quick you should it try out! Ps. I contribute to this project 🥰
>
> [Carlos Panato](https://twitter.com/comedordexis/status/1423339283713347587) - Staff engineer @ Mattermost

> arkade is the 'brew install' of Kubernetes. You can install and run an application in a single command. Finally! https://github.com/alexellis/arkade / by Alex Ellis
>
> [John Arundel](https://twitter.com/bitfield/status/1242385165445455872?s=20) - Cloud consultant, author

## Demo

![demo](https://vhs.charm.sh/vhs-7Fyg69mwbYHFuUtSKnWMYT.gif)

## Getting arkade

```bash
# Note: you can also run without `sudo` and move the binary yourself
curl -sLS https://get.arkade.dev | sudo sh

arkade --help
ark --help  # a handy alias

# Windows users with Git Bash
curl -sLS https://get.arkade.dev | sh
```

> Windows users: arkade requires bash to be available, therefore Windows users should [install and use Git Bash](https://git-scm.com/downloads)

An alias of `ark` is created at installation time, so you can also run `ark install APP`

## Usage overview

Arkade can be used to install Kubernetes apps or to download CLI tools.

* `arkade install` - install a Kubernetes app
* `arkade info` - see the post installation screen for a Kubernetes app
* `arkade get` - download a CLI tool
* `arkade update` - perform a self-update of arkade on MacOS and Linux

An arkade "app" could represent a helm chart such as `openfaas/faas-netes`, a custom CLI installer such as `istioctl`, or a set of static manifests (i.e. MetalLB).

An arkade "tool" is a CLI that can be downloaded for your operating system. Arkade downloads statically-linked binaries from their upstream locations on GitHub or the vendor's chosen URL such as with `kubectl` and `terraform`.

> Did you know? Arkade users run `arkade get` both on their local workstations, and on their CI runners such as GitHub Actions or Jenkins.

## Download CLI tools with arkade

arkade downloads the correct version of a CLI for your OS and CPU.

With automatic detection of: Windows / MacOS / Linux / Intel / ARM.

```bash
# Download a binary release of a tool

arkade get kubectl

# Download a specific version of that tool
arkade get kubectl@v1.22.0

# Download multiple tools at once
arkade get kubectl \
  helm \
  istioctl

# Download multiple specific versions
arkade get faas-cli@0.13.15 \
  kubectl@v1.22.0

# Override machine os/arch
arkade get faas-cli \
  --arch arm64 \
  --os linux

# Override machine os/arch
arkade get faas-cli \
  --arch arm64 \
  --os darwin
```
> This is a time saver compared to searching for download pages every time you need a tool.

Files are stored at `$HOME/.arkade/bin/`

Want to download tools to a custom path such as into the GitHub Actions cached tool folder?

```bash
arkade get faas-cli kubectl \
  --path $HOME/runner/_work/_tools

# Usage:
/runner/_work/_tools/faas-cli version

PATH=$PATH:$HOME/runner/_work/_tools
faas-cli version
```

Think of `arkade get TOOL` as a doing for CLIs, what `arkade install` does for helm.

Adding a new tool for download is as simple as editing [tools.go](https://github.com/alexellis/arkade/blob/master/pkg/get/tools.go).

[Click here for the full catalog of CLIs](#catalog-of-clis)

## Install System Packages

System packages are tools designed for installation on a Linux workstation, server or CI runner.

These are a more limited group of applications designed for quick setup, scripting and CI, and generally do not fit into the `arkade get` pattern, due to additional installation steps or system configuration.

```bash
# Show packages
arkade system install

# Show package flags
arkade system install go --help

# Install latest version of Go to /usr/local/bin/go
arkade system install go

# Install Go 1.18 to /tmp/go
arkade system install go \
  --version 1.18 \
  --path /tmp/

# Install containerd for ARM64, 32-bit ARM or x86_64
# with systemd enabled
arkade system install containerd \
  --systemd
```

Run the following to see what's available `arkade system install`:

```
  actions-runner  Install GitHub Actions Runner
  buildkitd       Install Buildkitd
  caddy           Install Caddy Server
  cni             Install CNI plugins
  containerd      Install containerd
  firecracker     Install Firecracker
  gitlab-runner   Install GitLab Runner
  go              Install Go
  node            Install Node.js
  prometheus      Install Prometheus
  pwsh            Install Powershell
  registry        Install registry
  tc-redirect-tap Install tc-redirect-tap
```

The initial set of system apps is now complete, learn more in the original proposal: [Feature: system packages for Linux servers, CI and workstations #654](https://github.com/alexellis/arkade/issues/654)

## Install Packages from OCI images

For packages distributed in Open Container Initiative (OCI) images, you can use `arkade oci install` to extract them to a given folder on your system.

vmmeter is one example of a package that is only published as a container image, which is not released on a GitHub releases page.

```bash
arkade oci install ghcr.io/openfaasltd/vmmeter \
  --path /usr/local/bin
```

* `--path` - the folder to extract the package to
* `--version` - the version of the package to extract, if not specified the `:latest` tag is used
* `--arch` - the architecture to extract, if not specified the host's architecture is used

## Install CLIs during CI with GitHub Actions

* [alexellis/arkade-get@master](https://github.com/alexellis/arkade-get)

Example of downloading faas-cli (specific version) and kubectl (latest), putting them into the PATH automatically, and executing one of them in a subsequent step.

```yaml
    - uses: alexellis/arkade-get@master
      with:
        kubectl: latest
        faas-cli: 0.14.10
    - name: check for faas-cli
      run: |
        faas-cli version
```

If you just need system applications, you could also try "setup-arkade":

* [alexellis/setup-arkade@master](https://github.com/alexellis/setup-arkade)

```yaml
    - uses: alexellis/setup-arkade@v2
    - name: Install containerd and go
      run: |
        arkade system install containerd
        arkade system install go
```

## Bump Helm chart versions

To bump the patch version of your Helm chart, run `arkade chart bump -f ./chart/values.yaml`. This updates the patch component of the version specified in Chart.yaml.

```bash
arkade chart bump -f ./charts/flagger/values.yaml
charts/flagger/Chart.yaml 1.36.0 => 1.37.0
```

By default, the new version is written to stdout. To bump the version in the file, run the above command with the `--write` flag.
To bump the version in the chart's Chart.yaml only if the chart has any changes, specify the `--check-for-updates` flag:

```bash
arkade chart bump -f ./charts/flagger/values.yaml --check-for-updates
no changes detected in charts/flagger/values.yaml; skipping version bump
```

The directory that contains the Helm chart should be a Git repository. If the flag is specified, the command runs `git diff --exit-code <file>` to figure out if the file has any changes.

## Verify and upgrade images in Helm charts

There are two commands built into arkade designed for software vendors and open source maintainers.

* `arkade helm chart upgrade` - run this command to scan for container images and update them automatically by querying a remote registry. 
* `arkade helm chart verify` - after changing the contents of a values.yaml or docker-compose.yaml file, this command will check each image exists on a remote registry

Whilst end-users may use a GitOps-style tool to deploy charts and update their versions, maintainers need to make conscious decisions about when and which images to change within a Helm chart or compose file.

These two features are used by OpenFaaS Ltd on projects and products like OpenFaaS CE/Pro (Serverless platform) and faasd (docker-compose file). 

### Upgrade images within a Helm chart

With the command `arkade chart upgrade` you can upgrade the image tags of a Helm chart from within a values.yaml file to the latest available semantically versioned image.

Original YAML file:

```yaml
stan:
  # Image used for nats deployment when using async with NATS-Streaming.
  image: nats-streaming:0.24.6
```

Running the command with `--verbose` prints the upgraded tags to stderr, allowing the output to stdout to be piped to a file.

```bash
arkade chart upgrade -f \
  ~/go/src/github.com/openfaas/faas-netes/chart/openfaas/values.yaml \
  --verbose

2023/01/03 10:12:47 Verifying images in: /home/alex/go/src/github.com/openfaas/faas-netes/chart/openfaas/values.yaml
2023/01/03 10:12:47 Found 18 images
2023/01/03 10:12:48 [natsio/prometheus-nats-exporter] 0.8.0 => 0.10.1
2023/01/03 10:12:50 [nats-streaming] 0.24.6 => 0.25.2
2023/01/03 10:12:52 [prom/prometheus] v2.38.0 => 2.41.0
2023/01/03 10:12:54 [prom/alertmanager] v0.24.0 => 0.25.0
2023/01/03 10:12:54 [nats] 2.9.2 => 2.9.10
```

Updated YAML file printed to console:

```yaml
stan:
  # Image used for nats deployment when using async with NATS-Streaming.
  image: nats-streaming:0.25.2
```

Write the updated image tags back to the file:

```bash
arkade chart upgrade -f \
  ~/go/src/github.com/openfaas/faasd/docker-compose.yaml \
  --write
```

Supported:

* `image:` - at the top level
* `component.image:` i.e. one level of nesting
* Docker Hub and GitHub Container Registry

Not supported yet:
* Custom strings that don't match the word "image": `clientImage: `
* Split fields for the image and tag name i.e. `image.name` and `image.tag`
* Third-level nesting `openfaas.gateway.image`

## Verify images within a helm chart

The `arkade chart verify` command validates that all images specified are accessible on a remote registry and takes a values.yaml file as its input.

Successful checking of a chart with `image: ghcr.io/openfaas/cron-connector:TAG`:

```bash
arkade chart verify  -f ~/go/src/github.com/openfaas/faas-netes/chart/cron-connector/values.yaml

echo $?
0
```

There is an exit code of zero and no output when the check passes.

You can pass `--verbose` to see a detailed view of what's happening.

Checking of nested components, where two of the images do not exist `autoscaler.image` and `dashboard.image`:

```bash
arkade chart verify  -f ~/go/src/github.com/openfaas/faas-netes/chart/openfaas/values.yamlecho $?
2 images are missing in /Users/alex/go/src/github.com/openfaas/faas-netes/chart/openfaas/values.yaml

COMPONENT           IMAGE
dashboard           ghcr.io/openfaasltd/openfaas-dashboard:0.9.8
autoscaler          ghcr.io/openfaasltd/autoscaler:0.2.5

Error: verifying failed

echo $?
1
```

Supported:

* `image:` - at the top level
* `component.image:` i.e. one level of nesting

Not supported yet:
* Custom strings that don't match the word "image": `clientImage: `
* Split fields for the image and tag name i.e. `image.name` and `image.tag`
* Third-level nesting `openfaas.gateway.image`

## Installing apps with arkade

You'll need a Kubernetes cluster to arkade. Unlike cloud-based marketplaces, arkade doesn't have any special pre-requirements and can be used with any private or public cluster.

### Create a Kubernetes cluster

If you have Docker installed, then you can install Kubernetes using KinD in a matter of moments:

```bash
arkade get kubectl@v1.22.0 \
  kind@v0.11.1

kind create cluster
```

You can also download k3d [k3s](https://github.com/rancher/k3s) in the same way with `arkade get k3d`.

### Install a Kubernetes app

No need to worry about whether you're installing to Intel or ARM architecture, the correct values will be set for you automatically.

```bash
arkade install openfaas \
  --gateways 2 \
  --load-balancer false
```

The post-installation message shows you how to connect. And whenever you want to see those details again, just run `arkade info openfaas`.

There are even more options you can choose with `arkade install openfaas --help` - the various flags you see map to settings from the helm chart README, that you'd usually have to look up and set via a `values.yaml` file.

If there's something missing from the list of flags that you need, arkade also supports `--set` for any arkade app that uses helm. Note that not every app uses helm.

Remember how awkward it was last time you installed the [Kubernetes dashboard](https://github.com/kubernetes/dashboard)? And how you could never remember the command to get the token to log in?

```bash
arkade install kubernetes-dashboard
```

Forgot your token? `arkade info kubernetes-dashboard`

This is an example of an arkade app that uses static YAML manifests instead of helm.

Prefer [Portainer](https://www.portainer.io)? Just run: `arkade install portainer`

### Uninstall an app

Run `arkade uninstall` or `arkade delete` for more information on how to remove applications from a Kubernetes cluster.

### Reduce the repetition

[Normally up to a dozen commands](https://cert-manager.io/docs/installation/kubernetes/) (including finding and downloading helm), now just one. No searching for the correct CRD to apply, no trying to install helm, no trying to find the correct helm repo to add:

```bash
arkade install cert-manager
```

Other common tools:

```bash
arkade install ingress-nginx

arkade install metrics-server
```

### Say goodbye to values.yaml and hello to flags

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
  -h, --help               help for ingress-nginx
      --host-mode          If we should install ingress-nginx in host mode.
  -n, --namespace string   The namespace used for installation (default "default")
      --update-repo        Update the helm repo (default true)
```

### Override with `--set`

You can also set helm overrides, for apps which use helm via `--set`

```bash
ark install openfaas --set faasIdler.dryRun=false
```

After installation, an info message will be printed with help for usage, you can get back to this at any time via:

```bash
arkade info <NAME>
```

### Compounding apps

Apps are easier to discover and install than helm chart which involve many more manual steps, however when you compound apps together, they really save you time.

#### Get a self-hosted TLS registry with authentication

Here's how you can get a self-hosted Docker registry with TLS and authentication in just 5 commands on an empty cluster:

Here's how you would bootstrap OpenFaaS with TLS:

```bash
arkade install ingress-nginx
arkade install cert-manager
arkade install openfaas
arkade install openfaas-ingress \
  --email web@example.com \
  --domain openfaas.example.com
```

And here's what it looks like for a private Docker registry with authentication enabled:

```bash
arkade install ingress-nginx
arkade install cert-manager
arkade install docker-registry
arkade install docker-registry-ingress \
  --email web@example.com \
  --domain reg.example.com
```

#### Get a public IP for a private cluster and your IngressController

And if you're running on a private cloud, on-premises or on your laptop, you can simply add the [inlets-operator](https://github.com/inlets/inlets-operator/) using [inlets](https://docs.inlets.dev/) to get a secure TCP tunnel and a public IP address.

```bash
arkade install inlets-operator \
  --access-token $HOME/digitalocean-token \
  --region lon1 \
  --provider digitalocean
```

This makes your cluster behave like it was on a public cloud and LoadBalancer IPs go from Pending to a real, functioning IP.

### Explore the apps

You can view the various apps available with `arkade install / --help`, more are available when you run the command yourself.

```bash
arkade install --help
ark --help

Examples:
  arkade install
  arkade install openfaas --helm3 --gateways=2
  arkade install inlets-operator --token-file $HOME/do-token
```

See the full catalog of apps: [See all apps](#catalog-of-apps)

## Community & contributing

### Tutorials & community blog posts

#### Watch a video walk-through by Alex Ellis

[![Install Apps and CLIs to Kubernetes](http://img.youtube.com/vi/8wU9s_mua8M/hqdefault.jpg)](https://www.youtube.com/watch?v=8wU9s_mua8M)

#### Featured tutorials

* [arkade by example — Kubernetes apps, the easy way 😎](https://medium.com/@alexellisuk/kubernetes-apps-the-easy-way-f06d9e5cad3c) - Alex Ellis
* [Walk-through — install Kubernetes to your Raspberry Pi in 15 minutes](https://medium.com/@alexellisuk/walk-through-install-kubernetes-to-your-raspberry-pi-in-15-minutes-84a8492dc95a)
* [Get a TLS-enabled Docker registry in 5 minutes](https://blog.alexellis.io/get-a-tls-enabled-docker-registry-in-5-minutes/) - Alex Ellis
* [Get TLS for OpenFaaS the easy way with arkade](https://blog.alexellis.io/tls-the-easy-way-with-openfaas-and-k3sup/) - Alex Ellis

#### Official blog posts

* [Two year update: Building an Open Source Marketplace for Kubernetes](https://blog.alexellis.io/kubernetes-marketplace-two-year-update/)
* [Why did the OpenFaaS community build arkade and what's in it for you?](https://www.openfaas.com/blog/openfaas-arkade/) - Alex Ellis

#### Community posts

* [A bit of Istio before tea-time](https://blog.alexellis.io/a-bit-of-istio-before-tea-time/) - Alex Ellis
* [Kubernetes: Automatic Let's Encrypt Certificates for Services with arkade](https://medium.com/@admantium/kubernetes-automatic-lets-encrypt-certificates-for-services-2a5f4aa7f886)
* [Introducing Arkade - The Kubernetes app installer](https://blog.heyal.co.uk/introducing-arkade/) - Alistair Hey
* [Portainer for kubernetes in less than 60 seconds!!](https://www.portainer.io/2020/04/portainer-for-kubernetes-in-less-than-60-seconds/) - by Saiyam Pathak
* [Video walk-through with DJ Adams - Pi & Kubernetes with k3s, k3sup, arkade and OpenFaaS](https://www.youtube.com/watch?v=ZiR3QEfBivk)
* [Coffee chat: Easy way to install Kubernetes Apps - arkade (ark)](https://sachcode.com/tech/coffee-chat-easy-way-install-kubernetes-apps/) by Sachin Jha
* [Arkade & OpenFaaS: serverless on the spot](https://zero2datadog.readthedocs.io/en/latest/faas.html) by Blaise Pabon
* ["Tool of the Day" with Adrian Goins from Rancher Labs](https://youtu.be/IWtEtfpqoRg?t=1425)

### Suggest a new app

To suggest a new app, please check past issues and [raise an issue for it](https://github.com/alexellis/arkade). Think also whether your app suggestion would be a good candidate for a Sponsored App.

## Sponsored apps

You can now propose your project or product as a Sponsored App. Sponsored Apps work just like any other app that we've curated, however they will have a note next to them in the app description `(sponsored)` and a link to your chosen site upon installation. An app sponsorship can be purchased for a minimum of 12 months and includes free development of the Sponsored App, with ongoing support via GitHub for the Sponsored App for the duration only. Ongoing support will be limited to a set amount of hours per month.

When your sponsorship expires the Sponsored App will be removed from arkade, and the ongoing support will cease. A Sponsored App can be renewed 60 days prior to expiration subject to a separate agreement and payment.

Example:

```bash
arkade VENDOR install PRODUCT
arkade acmeco install dashboard
```

[Contact OpenFaas Ltd](mailto:contact@openfaas.com) to find out how you can have your Sponsored App added to arkade.

## FAQ

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

Before contributing code, please see the [CONTRIBUTING guide](https://github.com/alexellis/arkade/blob/master/CONTRIBUTING.md). Note that arkade uses the same guide as [inlets.dev](https://inlets.dev/).

Both Issues and PRs have their own templates. Please fill out the whole template.

All commits must be signed-off as part of the [Developer Certificate of Origin (DCO)](https://developercertificate.org)

### Join us on Slack

Join `#contributors` at [slack.openfaas.io](https://slack.openfaas.io)

### License

MIT

## Catalog of apps and CLIs

An app is software or an add-on for your Kubernetes cluster.

A CLI or "tool" is a command line tool that you run directly on your own workstation or a CI runner.

### Catalog of Apps

|          TOOL           |                             DESCRIPTION                             |
|-------------------------|---------------------------------------------------------------------|
| argocd                  | Install argocd                                                      |
| cassandra               | Install cassandra                                                   |
| cert-manager            | Install cert-manager                                                |
| chart                   | Install the specified helm chart                                    |
| cockroachdb             | Install CockroachDB                                                 |
| consul-connect          | Install Consul Service Mesh                                         |
| cron-connector          | Install cron-connector for OpenFaaS                                 |
| crossplane              | Install Crossplane                                                  |
| docker-registry         | Install a community maintained Docker registry chart                |
| docker-registry-ingress | Install registry ingress with TLS                                   |
| falco                   | Install Falco                                                       |
| gitea                   | Install gitea                                                       |
| gitlab                  | Install GitLab                                                      |
| grafana                 | Install grafana                                                     |
| influxdb                | Install influxdb                                                    |
| ingress-nginx           | Install ingress-nginx                                               |
| inlets-operator         | Install inlets-operator                                             |
| istio                   | Install istio                                                       |
| jenkins                 | Install jenkins                                                     |
| kafka                   | Install Confluent Platform Kafka                                    |
| kafka-connector         | Install kafka-connector for OpenFaaS                                |
| kong-ingress            | Install kong-ingress for OpenFaaS                                   |
| kube-image-prefetch     | Install kube-image-prefetch                                         |
| kube-state-metrics      | Install kube-state-metrics                                          |
| kubernetes-dashboard    | Install kubernetes-dashboard                                        |
| kuma                    | Install Kuma                                                        |
| kyverno                 | Install Kyverno                                                     |
| linkerd                 | Install linkerd                                                     |
| loki                    | Install Loki for monitoring and tracing                             |
| metallb-arp             | Install MetalLB in L2 (ARP) mode                                    |
| metrics-server          | Install metrics-server                                              |
| minio                   | Install minio                                                       |
| mongodb                 | Install mongodb                                                     |
| mqtt-connector          | Install mqtt-connector for OpenFaaS                                 |
| nats-connector          | Install OpenFaaS connector for NATS                                 |
| nfs-provisioner         | Install nfs subdir external provisioner                             |
| opa-gatekeeper          | Install Open Policy Agent (OPA) Gatekeeper                          |
| openfaas                | Install openfaas                                                    |
| openfaas-ingress        | Install openfaas ingress with TLS                                   |
| openfaas-loki           | Install Loki-OpenFaaS and Configure Loki logs provider for OpenFaaS |
| portainer               | Install portainer to visualise and manage containers                |
| postgresql              | Install postgresql                                                  |
| prometheus              | Install Prometheus for monitoring                                   |
| qemu-static             | Install qemu-user-static                                            |
| rabbitmq                | Install rabbitmq                                                    |
| redis                   | Install redis                                                       |
| registry-creds          | Install registry-creds                                              |
| sealed-secret           | Install sealed-secrets                                              |
| tekton                  | Install Tekton pipelines and dashboard                              |
| traefik2                | Install traefik2                                                    |
| vault                   | Install vault                                                       |
| waypoint                | Install Waypoint                                                    |

There are 52 apps that you can install on your cluster.

> Note to contributors, run `go build && ./arkade install --print-table` to generate this list

### Catalog of CLIs

|                                     TOOL                                     |                                                                            DESCRIPTION                                                                            |
|------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [actions-usage](https://github.com/self-actuated/actions-usage)              | Get usage insights from GitHub Actions.                                                                                                                           |
| [actuated-cli](https://github.com/self-actuated/actuated-cli)                | Official CLI for actuated.dev                                                                                                                                     |
| [argocd](https://github.com/argoproj/argo-cd)                                | Declarative, GitOps continuous delivery tool for Kubernetes.                                                                                                      |
| [argocd-autopilot](https://github.com/argoproj-labs/argocd-autopilot)        | An opinionated way of installing Argo-CD and managing GitOps repositories.                                                                                        |
| [arkade](https://github.com/alexellis/arkade)                                | Portable marketplace for downloading your favourite DevOps CLIs and installing helm charts, with a single command.                                                |
| [atuin](https://github.com/atuinsh/atuin)                                    | Sync, search, and backup shell history with Atuin.                                                                                                                |
| [autok3s](https://github.com/cnrancher/autok3s)                              | Run Rancher Lab's lightweight Kubernetes distribution k3s everywhere.                                                                                             |
| [buildx](https://github.com/docker/buildx)                                   | Docker CLI plugin for extended build capabilities with BuildKit.                                                                                                  |
| [bun](https://github.com/oven-sh/bun)                                        | Bun is an incredibly fast JavaScript runtime, bundler, transpiler, and package manager – all in one.                                                              |
| [butane](https://github.com/coreos/butane)                                   | Translates human readable Butane Configs into machine readable Ignition Configs                                                                                   |
| [caddy](https://github.com/caddyserver/caddy)                                | Caddy is an extensible server platform that uses TLS by default                                                                                                   |
| [ch-remote](https://github.com/cloud-hypervisor/cloud-hypervisor)            | The ch-remote binary is used for controlling an running Virtual Machine.                                                                                          |
| [cilium](https://github.com/cilium/cilium-cli)                               | CLI to install, manage & troubleshoot Kubernetes clusters running Cilium.                                                                                         |
| [civo](https://github.com/civo/cli)                                          | CLI for interacting with your Civo resources.                                                                                                                     |
| [cloud-hypervisor](https://github.com/cloud-hypervisor/cloud-hypervisor)     | Cloud Hypervisor is an open source Virtual Machine Monitor (VMM) that runs on top of the KVM hypervisor and the Microsoft Hypervisor (MSHV).                      |
| [clusterawsadm](https://github.com/kubernetes-sigs/cluster-api-provider-aws) | Kubernetes Cluster API Provider AWS Management Utility                                                                                                            |
| [clusterctl](https://github.com/kubernetes-sigs/cluster-api)                 | The clusterctl CLI tool handles the lifecycle of a Cluster API management cluster                                                                                 |
| [cmctl](https://github.com/cert-manager/cmctl)                               | cmctl is a CLI tool that helps you manage cert-manager and its resources inside your cluster.                                                                     |
| [conftest](https://github.com/open-policy-agent/conftest)                    | Write tests against structured configuration data using the Open Policy Agent Rego query language                                                                 |
| [consul](https://github.com/hashicorp/consul)                                | A solution to connect and configure applications across dynamic, distributed infrastructure                                                                       |
| [copa](https://github.com/project-copacetic/copacetic)                       | CLI for patching container images                                                                                                                                 |
| [cosign](https://github.com/sigstore/cosign)                                 | Container Signing, Verification and Storage in an OCI registry.                                                                                                   |
| [cr](https://github.com/helm/chart-releaser)                                 | Hosting Helm Charts via GitHub Pages and Releases                                                                                                                 |
| [crane](https://github.com/google/go-containerregistry)                      | crane is a tool for interacting with remote images and registries                                                                                                 |
| [croc](https://github.com/schollz/croc)                                      | Easily and securely send things from one computer to another                                                                                                      |
| [crossplane](https://github.com/crossplane/crossplane)                       | Simplify some development and administration aspects of Crossplane.                                                                                               |
| [dagger](https://github.com/dagger/dagger)                                   | A portable devkit for CI/CD pipelines.                                                                                                                            |
| [devspace](https://github.com/devspace-sh/devspace)                          | Automate your deployment workflow with DevSpace and develop software directly inside Kubernetes.                                                                  |
| [dive](https://github.com/wagoodman/dive)                                    | A tool for exploring each layer in a docker image                                                                                                                 |
| [docker-compose](https://github.com/docker/compose)                          | Define and run multi-container applications with Docker.                                                                                                          |
| [doctl](https://github.com/digitalocean/doctl)                               | Official command line interface for the DigitalOcean API.                                                                                                         |
| [duplik8s](https://github.com/Telemaco019/duplik8s)                          | kubectl plugin to duplicate resources in a Kubernetes cluster.                                                                                                    |
| [eks-node-viewer](https://github.com/awslabs/eks-node-viewer)                | eks-node-viewer is a tool for visualizing dynamic node usage within an EKS cluster.                                                                               |
| [eksctl](https://github.com/eksctl-io/eksctl)                                | Amazon EKS Kubernetes cluster management                                                                                                                          |
| [eksctl-anywhere](https://github.com/aws/eks-anywhere)                       | Run Amazon EKS on your own infrastructure                                                                                                                         |
| [etcd](https://github.com/etcd-io/etcd)                                      | Distributed reliable key-value store for the most critical data of a distributed system.                                                                          |
| [faas-cli](https://github.com/openfaas/faas-cli)                             | Official CLI for OpenFaaS.                                                                                                                                        |
| [faasd](https://github.com/openfaas/faasd)                                   | faasd - a lightweight & portable faas engine                                                                                                                      |
| [firectl](https://github.com/firecracker-microvm/firectl)                    | Command-line tool that lets you run arbitrary Firecracker MicroVMs                                                                                                |
| [flux](https://github.com/fluxcd/flux2)                                      | Continuous Delivery solution for Kubernetes powered by GitOps Toolkit.                                                                                            |
| [flyctl](https://github.com/superfly/flyctl)                                 | Command line tools for fly.io services                                                                                                                            |
| [fstail](https://github.com/alexellis/fstail)                                | Tail modified files in a directory.                                                                                                                               |
| [fzf](https://github.com/junegunn/fzf)                                       | General-purpose command-line fuzzy finder                                                                                                                         |
| [gh](https://github.com/cli/cli)                                             | GitHub’s official command line tool.                                                                                                                              |
| [glab](https://github.com/gitlab-org/cli)                                    | A GitLab CLI tool bringing GitLab to your command line.                                                                                                           |
| [golangci-lint](https://github.com/golangci/golangci-lint)                   | Go linters aggregator.                                                                                                                                            |
| [gomplate](https://github.com/hairyhenderson/gomplate)                       | A flexible commandline tool for template rendering. Supports lots of local and remote datasources.                                                                |
| [goreleaser](https://github.com/goreleaser/goreleaser)                       | Deliver Go binaries as fast and easily as possible                                                                                                                |
| [gptscript](https://github.com/gptscript-ai/gptscript)                       | Natural Language Programming                                                                                                                                      |
| [grafana-agent](https://github.com/grafana/agent)                            | Grafana Agent is a telemetry collector for sending metrics, logs, and trace data to the opinionated Grafana observability stack.                                  |
| [grype](https://github.com/anchore/grype)                                    | A vulnerability scanner for container images and filesystems                                                                                                      |
| [hadolint](https://github.com/hadolint/hadolint)                             | A smarter Dockerfile linter that helps you build best practice Docker images                                                                                      |
| [helm](https://github.com/helm/helm)                                         | The Kubernetes Package Manager: Think of it like apt/yum/homebrew for Kubernetes.                                                                                 |
| [helmfile](https://github.com/helmfile/helmfile)                             | Deploy Kubernetes Helm Charts                                                                                                                                     |
| [hey](https://github.com/alexellis/hey)                                      | Load testing tool                                                                                                                                                 |
| [hostctl](https://github.com/guumaster/hostctl)                              | Dev tool to manage /etc/hosts like a pro!                                                                                                                         |
| [hubble](https://github.com/cilium/hubble)                                   | CLI for network, service & security observability for Kubernetes clusters running Cilium.                                                                         |
| [hugo](https://github.com/gohugoio/hugo)                                     | Static HTML and CSS website generator.                                                                                                                            |
| [influx](https://github.com/influxdata/influxdb)                             | InfluxDB’s command line interface (influx) is an interactive shell for the HTTP API.                                                                              |
| [inlets-pro](https://github.com/inlets/inlets-pro)                           | Cloud Native Tunnel for HTTP and TCP traffic.                                                                                                                     |
| [inletsctl](https://github.com/inlets/inletsctl)                             | Automates the task of creating an exit-server (tunnel server) on public cloud infrastructure.                                                                     |
| [istioctl](https://github.com/istio/istio)                                   | Service Mesh to establish a programmable, application-aware network using the Envoy service proxy.                                                                |
| [jq](https://github.com/jqlang/jq)                                           | jq is a lightweight and flexible command-line JSON processor                                                                                                      |
| [just](https://github.com/casey/just)                                        | Just a command runner                                                                                                                                             |
| [k0s](https://github.com/k0sproject/k0s)                                     | Zero Friction Kubernetes                                                                                                                                          |
| [k0sctl](https://github.com/k0sproject/k0sctl)                               | A bootstrapping and management tool for k0s clusters                                                                                                              |
| [k3d](https://github.com/k3d-io/k3d)                                         | Helper to run Rancher Lab's k3s in Docker.                                                                                                                        |
| [k3s](https://github.com/k3s-io/k3s)                                         | Lightweight Kubernetes                                                                                                                                            |
| [k3sup](https://github.com/alexellis/k3sup)                                  | Bootstrap Kubernetes with k3s over SSH < 1 min.                                                                                                                   |
| [k9s](https://github.com/derailed/k9s)                                       | Provides a terminal UI to interact with your Kubernetes clusters.                                                                                                 |
| [kail](https://github.com/boz/kail)                                          | Kubernetes log viewer.                                                                                                                                            |
| [keploy](https://github.com/keploy/keploy)                                   | Test generation for Developers. Generate tests and stubs for your application that actually work!                                                                 |
| [kgctl](https://github.com/squat/kilo)                                       | A CLI to manage Kilo, a multi-cloud network overlay built on WireGuard and designed for Kubernetes.                                                               |
| [kim](https://github.com/rancher/kim)                                        | Build container images inside of Kubernetes. (Experimental)                                                                                                       |
| [kind](https://github.com/kubernetes-sigs/kind)                              | Run local Kubernetes clusters using Docker container nodes.                                                                                                       |
| [kops](https://github.com/kubernetes/kops)                                   | Production Grade K8s Installation, Upgrades, and Management.                                                                                                      |
| [krew](https://github.com/kubernetes-sigs/krew)                              | Package manager for kubectl plugins.                                                                                                                              |
| [ktop](https://github.com/vladimirvivien/ktop)                               | A top-like tool for your Kubernetes cluster.                                                                                                                      |
| [kube-bench](https://github.com/aquasecurity/kube-bench)                     | Checks whether Kubernetes is deployed securely by running the checks documented in the CIS Kubernetes Benchmark.                                                  |
| [kube-burner](https://github.com/cloud-bulldozer/kube-burner)                | A tool aimed at stressing Kubernetes clusters by creating or deleting a high quantity of objects.                                                                 |
| [kube-linter](https://github.com/stackrox/kube-linter)                       | KubeLinter is a static analysis tool that checks Kubernetes YAML files and Helm charts to ensure the applications represented in them adhere to best practices.   |
| [kube-score](https://github.com/zegl/kube-score)                             | A tool that performs static code analysis of your Kubernetes object definitions.                                                                                  |
| [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)                | Framework for building Kubernetes APIs using custom resource definitions (CRDs).                                                                                  |
| [kubecm](https://github.com/sunny0826/kubecm)                                | Easier management of kubeconfig.                                                                                                                                  |
| [kubecolor](https://github.com/kubecolor/kubecolor)                          | KubeColor is a kubectl replacement used to add colors to your kubectl output.                                                                                     |
| [kubeconform](https://github.com/yannh/kubeconform)                          | A FAST Kubernetes manifests validator, with support for Custom Resources                                                                                          |
| [kubectl](https://github.com/kubernetes/kubernetes)                          | Run commands against Kubernetes clusters                                                                                                                          |
| [kubectx](https://github.com/ahmetb/kubectx)                                 | Faster way to switch between clusters.                                                                                                                            |
| [kubens](https://github.com/ahmetb/kubectx)                                  | Switch between Kubernetes namespaces smoothly.                                                                                                                    |
| [kubescape](https://github.com/kubescape/kubescape)                          | kubescape is the first tool for testing if Kubernetes is deployed securely as defined in Kubernetes Hardening Guidance by NSA and CISA                            |
| [kubeseal](https://github.com/bitnami-labs/sealed-secrets)                   | A Kubernetes controller and tool for one-way encrypted Secrets                                                                                                    |
| [kubetail](https://github.com/johanhaleby/kubetail)                          | Bash script to tail Kubernetes logs from multiple pods at the same time.                                                                                          |
| [kubetrim](https://github.com/alexellis/kubetrim)                            | Tidy up old Kubernetes clusters from kubeconfig.                                                                                                                  |
| [kubeval](https://github.com/instrumenta/kubeval)                            | Validate your Kubernetes configuration files, supports multiple Kubernetes versions                                                                               |
| [kubie](https://github.com/sbstp/kubie)                                      | A more powerful alternative to kubectx and kubens                                                                                                                 |
| [kumactl](https://github.com/kumahq/kuma)                                    | kumactl is a CLI to interact with Kuma and its data                                                                                                               |
| [kustomize](https://github.com/kubernetes-sigs/kustomize)                    | Customization of kubernetes YAML configurations                                                                                                                   |
| [kwok](https://github.com/kubernetes-sigs/kwok)                              | KWOK stands for Kubernetes WithOut Kubelet, responsible for simulating the lifecycle of fake nodes, pods, and other Kubernetes API resources                      |
| [kwokctl](https://github.com/kubernetes-sigs/kwok)                           | CLI tool designed to streamline the creation and management of clusters, with nodes simulated by `kwok`                                                           |
| [kyverno](https://github.com/kyverno/kyverno)                                | CLI to apply and test Kyverno policies outside a cluster.                                                                                                         |
| [labctl](https://github.com/iximiuz/labctl)                                  | iximiuz Labs control - start remote microVM playgrounds from the command line.                                                                                    |
| [lazydocker](https://github.com/jesseduffield/lazydocker)                    | A simple terminal UI for both docker and docker-compose, written in Go with the gocui library.                                                                    |
| [lazygit](https://github.com/jesseduffield/lazygit)                          | A simple terminal UI for git commands.                                                                                                                            |
| [linkerd2](https://github.com/linkerd/linkerd2)                              | Ultralight, security-first service mesh for Kubernetes.                                                                                                           |
| [mc](https://github.com/minio/mc)                                            | MinIO Client is a replacement for ls, cp, mkdir, diff and rsync commands for filesystems and object storage.                                                      |
| [metal](https://github.com/equinix/metal-cli)                                | Official Equinix Metal CLI                                                                                                                                        |
| [minikube](https://github.com/kubernetes/minikube)                           | Runs the latest stable release of Kubernetes, with support for standard Kubernetes features.                                                                      |
| [mixctl](https://github.com/inlets/mixctl)                                   | A tiny TCP load-balancer.                                                                                                                                         |
| [mkcert](https://github.com/FiloSottile/mkcert)                              | A simple zero-config tool to make locally trusted development certificates with any names you'd like.                                                             |
| [nats](https://github.com/nats-io/natscli)                                   | Utility to interact with and manage NATS.                                                                                                                         |
| [nats-server](https://github.com/nats-io/nats-server)                        | Cloud native message bus and queue server                                                                                                                         |
| [nerdctl](https://github.com/containerd/nerdctl)                             | Docker-compatible CLI for containerd, with support for Compose                                                                                                    |
| [nova](https://github.com/FairwindsOps/nova)                                 | Find outdated or deprecated Helm charts running in your cluster.                                                                                                  |
| [oc](https://github.com/openshift/oc)                                        | Client to use an OpenShift 4.x cluster.                                                                                                                           |
| [oh-my-posh](https://github.com/jandedobbeleer/oh-my-posh)                   | A prompt theme engine for any shell that can display kubernetes information.                                                                                      |
| [op](https://github.com/1password/)                                          | 1Password CLI enables you to automate administrative tasks and securely provision secrets across development environments.                                        |
| [opa](https://github.com/open-policy-agent/opa)                              | General-purpose policy engine that enables unified, context-aware policy enforcement across the entire stack.                                                     |
| [openshift-install](https://github.com/openshift/installer)                  | CLI to install an OpenShift 4.x cluster.                                                                                                                          |
| [operator-sdk](https://github.com/operator-framework/operator-sdk)           | Operator SDK is a tool for scaffolding and generating code for building Kubernetes operators                                                                      |
| [osm](https://github.com/openservicemesh/osm)                                | Open Service Mesh uniformly manages, secures, and gets out-of-the-box observability features.                                                                     |
| [pack](https://github.com/buildpacks/pack)                                   | Build apps using Cloud Native Buildpacks.                                                                                                                         |
| [packer](https://github.com/hashicorp/packer)                                | Build identical machine images for multiple platforms from a single source configuration.                                                                         |
| [polaris](https://github.com/FairwindsOps/polaris)                           | Run checks to ensure Kubernetes pods and controllers are configured using best practices.                                                                         |
| [popeye](https://github.com/derailed/popeye)                                 | Scans live Kubernetes cluster and reports potential issues with deployed resources and configurations.                                                            |
| [porter](https://github.com/getporter/porter)                                | With Porter you can package your application artifact, tools, etc. as a bundle that can distribute and install.                                                   |
| [promtool](https://github.com/prometheus/prometheus)                         | Prometheus rule tester and debugging utility                                                                                                                      |
| [rclone](https://github.com/rclone/rclone)                                   | 'rsync for cloud storage' - Google Drive, S3, Dropbox, Backblaze B2, One Drive, Swift, Hubic, Wasabi, Google Cloud Storage, Azure Blob, Azure Files, Yandex Files |
| [regctl](https://github.com/regclient/regclient)                             | Utility for accessing docker registries                                                                                                                           |
| [rekor-cli](https://github.com/sigstore/rekor)                               | Secure Supply Chain - Transparency Log                                                                                                                            |
| [replicated](https://github.com/replicatedhq/replicated)                     | CLI for interacting with the Replicated Vendor API                                                                                                                |
| [rosa](https://github.com/openshift/rosa)                                    | Red Hat OpenShift on AWS (ROSA) command line tool                                                                                                                 |
| [rpk](https://github.com/redpanda-data/redpanda)                             | Kafka compatible streaming platform for mission critical workloads.                                                                                               |
| [run-job](https://github.com/alexellis/run-job)                              | Run a Kubernetes Job and get the logs when it's done.                                                                                                             |
| [scaleway-cli](https://github.com/scaleway/scaleway-cli)                     | Scaleway CLI is a tool to help you pilot your Scaleway infrastructure directly from your terminal.                                                                |
| [seaweedfs](https://github.com/seaweedfs/seaweedfs)                          | SeaweedFS is a fast distributed storage system for blobs, objects, files, and data lake, for billions of files!                                                   |
| [skupper](https://github.com/skupperproject/skupper)                         | Skupper is an implementation of a Virtual Application Network, enabling rich hybrid cloud communication                                                           |
| [snowmachine](https://github.com/rgee0/snowmachine)                          | Festive cheer for your terminal.                                                                                                                                  |
| [sops](https://github.com/getsops/sops)                                      | Simple and flexible tool for managing secrets                                                                                                                     |
| [stern](https://github.com/stern/stern)                                      | Multi pod and container log tailing for Kubernetes.                                                                                                               |
| [syft](https://github.com/anchore/syft)                                      | CLI tool and library for generating a Software Bill of Materials from container images and filesystems                                                            |
| [talosctl](https://github.com/siderolabs/talos)                              | The command-line tool for managing Talos Linux OS.                                                                                                                |
| [task](https://github.com/go-task/task)                                      | A simple task runner and build tool                                                                                                                               |
| [tctl](https://github.com/temporalio/tctl)                                   | Temporal CLI.                                                                                                                                                     |
| [terraform](https://github.com/hashicorp/terraform)                          | Infrastructure as Code for major cloud providers.                                                                                                                 |
| [terraform-docs](https://github.com/terraform-docs/terraform-docs)           | Generate documentation from Terraform modules in various output formats.                                                                                          |
| [terragrunt](https://github.com/gruntwork-io/terragrunt)                     | Terragrunt is a thin wrapper for Terraform that provides extra tools for working with multiple Terraform modules                                                  |
| [terrascan](https://github.com/tenable/terrascan)                            | Detect compliance and security violations across Infrastructure as Code.                                                                                          |
| [tflint](https://github.com/terraform-linters/tflint)                        | A Pluggable Terraform Linter.                                                                                                                                     |
| [tfsec](https://github.com/aquasecurity/tfsec)                               | Security scanner for your Terraform code                                                                                                                          |
| [tilt](https://github.com/tilt-dev/tilt)                                     | A multi-service dev environment for teams on Kubernetes.                                                                                                          |
| [timoni](https://github.com/stefanprodan/timoni)                             | A package manager for Kubernetes powered by CUE.                                                                                                                  |
| [tkn](https://github.com/tektoncd/cli)                                       | A CLI for interacting with Tekton.                                                                                                                                |
| [tofu](https://github.com/opentofu/opentofu)                                 | OpenTofu lets you declaratively manage your cloud infrastructure                                                                                                  |
| [trivy](https://github.com/aquasecurity/trivy)                               | Vulnerability Scanner for Containers and other Artifacts, Suitable for CI.                                                                                        |
| [vagrant](https://github.com/hashicorp/vagrant)                              | Tool for building and distributing development environments.                                                                                                      |
| [vault](https://github.com/hashicorp/vault)                                  | A tool for secrets management, encryption as a service, and privileged access management.                                                                         |
| [vcluster](https://github.com/loft-sh/vcluster)                              | Create fully functional virtual Kubernetes clusters - Each vcluster runs inside a namespace of the underlying k8s cluster.                                        |
| [vhs](https://github.com/charmbracelet/vhs)                                  | CLI for recording demos                                                                                                                                           |
| [viddy](https://github.com/sachaos/viddy)                                    | A modern watch command. Time machine and pager etc.                                                                                                               |
| [waypoint](https://github.com/hashicorp/waypoint)                            | Easy application deployment for Kubernetes and Amazon ECS                                                                                                         |
| [yq](https://github.com/mikefarah/yq)                                        | Portable command-line YAML processor.                                                                                                                             |
| [yt-dlp](https://github.com/yt-dlp/yt-dlp)                                   | Fork of youtube-dl with additional features and fixes                                                                                                             |
There are 162 tools, use `arkade get NAME` to download one.

> Note to contributors, run `go build && ./arkade get --format markdown` to generate this list