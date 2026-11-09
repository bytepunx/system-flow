---
id: S-0629
type: story
nature: remediation
title: flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again
status: ready
owner: alex
created: 2026-11-07T23:40:41Z
updated: 2026-11-09T17:18:47Z
transitions:
  - to: ready
    at: 2026-11-09T17:13:07Z
    by: orchestrator
tags: [flai, serve, claude-code]
topics: [permissions, usage-limit]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go, design/adrs/README.md, design/system/agent-adapters.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0296-flai-serve-records-a-claude-code-usage-limit-refusal-as-a-failed-permission.md, design/issues/summary.md]
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
      seconds: 318
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 220
          output: 752
          cache_read: 11549396
          cache_write: 437757
          cost: 0.795
        - model: claude-opus-5-5
          input: 264
          output: 24868
          cache_read: 15219594
          cache_write: 309040
          cost: 5.3977
    - kind: orchestrator
      seconds: 112
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 268
          cache_read: 1264844
          cache_write: 123844
          cost: 0.3666
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
cost_of_delay:
  inputs:
    penalty_per_week: 20
    by: orchestrator
    at: 2026-11-09T17:08:35Z
  value: 20
  by: planner-S-0629
  at: 2026-11-09T17:10:09Z
forecast:
  duration: 40m
  delivery: 2026-11-09T18:08:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 5, behind S-0611, S-0616, S-0619, S-0620 and S-0624."
  by: flai
  at: 2026-11-09T17:18:47Z
finalized:
  by: orchestrator
  at: 2026-11-09T17:12:24Z
---
# S-0629 flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again

## Goal

This story remediates [I-0296](../../../design/issues/I-0296-flai-serve-records-a-claude-code-usage-limit-refusal-as-a-failed-permission.md), "flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0296 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0296 is closed with `flai issue close I-0296 --reason` saying what fixed it

## Tasks
- T-2749 An ADR and the design say a permission_prompt check that ends on the harness's usage limit records no outcome and is checked again after the limit resets
- T-2750 A permission_prompt check that ends on the usage limit records nothing, opens no thread, and is checked again after the reset
- T-2751 The user and operator docs say a usage-limit refusal is checked again, and I-0296 is closed

## Notes

The analyzer found I-0296: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0296: analyzer 41s, 0.8156 USD.

### Planning

The proposed remedy is that a check refused on the harness's usage limit has no outcome. flai serve records nothing for it, opens no thread, and checks the version again after the limit resets. Every other failure stays final, as ADR-0106 § 2 decided. Because that rule is in an accepted ADR, T-2749 writes a new ADR first. The story's agent may choose otherwise when it pulls the story; the plan's thread on S-0629 names the alternatives.

Touches:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/serve/claudecheck.go`, `flai/internal/serve/claudecheck_test.go` | layout | `check`, `checkImage`, `conclude`, `tell`, and the gate in `consider` decide and record the outcome. The tests beside them reproduce I-0296. |
| `flai/internal/usage/log.go`, `flai/internal/usage/log_test.go` | layout | The ADR-0214 limit reader takes only file paths. The check holds its output in memory, so it needs a reader that takes bytes. |
| `design/adrs/README.md` | design | The new ADR refining ADR-0106 § 2 is indexed there. Its own file is named once `flai adr new` writes it. |
| `design/system/agent-adapters.md`, `design/system/flai-cli.md` | design | § Permissions, § Commands, and § The container harness (S-0508) describe the record and how it is checked again. |
| `docs/users/flai.md`, `docs/operators/index.md` | design | Both say a version is never checked again, passed or failed. |
| `design/issues/I-0296-…md`, `design/issues/summary.md` | goal and co-change | Criterion 2 closes the issue, and `flai issue close` updates the summary. |

`flai touches suggest` found nothing else that bears on the story. Its co-change list is the general churn of the design documents. No folder touch is kept.

Figures:

- **Forecast.** `flai forecast` gave 24m at size 12. I raised it to 40m for the ADR, the new reader, and the host and image check paths. Delivery is set a day out because the story is a draft until the operator finalizes it.
- **Cost of delay.** The value is 20 USD a week, as `flai cod` works it out, kept as given. It rests on the 20 USD a week of penalty the orchestrator set on TH-0932, as recommended there. No time lost is set: the analysis measured none.
- **After.** None. The story reads nothing that another open story writes. ADR-0214's limit reader and the ADR-0106 check are already on main.
