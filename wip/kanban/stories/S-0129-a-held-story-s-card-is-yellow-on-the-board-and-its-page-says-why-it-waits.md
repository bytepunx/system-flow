---
id: S-0129
type: story
nature: feature
title: A held story's card is yellow on the board and its page says why it waits
status: in-progress
parent: E-0009
owner: alex
created: 2026-09-26T07:59:21Z
updated: 2026-09-26T17:50:35Z
transitions:
  - to: ready
    at: 2026-09-26T17:50:20Z
    by: alex
  - to: in-progress
    at: 2026-09-26T17:50:35Z
    by: agent-S-0129
tags: [dashboard]
touches: [flaiover]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0129 A held story's card is yellow on the board and its page says why it waits

## Goal

Show the operator the hold of [ADR-0046](../../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) on the dashboard: a held ready story's card is yellow, and the card and the story's page say why it waits and what clears it, naming the story it waits for as a link.

## Acceptance criteria

- [ ] A held ready story's card on the board shows a yellow dot and a short line with the reason code and the story it waits for, whether or not the story ever had an agent.
- [ ] The story's page shows the full reason from flai, with each story it names linked to that story's page.
- [ ] When the hold clears, the card and page stop showing it without a reload, as other agent states do.
- [ ] A story that is both blocked and held shows `BLOCKED` first and the hold second.
- [ ] `design/system/flaiover-dashboard.md` describes it; flaiover's tests and lint pass; it was tried in a browser against a flai serve with two scratch stories whose touches overlap.

## Tasks

## Notes

- The reason comes from flai's `agent.status` as built by the previous story under E-0009; the yellow dot is `AgentDot.svelte` and the card is `BoardCard.svelte`.
