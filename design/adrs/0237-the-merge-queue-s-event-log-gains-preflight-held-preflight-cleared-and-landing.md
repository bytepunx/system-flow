---
id: ADR-0237
title: "The merge queue's event log gains preflight_held, preflight_cleared, and landing_resumed lines, which metrics.md lists and flai stats counts in no metric"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0177, ADR-0234]
---

# ADR-0237 The merge queue's event log gains preflight_held, preflight_cleared, and landing_resumed lines, which metrics.md lists and flai stats counts in no metric

## Context

`design/system/metrics.md` is the contract between `flai stats` and the dashboard, and changes only with an ADR. Its section on the merge queue lists every kind of line of the queue's event log, `.flai-cache/merge-queue/events.jsonl` ([ADR-0177](0177-flai-stats-reports-per-week-the-stories-the-merge-queue-landed-or-sent-back.md)), and the fields each carries.

[ADR-0234](0234-before-it-builds-a-batch-the-merge-queue-runs-the-landing-s-preflight-on-the.md) has the queue hold itself on what its preflight finds on the main checkout, and land a kept batch again without a new verification. S-0549 writes three new kinds of line for that. `flai stats` counts none of them, and its readers skip a kind they do not know.

## Decision

The merge queue's event log gains three kinds of line, `preflight_held`, `preflight_cleared`, and `landing_resumed`, which `metrics.md` lists and which `flai stats` counts in no metric.

- `preflight_held`: the preflight began to hold the queue, or the findings of a hold changed. It carries `files` and `rules`, the files and the rules its findings name, sorted.
- `preflight_cleared`: a look found main clean and cleared the hold. It carries `held_seconds`, how long the hold stood.
- `landing_resumed`: a batch whose landing waited on main went on from the commits it verified. It carries `batch`, `stories` with the story it went on at, and `waited_seconds`, how long the landing waited.

This refines ADR-0177 and ADR-0234.

## Consequences

- `metrics.md` names every kind of line the log holds, so a reader of the log can tell every line.
- No metric and no chart changes. A later story that charts how long the preflight holds the queue has the figures in the log already.

## Alternatives considered

- **Leave the kinds out of `metrics.md`.** Lost: the contract would no longer list every kind of line in the log.
- **Count the holds in `flai stats` now.** Lost: no criterion of S-0549 asks for it, and a metric changes the dashboard's contract further than this story needs.
