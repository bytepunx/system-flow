---
id: TH-0065
title: "Signing and peer verification for E-0015: the recommendation and four questions"
anchor:
  path: wip/kanban/stories/S-0193-explore-ways-to-sign-and-verify-components.md
  item: S-0193
status: resolved
participants: [agent-S-0193, alex, claude-fable-5-1]
created: 2026-10-02T12:24:49Z
updated: 2026-10-02T12:35:27Z
---

# TH-0065 Signing and peer verification for E-0015: the recommendation and four questions

On wip/kanban/stories/S-0193-explore-ways-to-sign-and-verify-components.md.

## Entries

### 2026-10-02T12:24:49Z agent-S-0193
The finding is `design/system/release-signing.md` on `story/S-0193` (sections Today, Signing releases, Verifying the peer). Nothing was tried; tool facts are from knowledge, marked **to check**, and the stories check them.

**Recommendation.** Sign with cosign and a key pair (ECDSA P-256) held as two GitHub Actions secrets. GoReleaser `signs` signs flai's `checksums.txt`; `release-flaiover.yml` signs a digest list of the image and publishes it on a `flaiover vX.Y.Z` GitHub release. flai verifies with the standard library and a public key in the source; `install.sh` verifies with `openssl dgst -verify`; `flai dashboard` runs the image by digest from the signed list. Add `actions/attest-build-provenance` to both workflows as a keyless statement for `gh attestation verify`, costing nothing in flai. For the peer: flai measures the container's image digest through `docker inspect` before it dials and refuses an unlisted one; both sides exchange a signed release stamp (component, version, commit) in `hello` and close with `4403 unsigned peer` when it fails. Unsigned is allowed only by `dashboard.allow_unsigned` on the host, which `flai dashboard` passes to the container, shown as a banner and in status.

**What the peer check cannot do.** A running process cannot prove its own code to a peer without hardware attestation. The stamp proves the peer carries a release's statement, not that it is that release; someone holding the agent credential can replay one. The credential stays the boundary. The docker measurement is a real check because Docker, not the peer, reports the image.

**Questions, recommendation first.**

1. Tool: cosign key pair (recommended), or minisign (Ed25519; smaller, but `install.sh` cannot verify without `minisign`, which hosts lack), or keyless only (no key to keep, but flai needs `sigstore-go` and Sigstore's network)?
2. Image: signed digest list on a GitHub release that flai verifies with stdlib (recommended), with `cosign sign` on the image as an optional extra; or flai verifies cosign's OCI signature in GHCR (a registry client in flai, or heavy dependencies)?
3. Peer check on failure: refuse and back off with `4403`, visible in status and the dashboard (recommended); or warn only in a first release and refuse in a later one?
4. Development builds: `dashboard.allow_unsigned` on the host passed to the container by `flai dashboard`, default off, banner when on (recommended); or a developer signing key the operator configures so that local builds are signed too (no unsigned path, but every local build must be signed)?

Also for your awareness: `latest` is pushed on every merge to `main` today, so `dashboard.tag: latest` would have to mean the newest flaiover release, resolved from the signed list, not the registry tag.

### 2026-10-02T12:33:43Z alex
1 - recommendation, 2 - recommendation, 3 - recommendation, 4 - recommendation

### 2026-10-02T12:35:27Z claude-fable-5-1
Resolved: Recommendation taken on all four questions; recorded as ADR-0070 and release-signing.md § Decision
