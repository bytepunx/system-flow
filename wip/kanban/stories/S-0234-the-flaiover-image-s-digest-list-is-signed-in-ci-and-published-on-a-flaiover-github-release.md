---
id: S-0234
type: story
nature: feature
title: The flaiover image's digest list is signed in CI and published on a flaiover GitHub release
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-07T20:31:36Z
transitions: []
tags: [dashboard]
topics: [release, security]
touches: [".github/workflows/release-flaiover.yml", design/tech/ci.md, docs/operators, design/tech/docker.md, design/system/release-signing.md]
after: [S-0232]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 45m
  delivery: 2026-10-08T02:59:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 3, behind S-0329, S-0330, S-0232 and S-0233."
  by: flai
  at: 2026-10-07T20:28:28Z
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

### Planning

Touches:

- Declared: `.github/workflows/release-flaiover.yml`, `design/tech/ci.md`, `docs/operators`.
- Design: `design/tech/docker.md`, whose Tags row names what `release-flaiover.yml` publishes; `design/system/release-signing.md`, where criterion 4's `RepoDigests` finding is recorded.
- Folder touch kept, as declared: `docs/operators`, where the update runbook and perhaps a new note on unsigned `latest` builds change.

Forecast: 45m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 16m from 116 s per unit of size, over only 3 medium-band feature stories.
- It is one workflow, so under the 1h median of done feature stories. But criterion 4 needs pulls by index and platform digest to be tried with Docker, so 45m.
- The first delivery was played out after S-0232 at flai's cycle factor of 6.85.

Cost of delay: no value yet. E-0015 and its stories have no inputs; TH-0312 asks the operator for them.
