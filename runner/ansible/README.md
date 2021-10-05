# Ansible runner

Place to store Ansible playbooks that are executed by the Trento runner

## Implementing a check

# Table of contents

- [Structure]
   - [Check structure](#check-structure)
   - [Metadata files](#metadata-files)
   - [Check files](#check-files)
   - [Examples](#examples)

## Check structure

The checks folder is in `runner/ansible/roles/checks`. Each check has its own
number (e.g. 1.1.1) which refers to its place in the queue (1.1.2 will be 
executed after 1.1.1 etc.). Inside these numbered folders are two subfolders, both of these each has one file in them named `main.yml`:
- `defaults`
   The defaults directory (or better the main.yml file inside of it)
   contains all the required [metadata](#metadata-files) for the check. 

- `tasks`
  Contains the [check](#check-files). 


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

## Check files

The check files contain the actual task which the runner will execute. An
example is available at [tasks/main.yml](roles/checks/1.1.1/tasks/main.yml).
These files are written like a normal ansible task, but:

- The tasks are of `read only` nature, meaning that they are meant to check for things
  rather than something beeing executed.
- Their output is either `true` (check passed) or `false` (check failed).
  This output will then be passed to the `post-results` block with the `status` 
  field containing the result of the check. 

```
name: "{{ id }}.check"
```

```  
check_mode: no
register: config_updated
changed_when: config_updated.stdout != expected[id]
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

## Examples

As an example (we will use `dummy`), lets add a test which simply creates a simple file in /tmp.
Fist we want to define the metadata in `runner/ansible/roles/checks/dummy/defaults/main.yml`:

```
---

id: dummy
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
``` 

Then the actual task `runner/ansible/roles/checks/3.0.0/tast/main.yml`:

```
---

- name: "{{ id }}.check"
  stat:
      path: "/etc/os-release"
  check_mode: no
  register: config_updated
  changed_when: config_updated.stdout != expected[id]

- block:
    - import_role:
        name: post-results
  when:
    - ansible_check_mode
  vars:
    status: "{{ config_updated is not changed }}"
``` 



