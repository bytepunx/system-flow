---
id: ADR-0317
title: "flai stats reports a workflow's own timers under timers: each item's seconds in items[].timers, each group's distribution over the window in summary.timers, and each timer's definition in a top-level timers list, all left out when the workflow defines none beyond lead, cycle, queue, and review"
status: accepted
date: 2026-11-17
supersedes: []
superseded_by: []
refines: [ADR-0301, ADR-0305]
---

# ADR-0317 flai stats reports a workflow's own timers under timers: each item's seconds in items[].timers, each group's distribution over the window in summary.timers, and each timer's definition in a top-level timers list, all left out when the workflow defines none beyond lead, cycle, queue, and review

## Context

[ADR-0301](0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md) §5 lets an operator define timers in the workflow file. Each timer is read from an item's recorded transitions, so it applies to past items too. [ADR-0305](0305-lead-cycle-queue-and-review-time-and-time-in-state-are-readings-of-the-workflow.md) makes lead, cycle, queue, and review time readings of the baseline's timers `lead`, `cycle`, `queue`, and `review`. It keeps their field names, `lead_time_seconds`, `cycle_time_seconds`, `queue_time_seconds`, and `wait_review_seconds`, and says S-0752 reports a workflow's other timers beside them (§6). [metrics.md](../system/metrics.md) is the contract between `flai stats` and the dashboard, and changes only with an ADR. This ADR names where those timers go.

## Decision

**`flai stats --json` reports every timer the workflow defines, other than `lead`, `cycle`, `queue`, and `review`, under a key named `timers`, in three places. Each place is left out when the workflow defines no other timer.**

### 1. Which timers

A workflow's own timers are the timers its file defines, in file order, except `lead`, `cycle`, `queue`, and `review`. ADR-0305 §5 and §6 already report those four under their own field names, so listing them again would give one figure twice. On the baseline workflow there are none, and `flai stats` gives the same bytes as before.

### 2. The definitions: `timers`

The report gains a top-level list `timers`: one entry per own timer, with `name`, `title`, `types`, and `counts` as the file gives them, and the defaults filled in (ADR-0301 §5). The dashboard takes a timer's chart heading from its `title`.

### 3. Per item: `items[].timers`

Each item gains `timers`, a map from a timer's name to its value in seconds, as ADR-0305 §3 reads it from the item's `created` and transitions. It holds a timer only when the timer's `types` include the item's type and the timer has a value for the item. The map is left out when it is empty. Open items are included: a sum counts its open interval up to now, as `wait_review_seconds` does.

### 4. Over the window: `summary.timers` and `groups[].timers`

Each summary gains `timers`, a map from a timer's name to a distribution of the same shape as `cycle_time`: `count`, `p50_seconds`, `p85_seconds`, `max_seconds`, and `mean_seconds`. It holds every own timer whose `types` include the report's item type. The distribution takes the items of the group completed in the window, as ADR-0305 §4 says: with `counts: accepted`, those whose state is of kind `accepted`; with `counts: closed`, those of kind `cancelled` too. Only items that have a value count. A timer that takes the report's type but has no item to count gives a distribution of zeros, as `cycle_time` does.

### 5. The text report

`flai stats` without `--json` prints each own timer's distribution after the baseline's cycle time, lead time, and ready wait, under its title, in each summary it prints.

### 6. The host

`stats.get` answers the report `flai stats --json` prints, so the dashboard reads the same `timers` from flai and computes none itself. Both read the workflow the repository runs: its `design/workflows/board.md`, else the baseline.

## Consequences

- An operator who adds a timer to the workflow file sees it in `flai stats` and `stats.get` with no change to flai. Its values cover every item already on the board and in the archive.
- The baseline's figures, their field names, and the golden fixtures stay as they are. On the baseline workflow the three `timers` keys are absent.
- A workflow that redefines `lead`, `cycle`, `queue`, or `review` changes the baseline figure of that name and adds no `timers` entry for it.
- S-0754 draws the timers on the dashboard from these fields.

## Alternatives considered

- **List every timer, the baseline's four included.** The four would appear twice, under `timers` and under their own field names. The baseline's output would also change, and its fixtures with it.
- **Give each timer's window distribution in the top-level list.** The groups, by nature, type, or parent, would then have no distribution of their own. A summary already holds `cycle_time` per group, so `timers` sits beside it.
- **Name the section `workflow_timers`.** ADR-0301 calls them timers. `timers` is the name in the workflow file, in `flai stats`, and on the charts.
