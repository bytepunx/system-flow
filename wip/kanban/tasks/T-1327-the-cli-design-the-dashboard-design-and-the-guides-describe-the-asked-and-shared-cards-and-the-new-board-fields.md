---
id: T-1327
type: task
nature: feature
title: The CLI design, the dashboard design, and the guides describe the asked and shared cards and the new board fields
status: done
parent: S-0338
owner: alex
created: 2026-10-08T04:32:13Z
updated: 2026-10-08T10:41:36Z
transitions:
  - to: ready
    at: 2026-10-08T10:39:28Z
    by: agent-S-0338
  - to: in-progress
    at: 2026-10-08T10:39:28Z
    by: agent-S-0338
  - to: done
    at: 2026-10-08T10:41:36Z
    by: agent-S-0338
stream: S-0338
tags: [flai, flaiover]
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md]
after: [T-1325, T-1326]
usage:
  source: log
  seconds: 127
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 232
      cache_read: 2205670
      cache_write: 99716
      cost: 1.0325
---
# T-1327 The CLI design, the dashboard design, and the guides describe the asked and shared cards and the new board fields

## Work

Document what T-1324 to T-1326 built. It waits for T-1325 and T-1326, the last of them.

- `design/system/flai-cli.md` § `flai board`: the card's new fields beside `held: {code, reason}`.
- `docs/users/flai.md`: the same fields in `flai board --json`.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`: what an asked card and a shared card say, and the share in a conversation.

## Done when

- Each document names the fields or the card it describes.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
