#!/usr/bin/env sh
# Checks agent-image-digests.sh with no network and no registry (ADR-0295): a
# fake docker first on PATH answers `docker buildx imagetools inspect --raw
# <image>@<digest>` with a fixture image index chosen by the digest, each an
# OCI index of linux/amd64, linux/arm64 (with variant v8 in one), and an
# attestation manifest of platform unknown/unknown, in an unsorted order. The
# script must write flai-agent-bedrock_<version>.digests and
# flai-agent-claude-code_<version>.digests in the format of
# design/system/release-signing.md, "The digest list", byte for byte, print
# their paths, take the registry from AGENT_IMAGE_REGISTRY, and on a bad
# version, a bad digest, a missing directory, or a failed fetch exit non-zero
# and leave no list.
#
# It needs jq and sha256sum, as flaiover-digests.sh does, and skips without them.
set -eu
. "$(dirname "$0")/env.sh"

skip() {
  echo "agent-image-digests-test: skipped, $1"
  exit 0
}

command -v jq >/dev/null 2>&1 || skip "jq is not on PATH"
command -v sha256sum >/dev/null 2>&1 || skip "sha256sum is not on PATH"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
mkdir -p "$tmp/bin" "$tmp/fixtures"
LOG="$tmp/docker.log"
FIXTURES="$tmp/fixtures"
export LOG FIXTURES
PATH="$tmp/bin:$PATH"
unset AGENT_IMAGE_REGISTRY

# The fake docker logs the reference it is asked for and prints the fixture
# whose sha256 is the reference's digest, or fails as a registry would.
cat >"$tmp/bin/docker" <<'EOF'
#!/usr/bin/env sh
[ "$*" = "buildx imagetools inspect --raw $5" ] || { echo "fake docker: unexpected call: $*" >&2; exit 2; }
echo "$5" >>"$LOG"
hex=${5##*@sha256:}
[ -f "$FIXTURES/$hex" ] || { echo "fake docker: $5: not found" >&2; exit 1; }
cat "$FIXTURES/$hex"
EOF
chmod +x "$tmp/bin/docker"

checks=0
fail() {
  [ -z "${out:-}" ] || printf 'stdout:\n%s\n' "$out"
  [ ! -s "$tmp/err" ] || printf 'stderr:\n%s\n' "$(cat "$tmp/err")"
  echo "agent-image-digests-test: FAIL: $1"
  exit 1
}
ok() {
  checks=$((checks + 1))
  echo "agent-image-digests-test: ok: $1"
}

# h prints 64 lower-case hex characters made from the character given.
h() { printf "%064d" 0 | tr 0 "$1"; }

# fixture writes an OCI image index of the manifest entries given, joined by
# commas, and sets $index to its digest.
fixture() {
  printf '{"schemaVersion":2,"mediaType":"application/vnd.oci.image.index.v1+json","manifests":[%s]}' "$1" >"$tmp/index.json"
  index="sha256:$(sha256sum "$tmp/index.json" | cut -d' ' -f1)"
  mv "$tmp/index.json" "$FIXTURES/${index#sha256:}"
}
entry() {
  printf '{"mediaType":"application/vnd.oci.image.manifest.v1+json","digest":"sha256:%s","size":1234,"platform":{%s}}' "$1" "$2"
}
attestation() {
  printf '{"mediaType":"application/vnd.oci.image.manifest.v1+json","digest":"sha256:%s","size":567,"annotations":{"vnd.docker.reference.digest":"sha256:%s","vnd.docker.reference.type":"attestation-manifest"},"platform":{"architecture":"unknown","os":"unknown"}}' "$1" "$2"
}

fixture "$(entry "$(h a)" '"architecture":"arm64","os":"linux"'),$(attestation "$(h c)" "$(h a)"),$(entry "$(h b)" '"architecture":"amd64","os":"linux"')"
bedrock=$index
fixture "$(attestation "$(h f)" "$(h d)"),$(entry "$(h d)" '"architecture":"arm64","os":"linux","variant":"v8"'),$(entry "$(h e)" '"architecture":"amd64","os":"linux"')"
claude_code=$index
version=1.40.0

# expected prints the list the fixtures must give for the image $1, under the
# registry $2.
expected() {
  ref="$2/flai-agent-$1:$version"
  case $1 in
    bedrock) printf '%s %s index\nsha256:%s %s linux/amd64\nsha256:%s %s linux/arm64\n' "$bedrock" "$ref" "$(h b)" "$ref" "$(h a)" "$ref" ;;
    claude-code) printf '%s %s index\nsha256:%s %s linux/amd64\nsha256:%s %s linux/arm64/v8\n' "$claude_code" "$ref" "$(h e)" "$ref" "$(h d)" "$ref" ;;
  esac
}

# check_format checks the file $1, the list of $2 under the registry $3, against
# the format's rules, then against the expected bytes.
check_format() {
  [ -s "$1" ] || fail "$1 is missing or empty"
  [ "$(tail -c 1 "$1" | od -An -c | tr -d ' ')" = '\n' ] || fail "$1 has no final newline"
  ! grep -q "$(printf '\r')" "$1" || fail "$1 has a CR: line ends must be LF"
  ! grep -q '^$' "$1" || fail "$1 has a blank line"
  ref="$3/flai-agent-$2:$version"
  bad=$(grep -Ev "^sha256:[0-9a-f]{64} $(printf '%s' "$ref" | sed 's/[.]/[.]/g') (index|[a-z0-9]+/[a-z0-9]+(/[a-z0-9]+)?)$" "$1" || true)
  [ -z "$bad" ] || fail "$1 has a line not of digest, $ref, and kind: $bad"
  [ "$(sed -n '1s/.* //p' "$1")" = index ] || fail "$1 does not start with its index line"
  [ "$(sed -n '1s/ .*//p' "$1")" = "$(expected "$2" "$3" | sed -n '1s/ .*//p')" ] || fail "$1's index line does not name the index digest given"
  [ "$(sed 1d "$1" | cut -d' ' -f3)" = "$(sed 1d "$1" | cut -d' ' -f3 | LC_ALL=C sort)" ] || fail "$1's platform lines are not sorted by kind"
  ! grep -q unknown "$1" || fail "$1 lists the attestation manifest"
  expected "$2" "$3" >"$tmp/want"
  cmp -s "$tmp/want" "$1" || fail "$1 is not the expected list: want
$(cat "$tmp/want")
got
$(cat "$1")"
}

# run runs agent-image-digests.sh with the arguments given, keeping stdout in
# $out and stderr in $tmp/err, and sets $status.
run() {
  rm -f "$LOG"
  status=0
  out=$("$ROOT/scripts/agent-image-digests.sh" "$@" 2>"$tmp/err") || status=$?
}

# entries prints the names in the directory $1, hidden ones included, sorted,
# each followed by a space.
entries() {
  find "$1" -mindepth 1 -maxdepth 1 -exec basename {} \; | LC_ALL=C sort | tr '\n' ' '
}

# expect_refusal checks the last run failed, printed nothing on stdout, said why
# on stderr with the script's prefix, and left the directory $2 empty.
expect_refusal() {
  [ "$status" -ne 0 ] || fail "$1: agent-image-digests.sh exited 0"
  [ -z "$out" ] || fail "$1: agent-image-digests.sh printed paths"
  grep -q '^agent-image-digests: ' "$tmp/err" || fail "$1: no error prefixed agent-image-digests: on stderr"
  if [ -d "$2" ] && [ -n "$(entries "$2")" ]; then
    fail "$1: left $(entries "$2")in $2"
  fi
  ok "$1"
}

mkdir "$tmp/out"
run "$version" "$bedrock" "$claude_code" "$tmp/out"
[ "$status" -eq 0 ] || fail "agent-image-digests.sh failed with fixtures that are right"
ok "it succeeds against the fixture registry"
[ "$(entries "$tmp/out")" = "flai-agent-bedrock_$version.digests flai-agent-claude-code_$version.digests " ] \
  || fail "it wrote $(entries "$tmp/out"), not the two lists alone"
ok "it writes flai-agent-bedrock_$version.digests and flai-agent-claude-code_$version.digests and nothing else"
check_format "$tmp/out/flai-agent-bedrock_$version.digests" bedrock ghcr.io/bytepunx
ok "the bedrock list is in the digest list format, byte for byte"
check_format "$tmp/out/flai-agent-claude-code_$version.digests" claude-code ghcr.io/bytepunx
ok "the claude-code list is in the digest list format, with the variant, byte for byte"
[ "$out" = "$(printf '%s\n%s' "$tmp/out/flai-agent-bedrock_$version.digests" "$tmp/out/flai-agent-claude-code_$version.digests")" ] \
  || fail "stdout is not the two paths, one per line"
ok "stdout lists the paths it wrote"
[ "$(cat "$LOG")" = "$(printf 'ghcr.io/bytepunx/flai-agent-bedrock@%s\nghcr.io/bytepunx/flai-agent-claude-code@%s' "$bedrock" "$claude_code")" ] \
  || fail "it fetched $(cat "$LOG"), not each image at its index digest"
ok "it fetches each image from ghcr.io/bytepunx at its own index digest"

mkdir "$tmp/registry"
cd "$tmp/registry"
AGENT_IMAGE_REGISTRY=registry.example.test/team/ run "$version" "$bedrock" "$claude_code"
cd "$ROOT"
[ "$status" -eq 0 ] || fail "agent-image-digests.sh failed with AGENT_IMAGE_REGISTRY set"
check_format "$tmp/registry/flai-agent-bedrock_$version.digests" bedrock registry.example.test/team
check_format "$tmp/registry/flai-agent-claude-code_$version.digests" claude-code registry.example.test/team
[ "$out" = "$(printf './flai-agent-bedrock_%s.digests\n./flai-agent-claude-code_%s.digests' "$version" "$version")" ] \
  || fail "with no directory given, stdout is not the two paths in the current directory"
ok "AGENT_IMAGE_REGISTRY changes the reference, and the directory defaults to the current one"

mkdir "$tmp/v"
run "v$version" "$bedrock" "$claude_code" "$tmp/v"
expect_refusal "a version with a leading v is refused" "$tmp/v"

mkdir "$tmp/digest"
run "$version" "$bedrock" "sha256:$(h A)" "$tmp/digest"
expect_refusal "an upper-case digest is refused" "$tmp/digest"
run "$version" "${bedrock#sha256:}" "$claude_code" "$tmp/digest"
expect_refusal "a digest without sha256: is refused" "$tmp/digest"

run "$version" "$bedrock" "$claude_code" "$tmp/missing"
expect_refusal "a missing directory is refused" "$tmp/missing"
[ ! -e "$tmp/missing" ] || fail "a missing directory was made"

mkdir "$tmp/fetch"
run "$version" "$bedrock" "sha256:$(h 9)" "$tmp/fetch"
expect_refusal "a failed fetch of the second image leaves no list, for either image" "$tmp/fetch"
grep -q "flai-agent-claude-code" "$tmp/err" || fail "the fetch failure does not name flai-agent-claude-code"
ok "the fetch failure names the image"

echo "agent-image-digests-test: passed, $checks checks"
