---
id: TH-0872
title: "S-0587 cost of delay input: recommend time_lost_per_cycle 10m"
anchor:
  path: wip/kanban/stories/S-0587-default-harness-and-agent-should-work-out-of-the-box.md
  item: S-0587
status: resolved
participants: [planner-S-0587, orchestrator, alex]
created: 2026-11-07T10:05:08Z
updated: 2026-11-07T15:15:56Z
---

# TH-0872 S-0587 cost of delay input: recommend time_lost_per_cycle 10m

On wip/kanban/stories/S-0587-default-harness-and-agent-should-work-out-of-the-box.md.

## Entries

### 2026-11-07T10:05:08Z planner-S-0587
S-0587 has no cost of delay inputs and no epic, so `flai cod S-0587` cannot work out a value. I am planning the rest meanwhile.

**Recommendation:** set `time_lost_per_cycle: 10m`. That matches S-0589 and S-0590, the two other remediation stories on agent start-up. It also matches what the bug costs: the operator spends a few minutes per refused start on a hand fix, and that only happens in new projects.

Alternatives:

| Input | Value | When it fits |
|-------|-------|--------------|
| `time_lost_per_cycle` | 30m | Every new project hits this before its first story runs, and the fix means editing several stories or the host. |
| `penalty_per_week` | 250 | You count it as onboarding friction that loses would-be users, not as your own time. |

Which input and figure should S-0587 carry?

### 2026-11-07T10:05:16Z planner-S-0587
Correction to my recommendation. S-0573 is the earlier fix for this same failure ("Stories in new projects ignore default agent/harness settings"). You gave it `penalty_per_week: 50`. S-0587 is the same failure, reached by a path S-0573 did not cover, so I now recommend the same input: **`penalty_per_week: 50`**. Second choice: `time_lost_per_cycle: 10m`, as on S-0589 and S-0590.

### 2026-11-07T10:06:28Z orchestrator (recommendation)
**Recommended:** set `penalty_per_week: 50` on S-0587, as the planner corrected its recommendation to:

```bash
flai edit S-0587 --penalty-per-week 50
```

**Why:**
- S-0573, "Stories in new projects ignore default agent/harness settings", fixed the same failure by another path. You gave it `penalty_per_week: 50` on 2026-11-03, which I checked on S-0573. S-0587 is the same harm, so the same input keeps the two consistent.
- The planner then works out the value on its next run.

**Why you:** S-0587 is in ready, not the backlog, and has no epic and no inputs. So it is not a backlog story's first input, and setting it is not mine under `plan_backlog_stories`. Once you confirm this, or answer otherwise, I set what you decide under `set_cost_of_delay`.

**Alternative:** `time_lost_per_cycle: 10m`, as on S-0589 and S-0590, if you count it as your own time lost rather than as onboarding friction.

**What changes meanwhile:** nothing. The policy is `throughput`, so the ready column is ordered by forecast and this value does not move S-0587. It is held now for declaring no touches, and its planner is still running.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T15:15:56Z alex
Resolved: S-0587 was accepted
