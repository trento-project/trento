
# Hacks

This directory is a collection of scripts & utilities that can assist developers
and users interested in testing and experimenting with Trento.

## deploy.sh
`deploy.sh` is a very simple script that will attempt to copy the `trento` binary
and the `consul` binary to a remote server and start both services

### Requirements
The machines that we are deploying to require to have `rsync` as well as a running
SSH server.

### Usage
`./deploy.sh [username@]<target-server-ip> <consul-ip> [deploy-agent*|deploy-web]`

  - `[username]@<target-server-ip>`
    The IP address of the host where we are deploying `trento` and `consul` on.
    
  - `<consul-ip>`
    The IP of the consul server that we are connecting to. When `deploy-web` is
    used in the next field, this is ignored.

  - `[deploy-agent|deploy-web]`
    `deploy-agent` causes to deploy the `consul` and `trento` agents 
    while
    `deploy-web` causes to deploy the web server as well as a `consul` server instance
