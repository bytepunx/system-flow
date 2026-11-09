---
id: T-2708
type: task
nature: remediation
title: Close I-0230 with the reason S-0613's fix gives
status: done
parent: S-0613
owner: alex
created: 2026-11-09T08:12:23Z
updated: 2026-11-09T09:06:16Z
transitions:
  - to: ready
    at: 2026-11-09T09:05:57Z
    by: agent-S-0613
  - to: in-progress
    at: 2026-11-09T09:05:58Z
    by: agent-S-0613
  - to: done
    at: 2026-11-09T09:06:16Z
    by: agent-S-0613
stream: S-0613
tags: [flai, issues]
touches: [design/issues/I-0230-flai-guard-refuses-the-orchestrator-a-read-or-a-permitted-flai-command-when-the-command-carries-a-stderr-redirection-such-as-2-1-or-2-dev-null.md, design/issues/summary.md]
after: [T-2706]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 3
      output: 26
      cache_read: 428468
      cache_write: 1609
      cost: 0.1564
---
# T-2708 Close I-0230 with the reason S-0613's fix gives

## Work

In the story's worktree, run `flai issue close I-0230 --reason "<reason>"`, the reason saying that flai guard now judges the orchestrator's and every role file's shell commands by their arguments, with redirections and their targets left out, so `2>&1` and `2>/dev/null` no longer turn a permitted call into a refused one, fixed by S-0613 with the tests in `flai/internal/guard/orchestrate_test.go`. The command rewrites the issue file and `design/issues/summary.md`.

It waits for T-2706 because the reason names the fix and its tests, which T-2706 writes and proves.

## Done when

- [ ] I-0230's status is closed, with that reason, and `design/issues/summary.md` no longer lists it as open

## Notes
