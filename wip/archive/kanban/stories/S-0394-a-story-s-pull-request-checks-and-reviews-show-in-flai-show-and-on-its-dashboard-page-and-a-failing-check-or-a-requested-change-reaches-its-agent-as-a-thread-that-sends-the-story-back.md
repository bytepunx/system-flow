---
id: S-0394
type: story
nature: feature
title: A story's pull request checks and reviews show in flai show and on its dashboard page, and a failing check or a requested change reaches its agent as a thread that sends the story back
status: done
parent: E-0021
owner: alex
created: 2026-10-16T04:36:28Z
updated: 2026-10-17T17:33:41Z
transitions:
  - to: ready
    at: 2026-10-17T06:43:51Z
    by: alex
  - to: in-progress
    at: 2026-10-17T16:03:23Z
    by: agent-S-0394
  - to: review
    at: 2026-10-17T17:32:00Z
    by: agent-S-0394
  - to: done
    at: 2026-10-17T17:33:41Z
    by: orchestrator
tags: [cli, dashboard]
topics: [release, dashboard]
touches: [flai/internal/forge/state.go, flai/internal/forge/state_test.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/cmd/items.go, flai/cmd/show_test.go, flai/cmd/forge.go, flai/cmd/forge_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, flaiover/src/lib/server/agent.ts, "flaiover/src/routes/api/items/[id]/forge/+server.ts", "flaiover/src/routes/api/items/[id]/forge/forge.test.ts", flaiover/src/lib/components/ForgeStatus.svelte, flaiover/src/lib/components/ForgeStatus.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", design/adrs, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, flai/internal/serve/forge.go, flai/internal/serve/forge_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/harness/adapters.go, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, design/components/flaiover/src/lib/components/ForgeStatus.examples.json, design/system/agent-adapters.md, docs/operators/index.md]
after: [S-0393]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5334
  priced_by: harness
  turns:
    - day: 2026-10-17
      ceremony: 3
      hand_edits: 3
      work: 77
  models:
    - model: claude-haiku-5-5
      input: 124
      output: 47706
      cache_read: 8429901
      cache_write: 214196
      cost: 0.5703
    - model: claude-opus-5-5
      input: 672
      output: 284146
      cache_read: 40682297
      cache_write: 1228957
      cost: 20.777
    - model: claude-sonnet-5-5
      input: 38
      output: 14047
      cache_read: 1303893
      cache_write: 99792
      cost: 0.5204
  strategic:
    - kind: orchestrator
      seconds: 51
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 152
          cache_read: 3580190
          cache_write: 3779
          cost: 0.8985
forecast:
  duration: 1h15m
  delivery: 2026-10-17T10:14:00Z
  basis: "Its own forecast of 1h15m; 5th in the pull order with an in-progress limit of 5, behind S-0392, S-0395, S-0396, S-0397 and S-0393."
  by: flai
  at: 2026-10-17T06:44:14Z
finalized:
  by: alex
  at: 2026-10-17T06:41:42Z
---
# S-0394 A story's pull request checks and reviews show in flai show and on its dashboard page, and a failing check or a requested change reaches its agent as a thread that sends the story back

## Goal

The story's agent hears from CI and from reviewers without the operator relaying. While a story's pull request is open, flai serve reads its required checks and its reviews. `flai show`, the MCP `item_get`, and the story's and review's pages in the dashboard show each one's state. When a required check fails, or a reviewer requests changes, flai sends the story back to in-progress with a thread that names the check or the comment and links its log. flai serve then starts the story's agent, which fixes the branch and moves the story to review again, and that pushes the fix.

## Acceptance criteria

- [x] For each story in review with a pull request, flai serve reads the request's required checks and reviews on its minute loop, through the forge contract, and keeps the latest state with the time of each change under `.flai-cache`, the state S-0392's contract returns.
- [x] `flai show S-nnnn`, its `--json`, and the MCP `item_get` carry the pull request's URL, each required check's name, state, and log link, and each review's author and state, read from that kept state. They never call the forge, so an agent without the token sees them too.
- [x] The dashboard's story page and its review page show the same, with a link to the pull request, through a host read method and a `/api/items/[id]/forge` route.
- [x] A required check that fails, or a review that requests changes, opens one thread on the story. The thread names the check and links its log, or quotes the review's comments and links them. It moves the story from review back to in-progress, by `flai`, with the reason naming the same. The same failure on the same commit never opens a second thread.
- [x] flai serve starts the story's agent for such a thread, as it does when a thread the agent waits on is answered, and its first `inbox` lists the thread. The agent fixes the branch, closes out, and moves the story to review, which pushes again under S-0393.
- [x] An ADR records the reading, the kept state, and the send back by flai on a failing check or a requested change. `design/system/workflow.md` (review to in-progress), `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, and the reference describe them.

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
- Criterion 5 is met for a story whose agent flai serve ran, as an answered thread's resume is. A story flai serve never ran an agent for is worked from a session of its own. flai serve starts nothing beside that session, whose inbox lists the thread. ADR-0167 does not say this; `workflow.md` and `flai-cli.md` do.
- The criteria were verified against the tests with a fake forge. No real GitHub pull request was read here.

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

### Accepted by the orchestrator

- Verified: 6bd4ad96615de4698ee8343732ace14314ba7486
- At: 2026-10-17T17:33:41Z

Verdict: all six criteria met at 6bd4ad9; flai verify passed at that commit against main 0cef465, integration and smoke included; readers never call the forge; one low-severity note: sendBackReason puts a check name and reviewer login unquoted into the agent's prompt, set only by those who can already push; no convention breach.
- 1: flai/internal/forge/state.go, flai/internal/serve/forgestatus.go, flai/internal/serve/forge.go, flai/internal/forge/state_test.go, flai/internal/serve/forgestatus_test.go
- 2: flai/cmd/items.go, flai/cmd/forge.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/cmd/show_test.go, flai/internal/mcpserver/server_test.go
- 3: flai/internal/hostapi/reads.go, flaiover/src/routes/api/items/[id]/forge/+server.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/components/ForgeStatus.svelte, flaiover/src/routes/items/[id]/+page.svelte, flaiover/src/lib/components/Review.svelte
- 4: flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go
- 5: flai/internal/serve/agents.go, flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/serve/agents_test.go
- 6: design/adrs/0167-flai-serve-reads-each-pull-request-s-required-checks-and-reviews-into-kept.md, design/adrs/README.md, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md
