---
id: T-0184
type: task
nature: remediation
title: "Verify: all tiers, and the new server over stdio against a scratch project"
status: done
parent: S-0058
owner: alex
created: 2026-09-19T02:46:21Z
updated: 2026-09-19T03:06:23Z
transitions:
  - to: ready
    at: 2026-09-19T02:56:17Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T02:56:17Z
    by: system-flow
  - to: done
    at: 2026-09-19T03:06:23Z
    by: system-flow
stream: S-0058
tags: []
---

# T-0184 Verify: all tiers, and the new server over stdio against a scratch project

## Work
Run `make flai-test` and the flaiover checks. Drive the branch's `bin/flai mcp` over stdio with a small JSON-RPC script against a scratch project: `inbox` on a fresh project, then from a second process move a story to ready as someone else, and see `inbox` list it under `ready` and `changes` once; pull it with `item_move` and see no echo; hold `wait_for_events`, block the story from the second process, and see the event. Record what ran. Note for the operator that the session's MCP server is the installed `flai` from `PATH`, so this story reaches their sessions only when that binary is upgraded, which needs sudo. Record the issue this story fixes on the branch with `flai issue new`, since it was deliberately not recorded on main when queued. Tick the story's criteria for what was observed.

## Done when
- All tiers pass, with results in the narrative
- The stdio run is recorded with its outcome
- Every criterion on S-0058 is checked, or unchecked with the reason in the story notes

## Notes
