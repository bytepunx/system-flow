#!/usr/bin/env sh
# Print flaiover's digest list for a release on stdout, in the format of
# design/system/release-signing.md, "The digest list": the index line first, then
# one line per platform manifest of the index (attestation manifests, platform
# unknown/unknown, left out), sorted by kind in byte order. release-flaiover.yml
# writes it to flaiover_<version>.digests on a flaiover/v* tag and signs it.
#
# Usage: scripts/flaiover-digests.sh <image> <version> <index-digest>
#   e.g. scripts/flaiover-digests.sh ghcr.io/bytepunx/flaiover 0.42.3 sha256:0dce07ff...
# Needs docker with buildx (imagetools), jq, and sha256sum; a private image needs
# a docker login to its registry first.
set -eu

die() { echo "flaiover-digests: $*" >&2; exit 1; }

[ "$#" -eq 3 ] || die "usage: $0 <image> <version> <index-digest>, e.g. $0 ghcr.io/bytepunx/flaiover 0.42.3 sha256:0dce07ff..."
image=$1 version=$2 digest=$3

case $image in
  "" | *@* | */ | *[[:space:]]*) bad_image=1 ;;
  *) case ${image##*/} in *:*) bad_image=1 ;; *) bad_image=0 ;; esac ;;
esac
[ "$bad_image" -eq 0 ] || die "image '$image' must be a repository without a tag or digest, e.g. ghcr.io/bytepunx/flaiover"
case $version in
  "" | v* | *[!0-9A-Za-z.+-]*) die "version '$version' must be the release version without its v, e.g. 0.42.3" ;;
esac
hex=${digest#sha256:}
if [ "$hex" = "$digest" ] || [ "${#hex}" -ne 64 ] || [ -n "$(printf '%s' "$hex" | tr -d '0-9a-f')" ]; then
  die "digest '$digest' must be sha256: and 64 lower-case hex characters, the index digest docker/build-push-action outputs"
fi

tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
docker buildx imagetools inspect --raw "$image@$digest" >"$tmp" \
  || die "cannot fetch $image@$digest (pushed? logged in to its registry?)"
got=$(sha256sum "$tmp" | cut -d' ' -f1)
[ "$got" = "$hex" ] || die "the registry returned content with digest sha256:$got for $image@$digest"

media=$(jq -r '.mediaType // ""' "$tmp") || die "$image@$digest is not JSON"
case $media in
  application/vnd.oci.image.index.v1+json | application/vnd.docker.distribution.manifest.list.v2+json) ;;
  *) die "$image@$digest is '${media:-no media type}', not an image index: give the index digest, not a platform's manifest" ;;
esac

# One "<digest> <kind>" per platform manifest; an entry without a valid digest
# or os and architecture is an error rather than a line left out.
platforms=$(jq -r '
  [.manifests[]? | select((.platform.os // "") != "unknown" or (.platform.architecture // "") != "unknown")]
  | map(
      if (.digest | type) != "string" or (.digest | test("^sha256:[0-9a-f]{64}$") | not)
         or (.platform.os // "") == "" or (.platform.architecture // "") == ""
      then error("bad manifest entry: \(tojson)")
      else "\(.digest) \(.platform.os)/\(.platform.architecture)\(if (.platform.variant // "") != "" then "/\(.platform.variant)" else "" end)"
      end)
  | .[]' "$tmp") || die "cannot read the manifests of $image@$digest"
[ -n "$platforms" ] || die "$image@$digest lists no platform manifests"

ref="$image:$version"
printf '%s %s index\n' "$digest" "$ref"
printf '%s\n' "$platforms" | LC_ALL=C sort -k2,2 | while read -r d kind; do
  printf '%s %s %s\n' "$d" "$ref" "$kind"
done
