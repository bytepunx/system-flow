---
id: TH-1256
title: "Plan for S-0813: seven tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0813-a-tier-s-queue-and-when-are-read-as-the-suites-of-the-batch-and-the-move-to.md
  item: S-0813
status: resolved
participants: [planner-S-0813, orchestrator]
created: 2026-11-16T23:51:36Z
updated: 2026-11-16T23:56:22Z
---

# TH-1256 Plan for S-0813: seven tasks in four layers

On wip/kanban/stories/S-0813-a-tier-s-queue-and-when-are-read-as-the-suites-of-the-batch-and-the-move-to.md.

## Entries

### 2026-11-16T23:51:36Z planner-S-0813
S-0813 is planned: seven backlog tasks in four layers, with its touches, forecast, and cost of delay written. The reasons are under `### Planning` in its Notes.

## Tasks and layers

| Layer | Task | What it does | After |
|-------|------|--------------|-------|
| 1 | T-3450 | ADR: how `queue` and `when` map onto `verify:` and `verify_at.batch`, what of ADR-0174 and ADR-0233 it supersedes, and what `flai check` and `flai upgrade` do | none |
| 2 | T-3452 | `flai/internal/workflow/migrate.go`: works out the workflow lines a manifest's keys stand for | T-3450 |
| 3 | T-3455 | `flai check` gives the advisory `manifest.tests-queue` finding, with those lines (criterion 2) | T-3452 |
| 3 | T-3459 | while a project sets the keys, flai runs the workflow worked out from them; once its file names `batch` suites, it runs the file and ignores them (criterion 3) | T-3452 |
| 3 | T-3461 | `flai upgrade` writes the lines into a workflow file it adds to a project that has none | T-3452 |
| 4 | T-3464 | this repository: create `design/workflows/board.md`, remove the two keys, and compare what each step selects before and after (criterion 4) | T-3455, T-3459 |
| 4 | T-3467 | design, settings index, user guide, reference, close-out scripts, and the template changelog | T-3455, T-3459, T-3461 |

The three layer-3 tasks share no file, and neither do the two layer-4 tasks.

## Figures

| Figure | Value | Reason |
|--------|-------|--------|
| Forecast | 1h30m, delivery 2026-11-17T09:54:00Z | flai gave 1h3m. I raised it 27m for seven tasks in four Go packages and a close-out that runs `integration` and `smoke`. |
| Cost of delay | 137.06 USD a week | flai's share of E-0034's figure, from your `time_lost_per_cycle` of 10h. Not adjusted. |

The story keeps its `after` on S-0808 and S-0809, each of which writes code this story builds on. S-0807 reaches it through S-0808.

## Assumptions

1. The title says "flai upgrade writes them into a project's file". ADR-0307 says upgrade never writes a project's own file, and S-0806's criterion 4 says the same. I planned T-3461 to write the lines only into a workflow file upgrade adds to a project that has none. A project with its own file gets them from `flai check`'s advice. If you meant upgrade to change a project's own file, the ADR in T-3450 has to supersede part of ADR-0307. Say so here and I will re-plan T-3461.
2. Criterion 4 may not hold if the tiers' `paths` stand in for `when`. ADR-0312's `verify:` lists suites with no paths. Here, `smoke` selects `scripts/**/*.sh` and `integration` selects `scripts/integration.sh`, but their `when` is only `flai/**` and `template/**`. With `paths` standing in, a script-only story's close-out would start running both. My recommendation is that T-3450's ADR lets a `verify:` entry carry its own `when` (`{suite, when}`). T-3452 then parses it, and the migration is exact. T-3464 stops and asks here if the selections differ.
3. `flai/internal/manifest/manifest_tests.go` keeps parsing both keys, and its `manifest.tests` error stays, because criterion 3 needs a project that sets them to keep running as before. "Retired" means flai ignores the keys once the file names `batch` suites, and the check advises removing them. The code is not deleted.
4. The tasks read the code S-0805 to S-0809 put on main, such as a tier's `suite`, the parsed `verify:`, and the loader those stories choose. T-3459 names `flai/internal/workflow/load.go` as the likely place, and its agent may move it to where those stories settle.

## Proposals

None to split, merge, or drop.

### 2026-11-16T23:55:50Z orchestrator
Approved. The plan fits S-0813.

- The seven tasks cover the four criteria. The three layer-3 tasks share no file, and neither do the two layer-4 tasks.
- 30 file touches, no folder. The forecast is 1h30m, raised from flai's 1h3m for the reasons you give. The value of 137.06 USD a week is flai's figure.
- Assumption 1: your reading stands. ADR-0307 §4 says `flai upgrade` adds the workflow file only to a project that lacks it, and never writes it once the project has one. S-0806's criterion 4 says the same. T-3461 stays as planned.
- Assumption 2: deciding how `when` maps is T-3450's ADR, which criterion 1 asks for. Criterion 4 asks that nothing this repository runs changes, so the ADR should take the mapping that keeps it the same. T-3464 stops and asks if the selections differ.
- Assumptions 3 and 4 are readings inside the story's scope.

I finalize S-0813 now. It stays in the backlog, held by `after` on S-0808 and S-0809.

Source: design/adrs/0307-flai-reads-the-workflow-file-of-the-checkout-it-opens-refuses-to-open-a.md § 4. Upgrade keeps the project's file

### 2026-11-16T23:56:22Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0813; upgrade writes only a file it adds, per ADR-0307 §4.
