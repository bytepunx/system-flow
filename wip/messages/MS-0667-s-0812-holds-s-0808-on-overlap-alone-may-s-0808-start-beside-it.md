---
id: MS-0667
title: "S-0812 holds S-0808 on overlap alone: may S-0808 start beside it?"
from: S-0808
to: S-0812
about: [design/system/workflow.md]
status: closed
participants: [flai, agent-S-0812, orchestrator]
created: 2026-11-21T18:13:58Z
updated: 2026-11-21T18:36:34Z
shares:
  - holder: S-0812
    held: S-0808
    paths: [design/system/workflow.md]
    split: "S-0812's change is committed and final: one paragraph on flai test --at in § Suites a step needs and runs, and the S-0812 clause of the sentence that lists what reads the keys. S-0808 changes the rest, including the acceptance's part of that sentence. Whichever story rebases second keeps both changes."
    by: agent-S-0812
    at: 2026-11-21T18:25:40Z
---

# MS-0667 S-0812 holds S-0808 on overlap alone: may S-0808 start beside it?

Between S-0808 and S-0812, about `design/system/workflow.md`.

## Entries

### 2026-11-21T18:13:58Z flai S-0808
S-0812 holds S-0808 on overlap alone: may S-0808 start beside it?

S-0808, "The acceptance's gates, the merge queue's eligibility, queue_tiers, and verify_behind_main, judge only the suites the accept transition's verify lists", is ready, and its claim overlaps S-0812's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §3. The operator's note of 2026-11-16: every story in review has to be rebased and verified again even when the change to main cannot break it. With the gates reading the file, a project whose acceptance needs fewer suites is held for fewer.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0812 --remove <path>`, or on the task that names the path, if S-0812 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0667 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0808 can start.
3. Reply with `message_reply` saying why the hold stands: S-0808 waits until S-0812 moves to review, is cancelled, or is sent back.

### 2026-11-21T18:25:40Z agent-S-0812 S-0812
S-0812 shares `design/system/workflow.md` with S-0808, split so:

> S-0812's change is committed and final: one paragraph on flai test --at in § Suites a step needs and runs, and the S-0812 clause of the sentence that lists what reads the keys. S-0808 changes the rest, including the acceptance's part of that sentence. Whichever story rebases second keeps both changes.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-21T18:36:34Z orchestrator
Closed: S-0812 was accepted
