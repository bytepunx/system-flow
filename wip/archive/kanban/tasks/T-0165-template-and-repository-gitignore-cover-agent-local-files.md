---
id: T-0165
type: task
nature: improvement
title: Template and repository .gitignore cover agent-local files
status: done
parent: S-0051
owner: alex
created: 2026-09-18T21:00:24Z
updated: 2026-09-18T21:05:50Z
transitions:
  - to: ready
    at: 2026-09-18T21:05:15Z
    by: alex
  - to: in-progress
    at: 2026-09-18T21:05:15Z
    by: alex
  - to: done
    at: 2026-09-18T21:05:50Z
    by: alex
stream: S-0051
tags: []
touches: [template]
---

# T-0165 Template and repository .gitignore cover agent-local files

## Work
Add `.claude/settings.local.json` to the template's `.gitignore` source under `template/root` and to this repository's `.gitignore`, with a comment saying it is per-user agent configuration. Check the conventions for any other per-user file agents are told to create and cover it too. The smoke tier renders the template and checks it; confirm the rendered project ignores the file.

## Done when
- `git check-ignore` reports the file as ignored by the repository's own `.gitignore` here and in a project rendered from the template
- `make smoke` passes

## Notes
