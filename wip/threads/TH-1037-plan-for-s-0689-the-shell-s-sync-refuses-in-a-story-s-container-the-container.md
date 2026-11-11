---
id: TH-1037
title: "Plan for S-0689: the shell's sync refuses in a story's container, the container prompt names task_done and verify, and sync's stop and trial-merge reports stop misleading; four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0689-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md
  item: S-0689
status: resolved
participants: [planner-S-0689, alex]
created: 2026-11-10T10:34:32Z
updated: 2026-11-11T08:53:36Z
---

# TH-1037 Plan for S-0689: the shell's sync refuses in a story's container, the container prompt names task_done and verify, and sync's stop and trial-merge reports stop misleading; four tasks in three layers

On wip/kanban/stories/S-0689-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md.

## Entries

### 2026-11-10T10:34:32Z planner-S-0689
The plan for S-0689 is four tasks in three layers. No question is open.

## Proposed solution, from I-0272

- ADR-0246 says `flai stream sync` refuses in a story's container.
- The refusal keys on `MainAbsent`. ADR-0253's mount of `system-flow.yaml` leaves that empty, so the sync runs against the read-only `.git`.
- The fix: the shell's sync refuses there at once, the container prompt names the MCP tools that sync on the host, and sync's reports stop misleading wherever it runs.

## Tasks and layers

| Layer | Task | What | after |
|-------|------|------|-------|
| 1 | T-3026 | `flai stream sync` and the shell's `flai task done` refuse in a story's container. The error names the MCP tools `task_done` and `verify`. The signal is `FLAI_CONFIG` under `.flai-cache/agents` with `FLAI_STORY` set. | none |
| 1 | T-3027 | The container harness's start prompt says to close tasks with `task_done` and close out with `verify`, not the shell's sync. | none |
| 2 | T-3030 | A rebase that stops with no conflict reports git's error, with no "resolve them" and no abort that cannot restore the branch. A trial merge that cannot run is reported as not run. Includes a read-only reproduction test. | T-3026: both change `flai/cmd/stream_sync_test.go` |
| 3 | T-3033 | Updates `docs/users/flai.md` (and drops the `stream_sync` MCP tool it names, which does not exist), `design/system/flai-cli.md`, and the agent-container runbook, then closes I-0272. | T-3026, T-3027, T-3030: it documents what they built |

## Figures

| Figure | Value | Basis |
|--------|-------|-------|
| Touches | 20 files, no folder | Code layout from the explorer, the container docs, and the criteria. Co-change added nothing at 10 percent or less. |
| Forecast | 1h30m, delivery 2026-11-10T18:15Z | Raised from flai's 39m for four tasks over three packages, a read-only git test, and the integration and smoke tiers at close-out. |
| Cost of delay | 30 USD a week | `flai cod` from the operator's input of 12m per cycle. Unchanged. |

## Assumptions

- No `asked` thread exists on S-0689 or I-0272, and the story had no earlier plan. I planned it as a story with no tasks.
- No MCP tool syncs alone. Between tasks, `task_done` and the `verify` close-out are the host's sync for a container story. A `stream_sync` MCP tool, for an agent that must sync after an `overlapped` change with no task to close, is a story I would add apart. I created no draft for it.
- No ADR is needed: the refusal is what ADR-0246 already decided.
- The baseline `git.md` and its template copy stay unchanged. They say nothing of containers.
- The read-only reproduction follows `TestGitWorksReadOnlyOnTheMountsAndRefusesACommit`, which skips under root.

### 2026-11-11T08:53:36Z alex
Resolved.
