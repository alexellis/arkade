# arkade - The Open Source Kubernetes Marketplace

arkade provides a portable marketplace for downloading your favourite devops CLIs and installing helm charts, with a single command.

You can also download CLIs like `kubectl`, `kind`, `kubectx` and `helm` faster than you can type "apt-get/brew update".

<img src="docs/arkade-logo-sm.png" alt="arkade logo" width="150" height="150">

[![Build Status](https://travis-ci.com/alexellis/arkade.svg?branch=master)](https://travis-ci.com/alexellis/arkade)
[![GoDoc](https://godoc.org/github.com/alexellis/arkade?status.svg)](https://godoc.org/github.com/alexellis/arkade)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexellis/arkade)](https://goreportcard.com/report/github.com/alexellis/arkade)
![GitHub All Releases](https://img.shields.io/github/downloads/alexellis/arkade/total)

With over 40 helm charts and apps available for Kubernetes, gone are the days of contending with dozens of README files just to set up a development stack with the usual suspects like ingress-nginx, Postgres and cert-manager.

## Should you try arkade?

Here's what [Ivan Velichko](https://twitter.com/iximiuz/status/1422605221226860548?s=20), SRE @ Booking.com has to say about arkade:

> I was setting up a new dev environment yesterday. Kind, helm, kustomize, kubectl, all this stuff. My take is - arkade is highly underappreciated.
> I'd spend an hour in the past to install such tools. With arkade it was under ten minutes.

[Greg](https://twitter.com/cactusanddove) runs Fullstack JS and is a JavaScript developer, he says:

> This is real magic get #kubernetes up and going in a second; then launch #openfaas a free better than lambda solution that uses docker images.

[@arghzero](https://twitter.com/ArghZero/status/1346097288851070983?s=20) says:

> for getting the basics installed, nothing beats arkade
> it can install commonly used cli tools like kubectl locally for you, as well as common k8s pkgs like ingress-nginx or portainer

[@Yankexe](https://twitter.com/yankexe/status/1305427718050250754?s=20) says:

> It's hard to use K8s without Arkade these days. 
> My team at @lftechnology absolutely loves it. 

From [Michael Cade @ Kasten](https://twitter.com/MichaelCade1/status/1390403831167700995?s=20)

> I finally got around to installing Arkade, super simple! 
> quicker to install this than the argocli standalone commands, but there are lots of handy little tools in there.
> also, the neat little part about arkade, not only does it make it easy to install a ton of different apps and CLIs you can also get the info on them as well pretty quickly.

## Get arkade

```bash
# Note: you can also run without `sudo` and move the binary yourself
curl -sLS https://get.arkade.dev | sudo sh

arkade --help
ark --help  # a handy alias

# Windows users with Git Bash
curl -sLS https://get.arkade.dev | sh
```

> Windows users: arkade requires bash to be available, therefore Windows users can [install Git Bash](https://git-scm.com/downloads).

An alias of `ark` is created at installation time, so you can also run `ark install APP`

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
arkade get APP
```

|       TOOL       |                                                    DESCRIPTION                                                     |
|------------------|--------------------------------------------------------------------------------------------------------------------|
| argocd           | Declarative, GitOps continuous delivery tool for Kubernetes.                                                       |
| argocd-autopilot | An opinionated way of installing Argo-CD and managing GitOps repositories.                                         |
| arkade           | Portable marketplace for downloading your favourite devops CLIs and installing helm charts, with a single command. |
| buildx           | Docker CLI plugin for extended build capabilities with BuildKit.                                                   |
| civo             | CLI for interacting with your Civo resources.                                                                      |
| cosign           | Container Signing, Verification and Storage in an OCI registry.                                                    |
| docker-compose   | Define and run multi-container applications with Docker.                                                           |
| doctl            | Official command line interface for the DigitalOcean API.                                                          |
| faas-cli         | Official CLI for OpenFaaS.                                                                                         |
| flux             | Continuous Delivery solution for Kubernetes powered by GitOps Toolkit.                                             |
| gh               | GitHub’s official command line tool.                                                                               |
| helm             | The Kubernetes Package Manager: Think of it like apt/yum/homebrew for Kubernetes.                                  |
| helmfile         | Deploy Kubernetes Helm Charts                                                                                      |
| hugo             | Static HTML and CSS website generator.                                                                             |
| influx           | InfluxDB’s command line interface (influx) is an interactive shell for the HTTP API.                               |
| inlets-pro       | Cloud Native Tunnel for HTTP and TCP traffic.                                                                      |
| inletsctl        | Automates the task of creating an exit-server (tunnel server) on public cloud infrastructure.                      |
| istioctl         | Service Mesh to establish a programmable, application-aware network using the Envoy service proxy.                 |
| jq               | jq is a lightweight and flexible command-line JSON processor                                                       |
| k0s              | Zero Friction Kubernetes                                                                                           |
| k0sctl           | A bootstrapping and management tool for k0s clusters                                                               |
| k3d              | Helper to run Rancher Lab's k3s in Docker.                                                                         |
| k3sup            | Bootstrap Kubernetes with k3s over SSH < 1 min.                                                                    |
| k9s              | Provides a terminal UI to interact with your Kubernetes clusters.                                                  |
| kail             | Kubernetes log viewer.                                                                                             |
| kgctl            | A CLI to manage Kilo, a multi-cloud network overlay built on WireGuard and designed for Kubernetes.                |
| kim              | Build container images inside of Kubernetes. (Experimental)                                                        |
| kind             | Run local Kubernetes clusters using Docker container nodes.                                                        |
| kops             | Production Grade K8s Installation, Upgrades, and Management.                                                       |
| krew             | Package manager for kubectl plugins.                                                                               |
| kube-bench       | Checks whether Kubernetes is deployed securely by running the checks documented in the CIS Kubernetes Benchmark.   |
| kubebuilder      | Framework for building Kubernetes APIs using custom resource definitions (CRDs).                                   |
| kubectl          | Run commands against Kubernetes clusters                                                                           |
| kubectx          | Faster way to switch between clusters.                                                                             |
| kubens           | Switch between Kubernetes namespaces smoothly.                                                                     |
| kubeseal         | A Kubernetes controller and tool for one-way encrypted Secrets                                                     |
| kubetail         | Bash script to tail Kubernetes logs from multiple pods at the same time.                                           |
| kustomize        | Customization of kubernetes YAML configurations                                                                    |
| linkerd2         | Ultralight, security-first service mesh for Kubernetes.                                                            |
| mc               | MinIO Client is a replacement for ls, cp, mkdir, diff and rsync commands for filesystems and object storage.       |
| metal            | Official Equinix Metal CLI                                                                                         |
| minikube         | Runs the latest stable release of Kubernetes, with support for standard Kubernetes features.                       |
| nats             | Utility to interact with and manage NATS.                                                                          |
| nerdctl          | Docker-compatible CLI for containerd, with support for Compose                                                     |
| nova             | Find outdated or deprecated Helm charts running in your cluster.                                                   |
| opa              | General-purpose policy engine that enables unified, context-aware policy enforcement across the entire stack.      |
| osm              | Open Service Mesh uniformly manages, secures, and gets out-of-the-box observability features.                      |
| pack             | Build apps using Cloud Native Buildpacks.                                                                          |
| packer           | Build identical machine images for multiple platforms from a single source configuration.                          |
| polaris          | Run checks to ensure Kubernetes pods and controllers are configured using best practices.                          |
| popeye           | Scans live Kubernetes cluster and reports potential issues with deployed resources and configurations.             |
| porter           | With Porter you can package your application artifact, tools, etc. as a bundle that can distribute and install.    |
| rekor-cli        | Secure Supply Chain - Transparency Log                                                                             |
| stern            | Multi pod and container log tailing for Kubernetes.                                                                |
| terraform        | Infrastructure as Code for major cloud providers.                                                                  |
| tkn              | A CLI for interacting with Tekton.                                                                                 |
| trivy            | Vulnerability Scanner for Containers and other Artifacts, Suitable for CI.                                         |
| vagrant          | Tool for building and distributing development environments.                                                       |
| yq               | Portable command-line YAML processor.                                                                              |

> This is a time saver compared to searching for download pages every time you need a tool.

Think of `arkade get TOOL` as a doing for CLIs, what `arkade install` does for helm.

Adding a new tool for download is as simple as editing [tools.go](https://github.com/alexellis/arkade/blob/master/pkg/get/tools.go).

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
  -h, --help               help for ingress-nginx
      --host-mode          If we should install ingress-nginx in host mode.
  -n, --namespace string   The namespace used for installation (default "default")
      --update-repo        Update the helm repo (default true)
```

#### Override with `--set`

You can also set helm overrides, for apps which use helm via `--set`

```bash
ark install openfaas --set faasIdler.dryRun=false
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


#### Explore the apps

You can view the various apps available with `arkade install / --help`, more are available when you run the command yourself.

```bash
arkade install --help
ark --help

Examples:
  arkade install
  arkade install openfaas --helm3 --gateways=2
  arkade install inlets-operator --token-file $HOME/do-token

Available Commands:
  argocd                  Install argocd
  cassandra               Install cassandra
  cert-manager            Install cert-manager
  chart                   Install the specified helm chart
  consul-connect          Install Consul Service Mesh
  cron-connector          Install cron-connector for OpenFaaS
  crossplane              Install Crossplane
  docker-registry         Install a Docker registry
  docker-registry-ingress Install registry ingress with TLS
  falco                   Install Falco
  gitea                   Install gitea
  gitlab                  Install GitLab
  grafana                 Install grafana
  influxdb                Install influxdb
  info                    Find info about a Kubernetes app
  ingress-nginx           Install ingress-nginx
  inlets-operator         Install inlets-operator
  istio                   Install istio
  jenkins                 Install jenkins
  kafka                   Install Confluent Platform Kafka
  kafka-connector         Install kafka-connector for OpenFaaS
  kong-ingress            Install kong-ingress for OpenFaaS
  kube-image-prefetch     Install kube-image-prefetch
  kube-state-metrics      Install kube-state-metrics
  kubernetes-dashboard    Install kubernetes-dashboard
  kyverno                 Install Kyverno
  linkerd                 Install linkerd
  loki                    Install Loki for monitoring and tracing
  metrics-server          Install metrics-server
  minio                   Install minio
  mongodb                 Install mongodb
  mqtt-connector          Install mqtt-connector for OpenFaaS
  nats-connector          Install OpenFaaS connector for NATS
  nfs-client-provisioner  Install nfs client provisioner
  nginx-inc               Install nginx-inc for OpenFaaS
  opa-gatekeeper          Install Open Policy Agent (OPA) Gatekeeper
  openfaas                Install openfaas
  openfaas-ingress        Install openfaas ingress with TLS
  openfaas-loki           Install Loki-OpenFaaS and Configure Loki logs provider for OpenFaaS
  osm                     Install osm
  portainer               Install portainer to visualise and manage containers
  postgresql              Install postgresql
  rabbitmq                Install rabbitmq
  redis                   Install redis
  registry-creds          Install registry-creds
  sealed-secrets          Install sealed-secrets
  tekton                  Install Tekton pipelines and dashboard
  traefik2                Install traefik2
```

## Community & contributing

### Do you use this project? 👋

Alex created this project for developers just like yourself. If you use arkade, become a sponsor so that he can continue to grow and improve it for your future use.

<a href="https://github.com/sponsors/alexellis/">
<img alt="Sponsor this project" src="https://github.com/alexellis/alexellis/blob/master/sponsor-today.png" width="90%">
</a>

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

### Suggest a new app

To suggest a new app, please check past issues and [raise an issue for it](https://github.com/alexellis/arkade). Think also whether your app suggestion would be a good candidate for a Sponsored App.

### Sponsored apps

You can now propose your project or product as a Sponsored App. Sponsored Apps work just like any other app that we've curated, however they will have a note next to them in the app description `(sponsored)` and a link to your chosen site upon installation. An app sponsorship can be purchased for a minimum of 12 months and includes free development of the Sponsored App, with ongoing support via GitHub for the Sponsored App for the duration only. Ongoing support will be limited to a set amount of hours per month.

When your sponsorship expires the Sponsored App will be removed from arkade, and the ongoing support will cease. A Sponsored App can be renewed 60 days prior to expiration subject to a separate agreement and payment.

Example:

```bash
arkade VENDOR install PRODUCT
arkade acmeco install dashboard
```

Current sponsored apps include [Venafi](https://venafi.com) for Machine Identity:

```bash
arkade venafi install --help
arkade venafi info --help
```

[Contact us](mailto:contact@openfaas.com) to find out how you can have your Sponsored App added to arkade.

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

### Developer workflow

Here's the basics for contributing:

#### Cloning

```bash
mkdir $GOPATH/go/src/github.com/alexellis
git clone https://github.com/alexellis/arkade
cd arkade

go build
```

#### Adding your fork:

```bash
git remote add fork https://github.com/NAME/arkade
```

#### To verify changes:

```bash
gofmt -w -s ./pkg
gofmt -w -s ./cmd
go test ./...
```

#### Checkout a branch to start work

```bash
git checkout -b fork/add-NAME-of-APP
```

#### Push up your changes for a PR

```bash
git config user.name "Full name"
git config user.email "real@email.com"

git commit -s

git push fork add-NAME-of-APP
```

#### Test other people's PRs

You can also check out other people's PRs and test them:

```bash
arkade get gh
gh auth

$ gh pr list

Showing 10 of 10 open pull requests in alexellis/arkade

#477  Add comma escaping for --set flag        yankeexe:fix-openfaas-helm-set
#438  Added support to install kube-ps1        andreppires:master

gh checkout 477

go build

# Try the new version of arkade
./arkade
```

### Join us on Slack

Join `#contributors` at [slack.openfaas.io](https://slack.openfaas.io)

### License

MIT
