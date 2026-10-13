---
id: I-0139
title: "flai check finds `markdown.MD012` outside the story at close-out"
class: efficiency
status: closed
count: 5
first_reported: 2026-10-12T15:41:38Z
last_reported: 2026-10-12T16:21:18Z
updated: 2026-10-13T17:20:22Z
---

# I-0139 flai check finds `markdown.MD012` outside the story at close-out

## Description
flai check finds `markdown.MD012` outside the story at close-out

## Instances

### 2026-10-12T15:41:38Z
Story: S-0364.
flai check found outside the story:
`wip/archive/agents/S-0232.md`: MD012/no-multiple-blanks Multiple consecutive blank lines [Expected: 1; Actual: 2]

### 2026-10-12T15:42:13Z
Story: S-0366.
flai check found outside the story:
`wip/archive/agents/S-0232.md`: MD012/no-multiple-blanks Multiple consecutive blank lines [Expected: 1; Actual: 2]

### 2026-10-12T16:01:45Z
Story: S-0313.
flai check found outside the story:
`wip/archive/agents/S-0232.md`: MD012/no-multiple-blanks Multiple consecutive blank lines [Expected: 1; Actual: 2]

### 2026-10-12T16:06:55Z
Story: S-0234.
flai check found outside the story:
`wip/archive/agents/S-0232.md`: MD012/no-multiple-blanks Multiple consecutive blank lines [Expected: 1; Actual: 2]

### 2026-10-12T16:21:18Z
Story: S-0337.
flai check found outside the story:
`wip/archive/agents/S-0232.md`: MD012/no-multiple-blanks Multiple consecutive blank lines [Expected: 1; Actual: 2]

## Remediation

Story S-0376 remediates this issue, created from it at 2026-10-12T16:52:58Z.
Closed 2026-10-13T17:20:22Z: S-0376: flai stream answer (AnswerOpenQuestion) now leaves one blank line where it removes an open question and saves through the markdown lint guard, and flai archive writes each narrative it moves with every run of blank lines outside fenced code collapsed to one, so no archived narrative keeps an MD012 finding that no story's close-out can clear. TestAnswerOpenQuestionLeavesOneBlankLine and TestArchiveCollapsesANarrativesBlankRuns reproduce the S-0232 case.
