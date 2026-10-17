---
id: S-0394
type: story
nature: feature
title: A story's pull request checks and reviews show in flai show and on its dashboard page, and a failing check or a requested change reaches its agent as a thread that sends the story back
status: backlog
parent: E-0021
owner: alex
created: 2026-10-16T04:36:28Z
updated: 2026-10-17T03:34:53Z
transitions: []
tags: [cli, dashboard]
topics: [release, dashboard]
touches: [flai/internal/forge/state.go, flai/internal/forge/state_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/serve.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/cmd/items.go, flai/cmd/show_test.go, flai/cmd/forge.go, flai/cmd/forge_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, flaiover/src/lib/server/agent.ts, "flaiover/src/routes/api/items/[id]/forge/+server.ts", "flaiover/src/routes/api/items/[id]/forge/forge.test.ts", flaiover/src/lib/components/ForgeStatus.svelte, flaiover/src/lib/components/ForgeStatus.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", "flaiover/src/routes/review/[id]/+page.svelte", design/adrs, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
after: [S-0393]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h15m
  delivery: 2026-10-17T07:12:00Z
  basis: "Its own forecast of 1h15m; 6th in the pull order with an in-progress limit of 5, behind S-0442, S-0444, S-0443, S-0241, S-0392 and S-0393."
  by: flai
  at: 2026-10-17T03:34:53Z
---
# S-0394 A story's pull request checks and reviews show in flai show and on its dashboard page, and a failing check or a requested change reaches its agent as a thread that sends the story back

## Goal

The story's agent hears from CI and from reviewers without the operator relaying. While a story's pull request is open, flai serve reads its required checks and its reviews. `flai show`, the MCP `item_get`, and the story's and review's pages in the dashboard show each one's state. When a required check fails, or a reviewer requests changes, flai sends the story back to in-progress with a thread that names the check or the comment and links its log. flai serve then starts the story's agent, which fixes the branch and moves the story to review again, and that pushes the fix.

## Acceptance criteria

- [ ] For each story in review with a pull request, flai serve reads the request's required checks and reviews on its minute loop, through the forge contract, and keeps the latest state with the time of each change under `.flai-cache`, the state S-0392's contract returns.
- [ ] `flai show S-nnnn`, its `--json`, and the MCP `item_get` carry the pull request's URL, each required check's name, state, and log link, and each review's author and state, read from that kept state. They never call the forge, so an agent without the token sees them too.
- [ ] The dashboard's story page and its review page show the same, with a link to the pull request, through a host read method and a `/api/items/[id]/forge` route.
- [ ] A required check that fails, or a review that requests changes, opens one thread on the story. The thread names the check and links its log, or quotes the review's comments and links them. It moves the story from review back to in-progress, by `flai`, with the reason naming the same. The same failure on the same commit never opens a second thread.
- [ ] flai serve starts the story's agent for such a thread, as it does when a thread the agent waits on is answered, and its first `inbox` lists the thread. The agent fixes the branch, closes out, and moves the story to review, which pushes again under S-0393.
- [ ] An ADR records the reading, the kept state, and the send back by flai on a failing check or a requested change. `design/system/workflow.md` (review to in-progress), `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, and the reference describe them.

## Tasks
- T-1615 flai keeps each story's pull request state, its checks and reviews and when each changed, under .flai-cache, read without the forge
- T-1616 An ADR records reading a pull request's checks and reviews into kept state, and flai sending a story back to in-progress with a thread on a failing check or a requested change
- T-1617 flai serve reads each pull request's checks and reviews into the kept state, and on a failing check or a requested change opens one thread and sends the story back
- T-1618 flai show, the MCP item_get, and flai forge status carry a story's pull request with each required check's and review's state, read from the kept state
- T-1619 A host read method gives the dashboard a story's pull request, checks, and reviews from the kept state
- T-1620 flai serve starts a story's agent when flai sent the story back with a thread on a failing check or a requested change
- T-1621 The dashboard's story and review pages show a story's pull request, its required checks, and its reviews
- T-1622 The workflow, the CLI and dashboard design, and the user docs describe the checks and reviews read back and the send back on a failure

## Notes

- Drafted by the planner for E-0021, after S-0393, which opens the pull request this story reads.
- Sending the story back is a new maker of the review to in-progress transition, today the operator's alone. The planner assumes this rather than leaving the story in review while its agent changes the branch, so that a fix goes through the close-out and the move to review as every other change does. Raised on the plan's thread.
- The kept state is what S-0396 reads to measure the wait on checks and reviewers.

### Planning

- **Touches.** None were declared. Every touch is the planner's, and each is a file except `design/adrs`:
  - **Layout:**
    - `flai/internal/forge/state.go`, with its test: the kept state and its events. Criterion 1.
    - `flai/internal/serve/forgestatus.go`, with its test, and `serve.go`: the reading, the thread, and the send back. Criteria 1 and 4.
    - `flai/internal/serve/agents.go`, with its test: the start beside `answered`. Criterion 5.
    - `flai/cmd/items.go` with `show_test.go`, `flai/cmd/forge.go` with its test, and `flai/internal/mcpserver/server.go` with its test: criterion 2.
    - `flai/internal/hostapi/reads.go`, with its test, `flaiover/src/lib/server/agent.ts`, the new route and `ForgeStatus.svelte` with their tests, and the item and review pages: criterion 3.
  - **Co-change:** `flai/internal/mcpserver/folder.go`, at 8%, where `item_get`'s description lives. It should name the new field.
  - **Design:** the workflow, the CLI and dashboard design, and the user docs and reference, for criterion 6.
- **Folder touch kept.** `design/adrs`, a shared path, for an ADR `flai adr new` names.
- **Co-change left out.**
  - `docs/operators/index.md` and `settings.md`: this story adds no setting.
  - `flai/internal/hostapi/writes.go`: it adds a read, not a write.
- **Layers.**
  - First: T-1615 and T-1616 together.
  - Second: T-1617, T-1618, and T-1619, each after T-1615 alone.
  - Third: T-1620, after T-1617; and T-1621, after T-1619.
  - Fourth: T-1622, after T-1616, T-1618, T-1620, and T-1621.
- **Forecast: 1h15m.** `flai forecast` gave 59m: 96 s per unit over 61 done feature stories on claude-opus-5-5 in the large band, times size 37. Raised about a quarter, as S-0384 was, for spanning flai and flaiover with both test tiers and a new Svelte component, across eight tasks in four layers. flai's replan sets the delivery from this duration.
- **Cost of delay: not yet set.** `flai cod` refuses until E-0021 has inputs, asked for on TH-0463.
