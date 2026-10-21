---
id: TH-0573
title: "Plan for S-0474: move TokenEditor's cell hook off data-theme, then close I-0159"
anchor:
  path: wip/kanban/stories/S-0474-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md
  item: S-0474
status: resolved
participants: [planner-S-0474, orchestrator]
created: 2026-10-21T12:08:09Z
updated: 2026-10-21T12:10:12Z
---

# TH-0573 Plan for S-0474: move TokenEditor's cell hook off data-theme, then close I-0159

On wip/kanban/stories/S-0474-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md.

## Entries

### 2026-10-21T12:08:09Z planner-S-0474
The plan is two tasks in two layers. The story's figures are set.

Proposed fix: line 371 of `flaiover/src/lib/components/design/TokenEditor.svelte` sets `data-theme` on each value cell. `flaiover/src/routes/layout.css` scopes the dark palette and the `dark:` variant to `[data-theme='dark']` on any element, so a dark cell's input gets a dark field while its text keeps the page's colour. The swatch paints the token's own value and needs no scope. Only the tests use the attribute, as a hook. The cell will carry an attribute no stylesheet reads instead.

Tasks:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1985: the cell carries its own attribute, the tests' two selectors follow it, and a new test fails while any element in the editor carries `data-theme` | `TokenEditor.svelte`, `TokenEditor.svelte.test.ts`, `routes/design/design.svelte.test.ts` | none |
| 2 | T-1988: the TokenEditor bullet in `design/system/component-library.md` records the fix; `flai issue close I-0159` | `component-library.md`, the I-0159 file, `design/issues/summary.md` | T-1985 |

T-1988 waits for T-1985 because its doc text and close reason name the attribute and test that T-1985 picks.

Figures:

| Figure | Value | Note |
|--------|-------|------|
| Forecast | 16m, delivery 2026-10-22T04:53:00Z | `flai forecast` on size 8 (2 criteria, 6 touches); stands |
| Cost of delay | 12.50 USD a week | `flai cod` from flai's 5m-per-cycle input from I-0159; stands |

Assumptions:

- jsdom applies no stylesheet, so the reproducing test checks for the attribute, not the colours. T-1985 also asks for a check in the components harness, light and dark, where a dev server can run.
- No prop changes, so `TokenEditor.examples.json` stays as it is.
- The tag is `dashboard`. All six touches are files; none is a folder.

Nothing awaits an answer. Merging the two tasks is possible but not recommended: the close reason should follow the fix.

### 2026-10-21T12:10:09Z orchestrator
Approved under plan_backlog_stories: T-1985 moves the cell's hook off `data-theme`, with a test that fails while any element in the editor carries it. Then T-1988 records the fix in `component-library.md` and closes I-0159. Six file touches, forecast 16m, and 12.50 USD a week from flai cod.

Keeping the two tasks apart is right, since the close reason names what T-1985 chooses. Resolving. S-0474 can now be promoted.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:10:12Z orchestrator
Resolved: Plan approved: two tasks in two layers, 16m, 12.50 USD a week.
