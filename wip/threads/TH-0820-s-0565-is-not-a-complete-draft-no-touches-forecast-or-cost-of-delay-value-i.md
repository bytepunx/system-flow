---
id: TH-0820
title: "S-0565 is not a complete draft: no touches, forecast, or cost of delay value; I started its planner"
anchor:
  path: wip/kanban/stories/S-0565-one-story-of-this-repository-is-worked-end-to-end-by-flai-serve-on-harness.md
  item: S-0565
status: resolved
participants: [orchestrator]
created: 2026-11-03T01:19:43Z
updated: 2026-11-03T01:30:16Z
---

# TH-0820 S-0565 is not a complete draft: no touches, forecast, or cost of delay value; I started its planner

On wip/kanban/stories/S-0565-one-story-of-this-repository-is-worked-end-to-end-by-flai-serve-on-harness.md.

## Entries

### 2026-11-03T01:19:43Z orchestrator
**Recommended:** nothing for you yet. Let planner-S-0565 fill what is missing, then finalize under finalize_drafts once `flai promote --drafts` lists it complete.

**What `flai promote --drafts` says S-0565 lacks:**

- touches
- a forecast duration and delivery
- a cost of delay value

Its goal, criteria, and `after` (S-0508, S-0510, S-0517) are consistent with TH-0812 and TH-0814 (D-0466, D-0469).

**What I did:** `flai plan --candidates` lists S-0565, so I started its planner under plan_backlog_stories at 2026-11-03T01:19:38Z.

**One risk:** flai serve started the planner with the command `rmhclod`. TH-0819 says that binary failed flai's permission check with a shell syntax error. If the planner fails, that is why.

### 2026-11-03T01:30:16Z orchestrator
Resolved: planner-S-0565 filled what was missing: flai promote --drafts now lists S-0565 complete (touches, forecast 1h30m, value 52.12). Finalizing waits on the operator's answer to TH-0821 (proof stories, close-out option a or b).
