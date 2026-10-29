---
id: T-2339
type: task
nature: improvement
title: A thread reason loads its item as its title, status, goal, and criteria, a reason pack's digest is capped at a sixteenth of the budget, and the quarter test measures a real thread
status: done
parent: S-0537
owner: alex
created: 2026-10-29T16:46:05Z
updated: 2026-10-29T17:12:55Z
transitions:
  - to: ready
    at: 2026-10-29T16:46:30Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T16:46:31Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T17:12:55Z
    by: agent-S-0537
stream: S-0537
tags: []
touches: [flai/internal/context/reasons.go, flai/internal/context/reasons_test.go, flai/internal/context/digest.go, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, flai/cmd/prime.go, docs/users/flai-reference.md, flai/internal/context/testdata/thread/S-0542-claude-code-s-edit-tool-cannot-write-a-work-item.md, flai/internal/context/testdata/thread/TH-0750-plan-for-s-0542-cut-every-title-made-file-name-s-slug-at-80-bytes-in-four-tasks-over-three-layers.md]
usage:
  source: log
  seconds: 1584
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 153
      output: 56811
      cache_read: 11352507
      cache_write: 366823
      cost: 5.4921
---
# T-2339 A thread reason loads its item as its title, status, goal, and criteria, a reason pack's digest is capped at a sixteenth of the budget, and the quarter test measures a real thread

## Work

TH-0753: the operator confirmed the orchestrator's recommendation. Criterion 2's test passed on a one-entry fixture thread, while real thread packs here measured 31 to 33% of the general pack.

- In `reasons.go`, a `thread` reason loads the work item it is on as its title, status, goal, and acceptance criteria. The notes, tasks, and planning are left out; `item_get` gives them. Other reasons keep loading their items whole.
- A reason pack's digest is capped at a sixteenth of the budget. That is still within ADR-0194's eighth, and every other pack keeps its eighth.
- The quarter test measures a pack for a copy of a real thread of this repository, with its story, against the general pack.
- `agent-context.md`'s measured table, and the docs that say a thread loads its item whole, follow.
- The story's agent writes the ADR refining ADR-0217 §3 and §4.

No `after`: every other task is done.

## Done when

- The quarter test passes over a real thread and its story, and logs both sizes.
- A test asserts a thread's item loads its goal and criteria but not its notes.
- A test asserts the reason pack's digest stays within a sixteenth of the budget.
- `flai test` passes on the changed paths.

## Notes
