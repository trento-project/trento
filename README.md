# Trento

_Trento is a city on the Adige River in Trentino-Alto Adige/Südtirol in Italy. [...] It is one of the nation's wealthiest and most prosperous cities, with its province being one of the richest in Italy._

A cloud-native, web application aiming at bringing high quality of life to SAP Applications administrators.

This project is a reboot of the "SUSE Console for SAP Applications", also known as the [Blue Horizon for SAP](https://github.com/SUSE/blue-horizon-for-sap) prototype, which is focused on automated infrastructure deployment and provisioning for SAP Applications.

As opposed to that first iteration, this new one will focus more on operations of existing clusters, rather than deploying new one.

## Features

T.B.D.

## Requirements

To build the entire application you will need the following dependencies:

- Go ^1.16
- Node.js ^15.x

## Installation

This project is in development so, for the time being, you need to clone it and build it manually: 

```shell
git clone github.com/SUSE/console-for-sap-applications
cd console-for-sap
make
```

Pre-built binaries will be available soon.

## Usage

You can start the web application as follows:

```shell
console-for-sap webapp serve
```

## Development

We use GNU Make as a task manager; here are some common targets:
```shell
make clean
make test
make fmt
make webapp-assets
make build
```

Feel free to peek at the [Makefile](Makefile) to know more.

## Support

T.B.D.

## Contributing

T.B.D.

## License

Copyright 2021 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
