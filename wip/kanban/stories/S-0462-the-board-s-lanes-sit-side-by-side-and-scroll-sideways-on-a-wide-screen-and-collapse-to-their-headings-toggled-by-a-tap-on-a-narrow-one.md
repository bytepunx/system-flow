---
id: S-0462
type: story
nature: improvement
title: The board's lanes sit side by side and scroll sideways on a wide screen, and collapse to their headings, toggled by a tap, on a narrow one
status: ready
parent: E-0027
owner: alex
created: 2026-10-21T07:01:02Z
updated: 2026-10-21T07:58:04Z
transitions:
  - to: ready
    at: 2026-10-21T07:54:32Z
    by: alex
tags: [dashboard]
topics: [dashboard]
touches: [flaiover/src/lib/lanecollapse.svelte.ts, flaiover/src/lib/lanecollapse.svelte.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/+layout.svelte, flaiover/e2e/board-layout.e2e.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0458]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 59
          cache_read: 271604
          cache_write: 1064
          cost: 0.0684
forecast:
  duration: 50m
  delivery: 2026-10-21T10:33:00Z
  basis: "Its own forecast of 50m; 3rd in the pull order with an in-progress limit of 5, behind S-0456, S-0458, S-0463, S-0464, S-0465, S-0459 and S-0461."
  by: flai
  at: 2026-10-21T07:58:04Z
finalized:
  by: alex
  at: 2026-10-21T07:54:08Z
---
# S-0462 The board's lanes sit side by side and scroll sideways on a wide screen, and collapse to their headings, toggled by a tap, on a narrow one

## Goal

The board reads as a board at every width. Today its six lanes stack in one column below 640px and wrap into two columns up to 1024px (`grid-cols-1 sm:grid-cols-2 lg:grid-cols-6`), all inside the layout's `max-w-6xl`. On a wide screen the lanes never stack: they sit in one row as wide as the main area, and the operator scrolls sideways between them. On a phone each lane collapses to its heading, which stays visible, and tapping a heading opens or closes its lane.

## Acceptance criteria

- [ ] From Tailwind's `sm` width the six lanes are one row with a minimum lane width. The row scrolls sideways inside the board when the lanes do not fit, with the page itself never scrolling sideways, and the lanes never wrap or stack at any width from `sm`. The board takes the whole main area beside the rail, not the layout's `max-w-6xl`. Each lane's heading stays visible while the lane's cards scroll.
- [ ] Below `sm` each lane is its heading, with its name, its counts line, and its WIP count against the limit, as a button with `aria-expanded` and `aria-controls`. Tapping it, or pressing Enter or Space on it, opens or closes the lane's cards. Ready, In progress, and Review are open the first time, and the others are closed.
- [ ] The lanes the operator opened or closed are kept per browser in `localStorage` (`$lib/lanecollapse.svelte.ts`), as the board's type toggles are, and are used only below `sm`.
- [ ] The lane menu, the card menu, reordering within a lane, moving a card to another lane through the lane-move dialog, the merge queue in Review, and the publish banner in Done work as before at both widths. A collapsed lane still takes a moved card, and its counts change at once.
- [ ] A Playwright test at 390px and at 1440px checks the collapsed lanes and their toggling, and the one row that scrolls sideways with the document no wider than the viewport. The board's existing tests pass.
- [ ] `design/system/flaiover-dashboard.md` § Views (`/board`) and `docs/users/flaiover.md` § Board describe both layouts.

## Tasks
- T-1969 lanecollapse.svelte.ts keeps which lanes are open on a narrow board per browser, with Ready, In progress, and Review open at first
- T-1970 The board lays its lanes out in one row that scrolls sideways from sm and as headings that toggle their lanes below sm, across the whole main area
- T-1971 A Playwright test shows the board's toggling lanes at 390px and its one sideways-scrolling row at 1440px, with the document no wider than the viewport
- T-1972 The dashboard's design and user docs describe the board's two layouts and its per-browser open lanes

## Notes

- Drafted by the planner for E-0027. After S-0458, which changes the layout's main area that the board opts out of `max-w-6xl` in. It can run beside S-0459 and S-0460.
- Which lanes start open on a phone is the planner's assumption, named on the plan thread: the three where the work moves.

### Planning

- **Touches.** None were declared, so every touch is the planner's. All are files, and no folder touch was kept.
  - **Layout:** `flaiover/src/routes/board/+page.svelte` (860 lines, the lane grid at line 658 and each lane's `<h2>` at line 706) and `board.svelte.test.ts`. `flaiover/src/lib/lanecollapse.svelte.ts` and its test are new, after `boardtypes.svelte.ts`. `flaiover/src/routes/+layout.svelte` lets the board leave `max-w-6xl`. `flaiover/e2e/board-layout.e2e.ts` is new.
  - **Design:** `design/system/flaiover-dashboard.md` § Views (the `/board` row) and `docs/users/flaiover.md` § Board.
  - **Co-change:** `flai touches suggest` lists `BoardCard.svelte` and its test at 5%. The cards keep their shape inside a wider or a collapsed lane, so none was added. `LaneMenu.svelte`, `CardReorder.svelte`, and `LaneMoveDialog.svelte` are only checked, not changed.
- **Forecast.** 50m against flai's 21m. The 860-line board must keep its lane menu, card menu, reordering, and moves working in two layouts, checked by a Playwright test at two widths. The delivery is flai's play-out, moved by its own 29m raise and S-0458's 43m.
- **Cost of delay.** Not yet set. E-0027 has no inputs, so `flai cod` refuses. The input is asked on TH-0553 question 1.
