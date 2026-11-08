---
id: T-2625
type: task
nature: improvement
title: itemedit.Apply checks an edit with the item-write rule set before and after, in place of check.Run
status: done
parent: S-0598
owner: alex
created: 2026-11-07T23:47:55Z
updated: 2026-11-08T00:33:43Z
transitions:
  - to: ready
    at: 2026-11-08T00:08:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T00:08:42Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:33:43Z
    by: agent-S-0598
stream: S-0598
tags: [cli]
touches: [flai/internal/itemedit/itemedit.go, flai/internal/itemedit/rules_test.go, flai/internal/itemedit/criteria_test.go]
after: [T-2620]
usage:
  source: log
  seconds: 1501
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 15942
      cache_read: 7735227
      cache_write: 158230
      cost: 2.8282
---
# T-2625 itemedit.Apply checks an edit with the item-write rule set before and after, in place of check.Run

## Work

In `flai/internal/itemedit/itemedit.go`, `runCheck` is `check.Run`, which `Apply` runs before and after the write. Make it the item-write rule set T-2620 adds. The rule set takes the edited item, its parent, its old and new parent on a reparent, and every file the edit wrote, including a retitle's renamed file, its narrative, and the files whose links it rewrote. Keep `runCheck` a variable, so that I-0182's test can still bring another writer's finding between the two runs. Keep the scoping, the error on the item's own file, the advisory rules, and the restore from the snapshot as they are (ADR-0197).

An edit also lists the items again and rewrites `wip/agents/index.md`. Time `flai edit` once the rule set is in. If that second list keeps it over 3 s, reuse the list the check made.

This task waits for T-2620, whose rule set it calls. It shares no file with the creation's task, so the two run together. `flai criteria tick`, finalize, and serve's replan all save through `Apply`, so they get the change with no edit of their own.

## Done when

- `Apply` no longer calls `check.Run`.
- Tests in `rules_test.go` make an edit, a retitle, and a reparent that each bring a finding of each kept rule. Each is refused as before, with every file put back, and the advisory rules still come back as warnings.
- A test shows that a finding of a skipped rule already in the repository refuses nothing.
- `flai test` passes on the paths this task changed.

## Notes
