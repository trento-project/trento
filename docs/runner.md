# Ansible runner

The Trento runner is responsible of running the Trento checks among the installed Trento Agents.

## Implementing a check

To be defined...

## Metadata files

The metadata files provide information about the check themselves. They are used to get information
from the Trento Web GUI and render properly everything related with the Ansible tasks.

In order to use them properly, some fields are required. An example is available at [defaults/main.yml](roles/checks/1.1.1/defaults/main.yml).

These are the fields needed by Trento:

- `id`: The check unique identifier. This value must not be changed during the lifetime of the check
- `name`: A short name for the check. It is used as more user friendly identifier for the check for the developers (variables files with the expected values use this field for example).
  This `name` field provides the visualization order of the checks in the Trento web interface as well
- `group`: The group which the check belongs to. It is used to group the checks under different visual elements in Trento
- `labels`: A list of labels (separated by command) which helps to group the checks by execution groups. The difference between this and the `group` field
is that the labels are used for control purpose (select all the checks with this label e.g.), and the groups are used just for visual purposes
- `description`: A longer description about the check purpose. It can be written using markdown
- `implementation`: Usually the task `main.yml` content

### Creating a new id

The `id` must be unique in the check collection. It must be 6 hexadecimal digits string.
In order to create a new unique identifier, and to check if there are any duplicated entries, the
`id_checker.py` script can be used.

To use it:
```
# Check if the checks include all the required metadata values and if there is any duplicated id
python3 hack/id_checker.py
```

To add a new unique `id`:
```
python3 hack/id_checker.py --generate
```
