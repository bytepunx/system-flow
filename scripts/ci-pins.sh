#!/usr/bin/env sh
# Fails when what CI and the images build from can move under a tag (S-0630,
# I-0297). Given a repository root, default this script's parent folder, it
# prints one line per finding, <path>:<line>: <reason>, and exits 1 when there
# is any:
#   - a uses: in .github/workflows/*.yml whose ref is not a 40-character commit
#     SHA; a local ./ action is exempt
#   - a workflow with no top-level permissions:, or one whose top-level block
#     grants a write or is write-all: writes belong on the jobs that use them
#   - a FROM in flaiover/Dockerfile or flai/images/*/Dockerfile, or an image: in
#     flaiover/compose.yaml, that names latest, main-stable, no tag, or a tag of
#     fewer than three version parts, or an -alpine variant that names no Alpine
#     release; an ARG's default stands for ${NAME}, and a
#     :local image or an earlier stage is exempt
#   - a FROM in flaiover/Dockerfile without an @sha256: digest, since a signed
#     flaiover release builds from it
#   - scripts/litellm.sh's default LITELLM_IMAGE, checked as an image is
#   - more than one lockfile under flaiover/, or a package-lock.json at the root
# smoke.sh runs it, so CI's system-flow-check.yml and the close-out run it.
set -eu
ROOT=${1:-"$(cd "$(dirname "$0")/.." && pwd)"}
cd "$ROOT"
out=$(mktemp)
trap 'rm -f "$out"' EXIT

# The tag check, shared by every image reference: $1 is the reference with its
# ARGs expanded. It prints the reason the reference fails, or nothing.
image_awk='
function image_reason(ref,   name, tag, i, core, parts) {
  sub(/@sha256:[0-9a-f]+$/, "", ref)
  if (ref ~ /:local$/) return ""
  name = ref; tag = ""
  i = match(ref, /:[^:\/]+$/)
  if (i > 0) { name = substr(ref, 1, i - 1); tag = substr(ref, i + 1) }
  if (tag == "") return "image " ref " names no tag"
  if (tag == "latest" || tag == "main-stable") return "image " ref " names the moving tag " tag
  core = tag
  sub(/^v/, "", core)
  if (!match(core, /^[0-9]+(\.[0-9]+)*/)) return "image " ref " names no version"
  parts = split(substr(core, 1, RLENGTH), p, ".")
  if (parts < 3) return "image " ref " names " tag ", not an exact version"
  if (tag ~ /-alpine$/) return "image " ref " names alpine with no release, such as -alpine3.24"
  return ""
}'

# Workflows: the uses: refs and the top-level permissions: block.
for f in .github/workflows/*.yml .github/workflows/*.yaml; do
  [ -f "$f" ] || continue
  awk -v f="$f" '
    function report(line, msg) { print f ":" line ": " msg }
    /^[[:space:]]*(-[[:space:]]+)?uses:[[:space:]]*/ {
      ref = $0
      sub(/^[[:space:]]*(-[[:space:]]+)?uses:[[:space:]]*/, "", ref)
      sub(/[[:space:]]+#.*$/, "", ref)
      gsub(/["'\'']/, "", ref)
      sub(/[[:space:]]+$/, "", ref)
      if (ref ~ /^\.\//) next
      at = ref; sub(/^.*@/, "", at)
      if (ref !~ /@/ || length(at) != 40 || at !~ /^[0-9a-f]+$/) report(NR, "action " ref " is not pinned to a commit SHA")
      next
    }
    /^permissions:/ {
      seen = 1; inblock = 1; line = NR
      v = $0; sub(/^permissions:[[:space:]]*/, "", v); sub(/[[:space:]]*#.*$/, "", v)
      if (v == "write-all" || v ~ /write/) report(NR, "the top-level permissions grant " v)
      next
    }
    inblock && /^[^[:space:]#]/ { inblock = 0 }
    inblock && /:[[:space:]]*write/ {
      g = $0; sub(/^[[:space:]]+/, "", g); sub(/[[:space:]]*#.*$/, "", g)
      report(NR, "the top-level permissions grant " g "; grant it on the job that needs it")
    }
    END { if (!seen) report(1, "no top-level permissions: block; grant contents: read") }
  ' "$f" >>"$out"
done

# Dockerfiles: each FROM, with ARG defaults expanded.
for f in flaiover/Dockerfile flai/images/*/Dockerfile; do
  [ -f "$f" ] || continue
  awk -v f="$f" "$image_awk"'
    toupper($1) == "ARG" {
      a = $2
      if (index(a, "=") > 0) {
        n = substr(a, 1, index(a, "=") - 1); v = substr(a, index(a, "=") + 1)
        gsub(/["'\'']/, "", v)
        if (!(n in def)) def[n] = v
      }
      next
    }
    toupper($1) == "FROM" {
      i = 2
      while ($i ~ /^--/) i++
      ref = $i
      for (n in def) { gsub("\\$\\{" n "\\}", def[n], ref); gsub("\\$" n "([^A-Za-z0-9_]|$)", def[n], ref) }
      if ((i + 2) <= NF && toupper($(i + 1)) == "AS") stage[$(i + 2)] = 1
      if (ref in stage) next
      r = image_reason(ref)
      if (r != "") print f ":" NR ": " r
      else if (f == "flaiover/Dockerfile" && ref !~ /@sha256:[0-9a-f]+$/) print f ":" NR ": image " ref " has no @sha256: digest, and a signed release builds from this file"
    }
  ' "$f" >>"$out"
done

# The local stack's compose file.
if [ -f flaiover/compose.yaml ]; then
  awk "$image_awk"'
    /^[[:space:]]*image:[[:space:]]*/ {
      ref = $0; sub(/^[[:space:]]*image:[[:space:]]*/, "", ref); sub(/[[:space:]]*#.*$/, "", ref); gsub(/["'\'']/, "", ref)
      r = image_reason(ref)
      if (r != "") print "flaiover/compose.yaml:" NR ": " r
    }
  ' flaiover/compose.yaml >>"$out"
fi

# The LiteLLM proxy's default image.
if [ -f scripts/litellm.sh ]; then
  awk "$image_awk"'
    /^LITELLM_IMAGE=\$\{LITELLM_IMAGE:-/ {
      ref = $0; sub(/^LITELLM_IMAGE=\$\{LITELLM_IMAGE:-/, "", ref); sub(/\}.*$/, "", ref)
      r = image_reason(ref)
      if (r != "") print "scripts/litellm.sh:" NR ": " r
    }
  ' scripts/litellm.sh >>"$out"
fi

# Lockfiles: what git tracks, or what is on disk outside a git checkout.
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  tracked() { [ -n "$(git ls-files -- "$1")" ]; }
else
  tracked() { [ -f "$1" ]; }
fi
locks=""
for l in flaiover/pnpm-lock.yaml flaiover/package-lock.json flaiover/yarn.lock; do
  if tracked "$l"; then locks="$locks $l"; fi
done
if [ "$(echo "$locks" | wc -w)" -gt 1 ]; then
  echo "flaiover:1: more than one lockfile:$locks; flaiover/pnpm-lock.yaml is the one the build reads" >>"$out"
fi
if tracked package-lock.json; then
  echo "package-lock.json:1: a lockfile at the root, which no package.json there needs" >>"$out"
fi

if [ -s "$out" ]; then
  cat "$out"
  echo "ci-pins: $(wc -l <"$out" | tr -d ' ') findings: pin each to an exact version or commit (design/tech/ci.md)" >&2
  exit 1
fi
echo "ci-pins: every action, image, workflow permission, and lockfile is pinned"
