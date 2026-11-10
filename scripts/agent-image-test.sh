#!/usr/bin/env sh
# Reproduces I-0339 in a scratch copy of this tree's env.sh and agent-image.sh,
# with no Docker: when the credential helper the Docker configuration names in
# credsStore fails, as Docker Desktop's desktop.exe does from WSL when Windows
# interop is down, or is not on PATH, agent-image.sh builds with DOCKER_CONFIG
# at a scratch configuration with no login, keeping the current context and the
# folder's other entries, says so, leaves the operator's config.json as it was,
# and removes the scratch configuration. A working helper, or none named,
# leaves DOCKER_CONFIG as it is.
set -eu
. "$(dirname "$0")/env.sh"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
repo="$tmp/repo"
mkdir -p "$repo/scripts" "$repo/bin" "$tmp/home/.docker/cli-plugins" "$tmp/home/.docker/contexts"
cp "$ROOT/scripts/env.sh" "$ROOT/scripts/agent-image.sh" "$repo/scripts/"
cat >"$repo/system-flow.yaml" <<'EOF'
name: fixture
key: SF
agent:
  harness: claude-code
EOF
unset DOCKER_CONFIG
HOME="$tmp/home"
LOG="$tmp/docker.log"
export HOME LOG

# The fake docker records, for each call, its command, the DOCKER_CONFIG it
# sees, and that configuration's config.json and entries.
cat >"$repo/bin/docker" <<'EOF'
#!/usr/bin/env sh
{
  echo "call $1 DOCKER_CONFIG=${DOCKER_CONFIG:-}"
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
  echo "agent-image-test: $2 (I-0339)"
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
  out="$("$repo/scripts/agent-image.sh" 2>&1)" || fail "$out" "agent-image.sh failed: $1"
}

# expect_anonymous says both builds ran with a scratch configuration with no
# login, the current context kept and the folder's other entries linked, that
# the operator's config.json is unchanged, and that the scratch one is gone.
expect_anonymous() {
  [ "$(grep -c '^call build ' "$LOG")" -eq 2 ] || fail "$out" "$1: agent-image.sh did not run docker build twice"
  scratch="$(sed -n 's/^call build DOCKER_CONFIG=//p' "$LOG" | sort -u)"
  if [ -z "$scratch" ] || [ "$(printf '%s\n' "$scratch" | wc -l)" -ne 1 ]; then
    fail "$out" "$1: the builds did not run with one DOCKER_CONFIG"
  fi
  [ "$scratch" != "$HOME/.docker" ] || fail "$out" "$1: the builds ran with the operator's configuration"
  [ "$(grep -c '^config {"currentContext": "desktop-linux"}$' "$LOG")" -eq 2 ] || fail "$out" "$1: the scratch config.json is not the current context alone, with no login"
  grep -q '^entry cli-plugins$' "$LOG" || fail "$out" "$1: the scratch configuration lacks cli-plugins"
  grep -q '^entry contexts$' "$LOG" || fail "$out" "$1: the scratch configuration lacks contexts"
  [ ! -e "$scratch" ] || fail "$out" "$1: the scratch configuration $scratch was left behind"
  cmp -s "$HOME/.docker/config.json" "$tmp/config.json" || fail "$out" "$1: the operator's config.json changed"
  case "$out" in
    *"$2"*"building with no login"*"is unchanged"*) ;;
    *) fail "$out" "$1: agent-image.sh did not say it built with no login, naming $2" ;;
  esac
}

# expect_configured says both builds ran with DOCKER_CONFIG unset, and nothing
# was said of the login.
expect_configured() {
  [ "$(grep -c '^call build DOCKER_CONFIG=$' "$LOG")" -eq 2 ] || fail "$out" "$1: the builds did not run with the operator's configuration"
  case "$out" in
    *"no login"*) fail "$out" "$1: agent-image.sh said it built with no login" ;;
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
echo "agent-image-test: passed"
