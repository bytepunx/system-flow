#!/usr/bin/env sh
# A local LiteLLM proxy for this project's agents, in Docker, with Postgres
# behind it so that virtual keys and /spend/logs work (E-0019, S-0356, S-0357).
#
#   scripts/litellm.sh up        start (or update) the proxy and write the key file
#   scripts/litellm.sh down      stop it; --purge also drops the spend database
#   scripts/litellm.sh status    is it up, which models, the key's spend so far
#   scripts/litellm.sh key       print the key file; --rotate issues a new virtual key
#   scripts/litellm.sh spend     the virtual key's rows from /spend/logs, as JSON
#   scripts/litellm.sh logs      follow the proxy's logs
#
# Everything it writes lives in .flai-cache/litellm/ (git-ignored), in the main
# checkout when run from a story worktree:
#
#   admin.env     the master key, the salt key, the Postgres password, and the
#                 upstream ANTHROPIC_API_KEY; read by docker compose and this
#                 script only, never exported to agents (mode 0600). Only the
#                 master key may read /spend/logs (a virtual key gets 401), so
#                 `spend` reads it here for whoever holds the virtual key.
#   litellm.env   THE KEY FILE: LITELLM_BASE_URL, LITELLM_API_KEY (a virtual key
#                 with a budget), and FLAI_TEST_LITELLM_HAIKU, as the gateway
#                 tests and a provider entry want them (mode 0600). env.sh
#                 sources it, so every script, scripts/flai.sh serve, and each
#                 agent flai serve starts has them. In any other shell:
#                 . .flai-cache/litellm/litellm.env
#   config.yaml   the proxy's model list: every Anthropic model, upstream to
#                 Anthropic with the key in admin.env
#   compose.yaml  the two services, bound to 127.0.0.1 only
#
# Inputs, on the first `up` (kept in admin.env afterwards):
#   ANTHROPIC_API_KEY     required: the upstream key the proxy spends
#   LITELLM_PORT          host port, default 4000
#   LITELLM_IMAGE         default ghcr.io/berriai/litellm:main-stable; pin a tag
#   LITELLM_MAX_BUDGET    the virtual key's budget in USD per 30 days, default 5
#   LITELLM_HAIKU_MODEL   default anthropic/claude-haiku-4-5
#
# POSIX sh. Needs docker with the compose plugin, curl, jq, and openssl.
set -eu
. "$(dirname "$0")/env.sh"

DIR="$CACHE_ROOT/.flai-cache/litellm"
ADMIN="$DIR/admin.env"
KEYFILE="$DIR/litellm.env"
CONFIG="$DIR/config.yaml"
COMPOSE="$DIR/compose.yaml"
PROJECT="system-flow-litellm"
ALIAS="system-flow-agents"

say() { echo "litellm: $*"; }
die() { echo "litellm: $*" >&2; exit 1; }
need() { command -v "$1" >/dev/null 2>&1 || die "$1 is not installed"; }
compose() { docker compose -p "$PROJECT" --env-file "$ADMIN" -f "$COMPOSE" "$@"; }

# read VAR FILE: the value of VAR in an env file, or empty.
readvar() { sed -n "s/^$1=//p" "$2" 2>/dev/null | head -n1; }

# admin: the admin env file, written once with fresh secrets; later runs may
# change the port, image, budget, or upstream key by setting the variable.
admin() {
  mkdir -p "$DIR"
  if [ ! -f "$ADMIN" ]; then
    [ -n "${ANTHROPIC_API_KEY:-}" ] || die "set ANTHROPIC_API_KEY, the upstream key the proxy spends, for the first up"
    ( umask 077; cat > "$ADMIN" <<EOT
LITELLM_MASTER_KEY=sk-$(openssl rand -hex 24)
LITELLM_SALT_KEY=sk-$(openssl rand -hex 24)
POSTGRES_PASSWORD=$(openssl rand -hex 16)
ANTHROPIC_API_KEY=$ANTHROPIC_API_KEY
LITELLM_PORT=${LITELLM_PORT:-4000}
LITELLM_IMAGE=${LITELLM_IMAGE:-ghcr.io/berriai/litellm:main-stable}
LITELLM_MAX_BUDGET=${LITELLM_MAX_BUDGET:-5}
LITELLM_HAIKU_MODEL=${LITELLM_HAIKU_MODEL:-anthropic/claude-haiku-4-5}
EOT
    )
    say "wrote $ADMIN (master key, salt key, database password, upstream key)"
  else
    for v in ANTHROPIC_API_KEY LITELLM_PORT LITELLM_IMAGE LITELLM_MAX_BUDGET LITELLM_HAIKU_MODEL; do
      eval "val=\${$v:-}"
      if [ -n "$val" ] && [ "$val" != "$(readvar "$v" "$ADMIN")" ]; then
        tmp="$ADMIN.tmp"
        ( umask 077; grep -v "^$v=" "$ADMIN" > "$tmp"; echo "$v=$val" >> "$tmp" ); mv "$tmp" "$ADMIN"
        say "set $v in $ADMIN"
      fi
    done
  fi
  MASTER_KEY="$(readvar LITELLM_MASTER_KEY "$ADMIN")"
  PORT="$(readvar LITELLM_PORT "$ADMIN")"
  MAX_BUDGET="$(readvar LITELLM_MAX_BUDGET "$ADMIN")"
  HAIKU="$(readvar LITELLM_HAIKU_MODEL "$ADMIN")"
  BASE_URL="http://127.0.0.1:$PORT"
}

# files: the proxy's config and the compose file, rewritten every up so that a
# newer version of this script takes effect.
files() {
  cat > "$CONFIG" <<'EOT'
# Written by scripts/litellm.sh; edit the script, not this file.
model_list:
  # The names Claude Code sends, as this project's agents name their models,
  # each upstream to Anthropic with the key in admin.env.
  - model_name: claude-opus-5-5
    litellm_params: { model: anthropic/claude-opus-5-5, api_key: os.environ/ANTHROPIC_API_KEY }
  - model_name: claude-sonnet-5-5
    litellm_params: { model: anthropic/claude-sonnet-5-5, api_key: os.environ/ANTHROPIC_API_KEY }
  - model_name: claude-haiku-4-5
    litellm_params: { model: anthropic/claude-haiku-4-5, api_key: os.environ/ANTHROPIC_API_KEY }
  - model_name: claude-haiku-4-5-20251001
    litellm_params: { model: anthropic/claude-haiku-4-5-20251001, api_key: os.environ/ANTHROPIC_API_KEY }
  # Any other Anthropic model, by its provider-prefixed name (ADR-0129's models map).
  - model_name: "anthropic/*"
    litellm_params: { model: "anthropic/*", api_key: os.environ/ANTHROPIC_API_KEY }
general_settings:
  master_key: os.environ/LITELLM_MASTER_KEY
  database_url: os.environ/DATABASE_URL
  store_model_in_db: false
litellm_settings:
  drop_params: true
EOT
  cat > "$COMPOSE" <<'EOT'
# Written by scripts/litellm.sh; edit the script, not this file.
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: litellm
      POSTGRES_USER: litellm
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U litellm -d litellm"]
      interval: 2s
      timeout: 3s
      retries: 30
    restart: unless-stopped
  litellm:
    image: ${LITELLM_IMAGE}
    command: ["--config", "/app/config.yaml", "--port", "4000"]
    environment:
      LITELLM_MASTER_KEY: ${LITELLM_MASTER_KEY}
      LITELLM_SALT_KEY: ${LITELLM_SALT_KEY}
      ANTHROPIC_API_KEY: ${ANTHROPIC_API_KEY}
      DATABASE_URL: postgresql://litellm:${POSTGRES_PASSWORD}@db:5432/litellm
      STORE_MODEL_IN_DB: "False"
    volumes:
      - ./config.yaml:/app/config.yaml:ro
    ports:
      - "127.0.0.1:${LITELLM_PORT}:4000"
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped
volumes:
  pgdata:
EOT
}

# api METHOD PATH [JSON]: call the proxy as the admin; prints the body, fails on
# a non-2xx status with the body on stderr.
api() {
  method="$1"; path="$2"; body="${3:-}"
  tmp="$(mktemp)"
  code="$(curl -sS -o "$tmp" -w '%{http_code}' -X "$method" "$BASE_URL$path" \
    -H "Authorization: Bearer $MASTER_KEY" -H "Content-Type: application/json" \
    ${body:+-d "$body"})" || { rm -f "$tmp"; return 1; }
  case "$code" in
    2*) cat "$tmp"; rm -f "$tmp" ;;
    *) cat "$tmp" >&2; echo >&2; rm -f "$tmp"; return 1 ;;
  esac
}

wait_ready() {
  say "waiting for $BASE_URL to be ready"
  i=0
  while [ $i -lt 90 ]; do
    if api GET /health/readiness >/dev/null 2>&1; then say "ready"; return 0; fi
    i=$((i + 1)); sleep 2
  done
  die "the proxy did not become ready in 180s; see scripts/litellm.sh logs"
}

# keyfile: write litellm.env with a virtual key that still exists, or a new one.
keyfile() {
  rotate="${1:-}"
  key="$(readvar LITELLM_API_KEY "$KEYFILE")"
  if [ -n "$key" ] && [ -z "$rotate" ] && api GET "/key/info?key=$key" >/dev/null 2>&1; then
    say "keeping the virtual key in $KEYFILE"
  else
    if [ -n "$key" ]; then
      api POST /key/delete "{\"keys\":[\"$key\"]}" >/dev/null 2>&1 && say "deleted the old virtual key" || true
    fi
    key="$(api POST /key/generate "{\"key_alias\":\"$ALIAS-$(date -u +%Y%m%dT%H%M%SZ)\",\"max_budget\":$MAX_BUDGET,\"budget_duration\":\"30d\",\"metadata\":{\"project\":\"system-flow\",\"made_by\":\"scripts/litellm.sh\"}}" | jq -r .key)"
    [ -n "$key" ] && [ "$key" != null ] || die "no key in the proxy's answer to /key/generate"
    say "issued a virtual key with a budget of $MAX_BUDGET USD per 30 days"
  fi
  ( umask 077; cat > "$KEYFILE" <<EOT
# Written by scripts/litellm.sh up; sourced by scripts/env.sh. In another
# shell: . $KEYFILE
# The key is a LiteLLM virtual key with a budget, not the master key, which
# stays in admin.env beside it.
LITELLM_BASE_URL=$BASE_URL
LITELLM_API_KEY=$key
FLAI_TEST_LITELLM_HAIKU=$HAIKU
export LITELLM_BASE_URL LITELLM_API_KEY FLAI_TEST_LITELLM_HAIKU
EOT
  )
  say "wrote $KEYFILE"
}

cmd_up() {
  need docker; need curl; need jq; need openssl
  docker compose version >/dev/null 2>&1 || die "docker compose (the v2 plugin) is not installed"
  admin; files
  say "starting $PROJECT on $BASE_URL"
  compose up -d --remove-orphans
  wait_ready
  keyfile
  cat <<EOT

The proxy is up on $BASE_URL. Key file: $KEYFILE
Scripts in this repository, scripts/flai.sh serve, and the agents it starts see
it through scripts/env.sh. Elsewhere, run:  . $KEYFILE
Claude Code by hand:  ANTHROPIC_BASE_URL=\$LITELLM_BASE_URL ANTHROPIC_AUTH_TOKEN=\$LITELLM_API_KEY ANTHROPIC_API_KEY= claude
Provider entry for the host's providers map (ADR-0129; the command that sets it is S-0349's):
  litellm: { api: anthropic-messages, base_url: $BASE_URL, key_env: LITELLM_API_KEY }
Spend so far:  scripts/litellm.sh status;  the rows:  scripts/litellm.sh spend
EOT
}

cmd_down() {
  need docker
  [ -f "$ADMIN" ] || die "nothing to stop: $ADMIN does not exist"
  admin; files
  if [ "${1:-}" = "--purge" ]; then
    say "stopping and dropping the spend database"
    compose down --volumes --remove-orphans
  else
    say "stopping; the spend database is kept (--purge drops it)"
    compose down --remove-orphans
  fi
}

cmd_status() {
  need docker; need curl; need jq
  [ -f "$ADMIN" ] || die "not set up: run scripts/litellm.sh up"
  admin
  compose ps --format 'table {{.Name}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || true
  if ! api GET /health/liveliness >/dev/null 2>&1; then say "the proxy on $BASE_URL does not answer"; exit 1; fi
  say "proxy on $BASE_URL answers"
  say "models: $(api GET /v1/models | jq -r '[.data[].id] | join(", ")')"
  key="$(readvar LITELLM_API_KEY "$KEYFILE")"
  if [ -n "$key" ]; then
    api GET "/key/info?key=$key" | jq -r '"key \(.info.key_alias // "?"): spent \(.info.spend // 0) of \(.info.max_budget // "unlimited") USD"' | sed 's/^/litellm: /'
  else
    say "no key file yet: run scripts/litellm.sh up"
  fi
}

cmd_key() {
  need curl; need jq
  [ -f "$ADMIN" ] || die "not set up: run scripts/litellm.sh up"
  admin
  if [ "${1:-}" = "--rotate" ]; then keyfile rotate; fi
  [ -f "$KEYFILE" ] || die "no key file yet: run scripts/litellm.sh up"
  cat "$KEYFILE"
}

# spend: the virtual key's rows from /spend/logs, newest first, as JSON. The
# proxy files each row under the SHA-256 of the key, and only the master key
# may read the log, so this reads it for whoever holds the virtual key.
cmd_spend() {
  need curl; need jq; need openssl
  [ -f "$ADMIN" ] || die "not set up: run scripts/litellm.sh up"
  admin
  key="$(readvar LITELLM_API_KEY "$KEYFILE")"
  [ -n "$key" ] || die "no key file yet: run scripts/litellm.sh up"
  hash="$(printf '%s' "$key" | openssl dgst -sha256 | sed 's/^.*= *//')"
  api GET "/spend/logs?api_key=$hash" \
    | jq '[.[] | {startTime, status, model, spend, prompt_tokens, completion_tokens, total_tokens, session_id, request_id, request_tags}]'
}

cmd_logs() {
  need docker
  [ -f "$ADMIN" ] || die "not set up: run scripts/litellm.sh up"
  admin; files
  compose logs -f litellm
}

case "${1:-}" in
  up) shift; cmd_up "$@" ;;
  down) shift; cmd_down "$@" ;;
  status) shift; cmd_status "$@" ;;
  key) shift; cmd_key "$@" ;;
  spend) shift; cmd_spend "$@" ;;
  logs) shift; cmd_logs "$@" ;;
  *) sed -n '2,21p' "$0" | sed 's/^# \{0,1\}//'; exit 2 ;;
esac
