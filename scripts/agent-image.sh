#!/usr/bin/env sh
# Build this repository's agent images from the tree, in order: bedrock (flai and pilot compiled from
# this tree, git, tini) at its tree target, the claude-code harness image over it, then the project's
# layer, the Dockerfile system-flow.yaml names under agent.image.layer, over the harness image, with
# the build arguments the manifest names under agent.image.args (S-0510, S-0760). Not stamped, not
# signed. Each image is labelled as flai's own tree build labels it (ADR-0295 §2): its flai version,
# through bedrock's FLAI_VERSION; dev.system-flow.recipe, its recipe as flai computes it, from its
# Dockerfile, target, and build arguments; and above bedrock dev.system-flow.base, the ID of the image
# it was built over; so that flai image status judges the images current.
#
# It is this repository's build command for its agent images (ADR-0295 §7, S-0759):
#   flai serve agent image command --project <main checkout> -- scripts/agent-image.sh
# flai runs it in the main checkout with the FLAI_IMAGE_* variables, and it builds exactly the images
# FLAI_IMAGE_STALE names, of bedrock, harness, and layer: each with the tag FLAI_IMAGE_BEDROCK,
# FLAI_IMAGE_HARNESS, or FLAI_IMAGE_LAYER gives and the recipe label FLAI_IMAGE_RECIPE_BEDROCK,
# _HARNESS, or _LAYER gives (the recipe it computes itself when that is empty), flai built at
# FLAI_IMAGE_VERSION. The harness image is tagged flai-agent-claude-code:<version> too, as flai's own
# build tags it. An empty FLAI_IMAGE_STALE builds nothing. It refuses, naming the variable, an inner
# harness other than claude-code in FLAI_IMAGE_INNER, and an image asked whose tag is empty.
# With no FLAI_IMAGE_STALE, as make agent-images runs it, it builds all three for this tree's flai, a
# development build, version dev: flai-agent-bedrock:dev, also tagged flai-agent-bedrock:local, then
# flai-agent-claude-code:local, also tagged :dev, then, when the manifest names a layer,
# flai-agent-<project key>-claude-code:local.
# --layer builds the project's layer alone, over flai-agent-claude-code:local already built, labelled
# the same way, and fails when the manifest names none.
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
INNER="claude-code"
HARNESS="flai-agent-$INNER:local"
BEDROCK_FILE="flai/images/bedrock/Dockerfile"
HARNESS_FILE="flai/images/$INNER/Dockerfile"
LABEL_RECIPE="dev.system-flow.recipe"
LABEL_BASE="dev.system-flow.base"

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

# manifest_args prints the build arguments agent.image.args gives the layer, one NAME=value a line,
# their quotes stripped, read from the manifest as flai writes it (four-space, then six-space
# indents); nothing when it names none.
manifest_args() {
  awk '
    /^[^ #]/ { agent = ($0 ~ /^agent:/); image = 0; args = 0; next }
    agent && /^  [^ #]/ { image = ($0 ~ /^  image:/); args = 0; next }
    agent && image && /^    [^ #]/ { args = ($0 ~ /^    args:/); next }
    agent && image && args && /^      [^ #]/ {
      line = $0; sub(/^      /, "", line); sub(/[ ]+#.*$/, "", line)
      name = line; sub(/:.*$/, "", name)
      value = line; sub(/^[^:]*:[ ]*/, "", value); gsub(/["'\'']/, "", value)
      print name "=" value
    }
  ' system-flow.yaml
}

# layer_tag prints the layer's tag, flai-agent-<key>-claude-code:local, the manifest's key in lower
# case, as flai names the image (manifest.LayerImage), since docker refuses upper case.
layer_tag() {
  key="$(sed -n 's/^key:[ ]*//p' system-flow.yaml | tr -d "\"'" | tr '[:upper:]' '[:lower:]')"
  [ -n "$key" ] || { echo "agent-image: system-flow.yaml names no key, which the layer's tag needs" >&2; exit 1; }
  echo "flai-agent-$key-$INNER:local"
}

# recipe prints the recipe of the Dockerfile $1 built at the target $2 ("" for none) with the build
# arguments on standard input, one NAME=value a line, as flai computes it (agentimage.Recipe): the
# sha256 of the Dockerfile's bytes, then "\n--target=<target>\n", then NAME=value and a newline for
# each argument, by name in byte order, BASE and FLAI_VERSION left out. Bedrock's FLAI_COMMIT and
# FLAI_DATE are not given it, as flai's tree build leaves them out of the recipe (D-0608).
recipe() {
  {
    cat "$1"
    printf '\n--target=%s\n' "$2"
    awk -F= 'NF && $1 != "BASE" && $1 != "FLAI_VERSION"' | LC_ALL=C sort -t= -k1,1
  } | sha256_hex | sed 's/^/sha256:/'
}

# sha256_hex prints the sha256 of standard input in hex.
sha256_hex() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum
  else
    shasum -a 256
  fi | cut -d' ' -f1
}

# image_id prints the ID of the image $1 on this host, which an image built over it is labelled
# with, and fails when docker has none.
image_id() {
  id="$(docker image inspect --format '{{.Id}}' "$1" 2>/dev/null)" || id=""
  [ -n "$id" ] || { echo "agent-image: docker has no image $1 on this host to build over; build it first" >&2; exit 1; }
  echo "$id"
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

# build_bedrock builds bedrock's tree target from this tree, flai at the version $1, labelled with
# the recipe $2, or the one computed when it is empty, and tagged with each of the rest.
build_bedrock() {
  version="$1"
  want="$2"
  shift 2
  [ -n "$want" ] || want="$(recipe "$BEDROCK_FILE" tree </dev/null)"
  commit="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
  # The commit's date, not the clock's, so an unchanged tree rebuilds from the cache.
  date="$(TZ=UTC git log -1 --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo)"
  echo "building $* from this tree, flai $version at commit $commit"
  # Each tag becomes -t <tag>: the loop walks the tags as they were, shifting each off in turn.
  for tag; do
    set -- "$@" -t "$tag"
    shift
  done
  docker build -f "$BEDROCK_FILE" --target tree --build-arg "FLAI_VERSION=$version" \
    --build-arg "FLAI_COMMIT=$commit" --build-arg "FLAI_DATE=$date" \
    --label "$LABEL_RECIPE=$want" "$@" .
}

# build_harness builds the harness image over the bedrock image $1, labelled with the recipe $2, or
# the one computed when it is empty, and with bedrock's ID, and tagged with each of the rest.
build_harness() {
  bedrock="$1"
  want="$2"
  shift 2
  [ -n "$want" ] || want="$(recipe "$HARNESS_FILE" "" </dev/null)"
  base_id="$(image_id "$bedrock")"
  echo "building $* over $bedrock ($base_id)"
  for tag; do
    set -- "$@" -t "$tag"
    shift
  done
  docker build -f "$HARNESS_FILE" --build-arg "BASE=$bedrock" \
    --label "$LABEL_RECIPE=$want" --label "$LABEL_BASE=$base_id" "$@" .
}

# build_layer builds the layer, the Dockerfile $1, tagged $2, over the harness image $3, with the
# build arguments agent.image.args gives, labelled with the recipe $4, or the one computed when it
# is empty, and with the harness image's ID.
build_layer() {
  file="$1"
  tag="$2"
  base="$3"
  want="$4"
  [ -f "$file" ] || { echo "agent-image: agent.image.layer names $file, which is not a file; fix the path in system-flow.yaml" >&2; exit 1; }
  args="$(manifest_args | LC_ALL=C sort -t= -k1,1)"
  [ -n "$want" ] || want="$(printf '%s\n' "$args" | recipe "$file" "")"
  base_id="$(image_id "$base")"
  set -- build -f "$file" --build-arg "BASE=$base"
  while IFS= read -r arg; do
    [ -z "$arg" ] || set -- "$@" --build-arg "$arg"
  done <<EOF
$args
EOF
  shown="$(printf '%s\n' "$args" | tr '\n' ' ')"
  shown="${shown% }"
  echo "building $tag from $file over $base ($base_id) with the build arguments ${shown:-none}"
  docker "$@" --label "$LABEL_RECIPE=$want" --label "$LABEL_BASE=$base_id" -t "$tag" .
  echo "built $tag"
}

# make_images builds the images FLAI_IMAGE_STALE names, in order, as the FLAI_IMAGE_* variables
# say them (ADR-0295 §7), bedrock also tagged BEDROCK_ALSO when it is set.
make_images() {
  want_bedrock=0
  want_harness=0
  want_layer=0
  # shellcheck disable=SC2086 # FLAI_IMAGE_STALE splits into one word per image
  for part in $FLAI_IMAGE_STALE; do
    case "$part" in
      bedrock) want_bedrock=1 ;;
      harness) want_harness=1 ;;
      layer) want_layer=1 ;;
      *) echo "agent-image: FLAI_IMAGE_STALE names $part, which is none of bedrock, harness, and layer" >&2; exit 2 ;;
    esac
  done
  if [ "$want_bedrock$want_harness$want_layer" = 000 ]; then
    echo "agent-image: FLAI_IMAGE_STALE names no image, so there is none to build"
    return 0
  fi
  [ "${FLAI_IMAGE_INNER:-}" = "$INNER" ] || { echo "agent-image: FLAI_IMAGE_INNER is ${FLAI_IMAGE_INNER:-empty}, but this script builds the images of $INNER alone" >&2; exit 1; }
  version="${FLAI_IMAGE_VERSION:-}"
  bedrock_tag="${FLAI_IMAGE_BEDROCK:-}"
  harness_tag="${FLAI_IMAGE_HARNESS:-}"
  layer_ref="${FLAI_IMAGE_LAYER:-}"
  if [ "$want_bedrock$want_harness" != 00 ]; then
    [ -n "$version" ] || { echo "agent-image: FLAI_IMAGE_VERSION is empty, but FLAI_IMAGE_STALE asks for bedrock or the harness image, which need the flai version" >&2; exit 1; }
    [ -n "$bedrock_tag" ] || { echo "agent-image: FLAI_IMAGE_BEDROCK is empty, but FLAI_IMAGE_STALE asks for bedrock or the harness image over it" >&2; exit 1; }
  fi
  if [ "$want_harness$want_layer" != 00 ]; then
    [ -n "$harness_tag" ] || { echo "agent-image: FLAI_IMAGE_HARNESS is empty, but FLAI_IMAGE_STALE asks for the harness image or the layer over it" >&2; exit 1; }
  fi
  if [ "$want_layer" = 1 ]; then
    [ -n "$layer_ref" ] || { echo "agent-image: FLAI_IMAGE_LAYER is empty, but FLAI_IMAGE_STALE asks for the layer" >&2; exit 1; }
    [ -n "$layer" ] || { echo "agent-image: FLAI_IMAGE_STALE asks for the layer, but system-flow.yaml names no agent.image.layer" >&2; exit 1; }
  fi
  credentials_or_anonymous
  if [ "$want_bedrock" = 1 ]; then
    build_bedrock "$version" "${FLAI_IMAGE_RECIPE_BEDROCK:-}" "$bedrock_tag" ${BEDROCK_ALSO:+"$BEDROCK_ALSO"}
    echo "built $bedrock_tag"
  fi
  if [ "$want_harness" = 1 ]; then
    # flai's own build tags the harness image flai-agent-<inner>:<version> too (HarnessTags).
    own="flai-agent-$INNER:$version"
    if [ "$own" = "$harness_tag" ]; then
      build_harness "$bedrock_tag" "${FLAI_IMAGE_RECIPE_HARNESS:-}" "$harness_tag"
    else
      build_harness "$bedrock_tag" "${FLAI_IMAGE_RECIPE_HARNESS:-}" "$harness_tag" "$own"
    fi
    echo "built $harness_tag"
  fi
  if [ "$want_layer" = 1 ]; then
    build_layer "$layer" "$layer_ref" "$harness_tag" "${FLAI_IMAGE_RECIPE_LAYER:-}"
  fi
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
  tag="$(layer_tag)"
  credentials_or_anonymous
  build_layer "$layer" "$tag" "$HARNESS" ""
  exit 0
fi
[ "$#" -eq 0 ] || { echo "usage: scripts/agent-image.sh [--layer | --tag | --check [S-nnnn]], or with no argument as flai's build command, with the FLAI_IMAGE_* variables" >&2; exit 2; }
if [ -z "${FLAI_IMAGE_STALE+set}" ]; then
  # By hand, or make agent-images: all three, for this tree's flai, a development build. The
  # recipes are left empty, so that each build computes its own.
  FLAI_IMAGE_VERSION="dev"
  FLAI_IMAGE_INNER="$INNER"
  FLAI_IMAGE_BEDROCK="flai-agent-bedrock:dev"
  FLAI_IMAGE_HARNESS="$HARNESS"
  FLAI_IMAGE_LAYER=""
  FLAI_IMAGE_RECIPE_BEDROCK=""
  FLAI_IMAGE_RECIPE_HARNESS=""
  FLAI_IMAGE_RECIPE_LAYER=""
  FLAI_IMAGE_STALE="bedrock harness"
  if [ -n "$layer" ]; then
    FLAI_IMAGE_LAYER="$(layer_tag)"
    FLAI_IMAGE_STALE="$FLAI_IMAGE_STALE layer"
  fi
  # Today's name for bedrock, which the claude-code Dockerfile's BASE defaults to.
  BEDROCK_ALSO="flai-agent-bedrock:local"
fi
make_images
