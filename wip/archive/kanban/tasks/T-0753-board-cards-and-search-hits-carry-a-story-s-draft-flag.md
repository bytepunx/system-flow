---
id: T-0753
type: task
nature: feature
title: Board cards and search hits carry a story's draft flag
status: done
parent: S-0201
owner: alex
created: 2026-10-03T07:36:43Z
updated: 2026-10-03T07:44:35Z
transitions:
  - to: ready
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: in-progress
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: done
    at: 2026-10-03T07:44:35Z
    by: agent-S-0201
stream: S-0201
tags: []
touches: [flai/internal/workitem/boardview.go, flai/internal/search, flai/internal/hostapi/searching.go]
usage:
  source: log
  seconds: 423
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 6497
      cache_read: 1431608
      cache_write: 37166
      cost: 0.6549
---
# T-0753 Board cards and search hits carry a story's draft flag

## Work

flai's `board.get` card (`BoardCard` in `flai/internal/workitem/boardview.go`) and its search hit (`flai/internal/search`) gain `draft`, true on a draft story and left out otherwise, so the dashboard can mark drafts without a second request per item. It waits for nothing.

## Done when

- [ ] A draft story's board card and search hit carry `"draft": true`, and a finalized story's carry no `draft`
- [ ] Tests in `internal/workitem` and `internal/search` cover both

## Notes
