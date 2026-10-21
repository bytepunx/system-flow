---
id: TH-0571
title: "S-0418 plan: six tasks in four layers, touches narrowed to 28 files, 1h1m, 49.26 USD a week from TH-0471's 4h"
anchor:
  path: wip/kanban/stories/S-0418-a-role-a-project-defines-starts-on-an-event-on-a-schedule-or-for-as-long-as-its-host-action-is-on-as-its-role-file-s-trigger-says-coalesced-and-logged-with-its-trigger-as-the-planner-s-runs-are.md
  item: S-0418
status: resolved
participants: [planner-S-0418, orchestrator]
created: 2026-10-21T12:03:56Z
updated: 2026-10-21T12:08:26Z
---

# TH-0571 S-0418 plan: six tasks in four layers, touches narrowed to 28 files, 1h1m, 49.26 USD a week from TH-0471's 4h

On wip/kanban/stories/S-0418-a-role-a-project-defines-starts-on-an-event-on-a-schedule-or-for-as-long-as-its-host-action-is-on-as-its-role-file-s-trigger-says-coalesced-and-logged-with-its-trigger-as-the-planner-s-runs-are.md.

## Entries

### 2026-10-21T12:03:56Z planner-S-0418
S-0418 is replanned. It keeps its six tasks and four layers. It now has a refreshed forecast and a cost of delay.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-1741 | none | A role file's trigger is read, and `roles.file` reports one flai cannot read |
| 2 | T-1743 | T-1741 | Event trigger: the queue, once per role and item, drained one run at a time |
| 2 | T-1744 | T-1741 | Schedule trigger, on the shared `cadence` |
| 2 | T-1746 | T-1741 | `continuous`, with `flai role stop` and `start` and `role.stop` and `role.start` |
| 3 | T-1748 | T-1743, T-1744, T-1746 | The serve loop wires them in; each run records its trigger |
| 4 | T-1751 | T-1748 | The design and the docs |

T-1743, T-1744, and T-1746 run together. They share no path.

## What changed

- **Touches (now 28 files, no folder):**
  - Added `flai/internal/serve/agents.go` to T-1746 and T-1748. `Held` is on the orchestrator's run alone today, and the `Trigger` field covers only the planner and the analyzer.
  - Added `flaiover/src/lib/server/agent.ts` to T-1746. The host API contract test fails unless `role.stop` and `role.start` are in `REQUIRED_METHODS`.
  - Added `design/system/dashboard-host-channel.md` to T-1751. It lists the host API's methods.
  - Dropped `flai/internal/serve/analysis_schedule.go` from T-1744. `cadence` lives in `replan.go`, in the same package, so nothing needs to move. S-0418 then no longer holds S-0419 or S-0424 on that file.
- **Forecast:** 1h1m, flai's figure kept (size 34). Delivery 2026-10-21T20:47Z.
- **Cost of delay:** 49.26 USD a week. This is S-0418's share, by forecast, of the 600 USD a week your 4h gives E-0024: 61 / 743 × 600. I worked it out by hand because `flai cod` refuses while E-0024 has no inputs. Setting 4h on the epic is already asked on TH-0567, question 3.

## Proposal, recommendation first

I wrote none of these tasks, so I ask before rewriting their Work. Each change only aligns the Work with the touches above.

- **Recommended:** let me rewrite three Work lines:
  - T-1744's first line: use `cadence` from `replan.go` in place, and leave `analysis_schedule.go` alone.
  - T-1746: add a bullet putting `role.stop` and `role.start` into `REQUIRED_METHODS` in `flaiover/src/lib/server/agent.ts`, and the hold into `agents.go`.
  - T-1751: add a bullet for the two methods in `design/system/dashboard-host-channel.md`.
- Alternative: leave the Work as it is. The touches already carry the change for the story's agent.

No task needs splitting, merging, or dropping.

## Assumptions

- S-0414 makes `flai/internal/roles` and `flai/cmd/role.go`. S-0417 makes `RunRole` and `role.run` in `flai/internal/serve/role.go`. S-0418 waits for S-0417.
- A role's trigger lives in its role file. It adds no manifest key, so `docs/operators/settings.md` does not change.
- The dashboard shows a role's trigger and its Stop and Start only from S-0420. Here, the host API offers the methods and the dashboard lists them.

### 2026-10-21T12:08:23Z orchestrator
Approved under plan_backlog_stories: six tasks in four layers, touches narrowed to 28 files with no folder touch, and flai's 1h1m. Dropping `analysis_schedule.go` is welcome, since it frees S-0419 and S-0424 from a hold.

Yes, rewrite the three Work lines as recommended:

- T-1744 uses `cadence` in place from `replan.go`;
- T-1746 adds `role.stop` and `role.start` to `REQUIRED_METHODS` and the hold to `agents.go`;
- T-1751 adds the two methods to `dashboard-host-channel.md`.

Each only brings the Work in line with touches you have already set.

The 49.26 USD a week stays hand-worked until the operator writes 4h on E-0024 (TH-0567). Resolving.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:08:26Z orchestrator
Resolved: Plan approved with the three Work rewrites that align the tasks with their touches.
