---
id: TH-0335
title: "Plan for S-0309: bound permission_prompt's wait, keep the thread open, take the answer on the retry"
anchor:
  path: wip/kanban/stories/S-0309-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md
  item: S-0309
status: open
participants: [planner-S-0309, orchestrator]
created: 2026-10-07T23:28:14Z
updated: 2026-10-07T23:28:25Z
---

# TH-0335 Plan for S-0309: bound permission_prompt's wait, keep the thread open, take the answer on the retry

On wip/kanban/stories/S-0309-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md.

## Entries

### 2026-10-07T23:28:14Z planner-S-0309
S-0309 now has five tasks in three layers, 11 touches, a 35m forecast, and a cost of delay of 75 USD a week.

## Proposed remedy for I-0103 (recommended)

1. `permission_prompt` holds a call for at most 4 minutes. That is below Claude Code's idle timeout on both transports: 5 minutes over HTTP and 30 over stdio.
2. When the 4 minutes pass, or the session ends, with no answer, it refuses the write. The refusal names the thread, and the thread stays open. Today the thread is settled as refused.
3. A later request for the same story, tool, path, and input finds that thread. An allow already on it lets the write through at once. A refusal refuses it. With no answer yet, it waits again, with the same limit.
4. The open thread counts as the agent's own question, through `endWhy` in `flai/internal/mcpserver/server.go`. So `wait_for_events` already ends the agent once nothing else is left. flai serve starts it again on your answer, and it makes the same write again.

Alternative: send progress notifications while waiting, with `keepAlive` in `flai/internal/mcpserver/timing.go`. This keeps the agent blocked for as long as you are away, and it helps only if Claude Code sends a progress token.

## Tasks and layers

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-1274 Write an ADR refining ADR-0086 and ADR-0097 | none | `design/adrs/`, `design/adrs/README.md` |
| 2 | T-1275 Add the wait limit, the open thread, and the retry, with a test that reproduces I-0103 | T-1274 | `flai/internal/mcpserver/permission.go`, `permission_test.go` |
| 2 | T-1276 Change the start prompt and `delegation.md`: retry the write after the answer instead of staging files with `cp` commands | T-1274 | `flai/internal/harness/harness.go`, `harness_test.go`, `design/conventions/delegation.md` |
| 2 | T-1277 Update the docs | T-1274 | `docs/users/flai.md`, `design/system/flai-cli.md` |
| 3 | T-1278 Close I-0103 | T-1275, T-1276, T-1277 | the I-0103 file, `design/issues/summary.md` |

The three tasks in layer 2 share no paths, so they can run together.

## Assumptions

- The 4-minute limit is a constant in the code, not a manifest setting.
- The template's `delegation.md` lacks the workaround text, so it is not changed.
- `design/adrs/` is the only folder touch, because the ADR's file name is not known until it is written.
- If you want a different remedy, T-1274 records it, and the tasks in layer 2 change with it.

### 2026-10-07T23:28:25Z orchestrator (recommendation)
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
