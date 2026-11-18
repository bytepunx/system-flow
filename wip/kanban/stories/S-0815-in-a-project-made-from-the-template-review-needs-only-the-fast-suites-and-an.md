---
id: S-0815
type: story
nature: feature
title: In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:11Z
updated: 2026-11-18T08:13:57Z
transitions: []
tags: [flai, template, flaiover]
topics: [verify, workflow]
touches: [flai/cmd/suites_e2e_test.go, flai/cmd/testdata/e2e-suites-board.md, flai/internal/serve/suites_e2e_test.go, flaiover/src/lib/server/suites.test.ts, docs/operators/runbooks/change-the-board.md, docs/users/conventions.md, design/system/workflow.md]
after: [S-0807, S-0808, S-0809, S-0810, S-0811, S-0812, S-0814, S-0756]
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
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 556
          cache_read: 252710
          cache_write: 5265
          cost: 0.1038
    - kind: orchestrator
      seconds: 68
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 132
          cache_read: 1026213
          cache_write: 22654
          cost: 0.3215
cost_of_delay:
  value: 63.98
  by: planner-S-0815
  at: 2026-11-16T23:49:55Z
forecast:
  duration: 45m
  delivery: 2026-11-18T19:16:00Z
  basis: "Its own forecast of 45m; 9th in the pull order with an in-progress limit of 3, behind S-0811, S-0834, S-0842, S-0809, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813 and S-0814."
  by: flai
  at: 2026-11-18T08:13:57Z
finalized:
  by: alex
  at: 2026-11-17T07:47:11Z
---
# S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites

## Goal

The end-to-end proof of [ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md), as S-0756 is of the epic: the operator's example of 2026-11-16, made in the workflow file alone.

## Acceptance criteria

- [ ] A test makes a project from the template, edits its workflow file so that the move to review needs `static` and `behavior`, adds an `integration` state between `review` and `done` whose `runs:` names `kluster: [deploy, checks]` and the `integration` and `smoke` suites, and sets `publish` to the two slow suites.
- [ ] In that project a story's close-out runs only the fast suites, its entry to `integration` runs the slow ones once with the cluster's steps skipped for want of a cluster, a failing tier there sends it back, and its acceptance from `integration` is not held for a change to main that only the fast suites' paths do not select.
- [ ] The same project with a cluster configured deploys and runs its checks on entering `integration`.
- [ ] Anything that does not follow the file is recorded as an issue or a story under E-0034.
- [ ] The user guide has a worked example of the file.

## Tasks
- T-3463 An integration test makes a project from the template with the suites workflow file and checks the close-out, the acceptance from integration, and the publish
- T-3466 A serve test enters the integration state in the suites project and checks the slow suites run once, the cluster steps are skipped or run, and a failing tier sends the story back
- T-3468 A flaiover test of the built flai reads the suites project's board and review, with the integration column and the suites each step needs and ran
- T-3470 The board runbook gives the suites workflow file as a worked example of verify, runs, and verify_at
- T-3472 Record each step that did not follow the suites workflow file as an issue or a story under E-0034, and say in workflow.md that ADR-0312's outcome holds

## Notes

- After S-0807: the proof moves a story by a transition's verify.
- After S-0808: the proof accepts a story under the file's gates.
- After S-0809: the proof's batch runs the file's suites.
- After S-0810: the proof enters a state with runs.
- After S-0811: the proof publishes with verify_at.publish.
- After S-0812: the proof's CI command reads verify_at.ci.
- After S-0814: the proof reads what the dashboard is given for the new state.
- After S-0756: the proof adds a state as S-0756 shows one can be added.

### Planning

- Touches, all files; the story declared none:
  - `flai/cmd/suites_e2e_test.go`, `flai/cmd/testdata/e2e-suites-board.md`: layout. `TestRenderPrototypeTemplate` in `flai/cmd/new_test.go` renders `../../template` with `runCLI` and skips under `-short`, so a test beside it runs in the integration tier, which `system-flow.yaml` selects for `flai/**`. The fixture is a copy of the template's file: `template/root/design/workflows/board.md` must equal flai's baseline byte for byte (D-0648), so the test edits its own project, never the template.
  - `flai/internal/serve/suites_e2e_test.go`: layout. The state's job and the cluster are driven in process by the labs in `flai/internal/serve`: `newVerifyLab` (`verifyjob_test.go`), `newKlusterLab` and `fakeKluster` (`kluster_test.go`), and `withChecks` (`checks_test.go`). Only that package can use them.
  - `flaiover/src/lib/server/suites.test.ts`: layout. A test importing `testing.ts` reads `board.get` and the review data from the built flai, as `inbox.test.ts` does; S-0814 gives what it reads.
  - `docs/operators/runbooks/change-the-board.md`, `docs/users/conventions.md`: design and co-change. The runbook is the guide S-0755 wrote for the workflow file, and the worked example goes there after `## Add a rule to a move`; the users' guide links to it from `### States` instead of copying it.
  - `design/system/workflow.md`: design, where the workflow file's section states the outcome.
  - Left out of `flai touches suggest`'s co-change list: `design/system/flai-cli.md` and `docs/users/flai.md` (65% each) change with commands, and this story changes none; `flai/internal/serve/kluster.go` (7%) is S-0810's.
- `after`: kept as declared, each for the output the Notes above name. S-0805 and S-0806, the suite labels and the file's keys, come through S-0807 and S-0810, which name them.
- Tasks, in three layers:
  - Layer 1: T-3463, the fixture and the `flai/cmd` test of the close-out, the acceptance from `integration`, and the publish.
  - Layer 2: T-3466 (the serve job, with no cluster and with one), T-3468 (the dashboard's reads), and T-3470 (the worked example). Each reads T-3463's fixture; they touch no path in common, so they run together.
  - Layer 3: T-3472 records what T-3463, T-3466, and T-3468 found, and states the outcome.
- Assumed: the transition from `integration` to `done` gives `verify: [static, behavior]`, so the acceptance's gates judge only the fast suites, as criterion 2 needs; and the template's tiers carry the suites `static`, `behavior`, `integration`, and `smoke` that S-0805 gives them.
- Forecast 45m, up from the 21m `flai forecast` gives for size 12 (5 criteria, 7 touches) at 101 s a unit over 6 done stories in its band. Like S-0756, which has 45m, it drives the close-out, a state's job, the cluster, publish, and the dashboard end to end, and diagnoses each failure it finds. The delivery is flai's for 21m, moved later by the 24m added.
- Cost of delay: 63.98 USD a week, as `flai cod` gives it once the forecast is 45m. That is this story's share of E-0034's 1500.00 USD a week, from the operator's `time_lost_per_cycle: 10h` on the epic, by its 45m over the 17h35m forecast for the epic's 20 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits. It moves as the sibling stories' planners change their forecasts.
