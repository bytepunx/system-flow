---
id: T-2492
type: task
nature: feature
title: The dashboard design, the dashboard guide, and the kluster design describe the Deployments page
status: done
parent: S-0581
owner: alex
created: 2026-11-06T21:46:22Z
updated: 2026-11-07T03:05:56Z
transitions:
  - to: ready
    at: 2026-11-07T03:03:32Z
    by: agent-S-0581
  - to: in-progress
    at: 2026-11-07T03:03:33Z
    by: agent-S-0581
  - to: done
    at: 2026-11-07T03:05:56Z
    by: agent-S-0581
stream: S-0581
tags: [kluster, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/kluster.md]
after: [T-2491]
usage:
  source: log
  seconds: 143
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 17704
      cache_read: 4735426
      cache_write: 77430
      cost: 1.8449
---
# T-2492 The dashboard design, the dashboard guide, and the kluster design describe the Deployments page

## Work

It waits for the Deployments page task of this story, whose page it documents as built.

- `design/system/flaiover-dashboard.md` § Views: add Deployments (`/ops/deployments`) to the Ops row of the rail's group table, and the page's row to the routes table, with its data.
- `docs/users/flaiover.md`: Deployments under the Ops section, saying what each row shows, each action, and how to reach a deployment for a demo.
- `design/system/kluster.md`: add the Deployments page and its host API methods to its section on the dashboard.

Bump `updated` on each.

## Done when

- The three documents say what the story built.
- `flai check --strict` and the markdown lint pass on the three paths.

## Notes
