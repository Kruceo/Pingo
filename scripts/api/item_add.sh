#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

NAME="${1:-}"
TOOL="${2:-}"
TARGET="${3:-}"
TIMEOUT="${4:-}"

if [[ -z "$NAME" || -z "$TOOL" || -z "$TARGET" || -z "$TIMEOUT" ]]; then
  echo "uso: $0 <name> <tool> <target> <timeout_ms>" >&2
  exit 2
fi

NAME_ESCAPED="$(json_escape "$NAME")"
TOOL_ESCAPED="$(json_escape "$TOOL")"
TARGET_ESCAPED="$(json_escape "$TARGET")"

payload="$(cat <<JSON
{"name":"$NAME_ESCAPED","tool":"$TOOL_ESCAPED","target":"$TARGET_ESCAPED","timeout":$TIMEOUT}
JSON
)"

curl_json -X POST "$BASE_URL/items" -d "$payload"
echo
