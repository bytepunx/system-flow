---
id: S-0332
type: story
nature: improvement
title: flai tells two stories' agents of a trial-merge conflict or a grown overlap with a message between them, and asks the operator only when they do not agree
status: done
parent: E-0018
owner: alex
created: 2026-10-07T20:10:41Z
updated: 2026-10-07T23:20:33Z
transitions:
  - to: ready
    at: 2026-10-07T21:43:11Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T22:47:38Z
    by: agent-S-0332
  - to: review
    at: 2026-10-07T23:19:35Z
    by: agent-S-0332
  - to: done
    at: 2026-10-07T23:20:33Z
    by: orchestrator
tags: [flai, template]
topics: [cli, git, conventions, template]
touches: [design/adrs, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/touches.go, flai/internal/inbox/inbox.go, flai/internal/mcpserver/items_write_test.go, docs/operators/settings.md, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server_test.go, design/system/agent-coordination.md, design/system/continuous-improvement.md]
after: [S-0331]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1933
  turns:
    - day: 2026-10-07
      ceremony: 2
      hand_edits: 2
      work: 59
  models:
    - model: claude-opus-5-5
      input: 374
      output: 148336
      cache_read: 22700507
      cache_write: 719713
      cost: 11.8244
  strategic:
    - kind: orchestrator
      seconds: 545
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 103
          output: 1612
          cache_read: 11965333
          cache_write: 71267
          cost: 2.9658
        - model: claude-sonnet-5-5
          input: 12
          output: 79
          cache_read: 166041
          cache_write: 51964
          cost: 0.1879
cost_of_delay:
  value: 142.12
  by: planner-E-0018
  at: 2026-10-07T20:22:00Z
forecast:
  duration: 55m
  delivery: 2026-10-07T23:26:00Z
  basis: "Its own forecast of 55m; 1st in the pull order with an in-progress limit of 3, behind S-0232 and S-0333."
  by: flai
  at: 2026-10-07T22:25:33Z
finalized:
  by: orchestrator
  at: 2026-10-07T20:24:16Z
---
# S-0332 flai tells two stories' agents of a trial-merge conflict or a grown overlap with a message between them, and asks the operator only when they do not agree

## Goal

Route flai's own coordination notices to the agents that must act on them. Today `flai stream sync` opens a thread on one story when its branch conflicts with another open story's, and a claim that grows into another story's sends both an `overlapped` change telling them to coordinate on a thread. Both reach the operator, who has nothing to decide until the agents disagree. Send each as a message conversation between the two stories, with the paths, so the agents settle it while both are working, and escalate to a thread on the operator only when one of them asks.

## Acceptance criteria

- [x] A trial merge at `flai stream sync` that conflicts with another open story's branch opens, or adds an entry to, one conversation between the two stories, `about` the conflicting paths, instead of a thread; the conversation closes when a later sync finds the two merging cleanly or the other story is no longer open.
- [x] A claim that grows to overlap another in-progress story's (S-0244) opens, or adds to, the same pair's conversation, `about` the paths gained, beside the `overlapped` change both stories still get.
- [x] `flai message escalate <conversation> "<reason>"`, and the MCP tool `message_escalate`, open a thread on the operator that names both stories, links the conversation, and says what they could not agree; the conversation records it.
- [x] Conflict threads already open are left as they are; no new conflict thread is opened by a sync.
- [x] An ADR refining ADR-0046 records the change, and `design/system/workflow.md` § Branches and collisions, `agent-narrative.md`, and `work-management.md` in both copies describe it.

## Tasks

- T-1194 An ADR refining ADR-0046 sends conflicts and grown overlaps to the two stories as a message, with escalation to the operator
- T-1195 flai stream sync keeps one conversation per conflicting pair of stories instead of a conflict thread
- T-1196 A claim grown into another in-progress story's opens or adds to the pair's conversation
- T-1197 flai message escalate and the MCP tool message_escalate open a thread on the operator from a conversation
- T-1198 The workflow, the convention, the template, and the guides say conflicts and grown overlaps arrive as messages

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07. It waits for S-0331 (`after`), so agents can read and answer the conversations it opens.

Layers:

1. T-1194, the ADR.
2. T-1195, the sync; T-1196, the grown claim; T-1197, escalation. They share no path and run together.
3. T-1198, the docs.

Touches:

- **Declared:** none before planning.
- **Layout:**
  - `flai/internal/storygit/sync.go` and its test: `reportConflicts` opens the conflict thread today.
  - `flai/internal/itemedit/claim.go` and its test: `ClaimWatch.Grown` records the grown overlap (S-0244).
  - The messages package, `flai/cmd/message.go`, and `flai/internal/mcpserver/messages.go` with their tests, for escalation; `flai/internal/guard/guard.go` and its test, so a sub-agent cannot escalate.
- **Co-change:** `flai touches suggest` from `sync.go` and `claim.go` gave `flai/cmd/stream_sync_test.go` (67%), `sync_test.go` and `claim_test.go` (50%), and `flai/cmd/stream.go` and `stream_sync.go` (33%).
- **Design:** `design/system/workflow.md` § Branches and collisions and `agent-narrative.md`, which describe the conflict thread and the `overlapped` change; both copies of `work-management.md` and `template/CHANGELOG.md`; `flai-cli.md`, `docs/users/flai.md`, and the generated reference.
- **Folder touch kept:** `design/adrs`, for T-1194's ADR, whose file `flai adr new` names. Inside `claims.shared`.
- **Not taken:** `flai/cmd/branch.go` (50%) and `flai/internal/workitem/hold.go` (33%): neither the branch command nor the hold changes here.

Forecast 55m, delivery 2026-10-08T09:29Z.

- `flai forecast` gave 38m: 78 s per unit over 44 done large-band improvement stories, times size 29 (5 criteria, 24 touches).
- Raised by 17m to the feature rate, 114 s per unit: besides rerouting two notices it adds a command and an MCP tool, which is feature work. Delivery is shifted by the same 17m.

Cost of delay 142.12 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 55m of 6h27m. It stands.

### Accepted by the orchestrator

- Verified: 2afebf1715724360aba2750585cf3d6b67cdcca8
- At: 2026-10-07T23:20:33Z

Verdict: meets all criteria (verifier at 2afebf1715724360aba2750585cf3d6b67cdcca8; flai verify passed every step at that commit). Criterion 4: no sync opens or adds to a conflict thread; already-open ones keep the resolution main already gave them on a clean merge or a departed story.

- 1: flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go
- 2: flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/touches.go, flai/internal/inbox/inbox.go
- 3: flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/folder.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go
- 4: flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go
- 5: design/adrs/0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md, design/adrs/README.md, design/system/workflow.md, design/system/agent-narrative.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md
