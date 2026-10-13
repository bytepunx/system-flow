---
id: T-1504
type: task
nature: remediation
title: The design, the user guide, and the restore runbook say how a start recovers from a crash
status: done
parent: S-0375
owner: alex
created: 2026-10-13T16:37:02Z
updated: 2026-10-13T17:01:42Z
transitions:
  - to: ready
    at: 2026-10-13T16:58:34Z
    by: agent-S-0375
  - to: in-progress
    at: 2026-10-13T16:58:34Z
    by: agent-S-0375
  - to: done
    at: 2026-10-13T17:01:42Z
    by: agent-S-0375
stream: S-0375
tags: [docs]
touches: [design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/runbooks/restore.md]
after: [T-1500, T-1501, T-1502]
usage:
  source: log
  seconds: 188
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 12491
      cache_read: 1643710
      cache_write: 67709
      cost: 1.0041
---
# T-1504 The design, the user guide, and the restore runbook say how a start recovers from a crash

## Work

Document what T-1500, T-1501, and T-1502 changed, in the same story as the documentation convention requires:

- `design/system/flai-cli.md`: where it says narratives and the index are written atomically (S-0100), add that `atomicfile.WriteFile` syncs the file and its folder, so a host crash leaves the old content or the new, never an empty file. In the story branch section, say what `flai stream open` repairs: a broken `story/S-nnnn` ref that resolves to no commit, and a worktree folder git does not know while the branch holds no commit; and what it refuses.
- `design/system/agent-narrative.md` § Recovery procedure: an empty narrative is written afresh by `flai stream open`, and the start can be run again after a crash.
- `docs/users/flai.md` § Starting a story: the same, in the user's terms, with the refusal's message and what to do.
- `docs/operators/runbooks/restore.md`: next to putting back a story worktree, say that `flai stream open S-nnnn` repairs a start a host crash interrupted, and when it refuses instead.

Waits for T-1500, T-1501, and T-1502, so that it describes what they built, not what was planned.

## Done when

- The four documents describe the durable writes and the repair as built.
- `flai test` on the four paths passes the markdown lint.

## Notes

Drafted by the planner. Paths and sections from the explorer's read of the design and docs.
