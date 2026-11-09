---
id: S-0624
type: story
nature: remediation
title: One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats
status: review
owner: alex
created: 2026-11-07T23:40:09Z
updated: 2026-11-09T17:22:51Z
transitions:
  - to: ready
    at: 2026-11-09T17:09:51Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T17:10:26Z
    by: agent-S-0624
  - to: review
    at: 2026-11-09T17:22:51Z
    by: agent-S-0624
tags: [flai, stats]
topics: [metrics]
touches: [flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, design/system/metrics.md, design/adrs, design/adrs/README.md, docs/users/flai.md, design/issues/I-0290-one-line-holding-a-nul-byte-in-flai-serve-s-journal-drops-the-whole-close-outs.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 754
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 1
      work: 66
  models:
    - model: claude-opus-5-5
      input: 150
      output: 21206
      cache_read: 18309694
      cache_write: 295796
      cost: 6.453
  strategic:
    - kind: orchestrator
      seconds: 91
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 30
          output: 207
          cache_read: 1576922
          cache_write: 113014
          cost: 0.4458
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-11-09T17:08:32Z
  value: 12.5
  by: planner-S-0624
  at: 2026-11-09T17:09:16Z
forecast:
  duration: 20m
  delivery: 2026-11-09T17:35:00Z
  basis: "flai forecast: median 120 s per unit of size over 44 done remediation stories on claude-opus-5-5 in the medium band, times size 10 (2 criteria, 8 touches), kept as given."
  by: planner-S-0624
  at: 2026-11-09T17:05:24Z
finalized:
  by: orchestrator
  at: 2026-11-09T17:09:47Z
---
# S-0624 One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats

## Goal

This story remediates [I-0290](../../../design/issues/I-0290-one-line-holding-a-nul-byte-in-flai-serve-s-journal-drops-the-whole-close-outs.md), "One line holding a NUL byte in flai serve's journal drops the whole close-outs section of flai stats", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0290 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0290 is closed with `flai issue close I-0290 --reason` saying what fixed it

## Tasks
- T-2737 An ADR and metrics.md § Close-outs say a serve journal line that is not an entry is skipped and warned of, and the section stays
- T-2738 flai stats keeps the close-outs section when a serve journal line holds NUL bytes or is not an entry, skipping and warning of it
- T-2739 I-0290 is closed, naming the ADR and the change that keeps the close-outs past a damaged journal line

## Notes

The analyzer found I-0290: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0290: analyzer 42s, 0.8156 USD.

### Planning

Proposed fix, from the one instance (line 874 of the host journal, an `orchestrate` entry preceded by NUL bytes from a write the host's stop cut short). `journalEntries` in `flai/internal/statsread/statsread.go` stops at the first line it cannot decode, so `CloseOuts` drops the whole section. Instead, each line has its NUL bytes removed before it is decoded, and a line that still is not an entry is skipped and counted. `CloseOuts` keeps the entries it read and warns once, naming how many lines it skipped and the first one's number. The section is absent only when the journal cannot be opened or read at all. `flai serve journal` already skips a line it cannot decode, so both readers then agree. `metrics.md` is the contract with the dashboard, so the change to the `closeouts` row comes with an ADR. The JSON shape stays the same: no new field, so the dashboard has nothing to change.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/statsread/statsread.go` | layout | `journalEntries` and `CloseOuts` live here |
| `flai/internal/statsread/statsread_test.go` | co-change, layout | `TestCloseOutsReadsTheServeJournal` is here; the reproducing test goes beside it |
| `design/system/metrics.md` | design | § Close-outs (S-0546) says an unreadable journal leaves the section absent |
| `design/adrs` | layout | Folder touch kept: the new ADR's file name is set only when `flai adr new` allocates its number, so no task can name it yet. The folder is in `claims.shared`, so it holds no ready story |
| `design/adrs/README.md` | co-change (37%) | The ADR index flai regenerates |
| `docs/users/flai.md` | co-change (26%) | Its close-outs paragraph repeats that an unreadable journal leaves the section out |
| `design/issues/I-0290-…-close-outs.md` | criteria | Closed by criterion 2 |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Left out: `flai/cmd/stats.go`, `flai/internal/metrics/closeouts.go`, and `flai/internal/metrics/metrics.go` (co-change 23% and 33%). The warning is logged in `statsread`, and `JournalEntry` and the report's shape do not change. Also left out: `flai/cmd/serve_actions.go`, whose `journal` reader already skips a bad line and whose writer appends a whole line in one write. A torn write after a host stop is the filesystem's, and no writer change prevents it.

Forecast: 20m, as `flai forecast` gave it once the touches were written: median 120 s per unit of size over 44 done remediation stories on claude-opus-5-5, times size 10. Kept as given. Two of the eight touches are files flai regenerates and one is a shared folder, which inflates the size. Against that, the close-out runs the integration and smoke tiers for a story under `flai/`. The two roughly cancel. Before the touches, it gave 5m on size 2.

Cost of delay: 12.50 USD a week, as `flai cod` gave it. That is 5m of time lost per 168h cycle at 150 USD an hour, 1.00 cycle a week. The input `time_lost_per_cycle: 5m` is the one recommended on TH-0931, which the orchestrator set. Kept as given: the cost is lost evidence, not blocked work.

No `after`: the story reads nothing another open story writes.
