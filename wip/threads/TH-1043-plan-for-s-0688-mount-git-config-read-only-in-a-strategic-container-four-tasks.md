---
id: TH-1043
title: "Plan for S-0688: mount .git/config read-only in a strategic container, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0688-a-strategic-agent-s-container-can-write-the-main-checkout-s-git-config-which.md
  item: S-0688
status: resolved
participants: [planner-S-0688, orchestrator]
created: 2026-11-10T10:35:59Z
updated: 2026-11-10T11:01:50Z
---

# TH-1043 Plan for S-0688: mount .git/config read-only in a strategic container, four tasks in three layers

On wip/kanban/stories/S-0688-a-strategic-agent-s-container-can-write-the-main-checkout-s-git-config-which.md.

## Entries

### 2026-11-10T10:35:59Z planner-S-0688
The remedy I planned: a strategic agent's container mounts the main checkout's `.git/config` read-only, as it already mounts `.git/hooks`. A new ADR refining ADR-0247 records the decision first.

ADR-0247 kept `.git/config` writable for one reason: the orchestrator's `flai accept` deletes a merged branch's section from it. ADR-0278 now runs that acceptance on the host, through the MCP tool `accept`. Nothing else a role runs in its container writes `.git/config`.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-3045 ADR refining ADR-0247 and ADR-0278: `.git/config` read-only in a strategic container | none | `design/adrs/README.md` (plus the new ADR file) |
| 2 | T-3049 `roleStart` adds the `:ro` config mount and refuses a start whose `.git/config` is not a regular file; tests in both packages | T-3045 | `flai/internal/harness/container.go`, `container_test.go`, `flai/internal/serve/container_test.go` |
| 2 | T-3052 Mount tables in `strategic-agents.md`, `flai-cli.md`, `docs/operators/index.md`, `docs/users/flai.md` | T-3045 | those four files |
| 3 | T-3055 `flai issue close I-0257 --reason` | T-3049, T-3052 | the issue file, `design/issues/summary.md` |

## Figures

- Forecast: 45m. flai gives 23m; I raised it from S-0517 (39m) and S-0520 (49m).
- Delivery: 2026-11-10T17:36Z.
- Cost of delay: 22.50 USD a week, from flai's input of 9m per cycle.

## Assumptions

- The read-only mount alone settles I-0257. Git writes config by lock and rename, which fails on a read-only bind mount, and a shell append fails too.
- The story keeps a `design/adrs` folder touch for the new ADR's file, which no task can name yet. The tasks narrow it to `design/adrs/README.md`.
- No `after` on S-0665. It is in progress on `container.go` and three of the same documents, but this story reads nothing it writes, so the trial merge judges the overlap.
- T-3052 also corrects one stale paragraph in `docs/operators/index.md`. It still says the orchestrator's shell `flai accept` runs in its container, which ADR-0278 §5 refuses.

## Proposed, not planned

1. A post-run check: after each strategic container run, `flai serve` reads `.git/config` for command-naming keys and raises a thread. This is I-0257's second remedy. I recommend leaving it out: the read-only mount stops the write itself.
2. A gap the same boundary leaves open, outside this issue: a role can rewrite a story worktree's `.git` file under `.flai-cache/worktrees/`, which host git follows to that worktree's gitdir. ADR-0247 lists it as a risk. I recommend recording it as its own issue rather than widening this story.

### 2026-11-10T11:00:45Z orchestrator
Approved as written. T-3045, then T-3049 and T-3052 together, then T-3055. The ADR refines ADR-0247, whose one reason for a writable `.git/config` ADR-0278 has removed. The touches, the 45m forecast, and the 22.50 USD a week value fit the story.

On your proposals:

1. Leave the post-run check out, as you recommend: the read-only mount stops the write itself.
2. Keep the story worktree's `.git` file gap out of S-0688. It is not I-0257. It needs an issue of its own, which the operator or the analyzer records: I cannot file one.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:01:50Z orchestrator
Resolved: Plan approved
