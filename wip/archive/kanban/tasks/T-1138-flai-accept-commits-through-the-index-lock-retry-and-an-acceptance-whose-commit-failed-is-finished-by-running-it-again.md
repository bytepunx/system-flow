---
id: T-1138
type: task
nature: remediation
title: flai accept commits through the index-lock retry, and an acceptance whose commit failed is finished by running it again
status: done
parent: S-0307
owner: alex
created: 2026-10-07T01:12:45Z
updated: 2026-10-07T01:56:03Z
transitions:
  - to: ready
    at: 2026-10-07T01:44:49Z
    by: agent-S-0307
  - to: in-progress
    at: 2026-10-07T01:44:50Z
    by: agent-S-0307
  - to: done
    at: 2026-10-07T01:56:03Z
    by: agent-S-0307
stream: S-0307
tags: [flai]
touches: [flai/internal/preview/accept.go, flai/cmd/accept.go, flai/cmd/accept_lock_test.go]
after: [T-1137]
usage:
  source: log
  seconds: 673
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 89
      output: 37002
      cache_read: 4274149
      cache_write: 192165
      cost: 2.8112
---
# T-1138 flai accept commits through the index-lock retry, and an acceptance whose commit failed is finished by running it again

## Work

Three changes to the acceptance in `flai/cmd/accept.go` and `flai/internal/preview/accept.go`. It waits for T-1137, whose helper the commit step calls.

1. Step 3 runs `git add -A` and `git commit` through T-1137's helper, so a lock another process holds for a moment no longer stops the acceptance.
2. An item that is done and archived but whose acceptance commit is missing is resumable, as one that is done but never archived already is. Today `preview.AcceptWith` refuses it as already done (`it.Closed()`). Tell the two apart by what git holds: the item is done and archived, and its archived files, or anything under the wip folder, have changes the main checkout has not committed. Such an item resumes at the commit step: no merge, no transition, no archive, then the commit with the usual subject, `chore: [S-nnnn] accept and archive` (with the epic when one followed), then the overlap notices. The paths for the notices cannot come from the merge any more; take them from the story's commits (`storygit.Committed`), and say so in a comment. A done, archived item with nothing to commit stays refused as already done. The orchestrator's acceptance still refuses a resumed item, as the operator's to complete.
3. When the commit fails even so, the error keeps git's whole output and says what happened and how to finish: the branch is merged, the item is done and archived, the changes are staged, and `flai accept <id>` commits them. It exits non-zero, so the host's journal records the acceptance as failed rather than done; read how `hostapi` builds the journal entry from the run's error to confirm it, and widen this task's touches only if it does not.

Update `flai accept`'s help (`Long`) for the resumed commit.

## Done when

- `flai/cmd/accept_lock_test.go`, against a real repository, reproduces I-0100: an acceptance whose commit meets `.git/index.lock` held for a moment completes, with one acceptance commit.
- With the lock held throughout (the helper's wait shortened), the acceptance fails, names git's error and `flai accept <id>`, and leaves the story done and archived; with the lock removed, `flai accept <id>` commits it, with the usual subject, and a third run is refused as already done.
- `flai accept <id> --dry-run` on such an item says it would commit the acceptance and changes nothing.
- The existing accept tests (`accept_resume_test.go` among them) pass, and `flai test flai/cmd flai/internal/preview` passes.

## Notes
