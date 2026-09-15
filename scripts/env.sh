#!/usr/bin/env sh
# Sourced by the other scripts. Sets PATH and flai environment for this repo.
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export ROOT
export PATH="$ROOT/bin:/usr/local/go/bin:$HOME/go/bin:$PATH"
export FLAI_CONFIG="${FLAI_CONFIG:-$ROOT/.flai-cache/config.json}"
mkdir -p "$ROOT/.flai-cache" "$ROOT/bin"
