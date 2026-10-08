---
id: TH-0347
title: "Plan for S-0319: a per-file flaiover-lint tier, three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0319-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md
  item: S-0319
status: resolved
participants: [planner-S-0319, orchestrator]
created: 2026-10-08T00:15:59Z
updated: 2026-10-08T00:16:07Z
---

# TH-0347 Plan for S-0319: a per-file flaiover-lint tier, three tasks in three layers

On wip/kanban/stories/S-0319-flai-test-never-runs-flaiover-s-prettier-eslint-or-svelte-check-so-formatting-faults-surface-only-at-the-close-out.md.

## Entries

### 2026-10-08T00:15:59Z planner-S-0319
S-0319 is planned. It stays a draft in the backlog, for you to finalize.

## Fix proposed

In both I-0105 instances (S-0212, S-0329), prettier faults got past `flai test`, which ran only vitest on flaiover files. The fix:

- add a `flaiover-lint` tier to `system-flow.yaml`'s `tests`, which `flai test` runs outside `--all`;
- the tier runs a new `scripts/flaiover-lint.sh` with prettier `--check` and eslint on the changed flaiover files.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1304: `scripts/flaiover-lint.sh` runs prettier and eslint on the files it is given | `scripts/flaiover-lint.sh` | none |
| 2 | T-1305: the `flaiover-lint` tier, written with `flai manifest set tests=…`, and a test in `flai/internal/verify/select_test.go` that loads this repository's manifest and fails without the tier | `system-flow.yaml`, `flai/internal/verify/select_test.go` | T-1304 |
| 3 | T-1306: `design/system/devex.md` and the project addition in `design/conventions/tooling.md` name the tier; I-0105 is closed | `design/system/devex.md`, `design/conventions/tooling.md`, the I-0105 file, `design/issues/summary.md` | T-1305 |

Each task waits for the one before it: the tier runs the script, and the documents and the issue's close describe the tier.

## Figures

- Forecast: 30m, delivery 2026-10-08T07:15:00Z. `flai forecast` gave 13m. I raised it because each task plants and reverts a prettier fault, and the first run installs flaiover's dependencies in the new worktree.
- Cost of delay value: 25 USD a week, as `flai cod` computes it from the 10m input.

## Assumptions

- svelte-check stays in the `all_only` `flaiover` tier and at the close-out. It checks the whole project and takes no file list, and neither instance was a type error. The alternative is a script that runs it whole and reports only the changed files. That would add its full run time to every `flai test` on flaiover.
- The new tier runs before `vitest`, because prettier on a few files is cheaper.
- `flai manifest set tests=<JSON>` is how the tier is written, so no hand edit of the manifest is needed.

## Your input (optional)

I-0105 now counts 2 occurrences, but `time_lost_per_cycle` is still the 10m flai set from the first. Raising it to 20m would make the value 50 USD a week. That input is yours to change; the plan does not depend on it.

### 2026-10-08T00:16:07Z orchestrator
Approved.

- T-1305's tier and its selection test remove I-0105's cause (criterion 1). T-1306 closes I-0105 (criterion 2).
- The tier is declared in the manifest's `tests` with `flai manifest set`, not by hand. Its rule goes below the marker in `tooling.md`, where project additions belong. No recorded decision changes.
- Keeping svelte-check whole-project in the `all_only` tier is reasonable: neither instance was a type error.
- The serial layers are explained. The touches are file by file. The 30m forecast and the 25 USD a week value stand.
- S-0313 also touches `system-flow.yaml`, so whichever of the two runs second is held until the first is accepted.

For the operator: I-0105 now counts 2 occurrences against flai's 10m input. A 20m input would give 50 USD a week. That input is yours, with `flai edit S-0319 --time-lost-per-cycle 20m`.

Source: design/conventions/tooling.md

Resolved: Plan approved by the orchestrator under plan_backlog_stories
