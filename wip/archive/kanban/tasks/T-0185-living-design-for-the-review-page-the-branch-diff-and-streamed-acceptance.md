---
id: T-0185
type: task
nature: feature
title: Living design for the review page, the branch diff, and streamed acceptance
status: done
parent: S-0041
owner: alex
created: 2026-09-19T04:08:28Z
updated: 2026-09-19T04:08:52Z
transitions:
  - to: ready
    at: 2026-09-19T04:08:30Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:08:30Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:08:52Z
    by: system-flow
stream: S-0041
tags: []
touches: [design/system]
---

# T-0185 Living design for the review page, the branch diff, and streamed acceptance

## Work
No new ADR: this builds on ADR-0016 (the dashboard delegates to flai), ADR-0019 (story branches), and S-0046 (done is acceptance). In `design/system/flaiover-dashboard.md` add the `/review/<id>` route and the endpoints: `GET /api/items/:id/diff` (`flai stream diff`), and `POST /api/items/:id/accept`, which runs `flai accept <id> --by <designer>` and streams newline-delimited JSON: one `progress` line per step flai logs, then one `done` line with the result or one `error` line with flai's message verbatim. Say who the designer is (the manifest's owner, as for threads and moves, since one project token has one holder). In `design/system/flai-cli.md` add `flai stream diff` and that `flai accept` logs each step as an info event.

## Done when
- Both documents describe the route, the endpoints, and the command
- `flai check --strict` is clean

## Notes
