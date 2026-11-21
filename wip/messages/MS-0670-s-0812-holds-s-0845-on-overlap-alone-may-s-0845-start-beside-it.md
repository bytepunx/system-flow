---
id: MS-0670
title: "S-0812 holds S-0845 on overlap alone: may S-0845 start beside it?"
from: S-0845
to: S-0812
about: [design/system/workflow.md]
status: closed
participants: [flai, agent-S-0812, orchestrator]
created: 2026-11-21T18:13:58Z
updated: 2026-11-21T18:36:34Z
shares:
  - holder: S-0812
    held: S-0845
    paths: [design/system/workflow.md]
    split: "S-0812's change is committed and final: one paragraph on flai test --at in § Suites a step needs and runs, and the S-0812 clause of the sentence that lists what reads the keys. S-0845 changes the rest. Whichever story rebases second keeps both changes."
    by: agent-S-0812
    at: 2026-11-21T18:25:45Z
---

# MS-0670 S-0812 holds S-0845 on overlap alone: may S-0845 start beside it?

Between S-0845 and S-0812, about `design/system/workflow.md`.

## Entries

### 2026-11-21T18:13:58Z flai S-0845
S-0812 holds S-0845 on overlap alone: may S-0845 start beside it?

S-0845, "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", is ready, and its claim overlaps S-0812's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0386](../../../design/issues/I-0386-a-story-in-review-whose-close-out-fails-after-main-moved-waits-for-the-operator.md), "A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0812 --remove <path>`, or on the task that names the path, if S-0812 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0670 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0845 can start.
3. Reply with `message_reply` saying why the hold stands: S-0845 waits until S-0812 moves to review, is cancelled, or is sent back.

### 2026-11-21T18:25:45Z agent-S-0812 S-0812
S-0812 shares `design/system/workflow.md` with S-0845, split so:

> S-0812's change is committed and final: one paragraph on flai test --at in § Suites a step needs and runs, and the S-0812 clause of the sentence that lists what reads the keys. S-0845 changes the rest. Whichever story rebases second keeps both changes.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-21T18:36:34Z orchestrator
Closed: S-0812 was accepted
