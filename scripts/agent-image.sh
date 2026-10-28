#!/usr/bin/env sh
# Build the agent images locally from the repo root, in order: flai-agent-bedrock:local (flai from
# this tree, git, tini), then flai-agent-claude-code:local on top of it. Not stamped, not signed.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
BEDROCK="flai-agent-bedrock:local"
CLAUDE_CODE="flai-agent-claude-code:local"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
VERSION="$(git describe --tags --match 'flai/v*' --abbrev=0 2>/dev/null | sed 's|flai/v||' || echo dev)"
# The commit's date, not the clock's, so an unchanged tree rebuilds from the cache.
DATE="$(TZ=UTC git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo)"
echo "building $BEDROCK"
docker build -f flai/images/bedrock/Dockerfile -t "$BEDROCK" \
  --build-arg "FLAI_VERSION=${VERSION:-dev}" --build-arg "FLAI_COMMIT=$COMMIT" \
  --build-arg "FLAI_DATE=$DATE" .
echo "building $CLAUDE_CODE"
docker build -f flai/images/claude-code/Dockerfile -t "$CLAUDE_CODE" --build-arg "BASE=$BEDROCK" .
echo "built $BEDROCK and $CLAUDE_CODE"
