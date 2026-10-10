#!/usr/bin/env sh
# Sourced by the other scripts. Sets PATH and flai environment for this repo.
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export ROOT
# In a story worktree (ADR-0019) the caches, config, and token live in the
# main checkout's .flai-cache; only bin/ is built per tree.
CACHE_ROOT="$ROOT"
if [ -f "$ROOT/.git" ]; then
  gitdir="$(sed -n 's/^gitdir: //p' "$ROOT/.git")"
  case "$gitdir" in
    */.git/worktrees/*) CACHE_ROOT="${gitdir%%/.git/worktrees/*}" ;;
  esac
fi
export CACHE_ROOT
# This tree's bin/ first, then the main checkout's, where install-tools.sh puts
# the pinned tools, so a story worktree lints with them too (I-0001).
TOOLS_PATH="$ROOT/bin"
[ "$CACHE_ROOT" = "$ROOT" ] || TOOLS_PATH="$TOOLS_PATH:$CACHE_ROOT/bin"
export PATH="$TOOLS_PATH:$CACHE_ROOT/.flai-cache/pnpm/node_modules/.bin:$HOME/go/bin:$PATH"
# Node toolchain caches stay inside the repository (safety convention).
export npm_config_cache="$CACHE_ROOT/.flai-cache/npm"
export PNPM_HOME="$CACHE_ROOT/.flai-cache/pnpm-home"
export npm_config_store_dir="$CACHE_ROOT/.flai-cache/pnpm-store"
export FLAI_CONFIG="${FLAI_CONFIG:-$CACHE_ROOT/.flai-cache/config.json}"
export FLAI_CACHE_DIR="${FLAI_CACHE_DIR:-$CACHE_ROOT/.flai-cache/cache}"
mkdir -p "$CACHE_ROOT/.flai-cache" "$ROOT/bin"
# The local LiteLLM proxy's key file, written by scripts/litellm.sh up: the
# gateway tests, flai serve, and the agents it starts read LITELLM_BASE_URL and
# LITELLM_API_KEY from it. An environment that already names the proxy wins.
if [ -z "${LITELLM_BASE_URL:-}" ] && [ -f "$CACHE_ROOT/.flai-cache/litellm/litellm.env" ]; then
  . "$CACHE_ROOT/.flai-cache/litellm/litellm.env"
fi
# First run: create the config with the cache inside the repo, not under ~/.flai.
if [ ! -f "$FLAI_CONFIG" ] && [ -x "$ROOT/bin/flai" ]; then
  "$ROOT/bin/flai" config set cache_dir "$CACHE_ROOT/.flai-cache/cache" >/dev/null 2>&1 || true
fi
