---
title: "Runbook: the release signing key"
updated: 2026-10-12
status: active
---

# The release signing key

Step by step, for the operator who holds the key. Most operators never run this: they only compare a fingerprint, under [Current key](#current-key).

## What the key is for

One cosign key pair, ECDSA P-256, signs every release. In CI it signs flai's `checksums.txt`, and later flaiover's digest list. The public key is built into flai and flaiover, which check a release with it before they install or run it. The private key lives only in GitHub Actions, as two secrets of `bytepunx/system-flow`. Why this way, and what it does not promise, is in [release-signing.md](../../../design/system/release-signing.md) and [ADR-0070](../../../design/adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md).

## Before you start

- cosign v3.1.3, from its [release page](https://github.com/sigstore/cosign/releases/tag/v3.1.3) or with Go:

  ```bash
  go install github.com/sigstore/cosign/v3/cmd/cosign@v3.1.3
  cosign version
  ```

- `gh`, logged in as an admin of `bytepunx/system-flow` (`gh auth status`).
- `openssl`.
- A trusted machine. Work in a fresh, empty directory outside any repository, so that the private key can never be committed:

  ```bash
  cd "$(mktemp -d)"
  ```

## Generate the pair

1. Make the pair:

   ```bash
   cosign generate-key-pair
   ```

   It asks for a password twice. Choose a long random one, from a password manager. It reads `COSIGN_PASSWORD` instead when that is set.
2. Check what it wrote:

   ```bash
   head -1 cosign.key   # -----BEGIN ENCRYPTED SIGSTORE PRIVATE KEY-----
   head -1 cosign.pub   # -----BEGIN PUBLIC KEY-----
   ```

   `cosign.key` is the private key, encrypted with the password. `cosign.pub` is the public key.
3. Check that the password opens the key, by signing and verifying a file the way CI does:

   ```bash
   cosign sign-blob --key cosign.key \
     --use-signing-config=false --new-bundle-format=false --tlog-upload=false \
     --output-signature test.sig --yes cosign.pub
   cosign verify-blob --key cosign.pub --signature test.sig --insecure-ignore-tlog cosign.pub
   rm test.sig
   ```

   cosign warns that two of the flags are deprecated. That is expected; see [release-signing.md](../../../design/system/release-signing.md#cosign-with-a-key-pair).

Do not use `cosign generate-key-pair github://…`. It writes repository secrets, which any workflow run can read, a branch's included.

## Store the secrets

The two secrets go in an environment, `release`, that only the release tags can deploy to. Create it before you set the secrets: a workflow that names an environment that does not exist creates it with no rules.

1. Create the environment, with deployment limited to the rules you add next:

   ```bash
   gh api -X PUT repos/bytepunx/system-flow/environments/release \
     -F 'deployment_branch_policy[protected_branches]=false' \
     -F 'deployment_branch_policy[custom_branch_policies]=true'
   ```

2. Admit only the release tags:

   ```bash
   gh api -X POST repos/bytepunx/system-flow/environments/release/deployment-branch-policies \
     -f name='flai/v*' -f type=tag
   gh api -X POST repos/bytepunx/system-flow/environments/release/deployment-branch-policies \
     -f name='flaiover/v*' -f type=tag
   ```

3. Check it. The list names the two tag patterns and nothing else:

   ```bash
   gh api repos/bytepunx/system-flow/environments/release/deployment-branch-policies \
     --jq '.branch_policies[] | "\(.type) \(.name)"'
   ```

4. Set the secrets:

   ```bash
   gh secret set COSIGN_PRIVATE_KEY --env release -R bytepunx/system-flow < cosign.key
   gh secret set COSIGN_PASSWORD --env release -R bytepunx/system-flow   # asks for the password
   gh secret list --env release -R bytepunx/system-flow
   ```

5. Limit who may create the release tags, with a tag ruleset. Without it, anyone with write access can push a `flai/v*` tag and run the release workflow with the key. This one lets only organisation owners create, move, or delete them:

   ```bash
   gh api -X POST repos/bytepunx/system-flow/rulesets --input - <<'EOF'
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
   EOF
   ```

   Or, on GitHub: the repository's **Settings**, **Rules**, **Rulesets**, **New ruleset**, **New tag ruleset**, with the same two patterns, the three restrictions, and repository admins or organisation owners in the bypass list. Whoever pushes release tags, including an agent that publishes with your credentials, must be in the bypass list.

Why both: `bytepunx/system-flow` is a private repository on GitHub Team, where an environment cannot have required reviewers. The environment's tag rule keeps a branch's workflow from reading the secrets, and the ruleset keeps the release tags to the people you trust.

## Record and publish the public key

1. Print the fingerprint, the SHA-256 of the public key in DER form:

   ```bash
   openssl pkey -pubin -in cosign.pub -outform DER | openssl dgst -sha256
   ```

   `openssl ec -pubin -in cosign.pub -outform DER | openssl dgst -sha256` gives the same.
2. Add a row to [Current key](#current-key) with the fingerprint and today's date.
3. Put `cosign.pub` in the two sources, or give it to the agent working the story: copy the file to `flai/internal/buildinfo/cosign.pub`, which flai embeds, and paste the same PEM into `RELEASE_PUBLIC_KEY` in `flaiover/src/lib/server/release.ts`. A flaiover test fails when the two differ, and a test on each side fails when the fingerprint is not the one under [Current key](#current-key). flai's release publishes the file as the asset `cosign.pub`.

`cosign.pub` is not secret. Keep a copy wherever is convenient.

## Current key

An operator compares this fingerprint with the one a release's notes give, before installing it.

| Fingerprint (SHA-256 of the DER public key) | Created | Replaces |
|---------------------------------------------|---------|----------|
| `ee6df72313c9dbac1cf2270ed7f01cc523786d37bfce2dec12210215d6c5da40` | 2026-10-12 | none, the first key |

## Destroy the local private key

Once the secrets are set and the check above passed, delete `cosign.key`. The Actions secret is then the only copy.

```bash
shred -u cosign.key   # Linux
rm -P cosign.key      # macOS
```

Plain `rm` will do where neither exists. If you want a recovery copy, keep the file and its password offline, in a password manager. We recommend not keeping one: a lost key is rotated, not recovered, and a copy is one more place it can leak from.

## Rotate the key

1. Generate a new pair, as in [Generate the pair](#generate-the-pair), in a new empty directory.
2. Release flai, and flaiover, carrying both public keys, the old and the new. This bridging release is still signed with the old key.
3. Only then replace the two secrets with the new key and password, as in step 4 of [Store the secrets](#store-the-secrets). `gh secret set` overwrites a secret of the same name.
4. Add a row to [Current key](#current-key), naming the old fingerprint under Replaces.
5. Destroy the new local private key.

A flai older than the bridging release cannot verify a release signed with the new key. `flai self-upgrade` says so and names the bridging release to install first.

## Replace a compromised key

1. If the key may be in someone else's hands, delete the secrets at once, so that no release is signed until you are ready:

   ```bash
   gh secret delete COSIGN_PRIVATE_KEY --env release -R bytepunx/system-flow
   gh secret delete COSIGN_PASSWORD --env release -R bytepunx/system-flow
   ```

2. Rotate, as above, without waiting. An older flai accepts the bridging release only if the old key signs it. If you deleted the secrets and keep no recovery copy, skip the bridging release: set the new key's secrets and release with them. An older flai then refuses every new release, and operators install the new one by hand after comparing its fingerprint, as in [Check a release by hand](#check-a-release-by-hand).
3. Announce the compromise in `CHANGELOG.md` and in the new release's notes, with the new key's fingerprint.

A bridging release is signed with the compromised key, which an attacker also holds. Its signature therefore proves nothing; the fingerprint published beside it is what an operator checks before installing. There is no offline root key behind the release key. That is deliberate for a project of this size, and it means a compromise is recovered only by people comparing fingerprints.

## Check a release by hand

Download `checksums.txt` and `checksums.txt.sig` from the release, and `cosign.pub` from the same release or from your own copy. A key from your own copy is the stronger check: a key downloaded beside the release proves only what the release says about itself. Compare its fingerprint with [Current key](#current-key) first.

With cosign:

```bash
cosign verify-blob --key cosign.pub --signature checksums.txt.sig --insecure-ignore-tlog checksums.txt
```

With openssl alone:

```bash
openssl base64 -d -A -in checksums.txt.sig -out checksums.txt.der
openssl dgst -sha256 -verify cosign.pub -signature checksums.txt.der checksums.txt   # Verified OK
```

Then check the archive against `checksums.txt`:

```bash
sha256sum --ignore-missing -c checksums.txt   # Linux
shasum -a 256 --ignore-missing -c checksums.txt   # macOS
```
