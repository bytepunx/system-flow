---
id: MS-0267
title: "S-0550 holds S-0551 on overlap alone: may S-0551 start beside it?"
from: S-0551
to: S-0550
about: [flai/internal/serve/hold_test.go, flai/internal/sync, flai/internal/claims, design/system/workflow.md, flai/internal/check/check.go]
status: open
participants: [flai, agent-S-0550]
created: 2026-10-31T08:16:11Z
updated: 2026-10-31T08:16:58Z
shares:
  - holder: S-0550
    held: S-0551
    paths: [flai/internal/serve/hold_test.go, design/system/workflow.md, flai/internal/check/check.go]
    split: "S-0550 adds one test to serve/hold_test.go showing two ready stories that both edit a claims.shared path (docs/users/flai.md) start together, and changes no hold code; S-0551 changes serve's start and hold code and its own hold tests. In workflow.md S-0550 adds the overlap-notice before-figures and says the acceptance notice skips claims.shared; S-0551 rewrites the hold and after policy sections. In check.go S-0550 adds the generated-index checks and --fix; S-0551 adds the after-reason warning. S-0550 also dropped flai/internal/sync and flai/internal/claims from its touches. Whichever lands second rebases over the other."
    by: agent-S-0550
    at: 2026-10-31T08:16:58Z
---

# MS-0267 S-0550 holds S-0551 on overlap alone: may S-0551 start beside it?

Between S-0551 and S-0550, about `flai/internal/serve/hold_test.go`, `flai/internal/sync`, `flai/internal/claims`, `design/system/workflow.md`, `flai/internal/check/check.go`.

## Entries

### 2026-10-31T08:16:11Z flai S-0551
S-0550 holds S-0551 on overlap alone: may S-0551 start beside it?

S-0551, "A story starts on a trial merge instead of an overlap hold, and after holds a story for a dependency only, so the worktrees run stories in parallel as people on one remote do", is ready, and its claim overlaps S-0550's on `flai/internal/serve/hold_test.go`, `flai/internal/sync`, `flai/internal/claims`, `design/system/workflow.md`, `flai/internal/check/check.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The worktrees isolate every story, and the scheduler then serializes them: with an in-progress limit of 5, nothing was in progress for 42 percent of the fortnight, one story for 29 percent, two for 26 percent, three for 2 percent, and four or five never, mean 1.6 when anything ran. The journal records 124 overlap holds and 113 `after` holds against 62 real conflicts, the overlap holds on `check.go`, the strategic-agents design and convention, the manifest and dashboard designs, and the `after` chains planner-made, S-0417 to S-0422 each waiting on the one before, as I-0152 already noted. People sharing one remote merge optimistically and resolve the rare conflict; this story makes flai do the same. A ready story whose touches overlap an in-progress story's is started when a trial merge of the in-progress story's branch into a worktree at the ready story's base is clean, and held only when it is not, with the journal naming the conflicting file; `after` is honoured as today but the planner sets it only when a story reads what another writes, named in the story's notes, and `flai check` warns of an `after` with no reason; and the ready wait at p85, 5.4 hours in the fortnight, is a figure `flai stats` reports.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0550 --remove <path>`, or on the task that names the path, if S-0550 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0267 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0551 can start.
3. Reply with `message_reply` saying why the hold stands: S-0551 waits until S-0550 moves to review, is cancelled, or is sent back.

### 2026-10-31T08:16:58Z agent-S-0550 S-0550
S-0550 shares `flai/internal/serve/hold_test.go`, `design/system/workflow.md`, `flai/internal/check/check.go` with S-0551, split so:

> S-0550 adds one test to serve/hold_test.go showing two ready stories that both edit a claims.shared path (docs/users/flai.md) start together, and changes no hold code; S-0551 changes serve's start and hold code and its own hold tests. In workflow.md S-0550 adds the overlap-notice before-figures and says the acceptance notice skips claims.shared; S-0551 rewrites the hold and after policy sections. In check.go S-0550 adds the generated-index checks and --fix; S-0551 adds the after-reason warning. S-0550 also dropped flai/internal/sync and flai/internal/claims from its touches. Whichever lands second rebases over the other.

The overlap between the two on them no longer holds (ADR-0134).
