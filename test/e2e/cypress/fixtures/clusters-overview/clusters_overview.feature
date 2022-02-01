Feature: Pacemaker Clusters Overview
    This is where the user has an overview of the status of all the Pacemaker clusters registered with trento

    Background:
        Given an healthy SAP deployment of 9 pacemaker clusters having the following cluster IDs and names:
    # 9c832998801e28cd70ad77380e82a5c0 => { name: "hana_cluster", type: "HANA scale-up", status: "Passing" }
    # 8bca366a6cb7816555538092a1ddd5aa => { name: "netweaver_cluster",
    # 04b8f8c21f9fd8991224478e8c4362f8 => { name: "hana_cluster", type: "HANA scale-up", status: "Critical" }
    # a034a158905404befe08775682910ee1 => { name: "drbd_cluster", type: "Unknown" }
    # 238a4de1239aae2aa87433eed788b3ad => { name: "drbd_cluster", type: "Unknown" }
    # 04a81f89c847e82390e35bece2e25c9b => { name: "drbd_cluster", type: "Unknown" }
    # acf59e7a5338f76f55d5055af3273480 => { name: "netweaver_cluster", type: "Unknown" }
    # 057f083c3be591f4398eed816d4c8cd7 => { name: "netweaver_cluster", type: "Unknown" }
    # 4e905d706da85f5be14f85fa947c1e39 => { name: "hana_cluster", type: "HANA scale-up", status: "Warning" }

    Scenario: Registered Clusters are shown in the list
        When I navigate to the Pacemaker Clusters Overview (/clusters)
        Then the displayed clusters should be the ones listed above

    Scenario: Health Container information matches the status of the listed clusters
        Given I am in the Pacemaker Clusters Overview
        When the health container is ready
        Then there should 1 items in Passing status
        And there should be 1 items in Warning status
        And there should be 1 items in Critical status

    Scenario: Discovered Clusters in the paginated list (10 items) are reporting their status correctly
        Given I am in the Hosts Overview
        And the listing shows 10 items per page
        Then the cluster with id '04b8f8c21f9fd8991224478e8c4362f8' is in Critical status
        And the cluster with id '4e905d706da85f5be14f85fa947c1e39' is in Warning status
        And the cluster with id '9c832998801e28cd70ad77380e82a5c0' is in Passing status
        And all other clusters are in Unknown status

    Scenario: Filtering the Clusters Overview by Health
        Given I am in the Hosts Overview
        When I filter by Health Passing
        Then the cluster with id '9c832998801e28cd70ad77380e82a5c0' should be displayed

        When I filter by Health Warning
        Then the cluster with id '4e905d706da85f5be14f85fa947c1e39' should be displayed

        When I filter by Health Critical
        Then the cluster with id '04b8f8c21f9fd8991224478e8c4362f8' should be displayed

    Scenario: Filtering the Clusters Overview by Cluster name
        Given I am in the Hosts Overview
        When I filter by SAP system HDD
        Then 1 items should be displayed

        When I filter by SAP system HDP
        Then 1 items should be displayed

        When I filter by SAP system HDQ
        Then 1 items should be displayed

    Scenario: Filtering the Clusters Overview by SAP System
        Given I am in the Hosts Overview
        When I filter by SAP system HDD
        Then 1 items should be displayed

        When I filter by SAP system HDP
        Then 1 items should be displayed

        When I filter by SAP system HDQ
        Then 1 items should be displayed

    Scenario: Filtering the Host Overview by Tags
        Given all the hosts containing 'hana' in their name are tagged with 'env1'
        And all the hosts containing 'drbd' in their name are tagged with 'env2'
        And all the hosts containing 'netweaver' in their name are tagged with 'env3'
        When I filter by tag 'tag1'
        Then 1 items should be shown

        When I filter by tag 'tag2'
        Then 1 items should be shown

        When I filter by tag 'tag3'
        Then 1 items should be shown