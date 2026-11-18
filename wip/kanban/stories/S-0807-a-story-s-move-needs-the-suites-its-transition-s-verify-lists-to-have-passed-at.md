---
id: S-0807
type: story
nature: feature
title: A story's move needs the suites its transition's verify lists to have passed at its head, and the close-out runs the suites its next transitions need
status: ready
parent: E-0034
owner: alex
created: 2026-11-16T23:42:04Z
updated: 2026-11-18T08:13:57Z
transitions:
  - to: ready
    at: 2026-11-17T07:47:01Z
    by: alex
tags: [flai, template]
topics: [verify, workflow, conventions]
touches: [flai/internal/verifyrecord/record.go, flai/internal/verifyrecord/record_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/record.go, flai/internal/verify/record_test.go, flai/internal/verify/suites.go, flai/internal/verify/suites_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/workflow/rules.go, flai/internal/workflow/rules_test.go, flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/cmd/move_workflow_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go, flai/internal/mcpserver/verify.go, scripts/close-out.sh, template/root/scripts/close-out.sh, flai/cmd/verify_baseline_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/change-the-board.md, template/CHANGELOG.md]
after: [S-0806]
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
      seconds: 2
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 1039
          cache_read: 439869
          cache_write: 10653
          cost: 0.2051
    - kind: orchestrator
      seconds: 19
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 52
          cache_read: 260334
          cache_write: 15858
          cost: 0.085
cost_of_delay:
  value: 140.92
  by: planner-S-0807
  at: 2026-11-16T23:49:09Z
forecast:
  duration: 1h30m
  delivery: 2026-11-18T10:42:00Z
  basis: "Its own forecast of 1h30m; 4th in the pull order with an in-progress limit of 3, behind S-0811, S-0834, S-0842, S-0809, S-0812 and S-0808."
  by: flai
  at: 2026-11-18T08:13:57Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:56:26Z
---
# S-0807 A story's move needs the suites its transition's verify lists to have passed at its head, and the close-out runs the suites its next transitions need

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §2 and §7. Today no rule reads the verify report: a close-out before review is convention, and it runs every tier the story's paths select.

## Acceptance criteria

- [ ] The rule `verified` is in the vocabulary: a story's move by a transition with `verify:` is refused, or warned of at `warn`, unless its last `flai verify` passed at its branch's head every tier of those suites its paths select.
- [ ] `flai verify <S-nnnn>` and the close-out job run the suites the transitions out of the story's state need, and `--all` every suite.
- [ ] The verify report records each suite's tiers and the head and base they passed at, and a later run reuses them as ADR-0156 reuses tiers.
- [ ] On the baseline, a story's close-out and its move to review run and require what they did before this story.
- [ ] The conventions, the workflow design, and the user guide say what a close-out runs.

## Tasks
- T-3454 The verify report records each suite's tiers with the head and base they passed at, merges suites per commit, and reuses them across a clean rebase
- T-3457 The rule verified refuses, or warns of at warn, a story's move by a transition with verify unless its last flai verify passed those suites at its branch's head
- T-3460 flai verify and the close-out job run the suites the transitions out of the story's state need, and --all every suite
- T-3462 On the baseline workflow a story's close-out runs and its move to review requires what they did before this story
- T-3465 The conventions, the workflow design, the CLI design, and the user and operator guides say what a close-out runs and what the verified rule needs

## Notes

- After S-0806: the rule and the close-out read the transition's verify key S-0806 parses.

### Planning

- After: S-0806 only. This story reads the `verify:` key on a transition, with its suites and severity, that S-0806 parses in `flai/internal/workflow/workflow.go`. It also reads the tier's `suite` S-0805 adds to `flai/internal/verify/verify.go`, but S-0806 already waits on S-0805, so S-0805 is not named again.
- Touches, all files, no folder:
  - Declared: none before this run.
  - Layout, from the explorer's read of the code: `flai/internal/workflow/rules.go` holds the rule vocabulary, and `flai/internal/workitem/rules.go` runs a transition's rules in `Repo.Move`. `flai/internal/verify/select.go`, `story.go`, and `record.go` choose the tiers, write the report, and reuse it. `flai/cmd/verify.go`, `flai/internal/serve/verifyjob.go`, `verifyjob_container.go`, and `flai/internal/mcpserver/verify.go` run the close-out job. `scripts/close-out.sh` and its template copy say what the close-out runs. The test files sit beside each of them.
  - New files, from the explorer's design: `flai/internal/verifyrecord/record.go` and its test, a leaf package for the report's types and reader. `verify` imports `workitem`, so `Repo.Move` cannot import `verify`. `flai/internal/verify/suites.go` and its test work out the suites that a state's out-transitions need. `flai/cmd/verify_baseline_test.go` holds the baseline regression test.
  - Design, from criterion 5 and ADR-0312's consequences: `design/system/workflow.md` (rules table, close-out), `design/system/flai-cli.md` § The close-out job, `design/conventions/work-management.md` and `delegation.md` with their template copies, `docs/users/flai.md` § Verify a story before review, `docs/operators/runbooks/change-the-board.md` (rule list), `template/CHANGELOG.md`, and `docs/users/flai-reference.md`, which `make flai-reference` regenerates.
  - Co-change: `flai touches suggest` found no path above 19% (`docs/operators/index.md`). None is kept, since these are the repository's hub documents and no criterion reaches them.
- Forecast: `flai forecast` gave 1h12m (110 s per unit, size 39: 5 criteria and 34 touches). I raised it to 1h30m. E-0034's S-0738, S-0741, S-0744, and S-0755 took 63 to 90 minutes of agent time against forecasts of 29m to 1h. This story also changes the verify record's format and its reuse, which every acceptance reads. Delivery moves by the same 18 minutes, to 2026-11-17T06:41Z.
- Cost of delay: 140.92 USD a week, as `flai cod` gives it. That is this story's share of E-0034's 1500 USD a week, from the operator's `time_lost_per_cycle` of 10h on the epic: 1h30m of the 15h58m forecast over its 20 open stories without inputs. Unchanged, since the story has no inputs of its own.
- S-0804, in progress, touches the verify package, the close-out, and the serve job too. No `after` is set for that: a trial merge judges the overlap (ADR-0230).
