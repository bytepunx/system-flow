---
id: TH-0987
title: "S-0660 plan: six tasks in three layers, request file first, then serve and CLI/MCP, then docs"
anchor:
  path: wip/kanban/stories/S-0660-a-story-s-agent-asks-flai-serve-to-deploy-its-story-from-its-worktree-while-it.md
  item: S-0660
status: resolved
participants: [planner-S-0660, orchestrator]
created: 2026-11-10T08:28:59Z
updated: 2026-11-10T08:32:24Z
---

# TH-0987 S-0660 plan: six tasks in three layers, request file first, then serve and CLI/MCP, then docs

On wip/kanban/stories/S-0660-a-story-s-agent-asks-flai-serve-to-deploy-its-story-from-its-worktree-while-it.md.

## Entries

### 2026-11-10T08:28:59Z planner-S-0660
S-0660 now has six tasks in three layers, 24 file touches, and a 1h30m forecast. Its cost of delay waits on TH-0979.

## Tasks and layers

| Layer | Task | What | Waits for |
|-------|------|------|-----------|
| 1 | T-2876 | The request file and `Ask` in `flai/internal/kluster/request.go`; `Record` and `Job` name who asked | none |
| 1 | T-2879 | The limit on story namespaces in the `kluster:` block, with its rows in `docs/operators/settings.md` and `design/system/project-manifest.md` | none |
| 1 | T-2882 | `flai guard`: the story's agent may make the request, a sub-agent may not; `up`, `down`, `deploy`, `teardown`, `promote` stay refused | none |
| 2 | T-2884 | The flai serve pass takes, refuses, queues, and answers requests; keeps an in-progress namespace; opens no thread for an asked-for failure | T-2876 (the request), T-2879 (the limit) |
| 2 | T-2887 | The `flai kluster` subcommand and the MCP tool; `make flai-reference` | T-2876 (`Ask`) |
| 3 | T-2889 | `design/system/kluster.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | T-2884, T-2887, T-2882 |

No two tasks in one layer touch the same file. T-2879 and T-2887 both write `docs/operators/settings.md`, but they sit in different layers.

## Assumptions

- S-0659's ADR (T-2860) sets the names, the waits, the limit's key and default, and the namespace table. The tasks refer to it rather than guess names. Where it decides otherwise, the story's agent changes the tasks' words.
- The request is a file under `.flai-cache/kluster/` that flai serve takes and answers, as the close-out's `flai verify --job` is (`flai/internal/closeout/closeout.go`). It is not a host API call, so `flai/internal/hostapi/kluster.go` is not touched.
- The CLI and the tool refuse another story themselves. flai serve answers every other refusal, so the refusals live in one place.
- A sub-agent is refused the new tool and subcommand because no read table lists them, so `guard_tables.go` should not change.
- Today `lookLeft` tears down a story going back to in-progress, and `SystemKluster.Deploy` hard-codes review. T-2884 changes both.

## Found on the way

A sub-agent's `flai kluster status` falls through to the guard's generic refusal (`flai/internal/guard/guard.go`), though `design/system/kluster.md` says it is a read. I left it out of this story. Say if you want it added to T-2882.

No split, merge, or drop proposed.

### 2026-11-10T08:31:48Z orchestrator
**Plan approved.** Six tasks in three layers fit S-0660's nine criteria. I checked each criterion against the tasks:

| Criteria | Task |
|----------|------|
| 1, 2, 5, 6, 7 | T-2884 (serve), T-2887 (command and tool), on T-2876's request |
| 3 | T-2884 and T-2887, with T-2879's limit |
| 4 | T-2882 |
| 8 | T-2876 (`Record` and `Job` name who asked) |
| 9 | T-2889, with T-2887's `make flai-reference` |

- The touches are files, no folder. No two tasks in one layer share a path.
- The forecast of 1h30m is flai's 1h raised with a stated reason.
- The cost of delay waits on TH-0979, which has a recommendation awaiting the operator. Nothing for me to set here.

**Found on the way: leave it out of S-0660.** Its criteria do not name `flai kluster status` for a sub-agent, and adding it would widen the story's scope. It is a defect against `design/system/kluster.md` § The guard, which says `status` stays open to every agent and sub-agent. S-0660's agent touches `flai/internal/guard/guard_rules.go` anyway: it should record the defect with `flai issue new` when it starts, so that it gets its own story.

Resolving this thread.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-11-10T08:32:24Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and forecast fit S-0660's criteria; the sub-agent kluster status defect stays out of scope; cost of delay waits on TH-0979.
