#!/usr/bin/env sh
# Write the digest lists of a flai release's agent images, one per image, in the
# format of design/system/release-signing.md, "The digest list" (ADR-0295):
# <dir>/flai-agent-bedrock_<version>.digests and
# <dir>/flai-agent-claude-code_<version>.digests, one per image because the
# format holds one repository. Each list comes from flaiover-digests.sh, beside
# this script, for <registry>/flai-agent-<image>:<version> at the index digest
# given. Both lists are written to temporary files and moved in place only when
# both are complete, so a failure leaves neither. It prints the paths it wrote,
# one per line. release-flai.yml runs it on a flai/v* tag, then signs and
# uploads the lists to the flai vX.Y.Z release.
#
# Usage: scripts/agent-image-digests.sh <version> <bedrock-index-digest> <claude-code-index-digest> [<dir>]
#   e.g. scripts/agent-image-digests.sh 1.40.0 sha256:0dce07ff... sha256:b3c17aaf... dist
# <dir> must exist and defaults to the current directory. The registry prefix is
# AGENT_IMAGE_REGISTRY, ghcr.io/bytepunx when unset or empty. Needs what
# flaiover-digests.sh needs: docker with buildx (imagetools), jq, and sha256sum.
set -eu

die() { echo "agent-image-digests: $*" >&2; exit 1; }

usage="usage: $0 <version> <bedrock-index-digest> <claude-code-index-digest> [<dir>], e.g. $0 1.40.0 sha256:0dce07ff... sha256:b3c17aaf... dist"
[ "$#" -eq 3 ] || [ "$#" -eq 4 ] || die "$usage"
version=$1 bedrock=$2 claude_code=$3 dir=${4:-.}
registry=${AGENT_IMAGE_REGISTRY:-ghcr.io/bytepunx}
registry=${registry%/}

case $version in
  "" | v* | *[!0-9A-Za-z.+-]*) die "version '$version' must be the release version without its v, e.g. 1.40.0" ;;
esac
for digest in "$bedrock" "$claude_code"; do
  hex=${digest#sha256:}
  if [ "$hex" = "$digest" ] || [ "${#hex}" -ne 64 ] || [ -n "$(printf '%s' "$hex" | tr -d '0-9a-f')" ]; then
    die "digest '$digest' must be sha256: and 64 lower-case hex characters: give each image's index digest, as docker/build-push-action outputs it"
  fi
done
[ -d "$dir" ] || die "directory '$dir' does not exist: create it first, or leave it out to write to the current directory"
[ "$dir" = / ] || dir=${dir%/}

here=$(dirname "$0")
tmp_bedrock="" tmp_claude_code=""
trap 'rm -f -- ${tmp_bedrock:+"$tmp_bedrock"} ${tmp_claude_code:+"$tmp_claude_code"}' EXIT
tmp_bedrock=$(mktemp "$dir/.flai-agent-bedrock_$version.XXXXXX") || die "cannot write to '$dir': make it writable"
tmp_claude_code=$(mktemp "$dir/.flai-agent-claude-code_$version.XXXXXX") || die "cannot write to '$dir': make it writable"

# list writes the list of flai-agent-$1 at index digest $2 to the file $3.
list() {
  "$here/flaiover-digests.sh" "$registry/flai-agent-$1" "$version" "$2" >"$3" \
    || die "cannot write the digest list of $registry/flai-agent-$1:$version at $2: see the error above"
}

list bedrock "$bedrock" "$tmp_bedrock"
list claude-code "$claude_code" "$tmp_claude_code"

out_bedrock="$dir/flai-agent-bedrock_$version.digests"
out_claude_code="$dir/flai-agent-claude-code_$version.digests"
mv "$tmp_bedrock" "$out_bedrock" || die "cannot move the digest list in place as $out_bedrock"
mv "$tmp_claude_code" "$out_claude_code" || die "cannot move the digest list in place as $out_claude_code"
echo "$out_bedrock"
echo "$out_claude_code"
