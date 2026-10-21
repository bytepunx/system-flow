---
id: TH-0553
title: "E-0027 needs a cost of delay input and four layout answers: recommend 2h lost per cycle, @lucide/svelte for icons, Overview/Updates/License off the rail's list, and a settings split"
anchor:
  path: wip/kanban/epics/E-0027-rework-main-layout.md
  item: E-0027
status: resolved
participants: [planner-E-0027, orchestrator, alex]
created: 2026-10-21T06:58:10Z
updated: 2026-10-21T13:03:19Z
---

# TH-0553 E-0027 needs a cost of delay input and four layout answers: recommend 2h lost per cycle, @lucide/svelte for icons, Overview/Updates/License off the rail's list, and a settings split

On wip/kanban/epics/E-0027-rework-main-layout.md.

## Entries

### 2026-10-21T06:58:10Z planner-E-0027
I am drafting E-0027's stories now. Five of its inputs are yours. Each has my recommendation first. The stories go in as drafts either way, and I set the cost of delay values once question 1 is answered.

1. **Cost of delay input.** E-0027 has no inputs, so `flai cod` refuses its stories. I recommend `time_lost_per_cycle: 2h`, the operator time lost each week to navigating the top menu and scrolling the item page. At `planning.hour_rate` 150 over a 168h cycle that is 300 USD a week. The alternatives are 8h, as on E-0020, the design-system epic, which is 1200 USD a week, or `penalty_per_week: 500`, as on E-0019 and E-0023.
2. **Icons.** The rail needs one icon per group and per page, and `flaiover/package.json` has no icon library. I recommend adding `@lucide/svelte` (ISC licence, tree-shaken, one Svelte component per icon), with a `design/tech` entry. The alternative is a small set of hand-written inline SVG components under `flaiover/src/lib/icons/`, with no dependency.
3. **Pages the epic's list leaves out.** I recommend these three placements:
   - Overview (`/`) is reached by the flaiover name in the top bar.
   - Updates (`/host`) is reached by the host-flai status pill, which already links there.
   - License (`/license`) goes in a small link at the foot of the rail.

   Items (`/items/:id`), Review, New, and Edit keep their current group: Flow for the first three, Documentation for Edit.
4. **Configuration split.** The epic's five Configuration pages come out of today's one `/settings` page (`SettingsPanel.svelte`, ten sections). I recommend this split:
   - **Features**: Host actions, Checks.
   - **Agents**: Default agent, Providers, Planning again.
   - **File System**: Shared paths, Import folders.
   - **Projects**: Projects, Tokens.
   - **Harnesses**: Starting agents, one block per harness.

   `/settings` itself redirects to Features.
5. **Narrow widths.** At phone width a left rail takes too much room. I recommend that below Tailwind's `md` (768px) the rail becomes a drawer, opened from a menu button at the start of the top bar and closed when a page is chosen. At `md` and wider it stays docked.

Answer any you want differently. Silence on 2 to 5 means I keep my recommendations as written into the stories. Question 1 has to be answered before the values can be set.

### 2026-10-21T07:55:47Z planner-E-0027
S-0458 to S-0462 are finalized and in ready, and S-0458 is in progress. I take questions 2 to 5 as answered by my recommendations, as the stories are written.

Question 1 is still open. Without it, none of the five stories has a cost of delay value, so a `cod` or `wsjf` order ranks them last. Their pull order is now S-0459, S-0461, S-0462, S-0460, S-0458. Every one of them waits on S-0458, which is already in progress.

Set `time_lost_per_cycle: 2h` on E-0027 (300 USD a week), or an alternative, and I will work out and write each story's value.

### 2026-10-21T11:04:44Z orchestrator (recommendation)
Recommended: set `time_lost_per_cycle: 2h` on E-0027, which is 300 USD a week.

- Three of the five stories are accepted and released: S-0458, S-0459 and S-0462, in flaiover 0.53.0 and 0.53.1. S-0460 is in progress and S-0461 is ready.
- So the value now orders only S-0461, which is alone in ready. It also records what the epic was worth for `flai stats`. The planner's 2h figure, operator time spent navigating, fits that purpose.
- E-0020's 8h, or a 500 USD penalty, would overstate a layout rework that is mostly delivered.

Questions 2 to 5 are settled: the stories were built on the planner's recommendations, and you finalized them. The epic's input is yours to set, with `flai edit E-0027` and its cost of delay flags.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-21T12:46:52Z alex
Confirmed the recommendation of 2026-10-21T11:04:44Z orchestrator.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-21T13:03:19Z alex
Resolved.
