---
id: S-0274
type: story
nature: improvement
title: "Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together"
status: backlog
parent: E-0017
owner: alex
created: 2026-10-05T01:35:32Z
updated: 2026-10-06T11:52:17Z
transitions: []
tags: [cli, mcp]
topics: [automation, mcp, hostapi, conventions]
touches: [flai/cmd/items.go, flai/cmd/story_start.go, flai/cmd/story_start_test.go, flai/cmd/move.go, flai/cmd/stream.go, flai/cmd/prime.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/session-start.md, design/conventions/work-management.md, template/root/design/conventions/session-start.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0261]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  value: 64
  by: planner-E-0017
  at: 2026-10-06T11:36:17Z
forecast:
  duration: 35m
  delivery: 2026-10-07T01:55:00Z
  basis: "Its own forecast of 35m; 38th in the pull order with an in-progress limit of 3, behind S-0222, S-0224, S-0226, S-0223, S-0227, S-0229, S-0284, S-0278, S-0295, S-0296, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0261, S-0264, S-0265, S-0269, S-0270, S-0271, S-0272 and S-0273."
  by: flai
  at: 2026-10-06T11:52:17Z
---
# S-0274 Opening a story is one call: flai story start moves it to in-progress, opens the stream, primes, and answers the first inbox together

## Goal

A story agent's first minute is four turns: `flai move S-nnnn in-progress`, `flai stream open`, `prime`, and `inbox`, and in S-0248 the prime's result overflowed the harness and was read back from a file (S-0261 fixes the size). `flai story start S-nnnn`, `story_start` over MCP, and `story.start` on the host channel do the four in one call and answer the worktree path, the branch, the prime pack within its budget, and the inbox. flai serve's prompt tells the agent to begin with it.

## Acceptance criteria
- [ ] `flai story start S-nnnn` moves the story to in-progress, opens the stream, and answers the worktree, the branch, the prime pack, and the inbox in one result, as text and `--json`, refusing a story that is not ready or is held
- [ ] The same is `story_start` over MCP and `story.start` on the host channel
- [ ] The harness prompt, `design/conventions/session-start.md`, the template's copies, `design/system/flai-cli.md`, and the user guide begin the loop with it

## Tasks

## Notes

Depends on S-0261 for the prime pack to fit the harness's tool result limit when returned inline.

### Planning

Touches, none declared before. `flai touches suggest S-0274` was seeded with `flai/cmd/prime.go` and `flai/cmd/stream.go`, which 21 commits changed:

- `flai/cmd/items.go`, `flai/cmd/story_start.go`, `flai/cmd/story_start_test.go`: layout. The story command is built in `items.go`; the subcommand and its tests are new files.
- `flai/cmd/move.go`, `flai/cmd/stream.go`, `flai/cmd/prime.go`: layout. Moving, opening the stream, and priming live here as command code, which must become callable together.
- `flai/internal/mcpserver/folder.go`, `server.go`: co-change (`folder.go` 5 of 21, `server.go` 4 of 21). `inbox` and the tools are registered here.
- `flai/internal/hostapi/writes.go`, `writes_test.go`: layout. `story.start` is a write method.
- `flai/internal/harness/harness.go`, `harness_test.go`: design (criterion 3).
- `design/conventions/session-start.md`, `work-management.md`, their `template/root` copies, and `template/CHANGELOG.md`: design. `work-management.md` lists the order of steps for pulling a story.
- `design/system/flai-cli.md`, `docs/users/flai.md`, `flai-reference.md`: co-change (8, 7, and 4 of 21).
- Left out: `flai/internal/context` and `flai/internal/conventions`, co-changed with prime (5 and 4 of 21). They build the pack, which this story reuses as it is, and S-0261 changes its size.

Forecast: flai gave 35m (89 s per unit over 21 done large improvement stories, times size 23), and it stands: the four steps exist and are composed here. The delivery, 2026-10-07T01:51Z, is flai's. It comes after S-0261's, which this story waits for.

Cost of delay: 64 USD a week, against flai's 92.11. This is E-0017's 900 USD a week shared by the turns each story removes. This one removes about 432: three of four start turns in each of 108 runs, plus the read-back of an overflowed prime.
