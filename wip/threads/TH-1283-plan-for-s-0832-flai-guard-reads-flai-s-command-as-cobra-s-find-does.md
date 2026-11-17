---
id: TH-1283
title: "Plan for S-0832: flai guard reads flai's command as cobra's Find does"
anchor:
  path: wip/kanban/stories/S-0832-flai-guard-reads-a-subcommand-s-value-flag-before-flai-s-command-as-the-command.md
  item: S-0832
status: resolved
participants: [planner-S-0832, orchestrator]
created: 2026-11-17T08:30:46Z
updated: 2026-11-17T08:33:42Z
---

# TH-1283 Plan for S-0832: flai guard reads flai's command as cobra's Find does

On wip/kanban/stories/S-0832-flai-guard-reads-a-subcommand-s-value-flag-before-flai-s-command-as-the-command.md.

## Entries

### 2026-11-17T08:30:46Z planner-S-0832
S-0832 is planned in three tasks, one per layer.

## Proposed fix

The guard stops listing the flags that take a value before flai's command (`flaiValues`, `authorValues`). It reads the command and subcommand as cobra's `Find` does instead: a flag before them takes the next word unless flai's command tree defines it with no value, or it carries its value after `=`. A flag flai adds later is then covered without anyone updating a list.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3530: `subcommands` and the callers that read the command by position (`names`, `readsRoles`, `hostAPIMethod`, the planner's moves, a role's plans) use the cobra reading, with a built-in no-value table and tests in `guard_test.go` | 7 files under `flai/internal/guard/` | none |
| 2 | T-3532: `flai guard` fills the table from flai's command tree, and a test in `flai/cmd/guard_test.go` holds the guard's reading to the root's `Find` | `flai/cmd/guard.go`, `flai/cmd/guard_test.go` | T-3530 |
| 3 | T-3537: the `flai guard` row of `design/system/flai-cli.md` says how the guard reads the command, and I-0358 is closed | `flai-cli.md`, I-0358's file, `design/issues/summary.md` | T-3532 |

## Figures

- Forecast: 40m, delivered 2026-11-17T16:29:00Z. That is `flai forecast`'s 28m plus 12m, for proposing the fix and for T-3532's test over the whole command tree.
- Cost of delay: 25 USD a week, from the operator's input of 10m per cycle. It stands as given.

## Assumptions

- A second gap from the same cause is in scope. `authorValues` lacks `--wait`, so `flai --wait 5m thread reply TH-0001 --by other` passes the author check today. T-3530 fixes it with the rest.
- Out of scope: `positionals` stays as it is for the words after the command, where pflag parses the found command's own flags. Only the words before the command change.
- No ADR is needed. This changes how the guard parses a command line, not a decision any ADR records. Say so if you want one.

Nothing here needs your answer. I would neither split nor merge the tasks.

### 2026-11-17T08:33:08Z orchestrator
Approved. The plan fits S-0832:

- T-3530 makes the guard read the command as cobra's `Find` does, with tests that reproduce I-0358, which meets criterion 1. T-3532 fills the no-value table from flai's command tree and holds it to `Find`. T-3537 documents it and closes I-0358, which meets criterion 2. Each task waits for the one before because it builds on its reading.
- The `--wait` gap in `authorValues` has the same cause, so fixing it in T-3530 is within the story.
- The change is to how the guard parses a command line, not to a recorded decision, so no ADR is needed.
- Every touch is a file. There is no folder touch and no `after`.
- The 40m forecast is `flai forecast`'s 28m, raised for the proposal and the command-tree test. The 25 USD a week value comes from `flai cod`.

I am finalizing S-0832 under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T08:33:42Z orchestrator
Resolved: Plan approved under plan_backlog_stories.
