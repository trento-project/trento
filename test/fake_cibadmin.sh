#!/bin/sh

cat <<EOF
<cib crm_feature_set="3.3.0" validate-with="pacemaker-3.3" epoch="474" num_updates="0" admin_epoch="0" cib-last-written="Mon Apr 19 22:04:07 2021" update-origin="hana01" update-client="crm_attribute" update-user="root" have-quorum="1" dc-uuid="1084787210">
  <configuration>
    <crm_config>
      <cluster_property_set id="cib-bootstrap-options">
        <nvpair id="cib-bootstrap-options-have-watchdog" name="have-watchdog" value="true"/>
        <nvpair id="cib-bootstrap-options-dc-version" name="dc-version" value="2.0.4+20200616.2deceaa3a-3.3.1-2.0.4+20200616.2deceaa3a"/>
        <nvpair id="cib-bootstrap-options-cluster-infrastructure" name="cluster-infrastructure" value="corosync"/>
        <nvpair id="cib-bootstrap-options-cluster-name" name="cluster-name" value="hana_cluster"/>
        <nvpair name="stonith-enabled" value="true" id="cib-bootstrap-options-stonith-enabled"/>
      </cluster_property_set>
      <cluster_property_set id="SAPHanaSR">
        <nvpair id="SAPHanaSR-hana_prd_site_srHook_Site2" name="hana_prd_site_srHook_Site2" value="SOK"/>
      </cluster_property_set>
    </crm_config>
    <nodes>
      <node id="1084787210" uname="hana01">
        <instance_attributes id="nodes-1084787210">
          <nvpair id="nodes-1084787210-lpa_prd_lpt" name="lpa_prd_lpt" value="1618862646"/>
          <nvpair id="nodes-1084787210-hana_prd_vhost" name="hana_prd_vhost" value="hana01"/>
          <nvpair id="nodes-1084787210-hana_prd_site" name="hana_prd_site" value="Site1"/>
          <nvpair id="nodes-1084787210-hana_prd_op_mode" name="hana_prd_op_mode" value="logreplay"/>
          <nvpair id="nodes-1084787210-hana_prd_srmode" name="hana_prd_srmode" value="sync"/>
          <nvpair id="nodes-1084787210-hana_prd_remoteHost" name="hana_prd_remoteHost" value="hana02"/>
        </instance_attributes>
      </node>
      <node id="1084787211" uname="hana02">
        <instance_attributes id="nodes-1084787211">
          <nvpair id="nodes-1084787211-lpa_prd_lpt" name="lpa_prd_lpt" value="30"/>
          <nvpair id="nodes-1084787211-hana_prd_op_mode" name="hana_prd_op_mode" value="logreplay"/>
          <nvpair id="nodes-1084787211-hana_prd_vhost" name="hana_prd_vhost" value="hana02"/>
          <nvpair id="nodes-1084787211-hana_prd_remoteHost" name="hana_prd_remoteHost" value="hana01"/>
          <nvpair id="nodes-1084787211-hana_prd_site" name="hana_prd_site" value="Site2"/>
          <nvpair id="nodes-1084787211-hana_prd_srmode" name="hana_prd_srmode" value="sync"/>
        </instance_attributes>
      </node>
    </nodes>
    <resources>
      <primitive id="stonith-sbd" class="stonith" type="external/sbd">
        <instance_attributes id="stonith-sbd-instance_attributes">
          <nvpair name="pcmk_delay_max" value="30s" id="stonith-sbd-instance_attributes-pcmk_delay_max"/>
        </instance_attributes>
      </primitive>
      <primitive id="rsc_ip_PRD_HDB00" class="ocf" provider="heartbeat" type="IPaddr2">
        <!--#####################################################-->
        <!--# Fencing agents - Native agents for cloud providers-->
        <!--#####################################################-->
        <!--######################################-->
        <!--# Floating IP address resource agents-->
        <!--######################################-->
        <instance_attributes id="rsc_ip_PRD_HDB00-instance_attributes">
          <nvpair name="ip" value="192.168.138.12" id="rsc_ip_PRD_HDB00-instance_attributes-ip"/>
          <nvpair name="cidr_netmask" value="24" id="rsc_ip_PRD_HDB00-instance_attributes-cidr_netmask"/>
          <nvpair name="nic" value="eth1" id="rsc_ip_PRD_HDB00-instance_attributes-nic"/>
        </instance_attributes>
        <operations>
          <op name="start" timeout="20" interval="0" id="rsc_ip_PRD_HDB00-start-0"/>
          <op name="stop" timeout="20" interval="0" id="rsc_ip_PRD_HDB00-stop-0"/>
          <op name="monitor" interval="10" timeout="20" id="rsc_ip_PRD_HDB00-monitor-10"/>
        </operations>
      </primitive>
      <primitive id="rsc_exporter_PRD_HDB00" class="systemd" type="prometheus-hanadb_exporter@PRD_HDB00">
        <!--#######################################-->
        <!--# non-production HANA - Cost optimized-->
        <!--#######################################-->
        <!--###############################-->
        <!--# Active/Active HANA resources-->
        <!--###############################-->
        <!--######################################-->
        <!--# prometheus-hanadb_exporter resource-->
        <!--######################################-->
        <operations>
          <op name="start" interval="0" timeout="100" id="rsc_exporter_PRD_HDB00-start-0"/>
          <op name="stop" interval="0" timeout="100" id="rsc_exporter_PRD_HDB00-stop-0"/>
          <op name="monitor" interval="10" id="rsc_exporter_PRD_HDB00-monitor-10"/>
        </operations>
        <meta_attributes id="rsc_exporter_PRD_HDB00-meta_attributes">
          <nvpair name="target-role" value="Started" id="rsc_exporter_PRD_HDB00-meta_attributes-target-role"/>
        </meta_attributes>
      </primitive>
      <master id="msl_SAPHana_PRD_HDB00">
        <meta_attributes id="msl_SAPHana_PRD_HDB00-meta_attributes">
          <nvpair name="clone-max" value="2" id="msl_SAPHana_PRD_HDB00-meta_attributes-clone-max"/>
          <nvpair name="clone-node-max" value="1" id="msl_SAPHana_PRD_HDB00-meta_attributes-clone-node-max"/>
          <nvpair name="interleave" value="true" id="msl_SAPHana_PRD_HDB00-meta_attributes-interleave"/>
        </meta_attributes>
        <primitive id="rsc_SAPHana_PRD_HDB00" class="ocf" provider="suse" type="SAPHana">
          <instance_attributes id="rsc_SAPHana_PRD_HDB00-instance_attributes">
            <nvpair name="SID" value="PRD" id="rsc_SAPHana_PRD_HDB00-instance_attributes-SID"/>
            <nvpair name="InstanceNumber" value="00" id="rsc_SAPHana_PRD_HDB00-instance_attributes-InstanceNumber"/>
            <nvpair name="PREFER_SITE_TAKEOVER" value="True" id="rsc_SAPHana_PRD_HDB00-instance_attributes-PREFER_SITE_TAKEOVER"/>
            <nvpair name="AUTOMATED_REGISTER" value="False" id="rsc_SAPHana_PRD_HDB00-instance_attributes-AUTOMATED_REGISTER"/>
            <nvpair name="DUPLICATE_PRIMARY_TIMEOUT" value="7200" id="rsc_SAPHana_PRD_HDB00-instance_attributes-DUPLICATE_PRIMARY_TIMEOUT"/>
          </instance_attributes>
          <operations>
            <op name="start" interval="0" timeout="3600" id="rsc_SAPHana_PRD_HDB00-start-0"/>
            <op name="stop" interval="0" timeout="3600" id="rsc_SAPHana_PRD_HDB00-stop-0"/>
            <op name="promote" interval="0" timeout="3600" id="rsc_SAPHana_PRD_HDB00-promote-0"/>
            <op name="monitor" interval="60" role="Master" timeout="700" id="rsc_SAPHana_PRD_HDB00-monitor-60"/>
            <op name="monitor" interval="61" role="Slave" timeout="700" id="rsc_SAPHana_PRD_HDB00-monitor-61"/>
          </operations>
        </primitive>
      </master>
      <clone id="cln_SAPHanaTopology_PRD_HDB00">
        <meta_attributes id="cln_SAPHanaTopology_PRD_HDB00-meta_attributes">
          <nvpair name="is-managed" value="true" id="cln_SAPHanaTopology_PRD_HDB00-meta_attributes-is-managed"/>
          <nvpair name="clone-node-max" value="1" id="cln_SAPHanaTopology_PRD_HDB00-meta_attributes-clone-node-max"/>
          <nvpair name="interleave" value="true" id="cln_SAPHanaTopology_PRD_HDB00-meta_attributes-interleave"/>
        </meta_attributes>
        <primitive id="rsc_SAPHanaTopology_PRD_HDB00" class="ocf" provider="suse" type="SAPHanaTopology">
          <!--#####################-->
          <!--# SAP HANA resources-->
          <!--#####################-->
          <instance_attributes id="rsc_SAPHanaTopology_PRD_HDB00-instance_attributes">
            <nvpair name="SID" value="PRD" id="rsc_SAPHanaTopology_PRD_HDB00-instance_attributes-SID"/>
            <nvpair name="InstanceNumber" value="00" id="rsc_SAPHanaTopology_PRD_HDB00-instance_attributes-InstanceNumber"/>
          </instance_attributes>
          <operations>
            <op name="monitor" interval="10" timeout="600" id="rsc_SAPHanaTopology_PRD_HDB00-monitor-10"/>
            <op name="start" interval="0" timeout="600" id="rsc_SAPHanaTopology_PRD_HDB00-start-0"/>
            <op name="stop" interval="0" timeout="300" id="rsc_SAPHanaTopology_PRD_HDB00-stop-0"/>
          </operations>
        </primitive>
      </clone>
    </resources>
    <constraints>
      <rsc_colocation id="col_saphana_ip_PRD_HDB00" score="2000" rsc="rsc_ip_PRD_HDB00" rsc-role="Started" with-rsc="msl_SAPHana_PRD_HDB00" with-rsc-role="Master"/>
      <rsc_order id="ord_SAPHana_PRD_HDB00" kind="Optional" first="cln_SAPHanaTopology_PRD_HDB00" then="msl_SAPHana_PRD_HDB00"/>
      <rsc_colocation id="col_exporter_PRD_HDB00" score="+INFINITY" rsc="rsc_exporter_PRD_HDB00" rsc-role="Started" with-rsc="msl_SAPHana_PRD_HDB00" with-rsc-role="Master"/>
    </constraints>
    <rsc_defaults>
      <meta_attributes id="rsc-options">
        <nvpair name="resource-stickiness" value="1000" id="rsc-options-resource-stickiness"/>
        <nvpair name="migration-threshold" value="5000" id="rsc-options-migration-threshold"/>
      </meta_attributes>
    </rsc_defaults>
    <op_defaults>
      <meta_attributes id="op-options">
        <nvpair name="timeout" value="600" id="op-options-timeout"/>
        <nvpair name="record-pending" value="true" id="op-options-record-pending"/>
      </meta_attributes>
    </op_defaults>
  </configuration>
  <status>
    <node_state id="1084787210" uname="hana01" in_ccm="true" crmd="online" crm-debug-origin="do_state_transition" join="member" expected="member">
      <transient_attributes id="1084787210">
        <instance_attributes id="status-1084787210">
          <nvpair id="status-1084787210-master-rsc_SAPHana_PRD_HDB00" name="master-rsc_SAPHana_PRD_HDB00" value="150"/>
          <nvpair id="status-1084787210-hana_prd_version" name="hana_prd_version" value="2.00.040.00.1553674765"/>
          <nvpair id="status-1084787210-hana_prd_clone_state" name="hana_prd_clone_state" value="PROMOTED"/>
          <nvpair id="status-1084787210-hana_prd_sync_state" name="hana_prd_sync_state" value="PRIM"/>
          <nvpair id="status-1084787210-hana_prd_roles" name="hana_prd_roles" value="4:P:master1:master:worker:master"/>
        </instance_attributes>
      </transient_attributes>
      <lrm id="1084787210">
        <lrm_resources>
          <lrm_resource id="stonith-sbd" type="external/sbd" class="stonith">
            <lrm_rsc_op id="stonith-sbd_last_0" operation_key="stonith-sbd_start_0" operation="start" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="2:0:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;2:0:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="6" rc-code="0" op-status="0" interval="0" last-rc-change="1618834289" last-run="1618834289" exec-time="1270" queue-time="0" op-digest="265be3215da5e5037d35e7fe1bcc5ae0"/>
          </lrm_resource>
          <lrm_resource id="rsc_exporter_PRD_HDB00" type="prometheus-hanadb_exporter@PRD_HDB00" class="systemd">
            <lrm_rsc_op id="rsc_exporter_PRD_HDB00_last_0" operation_key="rsc_exporter_PRD_HDB00_start_0" operation="start" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="8:6:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;8:6:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="38" rc-code="0" op-status="0" interval="0" last-rc-change="1618834310" last-run="1618834310" exec-time="2267" queue-time="0" op-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
            <lrm_rsc_op id="rsc_exporter_PRD_HDB00_monitor_10000" operation_key="rsc_exporter_PRD_HDB00_monitor_10000" operation="monitor" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="9:7:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;9:7:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="41" rc-code="0" op-status="0" interval="10000" last-rc-change="1618834310" exec-time="2" queue-time="0" op-digest="0d721f3bcf63b8d121ad4839b260e42a"/>
          </lrm_resource>
          <lrm_resource id="rsc_SAPHanaTopology_PRD_HDB00" type="SAPHanaTopology" class="ocf" provider="suse">
            <lrm_rsc_op id="rsc_SAPHanaTopology_PRD_HDB00_last_0" operation_key="rsc_SAPHanaTopology_PRD_HDB00_start_0" operation="start" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="20:2:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;20:2:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="29" rc-code="0" op-status="0" interval="0" last-rc-change="1618834294" last-run="1618834294" exec-time="3852" queue-time="0" op-digest="2d8d79c3726afb91c33d406d5af79b53" op-force-restart="" op-restart-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
            <lrm_rsc_op id="rsc_SAPHanaTopology_PRD_HDB00_monitor_10000" operation_key="rsc_SAPHanaTopology_PRD_HDB00_monitor_10000" operation="monitor" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="24:3:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;24:3:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="32" rc-code="0" op-status="0" interval="10000" last-rc-change="1618834298" exec-time="4107" queue-time="0" op-digest="64db68ca3e12e0d41eb98ce63b9610d2"/>
          </lrm_resource>
          <lrm_resource id="rsc_ip_PRD_HDB00" type="IPaddr2" class="ocf" provider="heartbeat">
            <lrm_rsc_op id="rsc_ip_PRD_HDB00_last_0" operation_key="rsc_ip_PRD_HDB00_start_0" operation="start" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="7:1:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;7:1:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="25" rc-code="0" op-status="0" interval="0" last-rc-change="1618834290" last-run="1618834290" exec-time="90" queue-time="0" op-digest="6e3bbd07a422997302424264856a2840"/>
            <lrm_rsc_op id="rsc_ip_PRD_HDB00_monitor_10000" operation_key="rsc_ip_PRD_HDB00_monitor_10000" operation="monitor" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="8:1:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;8:1:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="26" rc-code="0" op-status="0" interval="10000" last-rc-change="1618834290" exec-time="57" queue-time="0" op-digest="8313e7cc541e6aee1c924e232d7f548b"/>
          </lrm_resource>
          <lrm_resource id="rsc_SAPHana_PRD_HDB00" type="SAPHana" class="ocf" provider="suse">
            <lrm_rsc_op id="rsc_SAPHana_PRD_HDB00_last_failure_0" operation_key="rsc_SAPHana_PRD_HDB00_monitor_0" operation="monitor" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="3:1:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;3:1:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="19" rc-code="0" op-status="0" interval="0" last-rc-change="1618834290" last-run="1618834290" exec-time="2801" queue-time="0" op-digest="ff4ff123bc6f906497ef0ef5e44dffd1"/>
            <lrm_rsc_op id="rsc_SAPHana_PRD_HDB00_last_0" operation_key="rsc_SAPHana_PRD_HDB00_promote_0" operation="promote" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="12:6:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;12:6:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="40" rc-code="0" op-status="0" interval="0" last-rc-change="1618834308" last-run="1618834308" exec-time="2347" queue-time="0" op-digest="ff4ff123bc6f906497ef0ef5e44dffd1" op-force-restart=" INSTANCE_PROFILE " op-restart-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
            <lrm_rsc_op id="rsc_SAPHana_PRD_HDB00_monitor_60000" operation_key="rsc_SAPHana_PRD_HDB00_monitor_60000" operation="monitor" crm-debug-origin="build_active_RAs" crm_feature_set="3.3.0" transition-key="14:7:8:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:8;14:7:8:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana01" call-id="42" rc-code="8" op-status="0" interval="60000" last-rc-change="1618834314" exec-time="4086" queue-time="0" op-digest="05b857e482ebd46019d347fd55ebbcdb"/>
          </lrm_resource>
        </lrm_resources>
      </lrm>
    </node_state>
    <node_state id="1084787211" in_ccm="true" crmd="online" crm-debug-origin="do_update_resource" uname="hana02" join="member" expected="member">
      <lrm id="1084787211">
        <lrm_resources>
          <lrm_resource id="stonith-sbd" type="external/sbd" class="stonith">
            <lrm_rsc_op id="stonith-sbd_last_0" operation_key="stonith-sbd_monitor_0" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="5:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:7;5:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="5" rc-code="7" op-status="0" interval="0" last-rc-change="1618834616" last-run="1618834616" exec-time="31" queue-time="0" op-digest="265be3215da5e5037d35e7fe1bcc5ae0"/>
          </lrm_resource>
          <lrm_resource id="rsc_ip_PRD_HDB00" type="IPaddr2" class="ocf" provider="heartbeat">
            <lrm_rsc_op id="rsc_ip_PRD_HDB00_last_0" operation_key="rsc_ip_PRD_HDB00_monitor_0" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="6:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:7;6:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="9" rc-code="7" op-status="0" interval="0" last-rc-change="1618834616" last-run="1618834616" exec-time="44" queue-time="0" op-digest="6e3bbd07a422997302424264856a2840"/>
          </lrm_resource>
          <lrm_resource id="rsc_exporter_PRD_HDB00" type="prometheus-hanadb_exporter@PRD_HDB00" class="systemd">
            <lrm_rsc_op id="rsc_exporter_PRD_HDB00_last_0" operation_key="rsc_exporter_PRD_HDB00_monitor_0" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="7:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:7;7:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="13" rc-code="7" op-status="0" interval="0" last-rc-change="1618834616" last-run="1618834616" exec-time="8" queue-time="0" op-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
          </lrm_resource>
          <lrm_resource id="rsc_SAPHana_PRD_HDB00" type="SAPHana" class="ocf" provider="suse">
            <lrm_rsc_op id="rsc_SAPHana_PRD_HDB00_last_0" operation_key="rsc_SAPHana_PRD_HDB00_monitor_0" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="8:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;8:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="18" rc-code="0" op-status="0" interval="0" last-rc-change="1618834616" last-run="1618834616" exec-time="3035" queue-time="0" op-digest="ff4ff123bc6f906497ef0ef5e44dffd1" op-force-restart=" INSTANCE_PROFILE " op-restart-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
            <lrm_rsc_op id="rsc_SAPHana_PRD_HDB00_last_failure_0" operation_key="rsc_SAPHana_PRD_HDB00_monitor_0" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="8:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;8:15:7:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="18" rc-code="0" op-status="0" interval="0" last-rc-change="1618834616" last-run="1618834616" exec-time="3035" queue-time="0" op-digest="ff4ff123bc6f906497ef0ef5e44dffd1"/>
            <lrm_rsc_op id="rsc_SAPHana_PRD_HDB00_monitor_61000" operation_key="rsc_SAPHana_PRD_HDB00_monitor_61000" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="17:16:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;17:16:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="24" rc-code="0" op-status="0" interval="61000" last-rc-change="1618834619" exec-time="3645" queue-time="0" op-digest="05b857e482ebd46019d347fd55ebbcdb"/>
          </lrm_resource>
          <lrm_resource id="rsc_SAPHanaTopology_PRD_HDB00" type="SAPHanaTopology" class="ocf" provider="suse">
            <lrm_rsc_op id="rsc_SAPHanaTopology_PRD_HDB00_last_0" operation_key="rsc_SAPHanaTopology_PRD_HDB00_start_0" operation="start" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="28:16:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;28:16:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="25" rc-code="0" op-status="0" interval="0" last-rc-change="1618834619" last-run="1618834619" exec-time="3555" queue-time="0" op-digest="2d8d79c3726afb91c33d406d5af79b53" op-force-restart="" op-restart-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
            <lrm_rsc_op id="rsc_SAPHanaTopology_PRD_HDB00_monitor_10000" operation_key="rsc_SAPHanaTopology_PRD_HDB00_monitor_10000" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.3.0" transition-key="29:16:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" transition-magic="0:0;29:16:0:f8dd70dd-608c-49ba-8126-85e2cbebc787" exit-reason="" on_node="hana02" call-id="26" rc-code="0" op-status="0" interval="10000" last-rc-change="1618834623" exec-time="3714" queue-time="0" op-digest="64db68ca3e12e0d41eb98ce63b9610d2"/>
          </lrm_resource>
        </lrm_resources>
      </lrm>
      <transient_attributes id="1084787211">
        <instance_attributes id="status-1084787211">
          <nvpair id="status-1084787211-hana_prd_clone_state" name="hana_prd_clone_state" value="DEMOTED"/>
          <nvpair id="status-1084787211-master-rsc_SAPHana_PRD_HDB00" name="master-rsc_SAPHana_PRD_HDB00" value="100"/>
          <nvpair id="status-1084787211-hana_prd_version" name="hana_prd_version" value="2.00.040.00.1553674765"/>
          <nvpair id="status-1084787211-hana_prd_roles" name="hana_prd_roles" value="4:S:master1:master:worker:master"/>
          <nvpair id="status-1084787211-hana_prd_sync_state" name="hana_prd_sync_state" value="SOK"/>
        </instance_attributes>
      </transient_attributes>
    </node_state>
  </status>
</cib>
EOF
