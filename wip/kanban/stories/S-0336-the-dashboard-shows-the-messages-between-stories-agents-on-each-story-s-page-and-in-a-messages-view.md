---
id: S-0336
type: story
nature: feature
title: The dashboard shows the messages between stories' agents on each story's page and in a Messages view
status: ready
parent: E-0018
owner: alex
created: 2026-10-07T20:11:22Z
updated: 2026-10-08T04:49:57Z
transitions:
  - to: ready
    at: 2026-10-08T04:35:26Z
    by: orchestrator
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
      seconds: 444
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 62
          output: 975
          cache_read: 10762158
          cache_write: 22429
          cost: 2.6572
cost_of_delay:
  value: 217.39
  by: planner-E-0018
  at: 2026-10-08T04:33:21Z
forecast:
  duration: 40m
  delivery: 2026-10-08T05:55:00Z
  basis: "Its own forecast of 40m; 3rd in the pull order with an in-progress limit of 3, behind S-0232, S-0318, S-0320, S-0319 and S-0309."
  by: flai
  at: 2026-10-08T04:49:57Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:35:22Z
---
# S-0336 The dashboard shows the messages between stories' agents on each story's page and in a Messages view

## Goal

Let the operator see how the agents are coordinating without being asked. Messages are kept apart from the operator's threads and inbox, so today they are visible only through `flai message list`. Serve them read-only through the host API and show them on each story's page and in a Messages view across the project. The conversations S-0330 to S-0335 already write are enough: this story waits for nothing still open. What a held card says about a hold request or a share is S-0338, split from this one, which waits for S-0334.

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

Planned by planner-E-0018 on 2026-10-07, and split on 2026-10-08 as the operator accepted on TH-0320: the held card's part went to S-0338, which waits for S-0334 and for this story. This story now waits only for S-0330 (`after`), which is done, so it can be pulled as soon as it is finalized.

Layers, one task each, since each builds on the one before:

1. T-1211, the host API reads.
2. T-1212, flaiover's repo and route.
3. T-1213, the component, the view, and the menu.
4. T-1214, the story page and the inbox count.
5. T-1215, the docs.

Touches:

- **Declared:** kept, less the two `BoardCard.svelte` files, which went to S-0338.
- **Layout:**
  - `flai/internal/hostapi/hostapi.go`, its test, and `contract_test.go`: `threads.list` is served there, and the contract lists every read. No `messages.*` read exists yet; `messages.View` in `flai/internal/messages/messages.go` already shapes a conversation for output.
  - `flaiover/src/lib/server/repo.ts` and its test: `threads()` and `threadsFor()`.
  - `flaiover/src/routes/api/messages/+server.ts` and its test, as `routes/api/threads` is.
  - `flaiover/src/lib/components/Messages.svelte`, `flaiover/src/routes/messages/+page.svelte`, and their tests, as `Threads.svelte` and `routes/threads` are; `flaiover/src/lib/sitemenu.ts` and its test for the menu.
- **Co-change:** `flai touches suggest` gave `design/system/flaiover-dashboard.md` (50%), `docs/users/flaiover.md` (40%), `item.svelte.test.ts` (36%), and `repo.ts` (21%) at the first plan; at the split it gave `design/system/flai-cli.md` (37%), which lists the host API's reads beside `threads.list`, so T-1215 and the story now touch it.
- **Not taken:** `docs/users/flai.md` (34%) and `docs/operators/index.md` (27%): no command and no setting changes. `flai/internal/hostapi/writes.go` (10%): the reads write nothing. `flaiover/src/routes/board/+page.svelte` and `StoryAgent.svelte`: the board is S-0338's, and S-0335 already changed the agent's state.
- No folder touch.

Forecast 40m: `flai forecast` gave it, 104 s per unit over 33 done large-band feature stories, times size 23 (5 criteria, 18 touches). It stands; the split took 8m off the 48m of the first plan.

Cost of delay 217.39 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 40m of the 3h4m forecast over the four open stories (S-0334, S-0336, S-0337, S-0338). It stands. It is the first of them that can be pulled, since it waits for nothing open.
