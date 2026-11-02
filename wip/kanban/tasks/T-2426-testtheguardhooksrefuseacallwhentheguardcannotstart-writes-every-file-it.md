---
id: T-2426
type: task
nature: feature
title: TestTheGuardHooksRefuseACallWhenTheGuardCannotStart writes every file it executes before t.Parallel, so a parallel test's fork cannot make one text file busy
status: done
parent: S-0518
owner: alex
created: 2026-11-02T04:11:53Z
updated: 2026-11-02T04:12:11Z
transitions:
  - to: ready
    at: 2026-11-02T04:12:00Z
    by: agent-S-0518
  - to: in-progress
    at: 2026-11-02T04:12:01Z
    by: agent-S-0518
  - to: done
    at: 2026-11-02T04:12:11Z
    by: agent-S-0518
stream: S-0518
tags: []
touches: [flai/cmd/guard_test.go]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 643
      cache_read: 357010
      cache_write: 6283
      cost: 0.1345
---

# T-2426 TestTheGuardHooksRefuseACallWhenTheGuardCannotStart writes every file it executes before t.Parallel, so a parallel test's fork cannot make one text file busy

## Work

The close-out at 2026-11-02T04:07Z failed the integration tier in this story's own test: the template's edit hook met a fake guard that exited 1 but got `flai: Text file busy`, exit 126. The test wrote the fake `flai` while other parallel tests were forking, so a child of one of them held the write descriptor open when the hook executed the file (Go's ETXTBSY race). Write the copied scripts and every fake guard before `t.Parallel()`, and run the cases against those folders.

## Done when

The test writes no executable after `t.Parallel()`, `flai test` passes on `flai/cmd/guard_test.go`, and the close-out passes.

## Notes

Found by the close-out after T-2425, so it is a task of its own.
