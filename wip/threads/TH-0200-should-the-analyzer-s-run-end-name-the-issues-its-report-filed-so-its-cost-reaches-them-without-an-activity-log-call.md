---
id: TH-0200
title: Should the analyzer's run end name the issues its report filed, so its cost reaches them without an activity_log call?
anchor:
  path: wip/kanban/stories/S-0227-the-analyzer-s-cost-is-recorded-on-the-issues-it-filed-and-the-stories-made-from-them.md
  item: S-0227
status: resolved
participants: [agent-S-0227, alex]
created: 2026-10-06T21:31:02Z
updated: 2026-10-06T22:03:19Z
---

# TH-0200 Should the analyzer's run end name the issues its report filed, so its cost reaches them without an activity_log call?

On wip/kanban/stories/S-0227-the-analyzer-s-cost-is-recorded-on-the-issues-it-filed-and-the-stories-made-from-them.md.

## Entries

### 2026-10-06T21:31:02Z agent-S-0227
Criterion 2 charges an analyzer activity to the issues its `activity_log` call names. But the analyzer runs once and ends: flai serve logs its run end with no items (`flai/internal/serve/analyze.go`, `logRunEndSaying(..., nil)`), and neither its prompt nor `analyzer.md` asks it to call `activity_log`. As written, nearly every analyzer run would go to the project total, and no issue would carry a cost.

Options:

1. **(Recommended)** At the analyzer run's end, flai names the issues whose instances name the report the run wrote (the `Report: <path>.` line S-0224 adds), and charges the run's activity to them, split evenly. An `activity_log` call naming issues is charged the same way. No prompt change, and the charge does not depend on the agent remembering a call.
2. Keep the charge to `activity_log` only, and add a line to the analyzer's prompt and `analyzer.md` asking it to call `activity_log` naming the issues it filed or bumped before it ends. This touches files S-0224 (in review) also changes.
3. Criteria only: `activity_log` alone, no prompt change. Most runs charge the project total.

I am going on with the parts that do not depend on this (the issue's `usage`, the carry-over to the story, `flai stats`). Option 1 needs S-0224's `Report:` lines on main to be useful; it reads them as text, so it compiles without S-0224.

### 2026-10-06T21:49:44Z alex
go with option 1

### 2026-10-06T22:03:19Z agent-S-0227
Done as option 1: at an analyzer run's end flai names the issues that name its report and splits the run's cost between them (issues.NamingReport, serve.analyzeEnded). Recorded in ADR-0100.

Resolved.
