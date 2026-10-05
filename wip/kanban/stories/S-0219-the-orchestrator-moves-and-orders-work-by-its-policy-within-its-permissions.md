---
id: S-0219
type: story
nature: feature
title: The orchestrator moves and orders work by its policy within its permissions
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-05T03:46:30Z
transitions: []
tags: [flai]
touches: [flai/internal/harness, ".claude/agents/orchestrator.md", template, design/system/strategic-agents.md, flai/internal/guard, flai/internal/serve]
after: [S-0218, S-0209]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 73.17
  by: planner-E-0016
  at: 2026-10-04T04:48:19Z
forecast:
  duration: 1h30m
  delivery: 2026-10-05T21:22:00Z
  basis: "Its own forecast of 1h30m; 12th in the pull order with an in-progress limit of 3, behind S-0253, S-0257, S-0244, S-0262, S-0258, S-0212, S-0213, S-0214, S-0215, S-0216, S-0217 and S-0218."
  by: flai
  at: 2026-10-05T03:46:30Z
---
# S-0219 The orchestrator moves and orders work by its policy within its permissions

## Goal

With its permissions on, the orchestrator keeps the flow: it asks the planner to draft stories for backlog epics, finalizes drafts it judges ready, promotes backlog stories to ready while the limits leave room, and orders the ready column by the policy.

## Acceptance criteria
- [ ] With `plan_backlog_epics`, it runs the planner for a backlog epic that has no stories, and for one whose stories are all done when the epic is not
- [ ] With `finalize_drafts`, it finalizes a draft whose criteria, touches, forecast, and value are complete and consistent, and leaves one that is not with a thread saying why
- [ ] With `promote_to_ready`, it promotes `flai promote --candidates` in order while the ready limit has room, never a draft, never a held story
- [ ] With `order_ready`, it applies `flai order --by <policy> --apply` after each change to the ready column, and does not reorder a story the operator ordered by hand in the last day
- [ ] Every action is logged with the policy figure that justified it; a refusal from `flai guard` ends the attempt and is logged
- [ ] Tests cover each permission on and off against a fixture board

## Tasks

## Notes
