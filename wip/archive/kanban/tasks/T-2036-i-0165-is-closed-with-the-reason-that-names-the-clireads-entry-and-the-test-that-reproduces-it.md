---
id: T-2036
type: task
nature: improvement
title: I-0165 is closed with the reason that names the cliReads entry and the test that reproduces it
status: done
parent: S-0477
owner: alex
created: 2026-10-22T20:59:03Z
updated: 2026-10-22T21:19:53Z
transitions:
  - to: ready
    at: 2026-10-22T21:19:40Z
    by: agent-S-0477
  - to: in-progress
    at: 2026-10-22T21:19:40Z
    by: agent-S-0477
  - to: done
    at: 2026-10-22T21:19:53Z
    by: agent-S-0477
stream: S-0477
tags: [flai, guard, issues]
touches: [design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/issues/summary.md]
after: [T-2032, T-2034]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1926
      cache_read: 454037
      cache_write: 13158
      cost: 0.2219
---
# T-2036 I-0165 is closed with the reason that names the cliReads entry and the test that reproduces it

## Work

Waits for T-2032 and T-2034: the issue is closed only once the fix and its documentation are on the branch.

In the story's worktree, run `flai issue close I-0165 --reason "<reason>"`, the reason saying that `cliReads` in `flai/internal/guard/guard.go` now lists `components list` and `serve`, that `TestASubAgentRunsChecksButNotWrites` reproduces the refusal, and that `components config` stays refused. `flai issue close` rewrites the issue file and regenerates `design/issues/summary.md`; commit both with the task. If the plan thread settled that this story also closes I-0156 (the same cause, filed from S-0438), close it in the same way with a reason naming S-0477.

## Done when

- `design/issues/I-0165-...md` has `status: closed` and a remediation line with the reason
- `design/issues/summary.md` no longer lists I-0165 as open
- `flai check --strict` reports nothing on the issue files

## Notes

Issue commands write under the current checkout: run them in the story's worktree, not the main checkout.
