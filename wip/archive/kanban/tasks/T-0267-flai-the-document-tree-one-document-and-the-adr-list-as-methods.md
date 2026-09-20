---
id: T-0267
type: task
nature: feature
title: "flai: the document tree, one document, and the ADR list as methods"
status: done
parent: S-0074
owner: alex
created: 2026-09-20T08:29:59Z
updated: 2026-09-20T08:32:28Z
transitions:
  - to: ready
    at: 2026-09-20T08:30:01Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T08:30:01Z
    by: system-flow
  - to: done
    at: 2026-09-20T08:32:28Z
    by: system-flow
stream: S-0074
tags: []
---
# T-0267 flai: the document tree, one document, and the ADR list as methods

## Work
docs.tree (the manifest's three folders, each Markdown file with its title and front matter), doc.get (one Markdown file under those folders: front matter, body, raw; anything else refused), adrs.list (id, title, status, date, supersedes, superseded by, refines). Front matter reaches the dashboard as the strings the files hold.

## Done when
- Tests on a scratch project, including a path that escapes, a file that is not Markdown, and a dated front matter

## Notes
