#!/usr/bin/env sh
# Build the agent images locally from the repo root, in order: flai-agent-bedrock:local (flai from
# this tree, git, tini), then flai-agent-claude-code:local on top of it, then, when system-flow.yaml
# names one under agent.image.layer, the project's own layer over the claude-code image, tagged
# flai-agent-<project key>-claude-code:local (S-0510). Not stamped, not signed.
# --layer builds the project's layer alone, over the harness image already built, and fails when
# the manifest names none.
# --tag prints the image a story's agent runs in, and builds nothing: the layer's tag when the
# manifest names a layer, else the harness image's (S-0670). scripts/agent-container.sh runs it.
# --check [S-nnnn] checks the layer's image already built, and builds nothing: it fails when the
# image lacks a tool the test tiers run, then runs the smoke tier inside it, in the story's worktree
# with the mounts the close-out's container has (ADR-0274), scoped to the story and recording no
# issue (I-0281). With no story it checks in the worktree it is run from. CHECK_IMAGE names another
# image to check in place of the layer's, such as an older build kept under another tag.
# Before it builds, it asks the credential helper the Docker configuration names in credsStore; when
# that helper is missing, fails, or does not answer within ten seconds, as Docker Desktop's
# desktop.exe does from WSL when Windows interop is down, it builds with no login, from a scratch
# configuration, since the base images are public, and says so; the operator's is unchanged (I-0339).
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

# credentials_or_anonymous points DOCKER_CONFIG at a scratch configuration with no login when the
# credential helper the Docker configuration names in credsStore is not on PATH, fails, or does not
# answer within ten seconds, so docker build pulls the public base images anonymously (I-0339). The
# scratch configuration keeps the original's current context and links every other entry of its
# folder, such as cli-plugins, contexts, and buildx, so the build runs on the same engine with the
# same builder; the original's config.json is not changed. A working helper, or none named, leaves
# the configuration as it is.
credentials_or_anonymous() {
  config_dir="${DOCKER_CONFIG:-$HOME/.docker}"
  config="$config_dir/config.json"
  [ -f "$config" ] || return 0
  store="$(sed -n 's/.*"credsStore"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$config" | head -n 1)"
  [ -n "$store" ] || return 0
  helper="docker-credential-$store"
  if command -v "$helper" >/dev/null 2>&1; then
    if command -v timeout >/dev/null 2>&1; then
      timeout 10 "$helper" list </dev/null >/dev/null 2>&1 && return 0
    else
      "$helper" list </dev/null >/dev/null 2>&1 && return 0
    fi
    why="failed or did not answer within ten seconds"
  else
    why="is not on PATH"
  fi
  scratch="$(mktemp -d)"
  trap 'rm -rf "$scratch"' EXIT
  for entry in "$config_dir"/* "$config_dir"/.[!.]*; do
    [ -e "$entry" ] || continue
    name="$(basename "$entry")"
    [ "$name" = config.json ] || ln -s "$entry" "$scratch/$name"
  done
  context="$(sed -n 's/.*"currentContext"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$config" | head -n 1)"
  if [ -n "$context" ]; then
    printf '{"currentContext": "%s"}\n' "$context" >"$scratch/config.json"
  else
    printf '{}\n' >"$scratch/config.json"
  fi
  DOCKER_CONFIG="$scratch"
  export DOCKER_CONFIG
  echo "agent-image: the credential helper $helper, which $config names in credsStore, $why; building with no login, from a scratch Docker configuration, since the base images are public. $config is unchanged" >&2
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

# The tools the manifest's test tiers run, which the layer's image must hold on its own PATH: the
# Go tiers', the scripts', and the smoke tier's install test's (curl, openssl, goreleaser, and a ps
# that takes -p, probed below).
CHECK_TOOLS="go gofmt golangci-lint shellcheck node npx pnpm git tar curl openssl goreleaser"

# check_layer checks the layer's image in the worktree of story $1, or of the worktree this runs
# in when $1 is empty.
check_layer() {
  image="${CHECK_IMAGE:-$(layer_tag)}"
  story="$1"
  if [ -z "$story" ]; then
    [ "$CACHE_ROOT" != "$ROOT" ] || { echo "agent-image: --check runs in a story's worktree; run it there, or name the story: scripts/agent-image.sh --check S-nnnn" >&2; exit 2; }
    story="$(basename "$ROOT")"
  fi
  case "$story" in
    S-[0-9][0-9][0-9][0-9]) ;;
    *) echo "agent-image: --check takes a story, S-nnnn, not $story" >&2; exit 2 ;;
  esac
  worktree="$CACHE_ROOT/.flai-cache/worktrees/$story"
  [ -d "$worktree" ] || { echo "agent-image: no worktree at $worktree: flai story start $story opens it" >&2; exit 1; }
  # The mounts of the close-out's container (ADR-0274, as scripts/agent-container.sh mounts them):
  # the worktree, the main checkout's .git read-only with the story's registration over it, its
  # wip folder and manifest read-only, and its verify folder.
  git_common="$CACHE_ROOT/.git"
  registration="$git_common/worktrees/$story"
  [ -d "$registration" ] || { echo "agent-image: no worktree registration at $registration: flai story start $story opens the worktree with it" >&2; exit 1; }
  manifest="$CACHE_ROOT/system-flow.yaml"
  wip_name="$(sed -n '/^layout:/,/^[^ ]/s/^  wip: *//p' "$manifest" | head -n 1)"
  wip="$CACHE_ROOT/${wip_name:-wip}"
  verify="$CACHE_ROOT/.flai-cache/verify"
  mkdir -p "$verify"
  docker image inspect "$image" >/dev/null 2>&1 || { echo "agent-image: no image $image on this host; build it with scripts/agent-image.sh --layer" >&2; exit 1; }

  echo "check: $image holds $CHECK_TOOLS"
  # shellcheck disable=SC2016,SC2086 # $t and $missing expand in the container's shell; CHECK_TOOLS splits into one argument per tool
  docker run --rm --network none --entrypoint sh "$image" -c '
    missing=""
    for t in "$@"; do command -v "$t" >/dev/null 2>&1 || missing="$missing $t"; done
    # scripts/release-server.sh watches its server with ps -o stat= -p, which BusyBox ps refuses.
    ps -o stat= -p "$$" >/dev/null 2>&1 || missing="$missing ps(-p)"
    [ -z "$missing" ] || { echo "agent-image: the image lacks$missing, which the test tiers run" >&2; exit 1; }
  ' sh $CHECK_TOOLS || { echo "agent-image: check of $image failed: add what it lacks to $layer and build it again with scripts/agent-image.sh --layer" >&2; exit 1; }

  echo "check: the smoke tier inside $image, in $worktree"
  # MERGE_QUEUE_STORIES scopes the repository check and the markdown lint to the story, as a
  # close-out's CLOSE_OUT_STORY does, without recording the notes outside it in issues (ADR-0183):
  # a check writes nothing into the worktree.
  docker run --rm --user "$(id -u):$(id -g)" --network host \
    -v "$worktree:$worktree" -w "$worktree" \
    -v "$git_common:$git_common:ro" \
    -v "$registration:$registration" \
    -v "$wip:$wip:ro" \
    -v "$manifest:$manifest:ro" \
    -v "$verify:$verify" \
    -e GIT_CONFIG_COUNT=1 -e GIT_CONFIG_KEY_0=safe.directory -e "GIT_CONFIG_VALUE_0=$CACHE_ROOT" \
    -e "MERGE_QUEUE_STORIES=$story" \
    --entrypoint "$worktree/scripts/smoke.sh" "$image" \
    || { echo "agent-image: the smoke tier failed inside $image; its output is above" >&2; exit 1; }
  echo "checked $image: it holds the tiers' tools and the smoke tier passes inside it"
}

layer="$(manifest_layer)"
if [ "${1:-}" = "--check" ]; then
  [ "$#" -le 2 ] || { echo "usage: scripts/agent-image.sh --check [S-nnnn]" >&2; exit 2; }
  [ -n "$layer" ] || { echo "agent-image: system-flow.yaml names no agent.image.layer, so there is no layer to check" >&2; exit 1; }
  check_layer "${2:-}"
  exit 0
fi
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
  credentials_or_anonymous
  build_layer "$layer"
  exit 0
fi
[ "$#" -eq 0 ] || { echo "usage: scripts/agent-image.sh [--layer | --tag | --check [S-nnnn]]" >&2; exit 2; }
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
VERSION="$(git describe --tags --match 'flai/v*' --abbrev=0 2>/dev/null | sed 's|flai/v||' || echo dev)"
# The commit's date, not the clock's, so an unchanged tree rebuilds from the cache.
DATE="$(TZ=UTC git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo)"
credentials_or_anonymous
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
