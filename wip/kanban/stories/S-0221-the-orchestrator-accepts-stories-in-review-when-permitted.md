---
id: S-0221
type: story
nature: feature
title: The orchestrator accepts stories in review when permitted
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-02T11:54:43Z
transitions: []
tags: [flai]
touches: [flai/internal/harness, flai/internal/hostapi, flai/internal/preview, flai/internal/mcpserver, design/adrs]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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
