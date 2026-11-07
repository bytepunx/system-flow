---
id: T-2488
type: task
nature: feature
title: The dashboard design, the dashboard guide, and the kluster design describe the Ops group and the Klusters page
status: done
parent: S-0580
owner: alex
created: 2026-11-06T21:45:43Z
updated: 2026-11-07T01:58:32Z
transitions:
  - to: ready
    at: 2026-11-07T01:56:03Z
    by: agent-S-0580
  - to: in-progress
    at: 2026-11-07T01:56:05Z
    by: agent-S-0580
  - to: done
    at: 2026-11-07T01:58:32Z
    by: agent-S-0580
stream: S-0580
tags: [kluster, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md]
after: [T-2487]
usage:
  source: log
  seconds: 147
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 12052
      cache_read: 5586408
      cache_write: 67201
      cost: 1.8206
---
# T-2488 The dashboard design, the dashboard guide, and the kluster design describe the Ops group and the Klusters page

## Work

It waits for T-2487, whose page it documents as built.

- `design/system/flaiover-dashboard.md` § Views: add Ops to the rail's group table, between Design and Documentation, with Klusters (`/ops/klusters`). Add the page's row to the routes table, with its data (`kluster.status`, `kluster.up`, `kluster.down`), and say that the rail remembers Ops.
- `docs/users/flaiover.md`: an Ops section with Klusters, saying what the page shows, what Up and Down do, and what it says when the host action is off or the block is missing.
- `design/system/kluster.md`: fill in its section on the dashboard with the Klusters page and its host API methods.

Bump `updated` on each.

## Done when

- The three documents say what the story built.
- `flai check --strict` and the markdown lint pass on the three paths.

## Notes
