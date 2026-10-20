---
id: T-1911
type: task
nature: feature
title: flai accept --batch leaves the main checkout's uncommitted files outside wip out of its refusal, its lint fix, and its commit
status: done
parent: S-0402
owner: alex
created: 2026-10-20T00:34:59Z
updated: 2026-10-20T00:35:10Z
transitions:
  - to: ready
    at: 2026-10-20T00:35:07Z
    by: agent-S-0402
  - to: in-progress
    at: 2026-10-20T00:35:07Z
    by: agent-S-0402
  - to: done
    at: 2026-10-20T00:35:09Z
    by: agent-S-0402
stream: S-0402
tags: []
touches: [flai/cmd/accept_markdown.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/cmd/accept.go]
usage:
  source: log
  seconds: 2
  estimated: true
  priced_by: harness
  models: []
---
# T-1911 flai accept --batch leaves the main checkout's uncommitted files outside wip out of its refusal, its lint fix, and its commit

## Work

- Under `--batch` only, `flai accept` no longer refuses on the main checkout's uncommitted files outside `wip/`. It reads them before anything changes and leaves them as they are.
- Leave those files out of the lint fix and its refusal over the main checkout's uncommitted files, and out of the acceptance commit's `git add`, so that the commit holds only the acceptance and the operator's files stay uncommitted.
- `flai accept` without `--batch` keeps its refusal and `--yes` as today.
- Test a `--batch` acceptance with an uncommitted analysis report and issue file, as on TH-0462, and the refusal by hand.
- First layer: it waits for no task. Either answer to TH-0521 needs it, since the acceptance commit is made in the main checkout either way.
- It also changes three lines of `flai/cmd/accept.go`, which S-0403 claims, under the split agreed on MS-0148. A share between two stories in progress cannot be recorded (I-0169), so `flai/cmd/accept.go` joins the touches when the task closes.

## Done when

A `--batch` acceptance lands with those files left uncommitted and out of its commit, the acceptance by hand refuses as before, and the tests pass.

## Notes
