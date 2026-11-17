---
id: T-3557
type: task
nature: remediation
title: TestPublishFormulaUnderJudgement reads the publish formula the template ships, which main's checkout no longer holds
status: done
parent: S-0827
owner: alex
created: 2026-11-17T08:45:40Z
updated: 2026-11-17T08:51:58Z
transitions:
  - to: ready
    at: 2026-11-17T08:45:48Z
    by: agent-S-0827
  - to: in-progress
    at: 2026-11-17T08:45:48Z
    by: agent-S-0827
  - to: done
    at: 2026-11-17T08:51:58Z
    by: agent-S-0827
stream: S-0827
tags: []
touches: [flai/cmd/publish_judgement_test.go]
usage:
  source: log
  seconds: 370
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1137
      cache_read: 709247
      cache_write: 23063
      cost: 0.3491
---
# T-3557 TestPublishFormulaUnderJudgement reads the publish formula the template ships, which main's checkout no longer holds

## Work

Main's 4052fc603 removed `design/workflows/formulas/publish.md` as a stopgap for S-0837. So `TestPublishFormulaUnderJudgement` in `flai/cmd/publish_judgement_test.go` cannot open it, and every story whose change selects `./cmd` fails `go-test` (I-0401). S-0827's change selects `./cmd`. Point `publishFormulaFile` at the copy the template ships, `template/root/design/workflows/formulas/publish.md`. That copy is the baseline formula the test means. S-0837 retires that copy and this test together. The task waits for no other task: it touches no path that T-3523 or T-3524 touch.

## Done when

`flai test` on `flai/cmd/publish_judgement_test.go` and `flai/internal/verify/select_test.go` passes its `go-test` tier.

## Notes

`flai test` on both paths passed every tier, `go-test` included (5m20s).
