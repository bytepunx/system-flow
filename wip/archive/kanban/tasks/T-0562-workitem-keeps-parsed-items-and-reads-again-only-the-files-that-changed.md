---
id: T-0562
type: task
nature: improvement
title: workitem keeps parsed items and reads again only the files that changed
status: done
parent: S-0156
owner: alex
created: 2026-09-29T19:20:02Z
updated: 2026-09-29T19:23:47Z
transitions:
  - to: ready
    at: 2026-09-29T19:20:08Z
    by: agent-S-0156
  - to: in-progress
    at: 2026-09-29T19:20:09Z
    by: agent-S-0156
  - to: done
    at: 2026-09-29T19:23:47Z
    by: agent-S-0156
stream: S-0156
tags: []
touches: [flai/internal/workitem]
usage:
  source: log
  seconds: 218
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 9135
      cache_read: 1814597
      cache_write: 30636
      cost: 0.7908
---
# T-0562 workitem keeps parsed items and reads again only the files that changed

## Work

Keep the items `Repo.List` and `Repo.Get` parse in a process-wide store keyed by file path, with each file's modification time and size. Each list walks the item folders, stats every file, and parses again only a file that is new or whose time or size changed; a file gone from the walk leaves the store. A file modified too recently for its time to tell two writes apart is read again until it is older (the racy-git rule). Callers get their own copy of each item, so that one that changes an item changes nothing another sees. One folder is refreshed by one caller at a time, so that two lists at once share one read.

## Done when

- Behaviour tests show an item created, edited, moved, archived, or deleted, by flai or by a hand edit (same size, same second included), is seen by the next list and get.
- A test shows two lists at once parse each file once.
- A test shows a caller changing a listed item does not change the next list.
- `make test` and lint pass.

## Notes
