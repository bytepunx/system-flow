---
id: T-3156
type: task
nature: improvement
title: flai stream sync's help and report, and every other git rebase --continue advice after a sync's stop, say to git add and run flai stream sync again
status: done
parent: S-0726
owner: alex
created: 2026-11-12T15:56:45Z
updated: 2026-11-12T19:26:18Z
transitions:
  - to: ready
    at: 2026-11-12T19:03:01Z
    by: agent-S-0726
  - to: in-progress
    at: 2026-11-12T19:03:01Z
    by: agent-S-0726
  - to: done
    at: 2026-11-12T19:26:18Z
    by: agent-S-0726
stream: S-0726
tags: [flai]
touches: [flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/branch.go, flai/cmd/stream_sync_test.go, flai/cmd/task_done.go, flai/cmd/storycommit.go, flai/cmd/stream_restore.go, flai/internal/storygit/commit.go, flai/internal/taskdone/taskdone.go, flai/internal/verify/story.go, flai/internal/harness/harness.go, flai/internal/mcpserver/task.go, scripts/close-out.sh, template/root/scripts/close-out.sh, docs/users/flai-reference.md, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/cmd/branch_test.go, flai/cmd/task_done_test.go, flai/internal/harness/harness_test.go, flai/internal/mcpserver/task_test.go, flai/internal/taskdone/taskdone_test.go, flai/internal/verify/story_test.go]
after: [T-3155]
usage:
  source: log
  seconds: 1397
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 124
      output: 43754
      cache_read: 8482068
      cache_write: 408186
      cost: 4.9579
---
# T-3156 flai stream sync's help and report, and every other git rebase --continue advice after a sync's stop, say to git add and run flai stream sync again

## Work

Bring the command layer and every other agent-facing advice in line with T-3155's behaviour:

- `flai/cmd/stream.go`: the `flai stream sync` long help (around lines 186–206). Drop "A rebase already in progress is never continued for you". Say that a sync run again continues a rebase whose stop was resolved and staged, and refuses one that was not.
- `flai/cmd/branch.go`: the fallback in `syncStopped.Error()` and the "To continue:" lines of `Report()`. `flai/cmd/stream_sync.go`, as the report needs.
- `flai/cmd/task_done.go`, `flai/internal/taskdone/taskdone.go` (`rebaseStopped` and the commit refusal), and `flai/internal/mcpserver/task.go`: `task_done`'s help, refusal, and MCP description. Say: resolve, `git add`, then call `task_done` again, which syncs.
  - Check whether `taskdone`'s refusal of a rebase in progress has to let the sync step run, so that `task_done` called again continues the rebase. Make it do so when the ADR says it should.
- `flai/internal/storygit/commit.go`, `flai/cmd/storycommit.go`, `flai/cmd/stream_restore.go`, and `flai/internal/verify/story.go`: "finish it with `flai stream sync`" in place of "finish it with git rebase --continue".
- `flai/internal/harness/harness.go`, `closeTasks`: the agent instruction.
- `scripts/close-out.sh` and `template/root/scripts/close-out.sh`: the stop hint.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Update the tests that assert the old texts.
- In `flai/cmd/stream_sync_test.go`, next to `testSyncConflictListsPathsAndHowToContinueOrAbort`, add a test of the command:
  - a stop
  - the path resolved and staged
  - `flai stream sync` again, which ends clean and reports the settled issue and summary.

Waits for T-3155: the help and messages describe its behaviour, and the command test needs it.

## Done when

- No agent-facing message, help, or script tells an agent to run `git rebase --continue` after a sync's stop. `git grep "rebase --continue"` under `flai/` and `scripts/` finds only the code that runs it, and tests.
- The command test passes, and `flai test` passes on the changed paths, `settings-index` and `TestReferenceIsCurrent` included.

## Notes

Added at the story's review of the plan (ADR-0294 §5, §6):

- `flai/internal/taskdone/taskdone.go`'s commit step passes a resolved stop to the sync step instead of refusing it, committing nothing while the rebase waits; an unresolved stop is still refused, saying to resolve, `git add`, and call `task_done` again. Use one exported helper from `storygit` (from T-3155) for "resolved", not a copy.
- `flai/cmd/accept.go`'s sent-back advice (around lines 817 and 824) says `git rebase --continue`: say `git add` it and run `flai stream sync` instead.
