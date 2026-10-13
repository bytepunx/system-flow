---
id: T-1500
type: task
nature: remediation
title: atomicfile.WriteFile syncs the file before the rename and the folder after it
status: done
parent: S-0375
owner: alex
created: 2026-10-13T16:36:23Z
updated: 2026-10-13T16:58:13Z
transitions:
  - to: ready
    at: 2026-10-13T16:39:12Z
    by: agent-S-0375
  - to: in-progress
    at: 2026-10-13T16:39:12Z
    by: agent-S-0375
  - to: done
    at: 2026-10-13T16:58:13Z
    by: agent-S-0375
stream: S-0375
tags: [flai]
touches: [flai/internal/atomicfile/atomicfile.go, flai/internal/atomicfile/atomicfile_test.go]
usage:
  source: log
  seconds: 1141
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 8403
      cache_read: 1105728
      cache_write: 45548
      cost: 0.6755
---
# T-1500 atomicfile.WriteFile syncs the file before the rename and the folder after it

## Work

`atomicfile.WriteFile` (`flai/internal/atomicfile/atomicfile.go`) writes a temporary file and renames it into place, but never calls `Sync`. A host crash can keep the rename and lose the data, which is how I-0137 left `wip/agents/S-0337.md` at zero bytes. Narratives, the narrative index, work items, and threads all write through it.

- Call `Sync` on the temporary file after the write and before the close and rename.
- After the rename, open the folder and `Sync` it, so the rename itself is durable. Ignore an error from the folder sync only where the platform does not support it (such as `EINVAL` on some filesystems), and say so in a comment.
- Keep the function's signature, so no caller changes.
- Add a test that the written file has the content and mode given and that no temporary file is left in the folder. A crash cannot be simulated in a unit test; the test guards the behaviour around the new calls.

Waits for nothing: no other task of the story touches these files.

## Done when

- `WriteFile` syncs the file before the rename and the folder after it.
- `flai test flai/internal/atomicfile` passes.

## Notes

Drafted by the planner from I-0137's instance and the explorer's read of `atomicfile.go` lines 14-40.
