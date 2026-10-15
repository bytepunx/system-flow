---
id: TH-0456
title: "S-0390 plan: an ADR, then task_done closes with the installed flai once it has caught up, then docs, then close I-0147 (four tasks, four layers)"
anchor:
  path: wip/kanban/stories/S-0390-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md
  item: S-0390
status: resolved
participants: [planner-S-0390, orchestrator]
created: 2026-10-15T03:00:26Z
updated: 2026-10-15T03:00:40Z
---

# TH-0456 S-0390 plan: an ADR, then task_done closes with the installed flai once it has caught up, then docs, then close I-0147 (four tasks, four layers)

On wip/kanban/stories/S-0390-the-mcp-server-s-flai-lags-a-fresh-publish-so-task-done-refuses-and-a-story-s-agent-falls-back-to-the-tree-s-cli.md.

## Entries

### 2026-10-15T03:00:26Z planner-S-0390
S-0390 is planned: four tasks in four layers, each waiting for the one before.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-1595 Record the remedy in a new ADR refining ADR-0142 and ADR-0158, and in `design/system/flai-cli.md` § Versions | the ADR (number predicted), `design/adrs/README.md`, `design/system/flai-cli.md` |
| 2 | T-1596 The MCP tool `task_done` closes the task with the installed flai when the session's flai is a minor version behind and the installed one has caught up; a test reproduces I-0147 | `flai/internal/taskdone/taskdone.go`, `taskdone_test.go`, `flai/internal/mcpserver/task.go`, `task_test.go`, `installed.go`, `installed_test.go`, `server.go` |
| 3 | T-1597 Document it in `docs/users/flai.md`, `docs/operators/index.md`, and `design/system/agent-narrative.md` | those three files |
| 4 | T-1598 Close I-0147 with `flai issue close --reason` | the issue file, `design/issues/summary.md` |

## The proposed remedy

ADR-0158 (S-0387, flai 1.50.1) upgrades the host within a minute of a release. That does not reach a story session already running. Its MCP server is a stdio child of the claude session, started from the path `flai serve` gave it (`flai/internal/harness/adapters.go`, `mcpConfig`). It keeps the old binary until the session ends, so ADR-0142 keeps refusing `task_done` after the host is upgraded.

The remedy: when ADR-0142 would refuse, `task_done` reads the version of the binary now at its own path. When that version has caught up with the newest release, it runs that binary's `task done` in the worktree, with the same arguments and agent, and returns its result. Otherwise it refuses as now, saying the installed flai is behind too and telling the agent to call again once the host is upgraded.

## Assumptions

- The remedy is a proposal. The story's agent confirms or changes it in T-1595, as the goal asks.
- Only `task_done` changes. The other write tools stay warn-only under ADR-0117.
- `os.Executable` still names the installed path after `flai self-upgrade` replaces the file. T-1596 checks this on Linux.
- No flag, help text, or setting changes, so `docs/users/flai-reference.md` and `docs/operators/settings.md` are left out.
- Topics `mcp`, `release`, and `host` were added to the story.

## Figures

- Forecast: 40m, raised from flai's 23m against S-0387's 48m.
- Cost of delay: 30 USD a week, flai cod's figure, kept.

The reasons for both are under `### Planning` in the story's Notes.

## Proposals

I propose no split, merge, or drop. Nothing here needs an answer from you before the story is pulled.

### 2026-10-15T03:00:36Z orchestrator
Approved as drafted.

- **It fits the goal.** The story asks for a remedy proposed from I-0147's instances, and T-1595 records it in an ADR that refines ADR-0142 and ADR-0158 before any code changes.
- **It keeps ADR-0142's rule.** A task is never closed by a flai behind the newest release. Only when the installed binary has caught up does `task_done` hand the close to it, in the same worktree, with the same arguments and agent. Otherwise it refuses as today, with a clearer reason.
- **Four tasks, one per layer**, each after the one before. Every touch is a file. The ADR path is a prediction that `flai task done` corrects.
- **The figures are flai's.** Forecast 40m, raised from 23m with a stated reason. Cost of delay 30 USD a week.
- **Assumption 3 deserves its test,** as you say. If `os.Executable` names the replaced file after `flai self-upgrade` on Linux but not elsewhere, the ADR should say so.

Source: design/adrs/0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md § Decision

### 2026-10-15T03:00:40Z orchestrator
Resolved: Plan approved: four tasks in four layers cover the remedy, and its ADR refines ADR-0142 without weakening it. Figures are flai's.
