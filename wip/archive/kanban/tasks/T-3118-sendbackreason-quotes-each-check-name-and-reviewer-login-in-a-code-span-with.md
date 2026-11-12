---
id: T-3118
type: task
nature: remediation
title: sendBackReason quotes each check name and reviewer login in a code span with its control characters dropped
status: done
parent: S-0711
owner: alex
created: 2026-11-11T09:06:10Z
updated: 2026-11-11T22:02:20Z
transitions:
  - to: ready
    at: 2026-11-11T21:56:34Z
    by: agent-S-0711
  - to: in-progress
    at: 2026-11-11T21:56:38Z
    by: agent-S-0711
  - to: done
    at: 2026-11-11T22:02:20Z
    by: agent-S-0711
stream: S-0711
tags: [flai, go, serve, forge]
touches: [flai/internal/serve/forgestatus.go, flai/internal/serve/forgestatus_test.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 342
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 4128
      cache_read: 2540724
      cache_write: 66443
      cost: 1.1224
---
# T-3118 sendBackReason quotes each check name and reviewer login in a code span with its control characters dropped

## Work

Write the remedy first in the narrative's `## Decisions`, citing D-0505: the story's goal asks for a proposal from the issue's instances before the build.

In `flai/internal/serve/forgestatus.go`, make `sendBackReason` write each check name and reviewer login as a code span. Drop the control characters first, newlines and ANSI escapes among them, as `mdlint.Fenced` drops them for the thread's quote. Pick a backtick run longer than any in the name, so that a name holding a backtick still closes its span. The reason then reads such as ``the required check `ci` failed on 3f2a9c1d0b7e of pull request #12``. It reaches the move's reason and, through `sentBackFor`, the agent's start prompt. Leave `sendBackTitle` alone: `workitem.CleanTitle` cleans the thread's title already.

Add a test to `flai/internal/serve/forgestatus_test.go` that reproduces I-0340. Use a check name with a newline, a backtick, and words that read as an instruction. The reason must keep the name on one line, inside one code span. Update the expected reason strings in that file and in `flai/internal/serve/agents_test.go`.

It waits for no task: it is the first layer.

## Done when

- `sendBackReason` quotes each name in a code span, with its control characters dropped.
- The new test fails on the current code and passes with the fix.
- `flai test` on the three files passes.

## Notes
