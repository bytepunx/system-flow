---
id: T-1749
type: task
nature: research
title: Explore alternatives that make flai's agents write less and plainer without a local model
status: cancelled
parent: S-0413
owner: alex
created: 2026-10-16T05:45:53Z
updated: 2026-10-16T06:01:29Z
transitions:
  - to: cancelled
    at: 2026-10-16T06:01:29Z
    by: agent-S-0413
stream: S-0413
tags: [prose, research]
touches: [design/system/agent-prose.md]
after: [T-1747]
---
# T-1749 Explore alternatives that make flai's agents write less and plainer without a local model

## Work

Add a section on alternatives to `design/system/agent-prose.md`, meeting the story's third criterion.

- Cover at least: tighter conventions and prompts, with examples of good and bad entries; length budgets per kind of text that flai enforces when it writes (thread entries, logs, notes, final reports); structured output in place of prose, such as fixed fields an entry fills; lint on write, where flai refuses or warns on prose findings as it does on markdown findings in `wip/` (ADR-0061); a rewrite by a hosted small model through an existing provider in place of a local one; summaries on read, where the dashboard or `inbox` shows a short form and keeps the full text; and writing fewer documents.
- For each, say what it would change in flai or the conventions, what it costs per run, and how much reading it would save, against T-1742's figures where they apply.
- Compare them with the linter and model route of T-1745 and T-1747 in one table.

Waits for T-1747: it compares against the model route, and it writes the same document.

## Done when

- The alternatives section names each alternative with its change, its cost, and its expected saving, and holds the comparison table.

## Notes
- 2026-10-16T06:01:29Z: moved to cancelled: Duplicates agent-S-0413's plan (T-1750, T-1755, T-1758, T-1760), as the planner recommends on TH-0472
