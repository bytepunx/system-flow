#!/usr/bin/env sh
# Run a story's agent headless in the agent image, started by hand (S-0505, S-0507, S-0670).
#   scripts/agent-container.sh S-nnnn                  pilot, the image's entrypoint, starts claude
#   scripts/agent-container.sh S-nnnn -- cmd [arg]...  cmd instead of pilot, in the same mounts
#   scripts/agent-container.sh --probe S-nnnn [-- cmd [arg]...]  a trial of the image (I-0224)
# Without --probe the run is the story's agent, agent-S-nnnn, working the live story. With --probe
# it is probe-S-nnnn in a new session: the start request carries "probe":true, so pilot asks it
# only to read and report, and starts it without FLAI_STORY; the -- form passes no FLAI_STORY
# either. The mounts stay the story's, so a probe sees what the story's agent would.
# Feeds pilot a start request on standard input (ADR-0212): the story, its worktree, the agent's
# name and session, inner claude-code, the cache folder, and the names of the variables holding the
# project's HTTP MCP server's URL, token, and Host. pilot writes the MCP file, which names those
# variables and holds no token, and starts claude with the argv flai serve gives on the host, the
# prompt included; its log is claude's stream-json, as on the host.
# Mounts, each at its host path, the story's worktree read-write (the working directory), the main
# checkout's .git read-only and the story's registration in it, .git/worktrees/S-nnnn, read-write
# over it (ADR-0246), the main checkout's wip folder and system-flow.yaml read-only and its
# .flai-cache/verify read-write, made first so it is the operator's, so flai verify and the
# close-out run inside against the main checkout's narrative and board and write the report the
# host reads (ADR-0253), and .flai-cache/agents/S-nnnn read-write as the cache (claude's
# configuration in its claude folder); and the operator's Claude login (FLAI_CLAUDE_LOGIN, default
# ~/.claude) read-only at /run/claude-login; nothing else, never the Docker socket. Git in the
# worktree reads and refreshes its index; a commit fails on the read-only objects, since commits
# are the host's, through task_done. The server is the one `flai mcp status` names, with the
# token from `flai mcp token`, passed by environment, never in an argument.
# The container's flai keeps its config and cache under .flai-cache/agents/S-nnnn/flai (FLAI_CONFIG,
# FLAI_CACHE_DIR), so nothing is written under the container's home.
# FLAI_BIN       the flai that answers mcp status and token (default: the flai on PATH, installed)
# FLAI_AGENT_IMAGE  the image (default: the one flai serve runs, as scripts/agent-image.sh --tag
#                   prints it: the project's layer image when system-flow.yaml names
#                   agent.image.layer, else flai-agent-claude-code:local); a run refuses an image
#                   not built. scripts/agent-container.sh --image prints it and runs nothing.
# FLAI_AGENT_NETWORK  bridge (default: host.docker.internal by host-gateway) or host
# FLAI_AGENT_MODEL    the agent's model, such as haiku (default: claude's own)
# FLAI_AGENT_MAX_BUDGET_USD  the agent's max_budget_usd, such as 0.05, for a trial run
# FLAI_AGENT_SESSION  the session; a second run with the same one resumes it (default: a new one);
#                     ignored with --probe, which always starts a new one
set -eu
# The installed flai, before env.sh puts this tree's bin/ first on PATH.
FLAI_BIN="${FLAI_BIN:-$(command -v flai || true)}"
. "$(dirname "$0")/env.sh"
IMAGE="${FLAI_AGENT_IMAGE:-$("$ROOT/scripts/agent-image.sh" --tag)}"
[ "${1:-}" != "--image" ] || { printf '%s\n' "$IMAGE"; exit 0; }
LOGIN_IN="/run/claude-login"
USAGE="usage: $0 [--probe] S-nnnn [-- cmd [arg]...]"
PROBE=""
if [ "${1:-}" = "--probe" ]; then
  PROBE=1
  shift
fi
STORY="${1:-}"
case "$STORY" in
  S-[0-9][0-9][0-9][0-9]) shift ;;
  *) echo "$USAGE" >&2; exit 2 ;;
esac
if [ "${1:-}" = "--" ]; then
  shift
  [ "$#" -gt 0 ] || { echo "nothing to run after --" >&2; exit 2; }
elif [ "$#" -gt 0 ]; then
  echo "$USAGE" >&2; exit 2
fi

# A run in an image not built would make docker pull it, and fail on a local tag (I-0279).
docker image inspect "$IMAGE" >/dev/null 2>&1 || { echo "no image $IMAGE: build it with scripts/agent-image.sh (make agent-images), or set FLAI_AGENT_IMAGE to one built" >&2; exit 1; }
WORKTREE="$CACHE_ROOT/.flai-cache/worktrees/$STORY"
AGENT_DIR="$CACHE_ROOT/.flai-cache/agents/$STORY"
LOGIN="${FLAI_CLAUDE_LOGIN:-$HOME/.claude}"
[ -d "$WORKTREE" ] || { echo "no worktree at $WORKTREE: flai story start $STORY opens it" >&2; exit 1; }
# The story's registration, which the worktree's .git file names: git in the container writes its
# index, HEAD, and locks there, and reads the rest of the main checkout's .git read-only.
GIT_COMMON="$CACHE_ROOT/.git"
REGISTRATION="$GIT_COMMON/worktrees/$STORY"
[ -d "$REGISTRATION" ] || { echo "no worktree registration at $REGISTRATION: git in the container needs it; flai story start $STORY opens the worktree with it" >&2; exit 1; }
# What flai verify reads and writes of the main checkout (ADR-0253): its wip folder, as the
# manifest's layout names it, its manifest, and the folder its report goes in.
MANIFEST="$CACHE_ROOT/system-flow.yaml"
WIP_NAME="$(sed -n '/^layout:/,/^[^ ]/s/^  wip: *//p' "$MANIFEST" | head -n 1)"
WIP="$CACHE_ROOT/${WIP_NAME:-wip}"
VERIFY="$CACHE_ROOT/.flai-cache/verify"
# The checked-write lock a checked write of a work item on the host holds, so flai check's and flai
# verify's record step in the container waits for it (S-0679).
LOCK="$CACHE_ROOT/.flai-cache/items-checked.lock"
[ -d "$WIP" ] || { echo "no wip folder at $WIP: flai verify in the container reads the main checkout's" >&2; exit 1; }
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
NEW_SESSION="$(cat /proc/sys/kernel/random/uuid 2>/dev/null || date +%s)"
if [ -n "$PROBE" ]; then
  # A probe never resumes the story's session, and works under a name of its own (I-0224).
  AGENT="probe-$STORY"
  SESSION="$NEW_SESSION"
  PROBE_JSON=',"probe":true'
  STORY_ENV=""
else
  AGENT="agent-$STORY"
  SESSION="${FLAI_AGENT_SESSION:-$NEW_SESSION}"
  PROBE_JSON=""
  STORY_ENV="FLAI_STORY=$STORY"
fi
case "$WORKTREE$AGENT_DIR" in
  *[\\\"]*) echo "the worktree or cache path holds a quote or a backslash, which the start request cannot carry" >&2; exit 2 ;;
esac
mkdir -p "$AGENT_DIR/flai" "$AGENT_DIR/claude" "$VERIFY"
# A file, not a folder: Docker would make a missing bind-mount source a folder owned by root.
[ -e "$LOCK" ] || { : > "$LOCK" && chmod 600 "$LOCK"; }
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
REQUEST="$(printf '{"Story":"%s","Root":"%s","Project":"%s","Name":"%s","Session":"%s","Agent":{%s},"inner":"claude-code","mcp_url_env":"FLAI_MCP_URL","mcp_token_env":"FLAI_MCP_TOKEN","mcp_host_env":"FLAI_MCP_HOST","cache":"%s"%s}' \
  "$STORY" "$WORKTREE" "$PROJECT" "$AGENT" "$SESSION" "$AGENT_JSON" "$AGENT_DIR" "$PROBE_JSON")"
# docker run -e NAME, with no value, copies NAME from this environment: the token stays out of
# every argument.
FLAI_MCP_URL="$URL" FLAI_MCP_TOKEN="$TOKEN" FLAI_MCP_HOST="$HOST_HEADER"
export FLAI_MCP_URL FLAI_MCP_TOKEN FLAI_MCP_HOST
# Docker makes the main checkout's own folder as root, as its mounts' parent, so git there, as
# flai verify runs it to find the main branch and go build to stamp a build of the worktree, is told
# by command-line configuration that the folder is safe (ADR-0253, I-0255).
# NET_ARGS splits on spaces by design: it holds no path. The registration's mount comes after the
# read-only .git's, the shorter destination first, so that Docker nests it over its parent.
run() {
  # shellcheck disable=SC2086
  docker run --rm -i --user "$(id -u):$(id -g)" $NET_ARGS \
    -v "$WORKTREE:$WORKTREE" -w "$WORKTREE" \
    -v "$GIT_COMMON:$GIT_COMMON:ro" \
    -v "$REGISTRATION:$REGISTRATION" \
    -v "$WIP:$WIP:ro" \
    -v "$MANIFEST:$MANIFEST:ro" \
    -v "$VERIFY:$VERIFY" \
    -v "$LOCK:$LOCK" \
    -v "$AGENT_DIR:$AGENT_DIR" \
    -v "$LOGIN:$LOGIN_IN:ro" \
    -e FLAI_MCP_URL -e FLAI_MCP_TOKEN -e FLAI_MCP_HOST \
    -e "FLAI_CONFIG=$AGENT_DIR/flai/config.json" -e "FLAI_CACHE_DIR=$AGENT_DIR/flai/cache" \
    -e GIT_CONFIG_COUNT=1 -e GIT_CONFIG_KEY_0=safe.directory -e "GIT_CONFIG_VALUE_0=$CACHE_ROOT" \
    "$@"
}
[ -z "$PROBE" ] || echo "a probe of $IMAGE as $AGENT, in $STORY's mounts: it works no item and has no FLAI_STORY" >&2
if [ "$#" -gt 0 ]; then
  echo "running $STORY in $IMAGE on the $NETWORK network: $1" >&2
  # FLAI_STORY only for the story's agent: with --probe, STORY_ENV is empty and no -e is passed.
  run -e "CLAUDE_CONFIG_DIR=$AGENT_DIR/claude" -e "FLAI_AGENT=$AGENT" -e "FLAI_SESSION=$SESSION" \
    ${STORY_ENV:+-e "$STORY_ENV"} --entrypoint /sbin/tini "$IMAGE" -- "$@"
else
  echo "running $STORY in $IMAGE on the $NETWORK network: pilot, session $SESSION" >&2
  printf '%s\n' "$REQUEST" | run "$IMAGE"
fi
