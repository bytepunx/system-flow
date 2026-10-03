---
id: T-0756
type: task
nature: feature
title: Finalizing a draft records who finalized it, and the hostapi op item.finalize runs it
status: done
parent: S-0201
owner: alex
created: 2026-10-03T07:36:59Z
updated: 2026-10-03T07:57:57Z
transitions:
  - to: ready
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: in-progress
    at: 2026-10-03T07:44:52Z
    by: agent-S-0201
  - to: done
    at: 2026-10-03T07:57:57Z
    by: agent-S-0201
stream: S-0201
tags: []
touches: [flai/internal/workitem, flai/cmd/edit.go, flai/cmd/items.go, flai/cmd/planning.go, flai/internal/hostapi, flai/internal/itemedit, flaiover/src/lib/server/agent.ts]
after: [T-0753]
usage:
  source: log
  seconds: 785
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 199
      output: 55815
      cache_read: 10313295
      cache_write: 259230
      cost: 4.7003
---
# T-0756 Finalizing a draft records who finalized it, and the hostapi op item.finalize runs it

## Work

Clearing a story's draft flag records who did it and when, in the form TH-0083 settles (recommended: a `finalized` block, `{by, at}`, stamped by `flai edit --no-draft` and `flai move <story> ready --yes` with `--by`, else `FLAI_AGENT`, else the config author). The hostapi write `item.finalize` takes `{id}` and runs `flai edit <id> --no-draft --by=<owner>`, refusing an item that is not a draft story. It waits for T-0753, which touches `flai/internal/workitem` too, and for the answer on TH-0083.

## Done when

- [ ] Finalizing with `flai edit --no-draft`, `flai move ready --yes`, or `item.finalize` records the finalizer and the time as TH-0083 settles, and `flai check` validates it
- [ ] `item.finalize` is registered in hostapi and refuses a story that is not a draft
- [ ] Go tests cover the record, the op's arguments, and the refusal

## Notes
