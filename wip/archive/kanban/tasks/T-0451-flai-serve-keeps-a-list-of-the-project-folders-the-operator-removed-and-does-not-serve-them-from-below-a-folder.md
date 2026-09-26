---
id: T-0451
type: task
nature: improvement
title: flai serve keeps a list of the project folders the operator removed and does not serve them from below a folder
status: done
parent: S-0123
owner: alex
created: 2026-09-26T07:06:21Z
updated: 2026-09-26T07:08:50Z
transitions:
  - to: ready
    at: 2026-09-26T07:06:28Z
    by: agent-S-0123
  - to: in-progress
    at: 2026-09-26T07:06:29Z
    by: agent-S-0123
  - to: done
    at: 2026-09-26T07:08:50Z
    by: agent-S-0123
stream: S-0123
tags: [cli]
touches: [flai/internal/serve, flai/cmd]
---
# T-0451 flai serve keeps a list of the project folders the operator removed and does not serve them from below a folder

## Work

- `serve.Dir` keeps `removed.json` beside `projects.json`: the roots the operator removed, read, added to, and taken off.
- `serve.Place` leaves a found project whose root is in the list out, with `removed: true` and a reason saying so and how to serve it again; flai serve reads the list at every tick, so a removal takes effect within a second.
- `flai serve status` and `flai serve project list` show it as not served, with why, whether flai serve runs or not.

## Done when

- Behaviour tests cover Place with a removed root below an import folder and below the start folder, and the list's round trip; `make test` passes.

## Notes
