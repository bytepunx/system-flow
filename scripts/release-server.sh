#!/usr/bin/env sh
# Serve flai releases from GoReleaser dist folders on 127.0.0.1, answering the
# GitHub endpoints install.sh and flai self-upgrade read (the release listing,
# a release by tag, an asset download), for a calling script (S-0340):
#
#   URL=$("$ROOT/scripts/release-server.sh" --owner $$ [--repo local/flai] <dist-dir>...)
#   FLAI_API="$URL" FLAI_REPO=local/flai sh install.sh
#   FLAI_RELEASES_API="$URL" bin/flai self-upgrade --check --repo local/flai
#
# Each dist folder holds flai_X.Y.Z_<os>_<arch>.tar.gz archives and their
# checksums.txt, as scripts/flai-snapshot.sh --local X.Y.Z writes into
# flai/dist; each version is served as the release tagged flai/vX.Y.Z, newest
# first. It builds the server into bin/release-server, starts it in the
# background, waits until it answers, prints its base URL on stdout, and
# returns. The server logs each request to this script's stderr and stops
# within a second of the --owner process exiting, so a caller that dies under
# set -e leaves no server behind. The server is plain HTTP: point no
# http_proxy at it.
set -eu
. "$(dirname "$0")/env.sh"

usage() {
  echo "usage: scripts/release-server.sh --owner PID [--repo OWNER/NAME] DIST_DIR..." >&2
  exit 2
}

owner=""
repo="local/flai"
while [ "$#" -gt 0 ]; do
  case "$1" in
    --owner) [ "$#" -ge 2 ] || usage; owner="$2"; shift 2 ;;
    --repo) [ "$#" -ge 2 ] || usage; repo="$2"; shift 2 ;;
    --) shift; break ;;
    -*) usage ;;
    *) break ;;
  esac
done
[ "$#" -ge 1 ] || usage
case "$owner" in
  '' | *[!0-9]*)
    echo "release-server: --owner wants the PID of the calling shell (\$\$), whose exit stops the server" >&2
    exit 2
    ;;
esac
kill -0 "$owner" 2>/dev/null || { echo "release-server: --owner $owner names no running process" >&2; exit 2; }

# Each dist folder as -dir, with the arguments in their order.
n=$#
for dir in "$@"; do
  [ -d "$dir" ] || { echo "release-server: $dir is not a folder; give the dist folder scripts/flai-snapshot.sh --local X.Y.Z builds (flai/dist)" >&2; exit 2; }
  set -- "$@" -dir "$dir"
done
shift "$n"

(cd "$ROOT/flai" && go build -o "$ROOT/bin/release-server" ./internal/releaseserver/serve) >&2

# The base URL goes to a file, not to this script's stdout: the caller reads
# that through $(...), which would wait for the server to close it.
tmp="$(mktemp -d)"
pid=""
ready=""
cleanup() {
  if [ -z "$ready" ] && [ -n "$pid" ]; then
    kill "$pid" 2>/dev/null || true
  fi
  rm -rf "$tmp"
}
trap cleanup EXIT

"$ROOT/bin/release-server" -repo "$repo" -owner-pid "$owner" "$@" >"$tmp/url" </dev/null &
pid=$!

url=""
tries=0
while [ -z "$url" ]; do
  # A child that exited stays a zombie until waited for, and kill -0 still
  # finds it, so its state is read instead.
  case "$(ps -o stat= -p "$pid" 2>/dev/null || true)" in
    '' | Z*)
      echo "release-server: the server exited before it listened; its log is above" >&2
      exit 1
      ;;
  esac
  tries=$((tries + 1))
  if [ "$tries" -gt 100 ]; then
    echo "release-server: the server printed no URL within 10 s" >&2
    exit 1
  fi
  sleep 0.1
  url="$(head -n 1 "$tmp/url")"
done

curl -fsS --noproxy '*' -H "Accept: application/vnd.github+json" "$url/repos/$repo/releases?per_page=1" -o /dev/null \
  || { echo "release-server: $url does not answer the release listing" >&2; exit 1; }
ready=1
echo "$url"
