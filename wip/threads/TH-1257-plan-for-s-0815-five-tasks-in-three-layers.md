---
id: TH-1257
title: "Plan for S-0815: five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0815-in-a-project-made-from-the-template-review-needs-only-the-fast-suites-and-an.md
  item: S-0815
status: open
participants: [planner-S-0815, orchestrator]
created: 2026-11-16T23:51:41Z
updated: 2026-11-16T23:56:37Z
---

# TH-1257 Plan for S-0815: five tasks in three layers

On wip/kanban/stories/S-0815-in-a-project-made-from-the-template-review-needs-only-the-fast-suites-and-an.md.

## Entries

### 2026-11-16T23:51:41Z planner-S-0815
S-0815 is planned: 7 file touches, a forecast of 45m, a cost of delay of 63.98 USD a week, and five tasks in three layers. It stays a draft in the backlog.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3463: the fixture workflow file and a `flai/cmd` integration test of the close-out, the acceptance from `integration`, and the publish | `flai/cmd/suites_e2e_test.go`, `flai/cmd/testdata/e2e-suites-board.md` | none |
| 2 | T-3466: a serve test of entering `integration`, slow suites once, cluster steps skipped (no cluster, action off, no namespace) or run (cluster configured), a failing tier sending the story back | `flai/internal/serve/suites_e2e_test.go` | T-3463 |
| 2 | T-3468: a flaiover test of the built flai reading `board.get` and the review data for the `integration` state and its suites | `flaiover/src/lib/server/suites.test.ts` | T-3463 |
| 2 | T-3470: the worked example in the board runbook, linked from the users' guide | `docs/operators/runbooks/change-the-board.md`, `docs/users/conventions.md` | T-3463 |
| 3 | T-3472: record each step that did not follow the file as an issue or a story under E-0034; state the outcome in `design/system/workflow.md` | `design/system/workflow.md` | T-3463, T-3466, T-3468 |

Layer 2's tasks share no path, so they run together.

## Assumptions

- The `integration` to `done` transition gives `verify: [static, behavior]`. That way the acceptance's gates judge only the fast suites, as criterion 2 needs.
- The template's tiers carry the suites `static`, `behavior`, `integration`, and `smoke`, from S-0805.
- The test edits its own copy of the workflow file, never the template's. The template's file must equal flai's baseline (D-0648).
- The serve job and the cluster are tested in `flai/internal/serve`, with that package's labs (`newVerifyLab`, `newKlusterLab`, `fakeKluster`, `withChecks`). Only that package can use them.
- "The user guide" (criterion 5) is read as the operators' runbook `docs/operators/runbooks/change-the-board.md`, which S-0755 wrote for the workflow file. `docs/users/conventions.md` links to it rather than copying it.
- T-3463 also checks that the publish runs the slow suites, since criterion 1 sets `publish` and no other criterion tests it.

## Figures

- Forecast 45m, up from `flai forecast`'s 21m for size 12. The reason is the same as S-0756's 45m: it drives every step end to end and diagnoses each failure.
- Cost of delay 63.98 USD a week, from `flai cod`. It is this story's share of E-0034's 1500.00 USD a week, from the operator's `time_lost_per_cycle: 10h`. It will move as the sibling planners change their forecasts.

## Proposals

- `after` names S-0809 (the batch) and S-0812 (CI), but no criterion exercises the batch or `flai test --at ci`, and the file sets neither `batch` nor `ci`. Either add a criterion that runs a batch and `flai test --at ci` in the project, or drop the two from `after`. I recommend dropping them, so this story waits only on what it reads; I changed neither.
- `flai check` advises `wip.after-hub` on S-0807, which S-0808, S-0814, and S-0815 name. S-0815 does read S-0807's output, the move by a transition's verify, so I kept it.

No input the operator owns is missing.

### 2026-11-16T23:56:37Z orchestrator (recommendation)
**Recommended:** drop S-0809 and S-0812 from S-0815's `after`, and drop their two lines from its Notes. Keep the criteria as they are. The planner makes the change on its next run.

**Why:**

- `work-management.md` says a story names in `after` only a story whose output it reads. No criterion of S-0815 runs a batch or `flai test --at ci`, and its fixture sets neither `batch` nor `ci`.
- S-0809 and S-0812 each prove their own point in their own tests.
- Dropping them takes nothing out of what S-0815 proves today. The other six stories in `after` stay, since the proof reads each one's output.

**Why I am asking you:** the Goal calls S-0815 the end-to-end proof of ADR-0312, and ADR-0312 §6 includes the batch and CI. Whether the proof must cover them is a scope decision.

**Alternative:** keep both in `after` and add a criterion that runs a batch and `flai test --at ci` in the made project. That adds a task and raises the forecast.

The rest of the plan fits S-0815: five tasks over five criteria, 7 file touches, a 45m forecast, and flai's value of 63.98 USD a week. I leave S-0815 a draft until this is settled. Its `after` holds it in the backlog either way, so the answer delays nothing.

Source: design/conventions/work-management.md
