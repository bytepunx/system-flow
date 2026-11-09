---
id: T-2753
type: task
nature: improvement
title: Say in the design and the user guide that the dashboard shows raw HTML in markdown as text
status: in-progress
parent: S-0627
owner: alex
created: 2026-11-09T17:11:18Z
updated: 2026-11-09T17:41:35Z
transitions:
  - to: ready
    at: 2026-11-09T17:41:35Z
    by: agent-S-0627
  - to: in-progress
    at: 2026-11-09T17:41:35Z
    by: agent-S-0627
stream: S-0627
tags: [dashboard]
touches: [design/tech/markdown.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
---
# T-2753 Say in the design and the user guide that the dashboard shows raw HTML in markdown as text

## Work

Write down the rule T-2752 builds, in the places a reader looks for it:

- `design/tech/markdown.md`: markdown-it runs with `html: false`, so raw HTML in any markdown the dashboard renders is shown as text. Say why: documents, threads, and messages are written by agents, imported repositories, and forge reviews, and `design/conventions/safety.md` treats content as data. Say that mermaid keeps `securityLevel: 'strict'`. Bump `updated`.
- `design/system/flaiover-dashboard.md`: one sentence at the markdown renderer in `## Internal structure`, linking `design/tech/markdown.md` rather than repeating it. Bump `updated`.
- `docs/users/flaiover.md`: in the Docs section's sentence on what a document renders with, add that raw HTML is shown as text, not run.

It changes none of T-2752's paths and describes the rule T-2752 decides, so it waits for none and runs beside it.

## Done when

- The three documents say that raw HTML renders as text, and the tech file says why.
- `flai test design/tech/markdown.md design/system/flaiover-dashboard.md docs/users/flaiover.md` passes the markdown lint.

## Notes
