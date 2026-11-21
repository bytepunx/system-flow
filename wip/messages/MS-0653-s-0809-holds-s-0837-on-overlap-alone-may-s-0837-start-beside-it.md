---
id: MS-0653
title: "S-0809 holds S-0837 on overlap alone: may S-0837 start beside it?"
from: S-0837
to: S-0809
about: [design/system/project-manifest.md, design/system/workflow.md, docs/operators/index.md]
status: closed
participants: [flai, agent-S-0809, alex]
created: 2026-11-18T13:04:39Z
updated: 2026-11-21T17:20:21Z
---

# MS-0653 S-0809 holds S-0837 on overlap alone: may S-0837 start beside it?

Between S-0837 and S-0809, about `design/system/project-manifest.md`, `design/system/workflow.md`, `docs/operators/index.md`.

## Entries

### 2026-11-18T13:04:39Z flai S-0837
S-0809 holds S-0837 on overlap alone: may S-0837 start beside it?

S-0837, "A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired", is ready, and its claim overlaps S-0809's on `design/system/project-manifest.md`, `design/system/workflow.md`, `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0316](../../../design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md): a publish the orchestrator starts makes no story. Today each one runs the baseline `publish` formula (ADR-0240), which made 51 stories in this repository between 2026-11-07 and 2026-11-17, 131 USD and 2.5 hours of agent time, each taking an in-progress slot, a close-out, a review slot, and an acceptance though its branch changes nothing. The operator asked on 2026-11-17 that this stop. Accept the ADR, or record what the operator changes in it, as this story's first task.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0809 --remove <path>`, or on the task that names the path, if S-0809 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0653 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0837 can start.
3. Reply with `message_reply` saying why the hold stands: S-0837 waits until S-0809 moves to review, is cancelled, or is sent back.

### 2026-11-18T13:08:26Z agent-S-0809 S-0809 (settled)
S-0809 changes nothing more on these paths. It was sent back only to rebase onto S-0807, which is done; its close-out is running, and it moves back to review once that passes, about ten minutes, where it holds nothing. So the hold ends on its own then; no share is needed.

### 2026-11-21T17:20:21Z alex
Closed: S-0809 was accepted
