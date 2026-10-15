---
id: S-0390
type: story
nature: improvement
title: The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI
status: ready
owner: alex
created: 2026-10-15T02:55:10Z
updated: 2026-10-15T03:00:44Z
transitions:
  - to: ready
    at: 2026-10-15T03:00:41Z
    by: orchestrator
tags: []
topics: [mcp, release, host]
touches: [design/adrs/0159-an-outdated-mcp-server-s-task-done-closes-the-task-with-the-installed-flai-once-it-has-caught-up.md, design/adrs/README.md, design/system/flai-cli.md, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, flai/internal/mcpserver/installed.go, flai/internal/mcpserver/installed_test.go, flai/internal/mcpserver/server.go, docs/users/flai.md, docs/operators/index.md, design/system/agent-narrative.md, design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 287
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 38
          output: 19250
          cache_read: 1065457
          cache_write: 87263
          cost: 0.0312
        - model: claude-opus-5-5
          input: 42
          output: 19352
          cache_read: 1999641
          cache_write: 128383
          cost: 1.8142
    - kind: orchestrator
      seconds: 299
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 568
          cache_read: 10748030
          cache_write: 11932
          cost: 2.6975
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-15T02:55:10Z
  value: 30
  by: planner-S-0390
  at: 2026-10-15T02:59:04Z
forecast:
  duration: 40m
  delivery: 2026-10-15T03:45:00Z
  basis: "Its own forecast of 40m; 2nd in the pull order with an in-progress limit of 5, behind S-0385 and S-0386."
  by: flai
  at: 2026-10-15T03:00:44Z
finalized:
  by: alex
  at: 2026-10-15T02:55:40Z
---
# S-0390 The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI

## Goal

This story remediates [I-0147](../../../design/issues/I-0147-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md), "The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0147 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0147 is closed with `flai issue close I-0147 --reason` saying what fixed it

## Tasks
- T-1595 Record in a new ADR that an outdated MCP server's task_done closes the task with the installed flai once it has caught up
- T-1596 The MCP tool task_done closes a task with the installed flai when the session's flai is a minor version behind and the installed one has caught up
- T-1597 Document how task_done closes a task with the installed flai, for users, operators, and agents
- T-1598 Close I-0147 with the remedy that fixed it

## Notes

Cost of delay inputs set by flai from I-0147. time_lost_per_cycle 12m: 3m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-10-13T21:48:16Z, 1.2 days before this story; under one cycle counts as one).

### Planning

This is the planner's proposed remedy. The story's agent confirms or changes it in T-1595's ADR. The plan's thread on S-0390 has the rest.

- **What S-0387 already fixed.** ADR-0158, published in flai 1.50.1, makes `flai serve` upgrade the host, or ask the operator, within a minute of a release. That shortens the lag but does not end it. The host here still runs 1.48.0, which predates ADR-0158.
- **The cause that remains.** A story's agent's MCP server is a stdio child of its claude session. It runs the binary path `flai serve` gave it (`flai/internal/harness/adapters.go`, `mcpConfig`), and keeps that process for the whole session. After the host upgrades, the session's `task_done` still runs the old flai, and ADR-0142 still refuses it. All four instances were story sessions that started before the release.
- **The remedy.** When ADR-0142 would refuse, the MCP tool `task_done` reads the version of the binary now at its own path. When that version has caught up with the newest release, it closes the task by running that binary's `task done`. Otherwise it refuses, says the installed flai is behind too, and tells the agent to call again once the host is upgraded.
- **Not proposed.** Re-executing the stdio server would lose its MCP session. Restarting the agent's session would cut its work. Widening the refusal to every write tool is broader than the issue.

Where each touch came from:

| Touch | Source |
|-------|--------|
| `design/adrs/0159-…-once-it-has-caught-up.md`, `design/adrs/README.md` | design: the remedy is a decision refining ADR-0142 and ADR-0158; the number and slug are predicted, and `flai task done` widens the touches with the real path |
| `design/system/flai-cli.md` | design: § Versions: the host's flai and the tree |
| `flai/internal/taskdone/taskdone.go`, `taskdone_test.go` | layout: `refuseOutdated`, ADR-0142's check |
| `flai/internal/mcpserver/task.go`, `task_test.go` | layout: the `task_done` handler and `TestTaskDoneRefusesAFlaiAMinorVersionBehind` |
| `flai/internal/mcpserver/installed.go`, `installed_test.go` | layout: a new file to read the installed binary's version and run it; names predicted |
| `flai/internal/mcpserver/server.go` | co-change (22 commits with these paths): the field a test sets for the installed path |
| `docs/users/flai.md`, `docs/operators/index.md`, `design/system/agent-narrative.md` | design: § Closing a task, § Keeping the host's flai current, § What an agent is told through MCP |
| `design/issues/I-0147-…md`, `design/issues/summary.md` | goal: the second criterion closes I-0147 |

No folder touch is kept. Co-change's `docs/users/flai-reference.md` and `docs/operators/settings.md` are left out because no flag, help text, or setting changes.

Figures:

- **Forecast 40m**, raised from flai's 23m (85 s per unit of size times 16). The plan has an ADR, a subprocess path in `task_done` with a test on a fake installed binary, docs, and the issue's close. S-0387's comparable plan took 48m from start to review. Delivery moves by the 17m added, to 2026-10-15T03:41:00Z.
- **Cost of delay 30 USD a week**, flai cod's figure from the 12m time_lost_per_cycle input at 150 USD an hour, kept. The four instances each cost about 3m to work around, and a project made from the template has no tree CLI to fall back to.
