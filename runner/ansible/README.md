# Ansible runner

Place to store Ansible playbooks that are executed by the Trento runner

## Implementing a check

To be defined...

## Metadata files

The metadata files provide information about the check themselves. They are used to get information
from the Trento Web GUI and render properly everything related with the Ansible tasks.

In order to use them properly, some fields are required. An example is available at [defaults/main.yml](roles/checks/1.1.1/defaults/main.yml).

These are the fields needed by Trento:

- `external_id`: The check unique identifier. This value must not be changed during the lifetime of the check
- `id`: The internal identifier. It is used to identify the checks internally (a more user-friendly nomenclature), to specify the expected values and to apply the ordering to the lists in the web side. This value can be changed to adapt to any new internal need.
- `name`: A short name for the check. It is used as more user friendly identifier for the check
- `group`: The group which the check belongs to. It is used to group the checks under different visual elements in Trento
- `labels`: A list of labels (separated by command) which helps to group the checks by execution groups. The difference between this and the `group` field
is that the labels are used for control purpose (select all the checks with this label e.g.), and the groups are used just for visual purposes
- `description`: A longer description about the check purpose. It can be written using markdown
- `implementation`: Usually the task `main.yml` content

## Creating a new external_id

The `external_id` must be unique in the check collection. It must be 6 hexadecimal digits string.
In order to create a new unique identifier, and to check if there are any duplicated entries, the
`id_checker.py` script can be used.

To use it:
```
# Check if the checks include all the required metadata values and if there is any duplicated external_id
python3 support/id_checker.py
```

To add a new unique `external_id`:
```
python3 support/id_checker.py --generate
```
