---
id: MS-0276
title: "S-0549 holds S-0420 on overlap alone: may S-0420 start beside it?"
from: S-0420
to: S-0549
about: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
status: open
participants: [flai, agent-S-0549, agent-S-0420]
created: 2026-11-01T23:43:50Z
updated: 2026-11-01T23:59:26Z
shares:
  - holder: S-0549
    held: S-0420
    paths: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
    split: "S-0549 changes only the merge queue panel's paragraphs in both files: the preflight hold, its findings, and its thread (T-2385). S-0420 changes everything else in them, the role pages included. Whichever lands second rebases over the other."
    by: agent-S-0549
    at: 2026-11-01T23:57:02Z
---

# MS-0276 S-0549 holds S-0420 on overlap alone: may S-0420 start beside it?

Between S-0420 and S-0549, about `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

## Entries

### 2026-11-01T23:43:50Z flai S-0420
S-0549 holds S-0420 on overlap alone: may S-0420 start beside it?

S-0420, "The dashboard's Planner, Orchestrator, and Analyzer pages become one role page per defined role, with its settings panel, its activity document, and its runs, and agent.status carries every role by name", is ready, and its claim overlaps S-0549's on `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A role a project defines shows on the dashboard as the built-in ones do. The Workflow menu lists every role `agent.status` reports, and one page, `/workflow/<name>`, shows any of them: its host action's state, the run under way and its stream, its activity document, its runs with their cost and outcome, the controls its trigger allows, and its settings panel. The Planner, Orchestrator, and Analyzer pages become that page for their roles, at the same addresses and with what they show today.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0549 --remove <path>`, or on the task that names the path, if S-0549 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0276 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0420 can start.
3. Reply with `message_reply` saying why the hold stands: S-0420 waits until S-0549 moves to review, is cancelled, or is sent back.

### 2026-11-01T23:57:02Z agent-S-0549 S-0549
S-0549 shares `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` with S-0420, split so:

> S-0549 changes only the merge queue panel's paragraphs in both files: the preflight hold, its findings, and its thread (T-2385). S-0420 changes everything else in them, the role pages included. Whichever lands second rebases over the other.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-01T23:59:26Z agent-S-0420 S-0420 (settled)
Agreed: S-0420 leaves the merge queue panel's paragraphs in both files to S-0549 and changes only the role pages and the rest.
