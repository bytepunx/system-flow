---
id: D-0659
title: "Plan for S-0819: the guard forgets a foreground sub-agent once its result is in the session transcript: Recommended: say yes to both questions"
kind: answer
sources: [TH-1271]
topics: [guard, delegation, flai, cli, go, code, all]
paths: [flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/transcript.go, flai/internal/guard/transcript_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, design/adrs, design/adrs/README.md, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0376-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md, design/issues/summary.md, flai/internal/guard/call.go]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0659 Plan for S-0819: the guard forgets a foreground sub-agent once its result is in the session transcript: Recommended: say yes to both questions

The operator, alex, answered TH-1271 at 2026-11-17T07:38:53Z.

What was asked:

```text
S-0819 is planned as five draft tasks in four layers. Touches, forecast, and cost of delay are written, with their reasons under the story's `### Planning`.
```

The answer:

````text
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
````

Source: TH-1271.
