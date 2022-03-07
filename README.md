# arkade - The Open Source Kubernetes Marketplace

arkade provides a portable marketplace for downloading your favourite devops CLIs and installing helm charts, with a single command.

You can also download CLIs like `kubectl`, `kind`, `kubectx` and `helm` faster than you can type "apt-get/brew update".

<img src="docs/arkade-logo-sm.png" alt="arkade logo" width="150" height="150">

[![Sponsor this](https://img.shields.io/static/v1?label=Sponsor&message=%E2%9D%A4&logo=GitHub&link=https://github.com/sponsors/alexellis)](https://github.com/sponsors/alexellis)
[![build](https://github.com/alexellis/arkade/actions/workflows/build.yml/badge.svg)](https://github.com/alexellis/arkade/actions/workflows/build.yml)
[![GoDoc](https://godoc.org/github.com/alexellis/arkade?status.svg)](https://godoc.org/github.com/alexellis/arkade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexellis/arkade)](https://goreportcard.com/report/github.com/alexellis/arkade)
![Downloads](https://img.shields.io/github/downloads/alexellis/arkade/total)

With over 52 helm charts and apps available for Kubernetes, gone are the days of contending with dozens of README files just to set up a development stack with the usual suspects like ingress-nginx, Postgres and cert-manager.

- [arkade - The Open Source Kubernetes Marketplace](#arkade---the-open-source-kubernetes-marketplace)
  - [Should you try arkade?](#should-you-try-arkade)
  - [Getting arkade](#getting-arkade)
  - [Usage overview](#usage-overview)
  - [Download CLI tools with arkade](#download-cli-tools-with-arkade)
  - [Installing apps with arkade](#installing-apps-with-arkade)
  - [Community & contributing](#community--contributing)
  - [Sponsored apps](#sponsored-apps)
  - [FAQ](#faq)

## Should you try arkade?

> I was setting up a new dev environment yesterday. Kind, helm, kustomize, kubectl, all this stuff. My take is - arkade is highly underappreciated.
> I'd spend an hour in the past to install such tools. With arkade it was under ten minutes.
>
> [Ivan Velichko](https://twitter.com/iximiuz/status/1422605221226860548?s=20), SRE @ Booking.com

> This is real magic get #kubernetes up and going in a second; then launch #openfaas a free better than lambda solution that uses docker images.
>
> [Greg](https://twitter.com/cactusanddove) runs Fullstack JS and is a JavaScript developer

> for getting the basics installed, nothing beats arkade
> it can install commonly used cli tools like kubectl locally for you, as well as common k8s pkgs like ingress-nginx or portainer
>
> [@arghzero](https://twitter.com/ArghZero/status/1346097288851070983?s=20)


> It's hard to use K8s without Arkade these days.
> My team at @lftechnology absolutely loves it.
>
> [@Yankexe](https://twitter.com/yankexe/status/1305427718050250754?s=20)

> I finally got around to installing Arkade, super simple!
> quicker to install this than the argocli standalone commands, but there are lots of handy little tools in there.
> also, the neat little part about arkade, not only does it make it easy to install a ton of different apps and CLIs you can also get the info on them as well pretty quickly.
> 
> [Michael Cade @ Kasten](https://twitter.com/MichaelCade1/status/1390403831167700995?s=20)

> You've to install latest and greatest tools for your daily @kubernetesio tasks? No problem, check out #arkade the open source #kubernetes marketplace 👍
>
> [Thorsten Hans](https://twitter.com/ThorstenHans/status/1457982292597608449?s=20) - Cloud Native consultant

> If you want to install quickly a new tool in your dev env or in your k8s cluster you can use the Arkade (https://github.com/alexellis/arkade) easy and quick you should try out! Ps. I contribute to this project 🥰
>
> [Carlos Panato](https://twitter.com/comedordexis/status/1423339283713347587) - Staff engineer @ Mattermost

> arkade is the 'brew install' of Kubernetes. You can install and run an application in a single command. Finally! https://github.com/alexellis/arkade /by Alex Ellis
>
> [John Arundel](https://twitter.com/bitfield/status/1242385165445455872?s=20) - Cloud consultant, author

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
* `arkade update` - update arkade itself

An arkade "app" could represent a helm chart such as `openfaas/faas-netes`, a custom CLI installer such as `istioctl` or a set of static manifests (i.e. MetalLB).

An arkade "tool" is a CLI that can be downloaded for your operating system. Arkade downloads statically-linked binaries from their upstream locations on GitHub or a the vendor's chosen URL such as with `kubectl` and `terraform`.

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


Think of `arkade get TOOL` as a doing for CLIs, what `arkade install` does for helm.

Adding a new tool for download is as simple as editing [tools.go](https://github.com/alexellis/arkade/blob/master/pkg/get/tools.go).

[Click here for the full catalog of CLIs](#catalog-of-apps)

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

There's even more options you can chose with `arkade install openfaas --help` - the various flags you see map to settings from the helm chart README, that you'd usually have to look up and set via a `values.yaml` file.

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

### Do you use this project? 👋

<a href="https://github.com/sponsors/alexellis/">Alex</a> created this project for developers just like yourself. If you use arkade, become a sponsor so that he can continue to grow and improve it for your future use.

<a href="https://github.com/sponsors/alexellis/">
<img alt="Sponsor this project" src="https://github.com/alexellis/alexellis/blob/master/sponsor-today.png" width="90%">
</a>

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

Current sponsored apps include:

* [Venafi](https://venafi.com) for Machine Identity:

  ```bash
  arkade venafi install --help
  arkade venafi info --help
  ```

* [Kasten (k10)](https://kasten.io) for Kubernetes backup and restore:

  ```bash
  arkade kasten install --help
  arkade install k10 \
    --set eula.accept=true \
    --set clusterName=my-k10 \
    --set prometheus.server.enabled=false
  arkade kasten info k10
  ```

[Contact us](mailto:contact@openfaas.com) to find out how you can have your Sponsored App added to arkade.

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
| istio                   | Install istio                                                       |
| docker-registry-ingress | Install registry ingress with TLS                                   |
| ingress-nginx           | Install ingress-nginx                                               |
| openfaas-ingress        | Install openfaas ingress with TLS                                   |
| docker-registry         | Install a Docker registry                                           |
| registry-creds          | Install registry-creds                                              |
| sealed-secret           | Install sealed-secrets                                              |
| influxdb                | Install influxdb                                                    |
| mongodb                 | Install mongodb                                                     |
| kube-image-prefetch     | Install kube-image-prefetch                                         |
| gitea                   | Install gitea                                                       |
| mqtt-connector          | Install mqtt-connector for OpenFaaS                                 |
| kafka                   | Install Confluent Platform Kafka                                    |
| cert-manager            | Install cert-manager                                                |
| loki                    | Install Loki for monitoring and tracing                             |
| redis                   | Install redis                                                       |
| argocd                  | Install argocd                                                      |
| falco                   | Install Falco                                                       |
| rabbitmq                | Install rabbitmq                                                    |
| waypoint                | Install Waypoint                                                    |
| chart                   | Install the specified helm chart                                    |
| metrics-server          | Install metrics-server                                              |
| linkerd                 | Install linkerd                                                     |
| cron-connector          | Install cron-connector for OpenFaaS                                 |
| kafka-connector         | Install kafka-connector for OpenFaaS                                |
| crossplane              | Install Crossplane                                                  |
| minio                   | Install minio                                                       |
| openfaas                | Install openfaas                                                    |
| portainer               | Install portainer to visualise and manage containers                |
| OSM                     | Install osm                                                         |
| kong-ingress            | Install kong-ingress for OpenFaaS                                   |
| opa-gatekeeper          | Install Open Policy Agent (OPA) Gatekeeper                          |
| cockroachdb             | Install CockroachDB                                                 |
| prometheus              | Install Prometheus for monitoring                                   |
| kanister                | Install kanister for application-level data management              |
| kube-state-metrics      | Install kube-state-metrics                                          |
| kubernetes-dashboard    | Install kubernetes-dashboard                                        |
| postgresql              | Install postgresql                                                  |
| kyverno                 | Install Kyverno                                                     |
| openfaas-loki           | Install Loki-OpenFaaS and Configure Loki logs provider for OpenFaaS |
| traefik2                | Install traefik2                                                    |
| inlets-operator         | Install inlets-operator                                             |
| nginx-inc               | Install nginx-inc for OpenFaaS                                      |
| cassandra               | Install cassandra                                                   |
| metallb-arp             | Install MetalLB in L2 (ARP) mode                                    |
| grafana                 | Install grafana                                                     |
| nfs-provisioner         | Install nfs client provisioner                                      |
| gitlab                  | Install GitLab                                                      |
| nats-connector          | Install OpenFaaS connector for NATS                                 |
| jenkins                 | Install jenkins                                                     |
| tekton                  | Install Tekton pipelines and dashboard                              |
| consul-connect          | Install Consul Service Mesh                                         |

There are 52 apps that you can install on your cluster.

Note to contributors, run `arkade install --print-table` to generate this list

### Catalog of CLIs


|       TOOL       |                                                                DESCRIPTION                                                                |
|------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| argocd           | Declarative, GitOps continuous delivery tool for Kubernetes.                                                                              |
| argocd-autopilot | An opinionated way of installing Argo-CD and managing GitOps repositories.                                                                |
| arkade           | Portable marketplace for downloading your favourite devops CLIs and installing helm charts, with a single command.                        |
| autok3s          | Run Rancher Lab's lightweight Kubernetes distribution k3s everywhere.                                                                     |
| buildx           | Docker CLI plugin for extended build capabilities with BuildKit.                                                                          |
| civo             | CLI for interacting with your Civo resources.                                                                                             |
| clusterctl       | The clusterctl CLI tool handles the lifecycle of a Cluster API management cluster                                                         |
| cosign           | Container Signing, Verification and Storage in an OCI registry.                                                                           |
| devspace         | Automate your deployment workflow with DevSpace and develop software directly inside Kubernetes.                                          |
| dive             | A tool for exploring each layer in a docker image                                                                                         |
| docker-compose   | Define and run multi-container applications with Docker.                                                                                  |
| doctl            | Official command line interface for the DigitalOcean API.                                                                                 |
| eksctl           | Amazon EKS Kubernetes cluster management                                                                                                  |
| faas-cli         | Official CLI for OpenFaaS.                                                                                                                |
| flux             | Continuous Delivery solution for Kubernetes powered by GitOps Toolkit.                                                                    |
| gh               | GitHub’s official command line tool.                                                                                                      |
| goreleaser       | Deliver Go binaries as fast and easily as possible                                                                                        |
| helm             | The Kubernetes Package Manager: Think of it like apt/yum/homebrew for Kubernetes.                                                         |
| helmfile         | Deploy Kubernetes Helm Charts                                                                                                             |
| hostctl          | Dev tool to manage /etc/hosts like a pro!                                                                                                 |
| hugo             | Static HTML and CSS website generator.                                                                                                    |
| influx           | InfluxDB’s command line interface (influx) is an interactive shell for the HTTP API.                                                      |
| inlets-pro       | Cloud Native Tunnel for HTTP and TCP traffic.                                                                                             |
| inletsctl        | Automates the task of creating an exit-server (tunnel server) on public cloud infrastructure.                                             |
| istioctl         | Service Mesh to establish a programmable, application-aware network using the Envoy service proxy.                                        |
| jq               | jq is a lightweight and flexible command-line JSON processor                                                                              |
| k0s              | Zero Friction Kubernetes                                                                                                                  |
| k0sctl           | A bootstrapping and management tool for k0s clusters                                                                                      |
| k10multicluster  | Multi-cluster support for K10.                                                                                                            |
| k10tools         | Tools for evaluating and debugging K10.                                                                                                   |
| k3d              | Helper to run Rancher Lab's k3s in Docker.                                                                                                |
| k3s              | Lightweight Kubernetes                                                                                                                    |
| k3sup            | Bootstrap Kubernetes with k3s over SSH < 1 min.                                                                                           |
| k9s              | Provides a terminal UI to interact with your Kubernetes clusters.                                                                         |
| kail             | Kubernetes log viewer.                                                                                                                    |
| kanctl           | Framework for application-level data management on Kubernetes.                                                                            |
| kgctl            | A CLI to manage Kilo, a multi-cloud network overlay built on WireGuard and designed for Kubernetes.                                       |
| kim              | Build container images inside of Kubernetes. (Experimental)                                                                               |
| kind             | Run local Kubernetes clusters using Docker container nodes.                                                                               |
| kops             | Production Grade K8s Installation, Upgrades, and Management.                                                                              |
| krew             | Package manager for kubectl plugins.                                                                                                      |
| kube-bench       | Checks whether Kubernetes is deployed securely by running the checks documented in the CIS Kubernetes Benchmark.                          |
| kubebuilder      | Framework for building Kubernetes APIs using custom resource definitions (CRDs).                                                          |
| kubecm           | Easier management of kubeconfig.                                                                                                          |
| kubectl          | Run commands against Kubernetes clusters                                                                                                  |
| kubectx          | Faster way to switch between clusters.                                                                                                    |
| kubens           | Switch between Kubernetes namespaces smoothly.                                                                                            |
| kubescape        | kubescape is the first tool for testing if Kubernetes is deployed securely as defined in Kubernetes Hardening Guidance by to NSA and CISA |
| kubeseal         | A Kubernetes controller and tool for one-way encrypted Secrets                                                                            |
| kubestr          | Kubestr discovers, validates and evaluates your Kubernetes storage options.                                                               |
| kubetail         | Bash script to tail Kubernetes logs from multiple pods at the same time.                                                                  |
| kustomize        | Customization of kubernetes YAML configurations                                                                                           |
| linkerd2         | Ultralight, security-first service mesh for Kubernetes.                                                                                   |
| mc               | MinIO Client is a replacement for ls, cp, mkdir, diff and rsync commands for filesystems and object storage.                              |
| metal            | Official Equinix Metal CLI                                                                                                                |
| minikube         | Runs the latest stable release of Kubernetes, with support for standard Kubernetes features.                                              |
| mkcert           | A simple zero-config tool to make locally trusted development certificates with any names you'd like.                                     |
| nats             | Utility to interact with and manage NATS.                                                                                                 |
| nerdctl          | Docker-compatible CLI for containerd, with support for Compose                                                                            |
| nova             | Find outdated or deprecated Helm charts running in your cluster.                                                                          |
| opa              | General-purpose policy engine that enables unified, context-aware policy enforcement across the entire stack.                             |
| osm              | Open Service Mesh uniformly manages, secures, and gets out-of-the-box observability features.                                             |
| pack             | Build apps using Cloud Native Buildpacks.                                                                                                 |
| packer           | Build identical machine images for multiple platforms from a single source configuration.                                                 |
| polaris          | Run checks to ensure Kubernetes pods and controllers are configured using best practices.                                                 |
| popeye           | Scans live Kubernetes cluster and reports potential issues with deployed resources and configurations.                                    |
| porter           | With Porter you can package your application artifact, tools, etc. as a bundle that can distribute and install.                           |
| rekor-cli        | Secure Supply Chain - Transparency Log                                                                                                    |
| stern            | Multi pod and container log tailing for Kubernetes.                                                                                       |
| terraform        | Infrastructure as Code for major cloud providers.                                                                                         |
| tfsec            | Security scanner for your Terraform code                                                                                                  |
| tilt             | A multi-service dev environment for teams on Kubernetes.                                                                                  |
| tkn              | A CLI for interacting with Tekton.                                                                                                        |
| trivy            | Vulnerability Scanner for Containers and other Artifacts, Suitable for CI.                                                                |
| vagrant          | Tool for building and distributing development environments.                                                                              |
| vcluster         | Create fully functional virtual Kubernetes clusters - Each vcluster runs inside a namespace of the underlying k8s cluster.                |
| waypoint         | Easy application deployment for Kubernetes and Amazon ECS                                                                                 |
| yq               | Portable command-line YAML processor.                                                                                                     |

Note to contributors, run `arkade get --output markdown` to generate this list

