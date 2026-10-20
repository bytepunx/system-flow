---
id: TH-0467
title: "E-0023 needs a cost of delay input and two host answers: recommend penalty 500 USD a week as E-0019, OpenCode installed on the host before the first story, a 5 USD cap on the OpenRouter key"
anchor:
  path: wip/kanban/epics/E-0023-an-opencode-adapter-runs-a-story-s-agent-on-opencode-with-the-guard-the-permission-prompt-the-roles-and-usage-measured-through-its-plugin-hooks-and-its-server-s-events-so-a-non-claude-model-works-a-story-end-to-end.md
  item: E-0023
status: open
participants: [planner-E-0023]
created: 2026-10-16T05:07:49Z
updated: 2026-10-20T04:34:26Z
---

# TH-0467 E-0023 needs a cost of delay input and two host answers: recommend penalty 500 USD a week as E-0019, OpenCode installed on the host before the first story, a 5 USD cap on the OpenRouter key

On wip/kanban/epics/E-0023-an-opencode-adapter-runs-a-story-s-agent-on-opencode-with-the-guard-the-permission-prompt-the-roles-and-usage-measured-through-its-plugin-hooks-and-its-server-s-events-so-a-non-claude-model-works-a-story-end-to-end.md.

## Entries

### 2026-10-16T05:07:49Z planner-E-0023
I am drafting E-0023's stories now. Three inputs are yours. Each has my recommendation first. I draft on the recommended answers and change the stories if you answer otherwise.

## 1. Cost of delay input (blocks the stories' values)

`flai cod E-0023` refuses: the epic has no inputs.

- **Recommended:** `penalty_per_week: 500` on E-0023, the figure you gave E-0019, whose contracts this epic proves on a second harness. The value is then 500 USD a week, shared among the stories by duration.
- Alternatives: `time_lost_per_cycle: 1h` (150 USD a week at `planning.hour_rate` 150), or a figure of your own.

## 2. OpenCode on the host

The first story checks the four open findings (sub-agent sign in a plugin, answering `ask` headless, `step_finish` fields, `--attach`) against a real `opencode` binary. flai never installs it (the epic's Notes).

- **Recommended:** you install OpenCode on this host before the first story starts, and the story records its version as the known version the integration tests pin.
- Alternative: the first story runs on recorded fixtures only, and the open findings wait for a later story with the binary. The guard capability then stays off for longer.

## 3. Spend for the end-to-end story

- **Recommended:** an OpenRouter key with a 5 USD limit on the key, named in the host's `providers` map, and the non-Claude model chosen by you when you finalize that story (its `agent.model`). The plan assumes one story's run fits in that limit.
- Alternative: a higher limit, which you name.

Meanwhile I plan everything that needs no answer. Once question 1 is answered, I write the stories' values.

### 2026-10-20T04:34:26Z planner-E-0023
Reminder: E-0023's plan is complete except question 1. S-0406 to S-0412 have their touches, forecasts, and 38 tasks, all drafts in the backlog. Without a cost of delay input on E-0023, `flai cod` refuses for every story, so they have no value and the orchestrator cannot finalize them. E-0022 is nearly done, so E-0023 would be next.

- **Still recommended:** `flai edit E-0023 --penalty-per-week 500`, as on E-0019.

Once the input is on the epic, I write each story's value and its Planning note. Questions 2 (OpenCode installed before S-0406) and 3 (a 5 USD OpenRouter key limit) can wait until S-0406 and S-0412 are finalized.
