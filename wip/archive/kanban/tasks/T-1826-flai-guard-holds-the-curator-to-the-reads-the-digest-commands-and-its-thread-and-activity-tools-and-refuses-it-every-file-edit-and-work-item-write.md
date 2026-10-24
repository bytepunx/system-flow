---
id: T-1826
type: task
nature: feature
title: flai guard holds the curator to the reads, the digest commands, and its thread and activity tools, and refuses it every file edit and work item write
status: done
parent: S-0431
owner: alex
created: 2026-10-16T06:05:50Z
updated: 2026-10-24T00:13:55Z
transitions:
  - to: ready
    at: 2026-10-24T00:01:38Z
    by: agent-S-0431
  - to: in-progress
    at: 2026-10-24T00:01:39Z
    by: agent-S-0431
  - to: done
    at: 2026-10-24T00:13:55Z
    by: agent-S-0431
stream: S-0431
tags: [cli]
touches: [flai/internal/guard/guard.go, flai/internal/guard/call.go, flai/internal/guard/authors.go, flai/internal/guard/shared.go, flai/internal/guard/curate_test.go, flai/cmd/guard.go]
after: [T-1824]
usage:
  source: log
  seconds: 736
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 14599
      cache_read: 2861703
      cache_write: 94345
      cost: 1.4332
---
# T-1826 flai guard holds the curator to the reads, the digest commands, and its thread and activity tools, and refuses it every file edit and work item write

## Work

In `flai/internal/guard/guard.go`, with `FLAI_ROLE=curate`, check the session's own calls (`Guard.curate`):

- Pass the reads a sub-agent has; `flai digest new`, `edit`, `harvest`, `list`, and `show` (`cliCurates`); and the MCP tools `inbox`, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events` beside the reads (`MCPCurates`).
- Refuse `Edit`, `Write`, `MultiEdit`, and `NotebookEdit` on every file; `item_new`, `item_edit`, `item_move`, `issue_new`, `issue_bump`, and every other flai tool and command that writes, `curate`, `analyze`, and `plan` among them; and git's writes.
- Hold it to its own name, as ADR-0143 holds the analyzer.

Each refusal says the curator writes entries only with `flai digest new` and `edit`. `flai/cmd/guard.go` passes what the curator needs.

Waits for T-1824, which makes the role known. Runs beside the harness task.

## Done when

`curate_test.go` covers each pass and each refusal above, and a call naming another author; `flai test flai/internal/guard flai/cmd` passes.

## Notes

The `.claude/settings.json` entry that runs the guard for `curate` is the definitions task's.
