---
id: T-0990
type: task
nature: remediation
title: Close I-0068 saying what fixed it
status: done
parent: S-0261
owner: alex
created: 2026-10-05T05:52:32Z
updated: 2026-10-06T23:18:29Z
transitions:
  - to: ready
    at: 2026-10-06T23:18:17Z
    by: agent-S-0261
  - to: in-progress
    at: 2026-10-06T23:18:18Z
    by: agent-S-0261
  - to: done
    at: 2026-10-06T23:18:29Z
    by: agent-S-0261
stream: S-0261
tags: [flai, docs]
touches: [design/issues/I-0068-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md, design/issues/summary.md]
after: [T-0988, T-0989]
usage:
  source: log
  seconds: 11
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 7
      output: 2905
      cache_read: 241389
      cache_write: 10044
      cost: 0.1681
---
# T-0990 Close I-0068 saying what fixed it

## Work

In the story's worktree, check the measurements T-0988's test makes: the MCP `prime` result for each pack in I-0068's instances, and for the planner's pack for S-0261, now fits a tool result.

Then close the issue with `flai issue close I-0068 --reason`. The reason names the ADR, the fix, and the tests that reproduce the cause. The command updates `design/issues/summary.md`.

It waits for T-0988 because the issue closes only once the reproduction test passes. It waits for T-0989 so that the reason can point at documents that already describe the fix.

## Done when

- I-0068 is closed, and its reason names the ADR, the fix, and the reproducing tests.
- `design/issues/summary.md` no longer lists I-0068 as open.
- `flai check --strict` and the markdown lint pass.

## Notes
