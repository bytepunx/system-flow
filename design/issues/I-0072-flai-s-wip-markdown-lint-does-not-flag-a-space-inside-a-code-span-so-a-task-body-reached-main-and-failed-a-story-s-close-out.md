---
id: I-0072
title: flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-04T21:31:32Z
last_reported: 2026-10-04T21:31:32Z
updated: 2026-10-04T21:31:32Z
---

# I-0072 flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out

## Description
flai's wip markdown lint does not flag a space inside a code span, so a task body reached main and failed a story's close-out

## Instances

### 2026-10-04T21:31:32Z
Story: S-0256.
2026-10-04: T-0822's body, archived by S-0211's acceptance (2330743), has the code span `- Trigger:` written with a space before its closing backtick, which markdownlint's MD038 rejects (line 40). flai's wip lint (flai/internal/mdlint) has no MD038, so it took the body without a finding, and S-0256's close-out stopped at the markdown lint outside the story. Fixed on main in 38d51a2 under TH-0017's answer, as I-0056 and I-0070 were. Remedy to consider: flai's mdlint checks MD038, spaces at either end inside a code span.

## Remediation
