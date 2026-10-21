---
id: TH-0558
title: "S-0465 plan: four tasks in three layers, T-1979 and T-1980, then T-1981, then T-1982"
anchor:
  path: wip/kanban/stories/S-0465-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md
  item: S-0465
status: resolved
participants: [planner-S-0465, orchestrator]
created: 2026-10-21T07:38:02Z
updated: 2026-10-21T07:38:17Z
---

# TH-0558 S-0465 plan: four tasks in three layers, T-1979 and T-1980, then T-1981, then T-1982

On wip/kanban/stories/S-0465-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md.

## Entries

### 2026-10-21T07:38:02Z planner-S-0465
The plan uses the remedy from I-0145's instance. `flai move <story> done --json` will refuse the way `flai accept --json` does: exit 4 with `{"refused": {...}}`, findings included. `item.move` will map exit 4 to `Refused`. The board and the item page will then list the findings as data, through a component taken out of `Review.svelte`. `linkSegments` will drop its match on flai's wording.

## Tasks and layers

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-1979 | `flai move` to done refuses through `acceptRefusal`; `item.move` gets `exits: {4: Refused}`; tests; `docs/users/flai.md`, `design/system/flai-cli.md` | none |
| 1 | T-1980 | `AcceptFindings.svelte` with its test and examples; `Review.svelte` uses it | none |
| 2 | T-1981 | The board and the item page show the findings from the move's data; `MAIN_CHECKOUT_REFUSAL` goes from `filelinks.ts`; tests; `flaiover-dashboard.md`, `docs/users/flaiover.md` | T-1979, T-1980 |
| 3 | T-1982 | `flai issue close I-0145 --reason` | T-1981 |

T-1979 and T-1980 share no path, so they can run together.

## Figures

- Touches: 22 files, no folders.
- Forecast: 1h. flai gave 36m; I raised it because the change spans two components and S-0378, the review-page version of this fix, took 2h.
- Cost of delay: 22.50 USD a week, from flai's input of 9m.

## Assumptions

- No new ADR is needed. ADR-0145 already says a refused acceptance gives its findings as data, and a move to done is an acceptance.
- `respond` already passes `RepoError` data into the body, so the move route changes only in its comment.
- The board and the item page show findings only. They keep `undoes_main` as text, though T-1979 makes that data available too.

## Proposals

None. I would not split, merge, or drop any task.

### 2026-10-21T07:38:12Z orchestrator
Approved. T-1979 to T-1981 remove the cause, so the dashboard no longer parses flai's wording, and the move and filelinks tests reproduce it. T-1982 closes I-0145. Using ADR-0145's refusal-as-data for the move to done needs no new ADR, as you say. The 22 file touches match the tasks, the 1h forecast is reasoned against S-0378, and the 22.50 USD value comes from flai's input.

Source: design/conventions/strategic-agents.md

### 2026-10-21T07:38:17Z orchestrator
Resolved: The plan for S-0465 is approved.
