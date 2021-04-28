# Trento Example Deployment Layout

    +------------------------------------------------------------------------------------------------------------------------------------------------------------+
    | Datacenter/Cloud                                                                                                                                           |
    |                     +------------------------------------------------------------------------------------------------------------------------------+       |
    |                     |SAP Environment                                                                                                               |       |
    |                     | +-----------------------------------------------------------------------------------------------------------------------+    |       |
    |                     | |SAP Landscape                                                                                                          |    |       |
    |                     | |                                                                                                                       |    |       |
    |                     | | +-----------------------------------------------------------------------------------------------------------------+   |    |       |
    | +--------------+    | | | SAP System                                                                                                      |   |    |       |
    | | Trento       |    | | | +--------------------------------------+   +--------------------------------------+                             |   |    |       |
    | |  +--------+  |    | | | | SAP HANA Cluster                     |   | SAP Application Patform Cluster      |                             |   |    |       |
    | |  |        |  |    | | | |  +--------------+   +--------------+ |   |  +--------------+   +--------------+ |      +--------------+       |   |    |       |
    | |  | Trento |  |     | |  |  |SAP Host      |   |SAP Host      | |   |  |SAP Host      |   |SAP Host      | |      |SAP Host      |       |   |    |       |
    | |  |  WebUI |  |    | | | |  | +---------+  |   | +---------+  | |   |  | +---------+  |   | +---------+  | |      | +---------+  |       |   |    |       |
    | |  |        +--+----+-+-+-+--+-+> Trento |  |   +-+> Trento |  | |   |  +-+> Trento |  |   +-+> Trento |  | |      +-+> Trento |  |       |   |    |       |
    | |  |   |    |  |    | | | |  | |   Agent |  |   | |   Agent |  | |   |  | |   Agent |  |   | |   Agent |  | |      | |   Agent |  |       |   |    |       |
    | |  +---+----+  |    | | | |  | |   |     |  |   | |   |     |  | |   |  | |   |     |  |   | |   |     |  | |      | |   |     |  |       |   |    |       |
    | |      |       |    | | | |  | +---+-----+  |   | +---+-----+  | |   |  | +---+-----+  |   | +---+-----+  | |      | +---+-----+  |       |   |    |       |
    | |  +---v-----+ |    | | | |  |     |        |   |     |        | |   |  |     |        |   |     |        | |      |     |        |       |   |    |       |
    | |  |Consul   | |    | | | |  | +---v-----+  |   | +---v-----+  | |   |  | +---v-----+  |   | +---v-----+  | |      | +---v-----+  |       |   |    |       |
    | |  | Agent   +-+----+-+-+-+--+->Consul   |  |   +->Consul   |  | |   |  +->Consul   |  |   +->Consul   |  | |      +->Consul   |  |       |   |    |       |
    | |  | (Server)<-+----+-+-+-+--+-+ Agent   |  |   +-+ Agent   |  | |   |  +-+ Agent   |  |   +-+ Agent   |  | |      +-+ Agent   |  |       |   |    |       |
    | |  |         | |    | | | |  | | (Client)|  |   | | (Client)|  | |   |  | | (Client)|  |   | | (Client)|  | |      | | (Client)|  |       |   |    |       |
    | |  +---------+ |    | | | |  | +---------+  |   | +---------+  | |   |  | +---------+  |   | +---------+  | |      | +---------+  |       |   |    |       |
    | |              |    | | | |  |              |   |              | |   |  |              |   |              | |      |              |       |   |    |       |
    | +--------------+    | | | |  +--------------+   +--------------+ |   |  +--------------+   +--------------+ |      +--------------+       |   |    |       |
    |                     | | | |                                      |   |                                      |                             |   |    |       |
    |                     | | | +--------------------------------------+   +--------------------------------------+                             |   |    |       |
    |                     | | |                                                                                                                 |   |    |       |
    |                     | | +-----------------------------------------------------------------------------------------------------------------+   |    |       |
    |                     | |                                                                                                                       |    |       |
    |                     | | +-----------------------------------------------------------------------------------------------------------------+   |    |       |
    |                     | | | SAP System                                                                                                      |   |    |       |
    |                     | | | +--------------------------------------+   +--------------------------------------+                             |   |    |       |
    |                     | | | | SAP HANA Cluster                     |   | SAP Application Patform Cluster      |                             |   |    |       |
    |                     | | | |  +--------------+   +--------------+ |   |  +--------------+   +--------------+ |      +--------------+       |   |    |       |
    |                     | | | |  |SAP Host      |   |SAP Host      | |   |  |SAP Host      |   |SAP Host      | |      |SAP Host      |       |   |    |       |
    |                     | | | |  | +---------+  |   | +---------+  | |   |  | +---------+  |   | +---------+  | |      | +---------+  |       |   |    |       |
    |                -----+-+-+-+--+-+> Trento |  |   +-+> Trento |  | |   |  +-+> Trento |  |   +-+> Trento |  | |      +-+> Trento |  |       |   |    |       |
    |                     | | | |  | |   Agent |  |   | |   Agent |  | |   |  | |   Agent |  |   | |   Agent |  | |      | |   Agent |  |       |   |    |       |
    |                     | | | |  | |   |     |  |   | |   |     |  | |   |  | |   |     |  |   | |   |     |  | |      | |   |     |  |       |   |    |       |
    |                     | | | |  | +---+-----+  |   | +---+-----+  | |   |  | +---+-----+  |   | +---+-----+  | |      | +---+-----+  |       |   |    |       |
    |                     | | | |  |     |        |   |     |        | |   |  |     |        |   |     |        | |      |     |        |       |   |    |       |
    |                     | | | |  | +---v-----+  |   | +---v-----+  | |   |  | +---v-----+  |   | +---v-----+  | |      | +---v-----+  |       |   |    |       |
    |                     | | | |  | |Consul   |  |   +->Consul   |  | |   |  +->Consul   |  |   +->Consul   |  | |      +->Consul   |  |       |   |    |       |
    |                     +-+-+-+--+-> Agent   |  |   +-+ Agent   |  | |   |  +-+ Agent   |  |   +-+ Agent   |  | |      +-+ Agent   |  |       |   |    |       |
    |                     | | | |  | | (Client)|  |   | | (Client)|  | |   |  | | (Client)|  |   | | (Client)|  | |      | | (Client)|  |       |   |    |       |
    |                     | | | |  | +---------+  |   | +---------+  | |   |  | +---------+  |   | +---------+  | |      | +---------+  |       |   |    |       |
    |                     | | | |  |              |   |              | |   |  |              |   |              | |      |              |       |   |    |       |
    |                     | | | |  +--------------+   +--------------+ |   |  +--------------+   +--------------+ |      +--------------+       |   |    |       |
    |                     | | | |                                      |   |                                      |                             |   |    |       |
    |                     | | | +--------------------------------------+   +--------------------------------------+                             |   |    |       |
    |                     | | |                                                                                                                 |   |    |       |
    |                     | | +-----------------------------------------------------------------------------------------------------------------+   |    |       |
    |                     | |                                                                                                                       |    |       |
    |                     | +-----------------------------------------------------------------------------------------------------------------------+    |       |
    |                     |                                                                                                                              |       |
    |                     | +-----------------------------------------------------------------------------------------------------------------------+    |       |
    |                     | | SAP Landscape                                                                                                         |    |       |
    |                     | |  +----------------------------------------------------------------------------------------------------------------+   |    |       |
    |                     | |  |SAP System                                                                                                      |   |    |       |
    |                     | |  |  +------------------------+                                                                                    |   |    |       |
    |                     | |  |  | SAP HANA Cluster       |                                                                                    |   |    |       |
    |                     | |  |  |                        |                                                                                    |   |    |       |
    |                     | |  |  | +--------+  +--------+ |      +--------+  +--------+                                                        |   |    |       |
    |                     | |  |  | |SAP Host|  |SAP Host| |      |SAP Host|  |SAP Host|                                                        |   |    |       |
    |                     | |  |  | |        |  |        | |      |        |  |        |                                                        |   |    |       |
    |                     | |  |  | |        |  |        | |      |        |  |        |                                                        |   |    |       |
    |                     | |  |  | +--------+  +--------+ |      +--------+  +--------+                                                        |   |    |       |
    |                     | |  |  |                        |                                                                                    |   |    |       |
    |                     | |  |  +------------------------+                                                                                    |   |    |       |
    |                     | |  |                                                                                                                |   |    |       |
    |                     | |  |                                                                                                                |   |    |       |
    |                     | |  +----------------------------------------------------------------------------------------------------------------+   |    |       |
    |                     | |                                                                                                                       |    |       |
    |                     | +-----------------------------------------------------------------------------------------------------------------------+    |       |
    |                     |                                                                                                                              |       |
    |                     +------------------------------------------------------------------------------------------------------------------------------+       |
    |                                                                                                                                                            |
    +------------------------------------------------------------------------------------------------------------------------------------------------------------+

## Trento Agent

The Trento Agent is to be set up on every node that should be handled by Trento,
and it needs to be combined with a Consul Agent. 

The Trento Agent will run several periodic job loops. 

### Discovery Loop

The Discovery loop is happening at a very low frequency and priority and is capturing record-worthy
Data in the Consul Key-Value Store. At the current implementation, only selected datapoints are
regularly updated into the Consul Key-Value Store. There is no conflict-resolution being implemented,
the data is unconditionally overwritten.

## Checker Loop

TBD


# Nomenclature

## SAP Host

This is a physical or a virtual server ("node") running a Linux Operating System on which SAP system components are installed and running on (potentially supervised
by cluster management software like for example pacemaker).

## SAP System

This uses a very traditional multi-tier client server setup:

* Presentation Layer: Web Browser, Dedicated GUI or 3rd Party Application between End-User and Application Servers. Out of control for Trento.

* Application Server Layer: Business applications which the end users interact with (can be written in ABAP programming language, for example)

* Database Server Layer: Contains the customer data and all the ABAP code which is run by the application servers. Most often this is SAP HANA.

## SAP Landscape

A SAP Landscape is one or more SAP Systems, usually combined in a Staging Systems Setup: DEV -> QA -> PRD

## SAP Environment

In cases when there are multiple SAP applications to manage, each business application can be deployed in its own SAP Landscape. It could also be
that the same business application is replicated (sharded) in multiple SAP Landscapes (one for each Business Unit for example). Often this is also referred to as a `SAP Landscape`.
