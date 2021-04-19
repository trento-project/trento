#!/bin/sh

cat <<EOF
<crm_mon version="2.0.4">
  <summary>
    <stack type="corosync"/>
    <current_dc present="true" version="2.0.4+20200616.2deceaa3a-3.3.1-2.0.4+20200616.2deceaa3a" name="hana01" id="1084787210" with_quorum="true"/>
    <last_update time="Mon Apr 19 22:07:33 2021"/>
    <last_change time="Mon Apr 19 22:07:16 2021" user="root" client="crm_attribute" origin="hana01"/>
    <nodes_configured number="2"/>
    <resources_configured number="7" disabled="0" blocked="0"/>
    <cluster_options stonith-enabled="true" symmetric-cluster="true" no-quorum-policy="stop" maintenance-mode="false"/>
  </summary>
  <nodes>
    <node name="hana01" id="1084787210" online="true" standby="false" standby_onfail="false" maintenance="false" pending="false" unclean="false" shutdown="false" expected_up="true" is_dc="true" resources_running="5" type="member"/>
    <node name="hana02" id="1084787211" online="true" standby="false" standby_onfail="false" maintenance="false" pending="false" unclean="false" shutdown="false" expected_up="true" is_dc="false" resources_running="2" type="member"/>
  </nodes>
  <resources>
    <resource id="stonith-sbd" resource_agent="stonith:external/sbd" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
      <node name="hana01" id="1084787210" cached="true"/>
    </resource>
    <resource id="rsc_ip_PRD_HDB00" resource_agent="ocf::heartbeat:IPaddr2" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
      <node name="hana01" id="1084787210" cached="true"/>
    </resource>
    <resource id="rsc_exporter_PRD_HDB00" resource_agent="systemd:prometheus-hanadb_exporter@PRD_HDB00" role="Started" target_role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
      <node name="hana01" id="1084787210" cached="true"/>
    </resource>
    <clone id="msl_SAPHana_PRD_HDB00" multi_state="true" unique="false" managed="true" failed="false" failure_ignored="false">
      <resource id="rsc_SAPHana_PRD_HDB00" resource_agent="ocf::suse:SAPHana" role="Master" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
        <node name="hana01" id="1084787210" cached="true"/>
      </resource>
      <resource id="rsc_SAPHana_PRD_HDB00" resource_agent="ocf::suse:SAPHana" role="Slave" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
        <node name="hana02" id="1084787211" cached="true"/>
      </resource>
    </clone>
    <clone id="cln_SAPHanaTopology_PRD_HDB00" multi_state="false" unique="false" managed="true" failed="false" failure_ignored="false">
      <resource id="rsc_SAPHanaTopology_PRD_HDB00" resource_agent="ocf::suse:SAPHanaTopology" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
        <node name="hana01" id="1084787210" cached="true"/>
      </resource>
      <resource id="rsc_SAPHanaTopology_PRD_HDB00" resource_agent="ocf::suse:SAPHanaTopology" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
        <node name="hana02" id="1084787211" cached="true"/>
      </resource>
    </clone>
  </resources>
  <node_attributes>
    <node name="hana01">
      <attribute name="hana_prd_clone_state" value="PROMOTED"/>
      <attribute name="hana_prd_op_mode" value="logreplay"/>
      <attribute name="hana_prd_remoteHost" value="hana02"/>
      <attribute name="hana_prd_roles" value="4:P:master1:master:worker:master"/>
      <attribute name="hana_prd_site" value="Site1"/>
      <attribute name="hana_prd_srmode" value="sync"/>
      <attribute name="hana_prd_sync_state" value="PRIM"/>
      <attribute name="hana_prd_version" value="2.00.040.00.1553674765"/>
      <attribute name="hana_prd_vhost" value="hana01"/>
      <attribute name="lpa_prd_lpt" value="1618862836"/>
      <attribute name="master-rsc_SAPHana_PRD_HDB00" value="150"/>
    </node>
    <node name="hana02">
      <attribute name="hana_prd_clone_state" value="DEMOTED"/>
      <attribute name="hana_prd_op_mode" value="logreplay"/>
      <attribute name="hana_prd_remoteHost" value="hana01"/>
      <attribute name="hana_prd_roles" value="4:S:master1:master:worker:master"/>
      <attribute name="hana_prd_site" value="Site2"/>
      <attribute name="hana_prd_srmode" value="sync"/>
      <attribute name="hana_prd_sync_state" value="SOK"/>
      <attribute name="hana_prd_version" value="2.00.040.00.1553674765"/>
      <attribute name="hana_prd_vhost" value="hana02"/>
      <attribute name="lpa_prd_lpt" value="30"/>
      <attribute name="master-rsc_SAPHana_PRD_HDB00" value="100"/>
    </node>
  </node_attributes>
  <node_history>
    <node name="hana01">
      <resource_history id="stonith-sbd" orphan="false" migration-threshold="5000">
        <operation_history call="6" task="start" last-rc-change="Mon Apr 19 14:11:29 2021" last-run="Mon Apr 19 14:11:29 2021" exec-time="1270ms" queue-time="0ms" rc="0" rc_text="ok"/>
      </resource_history>
      <resource_history id="rsc_exporter_PRD_HDB00" orphan="false" migration-threshold="5000">
        <operation_history call="38" task="start" last-rc-change="Mon Apr 19 14:11:50 2021" last-run="Mon Apr 19 14:11:50 2021" exec-time="2267ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="41" task="monitor" interval="10000ms" last-rc-change="Mon Apr 19 14:11:50 2021" exec-time="2ms" queue-time="0ms" rc="0" rc_text="ok"/>
      </resource_history>
      <resource_history id="rsc_SAPHanaTopology_PRD_HDB00" orphan="false" migration-threshold="5000">
        <operation_history call="29" task="start" last-rc-change="Mon Apr 19 14:11:34 2021" last-run="Mon Apr 19 14:11:34 2021" exec-time="3852ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="32" task="monitor" interval="10000ms" last-rc-change="Mon Apr 19 14:11:38 2021" exec-time="4107ms" queue-time="0ms" rc="0" rc_text="ok"/>
      </resource_history>
      <resource_history id="rsc_ip_PRD_HDB00" orphan="false" migration-threshold="5000">
        <operation_history call="25" task="start" last-rc-change="Mon Apr 19 14:11:30 2021" last-run="Mon Apr 19 14:11:30 2021" exec-time="90ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="26" task="monitor" interval="10000ms" last-rc-change="Mon Apr 19 14:11:30 2021" exec-time="57ms" queue-time="0ms" rc="0" rc_text="ok"/>
      </resource_history>
      <resource_history id="rsc_SAPHana_PRD_HDB00" orphan="false" migration-threshold="5000">
        <operation_history call="19" task="probe" last-rc-change="Mon Apr 19 14:11:30 2021" last-run="Mon Apr 19 14:11:30 2021" exec-time="2801ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="40" task="promote" last-rc-change="Mon Apr 19 14:11:48 2021" last-run="Mon Apr 19 14:11:48 2021" exec-time="2347ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="42" task="monitor" interval="60000ms" last-rc-change="Mon Apr 19 14:11:54 2021" exec-time="4086ms" queue-time="0ms" rc="8" rc_text="master"/>
      </resource_history>
    </node>
    <node name="hana02">
      <resource_history id="rsc_SAPHana_PRD_HDB00" orphan="false" migration-threshold="5000">
        <operation_history call="18" task="probe" last-rc-change="Mon Apr 19 14:16:56 2021" last-run="Mon Apr 19 14:16:56 2021" exec-time="3035ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="18" task="probe" last-rc-change="Mon Apr 19 14:16:56 2021" last-run="Mon Apr 19 14:16:56 2021" exec-time="3035ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="24" task="monitor" interval="61000ms" last-rc-change="Mon Apr 19 14:16:59 2021" exec-time="3645ms" queue-time="0ms" rc="0" rc_text="ok"/>
      </resource_history>
      <resource_history id="rsc_SAPHanaTopology_PRD_HDB00" orphan="false" migration-threshold="5000">
        <operation_history call="25" task="start" last-rc-change="Mon Apr 19 14:16:59 2021" last-run="Mon Apr 19 14:16:59 2021" exec-time="3555ms" queue-time="0ms" rc="0" rc_text="ok"/>
        <operation_history call="26" task="monitor" interval="10000ms" last-rc-change="Mon Apr 19 14:17:03 2021" exec-time="3714ms" queue-time="0ms" rc="0" rc_text="ok"/>
      </resource_history>
    </node>
  </node_history>
</crm_mon>
EOF
