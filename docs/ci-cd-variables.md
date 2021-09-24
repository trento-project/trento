== CI/CD Variables expected by Trento GH Action Runner


- `SSH_CONFIG`
Configuration file for the SSH client used to reach the agents
```
Host *
    StrictHostKeyChecking no
```

- `SSH_KEY` 
Pre-authorized SSH key in the target hosts in order to access all hosts
listed in `TRENTO_TARGET_AGENTS` to trigger the installation on them

```
-----BEGIN OPENSSH PRIVATE KEY-----
the key
-----END OPENSSH PRIVATE KEY-----
```

- `TRENTO_USER`
The user that will run the CI code to deploy trento in the target machines
```
cloudadmin
```

- `TRENTO_SERVER_IP` 
The IP of the machine where we should run install-server.sh on
```
10.x.x.x
```


- `TRENTO_TARGET_AGENTS`
A comma-separated list of hosts where to run install-agent.sh on (cluster nodes)
```
10.x.x.x,10.x.x.x,10.x.x.x
```