#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$ROOT_DIR/bin"
PID_FILE="$ROOT_DIR/.gserver.pid"
SERVER_BIN="$BIN_DIR/server"
CLIENT_BIN="$BIN_DIR/client"
CONFIG_SRC="$ROOT_DIR/internal/settings/settings.ini"
CONFIG_DST="$BIN_DIR/settings.ini"
SCHEMA_SQL="$ROOT_DIR/internal/schema/create.sql"

usage() {
  cat <<'EOF'
Usage: ./dev.sh [options]

Options:
  -b    Build server and client binaries
  -l    Launch server in background
  -s    Stop background server process
  -k    Stop server, stop redis/mysql containers, and prune containers
  -c    Run client in foreground
EOF
}

build_all() {
  mkdir -p "$BIN_DIR"
  echo "[build] server"
  go build -o "$SERVER_BIN" ./cmd/server
  echo "[build] client"
  go build -o "$CLIENT_BIN" ./cmd/client
  echo "[copy] settings.ini"
  cp "$CONFIG_SRC" "$CONFIG_DST"
}

is_server_running() {
  [[ -f "$PID_FILE" ]] || return 1
  local pid
  pid="$(cat "$PID_FILE")"
  [[ -n "$pid" ]] || return 1
  kill -0 "$pid" 2>/dev/null
}

is_container_running() {
  local name="$1"
  docker ps --format '{{.Names}}' | grep -Fxq "$name"
}

container_exists() {
  local name="$1"
  docker ps -a --format '{{.Names}}' | grep -Fxq "$name"
}

ensure_required_containers() {
  if is_container_running "mysql-container"; then
    echo "[ok] mysql-container is running"
  elif container_exists "mysql-container"; then
    echo "[start] starting existing mysql-container"
    docker start mysql-container >/dev/null
  else
    echo "[run] creating mysql-container"
    docker run --name mysql-container -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -e MYSQL_DATABASE=gserver -p 3306:3306 -v "$SCHEMA_SQL:/docker-entrypoint-initdb.d/create.sql:ro" -d mysql:8.0 >/dev/null
  fi

  if is_container_running "redis-container"; then
    echo "[ok] redis-container is running"
  elif container_exists "redis-container"; then
    echo "[start] starting existing redis-container"
    docker start redis-container >/dev/null
  else
    echo "[run] creating redis-container"
    docker run --name redis-container -p 6379:6379 -d redis:7.0.12 >/dev/null
  fi
}

launch_server_background() {
  ensure_required_containers

  if [[ ! -x "$SERVER_BIN" ]]; then
    echo "[info] server binary not found, building first"
    build_all
  fi

  if is_server_running; then
    echo "[skip] server already running (pid: $(cat "$PID_FILE"))"
    return
  fi

  nohup "$SERVER_BIN" >/tmp/gserver.log 2>&1 &
  local pid=$!
  echo "$pid" >"$PID_FILE"
  echo "[ok] server started in background (pid: $pid, log: /tmp/gserver.log)"
}

stop_server_background() {
  if ! [[ -f "$PID_FILE" ]]; then
    echo "[skip] no pid file, server may not be running"
    return
  fi

  local pid
  pid="$(cat "$PID_FILE")"

  if [[ -z "$pid" ]]; then
    rm -f "$PID_FILE"
    echo "[skip] empty pid file removed"
    return
  fi

  if kill -0 "$pid" 2>/dev/null; then
    kill "$pid"
    echo "[ok] stopped server process $pid"
  else
    echo "[skip] process $pid not running"
  fi

  rm -f "$PID_FILE"
}

stop_container_if_running() {
  local name="$1"
  if is_container_running "$name"; then
    docker stop "$name" >/dev/null
    echo "[ok] stopped container $name"
  else
    echo "[skip] container $name is not running"
  fi
}

stop_all_and_prune() {
  stop_server_background
  stop_container_if_running "redis-container"
  stop_container_if_running "mysql-container"
  docker container prune -f >/dev/null
  echo "[ok] docker container prune completed"
}

run_client_foreground() {
  if [[ ! -x "$CLIENT_BIN" ]]; then
    echo "[info] client binary not found, building first"
    build_all
  fi

  "$CLIENT_BIN"
}

if [[ $# -eq 0 ]]; then
  usage
  exit 1
fi

while getopts ":blsck" opt; do
  case "$opt" in
    b) build_all ;;
    l) launch_server_background ;;
    s) stop_server_background ;;
    k) stop_all_and_prune ;;
    c) run_client_foreground ;;
    *)
      usage
      exit 1
      ;;
  esac
done
