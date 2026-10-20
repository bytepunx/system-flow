---
id: S-0403
type: story
nature: feature
title: With accept_reviews on and the queue running, the orchestrator enqueues each story in review it has verified instead of accepting it, never lands or bisects by hand, and opens no gate thread on a story the queue will verify
status: done
parent: E-0022
owner: alex
created: 2026-10-16T04:58:45Z
updated: 2026-10-20T18:59:48Z
transitions:
  - to: ready
    at: 2026-10-19T04:53:57Z
    by: alex
  - to: in-progress
    at: 2026-10-19T23:41:28Z
    by: agent-S-0403
  - to: review
    at: 2026-10-20T04:41:04Z
    by: agent-S-0403
  - to: done
    at: 2026-10-20T18:59:48Z
    by: alex
tags: [cli]
topics: [orchestration, conventions, template]
touches: [flai/internal/mcpserver/queue.go, flai/internal/mcpserver/queue_test.go, flai/internal/mcpserver/folder.go, flai/internal/guard/writes.go, flai/internal/guard/orchestrate_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/root/.claude/agents/orchestrator.md, ".claude/agents/orchestrator.md", template/template.yaml, template/CHANGELOG.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server_test.go, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/cmd/guard.go, flai/internal/serve/orchestrate_hold.go, docs/operators/settings.md, design/issues/I-0169-message-share-refuses-a-split-between-two-stories-in-progress-so-an-agreed-split-cannot-lift-the-overlap-that-refuses-the-other-story-s-flai-task-new.md]
after: [S-0400]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4008
  priced_by: harness
  turns:
    - day: 2026-10-19
      ceremony: 1
      hand_edits: 2
      work: 30
    - day: 2026-10-20
      ceremony: 2
      hand_edits: 1
      work: 82
  models:
    - model: claude-opus-5-5
      input: 554
      output: 193484
      cache_read: 40591653
      cache_write: 1472644
      cost: 21.1268
  strategic:
    - kind: orchestrator
      seconds: 13995
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 121
          output: 2460
          cache_read: 26508813
          cache_write: 34957
          cost: 6.6546
        - model: claude-sonnet-5-5
          input: 8
          output: 35
          cache_read: 115986
          cache_write: 73480
          cost: 0.1551
forecast:
  duration: 42m
  delivery: 2026-10-20T00:27:00Z
  basis: "Its own forecast of 42m; 2nd in the pull order with an in-progress limit of 5, behind S-0404."
  by: flai
  at: 2026-10-19T23:41:20Z
finalized:
  by: alex
  at: 2026-10-19T04:53:39Z
---
# S-0403 With accept_reviews on and the queue running, the orchestrator enqueues each story in review it has verified instead of accepting it, never lands or bisects by hand, and opens no gate thread on a story the queue will verify

## Goal

The orchestrator stops accepting stories one by one once the queue runs. It still judges each story, with its verifier and the ADR-0093 evidence, but it hands the story to the queue rather than merging it, and it leaves to the queue the blockers the queue settles, so that it no longer opens threads on gates that a batch would clear.

## Acceptance criteria

- [x] The MCP tool `queue_add`, and `flai queue add S-nnnn --by orchestrator --evidence -`, enqueue a story with the orchestrator's evidence and the commit its verifier checked, refused unless `accept_reviews` is on, with the same evidence rules as `flai accept --by orchestrator`.
- [x] While the `merge-queue` host action is on, `flai accept --by orchestrator` is refused, naming the queue, and `flai guard` refuses the orchestrator every queue write but `queue_add`, so that it never lands, bisects, holds, or dequeues by hand.
- [x] The convention `strategic-agents.md` § As the orchestrator, in the repository and the template, says: with the queue on, verify each story as today, enqueue it with the evidence, log each enqueue with `activity_log` naming the story, its place in the queue, and the policy figure, and open no thread on a story for a blocker the queue settles (a verify not at the head, a dirty main checkout, a stop inside the shared paths).
- [x] The orchestrator's agent definition, in `.claude/agents/` and the template's, names `queue_add`. `template/CHANGELOG.md` names the convention's change.
- [x] `design/system/strategic-agents.md` § Accepting a story and `docs/users/flai.md` describe the orchestrator's enqueueing.

## Tasks
- T-1649 flai queue add takes the orchestrator's evidence and verified commit, and the MCP tool queue_add enqueues through it
- T-1650 While the queue runs, flai accept --by orchestrator is refused, naming the queue, and flai guard refuses the orchestrator every queue write but queue_add
- T-1663 The orchestrator's convention, agent definition, and prompt tell it to enqueue with evidence while the queue runs, and to open no thread on a blocker the queue settles
- T-1672 The strategic agents' design, the CLI docs, and the template's changelog describe the orchestrator's enqueueing
- T-1910 The orchestrator's agent definitions name queue_add

## Notes

- Drafted by the planner for E-0022, after S-0400, since an enqueued story lands as the orchestrator only once landing exists.
- The agent definitions are under `.claude/`, which a story's agent writes only through the operator's permission prompt.

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Layout, from the explorer's map:
    - MCP tools register in `flai/internal/mcpserver/folder.go`
    - the guard's role rules: `flai/internal/guard/writes.go` and its `orchestrate_test.go`
    - the orchestrator's prompt: `flai/internal/serve/orchestrate.go`
    - the agent definitions: `.claude/agents/orchestrator.md` and `template/root/.claude/agents/orchestrator.md`.
  - Design: the epic's criterion 7; the convention `strategic-agents.md` § As the orchestrator and its template copy.
  - Co-change: `flai touches suggest` added `design/system/flai-cli.md`, `docs/users/flai-reference.md`, and `template/template.yaml`.
  - No folder touch.
- **Forecast.** flai's 42m is kept: one flag, one MCP tool, two refusals, and convention text.
- **Cost of delay.** Waits on TH-0465 question 1.
