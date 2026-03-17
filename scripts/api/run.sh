#!/usr/bin/env bash
set -euo pipefail

ADDR="${ADDR:-:8080}"
CONFIG="${CONFIG:-config.json}"

exec go run . api --addr "$ADDR" "$CONFIG"
