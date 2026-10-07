---
id: T-1100
type: task
nature: improvement
title: The conventions and the harness prompt send a story's agent to the one call with --commit
status: done
parent: S-0275
owner: alex
created: 2026-10-06T22:53:11Z
updated: 2026-10-07T08:44:54Z
transitions:
  - to: ready
    at: 2026-10-07T08:34:28Z
    by: agent-S-0275
  - to: in-progress
    at: 2026-10-07T08:34:28Z
    by: agent-S-0275
  - to: done
    at: 2026-10-07T08:44:54Z
    by: agent-S-0275
stream: S-0275
tags: [conventions]
touches: [design/conventions/continuous-improvement.md, design/conventions/decisions.md, template/root/design/conventions/continuous-improvement.md, template/root/design/conventions/decisions.md, template/CHANGELOG.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/root/design/adrs/README.md]
after: [T-1077, T-1083]
usage:
  source: log
  seconds: 626
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 35
      output: 14332
      cache_read: 2208838
      cache_write: 60320
      cost: 1.0782
---
# T-1100 The conventions and the harness prompt send a story's agent to the one call with --commit

## Work

Criterion 3, for the conventions and the harness prompt. It waits for T-1077 and T-1083, so that it names the flags and tools as they were built. It shares no path with T-1093, so the two run together.

- In `continuous-improvement.md` and `decisions.md`, tell a story's agent to record an issue or an ADR in one call from its worktree: `flai issue new|bump|close --commit` and `flai adr new --commit`, or `issue_new`, `issue_bump`, `issue_close`, and `adr_new` with `commit`. Drop any instruction to follow them with `git add`, `git commit`, or `flai touches`. Change the baseline in the `template/root` copies, then the project's copies, as `conventions.md` says, and keep any project additions below the marker.
- Add a `template/CHANGELOG.md` entry for the convention change.
- In `flai/internal/harness/harness.go`, change the story agent's prompt, the line that names `flai issue new` and `bump`, so it gives `--commit`. Leave the analyzer's prompt alone: it works in the main checkout.

## Done when

- The four convention files and the changelog say the one call, and no convention tells a story's agent to commit or claim an issue or ADR file by hand.
- `harness_test.go` asserts the prompt names `--commit`.
- `scripts/flai-test.sh` passes, and `flai check --strict` is clean on the changed files.

## Notes

Written by the planner. If other conventions or `.claude/agents` files tell a story's agent to commit an issue or ADR file by hand, widen this task's touches to them.
