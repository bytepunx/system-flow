---
id: MS-0669
title: "S-0812 holds S-0837 on overlap alone: may S-0837 start beside it?"
from: S-0837
to: S-0812
about: [design/system/workflow.md, docs/operators/settings.md]
status: closed
participants: [flai, agent-S-0812, orchestrator]
created: 2026-11-21T18:13:58Z
updated: 2026-11-21T18:36:34Z
shares:
  - holder: S-0812
    held: S-0837
    paths: [design/system/workflow.md, docs/operators/settings.md]
    split: S-0812's change is committed and final. In workflow.md it is one paragraph on flai test --at and the S-0812 clause of the sentence that lists what reads the keys. In settings.md it is the generated flag row for flai test --at. S-0837 changes the rest. settings.md's flag table is regenerated with scripts/flai-reference.sh after a rebase. Whichever story rebases second keeps both changes.
    by: agent-S-0812
    at: 2026-11-21T18:25:44Z
---

# MS-0669 S-0812 holds S-0837 on overlap alone: may S-0837 start beside it?

Between S-0837 and S-0812, about `design/system/workflow.md`, `docs/operators/settings.md`.

## Entries

### 2026-11-21T18:13:58Z flai S-0837
S-0812 holds S-0837 on overlap alone: may S-0837 start beside it?

S-0837, "A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired", is ready, and its claim overlaps S-0812's on `design/system/workflow.md`, `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0316](../../../design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md): a publish the orchestrator starts makes no story. Today each one runs the baseline `publish` formula (ADR-0240), which made 51 stories in this repository between 2026-11-07 and 2026-11-17, 131 USD and 2.5 hours of agent time, each taking an in-progress slot, a close-out, a review slot, and an acceptance though its branch changes nothing. The operator asked on 2026-11-17 that this stop. Accept the ADR, or record what the operator changes in it, as this story's first task.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0812 --remove <path>`, or on the task that names the path, if S-0812 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0669 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0837 can start.
3. Reply with `message_reply` saying why the hold stands: S-0837 waits until S-0812 moves to review, is cancelled, or is sent back.

### 2026-11-21T18:25:44Z agent-S-0812 S-0812
S-0812 shares `design/system/workflow.md`, `docs/operators/settings.md` with S-0837, split so:

> S-0812's change is committed and final. In workflow.md it is one paragraph on flai test --at and the S-0812 clause of the sentence that lists what reads the keys. In settings.md it is the generated flag row for flai test --at. S-0837 changes the rest. settings.md's flag table is regenerated with scripts/flai-reference.sh after a rebase. Whichever story rebases second keeps both changes.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-21T18:36:34Z orchestrator
Closed: S-0812 was accepted
