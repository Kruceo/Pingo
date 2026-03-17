#!/usr/bin/env bash
set -euo pipefail

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/common.sh"

NAME="${1:-}"
if [[ -z "$NAME" ]]; then
  echo "uso: $0 <name>" >&2
  exit 2
fi

ENCODED="$(urlencode "$NAME")"
curl -sS "$BASE_URL/items/$ENCODED"
echo
