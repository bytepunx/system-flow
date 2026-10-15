---
id: T-1542
type: task
nature: feature
title: The dashboard's design and user docs describe the Design group, its page, and the live preview
status: done
parent: S-0382
owner: alex
created: 2026-10-13T17:47:19Z
updated: 2026-10-15T01:54:28Z
transitions:
  - to: ready
    at: 2026-10-15T01:53:42Z
    by: system-flow
  - to: in-progress
    at: 2026-10-15T01:53:42Z
    by: agent-S-0382
  - to: done
    at: 2026-10-15T01:54:28Z
    by: agent-S-0382
stream: S-0382
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/design-system-tool.md]
after: [T-1541]
usage:
  source: log
  seconds: 46
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 9000
      cache_read: 1364431
      cache_write: 39435
      cost: 0.6836
---
# T-1542 The dashboard's design and user docs describe the Design group, its page, and the live preview

## Work

Add the Design group to § Views of `design/system/flaiover-dashboard.md`, with `/design` in the route table and how the preview is isolated from the dashboard's theme, and describe the group and its page under the header section of `docs/users/flaiover.md`. It waits for T-1541, whose page it describes.

## Done when

- Both documents describe the group, the page, and the preview, and the markdown lint passes on them.

## Notes

- Layer 3.
