#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

urlencode() {
  local value="${1:-}"
  if command -v node >/dev/null 2>&1; then
    node -p 'encodeURIComponent(process.argv[1])' "$value"
    return 0
  fi
  if command -v python3 >/dev/null 2>&1; then
    python3 -c 'import sys,urllib.parse; print(urllib.parse.quote(sys.argv[1]))' "$value"
    return 0
  fi
  if command -v python >/dev/null 2>&1; then
    python -c 'import sys,urllib; print(urllib.quote(sys.argv[1]))' "$value"
    return 0
  fi
  printf '%s\n' "$value"
}

json_escape() {
  local s="${1:-}"
  s="${s//\\/\\\\}"
  s="${s//\"/\\\"}"
  s="${s//$'\n'/\\n}"
  printf '%s' "$s"
}

curl_json() {
  curl -sS -H 'Content-Type: application/json' "$@"
}
