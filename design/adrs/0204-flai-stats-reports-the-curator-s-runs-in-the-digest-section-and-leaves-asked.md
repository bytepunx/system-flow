---
id: ADR-0204
title: "flai stats reports the curator's runs in the digest section and leaves asked_again out of a week no curator run has reached"
status: proposed
date: 2026-10-26
supersedes: []
superseded_by: []
refines: [ADR-0200]
topics: [cli, dashboard, analysis, digest]
---

# ADR-0204 flai stats reports the curator's runs in the digest section and leaves asked_again out of a week no curator run has reached

## Context

[ADR-0200](0200-flai-stats-reports-per-week-the-threads-opened-to-the-operator-the-stories.md) has `flai stats` count, per week, the threads opened to the operator that asked what a digest entry already answered. A thread counts as a repeat only once the curator lists it in an entry's `asked_again` ([ADR-0199](0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md)). The section reads the threads and the entries alone, so it reports `asked_again` 0 both when the curator found no repeat and when no curator ever ran.

I-0207 found the second case on system-flow. The `curate` host action was on for another project only, no curator had run here, and S-0435 read 0 repeats in every week as if it were a measurement. Every project reads the same false 0 until its curator first runs, and the analyzer and the dashboard's Digest chart take it at its word. `design/system/metrics.md` changes only with an ADR.

## Decision

**`flai stats` reports the curator's runs in the `digest` section, from its activity document, and leaves `asked_again` out of a week no curator run has reached.** It refines ADR-0200, as `metrics.md` § Digest defines it.

- **A curator run.** An entry in the curator's activity document, `wip/agents/curator.md` ([ADR-0196](0196-the-curator-runs-behind-the-curate-host-action-and-takes-in-what-the-archive.md)), at the time it ended. A project with no such document has no run.
- **`digest.curator`.** `runs`, the curator's runs over its whole log, not the window, and `last_run`, when the newest ended, absent before the first. It is present in every report that has `digest`.
- **`digest.weeks[].curator_runs`.** The curator's runs that ended in the week.
- **A week no run has reached.** A week whose Monday is after the newest run's end, or every week when the curator has never run. Its `asked_again` is absent, as `per_story` is absent in a week with no completed story: no repeat could have been counted, so the gap is not 0. A week a run reached keeps `asked_again`, 0 included. Its threads opened after the run still read low until the next run, as ADR-0200 says.
- **The text.** `flai stats` prints `asked again –` in a week no run has reached, `curator runs n` on each week's line, and one line for the curator under the entries: `curator  runs 3 · last 2026-10-26T05:00:00Z`, or `curator  runs 0 · never run`.

## Consequences

- A project whose curator has never run reads no repeat count at all, rather than 0. The Digest chart marks those weeks, and the analyzer reads `curator` before it reads `asked_again`.
- `asked_again` is now optional in `digest.weeks[]`. A reader that took it as always present must treat its absence as not measured. A flai older than this release still sends 0 there, and sends no `curator`.
- A run reaches the whole of the week it ended in. A thread opened later that week reads as measured, though only the next run can count it.
- The count needs no new file. flai already reads the activity documents for the strategic section.
