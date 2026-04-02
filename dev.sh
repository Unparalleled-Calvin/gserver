#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN_DIR="$ROOT_DIR/bin"
PID_FILE="$ROOT_DIR/.gserver.pid"
SERVER_BIN="$BIN_DIR/server"
CLIENT_BIN="$BIN_DIR/client"

usage() {
  cat <<'EOF'
Usage: ./dev.sh [options]

Options:
  -b    Build server and client binaries
  -l    Launch server in background
  -s    Stop background server process
  -c    Run client in foreground
EOF
}

build_all() {
  mkdir -p "$BIN_DIR"
  echo "[build] server"
  go build -o "$SERVER_BIN" ./cmd/server
  echo "[build] client"
  go build -o "$CLIENT_BIN" ./cmd/client
}

is_server_running() {
  [[ -f "$PID_FILE" ]] || return 1
  local pid
  pid="$(cat "$PID_FILE")"
  [[ -n "$pid" ]] || return 1
  kill -0 "$pid" 2>/dev/null
}

launch_server_background() {
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

while getopts ":blsc" opt; do
  case "$opt" in
    b) build_all ;;
    l) launch_server_background ;;
    s) stop_server_background ;;
    c) run_client_foreground ;;
    *)
      usage
      exit 1
      ;;
  esac
done
