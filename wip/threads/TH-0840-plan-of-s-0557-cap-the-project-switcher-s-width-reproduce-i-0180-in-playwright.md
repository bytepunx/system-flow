---
id: TH-0840
title: "Plan of S-0557: cap the project switcher's width, reproduce I-0180 in Playwright, close I-0180"
anchor:
  path: wip/kanban/stories/S-0557-the-top-bar-s-project-switcher-has-no-maximum-width-so-with-two-projects.md
  item: S-0557
status: resolved
participants: [planner-S-0557, agent-S-0557, orchestrator]
created: 2026-11-07T03:47:18Z
updated: 2026-11-07T03:53:58Z
---

# TH-0840 Plan of S-0557: cap the project switcher's width, reproduce I-0180 in Playwright, close I-0180

On wip/kanban/stories/S-0557-the-top-bar-s-project-switcher-has-no-maximum-width-so-with-two-projects.md.

## Entries

### 2026-11-07T03:47:18Z planner-S-0557
S-0557 is planned: three tasks in three layers, 40m forecast, cost of delay 15.00 USD a week. No thread named `asked` exists on S-0557, so no operator request was found to plan against. This is the story's first plan.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2507 Cap the project switcher's width and let it shrink inside the top bar | `flaiover/src/lib/components/ProjectSwitcher.svelte`, `ProjectSwitcher.svelte.test.ts` | none |
| 2 | T-2509 Reproduce I-0180 in Playwright: two long-named projects at 390px leave the page no wider than the window | `flaiover/e2e/project-switcher.e2e.ts` | T-2507, the fix it checks |
| 3 | T-2510 Close I-0180 naming the switcher's width cap and the Playwright test | I-0180's file, `design/issues/summary.md` | T-2509, the test its reason cites |

## Proposed solution

- The select takes the width of its longest option.
- Cap its width with a Tailwind max width, wider from `sm`.
- Let its `span` and `label` shrink and wrap.
- Keep each full label in its option and in the select's `title`.

## Assumptions

- The test answers `/api/projects` in the browser with route interception, as `board-layout.e2e.ts` answers `/api/board`. This keeps a second project off the server for the other test files.
- Playwright's Chromium can run on the host where the story is worked (`scripts/flaiover-e2e.sh`). The e2e suite is not a `flai test` tier, so the close-out does not run it. T-2509 records its own runs in its Notes.
- No user docs change: the switcher behaves as before, only narrower.
- The forecast was raised from flai's 14m to 40m for the flaiover build and the Playwright run.
- The cost of delay value is flai's figure from the operator's input, unchanged.

No story to split, merge, add, or drop.

### 2026-11-07T03:49:33Z agent-S-0557
Plan taken as written: three tasks in three layers, T-2507, then T-2509, then T-2510. I read `ProjectSwitcher.svelte`, its unit test, `board-layout.e2e.ts`, and the top bar in `+layout.svelte`; the cause and the fix stand as proposed. I work T-2507 myself, since it is one component I have read, and hand T-2509's build and Playwright run to a task sub-agent.

### 2026-11-07T03:53:36Z orchestrator
**Plan approved.** Three tasks in three layers. The cap comes first, T-2507, then the Playwright reproduction, T-2509, then I-0180's close, T-2510. That fits S-0557's two criteria. agent-S-0557 has taken the plan as written and is working T-2507. Nothing here waits on the operator, so I resolve this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T03:53:58Z orchestrator
Resolved: Plan approved; agent-S-0557 took it as written and is working it.
