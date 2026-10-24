---
id: S-0433
type: story
nature: feature
title: An operator's answer on a thread becomes a digest entry at once, before the curator's next run, so the next story that meets the question finds it
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:01Z
updated: 2026-10-24T00:28:42Z
transitions: []
tags: [cli]
topics: [cli, orchestration, planning]
touches: [design/system/agent-narrative.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/thread.go, flai/cmd/thread_test.go, flai/internal/digest/onanswer.go, flai/internal/digest/onanswer_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
after: [S-0432]
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
    - kind: orchestrator
      seconds: 381
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 50
          output: 782
          cache_read: 6168673
          cache_write: 15792
          cost: 1.3603
        - model: claude-sonnet-5-5
          input: 4
          output: 19
          cache_read: 41178
          cache_write: 15419
          cost: 0.039
cost_of_delay:
  value: 46.88
  by: planner-S-0433
  at: 2026-10-23T17:17:11Z
forecast:
  duration: 35m
  delivery: 2026-10-24T15:16:00Z
  basis: "Its own forecast of 35m; 22nd in the pull order with an in-progress limit of 5, behind S-0431, S-0486, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425 and S-0432."
  by: flai
  at: 2026-10-24T00:28:42Z
finalized:
  by: orchestrator
  at: 2026-10-23T17:18:00Z
---
# S-0433 An operator's answer on a thread becomes a digest entry at once, before the curator's next run, so the next story that meets the question finds it

## Goal

The question the operator just answered is not asked again while the curator waits for its next run. When the operator (TH-0476 question 5) replies to a thread someone else opened, confirms a recommendation, or answers otherwise than a pending recommendation, flai writes the entry there and then, from the CLI, the MCP tool, or the dashboard, with the rule S-0429's seed uses: an `answer` entry, or a `recommendation` entry that says what was overruled and why, citing the thread, with the topics and paths of its anchored item. A later answer on the same thread updates that entry rather than writing a second. The curator's next run merges it with what the digest already holds.

## Acceptance criteria

- [ ] `threads.ReplyWith` and `threads.Confirm`, or the one place every answer passes, write the entry through `flai/internal/digest` when the author is the operator and the thread is not the operator's own; a thread entry that is a recommendation, an agent's reply, and the operator's reply on their own thread write none.
- [ ] `flai thread reply`, `flai thread confirm`, the MCP `thread_reply`, and the host API's `thread.reply` and `thread.confirm` each write it, shown by a test on each path; the reply's result names the entry written or updated.
- [ ] A second answer by the operator on the same thread edits the entry that cites the thread rather than writing another; an entry the curator has merged into another (its source moved) is not written again.
- [ ] A failure to write the entry never refuses the answer: the answer is saved, the failure is logged as a warning, and the next `flai digest seed` or curator run picks the answer up.
- [ ] The entry passes `digest.secret`; when the answer matches a secret's shape, no entry is written and the warning says so, and the thread is left as it is.
- [ ] `design/system/strategic-agents.md`, `design/system/agent-narrative.md`, `docs/users/flai.md`, and `docs/users/flaiover.md` say that an operator's answer becomes an entry at once.

## Tasks
- T-1837 digest.OnAnswer writes or updates the entry for an operator's answer on a thread, and never fails the answer
- T-1838 The MCP thread_reply writes the operator's answer as a digest entry and names it in its result
- T-1839 The host API's thread.reply and thread.confirm write the operator's answer from the dashboard as a digest entry
- T-1840 flai thread reply and flai thread confirm write the operator's answer as a digest entry and name it in their result
- T-1841 The design and the user docs say that an operator's answer on a thread becomes a digest entry at once

## Notes

- Drafted by the planner for E-0025, after S-0432 in the epic's order. It needs only S-0429's rule, so it could follow S-0429; the plan thread proposes that.
- Writes happen in the main checkout, where threads are written, and are not committed, as a thread entry is not.

### Planning

Revisited by planner-S-0433 on 2026-10-23, with S-0428 in progress (its `flai/internal/digest` on its branch, not yet on main) and E-0025's cost of delay input set on TH-0476.

- **Touches.** The 12 declared touches stand, file by file, and one is added. Every touch is on a task.
  - **Declared, from layout:** a new `flai/internal/digest/onanswer.go` with its test (T-1837); the two places that save an operator's answer, `afterThread` in `flai/cmd/thread.go`, which `reply` and `confirm` share, and `threadReply` in `flai/internal/mcpserver/server.go`, each with its test (T-1840, T-1838). The host API's `thread.reply` and `thread.confirm` in `flai/internal/hostapi/writes.go` do not save a thread themselves: they run `flai thread reply --by=<owner> --json` and `flai thread confirm --by=<owner> --json` through `ExecRunner` and pass the CLI's JSON through, so the dashboard's entry is written by the CLI's code, and `writes.go` changes only if the result needs shaping; `writes_test.go` shows the path (T-1839). `threads.ReplyWith` and `threads.Confirm` are left alone: S-0428's `digest.go` imports `threads`, so `threads` cannot import `digest`.
  - **Declared, from design:** `design/system/strategic-agents.md`, `design/system/agent-narrative.md`, `docs/users/flai.md`, and `docs/users/flaiover.md`, named by criterion 6 (T-1841).
  - **Co-change, added:** `design/system/flai-cli.md` (48% of the commits that changed the declared files), on T-1841: its `## Commands` section describes `flai thread reply` and `confirm`, the MCP `thread_reply`, and the host API's `thread.reply` and `thread.confirm`, whose results gain `digest`. It is in `claims.shared`, as `docs/users/flai.md` is, so it holds no story.
  - **Co-change, not added:** `design/system/flaiover-dashboard.md` (31%), since the dashboard's pages do not change; `docs/users/flai-reference.md` (20%), since no flag is added, only a field of `--json`; `docs/operators/index.md` and `settings.md` (18%, 14%), since no setting is added; `flai/internal/mcpserver/folder.go` (12%), since no tool is added or registered; `flaiover/src/lib/server/agent.ts` (8%), since the dashboard need not show the entry.
  - **No folder touch.** The entries are written under `design/digest` in the main checkout at run time, as a thread entry is, not on the story's branch.
- **Forecast.** Kept at flai's 35m: the median of 107 s per unit of size over 85 done feature stories on claude-opus-5-5 in the large band, for size 19 (6 criteria, 13 touches). The work is one function and two call sites, with the dashboard path covered by the CLI's; the docs in five files are short additions. Delivery is flai's play-out, 25th in the pull order, moved as the order ahead changes.
- **Cost of delay.** 46.88 USD a week, from `flai cod`: the story's share by forecast (35m of 7h28m) of E-0025's 600 USD a week, which the operator set on TH-0476 as 4h lost per 168h cycle at 150 USD an hour. No input of the story's own; none is needed while it shares the epic's.
- **Who is the operator.** The manifest's `owner` (TH-0476 question 5), compared with the entry's author: the dashboard replies as `--by=<owner>`, the CLI as `--by`, else `FLAI_AGENT`, else the config author, and the MCP server as `--agent`, else `FLAI_AGENT`. An agent `flai serve` started is never the owner.
- **Layers.** T-1837; then T-1838 and T-1840, which touch different files; then T-1839, which waits for T-1840 because the dashboard path runs the CLI's code; then T-1841.
