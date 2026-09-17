#!/usr/bin/env sh
# Sourced by the other scripts. Sets PATH and flai environment for this repo.
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export ROOT
export PATH="$ROOT/bin:/usr/local/go/bin:$HOME/go/bin:$PATH"
export FLAI_CONFIG="${FLAI_CONFIG:-$ROOT/.flai-cache/config.json}"
export FLAI_CACHE_DIR="${FLAI_CACHE_DIR:-$ROOT/.flai-cache/cache}"
mkdir -p "$ROOT/.flai-cache" "$ROOT/bin"
# First run: create the config with the cache inside the repo, not under ~/.flai.
if [ ! -f "$FLAI_CONFIG" ] && [ -x "$ROOT/bin/flai" ]; then
  "$ROOT/bin/flai" config set cache_dir "$ROOT/.flai-cache/cache" >/dev/null 2>&1 || true
fi
