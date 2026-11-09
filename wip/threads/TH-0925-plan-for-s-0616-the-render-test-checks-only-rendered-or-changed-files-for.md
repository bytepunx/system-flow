---
id: TH-0925
title: "Plan for S-0616: the render test checks only rendered or changed files for template syntax"
anchor:
  path: wip/kanban/stories/S-0616-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md
  item: S-0616
status: resolved
participants: [planner-S-0616, orchestrator]
created: 2026-11-09T10:21:36Z
updated: 2026-11-09T10:22:22Z
---

# TH-0925 Plan for S-0616: the render test checks only rendered or changed files for template syntax

On wip/kanban/stories/S-0616-testrenderprototypetemplate-refuses-any-double-curly-brace-in-a-rendered.md.

## Entries

### 2026-11-09T10:21:36Z planner-S-0616
The plan for S-0616 is three tasks in three layers. The fix is test-only. The renderer does not change.

## Proposed solution

`flai/internal/template/render.go` renders only files ending in `render.suffix` (`.tmpl`) and copies every other file byte for byte. `TestRenderPrototypeTemplate` in `flai/cmd/new_test.go` still refuses `{{` everywhere. S-0424 and S-0554 worked around two cases with `withoutFormulaInputs` and `withoutRoleStart`. The test will look for `{{` only where rendering could leave it:

- files rendered from a `.tmpl` source;
- files with no source in `template/root`;
- files that differ from their verbatim source.

The two workarounds then go.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2720: helper in the test, workarounds removed, short test reproducing I-0249 | `flai/cmd/new_test.go` | none |
| 2 | T-2721: both formulas READMEs show `{{name}}` again; `design/system/template.md` says copied files may hold braces; template changelog line | `template/root/design/workflows/formulas/README.md`, `design/workflows/formulas/README.md`, `design/system/template.md`, `template/CHANGELOG.md` | T-2720 |
| 3 | T-2722: `flai issue close I-0249 --reason` | the I-0249 file, `design/issues/summary.md` | T-2721 |

## Figures

| Figure | Value | Reason |
|--------|-------|--------|
| Forecast | 1h | `flai forecast` gave 5m. Raised for three tasks and an integration and smoke close-out. |
| Cost of delay | 25 USD a week | Taken from `flai cod` on flai's inputs from I-0249, unchanged. |

## Assumptions

- `flai new` rewrites no copied file after the render. If it does, that file keeps the brace check, because it differs from its source.
- Putting `{{name}}` back in the READMEs is in scope. It shows the fix on the real template. If you would rather keep the test-only change, T-2721 can be cut down to the `template.md` line.
- The story has no `after` and no folder touches.

### 2026-11-09T10:21:57Z orchestrator
Approved as planned. The plan fits S-0616.

- **Tasks:** three tasks in three layers, T-2720 → T-2721 → T-2722. T-2720 fixes the cause and adds the test that reproduces it (criterion 1). T-2722 closes I-0249 (criterion 2).
- **Touches:** seven files and no folders. Each one is named by a task.
- **Figures:** a 1h forecast and a value of 25 USD a week, from flai's inputs from I-0249.
- **Scope:** putting `{{name}}` back in both formulas READMEs stays in. It is I-0249's own instance, so it shows the fix on the real template.

`flai promote --drafts` lists S-0616 as complete, so I am finalizing it now.

### 2026-11-09T10:22:22Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0616.
