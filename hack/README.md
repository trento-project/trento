# Hacks

This directory is a collection of scripts & utilities that can assist developers and users interested in testing and
experimenting with Trento.

## deploy.sh

`deploy.sh` is a very simple script that will attempt to copy the `trento` binary and the `consul` binary to a remote
server and start both services

### Requirements

The machines that we are deploying to require to have `rsync` as well as a running SSH server.

### Usage

`./deploy.sh [username@]<target-server-ip> <consul-ip> [deploy-agent*|deploy-web]`

- `[username]@<target-server-ip>`
  The IP address of the host where we are deploying `trento` and `consul` on.

- `<consul-ip>`
  The IP of the consul server that we are connecting to. When `deploy-web` is used in the next field, this is ignored.

- `[deploy-agent|deploy-web]`
  `deploy-agent` causes to deploy the `consul` and `trento` agents while
  `deploy-web` causes to deploy the web server as well as a `consul` server instance

## Automatic deploy to private nodes with GitHub actions

The repository contains a GitHub actions [workflow](../.github/workflows/ci.yaml) to deploy Trento to private nodes by
using [self-hosted runners](https://docs.github.com/en/actions/hosting-your-own-runners/about-self-hosted-runners).

This is triggered automatically when pushing (or merging PRs) into the `main` branch of the upstream remote, or manually
by
the [workflow_dispatch](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#workflow_dispatch)
event
(e.g., when working on forks).

### Set up runners instances

#### Generate token

Refer
to [this guide](https://docs.github.com/en/actions/hosting-your-own-runners/adding-self-hosted-runners#adding-a-self-hosted-runner-to-a-repository)
to add a runner to the repository. The generated token will be passed as a parameter to the provisioning script,
see [Provision](#provision).

Programmatically registering self-hosted runners is also possible by using the rest API.
See [https://docs.github.com/en/rest/reference/actions#self-hosted-runners](https://docs.github.com/en/rest/reference/actions#self-hosted-runners)

#### Runner labels

The provided actions workflow expects 3 nodes in total: 2 hana nodes and a monitoring node. For the actions to run, the
runners must have a label specifing the node type (`vmhana01`,`vmhana02` and `vmmonitoring`)

#### Provision

[`provision.sh`](./provision.sh) provides an easy way to boostrap a GitHub self-hosted runner environment on a private
node.

The script sets up users, configuration files, systemd unit and a self-hosted runner instance required in order to run
Trento, Consul and the deployment process.

`# ./provision.sh <github-user> <github-repo> <actions-runner-token> <actions-runner-name> <consul-bind-ip> <consul-server-ip> [agent|web]`

- `<github-user>`
  The GitHub username (e.g, `trento_project`).

- `<github-repo>`
  The GitHub repository name (e.g, `trento`).

- `<actions-runner-token>`
  See [Generate token](#generate-token).

- `<actions-runner-name>`
  The actions runner name.

  One of `vmahana01`, `vmhana02`, `vmmonitoring` is expected if used in conjuction with the
  provided [workflow](../.github/workflows/ci.yaml).

- `[agent|web]`
  `agent` provisions an agent instance
  `web` provisions a web instance

#### Example

This example assumes a running cluster composed by two nodes and a monitoring instance. Refer
to [ha-sap-terraform-deployments](https://github.com/SUSE/ha-sap-terraform-deployments) for the setup.

[Generate 3 tokens](#generate-token), then run the provision script on each node.

`vmhana01:~ # sh provision.sh youruser trento <TOKEN_1> vmhana01 10.162.30.92 10.162.29.225 agent`

`vmhana02:~ # sh provision.sh youruser trento <TOKEN_2> vmhana02 10.162.30.93 10.162.29.225 agent`

`vmmonitoring:~ # sh provision.sh youruser trento <TOKEN_3> vmmonitoring 10.162.29.225 10.162.29.225 web`

### Manually triggering the deployment

It is possible to manually trigger the deployment job from the GitHub Actions UI.

See [https://github.blog/changelog/2020-07-06-github-actions-manual-triggers-with-workflow_dispatch/](https://github.blog/changelog/2020-07-06-github-actions-manual-triggers-with-workflow_dispatch/)
for reference.