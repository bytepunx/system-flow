---
id: T-1505
type: task
nature: remediation
title: Close I-0137 with what fixed it
status: done
parent: S-0375
owner: alex
created: 2026-10-13T16:37:08Z
updated: 2026-10-13T17:01:53Z
transitions:
  - to: ready
    at: 2026-10-13T17:01:49Z
    by: agent-S-0375
  - to: in-progress
    at: 2026-10-13T17:01:50Z
    by: agent-S-0375
  - to: done
    at: 2026-10-13T17:01:53Z
    by: agent-S-0375
stream: S-0375
tags: [issues]
touches: [design/issues/I-0137-a-host-crash-during-flai-story-start-leaves-the-story-s-branch-ref-worktree-and-narrative-as-empty-files-and-no-agent-can-repair-the-ref.md, design/issues/summary.md]
after: [T-1503, T-1504]
usage:
  source: log
  seconds: 3
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1770
      cache_read: 232943
      cache_write: 9596
      cost: 0.1423
---
# T-1505 Close I-0137 with what fixed it

## Work

Close the issue for the story's second criterion, from the story's worktree:

```bash
flai issue close I-0137 --reason "<what fixed it>"
```

The reason names the three fixes: `atomicfile.WriteFile` syncs before and after the rename (T-1500), `flai stream open` writes an empty narrative afresh (T-1501), and it repairs a broken story ref and a crashed worktree that hold no commits (T-1502), with the test that reproduces the crash (T-1503). It regenerates `design/issues/summary.md`.

Waits for T-1503 and T-1504: the issue is closed only once the fix is built, tested, and documented.

## Done when

- I-0137 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean for the story.

## Notes

Drafted by the planner.
