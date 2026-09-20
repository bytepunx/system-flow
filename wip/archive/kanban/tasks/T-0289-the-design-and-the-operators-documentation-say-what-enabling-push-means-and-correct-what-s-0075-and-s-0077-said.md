---
id: T-0289
type: task
nature: feature
title: The design and the operators' documentation say what enabling push means, and correct what S-0075 and S-0077 said
status: done
parent: S-0078
owner: alex
created: 2026-09-20T13:06:49Z
updated: 2026-09-20T13:20:21Z
transitions:
  - to: ready
    at: 2026-09-20T13:19:07Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T13:19:08Z
    by: system-flow
  - to: done
    at: 2026-09-20T13:20:21Z
    by: system-flow
stream: S-0078
tags: []
---
# T-0289 The design and the operators' documentation say what enabling push means, and correct what S-0075 and S-0077 said

## Work
Operators: what enabling means (a holder of the dashboard token can then publish any story an agent has put in review), how to enable and disable, the journal; the sentences that said nothing pushes unasked are made true and the interval in which they were not is stated. Design: flai-cli.md, flaiover-dashboard.md, dashboard-host-channel.md. Users' pages.

## Done when
- The markdown lint and flai check --strict pass, run after the last edit

## Notes
