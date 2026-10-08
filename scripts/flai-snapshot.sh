#!/usr/bin/env sh
# GoReleaser snapshot build of flai into flai/dist (no publish).
#
#   scripts/flai-snapshot.sh                  all six targets, version 0.0.0-<commit>,
#                                             archives flai_snapshot_<os>_<arch>
#   scripts/flai-snapshot.sh --local X.Y.Z    the host's target only (go env GOOS/GOARCH),
#                                             no before hooks, version X.Y.Z, archive
#                                             flai_X.Y.Z_<os>_<arch>.tar.gz, as a release
#                                             tagged flai/vX.Y.Z names it, for the smoke
#                                             test's local release server (S-0340)
#   ... --local X.Y.Z --dist DIR              the same, into DIR (relative to flai/)
#                                             instead of flai/dist, which it leaves alone
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"

usage() {
  echo "usage: scripts/flai-snapshot.sh [--local X.Y.Z [--dist DIR]]" >&2
  exit 2
}

if [ "$#" -eq 0 ]; then
  exec goreleaser release --snapshot --clean --skip=publish
fi
if [ "$#" -lt 2 ] || [ "$1" != "--local" ]; then
  usage
fi
version="$2"
shift 2
dist="dist"
if [ "$#" -eq 2 ] && [ "$1" = "--dist" ] && [ -n "$2" ]; then
  dist="$2"
elif [ "$#" -ne 0 ]; then
  usage
fi

not_bare() {
  echo "flai-snapshot: --local wants a bare X.Y.Z such as 9.9.9, not \"$version\"" >&2
  exit 2
}
# install.sh and flai self-upgrade take only a bare X.Y.Z (buildinfo.Bare):
# digits, two dots, no empty part, no leading zero.
case "$version" in
  *[!0-9.]* | .* | *. | *..* | 0[0-9]* | *.0[0-9]*) not_bare ;;
esac
[ "$(printf '%s' "$version" | tr -d '0-9')" = ".." ] || not_bare

goos="$(go env GOOS)"
goarch="$(go env GOARCH)"

# GoReleaser's release has no --single-target and does not template goos or
# goarch, so the host's target goes in through a copy of the configuration
# with the build's two target lines narrowed. GoReleaser resolves paths from
# the working directory, so the copy may live outside the tree.
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
config="$tmp/goreleaser-local.yaml"
sed -e "s/^    goos: .*/    goos: [$goos]/" \
  -e "s/^    goarch: .*/    goarch: [$goarch]/" \
  .goreleaser.yaml >"$config"
printf 'dist: "%s"\n' "$dist" >>"$config"
for line in "    goos: [$goos]" "    goarch: [$goarch]"; do
  if [ "$(grep -cxF "$line" "$config" || true)" != "1" ]; then
    echo "flai-snapshot: flai/.goreleaser.yaml no longer has one build with a goos and a goarch line to narrow to $goos/$goarch; update scripts/flai-snapshot.sh to match it" >&2
    exit 1
  fi
done

FLAI_LOCAL_VERSION="$version" goreleaser release --snapshot --clean \
  --skip=publish,before --config "$config"
