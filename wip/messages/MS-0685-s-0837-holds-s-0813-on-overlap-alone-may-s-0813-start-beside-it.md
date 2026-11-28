---
id: MS-0685
title: "S-0837 holds S-0813 on overlap alone: may S-0813 start beside it?"
from: S-0813
to: S-0837
about: [flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, design/system/project-manifest.md, design/system/workflow.md, docs/operators/settings.md]
status: open
participants: [flai]
created: 2026-11-28T00:43:34Z
updated: 2026-11-28T00:43:34Z
---

# MS-0685 S-0837 holds S-0813 on overlap alone: may S-0813 start beside it?

Between S-0813 and S-0837, about `flai/internal/upgrade/upgrade.go`, `flai/internal/upgrade/upgrade_test.go`, `design/system/project-manifest.md`, `design/system/workflow.md`, `docs/operators/settings.md`.

## Entries

### 2026-11-28T00:43:34Z flai S-0813
S-0837 holds S-0813 on overlap alone: may S-0813 start beside it?

S-0813, "A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired", is ready, and its claim overlaps S-0837's on `flai/internal/upgrade/upgrade.go`, `flai/internal/upgrade/upgrade_test.go`, `design/system/project-manifest.md`, `design/system/workflow.md`, `docs/operators/settings.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §8. `queue: true` leaves a tier to the batch and `when:` runs it at the close-out anyway; with the file naming the suites of each step, the two keys say the same thing twice.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0837 --remove <path>`, or on the task that names the path, if S-0837 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0685 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0813 can start.
3. Reply with `message_reply` saying why the hold stands: S-0813 waits until S-0837 moves to review, is cancelled, or is sent back.
