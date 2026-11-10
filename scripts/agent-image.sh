#!/usr/bin/env sh
# Build the agent images locally from the repo root, in order: flai-agent-bedrock:local (flai from
# this tree, git, tini), then flai-agent-claude-code:local on top of it, then, when system-flow.yaml
# names one under agent.image.layer, the project's own layer over the claude-code image, tagged
# flai-agent-<project key>-claude-code:local (S-0510). Not stamped, not signed.
# --layer builds the project's layer alone, over the harness image already built, and fails when
# the manifest names none.
# --tag prints the image a story's agent runs in, and builds nothing: the layer's tag when the
# manifest names a layer, else the harness image's (S-0670). scripts/agent-container.sh runs it.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
BEDROCK="flai-agent-bedrock:local"
INNER="claude-code"
HARNESS="flai-agent-$INNER:local"

# The Dockerfile agent.image.layer names, relative to the root, read from the manifest as flai
# writes it (two-space indents); empty when it names none.
manifest_layer() {
  awk '
    /^[^ #]/ { agent = ($0 ~ /^agent:/); image = 0; next }
    agent && /^  [^ #]/ { image = ($0 ~ /^  image:/); next }
    agent && image && /^    layer:/ {
      sub(/^    layer:[ ]*/, ""); sub(/[ ]+#.*$/, ""); gsub(/["'\'']/, ""); print; exit
    }
  ' system-flow.yaml
}

# pinned is the default scripts/install-tools.sh gives the version variable $1, NAME="${NAME:-v}".
pinned() {
  grep "^$1=" scripts/install-tools.sh | sed 's/.*:-\(.*\)}".*/\1/'
}

# layer_tag prints the layer's tag, flai-agent-<key>-claude-code:local, the manifest's key in lower
# case, as flai names the image (manifest.LayerImage), since docker refuses upper case.
layer_tag() {
  key="$(sed -n 's/^key:[ ]*//p' system-flow.yaml | tr -d "\"'" | tr '[:upper:]' '[:lower:]')"
  [ -n "$key" ] || { echo "agent-image: system-flow.yaml names no key, which the layer's tag needs" >&2; exit 1; }
  echo "flai-agent-$key-$INNER:local"
}

# build_layer builds the manifest's layer over the harness image, with Go at flai/go.mod's version
# and the tool versions scripts/install-tools.sh pins as build arguments.
build_layer() {
  layer="$1"
  [ -f "$layer" ] || { echo "agent-image: agent.image.layer names $layer, which is not a file; fix the path in system-flow.yaml" >&2; exit 1; }
  tag="$(layer_tag)"
  go_version="$(sed -n 's/^go //p' flai/go.mod)"
  golangci="$(pinned GOLANGCI_LINT_VERSION)"
  shellcheck="$(pinned SHELLCHECK_VERSION)"
  pnpm="$(pinned PNPM_VERSION)"
  goreleaser="$(pinned GORELEASER_VERSION)"
  echo "building $tag from $layer over $HARNESS"
  docker build -f "$layer" -t "$tag" --build-arg "BASE=$HARNESS" \
    --build-arg "GO_VERSION=$go_version" --build-arg "GOLANGCI_LINT_VERSION=$golangci" \
    --build-arg "SHELLCHECK_VERSION=$shellcheck" --build-arg "PNPM_VERSION=$pnpm" \
    --build-arg "GORELEASER_VERSION=$goreleaser" .
  echo "built $tag"
}

layer="$(manifest_layer)"
if [ "${1:-}" = "--tag" ]; then
  [ "$#" -eq 1 ] || { echo "usage: scripts/agent-image.sh --tag" >&2; exit 2; }
  if [ -n "$layer" ]; then
    layer_tag
  else
    echo "$HARNESS"
  fi
  exit 0
fi
if [ "${1:-}" = "--layer" ]; then
  [ -n "$layer" ] || { echo "agent-image: system-flow.yaml names no agent.image.layer, so there is no layer to build" >&2; exit 1; }
  build_layer "$layer"
  exit 0
fi
[ "$#" -eq 0 ] || { echo "usage: scripts/agent-image.sh [--layer | --tag]" >&2; exit 2; }
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
VERSION="$(git describe --tags --match 'flai/v*' --abbrev=0 2>/dev/null | sed 's|flai/v||' || echo dev)"
# The commit's date, not the clock's, so an unchanged tree rebuilds from the cache.
DATE="$(TZ=UTC git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo)"
echo "building $BEDROCK"
docker build -f flai/images/bedrock/Dockerfile -t "$BEDROCK" \
  --build-arg "FLAI_VERSION=${VERSION:-dev}" --build-arg "FLAI_COMMIT=$COMMIT" \
  --build-arg "FLAI_DATE=$DATE" .
echo "building $HARNESS"
docker build -f flai/images/claude-code/Dockerfile -t "$HARNESS" --build-arg "BASE=$BEDROCK" .
echo "built $BEDROCK and $HARNESS"
if [ -n "$layer" ]; then
  build_layer "$layer"
fi
