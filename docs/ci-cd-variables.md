## CI/CD Variables expected by Trento GH Action Runner

These are a list of the variables that our GitHub Workflow supports in order
to make a `trento` deployment.

### Secrets

#### `OBS_USER`
This the OBS user that will trigger build in `opensuse.build.org`
`obsuser`

#### `OBS_PASS`
This is the password for the OBS user that will trigger build in `opensuse.build.org`
`obspassword`


## Environments

### `AZURE_DEMO` 
The default name for the environment is `AZURE_DEMO`. 
All of the above variables belong to this environment.

#### `SSH_CONFIG`
Configuration file for the SSH client used to reach the agents
```
Host hostname01
    User user
    IdentityFile /path/to/ssh/id_rsa
    StrictHostKeyChecking no
    ...
```

#### `SSH_KEY` 
Pre-authorized SSH key in the target hosts in order to access all hosts
listed in `TRENTO_TARGET_AGENTS` to trigger the installation on them

```
-----BEGIN OPENSSH PRIVATE KEY-----
the key
-----END OPENSSH PRIVATE KEY-----
```

#### `TRENTO_USER`
The user that will run the CI code to deploy trento in the target machines
```
cloudadmin
```

#### `TRENTO_SERVER_HOST` 
The IP of the machine where we should run install-server.sh on
```
10.x.x.x
```


#### `TRENTO_AGENT_HOSTS`
A comma-separated list of hosts where to run install-agent.sh on (cluster nodes)
```
10.x.x.x,10.x.x.x,10.x.x.x
```