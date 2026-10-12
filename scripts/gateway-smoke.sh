#!/usr/bin/env sh
# Gateway checks: claude -p through each gateway whose variables are set, with
# the guard and permission_prompt on (S-0351, S-0356). Each run spends real
# money on that gateway's key, capped by FLAI_TEST_GATEWAY_BUDGET; a gateway
# with a variable unset is skipped, naming each one missing.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
# One line per gateway: its name, its test, and the variables it needs.
while read -r name test vars; do
  missing=""
  for var in $vars; do
    eval "val=\${$var:-}"
    [ -n "$val" ] || missing="${missing:+$missing and }$var"
  done
  if [ -n "$missing" ]; then
    echo "smoke: gateway $name skipped: $missing not set"
    continue
  fi
  echo "smoke: gateway $name"
  go test -count=1 -v -run "^$test\$" ./internal/serve/
done <<EOF
OpenRouter TestGatewayOpenRouter OPENROUTER_API_KEY
LiteLLM TestGatewayLiteLLM LITELLM_BASE_URL LITELLM_API_KEY
EOF
