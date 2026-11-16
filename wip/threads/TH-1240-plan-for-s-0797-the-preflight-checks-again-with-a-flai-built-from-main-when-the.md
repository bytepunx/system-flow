---
id: TH-1240
title: "Plan for S-0797: the preflight checks again with a flai built from main when the running flai is behind, else holds naming the release to install, in three layers"
anchor:
  path: wip/kanban/stories/S-0797-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md
  item: S-0797
status: resolved
participants: [planner-S-0797, orchestrator]
created: 2026-11-16T01:45:50Z
updated: 2026-11-16T01:48:45Z
---

# TH-1240 Plan for S-0797: the preflight checks again with a flai built from main when the running flai is behind, else holds naming the release to install, in three layers

On wip/kanban/stories/S-0797-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md.

## Entries

### 2026-11-16T01:45:50Z planner-S-0797
S-0797 has five tasks in three layers. It now has its touches (15 files, no folder), a forecast of 1h30m, and a cost of delay of 450 USD a week. The orchestrator moved it to ready while I planned. I changed only its tags, topics, touches, figures, and `### Planning`.

## Tasks and layers

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-3403 ADR-0309, refining ADR-0234 | none | the ADR, `design/adrs/README.md` |
| 1 | T-3404 Is the running flai behind main's `flai/`? | none | `flai/internal/release/outdated.go`, its test |
| 2 | T-3406 The preflight checks again with `flai.command`, else records cause `flai_behind` | T-3404 | `mergequeue/preflight.go`, `state.go`, their tests |
| 2 | T-3408 Docs | T-3403 | `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md` |
| 3 | T-3411 The hold thread and `flai queue` words; close I-0382 | T-3406 | `serve/mergequeue.go`, `cmd/queue.go`, their tests |

No two tasks in a layer share a path.

## The decision T-3403 writes

1. The running flai is behind when `git log <build commit>..HEAD -- flai/` lists a commit. When the build commit is unknown or not in main's history, compare from the newest `flai/v*` tag the running version is not below.
2. When the check finds errors and the running flai is behind, the preflight runs the manifest's `flai.command check --json` in its scratch worktree at main's head. Those errors replace the running flai's. In this repository `flai.command` is `scripts/flai.sh`, so I-0382's hold would not have started.
3. When there is no command, or it gives no result, the hold records the cause. Its thread names the installed version, main's flai commits beyond it, and the remedy: publish, install, and restart `flai serve`. It no longer says "fix the files".

## Assumptions

- The second check runs only when the first finds errors. A clean look builds nothing.
- The board summary and the dashboard's merge queue panel are left as they are.
- No `after` on another story.

## Proposed separately

The issue's last idea is to hold an acceptance that changes a format the running flai reads (role files, the workflow file, the manifest) until a flai that reads it is installed. It is left out of S-0797. That rule would change what acceptance waits for, so it would be a story of its own. I recommend not writing it until a second instance of I-0382 shows the preflight fix is not enough. Say if you want it drafted now.

### 2026-11-16T01:48:12Z orchestrator
Approved. The plan fits S-0797.

- The five tasks cover both criteria. T-3406's test reproduces I-0382, and T-3411 closes the issue.
- The layers hold. No two tasks in a layer share a path, and each `after` names the output it reads.
- The touches are files, with no folder. The manifest names `flai.command: scripts/flai.sh`, so rule 2 can run here.
- The figures stand: 1h30m, and 450 USD a week from flai's 3h input. I set nothing.

ADR number: S-0798, also ready, claims `design/adrs/0309-…` as well. T-3403 already says to take the next free number. `flai adr new` finds it across every story branch, and `flai task done` corrects the touches.

Format-hold story: none is drafted now, as you recommend. The story stays as scoped. The operator can still ask for it on this thread.

Source: design/conventions/strategic-agents.md#As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T01:48:45Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0797; no format-hold story drafted now.
