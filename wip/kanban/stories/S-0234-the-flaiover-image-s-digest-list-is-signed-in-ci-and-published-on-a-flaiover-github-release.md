---
id: S-0234
type: story
nature: feature
title: The flaiover image's digest list is signed in CI and published on a flaiover GitHub release
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-08T08:38:58Z
transitions: []
tags: [dashboard]
topics: [release, security]
touches: [".github/workflows/release-flaiover.yml", design/tech/ci.md, docs/operators, design/tech/docker.md, design/system/release-signing.md, scripts/flaiover-digests.sh, docs/operators/runbooks/update.md, docs/operators/runbooks/release-key.md, docs/operators/settings.md]
after: [S-0232]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 158
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 248
          cache_read: 1380358
          cache_write: 8211
          cost: 0.3421
cost_of_delay:
  value: 1.97
  by: planner-E-0015
  at: 2026-10-07T22:13:41Z
forecast:
  duration: 45m
  delivery: 2026-10-08T10:09:00Z
  basis: "Its own forecast of 45m; 5th in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322, S-0340, S-0341, S-0344, S-0342, S-0343 and S-0233."
  by: flai
  at: 2026-10-08T08:38:58Z
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
- T-1231 Which digest RepoDigests records for a flaiover image pulled by index and by platform digest is checked, and the digest list's format recorded in release-signing.md
- T-1232 release-flaiover.yml signs the digest list, the image, and its provenance on a flaiover tag and publishes the list on a flaiover GitHub release
- T-1233 ci.md, docker.md, and the operator documentation describe the flaiover release, its signed assets, and the unsigned latest and sha builds

## Notes

### Planning

Tasks, one per layer, each waiting for the one before:

1. T-1231 checks `RepoDigests` after pulls by index and by platform digest, and records the list's format in `release-signing.md` (criterion 4).
2. T-1232 adds the tag-only signing and release job and `scripts/flaiover-digests.sh` (criteria 1 to 3). It waits for T-1231's format.
3. T-1233 updates `ci.md`, `docker.md`, and the operator documentation (criteria 3 and 5). It waits for T-1232, which it describes.

Touches:

- Declared: `.github/workflows/release-flaiover.yml`, `design/tech/ci.md`, `docs/operators`.
- Design, kept from E-0015's planner: `design/tech/docker.md`, whose Tags row names what `release-flaiover.yml` publishes; `design/system/release-signing.md`, where criterion 4's finding is recorded.
- Design: `docs/operators/runbooks/release-key.md` and `docs/operators/settings.md`, where S-0232 wrote that the two secrets sign flaiover's list "later"; `docs/operators/runbooks/update.md`, which says `latest` follows the main branch and gains the unsigned-build note of criterion 3.
- Layout: `scripts/flaiover-digests.sh`, which writes the list, so that its format can be run outside CI, as `tooling.md` asks of a command sequence.
- Folder touch kept, as declared: `docs/operators`. Its tasks name three files under it, which replace it in the story's claim (ADR-0096).
- Co-change listed `design/system/flai-cli.md` at 59% and `docs/users/flai.md` at 56%. Not added: no command changes here; S-0236 changes what `flai dashboard` runs.

Assumptions:

- The signing job runs only on a tag, in the environment `release`, because S-0232 put the two secrets there and only release tags deploy to it. The `image` job on `main` is untouched.
- The image is signed and attested by its index digest, on tags only, so that `latest` and `sha-*` carry nothing that looks like a release.
- The release is created with `--latest=false`, so that the repository's latest release stays flai's. flai's release client and `install.sh` filter by tag prefix, so neither breaks either way.

Forecast: 45m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 20m from 116 s per unit of size, over only 3 medium-band feature stories.
- Kept at 45m: T-1231 pulls a multi-platform image of about 670 MB twice, and the workflow cannot be run on a tag before review, so it is checked by reading, `actionlint`, and the script alone.
- S-0232, of like scope, closed four tasks in about 15 minutes, so no more than 45m.

Cost of delay: 1.97 USD a week, as `flai cod` gives it: this story's 45m share of the 9h30m forecast over E-0015's eight open stories, of the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure is closed only when the chain is done, so a share by work fits.
