---
id: T-0196
type: task
nature: feature
title: "Verify: all tiers, the pages in a real container, and the webhook against a local listener"
status: done
parent: S-0042
owner: alex
created: 2026-09-19T04:46:18Z
updated: 2026-09-19T04:59:20Z
transitions:
  - to: ready
    at: 2026-09-19T04:53:39Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:53:39Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:59:20Z
    by: system-flow
stream: S-0042
tags: []
---

# T-0196 Verify: all tiers, the pages in a real container, and the webhook against a local listener

## Work
Run `make flai-test`, `make flaiover-test`, and a production build. In a container built from the branch against a scratch project: `/api/activity` and `/api/inbox` reflect a narrative, a thread awaiting the designer, a story in review, a blocked item, and an overlap; the pages render at desktop and phone widths and the badge shows the count. For the webhook, run the dev server on the host with `dashboard.notify_url` pointing at a local listener, make a change that creates an inbox entry, and see exactly one POST with the documented body. Do not touch the operator's dashboard for this. Tick the story's criteria for what was observed.

## Done when
- All tiers pass, with results in the narrative
- The container and webhook runs are recorded with their outcomes
- Every criterion on S-0042 is checked, or unchecked with the reason in the story notes

## Notes
