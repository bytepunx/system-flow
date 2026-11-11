---
id: T-3023
type: task
nature: remediation
title: I-0260 is closed with the reason that the MCP tool test lists a command only for the tier that failed
status: done
parent: S-0707
owner: alex
created: 2026-11-10T10:32:24Z
updated: 2026-11-11T18:34:54Z
transitions:
  - to: ready
    at: 2026-11-11T18:34:43Z
    by: agent-S-0707
  - to: in-progress
    at: 2026-11-11T18:34:46Z
    by: agent-S-0707
  - to: done
    at: 2026-11-11T18:34:54Z
    by: agent-S-0707
stream: S-0707
tags: [issues]
touches: [design/issues/I-0260-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md, design/issues/summary.md]
after: [T-3019, T-3022]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1358
      cache_read: 591337
      cache_write: 15559
      cost: 0.2699
---
# T-3023 I-0260 is closed with the reason that the MCP tool test lists a command only for the tier that failed

## Work

Close I-0260 with `flai issue close I-0260 --reason`, the reason naming the fix: the MCP tool `test` lists a tier's `command` only for the tier that failed, gives the others a `files` count, and fits its answer in 40,000 bytes, as the new ADR records, with the reproducing test in `flai/internal/mcpserver/test_test.go`. flai rewrites the issue's row in `design/issues/summary.md`. It waits for T-3019 and T-3022 so that the reason names the fix as built and the ADR by its number.

## Done when

- I-0260 is closed with that reason.
- `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` passes.

## Notes
