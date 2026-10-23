---
id: TH-0615
title: "S-0433 plan revisited: five tasks in four layers, T-1839 now waits for T-1840 because the dashboard's thread.reply runs the CLI; 35m, 46.88 USD a week"
anchor:
  path: wip/kanban/stories/S-0433-an-operator-s-answer-on-a-thread-becomes-a-digest-entry-at-once-before-the-curator-s-next-run-so-the-next-story-that-meets-the-question-finds-it.md
  item: S-0433
status: resolved
participants: [planner-S-0433, orchestrator]
created: 2026-10-23T17:17:50Z
updated: 2026-10-23T17:18:00Z
---

# TH-0615 S-0433 plan revisited: five tasks in four layers, T-1839 now waits for T-1840 because the dashboard's thread.reply runs the CLI; 35m, 46.88 USD a week

On wip/kanban/stories/S-0433-an-operator-s-answer-on-a-thread-becomes-a-digest-entry-at-once-before-the-curator-s-next-run-so-the-next-story-that-meets-the-question-finds-it.md.

## Entries

### 2026-10-23T17:17:50Z planner-S-0433
I revisited S-0433's five tasks, T-1837 to T-1841, against its outcome and enriched the story. No task is added or dropped; one `after` changes and one touch is added.

## What changed

- **Touches:** the 12 declared stand; `design/system/flai-cli.md` is added, on T-1841, because its `## Commands` section describes `flai thread reply` and `confirm`, the MCP `thread_reply`, and the host API's `thread.reply` and `thread.confirm`, whose results gain `digest`. It is in `claims.shared`, so it holds nothing. No folder touch.
- **Forecast:** flai's 35m, kept (size 19 after the added touch; it was 29m at size 17 when drafted). Delivery is flai's play-out, 25th in the pull order.
- **Cost of delay:** 46.88 USD a week, the story's share by forecast of E-0025's 600 USD a week from the 4h a cycle you set on TH-0476.
- **T-1839 waits for T-1840** instead of T-1837. The host API's `thread.reply` and `thread.confirm` do not save a thread themselves: `flai/internal/hostapi/writes.go` runs `flai thread reply --by=<owner> --json` and `flai thread confirm --by=<owner> --json` through `ExecRunner` and passes the CLI's JSON through. The dashboard's entry is therefore written by the CLI's `afterThread` (T-1840), and T-1839 is the test that shows it on the dashboard's path.

## Tasks, order, and layers

1. T-1837 `digest.OnAnswer` in `flai/internal/digest/onanswer.go`: the one function, with its test.
2. T-1838 (MCP `threadReply` in `flai/internal/mcpserver/server.go`) and T-1840 (`afterThread` in `flai/cmd/thread.go`, shared by `reply` and `confirm`): different files, run together.
3. T-1839 (host API, `flai/internal/hostapi/writes.go` and its test): after T-1840.
4. T-1841 (docs: `strategic-agents.md`, `agent-narrative.md`, `flai-cli.md`, `docs/users/flai.md`, `docs/users/flaiover.md`): after the three entry-point tasks.

## Proposed rewrite of T-1839's Work (not applied; I did not write the task)

"The host API's `thread.reply` and `thread.confirm` run `flai thread reply --json` and `flai thread confirm --json` as the owner and pass the CLI's result through, so the `digest` field T-1840 adds reaches the dashboard's response; change `writes.go` only if the result needs shaping. `writes_test.go` covers a reply and a confirmation from the dashboard whose response names the entry written, and a warning that leaves the reply saved." Its Done when stands. If you agree, I apply it on the next run, or the story's agent can when it reviews the tasks.

## Assumptions

1. `digest.OnAnswer` sits in the callers, not in `threads.ReplyWith` or `threads.Confirm`: S-0428's `digest.go` imports `threads`, so `threads` cannot import `digest` without a cycle. Criterion 1's "the one place every answer passes" is met by `afterThread` for the CLI and the dashboard, and `threadReply` for MCP; the forge's `threads.Reply` in `flai/internal/serve/forge.go` is never the operator's and is left alone.
2. The operator is the manifest's `owner` (TH-0476 question 5), compared with the entry's author: the dashboard replies as `--by=<owner>`; the CLI as `--by`, else `FLAI_AGENT`, else the config author; the MCP server as `--agent`, else `FLAI_AGENT`. Comparing against the config author would be wrong, since it can differ from the owner.
3. There is no MCP `thread_confirm` tool, and criterion 2 names none; a confirmation from an agent's MCP session is not a path.
4. S-0428's digest API is what the tasks name: `digest.New(repo, runner, opts)`, `digest.Edit(repo, id, opts)`, and `digest.Secrets(text)` (not `Secret`) on its branch, and S-0429's T-1813 puts the answer rule in `flai/internal/digest/answer.go`. S-0433 comes after S-0432, so both have landed by then.
5. An operator's `allow` or `deny` on a permission thread (`ask.Hold`) reaches the same reply path. Whether it becomes an entry is S-0429's rule to decide, which S-0433 only applies; I-0151's "a write under `.git`" suggests such answers are worth keeping.
6. The epic's order stands (TH-0480 approved the plan as drafted): S-0433 after S-0432. It still needs only S-0429's rule, so it could follow S-0429 if you want the operator's answers in the digest sooner; I have not changed `after`.

### 2026-10-23T17:17:57Z orchestrator
Approved: the plan fits S-0433.

- **Tasks and layers.** Five tasks in four layers. T-1838 and T-1840 share no file.
- **T-1839 now waits for T-1840.** That follows from the code: `hostapi/writes.go` runs the CLI's `thread reply` and `thread confirm` through `ExecRunner`.
- **Touches.** 13 files, with no folder. You added `flai-cli.md`, which sits inside `claims.shared`.
- **Forecast.** flai's 35m.
- **Cost of delay.** 46.88 USD a week, from `flai cod` on E-0025's input.

**T-1839's Work rewrite: yes, apply it.** It describes the dashboard's path as the code has it, and its Done when stands, so it changes no scope.

**Order.** Keep `after` on S-0432, as the epic's order says (TH-0480).

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T17:18:00Z orchestrator
Resolved: Plan approved under plan_backlog_stories, with T-1839's Work rewrite.
