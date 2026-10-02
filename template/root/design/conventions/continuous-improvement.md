---
title: Continuous improvement
updated: 2026-09-15
audience: agent
order: 100
status: active
topics: [all]
---

# Continuous improvement

How to record recurring friction, defects, blockers, and inefficiencies so they can be measured and remediated. The folder and schema are in `design/system/continuous-improvement.md`.

## Rules

- Keep a list of items that interfere with progress, quality, and efficacy filed under `design/issues`, one file per issue
- Classifications:
  - `defect`: something is wrong
  - `blocker`: progress stopped until someone acts
  - `efficiency`: it works but costs time
  - `impression`: a suspected problem without a measurement
- Provide a high-level description with more detailed instances of the identified issue or impression
- Front matter:
  - increment the `count` each time the issue comes up
  - keep `first_reported` and `last_reported` timestamps
  - if possible, record `cost`: average wall-clock time one occurrence costs, as a duration like `20m`
- `summary.md`
  - in the issues folder
  - a table of open issues with count, average cost, total cost, and last occurrence
  - order by issue cost
  - update it with every occurrence
- When a story moves to review, for each issue it recorded or bumped:
  - check whether a story in the backlog or in progress already remediates it
  - if none does, create a `remediation` or `improvement` story in the backlog (MCP `item_new`, or `flai story new`)
  - the story links the issue document and recommends a solution, and its last criterion closes the issue

## When in doubt

- Record the impression as such rather than ignoring it or over-reporting.
- Increment similar existing issues and add the new instance to its body if it differs enough from earlier ones.

<!-- system-flow:end-of-baseline -->

## Project additions
