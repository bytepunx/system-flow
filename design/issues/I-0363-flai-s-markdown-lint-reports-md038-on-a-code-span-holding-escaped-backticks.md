---
id: I-0363
title: flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-13T16:00:01Z
last_reported: 2026-11-13T16:00:01Z
updated: 2026-11-13T16:00:01Z
---

# I-0363 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry

## Description
flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry

## Instances

### 2026-11-13T16:00:01Z
Story: S-0751.
S-0751's manual race tier (`flai test race`, 2026-11-13) failed at one test, flai/internal/mdlint/repo_test.go:76 TestRepositoryLintsClean: "mdlint reports, markdownlint does not: design/digest/D-0592-flai-stream-sync-json-prints-a-map-that-flai-cmd-stream-go-builds-by-hand-not.md:17 MD038/no-space-in-code Spaces inside code span elements [Context: \"`). But `\"]". Line 17 holds a code span with backslash-escaped backticks inside it (`RegeneratedAtHead []string \`json:\"regenerated_at_head\"\``), which flai/internal/mdlint splits where markdownlint does not. The entry is on main, written by the curator through flai digest, so flai's own writer accepted text its own lint then rejects. No race was reported. The close-out skips this test under flai verify, so it fails only the integration and race tiers outside a close-out.

## Remediation
