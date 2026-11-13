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
S-0751's manual race tier (`flai test race`, 2026-11-13) failed at one test, `TestRepositoryLintsClean` in flai/internal/mdlint/repo_test.go, line 76. It said mdlint reports what markdownlint does not: MD038 (no-space-in-code) at line 17 of design/digest/D-0592-flai-stream-sync-json-prints-a-map-that-flai-cmd-stream-go-builds-by-hand-not.md. That line holds a code span written with backslash-escaped backticks inside it, around a Go struct tag on `RegeneratedAtHead`. flai/internal/mdlint ends the span at an escaped backtick, and markdownlint does not. The entry is on main. The curator wrote it through flai digest, so flai's own writer accepted text that its own lint then rejects. The race tier reported no race. The close-out skips this test under flai verify, so it fails only the integration and race tiers outside a close-out. Quoting the finding's context in an issue note trips markdownlint's MD038 too, so this note leaves it out.

## Remediation
