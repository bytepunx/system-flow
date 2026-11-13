#!/usr/bin/env sh
# Tests agent-image.sh in a scratch copy of this tree's env.sh and agent-image.sh,
# with no Docker: a fake docker logs each call's arguments and answers
# docker image inspect --format '{{.Id}}' <ref> with sha256:id-of-<ref>.
#
# I-0339: when the credential helper the Docker configuration names in
# credsStore fails, as Docker Desktop's desktop.exe does from WSL when Windows
# interop is down, or is not on PATH, agent-image.sh builds with DOCKER_CONFIG
# at a scratch configuration with no login, keeping the current context and the
# folder's other entries, says so, leaves the operator's config.json as it was,
# and removes the scratch configuration. A working helper, or none named,
# leaves DOCKER_CONFIG as it is.
#
# S-0760: as flai's build command (ADR-0295 §7) it builds exactly the images
# FLAI_IMAGE_STALE names, with the tags, the flai version, and the recipe labels
# the FLAI_IMAGE_* variables give, and each base label the ID of the image under
# it; an empty FLAI_IMAGE_STALE builds nothing; it refuses an inner harness
# other than claude-code, and a layer asked with no FLAI_IMAGE_LAYER. With no
# variables it builds all three for flai dev, labelled with the recipes flai
# computes, the layer with the build arguments agent.image.args gives. The
# expected recipes are agentimage.Recipe's for the fixture Dockerfiles below.
set -eu
. "$(dirname "$0")/env.sh"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
repo="$tmp/repo"
mkdir -p "$repo/scripts" "$repo/bin" "$tmp/home/.docker/cli-plugins" "$tmp/home/.docker/contexts"
mkdir -p "$repo/flai/images/bedrock" "$repo/flai/images/claude-code" "$repo/flai/images/layer"
cp "$ROOT/scripts/env.sh" "$ROOT/scripts/agent-image.sh" "$repo/scripts/"
# shellcheck disable=SC2016 # ${BASE} is the Dockerfiles' own
{
  printf 'FROM scratch AS tree\n' >"$repo/flai/images/bedrock/Dockerfile"
  printf 'ARG BASE\nFROM ${BASE}\n' >"$repo/flai/images/claude-code/Dockerfile"
  printf 'ARG BASE\nFROM ${BASE}\nARG TOOL\nARG TOOL2\n' >"$repo/flai/images/layer/Dockerfile"
}
# agentimage.Recipe of the fixtures: bedrock at its tree target, the harness
# image, and the layer with TOOL=x and TOOL2=y (TOOL before TOOL2 by name,
# though TOOL2=y sorts before TOOL=x as a line).
RECIPE_BEDROCK="sha256:236afc2a792ce0fbba4d53073e2b1715235cb2b4b81962127dcd4bb63a55ca45"
RECIPE_HARNESS="sha256:b85f7703cce21eb9e86e7c7ee0c0864400c86a06527ddc06fcc553936b6cb1c4"
RECIPE_LAYER="sha256:14c57b1540efacaad8420f70745488fcf619d9b92457ff6b6b7bd407481c0fcd"
cat >"$repo/system-flow.yaml" <<'EOF'
name: fixture
key: SF
agent:
  harness: claude-code
EOF
unset DOCKER_CONFIG FLAI_IMAGE_VERSION FLAI_IMAGE_INNER FLAI_IMAGE_BEDROCK FLAI_IMAGE_HARNESS \
  FLAI_IMAGE_LAYER FLAI_IMAGE_RECIPE_BEDROCK FLAI_IMAGE_RECIPE_HARNESS FLAI_IMAGE_RECIPE_LAYER FLAI_IMAGE_STALE
HOME="$tmp/home"
LOG="$tmp/docker.log"
export HOME LOG

# The fake docker answers image inspect with a fake ID, logged as inspect
# <ref>, and records, for every other call, its command, the DOCKER_CONFIG it
# sees, its arguments, and that configuration's config.json and entries.
cat >"$repo/bin/docker" <<'EOF'
#!/usr/bin/env sh
if [ "$1" = image ] && [ "$2" = inspect ]; then
  for ref; do :; done
  echo "inspect $ref" >>"$LOG"
  echo "sha256:id-of-$ref"
  exit 0
fi
{
  echo "call $1 DOCKER_CONFIG=${DOCKER_CONFIG:-}"
  echo "args $*"
  if [ -n "${DOCKER_CONFIG:-}" ]; then
    echo "config $(cat "$DOCKER_CONFIG/config.json")"
    ls -A "$DOCKER_CONFIG" | sed 's/^/entry /'
  fi
} >>"$LOG"
EOF
chmod +x "$repo/bin/docker"

fail() {
  printf '%s\n' "$1"
  [ ! -f "$LOG" ] || cat "$LOG"
  echo "agent-image-test: $2"
  exit 1
}

# helper writes the fake credential helper docker-credential-$1, exiting $2.
helper() {
  printf '#!/usr/bin/env sh\nexit %s\n' "$2" >"$repo/bin/docker-credential-$1"
  chmod +x "$repo/bin/docker-credential-$1"
}

# build runs agent-image.sh with a fresh log, and fails when it fails.
build() {
  rm -f "$LOG"
  out="$("$repo/scripts/agent-image.sh" 2>&1)" || fail "$out" "agent-image.sh failed: $1 (I-0339)"
}

# expect_anonymous says both builds ran with a scratch configuration with no
# login, the current context kept and the folder's other entries linked, that
# the operator's config.json is unchanged, and that the scratch one is gone.
expect_anonymous() {
  [ "$(grep -c '^call build ' "$LOG")" -eq 2 ] || fail "$out" "$1: agent-image.sh did not run docker build twice (I-0339)"
  scratch="$(sed -n 's/^call build DOCKER_CONFIG=//p' "$LOG" | sort -u)"
  if [ -z "$scratch" ] || [ "$(printf '%s\n' "$scratch" | wc -l)" -ne 1 ]; then
    fail "$out" "$1: the builds did not run with one DOCKER_CONFIG (I-0339)"
  fi
  [ "$scratch" != "$HOME/.docker" ] || fail "$out" "$1: the builds ran with the operator's configuration (I-0339)"
  [ "$(grep -c '^config {"currentContext": "desktop-linux"}$' "$LOG")" -eq 2 ] || fail "$out" "$1: the scratch config.json is not the current context alone, with no login (I-0339)"
  grep -q '^entry cli-plugins$' "$LOG" || fail "$out" "$1: the scratch configuration lacks cli-plugins (I-0339)"
  grep -q '^entry contexts$' "$LOG" || fail "$out" "$1: the scratch configuration lacks contexts (I-0339)"
  [ ! -e "$scratch" ] || fail "$out" "$1: the scratch configuration $scratch was left behind (I-0339)"
  cmp -s "$HOME/.docker/config.json" "$tmp/config.json" || fail "$out" "$1: the operator's config.json changed (I-0339)"
  case "$out" in
    *"$2"*"building with no login"*"is unchanged"*) ;;
    *) fail "$out" "$1: agent-image.sh did not say it built with no login, naming $2 (I-0339)" ;;
  esac
}

# expect_configured says both builds ran with DOCKER_CONFIG unset, and nothing
# was said of the login.
expect_configured() {
  [ "$(grep -c '^call build DOCKER_CONFIG=$' "$LOG")" -eq 2 ] || fail "$out" "$1: the builds did not run with the operator's configuration (I-0339)"
  case "$out" in
    *"no login"*) fail "$out" "$1: agent-image.sh said it built with no login (I-0339)" ;;
  esac
}

cat >"$HOME/.docker/config.json" <<'EOF'
{
	"auths": {
		"https://index.docker.io/v1/": {}
	},
	"credsStore": "desktop.exe",
	"currentContext": "desktop-linux"
}
EOF
cp "$HOME/.docker/config.json" "$tmp/config.json"

helper desktop.exe 1
build "the helper fails"
expect_anonymous "the helper fails" "docker-credential-desktop.exe"

rm "$repo/bin/docker-credential-desktop.exe"
sed 's/"desktop.exe"/"absent-in-this-test"/' "$tmp/config.json" >"$HOME/.docker/config.json"
cp "$HOME/.docker/config.json" "$tmp/config.json"
build "the helper is not on PATH"
expect_anonymous "the helper is not on PATH" "docker-credential-absent-in-this-test"

helper absent-in-this-test 0
build "the helper answers"
expect_configured "the helper answers"

printf '{"auths": {}}\n' >"$HOME/.docker/config.json"
build "no helper is named"
expect_configured "no helper is named"

# From here the manifest names a layer, with build arguments quoted both ways
# and a comment, and a top-level key after them.
cat >"$repo/system-flow.yaml" <<'EOF'
name: fixture
key: SF
agent:
  harness: claude-code
  image:
    layer: flai/images/layer/Dockerfile
    args:
      TOOL2: "y"
      TOOL: 'x' # the tool
planning:
  args: none
EOF
LAYER_TAG="flai-agent-sf-claude-code:local"

# run runs agent-image.sh with a fresh log and the environment NAME=value
# pairs given, and is what it said; its exit is in status.
run() {
  rm -f "$LOG"
  : >"$LOG"
  status=0
  out="$(env "$@" "$repo/scripts/agent-image.sh" 2>&1)" || status=$?
}

# expect_ok says the run of $1 passed.
expect_ok() {
  [ "$status" -eq 0 ] || fail "$out" "$1: agent-image.sh exited $status"
}

# expect_files says the run of $1 built the Dockerfiles $2, in that order.
expect_files() {
  files="$(sed -n 's/^args build -f \([^ ]*\) .*/\1/p' "$LOG" | tr '\n' ' ')"
  [ "$files" = "$2" ] || fail "$out" "$1: agent-image.sh built '${files% }', not '${2% }'"
}

# expect_build says the run of $1 built an image tagged $2, among whose
# arguments is each of the rest.
expect_build() {
  what="$1"
  tag="$2"
  shift 2
  line=""
  while IFS= read -r l; do
    case "$l " in
      "args build "*" -t $tag "*) line="$l" ;;
    esac
  done <"$LOG"
  [ -n "$line" ] || fail "$out" "$what: agent-image.sh built no image tagged $tag"
  for word; do
    case "$line " in
      *" $word "*) ;;
      *) fail "$out" "$what: the build of $tag lacks $word: $line" ;;
    esac
  done
}

# expect_refused says the run of $1 failed, saying $2, and ran no docker build.
expect_refused() {
  [ "$status" -ne 0 ] || fail "$out" "$1: agent-image.sh passed"
  case "$out" in
    *"$2"*) ;;
    *) fail "$out" "$1: agent-image.sh did not say $2" ;;
  esac
  ! grep -q '^call build ' "$LOG" || fail "$out" "$1: agent-image.sh ran docker build"
}

ALL="flai/images/bedrock/Dockerfile flai/images/claude-code/Dockerfile flai/images/layer/Dockerfile "
FLAI="FLAI_IMAGE_VERSION=9.9.9 FLAI_IMAGE_INNER=claude-code FLAI_IMAGE_BEDROCK=flai-agent-bedrock:9.9.9 FLAI_IMAGE_HARNESS=flai-agent-claude-code:local FLAI_IMAGE_LAYER=$LAYER_TAG"

what="flai asks for all three"
# shellcheck disable=SC2086 # FLAI splits into one NAME=value per variable
run $FLAI FLAI_IMAGE_RECIPE_BEDROCK=sha256:given-bedrock FLAI_IMAGE_RECIPE_HARNESS=sha256:given-harness \
  FLAI_IMAGE_RECIPE_LAYER=sha256:given-layer "FLAI_IMAGE_STALE=bedrock harness layer"
expect_ok "$what"
expect_files "$what" "$ALL"
expect_build "$what" flai-agent-bedrock:9.9.9 "--target tree" "--build-arg FLAI_VERSION=9.9.9" \
  "--label dev.system-flow.recipe=sha256:given-bedrock"
expect_build "$what" flai-agent-claude-code:local "-t flai-agent-claude-code:9.9.9" \
  "--build-arg BASE=flai-agent-bedrock:9.9.9" "--label dev.system-flow.recipe=sha256:given-harness" \
  "--label dev.system-flow.base=sha256:id-of-flai-agent-bedrock:9.9.9"
expect_build "$what" "$LAYER_TAG" "--build-arg BASE=flai-agent-claude-code:local" \
  "--build-arg TOOL=x --build-arg TOOL2=y" "--label dev.system-flow.recipe=sha256:given-layer" \
  "--label dev.system-flow.base=sha256:id-of-flai-agent-claude-code:local"
! grep -q -- '-t flai-agent-bedrock:local' "$LOG" || fail "$out" "$what: agent-image.sh tagged bedrock with a name flai did not give"

what="flai asks for the layer alone, with no recipe"
# shellcheck disable=SC2086 # FLAI splits into one NAME=value per variable
run $FLAI FLAI_IMAGE_STALE=layer
expect_ok "$what"
expect_files "$what" "flai/images/layer/Dockerfile "
expect_build "$what" "$LAYER_TAG" "--build-arg BASE=flai-agent-claude-code:local" \
  "--label dev.system-flow.recipe=$RECIPE_LAYER" "--label dev.system-flow.base=sha256:id-of-flai-agent-claude-code:local"

what="flai asks for none"
# shellcheck disable=SC2086 # FLAI splits into one NAME=value per variable
run $FLAI FLAI_IMAGE_STALE=
expect_ok "$what"
[ ! -s "$LOG" ] || fail "$out" "$what: agent-image.sh ran docker"

what="flai asks for another inner harness"
# shellcheck disable=SC2086 # FLAI splits into one NAME=value per variable
run $FLAI FLAI_IMAGE_INNER=codex "FLAI_IMAGE_STALE=bedrock harness"
expect_refused "$what" "FLAI_IMAGE_INNER is codex"

what="flai asks for a layer it names no tag for"
# shellcheck disable=SC2086 # FLAI splits into one NAME=value per variable
run $FLAI FLAI_IMAGE_LAYER= FLAI_IMAGE_STALE=layer
expect_refused "$what" "FLAI_IMAGE_LAYER is empty"

what="by hand, with no variables"
run
expect_ok "$what"
expect_files "$what" "$ALL"
expect_build "$what" flai-agent-bedrock:dev "-t flai-agent-bedrock:local" "--target tree" \
  "--build-arg FLAI_VERSION=dev" "--label dev.system-flow.recipe=$RECIPE_BEDROCK"
expect_build "$what" flai-agent-claude-code:local "-t flai-agent-claude-code:dev" \
  "--build-arg BASE=flai-agent-bedrock:dev" "--label dev.system-flow.recipe=$RECIPE_HARNESS" \
  "--label dev.system-flow.base=sha256:id-of-flai-agent-bedrock:dev"
expect_build "$what" "$LAYER_TAG" "--build-arg BASE=flai-agent-claude-code:local" \
  "--build-arg TOOL=x --build-arg TOOL2=y" "--label dev.system-flow.recipe=$RECIPE_LAYER" \
  "--label dev.system-flow.base=sha256:id-of-flai-agent-claude-code:local"

what="--layer"
rm -f "$LOG"
out="$("$repo/scripts/agent-image.sh" --layer 2>&1)" || fail "$out" "$what: agent-image.sh failed"
expect_files "$what" "flai/images/layer/Dockerfile "
expect_build "$what" "$LAYER_TAG" "--build-arg BASE=flai-agent-claude-code:local" \
  "--build-arg TOOL=x --build-arg TOOL2=y" "--label dev.system-flow.recipe=$RECIPE_LAYER" \
  "--label dev.system-flow.base=sha256:id-of-flai-agent-claude-code:local"

what="--tag"
out="$("$repo/scripts/agent-image.sh" --tag 2>&1)" || fail "$out" "$what: agent-image.sh failed"
[ "$out" = "$LAYER_TAG" ] || fail "$out" "$what: agent-image.sh printed $out, not $LAYER_TAG"
echo "agent-image-test: passed"
