---
id: S-0061
type: story
nature: remediation
title: An agent's first inbox is bounded, not the last 24 hours in full
status: backlog
parent: E-0006
owner: alex
created: 2026-09-19T05:44:37Z
updated: 2026-09-19T05:44:37Z
transitions: []
tags: [cli]
touches: [flai/internal/mcpserver]
---

# S-0061 An agent's first inbox is bounded, not the last 24 hours in full

## Goal
The first `inbox` an agent calls under a new name returns something it can read. Today it reports every change of the last 24 hours, which on a busy repository is too large for a tool result.

## Acceptance criteria
- [ ] `changes` in `inbox` and `events` in `wait_for_events` are capped at a stated number, newest kept, with `changes_omitted` (or the like) saying how many older ones were left out; the cap is a constant with a reason, not a setting
- [ ] A first look, with no cursor, reports story and epic changes only, not task transitions, and no more than the cap; ready work and threads are unaffected, since they are state
- [ ] After any look the cursor advances past everything, reported or omitted, so nothing omitted comes back later
- [ ] Tests: a repository with more changes than the cap, with and without a cursor; the server instructions and `docs/users/flai.md` say what a first look returns

## Tasks

## Notes
Found on 2026-09-19 in the first session whose MCP server was flai 1.2.9: the first `inbox` as `system-flow` returned 209 changes, 68,283 characters, and the harness saved it to a file because it exceeded what a tool result may hold. The second call returned no changes and the same three ready stories, as designed. 155 of the 209 were task transitions, and 208 were recorded as `alex`, because until S-0058 an agent's moves from the command line carried the config author's name.

This is a defect in S-0058 (its author's): "with no cursor, the last 24 hours are reported" was written without a bound. `catchUp` in `flai/internal/mcpserver/cursor.go` is where the cap belongs.

Record the issue on this story's branch when it is pulled; it was not recorded on main when queued, to keep the main checkout clean outside `wip/`.
