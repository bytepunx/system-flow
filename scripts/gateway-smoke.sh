#!/usr/bin/env sh
# Gateway checks: claude -p through each gateway whose key is set, with the
# guard and permission_prompt on (S-0351). Each run spends real money on that
# key, capped by FLAI_TEST_GATEWAY_BUDGET; a gateway whose key is unset is
# skipped, saying so.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
# One line per gateway: its name, the variable its key is in, and its test.
while read -r name var test; do
  eval "key=\${$var:-}"
  if [ -z "$key" ]; then
    echo "smoke: gateway $name skipped: $var is not set"
    continue
  fi
  echo "smoke: gateway $name"
  go test -count=1 -v -run "^$test\$" ./internal/serve/
done <<EOF
OpenRouter OPENROUTER_API_KEY TestGatewayOpenRouter
EOF
