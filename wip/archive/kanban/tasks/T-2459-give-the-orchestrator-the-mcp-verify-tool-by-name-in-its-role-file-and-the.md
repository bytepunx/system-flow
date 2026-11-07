---
id: T-2459
type: task
nature: remediation
title: Give the orchestrator the MCP verify tool by name, in its role file and the guard, before verify leaves the reads
status: done
parent: S-0556
owner: alex
created: 2026-11-05T01:36:58Z
updated: 2026-11-07T06:32:09Z
transitions:
  - to: ready
    at: 2026-11-07T06:20:23Z
    by: agent-S-0556
  - to: in-progress
    at: 2026-11-07T06:20:39Z
    by: agent-S-0556
  - to: done
    at: 2026-11-07T06:32:09Z
    by: agent-S-0556
stream: S-0556
tags: [flai, template]
touches: [flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/internal/roles/baseline/orchestrator.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md]
usage:
  source: log
  seconds: 654
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 27
      output: 2879
      cache_read: 3417077
      cache_write: 41264
      cost: 1.3052
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 65969
          cache_write: 1221
          cost: 0.016
---
# T-2459 Give the orchestrator the MCP verify tool by name, in its role file and the guard, before verify leaves the reads

## Work

First, propose the fix in the narrative's `## Decisions`, from I-0240's one instance (S-0546), and cite D-0446, which recorded why S-0546 kept `verify` among the reads. The proposal: the MCP tool `verify` now always asks flai serve for a close-out job (`flai/internal/mcpserver/verify.go`, ADR-0229), and that job commits in the story's worktree. So it leaves `guard.MCPReads`, and `flai verify --job` joins the forms a sub-agent may not run (the next task of this story). The story's own agent keeps both, since the guard passes its own calls. The orchestrator keeps the tool, because its role file's prompt tells it to verify a story in review with `verify`.

The orchestrator reaches `verify` today only through `MCPReads`: its role file's `tools:` list (`activity_log`, `inbox`, `thread_open`, `wait_for_events`, and two `own-thread` calls) does not name it, and `asRole` in `flai/internal/guard/roles.go` lets anything in `MCPReads` pass before it reads the file. So name it before it leaves the reads:

- add `verify` to the `tools:` list of the orchestrator's role file in all three copies: `flai/internal/roles/baseline/orchestrator.md`, `template/root/design/workflows/roles/orchestrator.md`, and `design/workflows/roles/orchestrator.md`, kept identical above the marker as the roles tests require
- add `verify` to `guard.MCPOrchestrates` in `flai/internal/guard/guard.go`, which mirrors the baseline file's tools for a project without a role file
- add a test to `flai/internal/guard/orchestrate_test.go`: the orchestrator's own MCP `verify` call passes with the role file held and without it

This task waits for nothing.

## Done when

- The orchestrator's role file names `verify` in its `tools:` list, in all three copies, and `guard.MCPOrchestrates` holds it
- The new test passes, and `flai test` on the changed paths passes
- The narrative's `## Decisions` holds the proposed fix, citing D-0446

## Notes

These decisions belong in the narrative's `## Decisions`. They are here because the narrative is read-only in this story's container and no flai tool writes that section.

- **Proposed fix (D-0446, I-0240):** `verify` leaves `guard.MCPReads`, and a sub-agent is refused `flai verify --job`. The story's agent keeps both, since the guard passes its own calls. The orchestrator keeps the MCP tool by name. This narrows what S-0546 kept, as D-0446 anticipated. Of the two calls S-0546 kept the tool for, the orchestrator's verify of a story in review stays through the tool, and a reader's look at the record stays through `flai verify --last`.
- **The orchestrator's command line:** it holds the MCP tool, not `flai verify --job`. Its prompt names only the tool and plain `flai verify`, so nothing it is told to do needs the flag.
- **The tool is in `tools:`, not under a permission:** this follows the plan and `guard.MCPOrchestrates`, which holds what the orchestrator may call whatever its permissions.
- **Plan reviewed and kept, four layers of one task each:** T-2459, then T-2460, then T-2461, then T-2462. T-2460 waits because the orchestrator must hold `verify` by name before it leaves the reads, and both change `guard.go`. T-2461 waits because it describes the refusals T-2460 builds. T-2462 waits so that the close reason describes what was built and documented. T-2462 is mine, since a sub-agent cannot close an issue.
- **Operator's note on S-0556:** the confirmation thread on the guard's read list comes before the close-out, after T-2462.
