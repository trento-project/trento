# Tagging the systems

In order to group and filter the systems a tagging mechanism can be used. This tags are placed as
meta-data in the agent nodes. Find information about how to set meta-data in the agents at: https://www.consul.io/docs/agent/options#node_meta

As an example, check the [meta-data file](./examples/trento-config.json) file. This file must be
located in the folder set as `-config-dir` during the agent execution.

The next items are reserved:
- `trento-ha-cluster`: Cluster which the system belongs to
- `trento-sap-environment`: Environment in which the system is running
- `trento-sap-landscape`: Landscape in which the system is running
- `trento-sap-environment`: SAP system (composed by database and application) in which the system is running

## Setting the tags from the KV storage

These reserved tags can be automatically set and updated using the [consul-template](https://github.com/hashicorp/consul-template).
To achieve this, the tags information will come from the KV storage.

The metadata will be set in the next paths:
- `trento/v0/hosts/$(hostname)/metadata/ha-cluster`
- `trento/v0/hosts/$(hostname)/metadata/sap-environment`
- `trento/v0/hosts/$(hostname)/metadata/sap-landscape`
- `trento/v0/hosts/$(hostname)/metadata/sap-system`

Notice that a new entry will be created for each node.

`consul-template` starts directly with the `trento` agent. It provides some configuration options to synchronize the utility with the consul agent.

- `config-dir`: Consul agent configuration directory; defaults to `consul.d`. The `consul-template` component running within `trento agent` will create the Trento node meta-data configuration file there. (i.e. `consul.d/trento-config.json`).
- `consul-template`: Template used to populate the trento meta-data configuration file (by default [meta-data file][./examples/trento-config.json] is used).
