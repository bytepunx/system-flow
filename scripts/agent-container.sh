#!/usr/bin/env sh
# Run Claude Code headless for one story in flai-agent-claude-code:local, started by hand (S-0505).
#   scripts/agent-container.sh S-nnnn [prompt-file]     the prompt from the file, or from stdin
#   scripts/agent-container.sh S-nnnn -- cmd [arg]...   cmd instead of claude, in the same mounts
# Mounts the story's worktree read-write at its host path (the working directory),
# .flai-cache/agents/S-nnnn read-write at its host path as CLAUDE_CONFIG_DIR, and the operator's
# Claude login (FLAI_CLAUDE_LOGIN, default ~/.claude) read-only at /run/claude-login; nothing else,
# never the Docker socket. The agent reaches the project's HTTP MCP server, the one `flai mcp
# status` names, with the token from `flai mcp token`, through .flai-cache/agents/S-nnnn/mcp.json.
# FLAI_BIN       the flai that answers mcp status and token (default: the flai on PATH, installed)
# FLAI_AGENT_NETWORK  bridge (default: host.docker.internal by host-gateway) or host
# FLAI_AGENT_MODEL    claude's --model, such as haiku (default: claude's own)
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
  *) echo "usage: $0 S-nnnn [prompt-file | -- cmd [arg]...]" >&2; exit 2 ;;
esac
PROMPT=""
if [ "${1:-}" = "--" ]; then
  shift
  [ "$#" -gt 0 ] || { echo "nothing to run after --" >&2; exit 2; }
elif [ "$#" -gt 1 ]; then
  echo "usage: $0 S-nnnn [prompt-file | -- cmd [arg]...]" >&2; exit 2
elif [ "$#" -eq 1 ]; then
  PROMPT="$1"
  [ -f "$PROMPT" ] || { echo "no prompt file $PROMPT" >&2; exit 2; }
  set --
elif [ -t 0 ]; then
  echo "no prompt: name a prompt file, or pipe the prompt on stdin" >&2; exit 2
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
SESSION="${FLAI_SESSION:-$(cat /proc/sys/kernel/random/uuid 2>/dev/null || date +%s)}"
mkdir -p "$AGENT_DIR"
chmod 700 "$AGENT_DIR"
# Claude Code finds its login in CLAUDE_CONFIG_DIR: a link to the read-only mount, not a copy.
ln -sfn "$LOGIN_IN/.credentials.json" "$AGENT_DIR/.credentials.json"
( umask 077 && printf '{"mcpServers":{"flai":{"type":"http","url":"%s","headers":{"Authorization":"Bearer %s","Host":"%s","X-Flai-Agent":"%s"}}}}\n' \
  "$URL" "$TOKEN" "$HOST_HEADER" "$AGENT" > "$AGENT_DIR/mcp.json" )
GIT_ARGS=""
if [ "${FLAI_AGENT_GIT_MOUNT:-0}" = 1 ]; then
  GIT_ARGS="-v $CACHE_ROOT/.git:$CACHE_ROOT/.git"
fi
if [ "$#" -eq 0 ]; then
  # As flai serve starts it (flai/internal/harness/adapters.go), with the server over HTTP.
  set -- claude -p --output-format stream-json --verbose \
    --mcp-config "$AGENT_DIR/mcp.json" --strict-mcp-config \
    --permission-prompt-tool mcp__flai__permission_prompt \
    --permission-mode acceptEdits --allowedTools "Bash,mcp__flai"
  [ -z "${FLAI_AGENT_MODEL:-}" ] || set -- "$@" --model "$FLAI_AGENT_MODEL"
fi
echo "running $STORY in $IMAGE on the $NETWORK network: $1" >&2
# NET_ARGS and GIT_ARGS split on spaces by design: neither holds a path with one.
# shellcheck disable=SC2086
docker run --rm -i --user "$(id -u):$(id -g)" $NET_ARGS \
  -v "$WORKTREE:$WORKTREE" -w "$WORKTREE" \
  -v "$AGENT_DIR:$AGENT_DIR" \
  -v "$LOGIN:$LOGIN_IN:ro" $GIT_ARGS \
  -e "CLAUDE_CONFIG_DIR=$AGENT_DIR" -e "FLAI_AGENT=$AGENT" -e "FLAI_SESSION=$SESSION" \
  -e "FLAI_STORY=$STORY" \
  "$IMAGE" "$@" < "${PROMPT:-/dev/stdin}"
