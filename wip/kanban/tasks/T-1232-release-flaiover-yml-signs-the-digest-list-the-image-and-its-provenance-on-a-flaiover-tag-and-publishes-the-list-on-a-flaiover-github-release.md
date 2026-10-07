---
id: T-1232
type: task
nature: feature
title: release-flaiover.yml signs the digest list, the image, and its provenance on a flaiover tag and publishes the list on a flaiover GitHub release
status: backlog
parent: S-0234
owner: alex
created: 2026-10-07T22:43:07Z
updated: 2026-10-07T22:43:07Z
transitions: []
stream: S-0234
tags: [dashboard]
touches: [".github/workflows/release-flaiover.yml", scripts/flaiover-digests.sh]
after: [T-1231]
---
# T-1232 release-flaiover.yml signs the digest list, the image, and its provenance on a flaiover tag and publishes the list on a flaiover GitHub release

## Work

Keep the `image` job as it is for every trigger, and give it the index digest of `docker/build-push-action` as a job output. Add a job that runs only on a `flaiover/v*` tag, needs `image`, and uses the environment `release`, which holds `COSIGN_PRIVATE_KEY` and `COSIGN_PASSWORD` and which only release tags deploy to (S-0232's release-key runbook). Give it `contents: write`, `packages: write`, `id-token: write`, and `attestations: write`. It:

- writes `flaiover_<version>.digests` with `scripts/flaiover-digests.sh`, which reads the index with `docker buildx imagetools inspect --raw` and prints the lines T-1231 fixed, one per digest with its reference, so that the format can be run and checked locally;
- fails when either secret is empty, before it signs anything;
- signs the list with `cosign sign-blob --key env://COSIGN_PRIVATE_KEY --tlog-upload=false`, writing `flaiover_<version>.digests.sig`;
- runs `cosign sign --key env://COSIGN_PRIVATE_KEY` on `ghcr.io/bytepunx/flaiover@<index digest>`;
- runs `actions/attest-build-provenance` on the image by that digest with `push-to-registry: true`;
- creates the GitHub release `flaiover vX.Y.Z` on the tag with `gh release create --latest=false`, so that the repository's latest release stays flai's, not a draft or prerelease, so that flai's release client lists it, with both files as assets.

Pin cosign and the attestation action at the versions `design/tech/ci.md` lists for `release-flai.yml`. A push to `main` runs the `image` job alone and publishes `latest` and `sha-*` as today. Check the workflow with `actionlint` if it is at hand, and the script against a published flaiover index.

Waits for T-1231: the lines the script writes are the format it records.

## Done when

- On a `flaiover/v*` tag the workflow writes, signs, and publishes the list and its signature on a `flaiover vX.Y.Z` release, signs and attests the image by digest, and fails without the key; a push to `main` runs only the image job.
- `scripts/flaiover-digests.sh` run against a published flaiover index prints the format `release-signing.md` records.

## Notes
