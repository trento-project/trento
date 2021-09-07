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
- [Installation](#installation)
- [Running Trento](#running-trento)
  - [Consul](#consul)
  - [Trento Agents](#trento-agents)
  - [Trento Web UI](#trento-web-ui)
- [Usage](#usage)
  - [Grouping and filtering in the Web UI](#grouping-and-filtering-in-the-web-ui)
- [Development](#development)
  - [Build system](#build-system)
  - [Mockery](#mockery)
- [Support](#support)
- [Contributing](#contributing)
- [License](#license)

# Features

- Automated discovery of SAP HANA HA clusters;
- SAP Systems and Instances overview;
- Grouping by Landscapes and Environments;
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

# Installation

## From binaries

Pre-built statically linked binaries are made available via [GitHub releases](https://github.com/trento-project/trento/releases).

## Manual

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

Trento Agents are responsible for discovering HA clusters. These Agents need to run in the same systems hosting the HA
Cluster services, so running them in isolated environments (e.g. serverless,
containers, etc.) makes little sense, as they won't be able as the discovery mechanisms will not be able to report any host information.

To start the trento agent:

```shell
./trento agent start
```

## Trento Web UI

At this point, we can start the web application as follows:

```shell
./trento web serve
```

Please consult the `help` CLI command for more insights on the various options.

# Usage

## Grouping and filtering in the Web UI

The web app provides the option to group and filter target Systems using tags.

This functionality will be implemented in the Web UI soon but, for the time being, the underlying `consul` KV storage can be used to populate these tags:

```shell
# To group SAP Systems into Landscapes and Environments, you can do the following,
# from any of the nodes where Trento is running:
export ENV=an-environment-name
export LAND=a-landscape-name
export SID=an-actual-SID
consul kv put trento/v0/environments/$ENV/name $ENV
consul kv put trento/v0/environments/$ENV/landscapes/$LAND/name $LAND
consul kv put trento/v0/environments/$ENV/landscapes/$LAND/sapsystems/$SID/name $SID
```

Keep in mind that the created Environments, Landscapes and SAP Systems are directories themselves, and there can be multiple of them.

It is possible to have multiple Landscapes with the same name in different Environments, and the same goes for SAP Systems.

Be aware that the node meta-data tags are not strictly linked to these names, they are soft relations (this means that only the string matches, there is no any real relationship between them).

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

To assist in testing & developing Trento, we have added a [Dockerfile](Dockerfile)
that will automatically fetch all the required compile-time dependencies to build
the binary and finally a container image with the fresh binary in it.

We also provide a [docker-compose.yml](docker-compose.yml) file that allows to
deploy other required services to run alongside `trento` by fetching
the images from the [dockerhub](https://hub.docker.com/) registry and running
the containers with your favourite container engine.

If you want to build & start `trento web` and it's dependencies, you can use `docker-compose`:

```shell
git clone https://github.com/trento-project/trento.git
cd trento
docker-compose up
```

The Web UI should then be reachable as defined in the `docker-compose.yml` file
(`localhost:8080` by default).

To only build the docker image, instead:

```shell
docker build -t trento .
```

> Please note that the `trento agent` component requires to be running on
> the OS (_not_ inside a container) so, while it is possible to hack the `docker-compose.yml`
> file to also run a Trento Agent, it makes little sense because most of its internals
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
