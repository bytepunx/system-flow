---
id: D-0236
title: "Plan for S-0309: bound permission_prompt's wait, keep the thread open, take the answer on the retry: Recommendation: take the planner's remedy, a bounded 4-minute wait, the thread kept open, and the…"
kind: answer
sources: [TH-0335]
topics: [cli, agents, flai, go, code, all]
paths: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/delegation.md, docs/users/flai.md, design/system/flai-cli.md, design/adrs, design/adrs/README.md, design/issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md, design/issues/summary.md, design/issues/I-0108-flai-task-done-commits-every-uncommitted-file-in-the-story-worktree-so-closing-one-task-of-a-layer-sweeps-its-siblings-work-into-its-commit.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md#decision"
asked_again: []
cited: []
---

# D-0236 Plan for S-0309: bound permission_prompt's wait, keep the thread open, take the answer on the retry: Recommendation: take the planner's remedy, a bounded 4-minute wait, the thread kept open, and the…

The operator, alex, answered TH-0335 at 2026-10-08T04:22:58Z.

What was asked:

```text
S-0309 now has five tasks in three layers, 11 touches, a 35m forecast, and a cost of delay of 75 USD a week.
```

The answer:

```text
Recommendation: take the planner's remedy, a bounded 4-minute wait, the thread kept open, and the answer taken on the retry. Approve the plan as it stands.

- **It removes I-0103's cause.** The call ends before Claude Code's idle timeout on both transports, so the write never fails on the timeout. Your later answer is still used.
- **It keeps who may answer.** The owners ADR-0097 names, and the first words that allow, are unchanged. Only how long the call waits changes, and what an unanswered thread becomes.
- **It leaves the operator's answer whole.** Today an unanswered thread is settled as refused, so an answer that comes late is lost. With the remedy, it is applied on the agent's retry. `endWhy` already ends the agent on its own open question, and flai serve restarts it on your answer.
- **The alternative depends on Claude Code.** Progress notifications help only if Claude Code sends a progress token, and they keep the agent blocked while you are away.

This refines ADR-0086 and ADR-0097, so it is yours to confirm. T-1274 records it.

The plan fits the story:

- T-1275 removes the cause, with a test that reproduces I-0103 (criterion 1).
- T-1278 closes I-0103 (criterion 2).
- The three tasks of layer 2 share no path.
- The 35m forecast rests on S-0299 and S-0257. The 75 USD a week value comes from flai's input from the issue.

S-0309 stays a draft until this thread settles.

Source: design/adrs/0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md § Decision
```

Source: TH-0335.
