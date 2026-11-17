---
id: TH-1271
title: "Plan for S-0819: the guard forgets a foreground sub-agent once its result is in the session transcript"
anchor:
  path: wip/kanban/stories/S-0819-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md
  item: S-0819
status: resolved
participants: [planner-S-0819, orchestrator, alex, agent-S-0819]
created: 2026-11-17T03:35:53Z
updated: 2026-11-17T07:40:46Z
---

# TH-1271 Plan for S-0819: the guard forgets a foreground sub-agent once its result is in the session transcript

On wip/kanban/stories/S-0819-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md.

## Entries

### 2026-11-17T03:35:53Z planner-S-0819
S-0819 is planned as five draft tasks in four layers. Touches, forecast, and cost of delay are written, with their reasons under the story's `### Planning`.

## Proposed fix

A foreground sub-agent blocks the story's agent. So once the session transcript holds the result of the Agent call that launched it, that sub-agent has ended, error or not. At the story agent's `wait_for_events`, `flai guard` reads the transcript the hook input names (`transcript_path`). Under the record's lock, it drops each sub-agent whose foreground launch has its result there, then decides on what is left. If the transcript cannot be read, the record stays as it is. No hook configuration changes, since the `PreToolUse` matcher already covers `mcp__flai__.*`.

## Tasks and layers

| Layer | Task | Waits for | What |
|-------|------|-----------|------|
| 1 | T-3504 (research) | none | Confirm the signal from S-0738's and S-0739's transcripts, then write a proposed ADR refining ADR-0092 |
| 2 | T-3505 | T-3504 | `flai/internal/guard`: the transcript reader (new `transcript.go`), the record's prune, `transcript_path` on the call, and unit tests |
| 2 | T-3507 | T-3504 | `design/system/agent-context.md`, `design/system/flai-cli.md`, `docs/users/flai.md` |
| 3 | T-3506 | T-3505 | `flai/cmd/guard.go` `holdWaits` prunes before deciding, plus the end-to-end test that reproduces I-0376 |
| 4 | T-3508 | T-3506, T-3507 | `flai issue close I-0376 --reason` |

## Figures

- Forecast: 1h30m. flai gave 34m with the touches; I raised it for T-3504's measurement and ADR.
- Cost of delay value: 25.00 USD a week, as `flai cod` gives from the inputs.

## Assumptions

- Claude Code writes the foreground Agent call's result into the main transcript even when the sub-agent ends on an API error. I-0376 says the agent got a partial result back. T-3504 confirms this before any code. If it does not hold, T-3504's ADR names the surer signal, and T-3505 builds that one.
- Only the foreground case is fixed. A background sub-agent cut off with no `SubagentStop` stays in the record, as ADR-0092 has it, because the prompt already forbids waiting on one.
- The story keeps the folder touch `design/adrs/`, because the ADR's file name is unknown until T-3504.
- T-3507 also corrects the `flai guard` row of `design/system/flai-cli.md`, which says the record's lock is "taken over after 5 s". The code never takes it over. It is the row the task edits anyway.

## Questions

1. I-0374 has the same cause, on S-0738. I recommend T-3508 also close I-0374 as a duplicate of I-0376 with the same reason. If you agree, I will add it to T-3508's Work and touches.
2. The cost of delay inputs count 1 occurrence (10m a cycle), but with I-0374 there are 2. I recommend raising `time_lost_per_cycle` to 20m, which makes the value 50.00 USD a week. The inputs are yours, so I have not changed them.

### 2026-11-17T03:36:56Z orchestrator (recommendation)
**Recommended:** say yes to both questions. Set `time_lost_per_cycle: 20m` on S-0819 (50 USD a week), and let T-3508 also close I-0374 as a duplicate of I-0376:

```bash
flai edit S-0819 --time-lost-per-cycle 20m
```

**The plan:** I approve it. Five tasks in four layers cover both criteria. T-3504 confirms the transcript signal before any code is written. The 15 touches name the guard's files, the design, and the issue files. The 1h30m forecast and the 25 USD a week value match the current input of 10m.

**Question 1, I-0374:** it is open with count 1, and it records the same cause on S-0738. Its title is I-0376's, almost word for word. Closing it with the same reason leaves no open issue for a cause the story has fixed. The planner adds it to T-3508's Work and touches on its next run.

**Question 2, the input:**

- **Why 20m:** flai set 10m from one occurrence. I-0374 is a second occurrence in the same cycle.
- **Why you:** S-0819 already has an input, the 10m flai set. I give only a backlog story's first input, so changing it is your call.
- **The alternative:** keep 10m (25 USD a week) if you count I-0374 and I-0376 as one incident.

**S-0819 itself:** `flai promote --drafts` lists it complete. I finalize it now, so it can go to ready. Your answer changes only its value and T-3508, not its criteria.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T07:38:53Z alex
Confirmed the recommendation of 2026-11-17T03:36:56Z orchestrator.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T07:39:48Z orchestrator
Set: `time_lost_per_cycle: 20m` on S-0819, as you confirmed. I gave the input only; its value is left for flai to work out again.

**For agent-S-0819:** S-0819 is in progress, so no planner runs for it again. Question 1 is yours to carry out: add I-0374 to T-3508's Work and touches, and close it as a duplicate of I-0376 with the same `--reason`, in the same task. Resolve this thread once T-3508 is done.

Source: design/digest/D-0659-plan-for-s-0819-the-guard-forgets-a-foreground-sub-agent-once-its-result-is-in.md

Source: design/digest/D-0659-plan-for-s-0819-the-guard-forgets-a-foreground-sub-agent-once-its-result-is-in.md

### 2026-11-17T07:40:46Z agent-S-0819
Resolved: T-3508 is done: I-0376 and its duplicate I-0374 are closed on story/S-0819 with reasons naming ADR-0315's prune (commits 2108fd7, 736a9ca); the 20m input is set.
