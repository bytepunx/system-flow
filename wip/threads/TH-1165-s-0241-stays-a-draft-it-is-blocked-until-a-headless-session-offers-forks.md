---
id: TH-1165
title: "S-0241 stays a draft: it is blocked until a headless session offers forks"
anchor:
  path: wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md
  item: S-0241
status: open
participants: [orchestrator]
created: 2026-11-13T15:45:14Z
updated: 2026-11-13T15:45:14Z
---

# TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks

On wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md.

## Entries

### 2026-11-13T15:45:14Z orchestrator
**Recommended:** leave S-0241 a draft and blocked until a headless Claude Code session offers `subagent_type: fork`. Then unblock it, and I finalize it under `finalize_drafts`.

**Why I did not finalize it:**

- `flai promote --drafts` lists it complete: criteria, touches, a 2h forecast, and a value of 25 USD a week.
- But the story is blocked, and its Notes say: "Do not move it to ready until a headless session offers forks: it cannot be worked before then."
- You moved it back to backlog on 2026-11-01. Finalizing it now would only add a blocked story to the pull order.

**What would change this:** a Claude Code release whose headless `claude -p` session offers forks. Then unblock the story with `flai unblock S-0241`.

**Board starvation (TH-1164):** S-0241 is the only complete draft. Every E-0034 draft lacks a cost of delay value until you decide TH-1097. E-0036's ready stories wait on S-0751, in progress.
