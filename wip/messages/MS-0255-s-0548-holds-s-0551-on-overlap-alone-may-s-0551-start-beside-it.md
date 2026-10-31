---
id: MS-0255
title: "S-0548 holds S-0551 on overlap alone: may S-0551 start beside it?"
from: S-0551
to: S-0548
about: [flai/internal/serve/agents_test.go, flai/internal/serve/serve_test.go, flai/internal/serve/forgemerge_test.go, flai/internal/serve/serve.go, flai/internal/serve/projects_test.go, flai/internal/serve/mcp_test.go, flai/internal/serve/budget_test.go]
status: open
participants: [flai, agent-S-0548]
created: 2026-10-31T07:26:05Z
updated: 2026-10-31T08:15:07Z
shares:
  - holder: S-0548
    held: S-0551
    paths: [flai/internal/serve/agents_test.go, docs/users/flai.md, design/system/flai-cli.md, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/mcpserver/work_test.go, template/CHANGELOG.md]
    split: "S-0548 has committed its edits to agents_test.go (waitFor's 30 s deadline, the send-back subtest's hold), work_test.go (t.Parallel), and check.go (the manifest.tests-timeout warning); it still makes TestMonorepoIsClean in check_test.go skip under -short and adds the flai test tier-name argument to flai.md and flai-cli.md. S-0551 changes serve's start and hold code and check's after warning, with its own tests. New tests S-0551 adds to work_test.go call t.Parallel unless they set an environment variable. Whichever lands second rebases over the other."
    by: agent-S-0548
    at: 2026-10-31T08:15:07Z
---

# MS-0255 S-0548 holds S-0551 on overlap alone: may S-0551 start beside it?

Between S-0551 and S-0548, about `flai/internal/serve/agents_test.go`, `flai/internal/serve/serve_test.go`, `flai/internal/serve/forgemerge_test.go`, `flai/internal/serve/serve.go`, `flai/internal/serve/projects_test.go`, `flai/internal/serve/mcp_test.go`, `flai/internal/serve/budget_test.go`.

## Entries

### 2026-10-31T07:26:05Z flai S-0551
S-0548 holds S-0551 on overlap alone: may S-0551 start beside it?

S-0551, "A story starts on a trial merge instead of an overlap hold, and after holds a story for a dependency only, so the worktrees run stories in parallel as people on one remote do", is ready, and its claim overlaps S-0548's on `flai/internal/serve/agents_test.go`, `flai/internal/serve/serve_test.go`, `flai/internal/serve/forgemerge_test.go`, `flai/internal/serve/serve.go`, `flai/internal/serve/projects_test.go`, `flai/internal/serve/mcp_test.go`, `flai/internal/serve/budget_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The worktrees isolate every story, and the scheduler then serializes them: with an in-progress limit of 5, nothing was in progress for 42 percent of the fortnight, one story for 29 percent, two for 26 percent, three for 2 percent, and four or five never, mean 1.6 when anything ran. The journal records 124 overlap holds and 113 `after` holds against 62 real conflicts, the overlap holds on `check.go`, the strategic-agents design and convention, the manifest and dashboard designs, and the `after` chains planner-made, S-0417 to S-0422 each waiting on the one before, as I-0152 already noted. People sharing one remote merge optimistically and resolve the rare conflict; this story makes flai do the same. A ready story whose touches overlap an in-progress story's is started when a trial merge of the in-progress story's branch into a worktree at the ready story's base is clean, and held only when it is not, with the journal naming the conflicting file; `after` is honoured as today but the planner sets it only when a story reads what another writes, named in the story's notes, and `flai check` warns of an `after` with no reason; and the ready wait at p85, 5.4 hours in the fortnight, is a figure `flai stats` reports.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0548 --remove <path>`, or on the task that names the path, if S-0548 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0255 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0551 can start.
3. Reply with `message_reply` saying why the hold stands: S-0551 waits until S-0548 moves to review, is cancelled, or is sent back.

### 2026-10-31T08:15:07Z agent-S-0548 S-0548
S-0548 shares `flai/internal/serve/agents_test.go`, `docs/users/flai.md`, `design/system/flai-cli.md`, `flai/internal/check/check.go`, `flai/internal/check/check_test.go`, `flai/internal/mcpserver/work_test.go`, `template/CHANGELOG.md` with S-0551, split so:

> S-0548 has committed its edits to agents_test.go (waitFor's 30 s deadline, the send-back subtest's hold), work_test.go (t.Parallel), and check.go (the manifest.tests-timeout warning); it still makes TestMonorepoIsClean in check_test.go skip under -short and adds the flai test tier-name argument to flai.md and flai-cli.md. S-0551 changes serve's start and hold code and check's after warning, with its own tests. New tests S-0551 adds to work_test.go call t.Parallel unless they set an environment variable. Whichever lands second rebases over the other.

The overlap between the two on them no longer holds (ADR-0134).
