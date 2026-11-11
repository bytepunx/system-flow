#!/usr/bin/env sh
# Reproduces I-0255 in the image a story's agent runs in: go build's VCS
# stamping in a story's worktree runs git in the main checkout, since Go takes
# only a .git folder for a repository's root and a worktree's .git is a file.
# In a container Docker makes the main checkout's folder as root, the parent of
# its mounts, so git refuses it for dubious ownership and go build stops with
# exit status 128. The harness's safe.directory for the main checkout
# (ADR-0253) is what lets it through. With a scratch main checkout and story
# worktree mounted as the harness mounts them, a stamped build must fail
# without that safe.directory and pass, stamped, with it.
#
# It needs Docker and the image (AGENT_VCS_IMAGE, else what agent-image.sh
# --tag prints) with go in it, and skips without them, as inside the image.
set -eu
. "$(dirname "$0")/env.sh"

skip() {
  echo "agent-vcs-test: skipped, $1"
  exit 0
}

command -v docker >/dev/null 2>&1 || skip "docker is not on PATH"
docker info >/dev/null 2>&1 || skip "docker does not answer"
image="${AGENT_VCS_IMAGE:-$("$ROOT/scripts/agent-image.sh" --tag)}"
docker image inspect "$image" >/dev/null 2>&1 || skip "no image $image: build it with scripts/agent-image.sh"
docker run --rm --network none --entrypoint /bin/sh "$image" -c 'command -v go' >/dev/null 2>&1 ||
  skip "$image has no go"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
main="$tmp/main"
story=S-0001
worktree="$main/.flai-cache/worktrees/$story"
registration="$main/.git/worktrees/$story"
case "$tmp" in
  *:*) skip "the scratch folder $tmp holds a colon, which a mount cannot carry" ;;
esac

fail() {
  printf '%s\n' "$1"
  echo "agent-vcs-test: $2 (I-0255)"
  exit 1
}

mkdir -p "$main"
git -C "$main" init -q -b main
git -C "$main" -c user.name=test -c user.email=test@example.invalid commit -q --allow-empty -m init
git -C "$main" worktree add -q -b "story/$story" "$worktree"
printf 'module example.invalid/stamp\n\ngo 1.21\n' >"$worktree/go.mod"
printf 'package main\n\nfunc main() {}\n' >"$worktree/main.go"

# build runs a stamped go build in the worktree with the harness's git mounts, as the operator's
# uid, with the -e options given; braces keep each path whole before its colon.
build() {
  docker run --rm --network none --user "$(id -u):$(id -g)" \
    -v "${worktree}:${worktree}" -w "$worktree" \
    -v "${main}/.git:${main}/.git:ro" \
    -v "${registration}:${registration}" \
    "$@" --entrypoint /bin/sh "$image" \
    -c 'go build -buildvcs=true -o /tmp/stamped . && go version -m /tmp/stamped' 2>&1
}

if got="$(build)"; then
  fail "$got" "a stamped build passed in $image with no safe.directory for the main checkout, so this no longer reproduces the cause"
fi
case "$got" in
  *"error obtaining VCS status: exit status 128"*) ;;
  *) fail "$got" "a stamped build with no safe.directory failed otherwise than on git's exit status 128" ;;
esac

got="$(build -e GIT_CONFIG_COUNT=1 -e GIT_CONFIG_KEY_0=safe.directory -e "GIT_CONFIG_VALUE_0=$main")" ||
  fail "$got" "a stamped build failed in $image with the harness's safe.directory for the main checkout"
case "$got" in
  *vcs.revision=*) ;;
  *) fail "$got" "the build with the harness's safe.directory recorded no vcs.revision" ;;
esac
echo "agent-vcs-test: passed in $image"
