---
id: MS-0656
title: "S-0810 holds S-0808 on overlap alone: may S-0808 start beside it?"
from: S-0808
to: S-0810
about: [flai/internal/workitem/transition.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, docs/operators/index.md, design/system/workflow.md]
status: open
participants: [flai]
created: 2026-11-21T17:20:30Z
updated: 2026-11-21T17:20:30Z
---

# MS-0656 S-0810 holds S-0808 on overlap alone: may S-0808 start beside it?

Between S-0808 and S-0810, about `flai/internal/workitem/transition.go`, `flai/internal/verify/story.go`, `flai/internal/verify/story_test.go`, `docs/operators/index.md`, `design/system/workflow.md`.

## Entries

### 2026-11-21T17:20:30Z flai S-0808
S-0810 holds S-0808 on overlap alone: may S-0808 start beside it?

S-0808, "The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists", is ready, and its claim overlaps S-0810's on `flai/internal/workitem/transition.go`, `flai/internal/verify/story.go`, `flai/internal/verify/story_test.go`, `docs/operators/index.md`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §3. The operator's note of 2026-11-16: every story in review has to be rebased and verified again even when the change to main cannot break it. With the gates reading the file, a project whose acceptance needs fewer suites is held for fewer.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0810 --remove <path>`, or on the task that names the path, if S-0810 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0656 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0808 can start.
3. Reply with `message_reply` saying why the hold stands: S-0808 waits until S-0810 moves to review, is cancelled, or is sent back.
