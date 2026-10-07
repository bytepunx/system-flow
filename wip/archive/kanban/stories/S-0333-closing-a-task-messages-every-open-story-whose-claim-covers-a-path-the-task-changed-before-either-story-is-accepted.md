---
id: S-0333
type: story
nature: feature
title: Closing a task messages every open story whose claim covers a path the task changed, before either story is accepted
status: done
parent: E-0018
owner: alex
created: 2026-10-07T20:10:50Z
updated: 2026-10-07T22:48:25Z
transitions:
  - to: ready
    at: 2026-10-07T21:43:11Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T22:09:26Z
    by: agent-S-0333
  - to: review
    at: 2026-10-07T22:47:32Z
    by: agent-S-0333
  - to: done
    at: 2026-10-07T22:48:25Z
    by: orchestrator
tags: [flai, template]
topics: [cli, git, conventions, template]
touches: [flai/internal/itemedit/covers.go, flai/internal/itemedit/covers_test.go, flai/cmd/accept_overlap.go, flai/cmd/accept_overlap_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task.go, flai/internal/mcpserver/task_test.go, design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, design/system/agent-coordination.md]
after: [S-0331]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2298
  turns:
    - day: 2026-10-07
      ceremony: 1
      test_runs: 1
      hand_edits: 2
      work: 63
  models:
    - model: claude-opus-5-5
      input: 278
      output: 93365
      cache_read: 16872294
      cache_write: 451143
      cost: 8.1585
  strategic:
    - kind: orchestrator
      seconds: 541
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 67
          output: 1103
          cache_read: 5852512
          cache_write: 20049
          cost: 1.4471
        - model: claude-sonnet-5-5
          input: 12
          output: 63
          cache_read: 182548
          cache_write: 50187
          cost: 0.2006
cost_of_delay:
  value: 113.7
  by: planner-E-0018
  at: 2026-10-07T20:22:02Z
forecast:
  duration: 44m
  delivery: 2026-10-08T08:30:00Z
  basis: "Its own forecast of 44m; 39th in the pull order with an in-progress limit of 3, behind S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0280, S-0287, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0309, S-0310, S-0312, S-0313, S-0314, S-0315, S-0316, S-0317, S-0318, S-0319, S-0320, S-0321, S-0322, S-0323, S-0324, S-0325, S-0326, S-0327 and S-0332."
  by: flai
  at: 2026-10-07T21:42:43Z
finalized:
  by: orchestrator
  at: 2026-10-07T20:24:17Z
---
# S-0333 Closing a task messages every open story whose claim covers a path the task changed, before either story is accepted

## Goal

Tell an agent what another story changed while both are still working, not only when one is accepted. Today an open story learns that an overlapping story changed its paths from the `overlapped` notice at acceptance (S-0132), when its own branch may already rely on the old shape. When `flai task done` commits a task, flai messages each other open story whose claim covers a path the commit changed, with the paths and the commit's subject, so its agent can adjust early or answer that the change breaks it.

## Acceptance criteria

- [x] `flai task done` and the MCP tool `task_done`, after committing, message each other story in progress or in review whose claim covers a path the commit changed, inside the shared paths too, `about` those paths, naming the task, the commit, and its subject; one conversation per pair of stories, reused when it is open.
- [x] A story with an empty claim is told of every path, as the notice at acceptance does.
- [x] A commit that changes no path another open story claims sends nothing, and a notice that cannot be sent is logged and does not fail the task's close.
- [x] `flai task done --json` and `task_done` return whom they told as `told`.
- [x] `design/system/workflow.md` § Branches and collisions, `agent-narrative.md`, and `work-management.md` in both copies say what a story's agent does with such a message.

## Tasks

- T-1199 Who an open story's change reaches is worked out in one place, shared by acceptance and task done
- T-1200 flai task done and task_done message each open story whose claim covers a path the task's commit changed, and return told
- T-1201 The workflow, the convention, the template, and the guides say what an agent does with a task's change notice

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07. It waits for S-0331 (`after`), so the agents it tells can read the message. It does not wait for S-0332; the two touch the same design documents, so the hold runs them one after the other.

This is the early notice of design K in `design/system/agent-coordination.md` (order plus notices, CoAgent and STALE): telling the later agent what the earlier one changed.

Layers:

1. T-1199, moving the coverage test out of `flai/cmd` so `taskdone` can call it, with no change of behaviour.
2. T-1200, the notice.
3. T-1201, the docs.

Touches:

- **Declared:** none before planning.
- **Layout:**
  - `flai/cmd/accept_overlap.go` and its test: `overlapsOf` and `covered` work out who the notice at acceptance tells.
  - `flai/internal/itemedit/covers.go` and its test: the new home of that test, beside `claim.go`; T-1199 may pick another file name.
  - `flai/internal/taskdone/taskdone.go`, `flai/cmd/task_done.go`, `flai/internal/mcpserver/task.go`, and their tests: the close of a task and its two front ends.
- **Co-change:** `flai touches suggest` from `taskdone.go` and `accept_overlap.go` found nothing changed with them often enough (3 commits).
- **Design:** `design/system/workflow.md` § Branches and collisions and `agent-narrative.md`; both copies of `work-management.md` and `template/CHANGELOG.md`; `flai-cli.md`, `docs/users/flai.md`, and the generated reference.

Forecast 44m, delivery 2026-10-08T09:32Z: `flai forecast` gave it, 114 s per unit over 29 done large-band feature stories, times size 23. It stands.

Cost of delay 113.70 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 44m of 6h27m. It stands.

### Accepted by the orchestrator

- Verified: f2e5abba7b6984c91cdc25bf273c6a782ba4a7f3
- At: 2026-10-07T22:48:25Z

Verdict: meets all criteria (verifier at f2e5abba7b6984c91cdc25bf273c6a782ba4a7f3; flai verify passed every step at that commit)

- 1: flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/task_done.go, flai/internal/mcpserver/task.go
- 2: flai/internal/itemedit/covers.go, flai/internal/itemedit/covers_test.go, flai/cmd/accept_overlap.go, flai/internal/taskdone/taskdone_test.go
- 3: flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go
- 4: flai/internal/taskdone/taskdone.go, flai/cmd/task_done.go, flai/cmd/task_done_test.go, flai/internal/mcpserver/task_test.go
- 5: design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/agent-coordination.md, docs/users/flai.md, docs/users/flai-reference.md
