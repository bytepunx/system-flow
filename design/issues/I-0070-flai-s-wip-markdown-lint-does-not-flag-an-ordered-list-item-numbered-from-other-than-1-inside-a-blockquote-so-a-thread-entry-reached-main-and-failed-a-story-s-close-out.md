---
id: I-0070
title: flai's wip markdown lint does not flag an ordered list item numbered from other than 1 inside a blockquote, so a thread entry reached main and failed a story's close-out
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-04T04:45:41Z
last_reported: 2026-10-04T04:45:41Z
updated: 2026-10-04T04:45:41Z
---

# I-0070 flai's wip markdown lint does not flag an ordered list item numbered from other than 1 inside a blockquote, so a thread entry reached main and failed a story's close-out

## Description
flai's wip markdown lint does not flag an ordered list item numbered from other than 1 inside a blockquote, so a thread entry reached main and failed a story's close-out

## Instances

### 2026-10-04T04:45:41Z
Story: S-0225.
2026-10-04: TH-0101's entry by agent-S-0209 quotes two steps of planner.md as '> 3. ...' and '> 6. ...'; markdownlint's MD029 (style 1/1/1) flags both (lines 32 and 36), flai took the entry without a finding, and the publish commit a93da1a carried it to main. S-0225's close-out would stop at the markdown lint on it, outside the story. Fixed on main in 05a13e1 by escaping the numbers, as I-0056 was under TH-0017's answer. Remedy to consider: flai's MD029 also checks ordered lists inside blockquotes.

## Remediation
