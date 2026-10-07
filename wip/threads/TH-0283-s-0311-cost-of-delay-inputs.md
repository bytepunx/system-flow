---
id: TH-0283
title: S-0311 cost of delay inputs
anchor:
  path: wip/kanban/stories/S-0311-the-orchestrator-needs-a-role-it-can-run-tests-in-review-under.md
  item: S-0311
status: open
participants: [planner-S-0311]
created: 2026-10-07T14:25:10Z
updated: 2026-10-07T14:25:10Z
---

# TH-0283 S-0311 cost of delay inputs

On wip/kanban/stories/S-0311-the-orchestrator-needs-a-role-it-can-run-tests-in-review-under.md.

## Entries

### 2026-10-07T14:25:10Z planner-S-0311
S-0311 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0311` cannot work out a value. The inputs are yours.

Recommendation: time_lost_per_cycle 30m, no revenue_per_week and no penalty_per_week. Each story the orchestrator verifies while FLAI_ROLE=orchestrate leaks into `flai verify` fails go-test (TH-0260 on S-0213), so the story is left in review and you verify it from your shell. That costs about one manual verify and a thread per affected acceptance.

Set it with `flai edit S-0311 --time-lost-per-cycle 30m`, or answer with other figures, and I will write the value from `flai cod`.
