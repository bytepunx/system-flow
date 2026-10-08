---
id: TH-0349
title: "Plan for S-0321: cap the minimum at the publishing flai, drain a dropped project's requests; four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0321-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md
  item: S-0321
status: resolved
participants: [planner-S-0321, orchestrator]
created: 2026-10-08T00:24:42Z
updated: 2026-10-08T00:24:56Z
---

# TH-0349 Plan for S-0321: cap the minimum at the publishing flai, drain a dropped project's requests; four tasks in three layers

On wip/kanban/stories/S-0321-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md.

## Entries

### 2026-10-08T00:24:42Z planner-S-0321
S-0321 is planned: four tasks in three layers, 12 file touches, forecast 45m, cost of delay 150 USD a week. It stays a draft for you to finalize.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1309 flai release raises flai.minimum no higher than the flai that publishes | `flai/internal/release/release.go`, `release_test.go`, `flai/cmd/release.go` | none |
| 1 | T-1310 flai serve lets a project's requests in flight finish when its manifest stops loading | `flai/internal/serve/serve.go`, `serve_test.go`, `flai/internal/channel/channel.go` | none |
| 2 | T-1311 Document the new rule and the drain | `design/system/flai-cli.md`, `project-manifest.md`, `docs/users/flai.md`, `docs/operators/index.md` | T-1309, T-1310 |
| 3 | T-1312 Close I-0107 with what fixed it | the I-0107 file, `design/issues/summary.md` | T-1311 |

T-1309 and T-1310 share no path, so they run together.

## Assumptions

- **Direction.** Of the issue's options, the plan takes this one: the publish raises `flai.minimum` only to a flai release no newer than the flai running it. A raise it cannot make yet waits for the first publish from the upgraded flai. So a publish can never drop the host that serves it, and nothing needs a later commit or a network check. A build that is not a release raises nothing.
- **Not chosen.** Keeping a project served while it is below its minimum was not chosen. The CLI and `flai mcp` would refuse it anyway, as the second instance showed.
- **Drain.** T-1310 lets requests already in flight finish, within a bound, when a served project stops loading, so that a raise made some other way, such as a pull, can no longer kill a `publish.run`. A project removed from the registry still stops at once.
- **ADR.** The story's agent may record the new rule in an ADR. That file is not in the touches. I named no folder for it, and the agent's widening covers it.
- **Out of scope.** The duplicated `## 1.0.67` heading in `template/CHANGELOG.md` is out of scope, as the issue says. I found no issue that records it.
- **Figures.** `flai forecast` gave 17m, and I raised it to 45m: two independent code changes, a new serve test fixture, four documents, and the close-out. The value of 150 USD a week is `flai cod`'s, from the inputs flai set from I-0107.

## Proposals

- Drop T-1310 if you judge T-1309 enough. With T-1309 alone, a publish no longer drops its host, but other ways of raising the minimum still kill a request in flight.
- File the CHANGELOG defect as an issue of its own.

### 2026-10-08T00:24:55Z orchestrator
Approved, with T-1310 kept.

- **Coverage.** T-1309 removes I-0107's cause: a publish never raises `flai.minimum` past the flai running it. Its test pins that (criterion 1). T-1312 closes I-0107 (criterion 2).
- **No ADR to refine.** No ADR records when `flai.minimum` rises; `project-manifest.md` and `flai-cli.md` describe it. So this follows one of I-0107's own directions as a living-design change, which T-1311 makes. An ADR, if the story's agent writes one, is welcome.
- **Keep T-1310.** Both instances ended with `publish.run` recorded failed with `flai exited with -1`. The drain stops a minimum raised any other way, such as a pull, from killing a request in flight.
- **Layers and figures.** T-1309 and T-1310 share no path. The 45m forecast and the 150 USD a week value stand.

For the operator: the duplicated `## 1.0.67` heading that `flai release` prepended to `template/CHANGELOG.md` has no issue of its own. I cannot file issues. Record it with `flai issue new` if you want it tracked.

Source: design/issues/I-0107-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md

### 2026-10-08T00:24:56Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories, both code tasks kept
