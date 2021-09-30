# Ansible runner

Place to store Ansible playbooks that are executed by the Trento runner

## Implementing a check

# Table of contents

- [Structure]
   - [Metadata files](#metadata-files)
   - [Check files](#check-files)
   - [Workflow](#workflow)

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
These files are written like a normal ansible task.

## Workflow
1. Adding checks
- Create a directory in `/runner/ansible/roles/checks` and name it
accordingly (e.g. `1.6.1`). In this newly created directory, add two more
called `defaults` and `tasks`. 

- In the `defaults` directory, create a file called `main.yml` and fill it
with the required metadata as described in [Metadata files](#metadata-files).

- Also add a `main.yml` file into the `tasks` directory and fill it with 
the task you wish to be executed according to the ansible syntax. 

2. Stopping trento processes
- Stop the running trento processes. 
E.g.: Execute `"ps aux | grep trento"` and kill the running trento services. 
- Stop the running kubernetes cluster by executing `"k3s-killall.sh"`.

3. Building trento and moving it to the machine
- Navigate to the root directory of trento and execute `make build`.
- Copy the newly build binary to the target machine, with a command
 like `scp`, to the according path (`/usr/bin/trento`).

4. Rerunning trento
- Start the trento services:
- Restart the k3s: `"systemctl start k3s"`. 
- `"./trento runner start --ara-server <ARA-SERVER-IP-AND-PORT> --consul-addr <TRENTO-CONSUL-SERVER-AND-PORT> -i 5"` (Runner).
- `"./trento web serve --ara-addr <ARA-SERVER-IP-AND-PORT> -p <ARA-SERVER-PORT>"` (Trento Web)

