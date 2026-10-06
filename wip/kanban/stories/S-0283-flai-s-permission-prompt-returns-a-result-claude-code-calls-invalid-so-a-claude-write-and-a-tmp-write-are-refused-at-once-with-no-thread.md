---
id: S-0283
type: story
nature: remediation
title: flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread
status: ready
owner: alex
created: 2026-10-06T03:45:18Z
updated: 2026-10-06T09:56:49Z
transitions:
  - to: ready
    at: 2026-10-06T06:20:20Z
    by: alex
tags: []
topics: [cli]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/folder.go, docs/users/flai.md, design/issues/I-0082-flai-s-permission-prompt-returns-a-result-claude-code-calls-invalid-so-a-claude-write-and-a-tmp-write-are-refused-at-once-with-no-thread.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 345
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 186
          output: 37
          cache_read: 1132049
          cache_write: 80274
          cost: 0.2433
        - model: claude-opus-5-5
          input: 129
          output: 20031
          cache_read: 3364565
          cache_write: 164307
          cost: 2.3277
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-06T03:45:18Z
  value: 25
  by: planner-S-0283
  at: 2026-10-06T06:27:36Z
forecast:
  duration: 25m
  delivery: 2026-10-06T10:50:00Z
  basis: "Its own forecast of 25m; 1st in the pull order with an in-progress limit of 3, with nothing ahead of it."
  by: flai
  at: 2026-10-06T09:56:49Z
finalized:
  by: alex
  at: 2026-10-06T06:19:28Z
---
# S-0283 flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread

## Goal

This story remediates [I-0082](../../../design/issues/I-0082-flai-s-permission-prompt-returns-a-result-claude-code-calls-invalid-so-a-claude-write-and-a-tmp-write-are-refused-at-once-with-no-thread.md), "flai's permission_prompt returns a result Claude Code calls invalid, so a .claude/ write and a /tmp write are refused at once with no thread". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0082 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0082 is closed with `flai issue close I-0082 --reason` saying what fixed it

## Tasks
- T-1001 permission_prompt answers Claude Code with one text block and no structured content, with a test reproducing I-0082
- T-1002 The users' guide says a write permission_prompt refuses reaches the agent with flai's reason
- T-1003 Close I-0082 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0082. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T03:19:21Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0283 on 2026-10-06.

Likely cause, for T-1001 to confirm: `flai/internal/mcpserver/folder.go` registers `permission_prompt` with `mcp.AddTool` and a typed handler (`ToolHandlerFor[PermissionIn, PermissionOut]`). The go-sdk v1.8.0 then gives the tool an output schema and sends `structuredContent` with the result. Claude Code wants one text block whose text is the JSON decision. The tests read only `Content[0]`, so they pass.

Touches, none declared before planning:

- `flai/internal/mcpserver/permission.go`: design. ADR-0086 puts the handler here, and its `PermissionOut`, `allow`, and `deny` build the result.
- `flai/internal/mcpserver/permission_test.go`: co-change. `flai touches suggest` from `permission.go` and `folder.go` lists it, and the first criterion asks for a test.
- `flai/internal/mcpserver/folder.go`: layout. It registers the tool, and the output schema comes from that registration.
- `docs/users/flai.md`: co-change (38% of the 45 commits) and design. Its `Writes under .claude/` section says what a refused write does.
- The I-0082 file and `design/issues/summary.md`: layout. The second criterion closes I-0082, which rewrites the issue and the summary.

Left out: `server.go` (69% co-change), because the fix is in the tool's registration and result, not the server. `design/system/flai-cli.md` and `design/adrs`, because the fix keeps ADR-0086's decision and changes only the wire shape.

Topic `cli` added: ADR-0086 and the `flai mcp` tools are filed under it.

Forecast 25m, delivery 2026-10-06T11:09Z. `flai forecast` gave 16m (median 116 s per unit over 10 done medium-band remediation stories on claude-opus-5-5, times size 8: 2 criteria and 6 touches), delivering at 11:00Z, 10th in the pull order. Raised by 9m because the cause must be confirmed over the wire, and S-0219's no-thread `.claude/` refusal traced, before the fix; the criteria count neither. Delivery shifted by the same 9m.

Cost of delay 25 USD a week, as `flai cod` gives from the operator's input (10m lost per 168h cycle at 150 USD an hour). It stands as computed. But I-0082 now counts 2 instances, and every refusal the handler gives is invalid while this lasts. Each `.claude/` write falls back to the operator pasting a file, so a recurrence costs more than the input counts.

Overlap: S-0284 (I-0081) also changes `permission.go`, `permission_test.go`, `docs/users/flai.md`, and `design/issues/summary.md`. S-0220, in progress, touches `flai/internal/mcpserver`. flai's overlap hold runs them one after another, so no `after` is set.

### Cause confirmed against Claude Code

Confirmed on 2026-10-06 against Claude Code 2.1.290, in the operator's session with Claude, outside this story. A stub MCP server stood in for `permission_prompt` in a headless `claude -p` run with `--permission-mode acceptEdits`. The run asked for one Write under a `.claude/agents/` folder in a scratch project, and the stub answered allow.

- With one text block and no output schema on the tool, the file was written.
- With an output schema on the tool and `structuredContent` beside the text block, as flai sends today, Claude Code answered `Permission prompt tool returned an invalid result. Expected a single text block param with type="text" and a string text value.` and wrote nothing.

flai 1.31.4's own answer to a deny, read over stdio, holds the text block and `structuredContent` both. The run did not tell the output schema from `structuredContent`; T-1001 removes both. So the planner's likely cause stands. T-1001 need not confirm it against Claude Code again; its Go test over the wire is still wanted.

The operator asked for this story at the top of ready (alex, 2026-10-06, in conversation).
