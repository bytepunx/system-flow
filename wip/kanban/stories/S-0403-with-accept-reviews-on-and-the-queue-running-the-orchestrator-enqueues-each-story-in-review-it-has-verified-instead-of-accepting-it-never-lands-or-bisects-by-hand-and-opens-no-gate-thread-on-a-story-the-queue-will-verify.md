---
id: S-0403
type: story
nature: feature
title: With accept_reviews on and the queue running, the orchestrator enqueues each story in review it has verified instead of accepting it, never lands or bisects by hand, and opens no gate thread on a story the queue will verify
status: backlog
parent: E-0022
owner: alex
created: 2026-10-16T04:58:45Z
updated: 2026-10-17T01:02:24Z
transitions: []
tags: [cli]
topics: [orchestration, conventions, template]
touches: [flai/internal/mcpserver/queue.go, flai/internal/mcpserver/queue_test.go, flai/internal/mcpserver/folder.go, flai/internal/guard/writes.go, flai/internal/guard/orchestrate_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/root/.claude/agents/orchestrator.md, ".claude/agents/orchestrator.md", template/template.yaml, template/CHANGELOG.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0400]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 42m
  delivery: 2026-10-17T08:36:00Z
  basis: "Its own forecast of 42m; 18th in the pull order with an in-progress limit of 5, behind S-0438, S-0440, S-0441, S-0439, S-0442, S-0443, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401 and S-0402."
  by: flai
  at: 2026-10-17T01:02:24Z
---
# S-0403 With accept_reviews on and the queue running, the orchestrator enqueues each story in review it has verified instead of accepting it, never lands or bisects by hand, and opens no gate thread on a story the queue will verify

## Goal

The orchestrator stops accepting stories one by one once the queue runs. It still judges each story, with its verifier and the ADR-0093 evidence, but it hands the story to the queue rather than merging it, and it leaves to the queue the blockers the queue settles, so that it no longer opens threads on gates that a batch would clear.

## Acceptance criteria

- [ ] The MCP tool `queue_add`, and `flai queue add S-nnnn --by orchestrator --evidence -`, enqueue a story with the orchestrator's evidence and the commit its verifier checked, refused unless `accept_reviews` is on, with the same evidence rules as `flai accept --by orchestrator`.
- [ ] While the `merge-queue` host action is on, `flai accept --by orchestrator` is refused, naming the queue, and `flai guard` refuses the orchestrator every queue write but `queue_add`, so that it never lands, bisects, holds, or dequeues by hand.
- [ ] The convention `strategic-agents.md` § As the orchestrator, in the repository and the template, says: with the queue on, verify each story as today, enqueue it with the evidence, log each enqueue with `activity_log` naming the story, its place in the queue, and the policy figure, and open no thread on a story for a blocker the queue settles (a verify not at the head, a dirty main checkout, a stop inside the shared paths).
- [ ] The orchestrator's agent definition, in `.claude/agents/` and the template's, names `queue_add`. `template/CHANGELOG.md` names the convention's change.
- [ ] `design/system/strategic-agents.md` § Accepting a story and `docs/users/flai.md` describe the orchestrator's enqueueing.

## Tasks
- T-1649 flai queue add takes the orchestrator's evidence and verified commit, and the MCP tool queue_add enqueues through it
- T-1650 While the queue runs, flai accept --by orchestrator is refused, naming the queue, and flai guard refuses the orchestrator every queue write but queue_add
- T-1663 The orchestrator's convention, agent definition, and prompt tell it to enqueue with evidence while the queue runs, and to open no thread on a blocker the queue settles
- T-1672 The strategic agents' design, the CLI docs, and the template's changelog describe the orchestrator's enqueueing

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
