---
id: MS-0662
title: "S-0811 holds S-0808 on overlap alone: may S-0808 start beside it?"
from: S-0808
to: S-0811
about: [design/system/workflow.md]
status: closed
participants: [flai, agent-S-0811, alex]
created: 2026-11-21T18:13:46Z
updated: 2026-11-27T22:47:30Z
---

# MS-0662 S-0811 holds S-0808 on overlap alone: may S-0808 start beside it?

Between S-0808 and S-0811, about `design/system/workflow.md`.

## Entries

### 2026-11-21T18:13:46Z flai S-0808
S-0811 holds S-0808 on overlap alone: may S-0808 start beside it?

S-0808, "The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists", is ready, and its claim overlaps S-0811's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §3. The operator's note of 2026-11-16: every story in review has to be rebased and verified again even when the change to main cannot break it. With the gates reading the file, a project whose acceptance needs fewer suites is held for fewer.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0811 --remove <path>`, or on the task that names the path, if S-0811 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0662 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0808 can start.
3. Reply with `message_reply` saying why the hold stands: S-0808 waits until S-0811 moves to review, is cancelled, or is sent back.

### 2026-11-21T18:15:10Z agent-S-0811 S-0811 (settled)
The hold stands, briefly. S-0811 is rebased onto main and its close-out job is running. It moves to review once that passes, which lifts the hold. S-0811's text in `design/system/workflow.md` § Suites a step needs and runs lists the acceptance's gates (S-0808) as still to come. When S-0808 lands, move them to the "done" list there.

### 2026-11-27T22:47:30Z alex
Closed: S-0811 was accepted
