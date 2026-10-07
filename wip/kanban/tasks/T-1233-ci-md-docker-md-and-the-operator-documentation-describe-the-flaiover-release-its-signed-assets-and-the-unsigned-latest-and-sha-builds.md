---
id: T-1233
type: task
nature: feature
title: ci.md, docker.md, and the operator documentation describe the flaiover release, its signed assets, and the unsigned latest and sha builds
status: backlog
parent: S-0234
owner: alex
created: 2026-10-07T22:43:16Z
updated: 2026-10-07T22:43:16Z
transitions: []
stream: S-0234
tags: [dashboard, docs]
touches: [design/tech/ci.md, design/tech/docker.md, docs/operators/runbooks/update.md, docs/operators/runbooks/release-key.md, docs/operators/settings.md]
after: [T-1232]
---
# T-1233 ci.md, docker.md, and the operator documentation describe the flaiover release, its signed assets, and the unsigned latest and sha builds

## Work

- `design/tech/ci.md`: the `release-flaiover.yml` row and the Releases row say that a `flaiover/v*` tag creates the GitHub release `flaiover vX.Y.Z` with `flaiover_<version>.digests` and its `.sig`, signs the image with cosign, and attests its provenance, and that a push to `main` does none of it.
- `design/tech/docker.md`: the Tags row says that only the semver and major tags of a release are signed and listed, and that `latest` and `sha-*` are not.
- `docs/operators/runbooks/update.md`: where it says `latest` follows the main branch, say that `latest` and `sha-*` builds are unsigned, that once `flai dashboard` verifies (S-0236) it runs only an image a signed list names, so those builds run only with `dashboard.allow_unsigned` (S-0239), and how to check a release image by hand with `cosign verify --key` and `gh attestation verify oci://`.
- `docs/operators/runbooks/release-key.md` and `docs/operators/settings.md`: replace "later `release-flaiover.yml`" and "later flaiover's digest list" with what the workflow now signs with the two secrets.

Waits for T-1232: it describes what the workflow does.

## Done when

- The five files describe the flaiover release and its assets as the workflow builds them, and the markdown lint passes on them.

## Notes
