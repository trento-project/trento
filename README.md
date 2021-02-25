# SUSE Console for SAP Applications v2

A cloud-native, web application to manage OS-related tasks for SAP Applications.

This project is a reboot of the "SUSE Console for SAP Applications", also known as the [Blue Horizon for SAP](https://github.com/SUSE/blue-horizon-for-sap) prototype, which is focused on automated infrastructure deployment and provisioning for SAP Applications.

As opposed to that first iteration, this new one will focus more on operations of existing clusters, rather than deploying new one.

## Features

T.B.D.

## Requirements

- Go ^1.16
- Go Modules (`export GO111MODULES=on`)

## Installation

You can get the project binary via `go get` as usual.

```shell
go get -u github.com/SUSE/console-for-sap
```

## Usage

You can start the web application as follows:

```shell
console-for-sap webapp serve
```

## Development

To build the entire application you will need the following additional dependencies:
- Node.js ^15.x
- [SASS](https://sass-lang.com/)

We provide a Makefile to build the entire app, including the frontend assets:
```shell
make
```

## Support

T.B.D.

## Contributing

T.B.D.

## License

Copyright 2021 SUSE LLC

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
