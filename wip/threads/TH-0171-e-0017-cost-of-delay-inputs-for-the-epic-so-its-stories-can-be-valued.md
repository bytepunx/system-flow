---
id: TH-0171
title: "E-0017: cost of delay inputs for the epic, so its stories can be valued"
anchor:
  path: wip/kanban/epics/E-0017-story-loop-work-whose-outcome-the-repository-determines-moves-from-the-agent-into-flai-reached-from-the-cli-the-host-channel-and-mcp.md
  item: E-0017
status: resolved
participants: [planner-E-0017, alex]
created: 2026-10-06T11:28:17Z
updated: 2026-10-06T11:36:46Z
---

# TH-0171 E-0017: cost of delay inputs for the epic, so its stories can be valued

On wip/kanban/epics/E-0017-story-loop-work-whose-outcome-the-repository-determines-moves-from-the-agent-into-flai-reached-from-the-cli-the-host-channel-and-mcp.md.

## Entries

### 2026-10-06T11:28:17Z planner-E-0017
Recommended: set `time_lost_per_cycle: 6h` on E-0017 (`flai edit E-0017 --time-lost-per-cycle 6h`). That gives 6h × 150 USD = 900 USD a week. Each of the seven stories, none of which has inputs, then gets a share of that value by its forecast duration among the epic's open stories, from `flai cod`.

Why 6h: the epic's evidence is 108 story runs and about 43 hours of agent time. The run logs under `~/.flai/serve/agents/` span 2026-09-23 to 2026-10-04, about 1.7 weeks, so that is about 25 hours a week. The epic counts the following as mechanical work it would remove:

- about 2,900 ceremony turns
- 1,449 hand test or lint runs whose logs the model reads
- 88 verifier minutes
- 650 empty wakes, 1,235 minutes waiting

That is about a quarter of the agent time, 6h a week.

Alternatives:

- 4h (600 USD a week) if you count only turn time and not the waits.
- A `penalty_per_week` if you would rather price the dollars (about $70 of pure-ceremony turns over the period, so about 40 USD a week, which undervalues the larger mechanical share).

flai cod refuses to value any story until the epic or the story has an input, and the inputs are yours. Meanwhile I am planning touches, forecasts, and the plan thread.

### 2026-10-06T11:32:49Z alex
6h

### 2026-10-06T11:34:49Z planner-E-0017
Recorded: `time_lost_per_cycle: 6h` on E-0017, which comes to 900.00 USD a week.

flai cod splits that among the seven stories by forecast duration. I set each story's value by the turns the epic's evidence says it removes instead, keeping the 900 total:

| Story | USD a week | Turns it removes |
|-------|-----------:|------------------|
| S-0269 | 435 | about 2,900 |
| S-0273 | 218 | 1,449 |
| S-0272 | 98 | 650 |
| S-0274 | 65 | about 432 |
| S-0270 | 44 | 88 verifier minutes, about 290 turns |
| S-0275 | 23 | about 150, my estimate |
| S-0271 | 17 | 111 |

Each story's `### Planning` note gives the reason. The plan summary follows in its own thread.

### 2026-10-06T11:36:46Z alex
Resolved.
