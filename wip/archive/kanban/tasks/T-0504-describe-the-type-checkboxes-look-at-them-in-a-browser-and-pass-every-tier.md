---
id: T-0504
type: task
nature: improvement
title: Describe the type checkboxes, look at them in a browser, and pass every tier
status: done
parent: S-0141
owner: alex
created: 2026-09-28T23:00:48Z
updated: 2026-09-28T23:07:28Z
transitions:
  - to: ready
    at: 2026-09-28T23:00:55Z
    by: agent-S-0141
  - to: in-progress
    at: 2026-09-28T23:03:46Z
    by: agent-S-0141
  - to: done
    at: 2026-09-28T23:07:28Z
    by: agent-S-0141
stream: S-0141
tags: [dashboard]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
---
# T-0504 Describe the type checkboxes, look at them in a browser, and pass every tier

## Work

Update the Board section of `docs/users/flaiover.md` and the `/board` row of `design/system/flaiover-dashboard.md`. Run the branch's dev server against a scratch project (never the operator's dashboard) and check in a browser that each checkbox toggles its type and that the choice survives a reload and a navigation away and back. Run `make test`, `make integration`, `make smoke`, `make flaiover-test`, and `flai check --strict`. Tick the story's criteria for what was seen.

## Done when

- Docs describe the checkboxes and that the choice is remembered in the browser.
- Every tier passes and the criteria are ticked for what was seen.

## Notes
