---
id: T-1221
type: task
nature: feature
title: release-signing.md's cosign, GoReleaser, and attestation claims are checked against their current documentation and corrected
status: done
parent: S-0232
owner: alex
created: 2026-10-07T22:15:23Z
updated: 2026-10-07T22:23:37Z
transitions:
  - to: ready
    at: 2026-10-07T22:15:47Z
    by: agent-S-0232
  - to: in-progress
    at: 2026-10-07T22:15:47Z
    by: agent-S-0232
  - to: done
    at: 2026-10-07T22:23:37Z
    by: agent-S-0232
stream: S-0232
tags: []
touches: [design/system/release-signing.md]
usage:
  source: log
  seconds: 470
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 77
      output: 22794
      cache_read: 6279874
      cache_write: 141454
      cost: 2.7077
---
# T-1221 release-signing.md's cosign, GoReleaser, and attestation claims are checked against their current documentation and corrected

## Work

Check each claim marked **to check** in `design/system/release-signing.md` § Signing releases that bears on cosign, GoReleaser, or the attestation action against their current documentation: cosign `generate-key-pair` (algorithm, files, `github://`), `sign-blob` flags in cosign v3 (`--output-signature`, `--tlog-upload=false`, `--key env://`, `COSIGN_PASSWORD`), the signature format and its `openssl` verification, GoReleaser `signs` (`artifacts: checksum`, failure behaviour), `actions/attest-build-provenance` (version, permissions, plans and visibility), and Actions environment secrets with required reviewers on the Team plan for a private repository. Correct the document where it is wrong, and record what was checked and where.

Waits for nothing: every other task builds on what it finds.

## Done when

- Each such claim is either confirmed, with its source, or corrected in the document.
- The document says on which date and against which versions it was checked.

## Notes
