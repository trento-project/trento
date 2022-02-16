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
- [Introduction](#introduction)
- [Installation](#installation)
  - [Requirements](#requirements)
  - [Quick-Start installation](#quick-start-installation)
    - [Trento Server installation](#trento-server-installation)
    - [Trento Agent installation](#trento-agent-installation)
      - [Starting Trento Agent service](#starting-trento-agent-service)
  - [Manual installation](#manual-installation)
    - [Pre-built binaries](#pre-built-binaries)
    - [Compile from source](#compile-from-source)
  - [Docker images](#docker-images)
  - [RPM Packages](#rpm-packages)
  - [Helm chart](#helm-chart)
    - [Install K3S](#install-k3s)
    - [Install Helm and chart dependencies](#install-helm-and-chart-dependencies)
    - [Install the Trento Server Helm chart](#install-the-trento-server-helm-chart)
    - [Other Helm chart usage examples:](#other-helm-chart-usage-examples-)
  - [Manually running Trento](#manually-running-trento)
    - [Trento Agents](#trento-agents)
    - [Trento Runner](#trento-runner)
      - [Starting the Trento Runner](#starting-the-trento-runner)
    - [Trento Web UI](#trento-web-ui)
- [Configuration](#configuration)
- [Development](#development)
  - [Helm development chart](#helm-development-chart)
  - [Build system](#build-system)
  - [Development dependencies](#development-dependencies)
  - [Docker](#docker)
  - [Dump a scenario from a running cluster](#dump-a-scenario-from-a-running-cluster)
  - [SAPControl web service](#sapcontrol-web-service)
- [Support](#support)
- [Contributing](#contributing)
- [License](#license)

# Features

- Automated discovery of SAP HANA HA clusters;
- SAP Systems and Instances overview;
- Configuration validation for Pacemaker, Corosync, SBD, SAPHanaSR and other generic _SUSE Linux Enterprise for SAP Application_ OS settings (a.k.a. the _HA Config Checks_);
- Specific configuration audits for SAP HANA Scale-Up Performance-Optimized scenarios deployed on MS Azure cloud.

# Introduction

_Trento_ is a comprehensive monitoring solution made by two main components, the _Trento Server_ and the _Trento Agent_.

The _Trento Server_ is an independent, cloud-native, distributed system and should run on dedicated infrastructure resources. It is in turn composed by the following sub-systems:

- The `trento web` application;
- The `trento runner` worker;

The _Trento Agent_ is a single background process (`trento agent`) running in each host of the target infrastructure the user desires to monitor.

Please note that, except for the third-party ones like Ansible, all the components are embedded within one single `trento` binary.

See the [architecture document](./docs/trento-architecture.md) for additional details.

> Being the project in development, all of the above might be subject to change!

# Installation

## Requirements

The _Trento Server_ is intended to run in many ways, depending on users' already existing infrastructure, but it's designed to be cloud-native and OS agnostic.
As such, our default installation method provisions a minimal, single node, [K3S] Kubernetes cluster to run its various components in Linux containers.  
The suggested physical resources for running all the _Trento Server_ components are 2GB of RAM and 2 CPU cores.
The _Trento Server_ needs to reach the target infrastructure.

The _Trento Agent_ component, on the other hand, needs to interact with a number of low-level system components
which are part of the [SUSE Linux Enterprise Server for SAP Applications](https://www.suse.com/products/sles-for-sap/) Linux distribution.
These could in theory also be installed and configured on other distributions providing the same functionalities, but this use case is not within the scope of the active development.

In addition to that, the _Trento Agent_ also requires the [Prometheus node_exporter component](https://github.com/prometheus/node_exporter) to be running to collect host information for the monitoring functionality.

The resource footprint of the _Trento Agent_ should not impact the performance of the host it runs on.

## Quick-Start installation

Installation scripts are provided to automatically install and update the latest version of Trento.
Please follow the instructions in the given order.

### Trento Server installation

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

### Trento Agent installation

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

- `ssh address`: the address to which the trento-agent should be reachable for ssh connection by the runner for check execution.

- `trento server IP`: the address where Trento server can be reached.

You can pass these arguments as flags or env variables too:

```
curl -sfL https://raw.githubusercontent.com/trento-project/trento/main/install-agent.sh | sudo bash -s - --ssh-address=192.168.33.10 --server-ip=192.168.33.1
```

```
SSH_ADDRESS=192.168.33.10 SERVER_IP=192.168.33.1 sudo ./install-agent.sh
```

#### Starting Trento Agent service

The installation script does not start the agent automatically.

You can enable boot startup and launch it with systemd:

```
sudo systemctl enable --now trento-agent
```

Please, make sure the server is running before starting the agent.

That's it! You can now reach the Trento web UI and start using it.

## Manual installation

### Pre-built binaries

Pre-built statically linked binaries are made available via [GitHub releases](https://github.com/trento-project/trento/releases).

### Compile from source

You clone also clone and build it manually:

```shell
git clone https://github.com/trento-project/trento.git
cd trento
make build
```

See the section below to know more about the build dependencies.

## Docker images

T.B.D.

## RPM Packages

T.B.D.

## Helm chart

The [packaging/helm](packaging/helm) directory contains the Helm chart for installing Trento Server in a Kubernetes cluster.

### Install K3S

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

Please refer to the [K3S official documentation](https://rancher.com/docs/k3s/latest/en/installation/) for more information about the installation.

### Install Helm and chart dependencies

Install Helm:

```
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```

Please refer to the [Helm official documentation](https://helm.sh/docs/intro/install/) for more information about the installation.

### Install the Trento Server Helm chart

Add third-party Helm repositories:

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

Install chart dependencies:

```
cd packaging/helm/trento-server/
helm dependency update
```

The Runner component of Trento server needs ssh access to the agent nodes to perform the checks.
You need to pass a valid private key used for ssh authentication to the Helm chart, and it will be stored
in the K3s cluster as a secret.
Please refer to the [Trento Runner](#trento-runner) section for more information.

Install the Trento Server chart:

```
helm install trento . --set-file trento-runner.privateKey=/your/path/id_rsa_runner
```

or perform a rolling update:

```
helm upgrade trento . --set-file trento-runner.privateKey=/your/path/id_rsa_runner
```

Now you can connect to the web server via `http://localhost` and point the agents to the cluster IP address.

### Other Helm chart usage examples:

Use a different container image (e.g. the `rolling` one):

```
helm install trento . --set trento-web.image.tag="rolling" --set trento-runner.image.tag="rolling" --set-file trento-runner.privateKey=id_rsa_runner
```

Use a different container registry:

```
helm install trento . --set trento-web.image.repository="ghcr.io/myrepo/trento-web" --set trento-runner.image.repository="ghcr.io/myrepo/trento-runner" --set-file trento-runner.privateKey=id_rsa_runner
```

Please refer to the the subcharts `values.yaml` for an advanced usage.

## Manually running Trento

What follows are explicit instructions on how to run all the various components without any automated orchestration.

### Trento Agents

Trento Agents are responsible for discovering SAP systems, HA clusters and some additional data. These Agents need to run in the same systems hosting the HA
Cluster services, so running them in isolated environments (e.g. serverless,
containers, etc.) makes little sense, as they won't be able as the discovery mechanisms will not be able to report any host information.

> NOTE: Suggested installation instructions for SUSE-based distributions, adjust accordingly

Install and start `node_exporter`:

```shell
zypper in -y golang-github-prometheus-node_exporter
systemctl start prometheus-node_exporter
```

To start the trento agent:

```shell
./trento agent start
```

Alternatively, you can use the `trento-agent.service` from this repository and start it, which will start
`node_exporter` automatically as a dependency:
```shell
cp packaging/systemd/trento-agent.service /etc/systemd/system
systemctl daemon-reload
systemctl start trento-agent.service
```

> If the discovery loop is being executed too frequently, and this impacts the Web interface performance, the agent
> has the option to configure the discovery loop mechanism using the `--discovery-period` flag. Increasing this value improves the overall performance of the application

#### Publishing discovery data

Trento Agents publish discovery data to a Collector on Trento Server.

The communication can be made secure by enabling `mTLS` with `--enable-mtls`

See [this tutorial](https://www.digitalocean.com/community/tutorials/openssl-essentials-working-with-ssl-certificates-private-keys-and-csrs) for extra information about SSL Certificates.

#### Server

```
$> ./trento web serve [...] --enable-mtls --cert /path/to/certs/server-cert.pem --key /path/to/certs/server-key.pem --ca /path/to/certs/ca-cert.pem
```

#### Agent

```
$> ./trento agent start [...] --enable-mtls --cert /path/to/certs/client-cert.pem --key /path/to/certs/client-key.pem --ca /path/to/certs/ca-cert.pem
```

**Development Note:** `./test/certs/` folder contains some dummy Server, Client and CA Certificates and Keys.
Those are useful in order to test `mTLS` communication between the Agent and the DataCollector.

---

### Trento Runner

The Trento Runner is a worker process responsible for driving automated configuration audits. It is based on [Ansible](https://docs.ansible.com/ansible/latest/index.html).
This component can be executed in the same host as the Web UI, but it is not mandatory: it can be executed in any other host with network access to the Trento Agents.

Find more information about how to create more Trento health checks [here](docs/runner.md).

In order to start them, some packages must be installed and started. Here a quick go through:

#### Starting the Trento Runner

The Runner needs the `ansible` Python package available locally:

```shell
pip install 'ansible~=4.6.0'
```

> The installed ansible components versions should be at least ansible~=4.6.0 and ansible-core~=2.11.5

Once dependencies are in place, you can start the Runner itself:

```shell
./trento runner start --api-host $WEB_IP --api-port $WEB_PORT -i 5
```

> _Note:_ The Trento Runner component must have SSH access to all the agents via a password-less SSH key pair.

### Trento Web UI

At this point, we can start the web application as follows:

```shell
./trento web serve
```

Please consult the `help` CLI command for more insights on the various options.

# Configuration

Trento can be run with a config file in replacement of command-line arguments.

## Locations

Configuration, if not otherwise specified by the `--config=/path/to/config.yaml` option, would be searched in following locations:

Note that order represents priority

- `/etc/trento/` <-- first location looked
- `/usr/etc/trento/` <-- fallback here if config not found in previous location
- `~/.config/trento/` aka user's home <-- fallback here

## Formats

`yaml` is the only supported format at the moment.

## Naming conventions

Each component of trento supports its own configuration, so the expected config files in the chosen location must be called after the component it is supporting.

So supported names are `(agent|web|runner).yaml`

Example locations:

`/etc/trento/agent.yaml`

`/etc/trento/web.yaml`

`/etc/trento/runner.yaml`

or

`/usr/etc/trento/agent.yaml`

`/usr/etc/trento/web.yaml`

`/usr/etc/trento/runner.yaml`

**Note**: `runner` still not supported for now

## Examples

```
# /etc/trento/agent.yaml

enable-mtls: true
cert: /path/to/certs/client-cert.pem
key: /path/to/certs/client-key.pem
ca: /path/to/certs/ca-cert.pem
collector-host: localhost
collector-port: 8443
```

```
# /etc/trento/web.yaml

db-port: 5432

collector-port: 8443
enable-mtls: true
cert: /path/to/certs/server-cert.pem
key: /path/to/certs/server-key.pem
ca: /path/to/certs/ca-cert.pem
```

## Environment Variables

All of the options supported by the command line and configuration file can be provided as environment variables as well.

The rule is: get the option name eg. `enable-mtls`, replace dashes `-` with underscores `_`, make it uppercase and add a `TRENTO_` prefix.

Examples:

`enable-mtls` -> `TRENTO_ENABLE_MTLS=true ./trento agent start`

`collector-host` -> `TRENTO_COLLECTOR_HOST=localhost ./trento web serve`

`cert` -> `TRENTO_CERT=/path/to/certs/server-cert.pem ./trento web serve`

# Development

## Helm development chart

A development Helm chart is available at [./hack/helm/trento-dev](./hack/helm/trento-dev).
The chart is based on the official Helm chart package and overrides certain values to provide a development environment.

```shell
# Update dependencies in the official Helm chart
cd ./packaging/helm/trento-server
helm dep update

# Install the development Helm chart
cd ./hack/helm/trento-dev
helm dep update
helm install trento-dev .
```

Since integration tests require a running PostgreSQL instance, please make sure the chart is installed prior to running the integration tests.
The PostgreSQL instance will be accessible at port `localhost:5432`.

If you want to skip database integration tests, you can use a dedicated environment variable as follows:

```shell
TRENTO_DB_INTEGRATION_TESTS=false make test
```

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

## Development dependencies

Additionally, for the development we use [`mockery`](https://github.com/vektra/mockery) for the `generate` target, which in turn is required for the `test` target.
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

## End-to-end testing
End-to-end testing inside Trento is achieved through [Cypress](https://cypress.io), a browser testing tool that can instrument both Chromium/Puppeteer and Firefox. It can be used also as an integration test framework since assertions can be run on JSON payloads too.

Setting the environment up to write or run end-to-end tests is quite simple. Just change directory to `test/e2e`, install the dependencies and open Cypress:

```sh
cd ./test/e2e/

npm install

npx cypress run # run it command line
npx cypress open # open the GUI tool to launch tests
```

### Adjusting environment variables to suit your own needs
Sometimes it can be useful to change some environment variable just to adjust some behavior to your own needs. Also, this will become handy if you need any new environment variable around.

There are two ways to do that. The first is to edit the `cypress.json` file and modify the `env` field:

```json
{
  "baseUrl": "http://localhost:8080",
  "viewportWidth": 1366,
  "viewportHeight": 768,
  "env": {
    "collector_host": "localhost",
    "collector_port": 8081,
    "heartbeat_interval": 5000,
    "db_host": "localhost",
    "db_port": 5432,
    "fixtures_path": "./cypress/fixtures",
    "trento_binary": "../../trento",
    "photofinish_binary": "photofinish"
  }
}
```

The second way to provide Cypress an environment variable is to state it as a proper environment variable, prepending `cypress_` to that:

```sh
cypress_foo=bar npx cypress run
```

### Further questions about E2E testing?
If you have any doubts, the team is here to help! Just make sure you also covered [the official Cypress documentation](https://docs.cypress.io/) first. :smile:

## Dump a scenario from a running cluster

A script to dump a scenario from a running cluster is available at [./hack/dump-scenario.sh](./hack/dump_scenario_from_k8s.sh).

Running this script in the k3s node where trento-server was installed will dump the current state of Trento to a local folder.

Please refer to the script usage help for more details:

```bash
./hack/dump-scenario.sh --help
```

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

[k3s]: https://k3s.io
