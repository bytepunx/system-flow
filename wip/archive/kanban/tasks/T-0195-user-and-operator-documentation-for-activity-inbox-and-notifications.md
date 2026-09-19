---
id: T-0195
type: task
nature: feature
title: User and operator documentation for activity, inbox, and notifications
status: done
parent: S-0042
owner: alex
created: 2026-09-19T04:46:17Z
updated: 2026-09-19T04:53:38Z
transitions:
  - to: ready
    at: 2026-09-19T04:53:26Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:53:26Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:53:38Z
    by: system-flow
stream: S-0042
tags: []
touches: [docs]
---

# T-0195 User and operator documentation for activity, inbox, and notifications

## Work
Update `docs/users/flaiover.md` (Activity, Inbox, the badge, desktop notifications) and `docs/operators/index.md` (`dashboard.notify_url`: what is posted, when, what is never sent, that it is off by default, failure behaviour).

## Done when
- Both documents state the behaviour as built
- Markdown lint and `flai check --strict` pass

## Notes
