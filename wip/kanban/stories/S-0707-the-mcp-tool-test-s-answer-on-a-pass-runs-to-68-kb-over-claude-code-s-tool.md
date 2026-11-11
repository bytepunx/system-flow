---
id: S-0707
type: story
nature: improvement
title: The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given
status: backlog
owner: alex
created: 2026-11-10T09:11:42Z
updated: 2026-11-11T08:58:19Z
transitions: []
tags: [flai, mcp]
touches: [flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, design/adrs, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0260-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md, design/issues/summary.md]
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
      seconds: 176
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 284
          output: 8420
          cache_read: 23454370
          cache_write: 498415
          cost: 1.6318
        - model: claude-opus-5-5
          input: 674
          output: 84857
          cache_read: 51538468
          cache_write: 912167
          cost: 19.2585
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T09:11:42Z
  value: 12.5
  by: planner-S-0707
  at: 2026-11-10T10:32:46Z
forecast:
  duration: 45m
  delivery: 2026-11-11T21:49:00Z
  basis: "Its own forecast of 45m; 19th in the pull order with an in-progress limit of 3, behind S-0694, S-0698, S-0702, S-0706, S-0719, S-0724, S-0716, S-0689, S-0705, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0699, S-0701 and S-0703."
  by: flai
  at: 2026-11-11T08:58:19Z
---
# S-0707 The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given

## Goal

This story remediates [I-0260](../../../design/issues/I-0260-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md), "The MCP tool test's answer on a pass runs to 68 KB, over Claude Code's tool-result limit, because each tier's command lists every file it was given". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0260 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0260 is closed with `flai issue close I-0260 --reason` saying what fixed it

## Tasks
- T-3019 The MCP tool test answers a tier's command only when it failed, and its answer fits 40,000 bytes
- T-3022 An ADR, flai-cli.md, and flai.md say what the MCP tool test's answer leaves out
- T-3023 I-0260 is closed with the reason that the MCP tool test lists a command only for the tier that failed

## Notes

Cost of delay inputs set by flai from I-0260. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:44:19Z, 7.6 days before this story).

### Planning

Proposed solution, from I-0260's one instance (S-0425: nine paths, every tier passed, a 68,784-character answer, of which gofmt's `command` was 14 KB, flaiover-lint's 10 KB, and vitest's 10 KB). The MCP tool `test` returns `verify.Result` as it is (`flai/internal/mcpserver/test.go`), and each `TierResult.Command` carries the tier's argv with `{files}` expanded (`flai/internal/verify/select.go`, `selectByPaths`; `run.go`). The fix is a projection in the MCP tool alone:

- Each tier keeps its name, state, duration, findings, omitted count, and log. Only the tier that failed keeps its `command`, where a rerun needs the exact argv. A tier that passed, or was not reached, gives in its place `files`, the number of files it was given.
- The whole answer fits 40,000 bytes of JSON (`context.PartLimit`), as `task_done`'s and `story_start`'s do ([ADR-0179](../../../design/adrs/0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md)): when it would not, `paths` is cut to its first entries and `paths_omitted` counts the rest, the count-and-omit form of [ADR-0198](../../../design/adrs/0198-task-done-s-answer-lists-at-most-ten-of-the-check-s-notes-outside-the-story-and.md).
- `flai test --json`, the dashboard's `test.run`, and the verify record keep every tier's whole `command`: a shell caller can write it to a file, and `flai/cmd/test_test.go` and `flai/internal/verify/run_test.go` rely on it. The tool's description says the answer is `flai test --json`'s with these cuts, and points to it for the whole commands.
- A new ADR records the answer's shape, as ADR-0179 and ADR-0198 did for `task_done`.

Touches, each with where it came from:

| Touch | From |
|-------|------|
| `flai/internal/mcpserver/test.go` | design: I-0260 names the MCP tool `test`; the handler and its description are here |
| `flai/internal/mcpserver/test_test.go` | co-change (2 of 3 commits of `test.go`): the reproducing test, a `{files}` tier over enough files that today's answer passes 40,000 bytes |
| `design/adrs` | layout: the new ADR, whose file name `flai adr new` gives; kept as a folder touch because no task can name that file before it exists |
| `design/adrs/README.md` | layout: the ADR index flai generates from the ADR files |
| `design/system/flai-cli.md` | design: its `flai test` row names the MCP tool and says what `--json` prints |
| `docs/users/flai.md` | design: "Run the tests for what changed" lists the JSON's fields and describes the MCP tool |
| `design/issues/I-0260-the-mcp-tool-test-s-answer-on-a-pass-runs-to-68-kb-over-claude-code-s-tool.md` | criterion 2 |
| `design/issues/summary.md` | criterion 2: closing the issue rewrites its row |

Co-change also listed `flai/cmd/test.go`, `flai/cmd/verify.go`, `flai/cmd/verify_test.go`, and `flai/internal/mcpserver/verify.go` (2 of 3 each). None is a touch: the CLI's JSON and the verify answer do not change, and the MCP tool `verify` carries no steps. `docs/users/flai-reference.md` is not one either, since `flai test`'s help text does not change.

No `after`: the story reads no other story's output.

Tasks, in three layers: T-3019 (the projection and the reproducing test) first; T-3022 (the ADR and the docs) after it, so they describe the shape it built; T-3023 (closing I-0260) last, after both, so its reason names the fix and the ADR.

Forecast: `flai forecast` gave 18m, at 107 s per unit of size over 30 done medium improvement stories, size 10. Raised to 45m: a story that changes `flai/` runs the integration and smoke tiers at its close-out, which took 8m to 21m in today's close-outs (S-0672 8m14s, S-0682 10m24s, S-0669 20m44s), and the story writes an ADR beside the code. The delivery is flai's, 2026-11-10T18:38Z, moved 27m later to match.

Cost of delay: 12.50 USD a week, as `flai cod` works it out from the operator's input, 5m lost per 168h cycle at 150 USD an hour. It stands: I-0260 has one instance, and the workaround it records, writing `flai test --json` to a file and reading it with jq, took about 5m.
