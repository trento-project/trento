"""
Trento Callback executor.
Find some documentation at:
https://ara.readthedocs.io/en/latest/api-usage.html
https://ara.readthedocs.io/en/latest/ansible-plugins-and-use-cases.html
https://github.com/ansible-community/ara/blob/master/ara/plugins/action/ara_record.py

:author: xarbulu
:organization: SUSE Linux GmbH
:contact: xarbulu@suse.com

:since: 2021-09-16
"""

import os
import yaml

from ansible.plugins.callback import CallbackBase
from ara.clients.http import AraHttpClient

TRENTO_TEST_LABEL_KEY = "ara_playbook_labels"
TRENTO_TEST_LABEL = "test"
TRENTO_RECORD_KEY = "trento-results"
TEST_RESULT_TASK_NAME = "set_test_result"
TEST_INCLUDE_TASK_NAME = "run_checks"
EXTERNAL_ID = "external_id"


class Results(object):
    """
    Object to store and user the execution results

    Result example:

    "results": {
        "clusterId": {
            "checks": {
                "ABCDEF": {
                    "hosts": {
                        "host1": {
                            "result": "passing"
                        }
                    }
                }
            }
        }
    }
    """
    def __init__(self):
        self.results = {"results": {}}

    def initialize_group(self, group):
        """
        Initialize the group on the results dictionary
        """
        if group not in self.results["results"]:
            self.results["results"][group] = {}
            self.results["results"][group]["checks"] = {}

    def add_result(self, group, test, host, result):
        """
        Add new result
        """
        # Add the group just in case it doesn't exist
        if group not in self.results["results"]:
            self.results["results"][group] = {}
            self.results["results"][group]["checks"] = {}

        checks = self.results["results"][group]["checks"]
        if test not in checks:
            checks[test] = {}
            checks[test]["hosts"] = {}

        hosts = checks[test]["hosts"]
        if host not in hosts:
            hosts[host] = {}

        hosts[host]["result"] = result


class CallbackModule(CallbackBase):
    """
    Trento Callback module
    """
    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'aggregate'
    CALLBACK_NAME = 'trento'

    def __init__(self):
        super(CallbackModule, self).__init__()
        self.playbook = None
        self.play = None
        self.results = Results()
        endpoint = os.getenv('ARA_API_SERVER')
        self.client = AraHttpClient(endpoint=endpoint, verify=False)

    def v2_playbook_on_start(self, playbook):
        """
        On start callback
        """
        self._display.banner("Trento callback loaded")
        self.playbook = playbook

    def v2_playbook_on_play_start(self, play):
        """
        On Play start callback
        """
        self.play = play
        self._initialize_results()

    def v2_runner_on_ok(self, result):
        """
        On task Ok
        """
        if self._is_check_include_loop(result):
            self._store_skipped(result)
            return

        if not self._is_test_result(result):
            return

        host = result._host.get_name()
        task_vars = self._all_vars(host=result._host, task=result._task)

        test_result = result._task_fields["args"]["test_result"]
        for group in task_vars["group_names"]:
            self.results.add_result(group, task_vars[EXTERNAL_ID], host, test_result)

    def v2_runner_on_failed(self, result):
        """
        On task Failed
        """
        if not self._is_test_result(result):
            return

        host = result._host.get_name()
        task_vars = self._all_vars(host=result._host, task=result._task)

        for group in task_vars["group_names"]:
            self.results.add_result(group, task_vars[EXTERNAL_ID], host, False)

    def v2_playbook_on_stats(self, _stats):
        """
        Upload ARA record at the end of the execution
        """
        if not self._is_test_execution():
            return

        self._display.banner("Publishing Trento results")
        self._create_or_update_record(
            self._get_playbook_id(),
            TRENTO_RECORD_KEY,
            self.results.results,
            "json")

    def _all_vars(self, host=None, task=None):
        """
        Get task vars

        host and task need to be specified in case 'magic variables' (host vars, group vars, etc)
        need to be loaded as well
        """
        return self.play.get_variable_manager().get_vars(
            play=self.play,
            host=host,
            task=task
        )

    def _initialize_results(self):
        """
        Initialize the results object
        """
        play_vars = self._all_vars()
        for _, host_data in play_vars["hostvars"].items():
            for group in host_data["group_names"]:
                self.results.initialize_group(group)

    def _is_test_execution(self):
        """
        Check if the current execution is a trento test execution
        """
        play_vars = self._all_vars()
        if TRENTO_TEST_LABEL_KEY not in play_vars or \
                 TRENTO_TEST_LABEL not in play_vars[TRENTO_TEST_LABEL_KEY]:
            self._display.banner("Not running a Trento test execution")
            return False
        return True

    def _is_test_result(self, result):
        """
        Check if the current task is a test result
        """
        if (result._task_fields.get("action") == "set_fact") and \
                (result._task_fields.get("name") == TEST_RESULT_TASK_NAME):
            return True
        return False

    def _is_check_include_loop(self, result):
        """
        Check if the current task is the checks include loop task
        """
        if (result._task_fields.get("action") == "include_role") and \
                (result._task_fields.get("name") == TEST_INCLUDE_TASK_NAME):
            return True
        return False

    def _store_skipped(self, result):
        """
        Store skipped checks
        """
        task_vars = self._all_vars(host=result._host, task=result._task)
        host = result._host.get_name()

        for check_result in result._result["results"]:
            skipped = check_result.get("skipped", False)
            if skipped:
                with open(os.path.join(
                    check_result["check_item"]["path"], "defaults/main.yml")) as file_ptr:

                    data = yaml.load(file_ptr, Loader=yaml.Loader)
                    check_id = data[EXTERNAL_ID]

                for group in task_vars["group_names"]:
                    self.results.add_result(group, check_id, host, "skipped")

    def _get_playbook_id(self):
        """
        Get current execution playbook id
        """
        play = self.client.get("/api/v1/plays?uuid=%s" % self.play._uuid)
        playbook_id = play["results"][0]["playbook"]
        return playbook_id

    def _create_or_update_record(self, playbook, key, value, record_type):
        """
        Create or update (if it already exists) an ARA record

        Based on:
        https://github.com/ansible-community/ara/blob/master/ara/plugins/action/ara_record.py#L145
        """
        changed = False
        record = self.client.get("/api/v1/records?playbook=%s&key=%s" % (playbook, key))
        if record["count"] == 0:
            record = self.client.post(
                "/api/v1/records", playbook=playbook, key=key, value=value, type=record_type)
            changed = True
        else:
            old = self.client.get("/api/v1/records/%s" % record["results"][0]["id"])
            if old["value"] != value or old["type"] != type:
                record = self.client.patch(
                    "/api/v1/records/%s" % old["id"], key=key, value=value, type=record_type)
                changed = True
            else:
                record = old
        return record, changed
