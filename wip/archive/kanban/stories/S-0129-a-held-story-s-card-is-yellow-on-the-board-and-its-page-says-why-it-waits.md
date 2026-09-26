---
id: S-0129
type: story
nature: feature
title: A held story's card is yellow on the board and its page says why it waits
status: done
parent: E-0009
owner: alex
created: 2026-09-26T07:59:21Z
updated: 2026-09-26T18:12:29Z
transitions:
  - to: ready
    at: 2026-09-26T17:50:20Z
    by: alex
  - to: in-progress
    at: 2026-09-26T17:50:35Z
    by: agent-S-0129
  - to: review
    at: 2026-09-26T18:03:52Z
    by: agent-S-0129
  - to: done
    at: 2026-09-26T18:12:29Z
    by: alex
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

- [x] A held ready story's card on the board shows a yellow dot and a short line with the reason code and the story it waits for, whether or not the story ever had an agent.
- [x] The story's page shows the full reason from flai, with each story it names linked to that story's page.
- [x] When the hold clears, the card and page stop showing it without a reload, as other agent states do.
- [x] A story that is both blocked and held shows `BLOCKED` first and the hold second.
- [x] `design/system/flaiover-dashboard.md` describes it; flaiover's tests and lint pass; it was tried in a browser against a flai serve with two scratch stories whose touches overlap.

## Tasks
- T-0470 Show a held story's yellow dot and a short hold line on its board card
- T-0471 Say on a held story's page why it waits, linking each story the reason names
- T-0472 Document the held card and try it in a browser against a scratch flai serve

## Notes

- The reason comes from flai's `agent.status` as built by the previous story under E-0009; the yellow dot is `AgentDot.svelte` and the card is `BoardCard.svelte`.
- The card cannot link the story it waits for: the whole card is one link, and a link inside a link is not valid HTML. The card names it; the page links it (T-0470).
- flai's hold carries no list of the stories it waits for, so the dashboard reads them from the reason's `starts when` clause (`holdWaitsFor`). A structured field in flai would be firmer; it is outside this story's claim.
- Tried on 2026-09-26 (T-0472) in headless Chrome against a flaiover dev server built from `story/S-0129` on port 5199 and a scratch flai host and serve from the same branch, with their own `HOME` and configuration, `FLAI_HOST_ADDR=127.0.0.1:4299`, and no `FLAI_HOST_*` or `FLAI_CONFIG` in their environment. The `agent` action was off. The scratch project had S-0002 in ready touching `app` and an open story in progress touching `app/api`. Every scratch process was stopped by PID; the operator's host and dashboard were untouched.
  - The card showed a yellow dot, labelled with flai's whole reason, and `held (overlap): S-0004` in the details row; blocked too, it read `BLOCKED held (overlap): S-0004`.
  - The page's header read `ready BLOCKED HELD`, and the side column gave `held (overlap): touches app, which holds app/api that S-0004 (in progress) touches; starts when S-0004 is accepted, cancelled, or sent back` with both `S-0004` linked to `/items/S-0004`.
  - Cancelling the open story from the command line cleared the card's dot and line and the page's `HELD` and reason within three seconds, without a reload.
  - The item page does not reload the item itself on a change event, so a block made while it is open shows only on reload; that predates this story and is not the hold.
