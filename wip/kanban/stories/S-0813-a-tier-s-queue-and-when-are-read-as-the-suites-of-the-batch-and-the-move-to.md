---
id: S-0813
type: story
nature: feature
title: A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:09Z
updated: 2026-11-17T09:04:31Z
transitions: []
tags: [flai, template]
topics: [verify, workflow]
touches: [design/adrs/README.md, design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md, design/adrs/0233-while-the-merge-queue-can-batch-the-close-out-runs-a-queue-tier-whose-when.md, flai/internal/workflow/migrate.go, flai/internal/workflow/migrate_test.go, flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go, flai/internal/workflow/load.go, flai/internal/workflow/load_test.go, flai/internal/check/check_layout.go, flai/internal/check/tests_test.go, flai/internal/verify/manifest.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, system-flow.yaml, design/workflows/board.md, design/conventions/tooling.md, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, design/system/project-manifest.md, design/system/workflow.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0808, S-0809]
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
    - kind: planner
      seconds: 19
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 3503
          cache_read: 1976477
          cache_write: 38906
          cost: 0.7453
    - kind: orchestrator
      seconds: 32
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 45
          cache_read: 428946
          cache_write: 13449
          cost: 0.1356
cost_of_delay:
  value: 137.06
  by: planner-S-0813
  at: 2026-11-16T23:48:49Z
forecast:
  duration: 1h30m
  delivery: 2026-11-17T14:00:00Z
  basis: "Its own forecast of 1h30m; 11th in the pull order with an in-progress limit of 3, behind S-0752, S-0806, S-0831, S-0833, S-0754, S-0811, S-0809, S-0756, S-0812, S-0808, S-0807, S-0810 and S-0241."
  by: flai
  at: 2026-11-17T09:04:31Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:56:27Z
---
# S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §8. `queue: true` leaves a tier to the batch and `when:` runs it at the close-out anyway; with the file naming the suites of each step, the two keys say the same thing twice.

## Acceptance criteria

- [ ] An ADR says how `queue` and `when` map onto `verify:` and `verify_at.batch`, and what of ADR-0174 and ADR-0233 it supersedes.
- [ ] `flai check` advises a manifest that sets `queue` or `when` of the lines to put in the workflow file in their place.
- [ ] A project that sets neither key and names suites in its file runs what the file says; one that still sets them runs what it ran before, until it removes them.
- [ ] This repository's manifest and file are migrated with no change to what any step runs.

## Tasks
- T-3450 An ADR maps a tier's queue and when onto the workflow file's verify and verify_at.batch, and says what of ADR-0174 and ADR-0233 it supersedes
- T-3452 The workflow package derives, from a manifest's queue and when, the verify and verify_at.batch lines that stand in their place
- T-3455 flai check advises a manifest that sets queue or when of the workflow file's lines to write in their place
- T-3459 While a project still sets queue or when, flai runs the workflow derived from them, and once its file names batch suites it runs the file and ignores the keys
- T-3461 flai upgrade writes the lines a project's queue and when stand for into the workflow file it adds to a project that has none
- T-3464 This repository's manifest drops queue and when and its workflow file names the suites in their place, with no change to what any step runs
- T-3467 The manifest design, the settings index, the user guide, the reference, and the close-out scripts say queue and when are retired and name the workflow lines in their place

## Notes

- After S-0808: the acceptance's gates must read the file before queue stops deciding them.
- After S-0809: the batch must read the file before queue stops deciding what it runs.

### Planning

- `after`:
  - S-0808 writes the code that has the acceptance's gates (`queue_tiers`, `verify_behind_main`, and the queue's eligibility) read the accept transition's `verify:`. T-3459 makes the keys reach them only through that file.
  - S-0809 writes the code that has the batch read `verify_at.batch`. T-3459 and T-3461 retire the keys behind it.
  - S-0808 waits for S-0807, so S-0807's close-out reading of the transitions' suites is on main first. This story reads that too, and needs no `after` of its own for it.
- Touches. None was declared. `flai touches suggest` had nothing to start from, so each touch is from the design and the code layout, found by the explorer:
  - Design: ADR-0312 §8 and the criteria give the ADR files (`design/adrs/README.md`, ADR-0174, ADR-0233), `system-flow.yaml`, and `design/workflows/board.md`. The main checkout has no `design/workflows/board.md` yet, so T-3464 creates it.
  - Layout:
    - `queue` and `when` are parsed in `flai/internal/manifest/manifest_tests.go`. Nothing there changes, so it is not a touch.
    - They are copied in `flai/internal/verify/manifest.go`, selected in `flai/internal/verify/select.go` (`SelectStoryLeavingQueue`), and read by `flai/internal/verify/story.go` (`QueueTiersNotRun`).
    - The check findings are in `flai/internal/check/check_layout.go`, with its tests in `tests_test.go`. The workflow loader and parser are `flai/internal/workflow/load.go` and `workflow.go`, and the kept workflow file is in `flai/internal/upgrade/upgrade.go`.
    - `flai/internal/workflow/migrate.go` and `migrate_test.go` are new files that T-3452 names.
  - Docs: the documents that describe the two keys are `design/system/project-manifest.md`, `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/operators/settings.md` (the `tests[].queue` and `tests[].when` rows), `docs/users/flai.md`, `docs/users/flai-reference.md`, and the project additions in `design/conventions/tooling.md`. `scripts/close-out.sh` and its template copy name `queue: true`, which brings in `template/CHANGELOG.md`.
  - No folder touch. The new ADR's file name is not known until `flai adr new` gives it, and `flai task done` widens T-3450's touches with it.
- Forecast: 1h30m. flai gave 1h3m, at 110 s per unit of size over 142 large-band feature stories, size 34. It is raised 27m for three reasons:
  - seven tasks across four Go packages;
  - an ADR whose mapping of `when` decides whether criterion 4 can hold;
  - a close-out that runs `integration` and `smoke`, about 12 minutes, since the story touches `flai/` and `template/`.
- The delivery is flai's 09:27Z plus the same 27m.
- Cost of delay: 137.06 USD a week, as `flai cod` gives it after the forecast was written. The story has no inputs of its own. The figure is its share of E-0034's 1500.00 USD a week from the operator's `time_lost_per_cycle` of 10h, by its 1h30m of the 16h25m forecast over the epic's open stories without inputs. It stands unadjusted: no input of the operator's is missing.
