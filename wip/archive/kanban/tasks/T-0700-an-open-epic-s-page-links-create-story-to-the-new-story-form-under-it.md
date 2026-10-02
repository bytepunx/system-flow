---
id: T-0700
type: task
nature: feature
title: An open epic's page links Create story to the new story form under it
status: done
parent: S-0191
owner: arobson
created: 2026-10-02T16:07:59Z
updated: 2026-10-02T16:09:18Z
transitions:
  - to: ready
    at: 2026-10-02T16:08:04Z
    by: agent-S-0191
  - to: in-progress
    at: 2026-10-02T16:08:04Z
    by: agent-S-0191
  - to: done
    at: 2026-10-02T16:09:18Z
    by: agent-S-0191
stream: S-0191
tags: [dashboard]
touches: ["flaiover/src/routes/items/[id]", design/system/flaiover-dashboard.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 74
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 4677
      cache_read: 914433
      cache_write: 17001
      cost: 0.4125
---
# T-0700 An open epic's page links Create story to the new story form under it

## Work

Beside the **New epic** link at the top of an epic's page, add a **Create story** link, after it, that opens `/new?type=story&parent=<epic>`, so the new story form starts with this epic chosen. Offer it, like **New epic**, only when the dashboard can write, and only on an epic that is open and not archived: the form offers only open epics as parents and drops a `?parent=` that names a closed one, so on a closed epic the link would not do what it says. Test it in `item.svelte.test.ts` and describe it in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`. It is the story's only task and waits for none.

## Done when

- [x] An open epic's page, when the dashboard can write, shows **Create story** after **New epic**, linking `/new?type=story&parent=<epic>`.
- [x] A done, cancelled, or archived epic, a story, a task, and a read-only dashboard show no **Create story**.
- [x] The item page tests cover both, and the design and user docs say it.

## Notes
