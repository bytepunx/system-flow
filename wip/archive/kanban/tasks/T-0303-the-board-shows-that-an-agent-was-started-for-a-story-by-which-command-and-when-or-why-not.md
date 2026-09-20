---
id: T-0303
type: task
nature: feature
title: The board shows that an agent was started for a story, by which command, and when, or why not
status: done
parent: S-0079
owner: alex
created: 2026-09-20T15:47:56Z
updated: 2026-09-20T15:57:04Z
transitions:
  - to: ready
    at: 2026-09-20T15:53:48Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:53:48Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:57:04Z
    by: system-flow
stream: S-0079
tags: []
---
# T-0303 The board shows that an agent was started for a story, by which command, and when, or why not

## Work
A read method of the channel reports the agent action for the project: enabled, the command's name, what runs, the last start or failure. The board shows it. Stopping an agent from the dashboard is out of scope and the page says so.

## Done when
- hostapi and component tests; the flaiover tests and build pass

## Notes
