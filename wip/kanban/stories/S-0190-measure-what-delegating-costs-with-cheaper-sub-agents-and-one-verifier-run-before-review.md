---
id: S-0190
type: story
nature: research
title: Measure what delegating costs with cheaper sub-agents and one verifier run before review
status: ready
owner: arobson
created: 2026-10-01T10:52:53Z
updated: 2026-10-01T11:46:27Z
transitions:
  - to: ready
    at: 2026-10-01T11:46:27Z
    by: alex
tags: [cli]
topics: [conventions]
touches: [design/system/agent-context.md]
after: [S-0189]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0190 Measure what delegating costs with cheaper sub-agents and one verifier run before review

## Goal

S-0189 ran the explorer on `haiku` and the verifier on `sonnet`, let roles carry their own model (`agent.roles`), and had the verifier's run replace the story's agent's own full-suite runs. Measure whether delegation now costs less than not delegating, the way S-0188 measured S-0184 and S-0185 (`design/system/agent-context.md § Sub-agents › Measured`).

## Acceptance criteria

- [ ] Two stories run by `flai serve` with a flai that includes S-0189 are measured as S-0188 measured S-0184 and S-0185, with the same reading of the logs and against the same non-delegating comparables
- [ ] The table in `agent-context.md § Sub-agents › Measured` gains their rows and shows whether the cost of delegation fell below the cost of not delegating, with each sub-agent's model and share of the cost
- [ ] The section says what to change next, if anything

## Tasks

## Notes

- Split from S-0189 criterion 4 on the designer's word (TH-0057): it can be measured only after S-0189 is accepted, flai is released and installed, and two more stories have run.
- Pull it once two stories have run with the release that includes S-0189.
