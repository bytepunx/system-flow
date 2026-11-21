---
id: MS-0650
title: "S-0809 holds S-0808 on overlap alone: may S-0808 start beside it?"
from: S-0808
to: S-0809
about: [flai/internal/verify/story.go, flai/internal/verify/story_test.go, docs/operators/index.md, design/system/workflow.md]
status: closed
participants: [flai, agent-S-0809, alex]
created: 2026-11-18T13:04:39Z
updated: 2026-11-21T17:20:21Z
---

# MS-0650 S-0809 holds S-0808 on overlap alone: may S-0808 start beside it?

Between S-0808 and S-0809, about `flai/internal/verify/story.go`, `flai/internal/verify/story_test.go`, `docs/operators/index.md`, `design/system/workflow.md`.

## Entries

### 2026-11-18T13:04:39Z flai S-0808
S-0809 holds S-0808 on overlap alone: may S-0808 start beside it?

S-0808, "The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists", is ready, and its claim overlaps S-0809's on `flai/internal/verify/story.go`, `flai/internal/verify/story_test.go`, `docs/operators/index.md`, `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §3. The operator's note of 2026-11-16: every story in review has to be rebased and verified again even when the change to main cannot break it. With the gates reading the file, a project whose acceptance needs fewer suites is held for fewer.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0809 --remove <path>`, or on the task that names the path, if S-0809 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0650 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0808 can start.
3. Reply with `message_reply` saying why the hold stands: S-0808 waits until S-0809 moves to review, is cancelled, or is sent back.

### 2026-11-18T13:08:23Z agent-S-0809 S-0809 (settled)
S-0809 changes nothing more on these paths. It was sent back only to rebase onto S-0807, which is done; its close-out is running, and it moves back to review once that passes, about ten minutes, where it holds nothing. So the hold ends on its own then; no share is needed.

### 2026-11-21T17:20:21Z alex
Closed: S-0809 was accepted
