#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

CURRENT_NAME="${1:-}"
NAME="${2:-}"
TOOL="${3:-}"
TARGET="${4:-}"
TIMEOUT="${5:-}"

if [[ -z "$CURRENT_NAME" || -z "$NAME" || -z "$TOOL" || -z "$TARGET" || -z "$TIMEOUT" ]]; then
  echo "uso: $0 <current_name> <name> <tool> <target> <timeout_ms>" >&2
  exit 2
fi

ENCODED="$(urlencode "$CURRENT_NAME")"

NAME_ESCAPED="$(json_escape "$NAME")"
TOOL_ESCAPED="$(json_escape "$TOOL")"
TARGET_ESCAPED="$(json_escape "$TARGET")"

payload="$(cat <<JSON
{"name":"$NAME_ESCAPED","tool":"$TOOL_ESCAPED","target":"$TARGET_ESCAPED","timeout":$TIMEOUT}
JSON
)"

curl_json -X PUT "$BASE_URL/items/$ENCODED" -d "$payload"
echo
