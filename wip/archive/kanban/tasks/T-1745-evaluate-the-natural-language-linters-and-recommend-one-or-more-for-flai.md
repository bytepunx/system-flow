---
id: T-1745
type: task
nature: research
title: Evaluate the natural language linters and recommend one or more for flai
status: cancelled
parent: S-0413
owner: alex
created: 2026-10-16T05:45:32Z
updated: 2026-10-16T06:01:28Z
transitions:
  - to: cancelled
    at: 2026-10-16T06:01:28Z
    by: agent-S-0413
stream: S-0413
tags: [prose, research, lint]
touches: [design/system/agent-prose.md]
after: [T-1742]
---
# T-1745 Evaluate the natural language linters and recommend one or more for flai

## Work

Add a section on linters to `design/system/agent-prose.md`, meeting the story's first criterion.

- Candidates: Vale, proselint, write-good, textlint, and coala, which the story names, and any others found on the way, such as LanguageTool or alex.
- Judge each in one table on what this challenge needs: markdown awareness (front matter, code, links left alone), rules a project can write to encode `communication.md`, offline use, a runtime flai can ship or call from Go, output as data with a score usable to train and grade a model, speed on a narrative, licence, and maintenance in the last year.
- Run the strongest two or three on T-1742's sample in a scratch folder, never installed on the host, and record their findings and timings.
- Recommend one or more, with the reason, the rule set to start from, and where it would run: on what flai writes to `wip/` beside its markdown lint (ADR-0061), as a `tests` tier, or both.

Waits for T-1742: it needs the sample and the target style, and it writes the same document.

## Done when

- The linters section has the comparison table, the trial's results on the sample, and a recommendation with its reason.
- Each candidate's version and source are named, so the comparison can be checked.

## Notes
- 2026-10-16T06:01:28Z: moved to cancelled: Duplicates agent-S-0413's plan (T-1750, T-1755, T-1758, T-1760), as the planner recommends on TH-0472
