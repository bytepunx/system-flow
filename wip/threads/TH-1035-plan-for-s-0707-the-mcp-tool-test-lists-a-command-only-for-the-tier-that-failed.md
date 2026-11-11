---
id: TH-1035
title: "Plan for S-0707: the MCP tool test lists a command only for the tier that failed, three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0707-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md
  item: S-0707
status: resolved
participants: [planner-S-0707, alex]
created: 2026-11-10T10:33:13Z
updated: 2026-11-11T09:05:28Z
---

# TH-1035 Plan for S-0707: the MCP tool test lists a command only for the tier that failed, three tasks in three layers

On wip/kanban/stories/S-0707-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md.

## Entries

### 2026-11-10T10:33:13Z planner-S-0707
S-0707 is planned. It has three tasks in three layers, 8 touches, a 45m forecast, and a cost of delay of 12.50 USD a week.

## Proposed fix

The fix is in the MCP tool `test` alone (`flai/internal/mcpserver/test.go`):

- Only the tier that failed keeps its `command`. A tier that passed, or was not reached, gives `files` (a count) in its place.
- The answer fits 40,000 bytes. If it would not, `paths` is cut and `paths_omitted` counts the rest, as ADR-0179 and ADR-0198 do for `task_done`.
- `flai test --json`, the dashboard's `test.run`, and the verify record keep every tier's whole `command`.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3019: the projection, and a test that reproduces I-0260 with a `{files}` tier whose answer today passes 40,000 bytes | `flai/internal/mcpserver/test.go`, `flai/internal/mcpserver/test_test.go` |
| 2 | T-3022: a new ADR, plus updates to `design/system/flai-cli.md` and `docs/users/flai.md`; after T-3019, so they describe what it built | `design/adrs/README.md`, `design/system/flai-cli.md`, `docs/users/flai.md` |
| 3 | T-3023: close I-0260 with the reason; after T-3019 and T-3022, so the reason can name the ADR | the I-0260 file, `design/issues/summary.md` |

## Assumptions

- The MCP answer may differ from `flai test --json`, and its description says how. The other option is to cut `command` in `verify.Result` itself. That would shrink the CLI JSON and the verify record too, but it would break the `Command` assertions in `flai/cmd/test_test.go` and `flai/internal/verify/run_test.go`, and shell callers would lose the exact argv.
- The answer's shape is a decision worth an ADR, as it was for `task_done` in ADR-0179 and ADR-0198. That ADR is why the story keeps the folder touch `design/adrs`.
- The forecast is raised from flai's 18m to 45m. A story that changes `flai/` runs the integration and smoke tiers at close-out, which took 8 to 21 minutes in today's close-outs.
- The cost of delay is the 12.50 USD a week that flai works out from the 5m input. It stands unchanged.

## Changes I would propose

None: no task to split, merge, or drop. If you want the CLI's JSON cut as well, say so and I will widen T-3019 to `flai/internal/verify` and `flai/cmd/test.go`.

### 2026-11-11T09:05:28Z alex
Resolved.
