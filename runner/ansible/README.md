# Ansible runner

Place to store Ansible playbooks that are executed by the Trento runner

## Implementing a check

To be defined...

## Metadata files

The metadata files provide information about the test themselves. They are used to get information
from the Trento Web GUI and render properly everything related with the Ansible tasks.

In order to use them properly, some fields are required. An example is available at [defaults/main.yml](roles/checks/1.1.1/defaults/main.yml).

These are the fields needed by Trento:

- `id`: The test unique identifier.
- `name`: A short name for the check. It is used as more user friendly identifier for the test
- `group`: The group which the check belongs to. It is used to group the checks under different visual elements in Trento
- `labels`: A list of labels (separated by command) which helps to group the checks by execution groups. The difference between this and the `group` field
is that the labels are used for control purpose (select all the checks with this label e.g.), and the groups are used just for visual purposes
- `description`: A longer description about the test purpose. It can be written using markdown.
- `implementation`: Usually the task `main.yml` content
