---
id: T-2846
type: task
nature: remediation
title: Confirm I-0309's tests under load and close I-0309
status: done
parent: S-0651
owner: alex
created: 2026-11-10T07:01:10Z
updated: 2026-11-10T07:51:34Z
transitions:
  - to: ready
    at: 2026-11-10T07:28:07Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:28:08Z
    by: system-flow
  - to: done
    at: 2026-11-10T07:51:34Z
    by: agent-S-0651
stream: S-0651
tags: [cli, testing]
touches: [design/issues/I-0309-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md, design/issues/summary.md]
after: [T-2839, T-2842]
usage:
  source: log
  seconds: 1406
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 4607
      cache_read: 1718125
      cache_write: 38727
      cost: 0.7457
---
# T-2846 Confirm I-0309's tests under load and close I-0309

## Work

- Run every test I-0309 names together under a CPU load like the one in its instances: `TestImportWithAHostRegistersTheProject`, `TestImportWithServeStoppedSaysToStartIt`, `TestMCPOverHTTPAnswersAWaitBeforeClaudeCodeGivesUp`, and `TestBacksOffWhenTheDashboardSaysTheProjectIsHeld`, with `-race` and `-count` of at least 10. The two import tests were fixed by S-0606's T-2709, which stamps `fakeRunning`'s state an hour ahead; this confirms that on the story's branch. Record the load, the command, and the result in the story's narrative.
- Close the issue with `flai issue close I-0309 --reason "<what fixed it>"`, naming S-0606's T-2709 for the import tests, T-2839 for the MCP-over-HTTP wait test, and T-2842 for the channel back-off test. It regenerates `design/issues/summary.md`.

It waits for T-2839 and T-2842: the reason names their fixes, and the run under load checks them.

## Done when

- The four tests pass together under load with `-race`, and the narrative records the run.
- I-0309 is closed with a reason naming the three fixes, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
