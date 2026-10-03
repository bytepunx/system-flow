---
id: TH-0081
title: Which baseline conventions should the planner, orchestrator, and analyzer read?
anchor:
  path: wip/kanban/stories/S-0207-the-conventions-and-prime-know-the-planner-orchestrator-and-analyzer-as-roles.md
  item: S-0207
status: resolved
participants: [agent-S-0207, alex]
created: 2026-10-03T07:07:21Z
updated: 2026-10-03T07:21:49Z
---

# TH-0081 Which baseline conventions should the planner, orchestrator, and analyzer read?

On wip/kanban/stories/S-0207-the-conventions-and-prime-know-the-planner-orchestrator-and-analyzer-as-roles.md.

## Entries

### 2026-10-03T07:07:21Z agent-S-0207
S-0207 adds the roles `plan`, `orchestrate`, and `analyze` to conventions' `roles` (ADR-0068: no roles means every agent reads the file). I need your confirmation of the baseline's roles before I write them, here and in the template.

Recommended: change only the files the story's goal names, and add a new `strategic-agents.md` for the three:

| File | Today | Proposed |
|------|-------|----------|
| session-start, git | [story] | unchanged |
| communication, safety, documentation, tooling, README | none (all) | unchanged: all three strategic agents read them |
| work-management | [story] | [story, plan, orchestrate]: the planner sizes stories and writes definitions of ready; the orchestrator pulls, orders, and transitions |
| decisions | [story] | [story, orchestrate] |
| continuous-improvement | [story] | [story, analyze] |
| code-quality | [story, verify] | [story, verify, analyze]: the analyzer judges technical risk against it |
| logging | [story, verify] | unchanged |
| telemetry | none (all) | [story, explore, verify, analyze]: the planner and orchestrator never touch a service |
| delegation | [story, explore, verify] | unchanged |
| strategic-agents (new) | none | [plan, orchestrate, analyze] |

Two alternatives, if you want them:
(a) Follow the goal to the letter: work-management gets `orchestrate` only, not `plan`, and code-quality and telemetry stay as they are.
(b) Leave the planner out of work-management and repeat the rules it needs on sizing and ready in strategic-agents.md.

On naming: S-0196 made `flai check` accept the nouns `orchestrator`, `planner`, and `analyzer`, but no file uses them yet. I'll follow this story and use the verbs `plan`, `orchestrate`, and `analyze`, which match `explore`, `verify`, and the host actions, and I'll drop the nouns. Tell me if you'd rather keep the nouns.

### 2026-10-03T07:14:33Z alex
use recommendation

### 2026-10-03T07:21:49Z agent-S-0207
Resolved: Applied the recommended roles, here and in the template (template 1.0.38), on story/S-0207 in 82d2cfd. The roles are the verbs plan, orchestrate, and analyze.
