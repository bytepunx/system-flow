---
id: S-0734
type: story
nature: feature
title: Each flai release builds, signs, and publishes the bedrock and claude-code agent images with a signed digest list, and ships pilot for Linux
status: review
parent: E-0036
owner: alex
created: 2026-11-12T15:57:16Z
updated: 2026-11-13T08:21:01Z
transitions:
  - to: ready
    at: 2026-11-13T07:28:49Z
    by: alex
  - to: in-progress
    at: 2026-11-13T07:53:20Z
    by: agent-S-0734
  - to: review
    at: 2026-11-13T08:21:01Z
    by: agent-S-0734
tags: [cli, ci]
topics: [release, docker, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/agent-image-digests.sh, scripts/agent-image-digests-test.sh, scripts/smoke.sh, scripts/README.md, scripts/ci-pins.sh, design/tech/ci.md, design/tech/docker.md, design/system/release-signing.md, docs/operators/runbooks/agent-container.md, scripts/flaiover-digests.sh, design/tech/go.md, design/system/agent-adapters.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/runbooks/install.md]
after: [S-0733]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1688
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 6
      hand_edits: 1
      work: 49
  models:
    - model: claude-opus-5-5
      input: 208
      output: 76573
      cache_read: 17727473
      cache_write: 488584
      cost: 8.531
  strategic:
    - kind: orchestrator
      seconds: 24
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 75
          cache_read: 225479
          cache_write: 20675
          cost: 0.0709
forecast:
  duration: 45m
  delivery: 2026-11-13T09:53:00Z
  basis: "Its own forecast of 45m; 5th in the pull order with an in-progress limit of 3, behind S-0660, S-0760, S-0735, S-0745 and S-0759."
  by: flai
  at: 2026-11-13T07:53:06Z
finalized:
  by: alex
  at: 2026-11-13T07:28:47Z
---
# S-0734 Each flai release builds, signs, and publishes the bedrock and claude-code agent images with a signed digest list, and ships pilot for Linux

## Goal

No workflow builds or publishes the agent images today, and a released flai has no `pilot` ([ADR-0212](../../../design/adrs/0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md), decision 8 and § Consequences). So in any project but this one, nothing can supply them.

After this story, each `flai/v*` release publishes the bedrock and claude-code images that hold that release's `flai` and `pilot`. They are signed with the release key and named by index digest in a digest list signed beside them, as `release-flaiover.yml` does for flaiover ([release-signing.md § The digest list](../../../design/system/release-signing.md#the-digest-list)). The release also carries `flai` and `pilot` for Linux, which a host build from flai's Dockerfiles needs. It follows the registry, the names, and the list S-0733's ADR decides.

## Acceptance criteria

- [x] A `flai/v*` tag's release workflow builds the bedrock image from the release's own `flai` and `pilot`, and the claude-code image over it, for `linux/amd64` and `linux/arm64`. It pushes both to the registry S-0733's ADR names, tagged with the release's version
- [x] Each image carries the label that ties it to its flai, as the ADR names it, holding the release's version and commit
- [x] Each image is signed with cosign and the release key by its index digest, and attested when the repository is public, as the flaiover image is
- [x] A digest list for the agent images, in the format the ADR names, is signed with the release key and attached to the GitHub release `flai vX.Y.Z`. `scripts/agent-image-digests.sh` writes it, and a test checks its output against the format
- [x] The release carries `pilot` for Linux beside `flai`, in `checksums.txt`, so a host can build the images from flai's Dockerfiles without this tree
- [x] Every base the images build from is pinned by tag and index digest, and `scripts/ci-pins.sh` fails a release Dockerfile that is not
- [x] `design/tech/ci.md`, `design/tech/docker.md`, `design/system/release-signing.md`, and `docs/operators/runbooks/agent-container.md` § The images say what is published and where, and no longer say that nothing publishes the images

## Tasks
- T-3159 GoReleaser builds pilot for Linux beside flai and lists it in checksums.txt
- T-3160 The bedrock Dockerfile takes flai and pilot already built, both images carry the label that ties them to their flai, and every base is pinned by digest
- T-3161 scripts/agent-image-digests.sh writes the agent images' digest list in the ADR's format, with a test in the smoke tier
- T-3162 release-flai.yml builds, pushes, signs, and attests the agent images and attaches their signed digest list to the release
- T-3163 The tech, signing, and agent-container documents say what each flai release publishes of the agent images and where

## Notes

Proposed slice 2 of E-0036. It holds only if S-0733's ADR chooses publication. If the ADR chooses a host build alone, only the `pilot` criterion stays, and the plan thread on E-0036 proposes folding it into S-0740.

The workflow cannot be run from a story's worktree. The story's agent checks it with `actionlint` where the image has it, and with a dry run of `scripts/agent-image-digests.sh` over a fixture. The first real run is the next `flai/v*` tag, which the operator publishes.

### How the criteria were verified

- Criteria 1, 3, and the signing and release half of 4 are `release-flai.yml`'s jobs, which cannot run from a worktree. They were checked by reading them against `release-flaiover.yml`, and with actionlint 1.7.7 run from its image. actionlint's one finding is the `artifact-metadata` scope, which `release-flaiover.yml` grants too and that version does not know. The first real run is the next `flai/v*` tag.
- Criterion 2: bedrock's `prebuilt` target and the claude-code image over it were built locally from linux binaries in the layout the workflow writes. Both carried `dev.system-flow.flai` and `dev.system-flow.commit` with the values given, and ran `flai version`. The `tree` target was built the same way. The workflow fails a push whose platforms lack the labels.
- Criterion 4's script half: `scripts/agent-image-digests-test.sh` passes in the smoke tier.
- Criterion 5: a GoReleaser snapshot wrote `pilot_snapshot_linux_amd64.tar.gz` and `pilot_snapshot_linux_arm64.tar.gz` into `checksums.txt`, with flai's archives holding flai alone.
- Criterion 6: `scripts/ci-pins.sh` passes on the branch, and fails a copy of the bedrock Dockerfile with alpine's digest removed.
- After the first release, the operator makes the GHCR packages `flai-agent-bedrock` and `flai-agent-claude-code` public, as `docs/operators/runbooks/agent-container.md` says. They start private, and ADR-0295 needs them public for a pull with no login.

### Planning

- Touches:
  - `.github/workflows/release-flai.yml`, `flai/.goreleaser.yaml`: design (ADR-0203, The image) and layout. No workflow or GoReleaser section builds the images today; `release-flaiover.yml` is the model.
  - `flai/images/bedrock/Dockerfile`, `flai/images/claude-code/Dockerfile`: layout. The label is at bedrock's line 36, and the binaries are built from the tree there today.
  - `scripts/agent-image-digests.sh`, `scripts/agent-image-digests-test.sh`: new files the third criterion names, modelled on `scripts/flaiover-digests.sh`. `scripts/smoke.sh` runs the script tests, and `scripts/README.md` lists every script (co-change, 17%).
  - `scripts/ci-pins.sh`: layout (lines 78 to 99 check the image Dockerfiles' `FROM`).
  - `design/tech/ci.md`, `design/tech/docker.md`, `design/system/release-signing.md`, `docs/operators/runbooks/agent-container.md`: named in the last criterion.
  - Left out from `flai touches suggest`: `docs/users/flai.md` (36%), `design/system/flai-cli.md` (35%), and `docs/operators/index.md` (25%). Nothing a user runs changes here. `flai-cli.md`'s pull is S-0740's.
- `after` S-0733: this story reads the ADR's registry, image names, label, and digest list format.
- Tasks, in three layers: T-3159, T-3160, and T-3161 share no file and run together; T-3162 waits for all three, since its job downloads `pilot`, builds with the new Dockerfile, and runs the list script; T-3163 describes T-3162's job.
- Forecast 45m: `flai forecast` gives 35m (109 s per unit over 118 done large feature stories, size 19). Raised 10m: five tasks against the usual three or four, and a release job that has to be checked by reading against `release-flaiover.yml`, since it cannot run here. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
