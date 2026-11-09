---
id: T-2620
type: task
nature: improvement
title: "check runs an item write's rule set: the written items and their parent, the graph rules, the board, the story's narrative, and the written files' lint and conflict markers"
status: done
parent: S-0598
owner: alex
created: 2026-11-07T23:47:28Z
updated: 2026-11-08T00:07:24Z
transitions:
  - to: ready
    at: 2026-11-07T23:51:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T23:51:57Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:07:24Z
    by: agent-S-0598
stream: S-0598
tags: [cli]
touches: [flai/internal/check/itemwrite.go, flai/internal/check/itemwrite_test.go, flai/internal/check/check.go]
usage:
  source: log
  seconds: 927
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 8665
      cache_read: 5967221
      cache_write: 108179
      cost: 2.0841
---
# T-2620 check runs an item write's rule set: the written items and their parent, the graph rules, the board, the story's narrative, and the written files' lint and conflict markers

## Work

Add to `flai/internal/check` a run for a work item write, beside `Run`, in a new `itemwrite.go`. It takes the written items, their parent, and the written paths, and it runs only the rules such a write can affect:

- `oneItem` and `history` on the written items and their parent, and on the items whose rules read them, such as the tasks of an edited story;
- the graph rules over every item: `overlap`, `componentTag`, `after`, `touchesFolder`, `afterHub`, `afterReason`, `unaccepted`, and `epicLags`;
- `board`;
- `narratives` on the narrative of the written item's story alone;
- `markdown` and `conflictMarkers` on the written files alone.

It skips the rules that read no work item: layout, instructions, activities, documentation, the ADR index, conventions, roles, formulas, issues, experiments, analysis, the digest, the design tokens, the component examples (the `node` run), threads, and messages. Where `markdown`, `conflictMarkers`, or `narratives` walk the whole tree in `check.go`, give them the set of paths to read, and have `Run` pass every path, so `flai check` finds what it found before.

This task waits for none: `itemnew` and `itemedit` call what it adds.

## Done when

- `check` exports the item-write run, and `Run` still runs every rule.
- A test in `itemwrite_test.go` writes an item in fixture repositories that hold a finding of each kept rule. After `ScopeToItem`, the item-write run reports the same findings inside the story as `Run` does, for a story, a task, and an epic.
- A test shows that the item-write run starts no `node` and lints no file but the written ones.
- `flai test` passes on the paths this task changed.

## Notes
