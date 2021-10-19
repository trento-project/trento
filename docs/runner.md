# Ansible runner

This health checks are written in ansible following a pre-defined files structure. Find how to write a new check in the next chapters

## Implementing a check

# Table of contents

- [Structure]
   - [Check structure](#check-structure)
   - [Metadata files](#metadata-files)
   - [Check files](#check-files)
   - [Creating a new ID](#add-id)
   - [Examples](#examples)

## Check structure

The checks folder is in `runner/ansible/roles/checks`. Each check is stored in an individual folder, which gives the name to the check (e.g. 1.1.1). The name refers to its place in the queue (1.1.2 will be 
executed after 1.1.1 etc.).These check folders contain other two subfolders. Both of these each has one file in them named `main.yml`:
- `defaults`
   The `main.yml`file in the `defaults` directory contains all the required [metadata](#metadata-files) for the check. 

- `tasks`
  The `main.yml` file in the `tasks` directory contains the [check](#check-files). 


## Metadata files

The metadata files provide information about the check's themselves. They are used to get information
from the Trento Web GUI and render properly everything related with the Ansible tasks.

In order to use them properly, some fields are required. An example is available at [defaults/main.yml](roles/checks/1.1.1/defaults/main.yml).

These are the fields needed by Trento:

- `id`: The check's unique identifier.
- `name`: A short name for the check. It is used as more user friendly identifier and defines the execution order of the checks.
- `group`: The group which the check belongs to. It is used to group the checks under different visual elements in Trento
- `labels`: A list of labels (separated by command) which helps to group the checks by execution groups. The difference between this and the `group` field
is that the labels are used for control purpose (select all the checks with this label e.g.), and the groups are used just for visual purposes
- `description`: A longer description about the check's purpose. It can be written using markdown.
- `implementation`: Usually the task `main.yml` content
- `on_failure` : This field is a boolean which decides if the test result has a warning state on failure rather than the critical state.

## Check files

The check files contain the actual task which the runner executes. An
example is available at [tasks/main.yml](roles/checks/1.1.1/tasks/main.yml).
These files are written like a normal ansible task, but:

The tasks are of `read only` nature, meaning that they are meant to check for things rather than executing something.

If the expected value is a variable (something that differs in the different cloud providers for example), this value can be added in the files available at the `runner/ansible/vars` folder. The entries in this file follow a key/value syntax, where the key is the check name and the value the expected value

```
name: "{{ id }}.check"
```

```  
check_mode: no
register: config_updated
changed_when: config_updated.stdout != expected[name]
```
```
- block:
    - import_role:
        name: post-results
  when:
    - ansible_check_mode
  vars:
    status: "{{ config_updated is not changed }}"
``` 

## Creating a new ID

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

## Examples

As an example (in this case `dummy`), lets add a check which checks if the /etc/os-release file exists.
First, define the metadata in `runner/ansible/roles/checks/dummy/defaults/main.yml`:

```
---

name: dummy
name: Check os-release file
group: Testing
labels: generic
description: |
  Checks for the os-release file in /etc/os-release
remediation: |
  ## Remediation
  This check is for exemplary purposes

  ## References
  Place for references
implementation: "{{ lookup('file', 'roles/checks/'+id+'/tasks/main.yml') }}"

# Test data
key_name: token
id: FFFFFF
``` 
(In this example the ID FFFFFF is used [(six hex digits)](#add-id). When creating a check the `hack/id_checker.py` script should be used.
It checks if the ID's are unique and adds ID's to the checks if they are missing.)

Then the actual task `runner/ansible/roles/checks/dummy/tasks/main.yml`:

```
---

- name: "{{ name }}.check"
  stat:
      path: "/etc/os-release"
  check_mode: no
  register: config_updated
  changed_when: config_updated.stdout != expected[name]

- block:
    - import_role:
        name: post-results
  when:
    - ansible_check_mode
  vars:
    status: "{{ config_updated is not changed }}"
```