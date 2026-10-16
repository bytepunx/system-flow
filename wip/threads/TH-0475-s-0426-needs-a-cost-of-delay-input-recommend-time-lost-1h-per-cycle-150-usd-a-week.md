---
id: TH-0475
title: "S-0426 needs a cost of delay input: recommend time lost 1h per cycle (150 USD a week)"
anchor:
  path: wip/kanban/stories/S-0426-agent-memory-system.md
  item: S-0426
status: open
participants: [planner-S-0426]
created: 2026-10-16T05:55:55Z
updated: 2026-10-16T05:55:55Z
---

# TH-0475 S-0426 needs a cost of delay input: recommend time lost 1h per cycle (150 USD a week)

On wip/kanban/stories/S-0426-agent-memory-system.md.

## Entries

### 2026-10-16T05:55:55Z planner-S-0426
S-0426 has no cost of delay inputs and no epic to take a share from, so `flai cod S-0426` refuses. The inputs are yours.

Recommendation: `time_lost_per_cycle: 1h`. At `planning.hour_rate` 150 and a 168h cycle, that is 150 USD a week.

Why 1h: every story's agent, sub-agent, and strategic agent primes and searches the corpus through `flai prime` and `doc_search` today. An hour a week of agent time lost to slow or missed recall across the board is a modest guess. The research itself saves nothing until a follow-on story builds what it recommends.

Alternatives:

| Input | Value a week | When it fits |
|-------|--------------|--------------|
| `time_lost_per_cycle: 10m` | 25 USD | As S-0413, the other research story you filed tonight |
| `time_lost_per_cycle: 4h` | 600 USD | If you see priming and search overhead as a main cost of every run |

Set it with:

```bash
flai edit S-0426 --time-lost-per-cycle 1h
```

I plan the touches, forecast, and tasks meanwhile, and write the value once the input is set.
