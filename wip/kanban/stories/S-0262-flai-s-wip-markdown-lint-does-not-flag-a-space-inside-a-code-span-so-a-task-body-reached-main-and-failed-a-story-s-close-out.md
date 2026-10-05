---
id: S-0262
type: story
nature: remediation
title: flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out
status: ready
owner: alex
created: 2026-10-04T21:41:59Z
updated: 2026-10-05T04:04:28Z
transitions:
  - to: ready
    at: 2026-10-04T23:14:03Z
    by: alex
tags: [flai]
topics: [cli, go]
touches: [flai/internal/mdlint, design/system/flai-cli.md, docs/users/flai.md, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-04T21:41:59Z
  value: 25
  by: planner-S-0262
  at: 2026-10-05T03:13:30Z
forecast:
  duration: 20m
  delivery: 2026-10-05T04:32:00Z
  basis: "Its own forecast of 20m; 3rd in the pull order with an in-progress limit of 3, behind S-0257 and S-0244."
  by: flai
  at: 2026-10-05T04:04:28Z
finalized:
  by: alex
  at: 2026-10-04T23:14:01Z
---
# S-0262 flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out

## Goal

This story remediates [I-0072](../../../design/issues/I-0072-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md), "flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0072 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0072 is closed with `flai issue close I-0072 --reason` saying what fixed it

## Tasks
- T-0856 mdlint reports MD038, spaces inside a code span, as markdownlint does
- T-0857 The design and the user guide name MD038 among the rules flai lints, and I-0072 is closed

## Notes

Cost of delay inputs set by flai from I-0072. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-04T21:31:32Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0072's one instance: `flai/internal/mdlint` has no MD038, so it let through a code span with a space before its closing backtick. The fix adds MD038 to mdlint and checks it against markdownlint-cli2 on fixtures. This is how S-0240 added MD007 for I-0055 (3446ac5). Every path flai writes wip markdown through, `item_new`, `item_edit`, thread entries, and `flai check`, then refuses the body or reports it.

Touches:

- `flai/internal/mdlint` (design and co-change): I-0072 names the package. `flai touches suggest` on its rule files lists `doc.go`, `mdlint_test.go`, and `testdata/cases/expected.txt`, each changed with them in 2 of 3 commits, so the whole package is claimed.
- `design/system/flai-cli.md` (design): it lists the rules `flai check` lints beyond ADR-0061's, MD007 among them, and the `mdlint/` line of the layout.
- `docs/users/flai.md` (design): the user guide lists what flai's lint checks under `flai check`.
- `design/issues` (layout): the second criterion closes I-0072 there and regenerates `summary.md`.
- The story declared no touches, so none were kept. ADR-0061 is accepted and stays unchanged, as it did in S-0240. `scripts/mdlint-fixtures.sh` is run, not changed.

Figures:

- Forecast 20m, delivery 2026-10-05T03:40:00Z. `flai forecast` gave 9m: 85 s per unit of size, the median over 19 done remediation stories, times size 6. It is raised to the 21m S-0240 took for the same kind of change, because markdownlint-cli2 has to be run to regenerate the reference. The story is no longer held for want of touches, and its touches overlap no story in progress, so it can start once it is pulled. The delivery is now plus the duration and a few minutes for the pull.
- Cost of delay 25.00 USD a week stands, as `flai cod` worked it out from the operator's input: 10m lost per 168h cycle at 150 USD an hour. I-0056 and I-0070 come from the same lint gap, but each was a different rule, so they do not raise this rule's rate.
- Topics `cli` and `go` are added: both tasks change the Go CLI.
