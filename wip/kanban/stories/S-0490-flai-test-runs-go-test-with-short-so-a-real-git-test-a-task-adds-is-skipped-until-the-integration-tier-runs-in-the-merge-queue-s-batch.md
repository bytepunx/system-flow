---
id: S-0490
type: story
nature: remediation
title: flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch
status: backlog
owner: alex
created: 2026-10-23T20:20:26Z
updated: 2026-10-23T23:39:23Z
transitions: []
tags: [flai, tests, manifest]
topics: [testing, flai, cli]
touches: [system-flow.yaml, flai/internal/verify/select_test.go, design/conventions/tooling.md, design/conventions/code-quality.md, design/system/project-manifest.md, docs/users/flai.md, scripts/README.md, scripts/flai-test.sh, design/issues/I-0201-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 1666
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 131
          output: 2332
          cache_read: 45021714
          cache_write: 23453
          cost: 9.9073
cost_of_delay:
  inputs:
    time_lost_per_cycle: 40m
    by: flai
    at: 2026-10-23T20:20:26Z
  value: 100
  by: planner-S-0490
  at: 2026-10-23T21:21:33Z
forecast:
  duration: 1h
  delivery: 2026-10-24T15:30:00Z
  basis: "flai forecast gave 24m (median 117 s per unit of size over 30 done remediation stories, size 12) and delivery 2026-10-24T14:48:00Z as 27th in the pull order; raised to 1h because the three tasks span ten files in six documents whose wording must be read before it is changed, and the close-out's go-test on verify now runs its real-git tests; delivery moved by the same 36m."
  by: planner-S-0490
  at: 2026-10-23T21:23:34Z
finalized:
  by: alex
  at: 2026-10-23T23:39:23Z
---
# S-0490 flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch

## Goal

This story remediates [I-0201](../../../design/issues/I-0201-flai-test-runs-go-test-with-short-so-a-real-git-test-a-task-adds-is-skipped-until-the-integration-tier-runs-in-the-merge-queue-s-batch.md), "flai test runs go test with -short, so a real-git test a task adds is skipped until the integration tier runs in the merge queue's batch". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0201 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0201 is closed with `flai issue close I-0201 --reason` saying what fixed it

## Tasks
- T-2099 Run the go-test tier without -short on the changed packages, skipping the gateway tests, with a repository-manifest test that reproduces I-0201
- T-2100 Document that the go-test tier runs the changed packages' tests whole, in the conventions, the manifest design, the user docs, and the scripts' README
- T-2101 Close I-0201 with the reason naming the go-test tier's change

## Notes

Cost of delay inputs set by flai from I-0201. time_lost_per_cycle 40m: 20m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-23T18:27:27Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0490 on 2026-10-23. The remedy proposed, from the issue's two instances: the `go-test` tier of `system-flow.yaml` drops `-short` and adds `-skip ^TestGateway`, so `flai test <paths>` and the close-out's `go-test` step run the changed packages' tests whole while `integration` is left to the merge queue (ADR-0174). No flai code changes; `make test`, `scripts/test.sh`, and CI keep `-short`. The alternatives and the measured cost are on the plan's thread.

Touches, file by file, none a folder:

- `system-flow.yaml`: the goal (the `go-test` tier's command).
- `flai/internal/verify/select_test.go`: layout; the repository-manifest tests (`TestRepositoryManifestLintsChangedFlaioverFiles`) live there, and the reproducing test joins them.
- `design/conventions/tooling.md`: design; its project addition names "the short Go tests" (also co-change, 24% through its template copy).
- `design/conventions/code-quality.md`: co-change (14%) and design; its project addition says `-short` selects the tiers.
- `design/system/project-manifest.md`: design; the `tests` bullet describes `go-test` and the `queue` tiers.
- `docs/users/flai.md`: co-change (19%) and design; § Run the tests for what changed describes the cover of `go-test`.
- `scripts/README.md`: co-change (14%); says "the short Go tests".
- `scripts/flai-test.sh`: layout; its comment says `go-test` is `test.sh`'s.
- `design/issues/I-0201-….md` and `design/issues/summary.md`: the second criterion.

Left out of the co-change list: `template/CHANGELOG.md`, `template/root/**`, and `template/template.yaml` (the template's manifest has no `go-test` tier and its conventions do not say "short"), `design/conventions/git.md` and `work-management.md` (nothing in them changes), `flai/internal/verify/select.go`, `verify.go`, and `flai/cmd/root.go` (no code changes).

Figures: cost of delay 100 USD a week, `flai cod`'s figure from the operator's input (40m lost per 168h cycle at 150 USD an hour), unchanged. Forecast: `flai forecast` gave 24m and delivery 2026-10-24T14:48:00Z (size 12, 27th in the pull order); raised to 1h and 2026-10-24T15:30:00Z, since the three tasks span ten files in six documents whose wording is read before it is changed, and the close-out's `go-test` on `verify` now runs its real-git tests.
