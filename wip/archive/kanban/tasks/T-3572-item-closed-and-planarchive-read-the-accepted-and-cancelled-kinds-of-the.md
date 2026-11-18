---
id: T-3572
type: task
nature: remediation
title: Item.Closed and PlanArchive read the accepted and cancelled kinds of the workflow the repository runs
status: done
parent: S-0842
owner: alex
created: 2026-11-17T12:29:38Z
updated: 2026-11-18T07:49:58Z
transitions:
  - to: ready
    at: 2026-11-18T07:32:41Z
    by: agent-S-0842
  - to: in-progress
    at: 2026-11-18T07:32:41Z
    by: agent-S-0842
  - to: done
    at: 2026-11-18T07:49:57Z
    by: agent-S-0842
stream: S-0842
tags: [flai, workflow]
touches: [flai/internal/workitem/item.go, flai/internal/workitem/repo.go, flai/internal/workitem/archive.go, flai/internal/workitem/archive_test.go, flai/internal/workitem/store_test.go, flai/internal/workitem/view.go]
usage:
  source: log
  seconds: 1036
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 5897
      cache_read: 2398661
      cache_write: 71959
      cost: 1.0642
---
# T-3572 Item.Closed and PlanArchive read the accepted and cancelled kinds of the workflow the repository runs

## Work

Settle the fix the story's Planning note proposes, or the narrower one it names, and record the choice in the narrative's `## Decisions`. With the proposed fix: `Repo.List` and `Repo.Get` give each item they return the repository's `Flow()`, without changing an item the parse cache shares between repositories in a way another repository would read; `Item.Closed`, and `CompletedAt` through it, answer by `Workflow.Closed`, and by the baseline's names on an item no repository gave a workflow. `PlanArchive` reads the same, and its refusal names the workflow's accepted and cancelled states in place of "done or cancelled". Test `PlanArchive` with IDs and without on a workflow whose accepted state is `shipped`, and that items `List` and `Get` return from the cache answer `Closed` by it. First layer: it waits for no task.

## Done when

- On a workflow whose accepted state is `shipped`, `PlanArchive` archives a shipped story and its tasks, and its refusal of an open item names `shipped` and the cancelled state.
- An item read through `Repo.List` or `Repo.Get` answers `Closed` by the repository's workflow; one parsed alone answers by the baseline's names.
- `flai test` on the changed paths passes.

## Notes
