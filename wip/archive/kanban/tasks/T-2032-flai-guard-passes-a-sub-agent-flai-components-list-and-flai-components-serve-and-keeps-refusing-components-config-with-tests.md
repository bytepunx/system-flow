---
id: T-2032
type: task
nature: improvement
title: flai guard passes a sub-agent flai components list and flai components serve and keeps refusing components config, with tests
status: done
parent: S-0477
owner: alex
created: 2026-10-22T20:58:21Z
updated: 2026-10-22T21:13:49Z
transitions:
  - to: ready
    at: 2026-10-22T21:13:37Z
    by: agent-S-0477
  - to: in-progress
    at: 2026-10-22T21:13:38Z
    by: agent-S-0477
  - to: done
    at: 2026-10-22T21:13:49Z
    by: agent-S-0477
stream: S-0477
tags: [flai, guard, components]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 435
      cache_read: 102639
      cache_write: 2974
      cost: 0.0502
---
# T-2032 flai guard passes a sub-agent flai components list and flai components serve and keeps refusing components config, with tests

## Work

Add `"components": {"list", "serve"}` to `cliReads` in `flai/internal/guard/guard.go`, beside the other reads, with a comment saying why `serve` counts as one: it runs the sub-project's own Vite dev server on 127.0.0.1 and writes only under the git-ignored `.flai-cache` and the worktree's `node_modules`, as `flai test` runs processes and writes caches already (I-0165, I-0156). Leave `components config` out: it writes a sub-project's `components` key in `system-flow.yaml`.

In `flai/internal/guard/guard_test.go`, add `components` to the test guard `g`'s `Commands` list (line 16), without which `refuse()` skips the command and the tests pass whatever `cliReads` says. In `TestASubAgentRunsChecksButNotWrites`, add to `allowed` the forms the issues' instances ran: `flai components list`, `flai components list --project flaiover --json`, `scripts/flai.sh components serve --project flaiover --story S-0449`, and `bin/flai components serve`; add to `refused` `flai components config flaiover --setup src/app.css`. Check the first allowed case fails before the fix. In `flai/cmd/guard_test.go` `TestGuard`, add one end-to-end case through the real command list: a verifier's `flai components list` passes.

This task waits for nothing. Run `flai test` on the files changed.

## Done when

- A sub-agent's `flai components list` and `flai components serve`, in any form, pass `flai guard`, through `flai`, `bin/flai`, and `scripts/flai.sh`
- A sub-agent's `flai components config` is still refused with the ADR-0060 reason
- The new test cases fail on the old `cliReads` and pass on the new one, and `flai test flai/internal/guard flai/cmd` passes

## Notes

The guard's write rule (ADR-0141, `writes.go`) does not judge `flai` commands, so once `serve` is a read a read-only verifier may start the dev server; the story's agent records that trade-off in `## Decisions`.
