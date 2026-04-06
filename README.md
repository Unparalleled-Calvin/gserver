## gserver

A toy game server implemented with golang.

### Usage

Use `dev.sh` to build, launch, and stop services for local development.

```shell
./dev.sh -b # build the server and client binary
./dev.sh -l # launch the server (auto-start mysql/redis if needed)
./dev.sh -c # run the client
./dev.sh -s # stop the server
./dev.sh -k # stop server + stop mysql/redis + docker container prune
```

The server address and other settings can be configured in `internal/settings/settings.ini`.

### Dependency

- Go 1.26.1 or later
- Redis 7.0.12
- MySQL 8.0