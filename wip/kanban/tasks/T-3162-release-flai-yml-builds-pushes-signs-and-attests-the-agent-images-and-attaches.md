---
id: T-3162
type: task
nature: feature
title: release-flai.yml builds, pushes, signs, and attests the agent images and attaches their signed digest list to the release
status: backlog
parent: S-0734
owner: alex
created: 2026-11-12T15:57:51Z
updated: 2026-11-12T15:57:51Z
transitions: []
stream: S-0734
tags: [cli]
touches: [".github/workflows/release-flai.yml"]
after: [T-3159, T-3160, T-3161]
---
# T-3162 release-flai.yml builds, pushes, signs, and attests the agent images and attaches their signed digest list to the release

## Work

Add a job to `.github/workflows/release-flai.yml` after the GoReleaser job. It:

1. downloads the release's Linux `flai` and `pilot`, and checks them against `checksums.txt`;
2. builds the bedrock image from them, then the claude-code image over its digest, with `docker/build-push-action` for `linux/amd64` and `linux/arm64`, with the label's build arguments;
3. pushes both to the registry S-0733's ADR names, tagged with the version;
4. signs each index digest with `cosign sign --key`, as `release-flaiover.yml` does, and attests it with `actions/attest-build-provenance` when the repository is public;
5. writes the digest list with `scripts/agent-image-digests.sh`, signs it with `cosign sign-blob`, and uploads both to the release.

Pin every action by SHA with its version comment, as the file's others are, and run the job in the `release` environment.

It waits for:

- T-3159, whose `pilot` it downloads;
- T-3160, whose Dockerfile takes the binaries and the label;
- T-3161, whose script writes the list.

## Done when

- `actionlint`, where the image has it, passes on the workflow, and `scripts/ci-pins.sh` passes.
- A dry read of the job against `release-flaiover.yml`'s image job finds the same signing, attestation, and environment.

## Notes
