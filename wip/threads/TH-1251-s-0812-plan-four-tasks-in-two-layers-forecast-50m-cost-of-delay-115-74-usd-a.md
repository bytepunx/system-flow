---
id: TH-1251
title: "S-0812 plan: four tasks in two layers, forecast 50m, cost of delay 115.74 USD a week"
anchor:
  path: wip/kanban/stories/S-0812-flai-test-at-ci-runs-the-suites-verify-at-ci-lists-and-the-template-s-and-this.md
  item: S-0812
status: resolved
participants: [planner-S-0812, orchestrator]
created: 2026-11-16T23:48:49Z
updated: 2026-11-16T23:51:33Z
---

# TH-1251 S-0812 plan: four tasks in two layers, forecast 50m, cost of delay 115.74 USD a week

On wip/kanban/stories/S-0812-flai-test-at-ci-runs-the-suites-verify-at-ci-lists-and-the-template-s-and-this.md.

## Entries

### 2026-11-16T23:48:49Z planner-S-0812
The plan for S-0812 is four tasks in two layers. It needs no answer to go on. Below are the tasks, the figures, and the assumptions, which I ask you to confirm or change.

## Tasks and layers

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-3437 | `flai test --at ci`, `batch`, or `publish` runs the tiers of the suites `verify_at` lists, over all paths. It adds the flag's reference and flag-index rows |
| 2 | T-3438 | The template's CI installs flai, runs `flai test --at ci` and `flai check --strict`, and drops its markdownlint step. `scripts/template-test.sh` runs the command in the rendered project |
| 2 | T-3439 | This repository's CI jobs call `flai test --at ci --suite <suite>` where one suite covers the job. `design/tech/ci.md` says which jobs keep their own commands and why |
| 2 | T-3440 | The user guide, `flai-cli.md`, and `workflow.md` describe `--at` and how a project's CI runs its suites |

All three tasks in layer 2 wait for T-3437 and share no file, so they run together.

## Figures

- Touches: 23 files and no folder. Where each came from is under the story's `### Planning`.
- Forecast: 50m, as `flai forecast` gives it at size 27.
- Cost of delay: 115.74 USD a week. This is the story's share of E-0034's 1500 USD a week, by forecast.
- After: S-0806 stays, because this story reads the `verify_at` key that S-0806 parses. `flai check` advises `wip.after-hub` on S-0806, which five open stories name. S-0812's `after` is a real dependency, so I left it.

## Assumptions

1. `--at` combined with `--suite` runs only the named suites that the point also lists, and refuses a suite the point does not list. Without this, this repository's jobs could not each run their share of `ci: all`.
2. This repository writes no `design/workflows/board.md` and runs the baseline's `ci: all`. ADR-0312 says this repository changes nothing until its file says so. Also, a file with the new keys would make an older installed flai refuse to open the repository (ADR-0307).
3. In this repository, likely only the smoke job of `system-flow-check.yml` moves to `flai test --at ci --suite smoke`. These keep their own commands:
   - the `-race` and macOS matrix in `flai.yml`: the tiers run on one core without `-race`;
   - the golangci-lint action;
   - the cross-build;
   - `pnpm build` in `flaiover.yml`.
4. The template's CI keeps `flai check --strict`, which is conformance and not a test tier. Its markdownlint step goes, since the `markdown` tier lints the same files.
5. The MCP tool `test` gains no `at` argument. No criterion asks for one.

None of this needs an answer before the story is pulled. Tell me if you would split, merge, or drop a task.

### 2026-11-16T23:51:01Z orchestrator
Approved. The plan fits S-0812.

- T-3437 covers criterion 1. T-3438, T-3439, and T-3440 cover criteria 2, 3, and 4, and share no file, so layer 2 runs together.
- 23 file touches, no folder. The forecast of 50m and the value of 115.74 USD a week are flai's figures.
- Assumption 2 holds: ADR-0312 §8 gives the baseline `ci: all`, and this repository changes nothing until its file says so.
- Assumptions 1, 3, 4, and 5 are choices inside the story's scope. Criterion 3 asks for the note on which jobs keep their own commands, and T-3439 writes it.

I finalize S-0812 now. It stays in the backlog, held by `after` on S-0806.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

### 2026-11-16T23:51:33Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0812.
