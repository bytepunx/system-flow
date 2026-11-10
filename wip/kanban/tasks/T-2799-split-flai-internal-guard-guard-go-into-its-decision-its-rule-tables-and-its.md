---
id: T-2799
type: task
nature: improvement
title: Split flai/internal/guard/guard.go into its decision, its rule tables, and its shell parsing
status: done
parent: S-0631
owner: alex
created: 2026-11-09T22:05:02Z
updated: 2026-11-10T06:38:54Z
transitions:
  - to: ready
    at: 2026-11-10T06:16:50Z
    by: agent-S-0631
  - to: in-progress
    at: 2026-11-10T06:16:50Z
    by: agent-S-0631
  - to: done
    at: 2026-11-10T06:38:54Z
    by: agent-S-0631
stream: S-0631
tags: [cli]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_tables.go, flai/internal/guard/guard_rules.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/guard_shell.go]
usage:
  source: log
  seconds: 1279
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 19
      output: 68
      cache_read: 2741452
      cache_write: 18883
      cost: 1.0151
---
# T-2799 Split flai/internal/guard/guard.go into its decision, its rule tables, and its shell parsing

## Work

`guard.go` is 2,051 lines. I-0298 does not name it, but it has gone over the limit since the analysis, and the size test this story adds would fail on it. Move its declarations into new files of package `guard` by responsibility. What the guard allows and refuses stays the same, and so do its refusal texts. The line ranges below are approximate, from the plan's read on 2026-11-09:

| Lines | Declarations | File |
|-------|--------------|------|
| 220–400 | the read lists and maps: `MCPReads`, `MCPStoryReads`, `cliReads`, `flagReads`, `roleReads`, `readsRoles`, `HostAPIReads`, `readsHostAPI`, `hostAPIMethod`, `writesHostAPI`, `gitReads` | `guard_tables.go` |
| 401–816 | each role's lists, texts, and `rules`: the planner, orchestrator, analyzer, curator, forge, and kluster blocks | `guard_rules.go` |
| 817–1416 | `Guard`, `Hook`, `DecideAll`, `Reader` and `ReaderFor`, `rules`, `subAgent`, `planning`, the orchestrator's and planner's predicates, `ProtectedWrite` | `guard.go`, which keeps them |
| 1420–1786 | `waitsOnSubAgents`, the refusal helpers (`refusedOrchestrator`, `off`, `nevers`, `resolved`), `subcommands`, `positionals` | `guard_refusals.go` |
| 1787–2051 | `shells`, `heredoc`, `commands`, `arguments`, `simple`, `parse`, `operator`, `descriptor`, `opens`, `bodies`, `runsShell` | `guard_shell.go` |

The package already has `roles.go`, `authors.go`, `subagents.go`, `shared.go`, and `writes.go`. Where a block belongs in one of those, move it there instead, keeping every file under 1,200.

Change no test file.

This task waits for nothing and shares no path with the other splits, so all of them run together.

## Done when

- No non-test file in `flai/internal/guard` is over 2,000 lines, and none of the new files is over 1,200.
- `git diff --stat` shows no `_test.go` file changed.
- `flai test flai/internal/guard` passes, gofmt, goimports, and golangci-lint included.

## Notes
