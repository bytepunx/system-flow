#!/usr/bin/env sh
# Install the latest published flai from GitHub with install.sh, then exercise
# the same path through `flai self-upgrade` from the source build. Needs network
# access to GitHub and a token: GITHUB_TOKEN, GH_TOKEN, or gh. No close-out runs
# it: the smoke tier installs from a release built from the tree instead
# (scripts/install-test.sh). CI runs it on main and on a schedule.
#
# When a step fails it prints the step and its output and exits non-zero. On a
# host it also records the failure in this checkout through flai issue: it bumps
# the open issue titled as ISSUE_TITLE below, or opens it, and leaves it
# uncommitted. In CI (CI set, as GitHub Actions sets it) it records nothing: the
# failed run is the record.
set -eu
. "$(dirname "$0")/env.sh"
ISSUE_TITLE="The latest published flai release does not install from GitHub"
DIR="$ROOT/.flai-cache/install-published-test"
SCRATCH_HOME="$ROOT/.flai-cache/install-published-test-home"
LOG="$ROOT/.flai-cache/install-published-test.log"
OUTSIDE=""
trap 'if [ -n "$OUTSIDE" ]; then rm -rf "$OUTSIDE"; fi' EXIT
rm -rf "$DIR" "$SCRATCH_HOME" "$LOG"
# A story worktree has no .flai-cache of its own until something makes it.
mkdir -p "$ROOT/.flai-cache"

say() { echo "install-published-test: $*"; }
warn() { echo "install-published-test: $*" >&2; }

in_ci() {
  case "${CI:-}" in
    "" | false | 0) return 1 ;;
    *) return 0 ;;
  esac
}

# open_issue_id prints the ID of the open issue titled ISSUE_TITLE, or nothing.
# flai issue list prints one open issue a line: ID, class, status, count, cost,
# then the title.
open_issue_id() {
  issues=$(cd "$ROOT" && "$ROOT/scripts/flai.sh" issue list) || return 1
  printf '%s\n' "$issues" | awk -v t="$ISSUE_TITLE" \
    '$1 ~ /^I-[0-9]+$/ { id = $1; sub(/^[^ ]+ +[^ ]+ +[^ ]+ +[^ ]+ +[^ ]+ +/, ""); if ($0 == t) { print id; exit } }'
}

# record_issue STEP STATUS bumps the open issue titled ISSUE_TITLE, or opens it,
# with an instance naming STEP and the tail of its output in LOG. It commits
# nothing.
record_issue() {
  tail_out=$(tr '\r' '\n' <"$LOG" | tr -cd '\11\12\40-\176' | tr '\t' ' ' \
    | sed -e 's/[[:space:]]*$//' -e "s/\`\`\`/'''/g" | grep -v '^$' | tail -n 20)
  [ -n "$tail_out" ] || tail_out="(no output)"
  note=$(printf 'Step "%s" of scripts/install-published-test.sh failed with exit status %s on %s. The last lines of its output:\n\n```text\n%s\n```' \
    "$1" "$2" "$(uname -n)" "$tail_out")
  id=$(open_issue_id) || return 1
  if [ -n "$id" ]; then
    (cd "$ROOT" && "$ROOT/scripts/flai.sh" issue bump "$id" --note "$note") || return 1
  else
    (cd "$ROOT" && "$ROOT/scripts/flai.sh" issue new "$ISSUE_TITLE" --class defect --note "$note") || return 1
    id=$(open_issue_id) || return 1
  fi
  for f in "$ROOT/design/issues/$id-"*.md; do
    warn "recorded the failure in $f and in design/issues/summary.md, uncommitted"
  done
}

# fail STEP STATUS prints the failing step and its output, records it on a
# host, and exits non-zero.
fail() {
  warn "FAILED: $1 (exit status $2); its output:"
  cat "$LOG" >&2
  if in_ci; then
    warn "in CI (CI=${CI:-}): no issue recorded; the failed run is the record"
  else
    record_issue "$1" "$2" || warn "could not record the failure through flai issue"
  fi
  exit 1
}

# step NAME CMD [ARG...] runs the command with its output in LOG, prints the
# output when it passes, and fails with NAME when it does not. A step's
# function runs where set -e does not apply, so each returns its own failure.
step() {
  name=$1
  shift
  say "$name"
  status=0
  "$@" >"$LOG" 2>&1 || status=$?
  if [ "$status" -ne 0 ]; then
    fail "$name" "$status"
  fi
  cat "$LOG"
}

# These are preconditions, not the published release's failures: they exit
# non-zero and record nothing.
if [ -z "${GITHUB_TOKEN:-}${GH_TOKEN:-}" ] && ! command -v gh >/dev/null 2>&1; then
  warn "no GITHUB_TOKEN, GH_TOKEN, or gh; cannot reach the private repository"
  exit 1
fi
# gh reads its own auth from $HOME, so the token is resolved now, before HOME
# is redirected for the default-path install, and passed through explicitly.
TOKEN="${GITHUB_TOKEN:-${GH_TOKEN:-}}"
[ -n "$TOKEN" ] || TOKEN=$(gh auth token 2>/dev/null || true)
"$ROOT/scripts/flai-build.sh" >/dev/null

install_explicit_dir() {
  FLAI_INSTALL_DIR="$DIR" sh "$ROOT/install.sh" || return 1
  "$DIR/flai" version | grep '^flai [0-9]' \
    || { echo "the installed binary does not report a version"; return 1; }
}

install_default_dir() {
  env -i HOME="$SCRATCH_HOME" PATH="$PATH" GITHUB_TOKEN="$TOKEN" \
    FLAI_API="${FLAI_API:-}" FLAI_REPO="${FLAI_REPO:-}" sh "$ROOT/install.sh" >"$LOG.install" 2>&1
  install_status=$?
  cat "$LOG.install"
  [ "$install_status" -eq 0 ] || return "$install_status"
  [ -x "$SCRATCH_HOME/.flai/bin/flai" ] \
    || { echo "the default install did not land at \$HOME/.flai/bin"; return 1; }
  "$SCRATCH_HOME/.flai/bin/flai" version | grep '^flai [0-9]' \
    || { echo "the default-path binary does not report a version"; return 1; }
  grep -q 'export PATH=' "$LOG.install" \
    || { echo "install.sh did not print a PATH line"; return 1; }
}

upgrade_check() {
  "$ROOT/bin/flai" self-upgrade --check
}

upgrade_dir() {
  "$ROOT/bin/flai" self-upgrade --dir "$DIR" || return 1
  "$DIR/flai" version | head -1
}

# $DIR is inside this project, where self-upgrade never installs (S-0111), so
# the binary's own path is checked on a copy outside any project. Resolved, as
# flai resolves its own path: macOS's temporary directory is under /var, a link
# to /private/var (I-0045).
upgrade_own_path() {
  cp "$DIR/flai" "$OUTSIDE/flai" || return 1
  out=$("$OUTSIDE/flai" self-upgrade --check 2>&1)
  check_status=$?
  printf '%s\n' "$out"
  [ "$check_status" -eq 0 ] || return "$check_status"
  printf '%s\n' "$out" | grep -q 'installed at '"$OUTSIDE"'/flai' \
    || { echo "self-upgrade with no --dir did not resolve to the installed binary's own path"; return 1; }
}

upgrade_in_project() {
  out=$(env -i HOME="$SCRATCH_HOME" PATH="$PATH" GITHUB_TOKEN="$TOKEN" \
    FLAI_RELEASES_API="${FLAI_RELEASES_API:-}" "$DIR/flai" self-upgrade --check 2>&1)
  check_status=$?
  printf '%s\n' "$out"
  [ "$check_status" -eq 0 ] || return "$check_status"
  printf '%s\n' "$out" | grep -q 'installed at '"$SCRATCH_HOME"'/.flai/bin/flai' \
    || { echo "self-upgrade inside a project did not resolve to \$HOME/.flai/bin"; return 1; }
}

step "install.sh with an explicit FLAI_INSTALL_DIR" install_explicit_dir
step "install.sh with no FLAI_INSTALL_DIR installs under a fresh HOME/.flai/bin, without sudo" install_default_dir
step "flai self-upgrade --check" upgrade_check
step "flai self-upgrade --dir" upgrade_dir
OUTSIDE=$(cd "$(mktemp -d)" && pwd -P)
step "flai self-upgrade with no --dir, outside a project, resolves to the binary's own path" upgrade_own_path
step "flai self-upgrade with no --dir, run from inside a project, resolves to HOME/.flai/bin" upgrade_in_project
rm -rf "$SCRATCH_HOME" "$LOG" "$LOG.install"
say "passed"
