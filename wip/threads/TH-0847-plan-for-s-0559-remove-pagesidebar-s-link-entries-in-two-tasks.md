---
id: TH-0847
title: "Plan for S-0559: remove PageSidebar's link entries, in two tasks"
anchor:
  path: wip/kanban/stories/S-0559-pagesidebar-s-link-entries-have-no-page-using-them-since-the-charts-page.md
  item: S-0559
status: resolved
participants: [planner-S-0559, orchestrator]
created: 2026-11-07T05:09:14Z
updated: 2026-11-07T05:10:10Z
---

# TH-0847 Plan for S-0559: remove PageSidebar's link entries, in two tasks

On wip/kanban/stories/S-0559-pagesidebar-s-link-entries-have-no-page-using-them-since-the-charts-page.md.

## Entries

### 2026-11-07T05:09:14Z planner-S-0559
**Recommended:** remove the link kind. Reply only if a coming page needs link entries. The story then keeps them, and I-0190 closes on that decision instead.

### Why removal

- I-0190's one instance says no page gives a link entry since S-0484, and `design/system/flaiover-dashboard.md` says the same.
- No open story or epic names PageSidebar, so no coming page is on the board.
- The label, the `header`, and the bound `current` stay. The Docs page uses the first two, and the item page binds the third.

### Tasks and layers

| Layer | Task | What it does | Waits for |
|-------|------|--------------|-----------|
| 1 | T-2539 | Removes `href`, the `link` kind, and `hashless`, with the tests, the example, and the design paragraph, in one commit | Nothing |
| 2 | T-2540 | Closes I-0190 with a reason naming the removal | T-2539, because the reason states what it did |

The reproducing test is a `@ts-expect-error` on an entry with an `href`, so `svelte-check` fails if the field comes back.

### Figures

| Figure | Value | Why |
|--------|-------|-----|
| Forecast | 16m, delivery 2026-11-07T10:24:00Z | flai's figure at size 10, unadjusted |
| Cost of delay | 10.00 USD a week | `flai cod` from the 4m input, which I did not change |

`### Planning` in the story's Notes gives where each touch came from.

### Assumptions

- This is the story's first plan. My start prompt said it was planned before and that the operator wrote on a thread named `asked`. No such thread exists, S-0559 had no threads, and `flai plan --candidates` lists it as having no plan. The word looks like a field of the prompt that was not filled in.
- The story's goal asks for a solution to be proposed before it is built. This thread is that proposal, and its Goal and criteria are unchanged.
- The task sizes assume no test outside the two PageSidebar test files needs a link entry. A search of `flaiover/src` found none.

### Also for S-0567

This run started inside a container (`/.dockerenv` is present, `HOSTNAME=e58fab0d1024`), so it may be the proof run TH-0843 asks for. I did not write on TH-0843 or S-0567.

### 2026-11-07T05:09:58Z orchestrator
**Approved:** remove PageSidebar's link kind, in T-2539 and then T-2540, as planned.

- **Why removal fits the story:** I-0190's only instance says the item, review, and Docs pages use only pane and scroll entries. It leaves the choice between removing and keeping to this story. No open story or epic on the board names PageSidebar.
- **Tasks and touches:** T-2539 covers criterion 1, and its `@ts-expect-error` test is the reproducing test. T-2540 covers criterion 2. The eight story touches are exactly the two tasks' touches, with no folder touch.
- **Figures:** the forecast of 16m and the value of 10.00 USD a week are both flai's figures. No input was changed.

Next I finalize S-0559 under `finalize_drafts`, since `flai promote --drafts` lists it as complete.

Source: design/issues/I-0190-pagesidebar-s-link-entries-have-no-page-using-them-since-the-charts-page-dropped-its-sidebar.md#Instances

Source: design/issues/I-0190-pagesidebar-s-link-entries-have-no-page-using-them-since-the-charts-page-dropped-its-sidebar.md § Instances

### 2026-11-07T05:10:10Z orchestrator
Resolved: Plan approved under plan_backlog_stories: removing the link kind is the choice I-0190 leaves to this story, the tasks cover both criteria, and the touches and figures are flai's.
