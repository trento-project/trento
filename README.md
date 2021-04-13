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

## Features

T.B.D.

## Requirements

To build the entire application you will need the following dependencies:

- Go ^1.16
- Node.js ^15.x

## Installation

This project is in development so, for the time being, you need to clone it and
build it manually:

```shell
git clone github.com/trento-project/trento.git
cd trento
make build
```

Pre-built binaries will be available soon.

## Usage

You can start the web application as follows:

```shell
./trento web serve
```

The web application needs one or more agents registered against against Consul
(https://consul.io/). A consul server needs to be paired with the Trento Web
Application.

Each agent node also needs a consul agent started. Follow the
[Running An Agent](https://www.consul.io/docs/agent#running-an-agent) steps for
starting it prior running the Trento Agent.

The Trento agent can then be started:

```shell
./trento agent start -n $name examples/azure-rules.yaml
```

See the [Deployment Architecture](./docs/trento-architecture.md) for details.

## Development

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

> **Note**  
> The [`mockery`](https://github.com/vektra/mockery) tool is required for the `generate` target, which in turn is required for the `test` target.
> You can install it with `go install github.com/vektra/mockery/v2@latest`

## Support

T.B.D.

## Contributing

T.B.D.

## License

Copyright 2021 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied. See the License for the
specific language governing permissions and limitations under the License.
