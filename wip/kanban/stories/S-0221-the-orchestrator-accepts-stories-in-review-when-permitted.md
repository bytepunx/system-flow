---
id: S-0221
type: story
nature: feature
title: The orchestrator accepts stories in review when permitted
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-04T04:48:20Z
transitions: []
tags: [flai]
touches: [flai/internal/harness, flai/internal/hostapi, flai/internal/preview, flai/internal/mcpserver, design/adrs, flai/cmd/accept.go, flai/internal/guard, flaiover/src/lib/components/Review.svelte, flaiover/src/routes/items, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 73.17
  by: planner-E-0016
  at: 2026-10-04T04:48:20Z
forecast:
  duration: 1h30m
  delivery: 2026-10-05T05:30:00Z
  basis: "An ADR refining ADR-0032, the orchestrator's identity on accept.run and flai accept, guard checks, and who-accepted on the review page; like S-0201 (3771 s); serial after S-0220 on harness."
  by: planner-E-0016
  at: 2026-10-04T04:44:56Z
---
# S-0221 The orchestrator accepts stories in review when permitted

## Goal

ADR-0032 and `item_move` make acceptance the operator's alone. With `accept_reviews` on, the orchestrator accepts a story in review when the verifier's run passed, every criterion is ticked, the diff stays within the story's touches, and no thread on it is open; otherwise it leaves it with a thread saying what is missing.

## Acceptance criteria
- [ ] An ADR refines ADR-0032: acceptance by the orchestrator under `accept_reviews`, with the conditions above, recorded as `by: orchestrator`
- [ ] The hostapi `accept.run` and `flai accept` take the orchestrator's identity when the permission is on; `flai guard` refuses otherwise
- [ ] Before accepting it runs the acceptance preview and refuses on any blocker, and it never accepts a story whose criteria it cannot check against the diff
- [ ] Each acceptance is logged with the evidence; the review page shows who accepted
- [ ] Tests cover an acceptance, a refusal on an open thread, and the permission off

## Tasks

## Notes
