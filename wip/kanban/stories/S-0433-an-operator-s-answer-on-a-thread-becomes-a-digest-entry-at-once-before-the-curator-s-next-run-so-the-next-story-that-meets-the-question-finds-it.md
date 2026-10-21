---
id: S-0433
type: story
nature: feature
title: An operator's answer on a thread becomes a digest entry at once, before the curator's next run, so the next story that meets the question finds it
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:01Z
updated: 2026-10-21T13:25:33Z
transitions: []
tags: [cli]
topics: [cli, orchestration, planning]
touches: [design/system/agent-narrative.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/thread.go, flai/cmd/thread_test.go, flai/internal/digest/onanswer.go, flai/internal/digest/onanswer_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
after: [S-0432]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 29m
  delivery: 2026-10-22T07:55:00Z
  basis: "Its own forecast of 29m; 26th in the pull order with an in-progress limit of 5, behind S-0479, S-0483, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429, S-0430, S-0431 and S-0432."
  by: flai
  at: 2026-10-21T13:25:33Z
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

- **Touches.** The story declared none. Every touch comes from its five tasks, file by file.
  - **From layout:** a new `flai/internal/digest/onanswer.go` with its test, called from the three places that save an answer, which the explorer found: `flai/cmd/thread.go`, `threadReply` in `flai/internal/mcpserver/server.go`, and `thread.reply` and `thread.confirm` in `flai/internal/hostapi/writes.go`, each with its test. `threads.ReplyWith` itself is left alone, so that `flai/internal/threads` does not import `digest`, which reads threads.
  - **From design:** `strategic-agents.md`, `agent-narrative.md`, `docs/users/flai.md`, and `docs/users/flaiover.md` for criterion 6.
  - **Co-change, added:** `flai/internal/hostapi/writes_test.go` (74% of the commits that changed `threads.go` or `writes.go`) and `docs/users/flai.md` (13%).
  - **Co-change, not added:** `flaiover/src/lib/server/agent.ts` (41%), since the dashboard need not show the entry; `flai/cmd/serve_actions.go` (27%), since no host action changes; `docs/users/flai-reference.md` (23%), since no flag is added, only a field of the result.
  - **No folder touch.** The entries are written in the main checkout at run time, not on the story's branch.
- **Forecast.** Kept at flai's 29m. It is one function and three call sites. Delivery is flai's play-out moved by S-0429's raise.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
