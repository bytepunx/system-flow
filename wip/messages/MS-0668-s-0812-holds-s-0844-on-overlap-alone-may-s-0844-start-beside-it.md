---
id: MS-0668
title: "S-0812 holds S-0844 on overlap alone: may S-0844 start beside it?"
from: S-0844
to: S-0812
about: [docs/operators/runbooks/change-the-board.md, design/system/workflow.md]
status: closed
participants: [flai, agent-S-0812, orchestrator]
created: 2026-11-21T18:13:58Z
updated: 2026-11-21T18:36:34Z
shares:
  - holder: S-0812
    held: S-0844
    paths: [design/system/workflow.md, docs/operators/runbooks/change-the-board.md]
    split: "S-0812's change is committed and final. In workflow.md it is one paragraph on flai test --at and the S-0812 clause of the sentence that lists what reads the keys. In the runbook it is the batch and flai test --at clauses of the \"Nothing else reads the keys yet\" sentence in § Name the suites a step runs. S-0844 changes the rest. Whichever story rebases second keeps both changes."
    by: agent-S-0812
    at: 2026-11-21T18:25:42Z
---

# MS-0668 S-0812 holds S-0844 on overlap alone: may S-0844 start beside it?

Between S-0844 and S-0812, about `docs/operators/runbooks/change-the-board.md`, `design/system/workflow.md`.

## Entries

### 2026-11-21T18:13:58Z flai S-0844
S-0812 holds S-0844 on overlap alone: may S-0844 start beside it?

S-0844, "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds", is ready, and its claim overlaps S-0812's on `docs/operators/runbooks/change-the-board.md`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0380](../../../design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md), "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0812 --remove <path>`, or on the task that names the path, if S-0812 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0668 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0844 can start.
3. Reply with `message_reply` saying why the hold stands: S-0844 waits until S-0812 moves to review, is cancelled, or is sent back.

### 2026-11-21T18:25:42Z agent-S-0812 S-0812
S-0812 shares `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md` with S-0844, split so:

> S-0812's change is committed and final. In workflow.md it is one paragraph on flai test --at and the S-0812 clause of the sentence that lists what reads the keys. In the runbook it is the batch and flai test --at clauses of the "Nothing else reads the keys yet" sentence in § Name the suites a step runs. S-0844 changes the rest. Whichever story rebases second keeps both changes.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-21T18:36:34Z orchestrator
Closed: S-0812 was accepted
