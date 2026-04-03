## gserver

A toy game server implemented with golang.

### Usage

Generally, you can use `dev.sh` to launch a dump server for testing.

```shell
./dev.sh -b # build the server and client binary
./dev.sh -l # launch the server
./dev.sh -c # run the client
./dev.sh -s # stop the server
```

The server address and redis address can be configured in `internal/settings/settings.ini`.

### Dependency

- Go 1.26.1 or later
- Redis 7.0.12