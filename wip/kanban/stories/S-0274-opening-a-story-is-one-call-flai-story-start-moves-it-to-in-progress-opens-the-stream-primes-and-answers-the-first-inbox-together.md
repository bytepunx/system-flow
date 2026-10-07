---
id: S-0274
type: story
nature: improvement
title: "Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together"
status: ready
parent: E-0017
owner: alex
created: 2026-10-05T01:35:32Z
updated: 2026-10-07T03:13:45Z
transitions:
  - to: ready
    at: 2026-10-06T23:59:35Z
    by: alex
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions, template]
touches: [flai/cmd/items.go, flai/cmd/story_start.go, flai/cmd/story_start_test.go, flai/cmd/move.go, flai/cmd/stream.go, flai/cmd/prime.go, flai/cmd/branch.go, flai/internal/storygit/open.go, flai/internal/storygit/open_test.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/inbox.go, flai/internal/mcpserver/inbox_test.go, flai/internal/mcpserver/story_start.go, flai/internal/mcpserver/story_start_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/session-start.md, design/conventions/work-management.md, template/root/design/conventions/session-start.md, template/root/design/conventions/work-management.md, CLAUDE.md, template/root/CLAUDE.md.tmpl, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0261]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 64
  by: planner-E-0017
  at: 2026-10-06T11:36:17Z
forecast:
  duration: 60m
  delivery: 2026-10-07T04:16:00Z
  basis: "Its own forecast of 1h; 3rd in the pull order with an in-progress limit of 3, behind S-0269, S-0270 and S-0271."
  by: flai
  at: 2026-10-07T03:13:45Z
finalized:
  by: alex
  at: 2026-10-06T22:48:57Z
---
# S-0274 Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together

## Goal

A story agent's first minute is four turns: `flai move S-nnnn in-progress`, `flai stream open`, `prime`, and `inbox`, and in S-0248 the prime's result overflowed the harness and was read back from a file (S-0261 fixes the size). `flai story start S-nnnn`, `story_start` over MCP, and `story.start` on the host channel do the four in one call and answer the worktree path, the branch, the prime pack within its budget, and the inbox. flai serve's prompt tells the agent to begin with it.

## Acceptance criteria
- [ ] `flai story start S-nnnn` moves the story to in-progress, opens the stream, and answers the worktree, the branch, the prime pack, and the inbox in one result, as text and `--json`, refusing a story that is not ready or is held
- [ ] The same is `story_start` over MCP and `story.start` on the host channel
- [ ] The harness prompt, `design/conventions/session-start.md`, the template's copies, `design/system/flai-cli.md`, and the user guide begin the loop with it

## Tasks
- T-1091 Opening a story's branch and worktree is a storygit function that the CLI and the MCP server can both call
- T-1099 The inbox is an exported mcpserver function over the agent's on-disk cursor, so the CLI can answer it as MCP does
- T-1103 Package storystart moves a ready story to in-progress, opens its stream, and builds its prime pack in one function
- T-1109 flai story start S-nnnn answers the worktree, the branch, the prime pack, and the inbox in one result, as text and --json
- T-1111 The MCP tool story_start starts a story and answers its worktree, branch, pack, and inbox, and the server's instructions pull with it
- T-1115 The host channel's write method story.start runs flai story start --json and answers its result
- T-1117 The story agent's start prompt, session-start.md, work-management.md, CLAUDE.md, and the template's copies begin the loop with flai story start
- T-1118 flai-cli.md, the user guide, and the reference describe flai story start, story_start, and story.start

## Notes

Depends on S-0261 for the prime pack to fit the harness's tool result limit when returned inline.

### Planning

The 20 touches planner-E-0017 declared are all kept. `flai touches suggest S-0274`, run from them, ranks nothing above 19% (`design/system/flaiover-dashboard.md`, `docs/operators/index.md`), and none of it is reached by the command. Where each touch comes from:

- `flai/cmd/items.go`, `flai/cmd/story_start.go`, `flai/cmd/story_start_test.go`: layout. The `story` command is built in `newItemCmd` in `items.go`, and the subcommand and its tests are new files (T-1109).
- `flai/cmd/move.go`, `flai/cmd/prime.go`: layout. T-1109 shares `flai move`'s refusal wording and prints the pack through `printPack`.
- `flai/cmd/stream.go`, `flai/cmd/branch.go`, `flai/internal/storygit/open.go`, `open_test.go`: layout. Branch and worktree opening is `(*app).openStoryBranch` in `branch.go`, which the MCP server cannot call. T-1091 moves it into `storygit`. `branch.go` and the two `storygit` files are added by this plan.
- `flai/internal/storystart/start.go`, `start_test.go`: layout, added by this plan. This is a new package for the composition the three entry points share (T-1103). It cannot live in `cmd`, which MCP cannot import, or in `mcpserver`, which `cmd` would then need for more than the inbox.
- `flai/internal/mcpserver/server.go`, `folder.go`: co-change and layout. The tools and the server's instructions are here (T-1099, T-1111).
- `flai/internal/mcpserver/inbox.go`, `inbox_test.go`, `cursor.go`: layout, added by this plan. The inbox is computed only in `(*server).inbox`, over a per-agent cursor file in `cursor.go`. T-1099 exports it, so that `flai story start` advances the same cursor.
- `flai/internal/mcpserver/story_start.go`, `story_start_test.go`: layout, added by this plan. These are the tool's handler and its tests (T-1111).
- `flai/internal/hostapi/writes.go`, `writes_test.go`: layout. `story.start` is a write method that runs the CLI, as `item.move` does (T-1115).
- `flai/internal/harness/harness.go`, `harness_test.go`: design (criterion 3, T-1117).
- `design/conventions/session-start.md`, `work-management.md`, their `template/root` copies, and `template/CHANGELOG.md`: design (criterion 3, T-1117).
- `CLAUDE.md`, `template/root/CLAUDE.md.tmpl`: layout, added by this plan. The start prompt says to follow `CLAUDE.md`, whose "Prime your session" step 3 pulls with `flai stream open` (T-1117).
- `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`: co-change and criterion 3 (T-1118).
- No folder touches.
- Left out:
  - `flai/internal/context` and `flai/internal/conventions`: `ctxpack.ForStory` is reused as it is.
  - `flai/internal/workitem`: its transition and holds are called, not changed.
  - `flai/internal/guard`: `story start` and `story_start` are not in its read lists, so sub-agents stay refused them.
  - `design/system/dashboard-host-channel.md`: it lists no method by name.

Forecast: 60m, against flai's 49m (83 s per unit of size over 25 done large improvement stories, times size 35 with the widened touches). I raised it because the steps do not compose as they stand: stream open must first leave package `cmd` and the inbox the MCP handler. S-0217, which composed commands across the CLI, MCP, and the host channel, took 68m. The delivery is flai's 2026-10-07T09:50Z plus the added 11m. It comes after S-0261, which this story waits for.

Cost of delay: 64 USD a week, against flai's 72.75. The value is planner-E-0017's and stands. It shares E-0017's 900 USD a week by the turns each story removes, not by forecast duration. This story removes about 432: three of four start turns in each of 108 runs, plus the read-back of an overflowed prime. The longer forecast changes neither figure.
