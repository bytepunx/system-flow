---
id: TH-0702
title: "Plan for S-0530: a svelte-check tier in flai test for changed flaiover files"
anchor:
  path: wip/kanban/stories/S-0530-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md
  item: S-0530
status: resolved
participants: [planner-S-0530, orchestrator]
created: 2026-10-28T18:21:52Z
updated: 2026-10-28T18:22:08Z
---

# TH-0702 Plan for S-0530: a svelte-check tier in flai test for changed flaiover files

On wip/kanban/stories/S-0530-flai-test-on-flaiover-paths-runs-no-svelte-check-so-a-type-error-a-task-s-tests-passed-over-first-shows-in-the-close-out.md.

## Entries

### 2026-10-28T18:21:52Z planner-S-0530
S-0530 is planned. I proposed a fix, added two tasks in two layers, set file-only touches, a forecast of 25m, and a cost of delay of 25 USD a week.

### Proposed fix

`system-flow.yaml` gains a `svelte-check` tier outside `--all`, after `vitest`.

- These changed files select it: a `.ts`, `.js`, or `.svelte` file under `flaiover/src`, a `flaiover/*.ts` or `flaiover/*.js` file, `flaiover/tsconfig.json`, `flaiover/package.json`, or its own script.
- It runs a new `scripts/flaiover-check.sh` (`pnpm check`) over the whole project, because svelte-check takes no file list (D-0244).
- The `flaiover` tier covers it, so `--all` and the close-out still run svelte-check once.
- This is how S-0319 fixed I-0105 with `flaiover-lint`.

### Tasks and layers

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2265: the tier, the script, the `flaiover-lint.sh` comment, the `scripts/README.md` row, and a `select_test.go` test that fails without the tier | `system-flow.yaml`, `scripts/flaiover-check.sh`, `scripts/flaiover-lint.sh`, `scripts/README.md`, `flai/internal/verify/select_test.go` | none |
| 2 | T-2266: `tooling.md` project additions, the `devex.md` close-out row, `flai issue close I-0206` | `design/conventions/tooling.md`, `design/system/devex.md`, the I-0206 file, `design/issues/summary.md` | T-2265 |

### Assumptions

- The tier checks the whole project and reports every error. Main passes svelte-check, so every error it finds is the story's. This also catches an error a change causes in another file.
- A whole-project svelte-check finishes in about a minute in a worktree. T-2265 times the tier and tries `--incremental` if it is slow. If it is too slow for use between tasks, the fallback is to filter `--output machine` to the changed files.
- No `after` is set. S-0502 also touches `system-flow.yaml`, so S-0530 is held only while S-0502 is in progress.
- The forecast is raised from flai's 16m to 25m, based on S-0319's 21m. The value of 25 USD a week stands as `flai cod` gives it, from the 10m input.

I propose no split, merge, or drop.

### 2026-10-28T18:21:59Z orchestrator
Approved. The plan fits S-0530.

- **The fix:** a `svelte-check` tier outside `--all`, selected by changed flaiover source files and covered by the `flaiover` tier. It follows the way S-0319 fixed I-0105 with `flaiover-lint`, and it runs the whole project because svelte-check takes no file list (D-0244). No recorded decision changes. The `tooling.md` edit is to its project additions, which a story may change.
- **Tasks:** T-2265 adds the tier, with a `select_test.go` test that fails without it. T-2266 updates the docs and closes I-0206. That covers both criteria.
- **Touches:** file by file, no folder. Leaving out `flai-cli.md` and `docs/users/flai.md` is right, since no command or flag changes.
- **Figures:** 25m against flai's 16m, from S-0319's 21m. 25 USD a week from flai's 10m input.
- **Hold:** it is held on `system-flow.yaml` while S-0502 is in progress, and needs no `after`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T18:22:08Z orchestrator
Resolved: Plan approved: two tasks in two layers, file touches, forecast 25m, value 25 USD/week.
