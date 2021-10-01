# Trento

An open cloud-native web console improving on the life of SAP Applications
administrators.

_Trento is a city on the Adige River in Trentino-Alto Adige/Südtirol in Italy.
[...] It is one of the nation's wealthiest and most prosperous cities, [...]
often ranking highly among Italian cities for quality of life, standard of
living, and business and job opportunities._ ([source](https://en.wikipedia.org/wiki/Trento))

This project is a reboot of the "SUSE Console for SAP Applications", also known
as the [Blue Horizon for SAP](https://github.com/SUSE/blue-horizon-for-sap)
prototype, which is focused on automated infrastructure deployment and
provisioning for SAP Applications.

As opposed to that first iteration, this new one will focus more on operations
of existing clusters, rather than deploying new one.

# Table of contents

- [Features](#features)
- [Requirements](#requirements)
  - [Runtime dependencies](#runtime-dependencies)
  - [Build dependencies](#build-dependencies)
  - [Development dependencies](#development-dependencies)
- [Quick-start installation](#quick-start-installation)
  - [Trento Server installation](#trento-server-installation)
  - [Trento Agent installation](#trento-agent-installation)
- [Manual installation](#manual-installation)
  - [Pre-built binaries](#pre-built-binaries)
  - [Compile from source](#compile-from-source)
  - [Helm chart](#helm-chart)
- [Running Trento](#running-trento)
  - [Consul](#consul)
  - [Trento Agents](#trento-agents)
  - [Trento Runner](#trento-runner)
  - [Trento Web UI](#trento-web-ui)
- [Development](#development)
  - [Build system](#build-system)
  - [Mockery](#mockery)
- [Support](#support)
- [Contributing](#contributing)
- [License](#license)

# Features

- Automated discovery of SAP HANA HA clusters;
- SAP Systems and Instances overview;
- Configuration validation for Pacemaker, Corosync, SBD, SAPHanaSR and other generic _SUSE Linux Enterprise for SAP Application_ OS settings (a.k.a. the _HA Checks_);
- Specific configuration audits for SAP HANA Scale-Up Performance-Optimized scenarios deployed on MS Azure cloud.

# Base concepts

The entire Trento application is composed of the following parts:

- One or more Consul Agents in server mode;
- The Trento Web UI (`trento web`);
- A Consul Agent in client mode for each target node;
- A Trento Agent (`trento agent`) for each target node.

> See the [architecture document](./docs/trento-architecture.md) for additional details.

# Requirements

While the `trento web` component has only been tested on openSUSE 15.2
and SLES 15SP2 so far, it should be able to run on most modern Linux distributions.

The `trento agent` component could in theory also run on openSUSE, but it does not make much sense as it
needs to interact with different low-level SAP applications components
which are expected to be run in a
[SLES for SAP](https://www.suse.com/products/sles-for-sap/) installation.

## Runtime dependencies

Running the application will require:

- A running [Consul](https://www.consul.io/downloads) cluster.

> We have only tested version Consul version `1.9.x` and, while it _should_ work with any version implementing Consul Protocol version 3, we can´t make any guarantee in that regard.

## Build dependencies

To build the entire application you will need the following dependencies:

- [`Go`](https://golang.org/) ^1.16
- [`Node.js`](https://nodejs.org/es/) ^15.x

## Development dependencies

Additionally, for the development we use:

- [`Mockery`](https://github.com/vektra/mockery) ^2

> See the [Development](#development) section for details on how to install `mockery`.

# Quick-Start installation

Installation scripts are provided to automatically install and update the latest version of Trento.
Please follow the installation in the given order

## Trento Server installation

The script installs a single node K3s cluster and uses the [trento-server Helm chart](packaging/helm/trento-server)
to bootstrap a complete Trento server component.

You can `curl | bash` if you want to live on the edge.

```
curl -sfL https://raw.githubusercontent.com/trento-project/trento/main/install-server.sh | bash
```

Or you can fetch the script, and then execute it manually.

```
curl -O https://raw.githubusercontent.com/trento-project/trento/main/install-server.sh
chmod 700 install-server.sh
sudo ./install-server.sh
```

The script will ask you for a private key that is used by the runner service to perform checks in the agent hosts via ssh.

_Note: if a Trento server is already installed in the host, it will be updated._

Please refer to the [Trento Runner](#trento-runner) section for more information.
Please refer to the [Helm chart](#helm-chart) section for more information about the Helm chart.

## Trento Agent installation

After the server installation, you might want to install Trento agents in a running cluster.
Please add the public key to the ssh authorized_keys to enable the runner checks in the agent host, 
as mentioned in the server installation above.

As for the server component an installation script is provided, 
you can `curl | bash` it if you want to live on the edge.

```
curl -sfL https://raw.githubusercontent.com/trento-project/trento/main/install-agent.sh | sudo bash
```

Or you can fetch the script, and then execute it manually.

```
curl -O https://raw.githubusercontent.com/trento-project/trento/main/install-agent.sh
chmod 700 install-agent.sh
sudo ./install-agent.sh
```

The script will ask you for two IP addresses.

- `agent bind IP`: the private address to which the trento-agent should be bound for internal communications.
  This is an IP address that should be reachable by the other hosts, including the trento server.
  Note for the Pacemaker users: this IP address _should not be_ a floating IP.

- `trento server IP`: the address where Trento server can be reached.

You can pass these arguments as flags or env variables too:

```
curl -sfL https://raw.githubusercontent.com/trento-project/trento/main/install-agent.sh | sudo bash -s - --agent-bind-ip=192.168.33.10 --server-ip=192.168.33.1
```

```
AGENT_BIND_IP=192.168.33.10 SERVER_IP=192.168.33.1 sudo ./install-agent.sh
```

### Start Trento Agent

The installation script does not start the agent automatically.

You can start it by simply:

```
sudo systemctl start trento-agent
```

Please make sure the server is running before starting the agent

To enable the service execute:

```
sudo systemctl enable trento-agent
```

# Manual Installation

## Pre-built binaries

Pre-built statically linked binaries are made available via [GitHub releases](https://github.com/trento-project/trento/releases).

## Compile from source

You clone also clone and build it manually:

```shell
git clone https://github.com/trento-project/trento.git
cd trento
make build
```

## Docker images

T.B.D.

## RPM Packages

T.B.D.

## Helm chart

The [packaging/helm](packaging/helm) directory contains the Helm chart for installing Trento Server in a Kubernetes cluster.

### Install K3s

If installing as root:

```
# curl -sfL https://get.k3s.io | sh
```

If installing as non-root user:

```
curl -sfL https://get.k3s.io | sh -s - --write-kubeconfig-mode 644
```

Export KUBECONFIG env variable:

```
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
```

Please refer to the [K3s official documentation](https://rancher.com/docs/k3s/latest/en/installation/) for more information about the installation.

### Install Helm and chart dependencies

Install Helm:

```
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```

Please refer to the [Helm official documentation](https://helm.sh/docs/intro/install/) for more information about the installation.

### Install the Trento Server Helm chart

Add HashiCorp Helm repository:

```
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
```

Install chart dependencies:

```
cd packaging/helm/trento-server/
helm dependency update
```

The runner component of Trento server needs ssh access to the agent nodes to perform the checks.
You need to pass a valid private key used for ssh authentication to the Helm chart, and it will be stored
in the K3s cluster as a secret.
Please refer to the [Trento Runner](#trento-runner) section for more information.

Install Trento server chart:

```
helm install trento . --set-file trento-runner.privateKey=/your/path/id_rsa_runner
```

or perform a rolling update:

```
helm upgrade trento . --set-file trento-runner.privateKey=/your/path/id_rsa_runner
```

Now you can connect to the web server via `http://localhost` and point the agents to the cluster IP address.

### Other Helm chart usage examples:

Use a different container image:

```
helm install trento . --set trento-web.tag="runner" --set trento-runner.tag="runner" --set-file trento-runner.privateKey=id_rsa_runner
```

Use a different container registry:

```
helm install trento . --set trento-web.image.repository="ghcr.io/myrepo/trento" --set trento-runner.image.repository="ghcr.io/myrepo/trento" --set-file trento-runner.privateKey=id_rsa_runner
```

Please refer to the the subcharts `values.yaml` for an advanced usage.

# Running Trento

## Consul

The Trento application needs to be paired with a [Consul](https://consul.io/) deployment, which is leveraged for service discovery and persistent data storage.

Consul processes are all called "agents", and they can run either in server mode or in client mode.

#### Server Consul Agent

To start a Consul agent in server mode:

```shell
mkdir consul.d
./consul agent -server -ui -bootstrap-expect=1 -client=0.0.0.0 -data-dir=consul-data -config-dir=consul.d
```

This will start Consul listening to connections from any IP address on the default network interface, using `consul-data` as directory to persist data, and `consul.d` to autoload configuration files from.

#### Client Consul Agent

Each [Trento Agent](#trento-agents) instance also needs a Consul agent in client mode, on each target node we want to connect Trento to.

You can start Consul in client mode as follows:

```shell
export SERVER_IP_ADDRESS=#your Consul server IP address here
mkdir consul.d
./consul agent -retry-join=$SERVER_IP_ADDRESS -bind='{{ GetInterfaceIP "eth0" }}' -data-dir=consul-agent-data -config-dir=consul.d
```

Since the client Consul Agent will most likely run on a machine with multiple IP addresses and/or network interfaces, you will need to specify one with the `-bind` flag.

> Production deployments will require at least three instances of the Consul agent in server mode to ensure fault-tolerance. Be
> sure to check [Consul's deployment guide](https://learn.hashicorp.com/tutorials/consul/deployment-guide#configure-consul-agents).

> While Consul provides a `-dev` flag to run a standalone, stateless server agent, Trento does not support this mode: it needs a persistent server even during development.

> For development purposes, when running everything on a single host machine, no Client Consul Agent is required: the Server Agent exposes the same API and can be consumed by both `trento web` and `trento agent`.

## Trento Agents

Trento Agents are responsible for discovering SAP systems, HA clusters and some additional data. These Agents need to run in the same systems hosting the HA
Cluster services, so running them in isolated environments (e.g. serverless,
containers, etc.) makes little sense, as they won't be able as the discovery mechanisms will not be able to report any host information.

To start the trento agent:

```shell
./trento agent start
```

> If the discovery loop is being executed too frequently, and this impacts the Web interface performance, the agent
> has the option to configure the discovery loop mechanism using the `--discovery-period` flag. Increasing this value improves the overall performance of the application

## Trento Runner

The Trento Runner is responsible for running the health checks. It is based on [Ansible](https://docs.ansible.com/ansible/latest/index.html) and [ARA](https://ara.recordsansible.org/).
These 2 components (the Runner and ARA) can be executed in the same machine as the Web UI, but it is not mandatory, they can be executed in any other machine that has network access to the agents (the Runner and ARA can be even executed in different machines too, as long as the network connection is available between them).

In order to start them, some packages must be installed and started. Here a quick go through:

### ARA server

```shell
# Install ARA with server dependencies
pip install "ara[server]"
# Setup ARA database
ara-manage migrate
# Start ARA server. This process can be started in background or in other shell terminal
ara-manage runserver ip:port
```

If the requests to ARA server fail with a message like the next one, it means that the server address must be allowed:

```
2021-09-02 07:13:48,715 ERROR django.security.DisallowedHost: Invalid HTTP_HOST header: '10.74.1.5:8000'. You may need to add '10.74.1.5' to ALLOWED_HOSTS.
2021-09-02 07:13:48,732 WARNING django.request: Bad Request: /api/
```

To fix it run:

```
export ARA_ALLOWED_HOSTS="['10.74.1.5']"
# Or allow all the addresses with
export ARA_ALLOWED_HOSTS=['*']
```

### Runner

```shell
# Install ansible and ARA. Tested with version 2.11.2
pip install ansible ara

# Start the Trento runner
./trento runner start --ara-server http://araIP:port --consul-addr consulIP:port -i 5
# Find additional help with the -h flag
./trento runner start -h
```

**In order to use the runner component, the machine running it must have ssh authorization to all the
agents with a passwordless ssh key pair. Otherwise, the checks result is set as unreachable.**

## Trento Web UI

At this point, we can start the web application as follows:

```shell
./trento web serve
# If ARA server is not running in the same machine set the ara-addr flag
./trento web serve --ara-addr araIP:port
```

Please consult the `help` CLI command for more insights on the various options.

# Development

## Build system

We use GNU Make as a task manager; here are some common targets:

```shell
make # clean, test and build everything

make clean # removes any build artifact
make test # executes all the tests
make fmt # fixes code formatting
make web-assets # invokes the frontend build scripts
make generate # refresh automatically generated code (e.g. static Go mocks)
```

Feel free to peek at the [Makefile](Makefile) to know more.

## Mockery

As stated above, we use [`mockery`](https://github.com/vektra/mockery) for the `generate` target, which in turn is required for the `test` target.
You can install it with `go install github.com/vektra/mockery/v2`.

> Be sure to add the `mockery` binary to your `$PATH` environment variable so that `make` can find it. That usually comes with configuring `$GOPATH`, `$GOBIN`, and adding the latter to your `$PATH`.

## Docker

The [Dockerfile](Dockerfile) will automatically fetch all the required compile-time dependencies to build
the binary and finally a container image with the fresh binary in it.

We use a multi-stage build with two main targets: `trento-runner` and `trento-web`. The latter is the default.

You can build the component like follows:

```shell
docker build --target trento-runner -t trento-runner .
docker build -t trento-web . # same as specifying --target trento-web
```

> Please note that the `trento agent` component requires to be running on
> the OS (_not_ inside a container) so, while it is technically possible to run `trento agent` 
> commands in the container, it makes little sense because most of its internals
> require direct access to the host of the HA Cluster components.

## SAPControl web service

The SAPControl web service soap client was generated by [hooklift/gowsdl](https://github.com/hooklift/gowsdl),
then the methods and structs needed were cherry-picked and adapted.
For reference, you can find the full, generated, web service code [here](docs/_generated_soap_wsdl.go).

# Support

Please only report bugs via [GitHub issues](https://github.com/trento-project/trento/issues);
for any other inquiry or topic use [GitHub discussion](https://github.com/trento-project/trento/discussions).

# Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

# License

Copyright 2021 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
