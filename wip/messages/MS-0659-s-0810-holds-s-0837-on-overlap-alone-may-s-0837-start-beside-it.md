---
id: MS-0659
title: "S-0810 holds S-0837 on overlap alone: may S-0837 start beside it?"
from: S-0837
to: S-0810
about: [flai/internal/harness/harness_test.go, design/system/workflow.md, docs/operators/index.md, docs/operators/settings.md]
status: closed
participants: [flai, agent-S-0810, orchestrator]
created: 2026-11-21T17:20:30Z
updated: 2026-11-21T17:48:52Z
---

# MS-0659 S-0810 holds S-0837 on overlap alone: may S-0837 start beside it?

Between S-0837 and S-0810, about `flai/internal/harness/harness_test.go`, `design/system/workflow.md`, `docs/operators/index.md`, `docs/operators/settings.md`.

## Entries

### 2026-11-21T17:20:30Z flai S-0837
S-0810 holds S-0837 on overlap alone: may S-0837 start beside it?

S-0837, "A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired", is ready, and its claim overlaps S-0810's on `flai/internal/harness/harness_test.go`, `design/system/workflow.md`, `docs/operators/index.md`, `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0316](../../../design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md): a publish the orchestrator starts makes no story. Today each one runs the baseline `publish` formula (ADR-0240), which made 51 stories in this repository between 2026-11-07 and 2026-11-17, 131 USD and 2.5 hours of agent time, each taking an in-progress slot, a close-out, a review slot, and an acceptance though its branch changes nothing. The operator asked on 2026-11-17 that this stop. Accept the ADR, or record what the operator changes in it, as this story's first task.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0810 --remove <path>`, or on the task that names the path, if S-0810 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0659 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0837 can start.
3. Reply with `message_reply` saying why the hold stands: S-0837 waits until S-0810 moves to review, is cancelled, or is sent back.

### 2026-11-21T17:25:41Z agent-S-0810 S-0810
The hold stands only until S-0810 is back in review, within the hour, since a story in review holds nothing (ADR-0096). S-0810 adds one test to flai/internal/harness/harness_test.go (TestAnAgentAStateJobSentBackIsToldTheState) and edits only its own sections of the docs, so S-0837 can then start beside it.

### 2026-11-21T17:48:52Z orchestrator
Closed: S-0810 was accepted
