---
id: T-2269
type: task
nature: improvement
title: flai-cli.md, the flai stats guide, and the backup runbook say the stories' commits are kept in a file under .flai-cache between processes
status: done
parent: S-0531
owner: alex
created: 2026-10-28T18:24:45Z
updated: 2026-10-28T18:36:29Z
transitions:
  - to: ready
    at: 2026-10-28T18:26:35Z
    by: agent-S-0531
  - to: in-progress
    at: 2026-10-28T18:31:18Z
    by: agent-S-0531
  - to: done
    at: 2026-10-28T18:36:29Z
    by: agent-S-0531
stream: S-0531
tags: [flai, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/operators/runbooks/backup.md]
after: [T-2267]
usage:
  source: log
  seconds: 311
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 14
      output: 82
      cache_read: 610944
      cache_write: 44615
      cost: 0.222
---
# T-2269 flai-cli.md, the flai stats guide, and the backup runbook say the stories' commits are kept in a file under .flai-cache between processes

## Work

Write what T-2267 built in the three documents that describe `storygit.Committed`'s keeping:

- `design/system/flai-cli.md`, the paragraph on `storygit.Committed` (S-0500, I-0184): the kept file, its path, its key, what a miss is, and that a failed write never fails a call. Replace the sentence that a one-shot `flai stats` reads the log once.
- `docs/users/flai.md`, the touches drift paragraph under `flai stats`: a one-shot `flai stats` reads `git log` only when a branch tip moved since flai last read it.
- `docs/operators/runbooks/backup.md`, the row for a project's `.flai-cache/` apart from worktrees: add the kept stories' commits, made again.

Use the file name and behaviour T-2267 wrote; read its diff first.

This task waits for T-2267: it documents the file that task names. It runs beside the measurement task, which changes none of these files.

## Done when

- Each of the three documents says what the code does, with the file's path as T-2267 wrote it.
- `flai test` on the three files passes.

## Notes
