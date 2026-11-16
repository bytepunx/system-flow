---
id: T-3409
type: task
nature: remediation
title: A running close-out job records the head it runs at, and closeout reads a story's running job with that head and the step it is on
status: backlog
parent: S-0799
owner: alex
created: 2026-11-16T01:45:03Z
updated: 2026-11-16T01:45:03Z
transitions: []
stream: S-0799
tags: [flai, verify]
touches: [flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go]
---
# T-3409 A running close-out job records the head it runs at, and closeout reads a story's running job with that head and the step it is on

## Work

`closeout.Job` (`flai/internal/closeout/closeout.go`, lines 97-114) is written at a job's start with no commit and no step. It gets `Commit` only at its end (`runJob`, `flai/internal/serve/verifyjob.go`, line 377). So nothing can tell which head a running job verifies, or how far it has got.

- Give `Job` a `Head` (JSON `head`, omitted when empty): the story's head when the job's run began. In `runJob`, write the job again with `Head` set to `before` (line 357) once it is read, before the close-out starts. Keep `Commit` as the head it ended at.
- Add a reader in `closeout`, such as `Running(dir, story) (Progress, bool, error)`. It answers the story's job while `Running()` is true: its `ID`, `Started`, `Head`, and `Step`, the step the close-out is on. Read `Step` from the job's log (`LogPath`, line 151): `flai verify` prints `passed <step>`, or `passed tier <name>`, as each step ends. The step running is the one after the last such line, or the first step when there is none. When the log cannot say, leave `Step` empty and say so in the comment. If writing the step to the job file is simpler and reliable, do that instead, and say why in the narrative.
- Test the reader in `closeout_test.go`: a running job with a log that passed rebase and sync names the next step; an ended job answers false; a job written before this change, with no `head`, still reads. In `verifyjob_test.go`, a job's file holds `head` while it runs.

This task waits for none: every other task reads what it writes.

## Done when

- A running job's file names the head it runs at, and `closeout` answers a story's running job with its id, start, head, and step.
- `flai test flai/internal/closeout flai/internal/serve` passes.

## Notes
