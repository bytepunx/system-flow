---
id: T-2786
type: task
nature: remediation
title: flai serve starts no agent for a ready story while its hold stands, with a test that reproduces I-0301
status: done
parent: S-0633
owner: alex
created: 2026-11-09T22:01:24Z
updated: 2026-11-09T23:13:57Z
transitions:
  - to: ready
    at: 2026-11-09T22:21:08Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T22:21:09Z
    by: system-flow
  - to: done
    at: 2026-11-09T23:13:57Z
    by: agent-S-0633
stream: S-0633
tags: [flai, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/hold_test.go, flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/internal/workitem/itemlock.go]
usage:
  source: log
  seconds: 3168
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 75
      output: 617
      cache_read: 10941564
      cache_write: 1759480
      cost: 4.6441
---
# T-2786 flai serve starts no agent for a ready story while its hold stands, with a test that reproduces I-0301

## Work

The cause is found (see Notes): the look read a task that a checked creation had written and was about to remove, and judged the story's claim by that task's touches.

1. In `flai/internal/workitem/itemlock.go`, add `(*Repo).LockCheckedWrites() (unlock func(), err error)`: an operating system lock on `.flai-cache/items-checked.lock`, through `lockFile` with `itemsLockWait`, apart from the items lock that `Save` takes, so that a holder may still save.
2. In `flai/internal/itemnew/itemnew.go`, `Create` takes it before `repo.Create` writes the new item and holds it until the check has kept the item or `undo` has removed it, releasing it before `opt.Then` runs. A lock that cannot be taken refuses the creation, writing nothing.
3. In `flai/internal/serve/agents.go`, the look reads the items for its holds (`readyStories`, both calls) under the same lock, so it never judges a claim by an item a checked creation may yet remove. A lock it cannot take ends the look with a warning, as a failed read does.
4. Tests: in `flai/internal/serve/hold_test.go`, a ready story with no touches, a story in progress, and a task of the ready story with touches written while another holder holds the lock, removed before it lets go: the look starts no agent. It fails on the code before the change. In `flai/internal/itemnew/itemnew_test.go`, the lock is held while the check runs and free once `Create` returns, kept or refused.

The operator's own start (`flai/internal/serve/start.go`) is left as it is.

Waits for no task: it shares no path with T-2787.

## Done when

- The new test in `flai/internal/serve/hold_test.go` fails on the code before the change and passes after it
- A checked creation holds the lock from its write until it keeps or removes the item, as a test in `flai/internal/itemnew/itemnew_test.go` shows
- `TestAHeldStoryIsSkippedAndStartedFirstOnceClear` and `TestOneLookStartsNoTwoThatOverlap` still pass
- `flai test flai/internal/serve flai/internal/itemnew flai/internal/workitem` passes

## Notes

The path that started S-0568: the look itself, not another start path. At 2026-11-07T08:18:43Z planner-S-0568 called `item_new` for its first task, with touches under `flai/internal/guard/`. `itemnew.Create` wrote it as T-2550 in `wip/kanban/tasks/`, ran the checked creation's check, found `item.heading` (no `## Notes`), and removed it, answering "nothing was created" at 08:18:48Z. The look at 08:18:46Z read T-2550 between the write and the removal, so S-0568's claim was T-2550's touches, which overlap no story in progress: no hold, and `l.start` launched agent-S-0568. Its `story_start` at 08:18:50Z read the items again, without T-2550, and found it held `no-touches`. The planner wrote it again as T-2551 at 08:18:56Z. Sources: `~/.flai/serve/agents/sf-planner-20261107T081614Z.log`, `~/.flai/serve/agents/sf-S-0568-20261107T081846Z.log`, `~/.flai/serve/journal.jsonl`.
