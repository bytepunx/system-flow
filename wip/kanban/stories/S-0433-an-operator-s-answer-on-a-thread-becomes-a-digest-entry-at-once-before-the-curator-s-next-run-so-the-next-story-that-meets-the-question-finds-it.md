---
id: S-0433
type: story
nature: feature
title: An operator's answer on a thread becomes a digest entry at once, before the curator's next run, so the next story that meets the question finds it
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:01Z
updated: 2026-10-16T06:07:29Z
transitions: []
tags: [cli]
topics: [cli, orchestration, planning]
after: [S-0432]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

## Notes

- Drafted by the planner for E-0025, after S-0432 in the epic's order. It needs only S-0429's rule, so it could follow S-0429; the plan thread proposes that.
- Writes happen in the main checkout, where threads are written, and are not committed, as a thread entry is not.
