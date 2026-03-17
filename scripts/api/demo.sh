#!/usr/bin/env bash
set -euo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$HERE/../.." && pwd)"

export BASE_URL="${BASE_URL:-http://localhost:8080}"
export ADDR="${ADDR:-:8080}"

CONFIG_DIR="$(mktemp -d)"
CONFIG="$CONFIG_DIR/config.json"
export CONFIG

cp "$ROOT/config.example.json" "$CONFIG"

api_pid=""
cleanup() {
  if [[ -n "${api_pid:-}" ]]; then
    kill "$api_pid" >/dev/null 2>&1 || true
    wait "$api_pid" >/dev/null 2>&1 || true
  fi
  rm -rf "$CONFIG_DIR"
}
trap cleanup EXIT

(
  cd "$ROOT"
  bash "$HERE/run.sh" >/dev/null 2>&1 &
  api_pid="$!"
  echo "$api_pid" >"$CONFIG_DIR/pid"
)
api_pid="$(cat "$CONFIG_DIR/pid")"

for _ in {1..50}; do
  if curl -sS "$BASE_URL/health" >/dev/null 2>&1; then
    break
  fi
  sleep 0.1
done

echo "== health"
bash "$HERE/health.sh"

echo "== items (before)"
bash "$HERE/items_list.sh"

NAME="Demo Item $(date +%s)"
echo "== add"
bash "$HERE/item_add.sh" "$NAME" "pingv4" "1.1.1.1" "1500"

echo "== get"
bash "$HERE/item_get.sh" "$NAME"

echo "== update (timeout 2500)"
bash "$HERE/item_update.sh" "$NAME" "$NAME" "pingv4" "1.1.1.1" "2500"

echo "== delete"
bash "$HERE/item_delete.sh" "$NAME"
echo

echo "== items (after)"
bash "$HERE/items_list.sh"
