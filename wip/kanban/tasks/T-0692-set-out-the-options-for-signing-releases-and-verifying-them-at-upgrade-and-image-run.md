---
id: T-0692
type: task
nature: research
title: Set out the options for signing releases and verifying them at upgrade and image run
status: done
parent: S-0193
owner: arobson
created: 2026-10-02T12:16:02Z
updated: 2026-10-02T12:22:35Z
transitions:
  - to: ready
    at: 2026-10-02T12:19:02Z
    by: claude-fable-5-1
  - to: in-progress
    at: 2026-10-02T12:19:02Z
    by: claude-fable-5-1
  - to: done
    at: 2026-10-02T12:22:35Z
    by: claude-fable-5-1
stream: S-0193
tags: []
touches: [design/system/release-signing.md]
after: [T-0691]
usage:
  source: log
  seconds: 213
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 192
      output: 27
      cache_read: 758446
      cache_write: 17462
      cost: 0
---
# T-0692 Set out the options for signing releases and verifying them at upgrade and image run

## Work

Add to `design/system/release-signing.md` the options for signing what CI releases and verifying it where it is consumed: Ed25519 signatures in the minisign format, cosign with a key pair, cosign keyless with the Sigstore transparency log, GitHub artifact attestations (SLSA provenance), and GPG, as GoReleaser's `signs` supports them; for the image, cosign on the image digest versus a signed digest list published as a release asset. For each: what is signed, where the private key lives, how `flai self-upgrade` and `install.sh` verify with the public key and what that costs in dependencies and network, how `flai dashboard` would verify an image before running it, how keys rotate, and what a compromise means. End with a recommendation. Waits for T-0691 because both edit the same file and the options refer to what `## Today` states.

## Done when

- The document has a `## Signing releases` section with one subsection per option, a comparison table, and a recommendation with its reasons.

## Notes
