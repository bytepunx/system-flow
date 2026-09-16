---
title: Continuous improvement
updated: 2026-09-15
status: active
---

# Continuous improvement

`design/issues/` records the friction, defects, blockers, and inefficiencies that agents and operators hit while working, with a count and an average cost per occurrence, so remediation can be prioritised by evidence. Decided in [ADR-0014](../adrs/0014-design-issues.md); the agent-facing rules are in [conventions/continuous-improvement.md](../conventions/continuous-improvement.md).

## Layout

```text
design/issues/
├── README.md              # what the folder is
├── summary.md             # generated table of open issues
└── I-001-slug.md          # one file per issue
```

## Issue file

```yaml
---
id: I-001
title: golangci-lint on the host is v1 but the config is v2
class: efficiency          # defect | blocker | efficiency | impression
status: open               # open | closed
count: 3
cost: 5m                   # average wall-clock per occurrence, Go duration
first_reported: 2026-09-15T17:00:00Z
last_reported: 2026-09-15T20:10:00Z
updated: 2026-09-15T20:10:00Z
---

# I-001 golangci-lint on the host is v1 but the config is v2

## Description
What goes wrong, for whom, and what it costs.

## Instances
### 2026-09-15T17:00:00Z
Where it happened and what was done about it.

## Remediation
Proposed fix, or the story ID once one exists. Closed issues say what closed them.
```

| Class | Meaning |
|-------|---------|
| `defect` | Something is wrong and produced a wrong result |
| `blocker` | Progress stopped until someone acted |
| `efficiency` | It works but costs time every occurrence |
| `impression` | A suspected problem without a measured instance yet |

## Summary

`summary.md` is a table of open issues: ID, class, title, count, average cost per occurrence, total cost (count times average), last reported. Sorted by total cost descending so the most expensive friction is first. Every `flai issue` command regenerates it; `flai check` warns when it is stale, and `flai prime` prints it after the conventions when anything is open.

## Cadence

- Record an occurrence when it happens, not at the end of the story.
- When a story moves to review, the report includes the summary table if it changed during the story.
- When an epic completes, the summary is presented and the operator decides whether to spend time on remediation. Remediation is a story with nature `remediation` or `improvement` that closes the issue.

## Metrics

Issue counts and costs are inputs to the process-improvement loop in [workflow.md](workflow.md) alongside the flow metrics: they show where time goes outside the value stream.
