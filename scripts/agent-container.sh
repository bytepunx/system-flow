#!/usr/bin/env sh
# Run a story's agent headless in flai-agent-claude-code:local, started by hand (S-0505, S-0507).
#   scripts/agent-container.sh S-nnnn                  pilot, the image's entrypoint, starts claude
#   scripts/agent-container.sh S-nnnn -- cmd [arg]...  cmd instead of pilot, in the same mounts
# Feeds pilot a start request on standard input (ADR-0212): the story, its worktree, the agent's
# name and session, inner claude-code, the cache folder, and the names of the variables holding the
# project's HTTP MCP server's URL, token, and Host. pilot writes the MCP file, which names those
# variables and holds no token, and starts claude with the argv flai serve gives on the host, the
# prompt included; its log is claude's stream-json, as on the host.
# Mounts the story's worktree read-write at its host path (the working directory),
# .flai-cache/agents/S-nnnn read-write at its host path as the cache (claude's configuration in its
# claude folder), and the operator's Claude login (FLAI_CLAUDE_LOGIN, default ~/.claude) read-only
# at /run/claude-login; nothing else, never the Docker socket. The server is the one `flai mcp
# status` names, with the token from `flai mcp token`, passed by environment, never in an argument.
# The container's flai keeps its config and cache under .flai-cache/agents/S-nnnn/flai (FLAI_CONFIG,
# FLAI_CACHE_DIR), so nothing is written under the container's home.
# FLAI_BIN       the flai that answers mcp status and token (default: the flai on PATH, installed)
# FLAI_AGENT_NETWORK  bridge (default: host.docker.internal by host-gateway) or host
# FLAI_AGENT_MODEL    the agent's model, such as haiku (default: claude's own)
# FLAI_AGENT_MAX_BUDGET_USD  the agent's max_budget_usd, such as 0.05, for a trial run
# FLAI_AGENT_SESSION  the session; a second run with the same one resumes it (default: a new one)
# FLAI_AGENT_GIT_MOUNT=1  also mounts the main checkout's .git at its host path, for git inside
set -eu
# The installed flai, before env.sh puts this tree's bin/ first on PATH.
FLAI_BIN="${FLAI_BIN:-$(command -v flai || true)}"
. "$(dirname "$0")/env.sh"
IMAGE="flai-agent-claude-code:local"
LOGIN_IN="/run/claude-login"
STORY="${1:-}"
case "$STORY" in
  S-[0-9][0-9][0-9][0-9]) shift ;;
  *) echo "usage: $0 S-nnnn [-- cmd [arg]...]" >&2; exit 2 ;;
esac
if [ "${1:-}" = "--" ]; then
  shift
  [ "$#" -gt 0 ] || { echo "nothing to run after --" >&2; exit 2; }
elif [ "$#" -gt 0 ]; then
  echo "usage: $0 S-nnnn [-- cmd [arg]...]" >&2; exit 2
fi
WORKTREE="$CACHE_ROOT/.flai-cache/worktrees/$STORY"
AGENT_DIR="$CACHE_ROOT/.flai-cache/agents/$STORY"
LOGIN="${FLAI_CLAUDE_LOGIN:-$HOME/.claude}"
[ -d "$WORKTREE" ] || { echo "no worktree at $WORKTREE: flai story start $STORY opens it" >&2; exit 1; }
[ -f "$LOGIN/.credentials.json" ] || { echo "no Claude login at $LOGIN/.credentials.json: sign in with claude on the host, or set FLAI_CLAUDE_LOGIN" >&2; exit 1; }
[ -x "${FLAI_BIN:-/nonexistent}" ] || { echo "no flai to ask for the MCP server: install it, or set FLAI_BIN" >&2; exit 1; }
# The operator's host's server: its state is the main checkout's .flai-cache/mcp-http.json, whatever
# FLAI_CONFIG says.
URL="$(cd "$CACHE_ROOT" && "$FLAI_BIN" --json mcp status | sed -n 's/^ *"url": *"\([^"]*\)".*/\1/p' | head -n 1)"
[ -n "$URL" ] || { echo "the project's HTTP MCP server is not running: flai mcp start, or turn it on in flai serve" >&2; exit 1; }
TOKEN="$(cd "$CACHE_ROOT" && "$FLAI_BIN" mcp token)"
case "$TOKEN" in
  "" | *[!A-Za-z0-9._~+/=-]*) echo "flai mcp token printed no usable token" >&2; exit 1 ;;
esac
# The server's MCP SDK refuses, as DNS rebinding, a request on its loopback listener whose Host is
# not loopback, so the agent sends the Host the server listens on, whatever address it dials.
HOST_HEADER="$(echo "$URL" | sed 's|^https\{0,1\}://\([^/]*\).*|\1|')"
NETWORK="${FLAI_AGENT_NETWORK:-bridge}"
case "$NETWORK" in
  # Docker Desktop's host-gateway forwards to the WSL host's loopback, where the server listens;
  # its host network is its VM's, where nothing listens.
  bridge) NET_ARGS="--add-host host.docker.internal:host-gateway"
          URL="$(echo "$URL" | sed 's|^\(https\{0,1\}://\)[^/:]*|\1host.docker.internal|')" ;;
  host) NET_ARGS="--network host" ;;
  *) echo "FLAI_AGENT_NETWORK is bridge or host, not $NETWORK" >&2; exit 2 ;;
esac
AGENT="agent-$STORY"
SESSION="${FLAI_AGENT_SESSION:-$(cat /proc/sys/kernel/random/uuid 2>/dev/null || date +%s)}"
case "$WORKTREE$AGENT_DIR" in
  *[\\\"]*) echo "the worktree or cache path holds a quote or a backslash, which the start request cannot carry" >&2; exit 2 ;;
esac
mkdir -p "$AGENT_DIR/flai" "$AGENT_DIR/claude"
chmod 700 "$AGENT_DIR"
# Claude Code finds its login in CLAUDE_CONFIG_DIR, the cache's claude folder: a link to the
# read-only mount, not a copy.
ln -sfn "$LOGIN_IN/.credentials.json" "$AGENT_DIR/claude/.credentials.json"
AGENT_JSON=""
[ -z "${FLAI_AGENT_MODEL:-}" ] || AGENT_JSON="\"model\":\"$FLAI_AGENT_MODEL\""
if [ -n "${FLAI_AGENT_MAX_BUDGET_USD:-}" ]; then
  [ -z "$AGENT_JSON" ] || AGENT_JSON="$AGENT_JSON,"
  AGENT_JSON="$AGENT_JSON\"config\":{\"max_budget_usd\":\"$FLAI_AGENT_MAX_BUDGET_USD\"}"
fi
case "$AGENT_JSON" in
  *[\\]*) echo "FLAI_AGENT_MODEL or FLAI_AGENT_MAX_BUDGET_USD holds a backslash" >&2; exit 2 ;;
esac
PROJECT="$(sed -n 's/^key: *//p' "$ROOT/system-flow.yaml" | head -n 1)"
REQUEST="$(printf '{"Story":"%s","Root":"%s","Project":"%s","Name":"%s","Session":"%s","Agent":{%s},"inner":"claude-code","mcp_url_env":"FLAI_MCP_URL","mcp_token_env":"FLAI_MCP_TOKEN","mcp_host_env":"FLAI_MCP_HOST","cache":"%s"}' \
  "$STORY" "$WORKTREE" "$PROJECT" "$AGENT" "$SESSION" "$AGENT_JSON" "$AGENT_DIR")"
GIT_ARGS=""
if [ "${FLAI_AGENT_GIT_MOUNT:-0}" = 1 ]; then
  GIT_ARGS="-v $CACHE_ROOT/.git:$CACHE_ROOT/.git"
fi
# docker run -e NAME, with no value, copies NAME from this environment: the token stays out of
# every argument.
FLAI_MCP_URL="$URL" FLAI_MCP_TOKEN="$TOKEN" FLAI_MCP_HOST="$HOST_HEADER"
export FLAI_MCP_URL FLAI_MCP_TOKEN FLAI_MCP_HOST
# NET_ARGS and GIT_ARGS split on spaces by design: neither holds a path with one.
run() {
  # shellcheck disable=SC2086
  docker run --rm -i --user "$(id -u):$(id -g)" $NET_ARGS \
    -v "$WORKTREE:$WORKTREE" -w "$WORKTREE" \
    -v "$AGENT_DIR:$AGENT_DIR" \
    -v "$LOGIN:$LOGIN_IN:ro" $GIT_ARGS \
    -e FLAI_MCP_URL -e FLAI_MCP_TOKEN -e FLAI_MCP_HOST \
    -e "FLAI_CONFIG=$AGENT_DIR/flai/config.json" -e "FLAI_CACHE_DIR=$AGENT_DIR/flai/cache" \
    "$@"
}
if [ "$#" -gt 0 ]; then
  echo "running $STORY in $IMAGE on the $NETWORK network: $1" >&2
  run -e "CLAUDE_CONFIG_DIR=$AGENT_DIR/claude" -e "FLAI_AGENT=$AGENT" -e "FLAI_SESSION=$SESSION" \
    -e "FLAI_STORY=$STORY" --entrypoint /sbin/tini "$IMAGE" -- "$@"
else
  echo "running $STORY in $IMAGE on the $NETWORK network: pilot, session $SESSION" >&2
  printf '%s\n' "$REQUEST" | run "$IMAGE"
fi
