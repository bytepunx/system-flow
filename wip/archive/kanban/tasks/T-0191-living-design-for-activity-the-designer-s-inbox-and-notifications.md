---
id: T-0191
type: task
nature: feature
title: Living design for activity, the designer's inbox, and notifications
status: done
parent: S-0042
owner: alex
created: 2026-09-19T04:46:17Z
updated: 2026-09-19T04:46:40Z
transitions:
  - to: ready
    at: 2026-09-19T04:46:18Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:46:18Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:46:40Z
    by: system-flow
stream: S-0042
tags: []
touches: [design/system]
---

# T-0191 Living design for activity, the designer's inbox, and notifications

## Work
In `design/system/flaiover-dashboard.md` add the `/activity` and `/inbox` routes, `GET /api/activity` and `GET /api/inbox` with their shapes, the navigation badge, and notifications: desktop notifications are a per-browser choice, off by default; the webhook is `dashboard.notify_url` in `system-flow.yaml`, unset by default, to which the server posts JSON for each new inbox entry. In `project-manifest.md` add the key. Say what is sent and what never is (the token, file contents), and that both are derived from the files: nothing new is written to the repository. No ADR: the story's note already fixes that presence is derived, and the webhook is an optional outbound call described in the operator documentation.

## Done when
- The design documents describe the routes, the endpoints, the badge, and both notifications
- `flai check --strict` is clean

## Notes
