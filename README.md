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
  * [Build dependencies](#build-dependencies)
  * [Runtime dependencies](#runtime-dependencies)
  * [Development dependencies](#development-dependencies)
- [Installation](#installation)
- [Running `trento`](#running-trento)
  * [Consul](#consul)
  * [Trento agents](#trento-agents)
  * [Web server](#web-server)  
- [Usage](#usage)
  * [Tagging the systems](#tagging-the-systems)
- [Development](#development)
  * [Build system](#build-system)
  * [Mockery](#mockery)
- [Support](#support)
- [Contributing](#contributing)
- [License](#license)
  
# Features

T.B.D.

# Requirements

## Build dependencies

To build the entire application you will need the following dependencies:

- [`Go`](https://golang.org/) ^1.16
- [`Node.js`](https://nodejs.org/es/) ^15.x

## Runtime dependencies

Running the application will require:
  - A running [`consul`](https://www.consul.io/downloads) cluster. 

>We have only tested version `1.9.x` and while it *should* work with any consul agents that implement consul protocol version 3, we can´t guarantee it at the moment.

## Development dependencies

Additionally, for the development we use:
  - [`Mockery`](https://github.com/vektra/mockery) ^2

> See [Development section](#development) for details on how to configure `mockery`
# Installation

This project is in development so, for the time being, you need to clone it and
build it manually:

```shell
git clone https://github.com/trento-project/trento.git
cd trento
make build
```

Pre-built binaries will be available soon.

# Running trento

To run trento in our development environment we need at least:
  - Consul agent in server mode
  - Consul agent in client mode
  - Trento agent
  - Trento web server

## Consul

The web application needs one or more agents registered against against 
[Consul](https://consul.io/). A consul server needs to be paired with the 
Trento Web Application.

Each [Trento agent node](##trento-agents) also needs a consul agent started. 
Follow the [Running An Agent](https://www.consul.io/docs/agent#running-an-agent)
more detailed steps for starting it prior running the Trento Agent.

To start the consul agent as server:

```shell
./consul agent -server -bootstrap-expect=1 -bind=127.0.0.1 -data-dir=consul-data -ui
```

This will start consul, binding to `127.0.0.1` and use `consul-data` as 
directory to persist data.

Another agent in client mode is also required. For development purposes, to be
able to run both on the same host, we need to bind it to a new `lo` address:

```shell
sudo ip address add 127.0.0.2/32 dev lo
```

Create the directory for the consul agent client:
```shell
mkdir consul.d
```

Now we can start the agent:
```shell
./consul agent -node=test -data-dir=consul-agent-data -bind=127.0.0.2 -client=127.0.0.2 -retry-join=127.0.0.1 -ui -config-dir=./consul.d/test
```


> Production deployments require multiple instances for the consult agents. Be 
> sure to check [consul's deployment guide](https://learn.hashicorp.com/tutorials/consul/deployment-guide#configure-consul-agents)

## Trento agents

To start the trento agent:

```shell
./trento agent start -n $name examples/azure-rules.yaml
```

> Note that we are using `azure-rules.yaml` in this example which collect azure
recommendations on cluster settings and state for many HA components. New rules
can be created and tuned to adjust to different requirements.

> See the [Deployment Architecture](./docs/trento-architecture.md) for additional
details.

## Web server

At this point, we can start the web application as follows:

```shell
./trento web serve [flags]
```

The supported flags are as follows:
```shell
  -h, --help          help for serve
      --host string   The host to bind the HTTP service to (default "0.0.0.0")
  -p, --port int      The port for the HTTP service to listen at (default 8080)
```

Additionally, `trento` supports the next Global Flags:

```shell
Global Flags:
      --config string   config file (default is $HOME/.trento.yaml)
```

# Usage
## Tagging the systems

In order to group and filter the systems a tagging mechanism can be used. This tags are placed as
meta-data in the agent nodes. Find information about how to set meta-data in the agents at: https://www.consul.io/docs/agent/options#node_meta

As an example, check the [meta-data file](./examples/trento-config.json) file. This file must be
located in the folder set as `-config-dir` during the agent execution.

The next items are reserved:
- `trento-ha-cluster`: Cluster which the system belongs to
- `trento-sap-environment`: Environment in which the system is running
- `trento-sap-landscape`: Landscape in which the system is running
- `trento-sap-environment`: SAP system (composed by database and application) in which the system is running

### Setting the tags from the KV storage

These reserved tags can be automatically set and updated using the [consul-template](https://github.com/hashicorp/consul-template).
To achieve this, the tags information will come from the KV storage.

Set the metadata in the next paths:
- `trento/nodename/metadata/ha-cluster`
- `trento/nodename/metadata/sap-environment`
- `trento/nodename/metadata/sap-landscape`
- `trento/nodename/metadata/sap-system`

Notice that a new entry must exists for every node.

`consul-template` starts directly with the `trento` agent. It provides some configuration options to synchronize the utility with the consul agent.

- `config-dir`: Consul agent configuration files directory. It must be the same used by the consul agent. The `trento` agent creates a new folder with the node name where the trento meta-data configuration file is stored (e.g. `consul.d/node1/trento-config.json`).
- `consul-template`: Template used to populate the trento meta-data configuration file (by default [meta-data file][./examples/trento-config.json] is used).

### Filtering the nodes in the wep app

The web app provides the option to filter the systems using the previously commented reserved tags. To achieve this, the tags must be stored in the KV storage.
Use the next path:
- `trento/filters/sap-environments`
- `trento/filters/sap-landscapes`
- `trento/filters/sap-systems`

Each of them must have a json list format. As example: `["land1", "land2"]`.
These entries will be available in the filters on the `/environments` page.

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
You can install it with `go install github.com/vektra/mockery/v2@latest`.

> **Note**  
> Be sure to add the `mockery` binary to your `$PATH` environment variable so that `make` can find it.

# Support

As the project is currently in its early stages, we suggest that any question or
issue is directed to our [Issues](https://github.com/trento-project/trento/issues) 
section in GitHub.

# Contributing

T.B.D.

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
