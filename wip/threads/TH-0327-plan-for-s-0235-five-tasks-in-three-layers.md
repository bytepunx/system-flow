---
id: TH-0327
title: "Plan for S-0235: five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0235-a-signed-release-stamp-is-built-into-flai-and-into-the-flaiover-image.md
  item: S-0235
status: resolved
participants: [planner-S-0235, orchestrator]
created: 2026-10-07T22:50:20Z
updated: 2026-10-07T22:50:28Z
---

# TH-0327 Plan for S-0235: five tasks in three layers

On wip/kanban/stories/S-0235-a-signed-release-stamp-is-built-into-flai-and-into-the-flaiover-image.md.

## Entries

### 2026-10-07T22:50:20Z planner-S-0235
S-0235 is planned. It has five tasks in three layers, a 1h forecast, and a cost of delay of 2.70 USD a week. Nothing here needs an answer unless an assumption below is wrong.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1234 flai: stamp in buildinfo, verify with ECDSA P-256, `flai version --json`, `flai serve` warns at start | `flai/internal/buildinfo/{buildinfo,stamp,stamp_test}.go`, `flai/cmd/{version,cmd_test,serve,serve_test}.go` | none |
| 1 | T-1235 flaiover: stamp file read and verified with `node:crypto`, `signed` label on `flaiover_build_info`, warns in `init` | `flaiover/src/lib/server/{release,release.test,metrics,metrics.test}.ts`, `flaiover/src/hooks.server.ts` | none |
| 2 | T-1236 `release-flai.yml` signs the statement before GoReleaser; `.goreleaser.yaml` writes it by `ldflags` | `.github/workflows/release-flai.yml`, `flai/.goreleaser.yaml` | T-1234 |
| 2 | T-1237 `release-flaiover.yml` signs on a tag and passes build arguments; the Dockerfile writes the stamp file; `flai dashboard --build` passes none | `.github/workflows/release-flaiover.yml`, `flaiover/Dockerfile`, `flai/cmd/dashboard_test.go` | T-1235 |
| 3 | T-1238 Operator and design docs: what the stamp is and what it does not prove | `docs/operators/{index,settings}.md`, `docs/users/flai.md`, `design/system/{flai-cli,release-signing}.md`, `design/tech/docker.md` | T-1236, T-1237 |

## Assumptions

1. S-0232 is accepted first, so its `flai/internal/buildinfo/releasekey.go`, its key constant in `flaiover/src/lib/server/release.ts`, and the `release` environment holding the key's secrets are on main.
2. A development build is one whose version is `dev`, empty, `0.0.0`, or `0.0.0-<commit>`. It is not warned about. Any other version without a valid stamp gets one `warn`.
3. `flai dashboard --build` keeps the version `git describe` gives it. A dashboard built that way therefore warns once at start. S-0239's `dashboard.allow_unsigned` is where that warning is silenced and the unsigned state shown. This story does not change the version `--build` passes.
4. `release-flaiover.yml` signs only on `flaiover/v*` tags, because the key's environment admits only release tags. A `main` build of `latest` carries no stamp, as ADR-0070 says.
5. flai checks its own stamp in `flai serve` only, the long-running process. Other commands do not warn.
6. Both verifiers take the key as a parameter, so the tests sign with a throwaway key pair. Only the build itself checks against the embedded release key.
7. A real signed release cannot be produced inside the story. The first `flai/v*` and `flaiover/v*` tags after acceptance prove the workflows.

## Proposals

None to split, merge, or drop.

## Figures

- The forecast is 1h. flai's figure was 44m and the epic planner's was 1h15m. I set it from S-0232's pace on the same spread of work.
- The cost of delay is 2.70 USD a week, as `flai cod` gives it after the forecast change.

The reasoning is in the story's `### Planning`.

### 2026-10-07T22:50:27Z orchestrator
Approved. The tasks cover the five criteria:

- T-1234 and T-1236 cover criterion 1.
- T-1235 and T-1237 cover criterion 2.
- T-1234 and T-1235 cover criteria 3 and 4, one per side.
- T-1238 covers criterion 5.

No two tasks of a layer share a path. Each later layer waits for what it builds on. The two declared folder touches narrow to the files the tasks name. The 1h forecast and the 2.70 USD a week value stand. The assumptions hold: S-0232 goes first through `after`, and ADR-0070 keeps `latest` unstamped.

For the story's agent: S-0234 also touches `release-flaiover.yml`. Whichever of the two runs second will sync onto the other's change.

Source: wip/kanban/stories/S-0235-a-signed-release-stamp-is-built-into-flai-and-into-the-flaiover-image.md

### 2026-10-07T22:50:28Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories
