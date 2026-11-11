#!/usr/bin/env sh
# Reproduces I-0348 on the host: go build's VCS stamping in a story's worktree
# records the main checkout's commit, since Go takes only a .git folder for a
# repository's root and a worktree's .git is a file, so it walks up past the
# worktree to the main checkout and stamps that HEAD. buildinfo.Get falls back
# to that stamp when -ldflags -X set no Commit, so a bin/flai that
# flai-build.sh builds in a story worktree names the main checkout's commit.
# With a scratch main checkout, a story worktree below it whose HEAD differs,
# and a small flai module whose buildinfo falls back as the real one does, the
# worktree's bin/flai must print the worktree's own commit and carry no
# vcs.revision of Go's.
#
# It needs go and git, and skips without them. It runs offline: the module has
# no dependencies and GOTOOLCHAIN=local keeps go from fetching a toolchain.
set -eu
. "$(dirname "$0")/env.sh"

skip() {
  echo "flai-build-stamp-test: skipped, $1"
  exit 0
}

command -v go >/dev/null 2>&1 || skip "go is not on PATH"
command -v git >/dev/null 2>&1 || skip "git is not on PATH"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
main="$tmp/main"
story=S-0001
worktree="$main/.flai-cache/worktrees/$story"

fail() {
  echo "flai-build-stamp-test: $1 (I-0348)"
  exit 1
}

# commit commits in the checkout given first, with the git commit options after it.
commit() {
  dir="$1"
  shift
  git -C "$dir" -c user.name=test -c user.email=test@example.invalid commit -q "$@"
}

mkdir -p "$main/scripts" "$main/flai/internal/buildinfo"
cp "$ROOT/scripts/env.sh" "$ROOT/scripts/flai-build.sh" "$main/scripts/"
printf 'bin/\n.flai-cache/\n' >"$main/.gitignore"
printf 'module github.com/bytepunx/system-flow/flai\n\ngo 1.21\n' >"$main/flai/go.mod"
cat >"$main/flai/internal/buildinfo/buildinfo.go" <<'EOF'
// Package buildinfo resolves the commit as the real buildinfo.Get does.
package buildinfo

import "runtime/debug"

// Commit and Date are set via -ldflags -X.
var Commit, Date string

// Get returns Commit and Date, falling back to Go's VCS stamp when empty.
func Get() (string, string) {
	commit, date := Commit, Date
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			switch {
			case s.Key == "vcs.revision" && commit == "":
				commit = s.Value
			case s.Key == "vcs.time" && date == "":
				date = s.Value
			}
		}
	}
	return commit, date
}
EOF
cat >"$main/flai/main.go" <<'EOF'
package main

import (
	"fmt"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

func main() {
	commit, _ := buildinfo.Get()
	fmt.Println(commit)
}
EOF
git -C "$main" init -q -b main
git -C "$main" add -A
commit "$main" -m init
git -C "$main" worktree add -q -b "story/$story" "$worktree"
commit "$worktree" --allow-empty -m "story commit"
main_head="$(git -C "$main" rev-parse HEAD)"
want="$(git -C "$worktree" rev-parse HEAD)"

# The builds stay local and offline, and an inherited GOFLAGS such as
# -buildvcs=false cannot hide the cause.
GOTOOLCHAIN=local GOFLAGS='' GOWORK=off GOPROXY=off "$worktree/scripts/flai-build.sh" >/dev/null ||
  fail "the worktree's scripts/flai-build.sh failed"
got="$("$worktree/bin/flai")" || fail "the worktree's bin/flai failed to run"

if [ "$got" != "$want" ]; then
  if [ "$got" = "$main_head" ]; then
    got="$got, the main checkout's HEAD"
  fi
  fail "the worktree's bin/flai printed commit ${got:-(none)}, not the worktree's own $want"
fi

stamp="$(go version -m "$worktree/bin/flai")"
case "$stamp" in
  *vcs.revision=*) fail "the worktree's bin/flai carries Go's vcs.revision stamp: $(printf '%s\n' "$stamp" | grep 'vcs.revision=')" ;;
esac
echo "flai-build-stamp-test: passed"
