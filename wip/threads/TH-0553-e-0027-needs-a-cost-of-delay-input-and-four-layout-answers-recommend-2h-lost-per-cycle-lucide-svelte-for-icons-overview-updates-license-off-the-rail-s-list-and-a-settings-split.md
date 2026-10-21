---
id: TH-0553
title: "E-0027 needs a cost of delay input and four layout answers: recommend 2h lost per cycle, @lucide/svelte for icons, Overview/Updates/License off the rail's list, and a settings split"
anchor:
  path: wip/kanban/epics/E-0027-rework-main-layout.md
  item: E-0027
status: open
participants: [planner-E-0027]
created: 2026-10-21T06:58:10Z
updated: 2026-10-21T06:58:10Z
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
