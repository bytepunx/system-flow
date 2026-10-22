---
id: T-2029
type: task
nature: feature
title: writes.test.ts's temp projects ignore .flai-cache/ as a project made from the template does, so flai's items lock does not trip layout.gitignore
status: done
parent: S-0484
owner: alex
created: 2026-10-22T20:39:24Z
updated: 2026-10-22T20:40:12Z
transitions:
  - to: ready
    at: 2026-10-22T20:40:05Z
    by: agent-S-0484
  - to: in-progress
    at: 2026-10-22T20:40:06Z
    by: agent-S-0484
  - to: done
    at: 2026-10-22T20:40:12Z
    by: agent-S-0484
stream: S-0484
tags: []
touches: [flaiover/src/lib/server/writes.test.ts]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1348
      cache_read: 210963
      cache_write: 7563
      cost: 0.1152
---

# T-2029 writes.test.ts's temp projects ignore .flai-cache/ as a project made from the template does, so flai's items lock does not trip layout.gitignore

## Work

Since S-0479, every work item write takes `.flai-cache/items.lock` in the project. The temp copies of the fixture that `writes.test.ts` edits an item in and lints a new story in have no `.gitignore`, so flai's scoped check finds `layout.gitignore` and refuses the write: two tests fail on main as on this branch, and stop the close-out's flaiover tier. Write a `.gitignore` naming `.flai-cache/` into those two copies, as a project made from the template has. Found at S-0484's close-out; it waits for no task.

## Done when

`flai test` on `flaiover/src/lib/server/writes.test.ts` passes.

## Notes

The flai behaviour itself, a write refused for the lock file it made, is recorded as an issue for the operator.
