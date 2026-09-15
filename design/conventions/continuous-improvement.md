---
title: Continuous improvement
updated: 2026-09-15
audience: agent
order: 100
status: active
---

# Continuous improvement

How recurring friction, defects, blockers, and inefficiencies are recorded so they can be measured and remediated instead of re-experienced. The folder and schema are in `design/system/continuous-improvement.md`.

## Rules

- Keep a list of items that interfere with progress, quality, and efficacy filed under `design/issues`, one file per issue
- Classify each issue as `defect` (something is wrong), `blocker` (progress stopped until someone acts), `efficiency` (it works but costs time), or `impression` (a suspected problem without a measured instance yet), with further qualification in the body
- Provide a high-level description with more detailed instances of the identified issue or impression
- Increment the `count` in front matter each time the issue comes up, and keep `first_reported` and `last_reported` timestamps so frequency is visible
- Where possible, record `cost` in front matter: the average wall-clock time one occurrence costs, as a duration like `20m`
- Keep `summary.md` in the issues folder: a table of open issues with count, cost, and last occurrence. Update it with every occurrence.
- When a story moves to review, include the summary table in the report if it changed during the story. When an epic completes, present the table and ask the operator whether to spend time on remediation.

## When in doubt

- Record the impression as such rather than ignoring it or over-reporting.
- Increment an existing issue that is similar and add the new instance to its body if it differs enough from earlier ones to teach something.

<!-- system-flow:end-of-baseline -->

## Project additions
- Issues seeded on 2026-09-15 from the first sessions of this project; see `design/issues/summary.md`.
- Until S-027 lands, maintain the issue files and the summary by hand, following `design/system/continuous-improvement.md`.
