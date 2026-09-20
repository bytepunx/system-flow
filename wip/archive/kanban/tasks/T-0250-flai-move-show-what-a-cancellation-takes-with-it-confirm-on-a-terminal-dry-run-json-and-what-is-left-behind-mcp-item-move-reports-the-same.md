---
id: T-0250
type: task
nature: improvement
title: "flai move: show what a cancellation takes with it, confirm on a terminal, --dry-run, --json, and what is left behind; MCP item_move reports the same"
status: done
parent: S-0070
owner: alex
created: 2026-09-20T06:16:43Z
updated: 2026-09-20T06:23:29Z
transitions:
  - to: ready
    at: 2026-09-20T06:23:29Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T06:23:29Z
    by: system-flow
  - to: done
    at: 2026-09-20T06:23:29Z
    by: system-flow
stream: S-0070
tags: []
---
# T-0250 flai move: show what a cancellation takes with it, confirm on a terminal, --dry-run, --json, and what is left behind; MCP item_move reports the same

## Work
`flai move <id> cancelled` prints the items that go with it by type and state, asks on a terminal unless `--yes`, and changes nothing under `--dry-run`. `--json` lists every item cancelled. For each cancelled story with a narrative, a `story/<id>` branch, or a worktree, the output says what is left for a person to keep or remove; nothing is deleted. MCP `item_move` returns the cancelled list too. `inbox` and `wait_for_events` say which parent a cascaded cancellation came from.

## Done when
- Command tests for the preview, the refusal to proceed without confirmation, `--yes`, `--dry-run`, and `--json`
- An MCP test that an epic cancelled by the designer shows each child's cancellation with its cause

## Notes
