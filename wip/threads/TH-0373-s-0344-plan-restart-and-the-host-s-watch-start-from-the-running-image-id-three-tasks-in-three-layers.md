---
id: TH-0373
title: "S-0344 plan: restart and the host's watch start from the running image ID; three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0344-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md
  item: S-0344
status: resolved
participants: [planner-S-0344, orchestrator]
created: 2026-10-08T08:36:32Z
updated: 2026-10-08T08:36:43Z
---

# TH-0373 S-0344 plan: restart and the host's watch start from the running image ID; three tasks in three layers

On wip/kanban/stories/S-0344-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md.

## Entries

### 2026-10-08T08:36:32Z planner-S-0344
S-0344 is planned: 11 file touches, three tasks, forecast 30m, cost of delay 12.50 USD a week.

## Proposed fix

A restart starts from the image ID the container runs (`{{.Image}}`), not `{{.Config.Image}}`, the tag that `flai dashboard check` may have moved. S-0316's upgrade fallback already does this. Status and output keep naming the tag. The recommended way is a container label holding the ref, which `containerInfo` prefers.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1360 restart from the image ID, with a reproducing test | none | `flai/cmd/dashboard_upgrade.go`, `flai/cmd/dashboard.go`, `flai/cmd/dashboard_test.go`, `docs/users/flai-reference.md` |
| 2 | T-1361 the host's watch restarts from the recorded image ID | T-1360 | `flai/cmd/dashboard_watch.go`, `flai/cmd/dashboard_watch_test.go` |
| 3 | T-1362 docs, and `flai issue close I-0116` | T-1360, T-1361 | `docs/users/flai.md`, `docs/operators/index.md`, `design/system/flai-cli.md`, the I-0116 file, `design/issues/summary.md` |

T-1361 waits for T-1360 because it reuses T-1360's way of running by ID and naming the tag. Both change how `startDashboard` is called.

## Assumptions

1. The host's watch is in scope. It restarts from the recorded tag for the same reason, and ADR-0118 §1 says both restarts keep the release that runs. Criterion 1 asks that the cause no longer occur. If you want the story kept to `flai dashboard restart` alone, I would drop T-1361 and file the watch as its own issue.
2. When the recorded image ID is gone, such as after `docker image prune`, the watch falls back to the recorded ref with a warning. It does not leave the dashboard down.
3. The story's agent records the chosen fix in `## Decisions` before building. The story asks for a proposal first. The label is a recommendation, not a requirement.
4. Forecast is raised from flai's 17m to 30m. The size median does not count choosing the fix or the two reproducing tests.

No split, merge, or drop proposed.

### 2026-10-08T08:36:42Z orchestrator
Plan approved: three tasks in three layers, 11 file touches, a 30m forecast, and 12.50 USD a week as `flai cod` gives it. The tasks match the criteria. T-1360 and T-1361 reproduce and fix I-0116, and T-1362 closes it.

On assumption 1: keep T-1361, the watch, in this story. ADR-0118 § Decision ¶1 says "A restart, `flai dashboard restart` or `flai host`'s watch, keeps the release that runs." The watch restarting from the recorded tag breaks that rule for the same cause, and criterion 1 asks that the cause no longer occur. Assumptions 2 (fall back to the ref, with a warning, when the image ID is gone) and 3 (the agent records its chosen fix in `## Decisions` first) are sound as written.

S-0344 is still a draft. Finalizing it is a separate step that `flai promote --drafts` covers.

Source: design/adrs/0118-a-dashboard-release-chosen-on-the-updates-page-keeps-running-through-restarts.md § Decision

### 2026-10-08T08:36:43Z orchestrator
Resolved: Plan approved; T-1361 kept in scope under ADR-0118 § Decision ¶1.
