---
id: T-0703
type: task
nature: feature
title: A story's page offers New story on its own and New sibling under its open epic
status: done
parent: S-0192
owner: arobson
created: 2026-10-02T16:31:05Z
updated: 2026-10-02T16:33:38Z
transitions:
  - to: ready
    at: 2026-10-02T16:31:08Z
    by: agent-S-0192
  - to: in-progress
    at: 2026-10-02T16:31:08Z
    by: agent-S-0192
  - to: done
    at: 2026-10-02T16:33:38Z
    by: agent-S-0192
stream: S-0192
tags: []
touches: ["flaiover/src/routes/items/[id]", design/system/flaiover-dashboard.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 150
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 9250
      cache_read: 1150828
      cache_write: 53830
      cost: 0.7837
---
# T-0703 A story's page offers New story on its own and New sibling under its open epic

## Work

On a story's page in flaiover (`flaiover/src/routes/items/[id]/+page.svelte`), **New story** opens `/new?type=story` with no epic, and after it **New sibling** opens `/new?type=story&parent=<the story's epic>`. New sibling shows only on a writable dashboard, when the story has an epic that is open: not done, cancelled, or archived, since the form offers no other as a parent (as Create story on an epic, S-0191). The designer chose this on TH-0068. The page's test covers both links, and the design and the user guide describe them.

It waits for no other task: it is the story's only one.

## Done when

- [x] New story on a story's page opens `/new?type=story`, with no parent.
- [x] New sibling follows New story and opens `/new?type=story&parent=<epic>` when the story's epic is open, and is not shown when the story has no epic, its epic is done, cancelled, or archived, or the dashboard cannot write.
- [x] `item.svelte.test.ts` covers these and passes.
- [x] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say what the two links do.

## Notes
