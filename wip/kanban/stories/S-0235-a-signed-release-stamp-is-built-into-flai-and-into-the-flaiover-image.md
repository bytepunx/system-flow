---
id: S-0235
type: story
nature: feature
title: A signed release stamp is built into flai and into the flaiover image
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-07T19:34:53Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flai/internal/buildinfo, ".github/workflows/release-flaiover.yml", flaiover/Dockerfile, flaiover/src/lib/server/release.ts, docs/operators, flai/cmd/version.go, flai/cmd/serve.go, flaiover/src/lib/server/release.test.ts, flaiover/src/lib/server/metrics.ts, flaiover/src/lib/server/metrics.test.ts, flaiover/src/hooks.server.ts, design/tech/docker.md]
after: [S-0232]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 1h15m
  delivery: 2026-10-08T10:58:00Z
  basis: "flai's 116 s per unit of size rests on 3 medium-band stories and gave 24m; raised to 1h15m, above the 1h median of done feature stories, for a stamp built, verified, and tested in both components and both workflows; delivery played out after S-0232 at a cycle factor of 6.85"
  by: planner-E-0015
  at: 2026-10-07T19:33:55Z
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

### Planning

Touches:

- Declared: `.github/workflows/release-flai.yml`, `flai/.goreleaser.yaml`, `flai/internal/buildinfo`, `.github/workflows/release-flaiover.yml`, `flaiover/Dockerfile`, `flaiover/src/lib/server/release.ts`, `docs/operators`.
- Layout: `flai/cmd/version.go` for `flai version --json`; `flai/cmd/serve.go`, where flai checks its own stamp at start; `flaiover/src/hooks.server.ts`, where flaiover does; `flaiover/src/lib/server/metrics.ts` and its test for the `signed` label of `flaiover_build_info`; `flaiover/src/lib/server/release.test.ts`.
- Design: `design/tech/docker.md`, which describes the image.
- Folder touches kept, as declared: `flai/internal/buildinfo`, which may gain a stamp file beside `buildinfo.go`; `docs/operators`.

Forecast: 1h15m, delivery 2026-10-08T10:58Z.

- `flai forecast` gave 24m from 116 s per unit of size, over only 3 medium-band feature stories.
- Done feature stories of this size took a median of about 1h of agent time. This one builds, verifies, and tests a stamp in both components and both workflows, so 1h15m.
- Delivery is played out after S-0232 at flai's cycle factor of 6.85.

Cost of delay: no value yet. E-0015 and its stories have no inputs; TH-0312 asks the operator for them.
