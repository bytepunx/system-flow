---
id: S-0235
type: story
nature: feature
title: A signed release stamp is built into flai and into the flaiover image
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-07T22:25:33Z
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
cost_of_delay:
  value: 3.29
  by: planner-E-0015
  at: 2026-10-07T22:13:43Z
forecast:
  duration: 1h15m
  delivery: 2026-10-08T01:03:00Z
  basis: "Its own forecast of 1h15m; 4th in the pull order with an in-progress limit of 3, behind S-0232, S-0333, S-0332, S-0233 and S-0234."
  by: flai
  at: 2026-10-07T22:25:33Z
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

Forecast: 1h15m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 24m from 116 s per unit of size, over only 3 medium-band feature stories.
- Done feature stories of this size took a median of about 1h of agent time. This one builds, verifies, and tests a stamp in both components and both workflows, so 1h15m.
- The first delivery was played out after S-0232 at flai's cycle factor of 6.85.

Cost of delay: 3.29 USD a week, as `flai cod` gives it: this story's 1h15m share of the 9h30m forecast over E-0015's eight open stories, of the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, and that exposure is closed only when the chain is done, so a share by work fits.
