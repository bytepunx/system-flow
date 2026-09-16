---
id: T-060
type: task
nature: feature
title: "issues package: parse, write, list, IDs, new, bump, close, summary"
status: done
parent: S-027
owner: alex
created: 2026-09-16T23:54:24Z
updated: 2026-09-16T23:59:08Z
transitions:
  - to: ready
    at: 2026-09-16T23:59:08Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:59:08Z
    by: agent
  - to: done
    at: 2026-09-16T23:59:08Z
    by: agent
stream: S-027
tags: [cli, issues]
---

# T-060 issues package: parse, write, list, IDs, new, bump, close, summary

## Work
internal/issues: Issue schema, Parse and Marshal in the hand-written style, List, Get, NextID, New (count 1, both timestamps, first instance), Bump (count, last_reported, running average cost, instance inserted before Remediation), Close (status closed, reason appended), Summary and WriteSummary (average and total cost, most expensive first), SummaryTable for embedding, Validate.

## Done when
Unit tests cover the lifecycle, averaging, ordering, closing, and validation.

## Notes
