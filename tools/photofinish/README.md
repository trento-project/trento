# photofinish - a little, handy tool to replay events

This tiny CLI tool aims to fulfill the need to replay some events and get fixtures.

Photofinish reads a `.photofinish.toml` file in the current working directory and:

- It outputs the fixture sets in the TOML file;
- It issues POST requests against the endpoint we give (default: `http://localhost:8081/api/collect`) with the content of the fixture files as request body.

## Usage

```sh
$ photofinish help
photofinish 1.0.0

USAGE:
    photofinish [SUBCOMMAND]

FLAGS:
    -h, --help       Prints help information
    -V, --version    Prints version information

SUBCOMMANDS:
    help    Prints this message or the help of the given subcommand(s)
    list    list available event sets
    run     injects a specific set of events
```

## Example of `.photofinish.toml`
```toml
[first-test-scenario]
files = [
  "../../test/fixtures/discovery/host/expected_published_host_discovery.json",
  "../../test/fixtures/discovery/sap_system/sap_system_discovery_application.json",
  "../../test/fixtures/discovery/subscriptions/expected_published_subscriptions_discovery.json",
]

[second-test-scenario]
files = [
  "third file",
  "fourth-file"
]
```

## "How do I run a fixture set?"
```sh
$ photofinish run first-test-scenario
Successfully POSTed file: ../../test/fixtures/discovery/host/expected_published_host_discovery.json
Successfully POSTed file: ../../test/fixtures/discovery/sap_system/sap_system_discovery_application.json
Successfully POSTed file: ../../test/fixtures/discovery/subscriptions/expected_published_subscriptions_discovery.json
```


