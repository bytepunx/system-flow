---
id: S-0235
type: story
nature: feature
title: A signed release stamp is built into flai and into the flaiover image
status: backlog
parent: E-0015
owner: arobson
created: 2026-10-02T12:37:23Z
updated: 2026-10-08T04:49:57Z
transitions: []
tags: [cli, dashboard]
topics: [release, security]
touches: [".github/workflows/release-flai.yml", flai/.goreleaser.yaml, flai/internal/buildinfo, ".github/workflows/release-flaiover.yml", flaiover/Dockerfile, flaiover/src/lib/server/release.ts, docs/operators, flai/cmd/version.go, flai/cmd/serve.go, flaiover/src/lib/server/release.test.ts, flaiover/src/lib/server/metrics.ts, flaiover/src/lib/server/metrics.test.ts, flaiover/src/hooks.server.ts, design/tech/docker.md, flai/cmd/cmd_test.go, flai/cmd/serve_test.go, flai/cmd/dashboard_test.go, docs/users/flai.md, design/system/flai-cli.md, design/system/release-signing.md]
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
      seconds: 88
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 26
          output: 397
          cache_read: 2383459
          cache_write: 18117
          cost: 0.5917
cost_of_delay:
  value: 2.7
  by: planner-S-0235
  at: 2026-10-07T22:50:03Z
forecast:
  duration: 1h
  delivery: 2026-10-08T07:53:00Z
  basis: "Its own forecast of 1h; 9th in the pull order with an in-progress limit of 3, behind S-0232, S-0318, S-0320, S-0319, S-0309, S-0336, S-0326, S-0287, S-0322, S-0233 and S-0234."
  by: flai
  at: 2026-10-08T04:49:57Z
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
- T-1234 flai carries its release stamp in buildinfo, shows it in flai version --json, and flai serve warns at start without a valid one
- T-1235 flaiover reads and verifies its release stamp from a file, labels flaiover_build_info signed, and warns at start without a valid one
- T-1236 release-flai.yml signs flai's release statement before GoReleaser, and .goreleaser.yaml writes it into buildinfo by ldflags
- T-1237 release-flaiover.yml signs flaiover's release statement and passes it to the image, which keeps it as the stamp file, and flai dashboard --build passes none
- T-1238 The operator and design documentation say what the release stamp is, where each build carries it, and what it does not prove

## Notes

### Planning

Tasks, in three layers:

| Layer | Tasks | Why |
|-------|-------|-----|
| 1 | T-1234 (flai stamp, version, serve warn), T-1235 (flaiover stamp, metrics, init warn) | No path in common; both use S-0232's public key, on main first |
| 2 | T-1236 (flai release workflow and ldflags) after T-1234; T-1237 (flaiover workflow and image) after T-1235 | Each builds what its layer-1 task reads: the variable names, the stamp file |
| 3 | T-1238 (docs) after T-1236 and T-1237 | Documents the stamp as built |

Touches:

- Declared: `.github/workflows/release-flai.yml`, `flai/.goreleaser.yaml`, `flai/internal/buildinfo`, `.github/workflows/release-flaiover.yml`, `flaiover/Dockerfile`, `flaiover/src/lib/server/release.ts`, `docs/operators`, `flai/cmd/version.go`, `flai/cmd/serve.go`, `flaiover/src/lib/server/release.test.ts`, `flaiover/src/lib/server/metrics.ts`, `flaiover/src/lib/server/metrics.test.ts`, `flaiover/src/hooks.server.ts`, `design/tech/docker.md`. All kept.
- Layout, added by this run: `flai/cmd/cmd_test.go`, which holds `TestVersionPlainAndJSON`; `flai/cmd/serve_test.go`, for the start warning `flai/cmd/serve.go` logs before `serve.Run`; `flai/cmd/dashboard_test.go`, for criterion 2's `flai dashboard --build` with no stamp.
- Design, added by this run: `docs/users/flai.md` and `design/system/flai-cli.md`, which document `flai version`; `design/system/release-signing.md`, which records the stamp as built.
- Co-change listed `design/system/flai-cli.md` (53%) and `docs/users/flai.md` (49%), added above; `docs/users/flai-reference.md` (35%) is generated from flags and no flag changes; `design/system/flaiover-dashboard.md` (32%) does not describe the build info, so not added.
- Folder touches kept, as declared: `flai/internal/buildinfo` and `docs/operators`. Their tasks name the files (`buildinfo.go`, `stamp.go`, `stamp_test.go`; `docs/operators/index.md`, `docs/operators/settings.md`), so the claim narrows to them (ADR-0096).

Forecast: 1h, adjusted from 1h15m. flai replays the delivery from the pull order whenever it changes.

- `flai forecast` gave 34m before this run's touches and 44m after, from about 105 s per unit of size over 33 done large-band feature stories.
- S-0232, of the same spread across CI, Go, and TypeScript, closed four tasks in 14 minutes and its one code task ran more than 20 minutes more. This story has two code tasks with tests on both sides, two workflow tasks, and docs, so 1h: above flai's figure, below the earlier 1h15m.

Cost of delay: 2.70 USD a week, as `flai cod` gives it after the forecast changed: this story's 1h share of the 9h15m forecast over E-0015's eight open stories, of the epic's 25 USD a week penalty, which the operator set on TH-0312. Kept as given: each story closes part of one exposure, closed only when the chain is done, so a share by work fits.
