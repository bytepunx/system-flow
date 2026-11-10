---
id: TH-1016
title: "Plan for S-0700: the new item form scrolls to and focuses its refusal, three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0700-the-new-item-form-s-lint-refusal-is-easy-to-miss-so-a-refused-epic-looks-like-a.md
  item: S-0700
status: resolved
participants: [planner-S-0700, orchestrator]
created: 2026-11-10T10:28:40Z
updated: 2026-11-10T11:01:42Z
---

# TH-1016 Plan for S-0700: the new item form scrolls to and focuses its refusal, three tasks in three layers

On wip/kanban/stories/S-0700-the-new-item-form-s-lint-refusal-is-easy-to-miss-so-a-refused-epic-looks-like-a.md.

## Entries

### 2026-11-10T10:28:40Z planner-S-0700
The plan for S-0700 is three tasks in three layers. Its touches, forecast, and cost of delay are written, with the reasons under `### Planning` in the story's Notes.

## Proposed fix

`flaiover/src/lib/components/NewItemForm.svelte` draws the refusal alert below the body and preview, which are each `h-[55vh]`, and does nothing else on a 422. After any refusal, the form scrolls the alert into view and focuses it. The Create button says `Refused: see above` until the next submit.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2966 Scroll to and focus the refusal, mark the Create button, and add the test that reproduces I-0263 | `NewItemForm.svelte`, `NewItemForm.svelte.test.ts` | none |
| 2 | T-2968 Say so in `design/system/flaiover-dashboard.md` (`/new` row) and `docs/users/flaiover.md` | the two documents | T-2966 |
| 3 | T-2971 Close I-0263 with `flai issue close` | the issue file, `design/issues/summary.md` | T-2966, T-2968 |

## Figures

- Forecast: 45m, raised from flai's 18m. The flaiover close-out tier took 20m44s on S-0669 today.
- Cost of delay: 50 USD a week, as `flai cod` gives it from the 20m input.

## Assumptions

- The fix is in the browser only. flai's 422 and its findings stay as they are, and the component gets no new prop, so `NewItemForm.examples.json` is untouched.
- The missed refusal in I-0263 happened because the form was submitted from above the fold, for example with Enter in the title field. In that case neither the button's `Creating…` nor the alert can be seen.
- I was started to plan this story again for a thread named `asked`. No such thread exists on S-0700, and the story had no earlier plan, so this is a first plan.

If you would rather move the alert above the body, or show it at the top of the form, say so and I will change T-2966.

### 2026-11-10T10:57:41Z orchestrator
Approved as written. T-2966, T-2968, and T-2971 meet both criteria in order: the fix and its reproducing test, then the docs, then the close. The touches are the six files the plan names. The 45m forecast and the 50 USD a week value state their reasons. The alert stays where it is, and the form scrolls to it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:01:42Z orchestrator
Resolved: Plan approved
