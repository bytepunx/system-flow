---
id: T-0296
type: task
nature: feature
title: An agent is told when someone else edits its story, and what changed
status: done
parent: S-0085
owner: alex
created: 2026-09-20T15:21:42Z
updated: 2026-09-20T15:30:05Z
transitions:
  - to: ready
    at: 2026-09-20T15:30:05Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:30:05Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:30:05Z
    by: system-flow
stream: S-0085
tags: []
---
# T-0296 An agent is told when someone else edits its story, and what changed

## Work
flai edit notes who changed what (title, nature, tags, touches, parent, goal, criteria, notes, body) in a notification log under .flai-cache, outside git. inbox and wait_for_events report entries by others as edited events. No front matter key is added: item parsing is strict, so an older flai serving the same repository would refuse every item that carried one.

## Done when
- MCP tests: an edit by the designer reaches the agent once, the agent's own does not, the log stays bounded
- The Go tests and lint pass

## Notes
