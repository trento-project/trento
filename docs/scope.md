# Table of contents
- [About this document](#about-this-document)
- [Goal](#goal)
- [Scope](#scope)
  - [Hosts](#hosts)
    - [Discovery](#discovery)
    - [Checks](#checks)
      - [HA Checker: Scenario SAP HANA Performance Optimized on Azure](#ha-checker-scenario-sap-hana-performance-optimized-on-azure)
      - [HA Checker: Scenario Pacemaker on Azure](#ha-checker-scenario-pacemaker-on-azure)
      - [More coming](#more-coming)
  - [Clusters](#clusters)
  - [Systems](#systems)
  - [Landscapes](#landscapes)
  - [Environments](#environments)
  - [Agents](#agents)

# About this document
This document aims to provide information about the current features that `trento`
tries to cover and its main scope. Further, additional details are provided on
the functionality of each view.

# Goal
>Provide simple to use front end for all relevant OS related tasks for SAP 
>workloads

This means that `trento` focuses on the OS-related tasks without stepping into 
specific functionality that is already covered by the relevant application layer,
centered around SAP workloads and simplicity.

# Scope
The scope of each component of `trento` is as follows:

## Hosts
In the hosts overview, the user can get a list of all the hosts that are running
the trento agent. For each host it will be possible to check basic information
such as the hostname, its IP address, to what cluster it belongs and a list of
all the tags that have been set by the discovery mechanisms to classify each
host accordingly. 

When accessing the details of each host, it will be possible to access in-depth
information about the role each host has in the cluster such as the status of the
core HDB processes, SAP System ID (SID), etc. Here, the SAP administrator can
get a list of potential problems and improvements on the configuration of the
host that the discovery mechanisms have found.

### Discovery
The discovery mechanisms are the base of the checkers. These implement the code
in the trento agent to gather information about the nodes, their roles, the
configuration files that each node uses (depending on their role).


### Checks
Under the checks view it's possible to see a representation on the data collected
by the agents against predefined values that are provided by our own recommendations
and our partners such as SAP. The base for the checks are the above described
[Discovery mechanisms](#discovery)

There are currently 2 checkers supported:

#### HA Checker: Scenario SAP HANA Performance Optimized on Azure
This checker uses Azure's recommendations to provide an overview of how well
the deployed SAP HANA cluster meets their indicators. More details about this
recommendation can be seen [here](https://docs.microsoft.com/en-us/azure/virtual-machines/workloads/sap/high-availability-guide-suse-pacemaker)

#### HA Checker: Scenario Pacemaker on Azure
This check is also comparing against Azure's recommendations for a Pacemaker
cluster.

#### More coming
Additional checkers are being implemented for other cloud providers.

## Clusters
This view currently allows to check the status of the all the discovered clusters,
including the information of the status of each node and the role associated to
each node. This view also allows to see the resources, their type and distribution
within the cluster nodes.

## Systems
In this view are listed the systems, usually identified by a SID or 
`SAP System Idenfication` string such as `PRD`, `DEV`, or `QAS`. These are used
often to identify productive, development and quality assurance systems.
In this view we get an overview of the distribution of each SID, the number of
hosts disvered of each of them and additional details.

## Landscapes
The landscapes view shows a list of the discovered landscapes, detailing the
count of systems inside each of them as well as the total number of hosts
that belong to these systems.

## Environments
The environments view currently shows a detail of the environments that have
been found and the landscapes that each one of these environments has as well
as the number of systems, hosts and the global status.

## Agents
Though not specifically a view/section, as it is the core component of `trento`,
it is important to understand what they are.

The agents are the `trento` processes that run in each of the nodes that conform
a highly available SAP Applications cluster. On each of these, the agents attempt
to discover the role of the node and run a set of checks which result in the
recommendations described above in the [Checkers](#checkers).
The information that is gathered by them is stored through a distributed KV store
and is made visible on the `trento` web UI through the `trento` web UI component.

The agents implement the core functionality of trento. They are responsible for
the discovery of all the clustered components that are required in order to run
highly available SAP Applications. These agents implement discovery for:
  - Pacemaker
  - Corosync (soon)
  - SAP components