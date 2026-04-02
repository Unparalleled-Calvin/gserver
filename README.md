## gserver

A toy game server implemented with golang.

### Usage

Generally, you can use `dev.sh` to launch a dump server for testing.

```shell
./dev.sh -b # build the server and client binary
./dev.sh -l # launch the server
./dev.sh -c # run the client
# ./dev.sh -s # stop the server
```

The server address is set by environment variable `GSERVER_ADDR` and the default value is `localhost:8000`.

The redis address is set by environment variable `GSERVER_REDIS_ADDR` and the default value is `localhost:6379`.

### Dependency

- Go 1.26.1 or later
- Redis 7.0.12