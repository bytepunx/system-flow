#!/usr/bin/env sh
# Reproduces I-0279 in a scratch copy of this tree's env.sh, agent-image.sh,
# and agent-container.sh, with no Docker: agent-container.sh runs the image
# flai serve runs, the project's layer image when system-flow.yaml names
# agent.image.layer, not always flai-agent-claude-code:local; FLAI_AGENT_IMAGE
# overrides it; and a run refuses an image not built, naming the build script.
set -eu
. "$(dirname "$0")/env.sh"
unset FLAI_AGENT_IMAGE
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/scripts" "$tmp/bin"
cp "$ROOT/scripts/env.sh" "$ROOT/scripts/agent-image.sh" "$ROOT/scripts/agent-container.sh" "$tmp/scripts/"

fail() {
  printf '%s\n' "$1"
  echo "agent-container-test: $2 (I-0279)"
  exit 1
}

# expect_image says agent-container.sh --image prints $1.
expect_image() {
  got="$("$tmp/scripts/agent-container.sh" --image 2>&1)" || fail "$got" "agent-container.sh --image failed"
  [ "$got" = "$1" ] || fail "$got" "agent-container.sh --image printed the above, not $1"
}

cat >"$tmp/system-flow.yaml" <<'EOF'
name: fixture
key: SF
agent:
  harness: claude-code
  image:
    layer: flai/images/x/Dockerfile
EOF
expect_image "flai-agent-sf-claude-code:local"

got="$(FLAI_AGENT_IMAGE=my-image:1 "$tmp/scripts/agent-container.sh" --image 2>&1)" || fail "$got" "agent-container.sh --image failed with FLAI_AGENT_IMAGE"
[ "$got" = "my-image:1" ] || fail "$got" "agent-container.sh --image ignored FLAI_AGENT_IMAGE=my-image:1"

cat >"$tmp/system-flow.yaml" <<'EOF'
name: fixture
key: SF
agent:
  harness: claude-code
EOF
expect_image "flai-agent-claude-code:local"

# A run refuses an image docker does not have, before anything else is read.
cat >"$tmp/system-flow.yaml" <<'EOF'
name: fixture
key: SF
agent:
  image:
    layer: flai/images/x/Dockerfile
EOF
printf '#!/usr/bin/env sh\nexit 1\n' >"$tmp/bin/docker"
chmod +x "$tmp/bin/docker"
if got="$("$tmp/scripts/agent-container.sh" S-0001 2>&1)"; then
  fail "$got" "agent-container.sh ran with an image docker does not have"
fi
case "$got" in
  *"no image flai-agent-sf-claude-code:local"*scripts/agent-image.sh*) ;;
  *) fail "$got" "agent-container.sh did not refuse the image not built, naming scripts/agent-image.sh" ;;
esac

# With the image there, the run goes past the check and stops at the missing worktree.
printf '#!/usr/bin/env sh\nexit 0\n' >"$tmp/bin/docker"
if got="$("$tmp/scripts/agent-container.sh" S-0001 2>&1)"; then
  fail "$got" "agent-container.sh ran with no worktree"
fi
case "$got" in
  *"no worktree at "*) ;;
  *) fail "$got" "agent-container.sh stopped before the worktree check with the image built" ;;
esac
echo "agent-container-test: passed"
