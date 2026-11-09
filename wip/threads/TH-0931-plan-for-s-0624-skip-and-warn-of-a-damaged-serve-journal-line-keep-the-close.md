---
id: TH-0931
title: "Plan for S-0624: skip and warn of a damaged serve journal line, keep the close-outs section"
anchor:
  path: wip/kanban/stories/S-0624-one-line-holding-a-nul-byte-in-flai-serve-s-journal-drops-the-whole-close-outs.md
  item: S-0624
status: resolved
participants: [planner-S-0624, orchestrator, agent-S-0624]
created: 2026-11-09T17:05:44Z
updated: 2026-11-09T17:11:19Z
---

# TH-0931 Plan for S-0624: skip and warn of a damaged serve journal line, keep the close-outs section

On wip/kanban/stories/S-0624-one-line-holding-a-nul-byte-in-flai-serve-s-journal-drops-the-whole-close-outs.md.

## Entries

### 2026-11-09T17:05:44Z planner-S-0624
**Question: cost of delay inputs.** S-0624 has no cost of delay inputs and no epic, and I-0290 measured no time lost per cycle, so `flai cod` cannot give a value.

Recommended: `time_lost_per_cycle: 5m`. Its cost is lost evidence, not blocked work. The close-out verdicts and durations are unavailable on this host until the journal is edited by hand. 5m is the smallest figure flai set on this analysis's sibling stories (S-0620).

Alternatives: `0` (nothing lost per cycle, so the value is 0), or `15m` (each analysis run that reads close-outs spends time working around the gap).

**Proposed fix.** `journalEntries` in `flai/internal/statsread/statsread.go` stops at the first line it cannot decode, so one torn write drops the whole section. With this fix:

- Each line has its NUL bytes removed before it is decoded.
- A line that still does not decode is skipped and counted.
- `flai stats` warns once, naming the count and the first skipped line, and keeps the section.
- The section is absent only when the journal cannot be opened or read.

This matches what `flai serve journal` already does. The JSON shape does not change. `metrics.md` changes only with an ADR, so the decision comes first.

**Tasks, in order.**

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2737: an ADR, plus `metrics.md` § Close-outs and `docs/users/flai.md` | `design/system/metrics.md`, `design/adrs/README.md`, `docs/users/flai.md` | none |
| 2 | T-2738: the reader change and a test that reproduces I-0290 | `flai/internal/statsread/statsread.go`, `statsread_test.go` | T-2737 |
| 3 | T-2739: close I-0290 with its reason | the I-0290 file, `design/issues/summary.md` | T-2738 |

**Assumptions.**

- No change to the journal writer. A whole line is appended in one write, and NUL padding after a host stop is the filesystem's doing.
- No new JSON field for skipped lines. A warning on stderr is enough, so the dashboard needs no change.
- The story keeps the folder touch `design/adrs`, because the ADR's file name is not known yet. That folder is in `claims.shared`.
- Forecast: 20m, kept as `flai forecast` gave it.

I would not split, merge, or drop anything.

### 2026-11-09T17:08:39Z orchestrator
**Set:** `time_lost_per_cycle: 5m` on S-0624, the figure you recommend. Neither the thread nor the story gives a reason for `0` or `15m`. I set the input only; your next run works out the value.

**Plan approved.** The three tasks cover both criteria in order: the ADR first, since `metrics.md` changes only with one; the reader change and its reproducing test; then the close of I-0290. The touches name the files each task changes, and the `design/adrs` folder touch is explained under `### Planning` and sits in `claims.shared`. The 20m forecast is flai's, kept with a reason.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T17:09:25Z orchestrator
Resolved: time_lost_per_cycle 5m set on S-0624 as recommended; plan approved

### 2026-11-09T17:09:29Z planner-S-0624
Value written: 12.50 USD a week on S-0624, as `flai cod` gives it from `time_lost_per_cycle: 5m`. The Planning note records it. Nothing else changed.

### 2026-11-09T17:11:19Z agent-S-0624
Resolved: Plan approved by the orchestrator and cost of delay set; S-0624's agent works the three tasks as planned.
