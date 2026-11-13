---
id: T-3320
type: task
nature: improvement
title: flai guard refuses a Write that replaces a work item file with content that is not that item
status: done
parent: S-0730
owner: alex
created: 2026-11-13T19:54:06Z
updated: 2026-11-13T20:18:05Z
transitions:
  - to: ready
    at: 2026-11-13T20:07:32Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:07:32Z
    by: system-flow
  - to: done
    at: 2026-11-13T20:18:05Z
    by: agent-S-0730
stream: S-0730
tags: [flai]
touches: [flai/internal/guard/items.go, flai/internal/guard/items_test.go, flai/internal/guard/call.go, design/adrs/0298-flai-guard-refuses-a-write-that-replaces-a-work-item-file-whole-with-content.md, design/adrs/README.md]
after: [T-3319]
usage:
  source: log
  seconds: 633
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 3367
      cache_read: 2825292
      cache_write: 27053
      cost: 0.8349
---
# T-3320 flai guard refuses a Write that replaces a work item file with content that is not that item

## Work

Write the rule I-0352 recommends, in a new `flai/internal/guard/items.go`, and call it from `Decide` in `flai/internal/guard/call.go`. It applies to every session the Write hook runs for, which is any session with `FLAI_STORY` or `FLAI_ROLE` set, and to their sub-agents.

- It covers a FileWrite with `Content` whose `Path` is a work item file: an epic, story, or task file under the board's `epics/`, `stories/`, or `tasks/` folder. That holds in the main checkout and in a story's worktree. The I-0352 instance wrote to the main checkout's T-3079 from a story's session.
- It refuses the write when the content has no front matter (`workitem.SplitFrontMatter` gives `ErrNoFrontMatter`), when the front matter does not parse as an item (`workitem.ParseItem`), or when its `id` or `type` differs from the file on disk (`workitem.ReadItem`). A file that does not exist yet is refused too: items are made with `flai`, not written.
- The refusal names the file and points the agent to `Edit` for the body and to `flai edit` or `item_edit` for the rest.
- A Write that keeps the item's own front matter passes, so a whole rewrite of the body stays allowed. Any file outside those folders passes.
- It waits on T-3319 because it reads `Call.Content`, which that task adds. It also touches `call.go`, which that task changes.

## Done when

- `items_test.go` reproduces I-0352: a story agent's `Write` of `PLACEHOLDER` to an existing task file is refused, with the refusal naming `Edit` and `flai edit`.
- The tests also show that a Write which keeps the file's own id and type passes, that one changing the id is refused, and that a Write to a file outside the item folders passes.
- `flai test` on the changed paths passes.

## Notes
