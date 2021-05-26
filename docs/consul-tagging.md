# Tagging the systems

In order to group and filter the various Hosts into a hierarchy of Systems, Landscapes and Environments we leverage the Consul tagging mechanism. These tags are stored as
node meta-data in a static Consul configuration file ([read more in the Consul docs](https://www.consul.io/docs/agent/options#node_meta)).

The Trento Agent uses the `--consul-config-dir` flag to set the path where Consul auto-loads configuration files: the node meta-data will be stored there as `trento-metadata.json`.

The following keys will be used for the node meta-data:
- `trento-ha-cluster`: Cluster which the system belongs to
- `trento-sap-environment`: Environment in which the system is running
- `trento-sap-landscape`: Landscape in which the system is running
- `trento-sap-environment`: SAP system (composed by database and application) in which the system is running

These tags are automatically set and updated by [consul-template](https://github.com/hashicorp/consul-template), which we run embedded within our Trento Agent.

The metadata will be set under the following Consul KV store paths:
- `trento/v0/hosts/$(hostname)/metadata/ha-cluster`
- `trento/v0/hosts/$(hostname)/metadata/sap-environment`
- `trento/v0/hosts/$(hostname)/metadata/sap-landscape`
- `trento/v0/hosts/$(hostname)/metadata/sap-system`

Note that a new entry will be created for each node.
