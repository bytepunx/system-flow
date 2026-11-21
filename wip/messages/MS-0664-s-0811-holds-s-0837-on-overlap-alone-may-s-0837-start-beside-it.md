---
id: MS-0664
title: "S-0811 holds S-0837 on overlap alone: may S-0837 start beside it?"
from: S-0837
to: S-0811
about: [flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_test.go, flai/cmd/release.go, template/root/design/workflows/formulas/publish.md, design/system/workflow.md, design/system/pushing-from-the-board.md]
status: open
participants: [flai, agent-S-0811]
created: 2026-11-21T18:13:46Z
updated: 2026-11-21T18:15:07Z
---

# MS-0664 S-0811 holds S-0837 on overlap alone: may S-0837 start beside it?

Between S-0837 and S-0811, about `flai/internal/hostapi/writes_outcome.go`, `flai/internal/hostapi/writes_test.go`, `flai/cmd/release.go`, `template/root/design/workflows/formulas/publish.md`, `design/system/workflow.md`, `design/system/pushing-from-the-board.md`.

## Entries

### 2026-11-21T18:13:46Z flai S-0837
S-0811 holds S-0837 on overlap alone: may S-0837 start beside it?

S-0837, "A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired", is ready, and its claim overlaps S-0811's on `flai/internal/hostapi/writes_outcome.go`, `flai/internal/hostapi/writes_test.go`, `flai/cmd/release.go`, `template/root/design/workflows/formulas/publish.md`, `design/system/workflow.md`, `design/system/pushing-from-the-board.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0316](../../../design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md): a publish the orchestrator starts makes no story. Today each one runs the baseline `publish` formula (ADR-0240), which made 51 stories in this repository between 2026-11-07 and 2026-11-17, 131 USD and 2.5 hours of agent time, each taking an in-progress slot, a close-out, a review slot, and an acceptance though its branch changes nothing. The operator asked on 2026-11-17 that this stop. Accept the ADR, or record what the operator changes in it, as this story's first task.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0811 --remove <path>`, or on the task that names the path, if S-0811 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0664 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0837 can start.
3. Reply with `message_reply` saying why the hold stands: S-0837 waits until S-0811 moves to review, is cancelled, or is sent back.

### 2026-11-21T18:15:07Z agent-S-0811 S-0811 (settled)
The hold stands, briefly. S-0811 is rebased onto main and its close-out job is running now. It moves to review as soon as that passes, which lifts the hold. The overlap is real code, not wording: S-0811 changes `flai/cmd/release.go`, the `hostapi` publish outcome, and the template's publish formula, which S-0837 retires. Starting S-0837 on main after S-0811's change saves both stories a rebase over those files.
