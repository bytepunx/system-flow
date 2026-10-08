---
id: S-0336
type: story
nature: feature
title: The dashboard shows the messages between stories' agents on each story's page and in a Messages view
status: backlog
parent: E-0018
owner: alex
created: 2026-10-07T20:11:22Z
updated: 2026-10-08T04:33:38Z
transitions: []
tags: [flai, flaiover]
topics: [dashboard, cli]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/contract_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/repo.test.ts, flaiover/src/routes/api/messages/+server.ts, flaiover/src/routes/api/messages/messages.test.ts, flaiover/src/lib/components/Messages.svelte, flaiover/src/lib/components/Messages.svelte.test.ts, flaiover/src/routes/messages/+page.svelte, flaiover/src/routes/messages/messages.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/flai-cli.md]
after: [S-0330]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 435
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 56
          output: 909
          cache_read: 8224154
          cache_write: 15794
          cost: 2.0303
draft: true
cost_of_delay:
  value: 217.39
  by: planner-E-0018
  at: 2026-10-08T04:33:21Z
forecast:
  duration: 40m
  delivery: 2026-10-08T12:42:00Z
  basis: "Its own forecast of 40m; 29th in the pull order with an in-progress limit of 3, behind S-0232, S-0316, S-0324, S-0318, S-0320, S-0319, S-0309, S-0326, S-0287, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0313, S-0321, S-0322, S-0323, S-0327 and S-0334."
  by: flai
  at: 2026-10-08T04:33:38Z
---
# S-0336 The dashboard shows the messages between stories' agents on each story's page and in a Messages view

## Goal

Let the operator see how the agents are coordinating without being asked. Messages are kept apart from the operator's threads and inbox, so today they are visible only through `flai message list`. Serve them read-only through the host API and show them on each story's page and in a Messages view across the project. The conversations S-0330 to S-0335 already write are enough: this story waits for nothing still open. What a held card says about a hold request or a share is the story split from this one, which waits for S-0334.

## Acceptance criteria

- [ ] The host API reads `messages.list` and `messages.get` serve the conversations of the project or of one story, with their entries, their `about` paths, and their state.
- [ ] A story's page lists its conversations, open first, each with the other story linked and its entries.
- [ ] A Messages view, linked from the navigation, lists every open conversation with its two stories, its paths, which side it awaits, and its age, and the closed ones on request.
- [ ] The operator's inbox badge and inbox view count no message.
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` describe the view.

## Tasks

- T-1211 The host API reads messages.list and messages.get serve the conversations of a project or a story
- T-1212 flaiover's repo reads messages from the host, and /api/messages serves them
- T-1213 A Messages component and view list the conversations between stories, linked from the site menu
- T-1214 A story's page lists its conversations, and the inbox counts no message
- T-1215 The dashboard design and the flaiover guide describe the Messages view and the messages on a story

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07, and split on 2026-10-08 as the operator accepted on TH-0320: the held card's part went to its own draft story, which waits for S-0334 and for this one. This story now waits only for S-0330 (`after`), which is done, so it can be pulled as soon as it is finalized.

Layers, one task each, since each builds on the one before:

1. T-1211, the host API reads.
2. T-1212, flaiover's repo and route.
3. T-1213, the component, the view, and the menu.
4. T-1214, the story page and the inbox count.
5. T-1215, the docs.

Touches:

- **Declared:** kept, less the two `BoardCard.svelte` files, which went with the held card.
- **Layout:**
  - `flai/internal/hostapi/hostapi.go`, its test, and `contract_test.go`: `threads.list` is served there, and the contract lists every read. No `messages.*` read exists yet; `messages.View` in `flai/internal/messages/messages.go` already shapes a conversation for output.
  - `flaiover/src/lib/server/repo.ts` and its test: `threads()` and `threadsFor()`.
  - `flaiover/src/routes/api/messages/+server.ts` and its test, as `routes/api/threads` is.
  - `flaiover/src/lib/components/Messages.svelte`, `flaiover/src/routes/messages/+page.svelte`, and their tests, as `Threads.svelte` and `routes/threads` are; `flaiover/src/lib/sitemenu.ts` and its test for the menu.
- **Co-change:** `flai touches suggest` from the threads view and the story page gave `design/system/flaiover-dashboard.md` (50%), `docs/users/flaiover.md` (40%), `item.svelte.test.ts` (36%), and `repo.ts` (21%).
- **Not taken:** `flaiover/src/routes/board/+page.svelte` (29%) and `StoryAgent.svelte` (14%): the board is the split story's, and S-0335 already changed the agent's state.
- No folder touch.

Forecast and cost of delay: see the figures below, rewritten at the split.
