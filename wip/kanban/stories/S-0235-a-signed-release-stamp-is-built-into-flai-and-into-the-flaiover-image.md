---
id: S-0235
type: story
nature: feature
title: A signed release stamp is built into flai and into the flaiover image
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-02T12:37:23Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flai/internal/buildinfo, ".github/workflows/release-flaiover.yml", flaiover/Dockerfile, flaiover/src/lib/server/release.ts, docs/operators]
after: [S-0232]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0235 A signed release stamp is built into flai and into the flaiover image

## Goal

Each release build carries a statement CI signed with the release key before the build, the component's name, version, and commit (ADR-0070), so that a running flai or flaiover can show its peer which release it is. Development builds carry none.

## Acceptance criteria
- [ ] `release-flai.yml` signs the line `flai <version> <commit>` with `cosign sign-blob` before GoReleaser runs and hands the signature to it through the environment; `flai/.goreleaser.yaml` writes the statement and the signature into `buildinfo` by `ldflags`; `flai version --json` shows them, and a snapshot or `scripts/flai.sh` build shows none.
- [ ] `release-flaiover.yml` signs `flaiover <version> <commit>` and passes both as build arguments; the image keeps them as a file the server reads; `flaiover_build_info` gains a `signed` label, and `flai dashboard --build` produces an image with no stamp.
- [ ] Both components verify their own stamp at start with the embedded public key and log one `warn` when a build that has a version carries no valid stamp.
- [ ] Tests cover a valid stamp, a missing one, and one signed by another key, on both sides.
- [ ] The operator documentation says what the stamp is and what it does not prove, in the words of `release-signing.md § Verifying the peer`.

## Tasks

## Notes
