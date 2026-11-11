#!/usr/bin/env sh
# Build bin/flai from flai/, statically (CGO_ENABLED=0), so the one binary in a story worktree runs
# both on the host and in its agent container, whose Alpine has no glibc (S-0510).
# It builds to a temporary file beside bin/flai and renames it into place, so tiers that run at
# once (ADR-0266), each of which may build bin/flai, never run a half-written binary (S-0610).
# It links the built tree's own commit and commit time into buildinfo.Commit and Date, and builds
# with -buildvcs=false: Go's VCS stamp takes only a .git folder for a repository's root, so in a
# story's worktree, whose .git is a file, it would name the main checkout's HEAD (I-0348). git
# resolves the worktree's .git file, so these are the story branch's. Where git cannot answer, the
# two are left unset and flai version prints unknown for them.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "building bin/flai"
mkdir -p "$ROOT/bin"
tmp="$ROOT/bin/.flai.$$"
trap 'rm -f "$tmp"' EXIT
pkg=github.com/bytepunx/system-flow/flai/internal/buildinfo
ldflags=
if commit="$(git -C "$ROOT" rev-parse HEAD 2>/dev/null)" &&
  date="$(TZ=UTC git -C "$ROOT" show -s --format=%cd --date=format-local:%Y-%m-%dT%H:%M:%SZ HEAD 2>/dev/null)"; then
  ldflags="-X $pkg.Commit=$commit -X $pkg.Date=$date"
fi
CGO_ENABLED=0 go build -buildvcs=false -ldflags "$ldflags" -o "$tmp" .
mv -f "$tmp" "$ROOT/bin/flai"
