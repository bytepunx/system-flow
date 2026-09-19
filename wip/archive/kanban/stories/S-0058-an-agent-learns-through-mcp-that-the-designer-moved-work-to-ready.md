---
id: S-0058
type: story
nature: remediation
title: An agent learns through MCP that the designer moved work to ready
status: done
parent: E-0006
owner: alex
created: 2026-09-19T02:09:43Z
updated: 2026-09-19T04:04:39Z
transitions:
  - to: ready
    at: 2026-09-19T02:42:42Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:44:49Z
    by: alex
  - to: review
    at: 2026-09-19T03:07:17Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:04:39Z
    by: alex
tags: [cli]
touches: [flai/internal/mcpserver, design/conventions, template]
---

# S-0058 An agent learns through MCP that the designer moved work to ready

## Goal
When the designer moves a story to ready, or makes any other change to the board, an agent connected through `flai mcp` finds out the next time it looks, without having had to be waiting at that moment and without shelling out to `flai board`. Today the move is invisible to it.

## Acceptance criteria
- [x] `inbox` reports work as well as threads: the stories that are ready to pull, in pull order, with whether the WIP limit allows a pull, and the items the designer changed since this agent last looked (moved, blocked or unblocked, reordered), each with what happened, to what, by whom, and when
- [x] "Since this agent last looked" survives between calls and between sessions: the server keeps a cursor per agent (`FLAI_AGENT`) outside git, under `.flai-cache`, and a change made while no call was in flight is reported by the next `inbox` or `wait_for_events`, once
- [x] `wait_for_events` reports events, not only paths: for a work item, the ID, the transition or block, and who made it; it accepts the cursor so nothing between two calls is lost, and it still returns within a second of a change while it is held
- [x] A `board` tool returns what `flai board --json` does, so an agent can see the columns, limits, pull order, and overlaps over MCP
- [x] The agent's own writes through MCP or the CLI are not reported back to it as the designer's changes
- [x] The server instructions, the `CLAUDE.md` template, and the work-management and session-start conventions (template baseline first) say when to call what for an agent that ends its turn between messages as well as for one that holds `wait_for_events`, and that a ready story found in `inbox` is pulled without waiting to be told
- [x] Tests in `flai/internal/mcpserver` with the in-memory transports: a move to ready made between two calls appears in the next `inbox` and not in the one after; a held `wait_for_events` returns the event; the agent's own move does not appear; the board tool matches `flai board --json`. `design/system/flai-cli.md`, ADR-0020's living design, and `docs/users` updated

## Tasks
- T-0179 Living design: what inbox, wait_for_events, and board report, and the per-agent cursor
- T-0180 Work item changes since a cursor, and the board view shared by the CLI and MCP
- T-0181 MCP: inbox reports ready work and the designer's changes, wait_for_events reports events from the cursor, and a board tool
- T-0182 Moves and blocks are attributed to the agent that made them, and to the designer from the dashboard
- T-0183 Conventions, server instructions, and documentation for agents that end their turn
- T-0184 Verify: all tiers, and the new server over stdio against a scratch project

## Notes
Raised by the operator on 2026-09-19: "it doesn't appear the MCP is correctly reflecting when I have moved cards to ready". They had moved S-0040 to S-0043 to ready from the board; twelve minutes later the agent in the session still had not noticed, and only saw it by running `flai board` while queuing this story.

What was checked. The server's data is live: `item_get` on S-0052 returned the block the operator had just set. The gap is in what the tools report and when, in `flai/internal/mcpserver/server.go`:

- `inbox` lists threads and nothing else (`InboxOut` is agent, awaiting count, threads). The conventions tell the agent to call `inbox` at session start, at every task transition, and before review, and it did, each time getting an empty result while work had become ready.
- `wait_for_events` takes its baseline (`snapshot`) when the call starts and compares against it while the call is held, so anything that changed between two calls is never reported; and it returns changed file paths (`WaitOut.Changed`), so even a held call says `wip/kanban/stories/S-0040-....md` changed, not that S-0040 was moved to ready by alex.
- There is no board or ready-work tool; the tools are `doc_get`, `inbox`, `item_get`, `item_move`, the four thread tools, `wait_for_events`, and `who_touches`.

Part of this is how the agent was working, not only the server: the conventions say to hold `wait_for_events` when idle, and an agent in a conversational session ends its turn when it has reported, so it holds nothing between the operator's messages. The design in ADR-0020 assumed an agent that waits. The fix has to serve both, which is why the cursor matters more than the blocking call.

The MCP server in `.mcp.json` runs the installed `flai` from `PATH`, which on the operator's machine is upgraded only with sudo and was 1.2.3 when this was found; the behaviour described is the same in 1.2.4. When verifying the fix, check which binary the session's server is.

Not recorded under `design/issues` when queued, to keep the main checkout clean outside `wip/` for acceptances from the board; the agent that pulls this records the issue on its branch.

Not in the board `order`; the operator has not placed it.

Decided when pulled, 2026-09-19. Changes are derived from the files and nothing new is recorded in git; the cursor is a read marker under `.flai-cache/mcp/`. Item creation is not reported, because the files do not say who created an item and an agent's own new stories would come back to it. A blocked interval has no author field, and adding one would make older flai binaries refuse the item, so a block is reported to every agent, including the one that set it; this is the one place the fifth criterion is met only in part, for moves and not for blocks. A change to the pull order is reported only when the stories two orders share come in a different sequence, since `flai move` itself adds and removes stories from the order.

Verification, 2026-09-19. `make flai-test` passed (golangci-lint 0 issues; behavior, integration, smoke; markdown lint) and flaiover lint, svelte-check 0 errors, 20 files and 131 tests passed. The branch's `flai mcp` was driven over stdio by a JSON-RPC client against a scratch project, with a second process without `FLAI_AGENT` acting as the designer: a fresh `inbox` listed nothing; after `flai move S-0001 ready` it listed the story under `ready` and reported "moved to ready by alex" once; the next `inbox` reported no change and still listed it; the agent's `item_move` to in-progress was recorded as `claude` and was not echoed; a held `wait_for_events` returned the designer's block 1.3 seconds after it was made, with the changed path; a `wait_for_events` called after an unblock returned at once with it; `board` returned the columns and limits; the cursor file was `.flai-cache/mcp/claude.json`. Not verified in the operator's own session: its MCP server is the installed `flai` from `PATH`, which was 1.2.3 and is upgraded only with sudo, so this story reaches their sessions when that binary is upgraded.
