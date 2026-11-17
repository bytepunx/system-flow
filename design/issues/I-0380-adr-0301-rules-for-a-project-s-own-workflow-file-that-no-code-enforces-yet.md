---
id: I-0380
title: "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds"
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-15T23:00:06Z
last_reported: 2026-11-15T23:00:06Z
updated: 2026-11-17T08:32:55Z
---

# I-0380 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds

## Description
ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds

## Instances

### 2026-11-15T23:00:06Z
Story: S-0755.
S-0755 made flai read a project's own design/workflows/board.md (ADR-0307). Its docs task found four rules of ADR-0301 that no code applies to such a file. (1) §7 says a story whose branch changes the file is accepted by the operator alone, never by the orchestrator or the merge queue. Today the orchestrator's acceptance block covers only the paths Claude Code protects (ADR-0106). (2) §7 says the dashboard's document and file saves refuse the file on main. No code does. (3) §6 names the check board.wip-state, a limit in wip/kanban/board.md on a state that takes none. No Go code implements it. (4) workitem.LimitedColumns is fixed to ready, in-progress, and review. Board.SetLimit refuses a limit on any other state, and Board.Save drops one written by hand, so a state the project adds with limit: true can get no limit. The runbook docs/operators/runbooks/change-the-board.md states (1) and (2) as rules and (4) as a present limitation. S-0756, the epic's end-to-end story, will meet (1) and (4).

## Remediation

Story S-0844 remediates this issue, created from it at 2026-11-17T08:32:55Z.
