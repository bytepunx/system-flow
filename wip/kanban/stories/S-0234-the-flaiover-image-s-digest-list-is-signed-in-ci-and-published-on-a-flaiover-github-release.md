---
id: S-0234
type: story
nature: feature
title: The flaiover image's digest list is signed in CI and published on a flaiover GitHub release
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-02T12:37:23Z
transitions: []
tags: [dashboard]
topics: [release, security]
touches: [".github/workflows/release-flaiover.yml", design/tech/ci.md, docs/operators]
after: [S-0232]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0234 The flaiover image's digest list is signed in CI and published on a flaiover GitHub release

## Goal

Every `flaiover/v*` tag produces a GitHub release `flaiover vX.Y.Z` carrying a digest list, the index digest and each platform's manifest digest the workflow pushed to GHCR, signed with the release key (ADR-0070), so that flai can choose and verify an image without a registry client.

## Acceptance criteria
- [ ] On a `flaiover/v*` tag, `release-flaiover.yml` writes `flaiover_<version>.digests` from `docker/build-push-action`'s outputs with one line per digest and the reference it was pushed under, signs it with `cosign sign-blob --tlog-upload=false`, and creates the GitHub release with both files as assets; a run without the key fails.
- [ ] The workflow also runs `cosign sign --key` on the image by digest and `actions/attest-build-provenance` with `push-to-registry`, so that `cosign verify` and `gh attestation verify oci://` work for anyone.
- [ ] A push to `main` still publishes `latest` and `sha-*` as today, and the operator documentation says those builds are unsigned and what that means once `flai dashboard` verifies.
- [ ] Which digest `RepoDigests` records for an image pulled by index digest and by platform digest was checked, and the list carries every digest a pull can leave behind.
- [ ] `design/tech/ci.md` describes the release and its assets.

## Tasks

## Notes
