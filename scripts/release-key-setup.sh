#!/usr/bin/env bash
# Runs the release signing key runbook (docs/operators/runbooks/release-key.md on
# story/S-0232) end to end, for TH-0324: generate the cosign pair, store it in the
# `release` environment, add the release-tag ruleset, print the public key and its
# fingerprint, destroy the local private key, and optionally reply on TH-0324.
#
# It follows the thread's recommendations 1 and 2 (conditional attestation step,
# environment secrets behind a tag rule).
#
# Usage: scripts/release-key-setup.sh [--yes] [--keep-key] [--own-password] [--no-reply]
#   --yes           do not pause before each section
#   --keep-key      keep cosign.key and show its password, for an offline recovery copy
#                   (the runbook recommends against it)
#   --own-password  type the password instead of generating a random one
#   --no-reply      do not offer to post the result on TH-0324
set -euo pipefail

REPO=bytepunx/system-flow
ENV_NAME=release
COSIGN_WANT=v3.1.3
THREAD=TH-0324
ROOT=$(cd "$(dirname "$0")/.." && pwd)

yes=0 keep_key=0 own_password=0 reply=1
for arg in "$@"; do
  case $arg in
    --yes|-y) yes=1 ;;
    --keep-key) keep_key=1 ;;
    --own-password) own_password=1 ;;
    --no-reply) reply=0 ;;
    -h|--help) sed -n '2,15p' "$0"; exit 0 ;;
    *) echo "unknown argument: $arg" >&2; exit 2 ;;
  esac
done

WORK=""
say()  { printf '\n\033[1m== %s\033[0m\n' "$*"; }
ok()   { printf '   \033[32mok\033[0m %s\n' "$*"; }
warn() { printf '   \033[33mwarn\033[0m %s\n' "$*" >&2; }
die()  { printf '   \033[31merror\033[0m %s\n' "$*" >&2; [ -n "$WORK" ] && printf '   work directory: %s\n' "$WORK" >&2; exit 1; }
pause() {
  [ "$yes" = 1 ] && return 0
  read -r -p "   Continue with: $1? [Y/n] " a
  case $a in n|N|no) echo "   stopped. Work directory: $WORK"; exit 1 ;; esac
}

# ---------------------------------------------------------------- preflight
say "Before you start"
command -v cosign  >/dev/null || die "cosign not found. go install github.com/sigstore/cosign/v3/cmd/cosign@$COSIGN_WANT"
command -v gh      >/dev/null || die "gh not found"
command -v openssl >/dev/null || die "openssl not found"
have=$(cosign version 2>/dev/null | awk '/GitVersion/ {print $2}')
if [ "$have" = "$COSIGN_WANT" ]; then ok "cosign $have"; else warn "cosign is ${have:-unknown}; the runbook was checked with $COSIGN_WANT"; fi
gh auth status >/dev/null 2>&1 || die "gh is not logged in (gh auth login)"
me=$(gh api user --jq .login)
[ "$(gh api "repos/$REPO" --jq .permissions.admin)" = true ] || die "$me is not an admin of $REPO"
ok "gh logged in as $me, admin of $REPO"
ok "openssl $(openssl version | awk '{print $2}')"

WORK=$(mktemp -d)
chmod 700 "$WORK"
cd "$WORK"
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then die "$WORK is inside a git repository"; fi
ok "working in $WORK"
trap 'unset COSIGN_PASSWORD' EXIT

# ---------------------------------------------------------------- generate
say "Generate the pair"
pause "generate a cosign key pair"
if [ "$own_password" = 1 ]; then
  read -r -s -p "   Password: " p1; echo
  read -r -s -p "   Again:    " p2; echo
  [ "$p1" = "$p2" ] || die "the passwords differ"
  [ ${#p1} -ge 20 ] || warn "a password under 20 characters is weak for this"
  COSIGN_PASSWORD=$p1; unset p1 p2
else
  COSIGN_PASSWORD=$(openssl rand -base64 36 | tr -d '\n')
  ok "generated a random password; it goes only into the Actions secret"
fi
export COSIGN_PASSWORD
cosign generate-key-pair >/dev/null
[ "$(head -1 cosign.key)" = "-----BEGIN ENCRYPTED SIGSTORE PRIVATE KEY-----" ] || die "cosign.key has an unexpected header"
[ "$(head -1 cosign.pub)" = "-----BEGIN PUBLIC KEY-----" ] || die "cosign.pub has an unexpected header"
ok "cosign.key and cosign.pub written"

cosign sign-blob --key cosign.key \
  --use-signing-config=false --new-bundle-format=false --tlog-upload=false \
  --output-signature test.sig --yes cosign.pub 2>/dev/null
cosign verify-blob --key cosign.pub --signature test.sig --insecure-ignore-tlog cosign.pub >/dev/null 2>&1 \
  || die "the test signature did not verify"
rm test.sig
ok "the password opens the key; a test signature verifies"

# ---------------------------------------------------------------- environment
say "Store the secrets"
pause "create environment '$ENV_NAME' limited to flai/v* and flaiover/v* tags, and set the two secrets"
gh api -X PUT "repos/$REPO/environments/$ENV_NAME" \
  -F 'deployment_branch_policy[protected_branches]=false' \
  -F 'deployment_branch_policy[custom_branch_policies]=true' >/dev/null
ok "environment $ENV_NAME exists with custom deployment rules"

list_policies() {
  gh api "repos/$REPO/environments/$ENV_NAME/deployment-branch-policies" \
    --jq '.branch_policies[] | "\(.type) \(.name)"'
}
existing=$(list_policies)
for pat in 'flai/v*' 'flaiover/v*'; do
  if grep -qxF "tag $pat" <<<"$existing"; then
    ok "tag rule $pat already present"
  else
    gh api -X POST "repos/$REPO/environments/$ENV_NAME/deployment-branch-policies" \
      -f name="$pat" -f type=tag >/dev/null
    ok "added tag rule $pat"
  fi
done
policies=$(list_policies | LC_ALL=C sort)
# shellcheck disable=SC2001 # indents every line; ${var//} cannot anchor at each line's start
sed 's/^/     /' <<<"$policies"
[ "$policies" = "$(printf 'tag flai/v*\ntag flaiover/v*')" ] \
  || die "the environment admits more than the two release tag patterns; remove the others in Settings > Environments > $ENV_NAME"

gh secret set COSIGN_PRIVATE_KEY --env "$ENV_NAME" -R "$REPO" < cosign.key
printf '%s' "$COSIGN_PASSWORD" | gh secret set COSIGN_PASSWORD --env "$ENV_NAME" -R "$REPO"
secrets=$(gh secret list --env "$ENV_NAME" -R "$REPO")
# shellcheck disable=SC2001 # indents every line; ${var//} cannot anchor at each line's start
sed 's/^/     /' <<<"$secrets"
for s in COSIGN_PRIVATE_KEY COSIGN_PASSWORD; do
  awk '{print $1}' <<<"$secrets" | grep -qxF "$s" || die "secret $s is missing"
done
ok "both secrets set in $ENV_NAME"

# ---------------------------------------------------------------- ruleset
say "Limit who may create the release tags"
if gh api "repos/$REPO/rulesets" --jq '.[].name' | grep -qxF 'release tags'; then
  ok "ruleset 'release tags' already exists; left as it is"
else
  pause "add tag ruleset 'release tags' (only organisation owners may create, move, or delete flai/v* and flaiover/v*)"
  if gh api -X POST "repos/$REPO/rulesets" --input - >/dev/null <<'JSON'
{
  "name": "release tags",
  "target": "tag",
  "enforcement": "active",
  "bypass_actors": [
    {"actor_id": 1, "actor_type": "OrganizationAdmin", "bypass_mode": "always"}
  ],
  "conditions": {
    "ref_name": {"include": ["refs/tags/flai/v*", "refs/tags/flaiover/v*"], "exclude": []}
  },
  "rules": [{"type": "creation"}, {"type": "update"}, {"type": "deletion"}]
}
JSON
  then
    ok "ruleset 'release tags' created"
  else
    warn "gh api refused the ruleset. Add it on GitHub instead:"
    warn "Settings > Rules > Rulesets > New ruleset > New tag ruleset; patterns flai/v* and flaiover/v*;"
    warn "restrict creations, updates, deletions; bypass: organisation owners and whoever pushes release tags."
  fi
fi
role=$(gh api "orgs/${REPO%%/*}/memberships/$me" --jq .role 2>/dev/null || echo unknown)
if [ "$role" = admin ]; then
  ok "you are an organisation owner, so you can still push release tags"
else
  warn "your organisation role is '$role'; add yourself to the ruleset's bypass list or your release tags will be refused"
fi

# ---------------------------------------------------------------- public key
say "Record the public key"
FPR=$(openssl pkey -pubin -in cosign.pub -outform DER | openssl dgst -sha256 | awk '{print $NF}')
CREATED=$(date -u +%Y-%m-%d)
OUT="$HOME/cosign-$CREATED.pub"
cp cosign.pub "$OUT"
ok "public key copied to $OUT"

# ---------------------------------------------------------------- destroy
say "Destroy the local private key"
if [ "$keep_key" = 1 ]; then
  warn "kept $WORK/cosign.key. Its password is below. Move both to a password manager, then delete the file."
  printf '     %s\n' "$COSIGN_PASSWORD"
else
  if command -v shred >/dev/null; then shred -u cosign.key; else rm -P cosign.key 2>/dev/null || rm -f cosign.key; fi
  ok "cosign.key destroyed; the Actions secret is the only copy"
fi
unset COSIGN_PASSWORD

# ---------------------------------------------------------------- result
say "Result"
echo "Fingerprint (SHA-256 of the DER public key): $FPR"
echo "Created: $CREATED"
echo
cat "$OUT"

REPLY_TEXT=$(cat <<TXT
Ran the release-key runbook. I confirm recommendations 1 and 2: the attestation step skipped while the repository is private, and the \`release\` environment limited to \`flai/v*\` and \`flaiover/v*\` tags, with the tag ruleset.

1. cosign.pub:

\`\`\`
$(cat "$OUT")
\`\`\`

2. Fingerprint (SHA-256 of the DER public key): \`$FPR\`, created $CREATED.
TXT
)

if [ "$reply" = 1 ]; then
  echo
  read -r -p "   Post this public key and fingerprint as your reply on $THREAD? [y/N] " a
  if [[ $a =~ ^[yY] ]]; then
    FLAI=$(command -v flai || echo "$HOME/.flai/bin/flai")
    (cd "$ROOT" && "$FLAI" thread reply "$THREAD" "$REPLY_TEXT")
    ok "replied on $THREAD"
  else
    echo "   Not posted. Reply on $THREAD with:"; echo; echo "$REPLY_TEXT"
  fi
fi
cd "$ROOT"
rmdir "$WORK" 2>/dev/null || true
