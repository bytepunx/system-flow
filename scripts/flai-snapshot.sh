#!/usr/bin/env sh
# GoReleaser snapshot build of flai into flai/dist (no publish). Unsigned:
# --snapshot does not skip signing, and only the release workflow has cosign
# and the release key.
set -eu
. "$(dirname "$0")/env.sh"

# scripts/flai-snapshot.sh --local X.Y.Z [--dist DIR] builds a release for the
# smoke test's local release server (S-0340): the host's target only (go env
# GOOS/GOARCH), no before hooks, unsigned, version X.Y.Z, and the archive
# flai_X.Y.Z_<os>_<arch>.tar.gz, as a release tagged flai/vX.Y.Z names it, with
# its checksums.txt, into DIR (relative to flai/) or else flai/dist. With no
# arguments it builds every target as a snapshot, below.
if [ "$#" -gt 0 ]; then
  usage() {
    echo "usage: scripts/flai-snapshot.sh [--local X.Y.Z [--dist DIR]]" >&2
    exit 2
  }
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

  cd "$ROOT/flai"
  goos="$(go env GOOS)"
  goarch="$(go env GOARCH)"

  # GoReleaser's release has no --single-target and does not template goos or
  # goarch, so the host's target goes in through a copy of the configuration
  # with the build's two target lines narrowed, and the dist folder with it.
  # GoReleaser resolves paths from the working directory, so the copy may live
  # outside the tree.
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
    --skip=publish,before,sign --config "$config"
  exit 0
fi
cd "$ROOT/flai"
goreleaser release --snapshot --clean --skip=publish,sign
