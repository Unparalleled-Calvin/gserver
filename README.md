## gserver

A toy game server implemented with golang.

### Usage

Use `dev.sh` to build, launch, and stop services for local development.

```shell
./dev.sh -b # build the server and client binary
./dev.sh -l # auto-rebuild if needed, then run server in foreground
./dev.sh -L # auto-rebuild if needed, ensure mysql/redis, then launch server in background
./dev.sh -c # auto-rebuild if needed, then run client
./dev.sh -k # stop server only
./dev.sh -K # stop server + stop mysql/redis + docker container prune
```

Lowercase options operate on the server process only, while uppercase options include container management.

The server address and other settings can be configured in `internal/settings/settings.ini`.

### Dependency

- Go 1.26.1 or later
- Redis 7.0.12
- MySQL 8.0