---
id: I-0363
title: flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-13T16:00:01Z
last_reported: 2026-11-13T16:00:01Z
updated: 2026-11-17T11:42:33Z
---

# I-0363 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry

## Description
flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry

## Instances

### 2026-11-13T16:00:01Z
Story: S-0751.
S-0751's manual race tier (`flai test race`, 2026-11-13) failed at one test, `TestRepositoryLintsClean` in flai/internal/mdlint/repo_test.go, line 76. It said mdlint reports what markdownlint does not: MD038 (no-space-in-code) at line 17 of design/digest/D-0592-flai-stream-sync-json-prints-a-map-that-flai-cmd-stream-go-builds-by-hand-not.md. That line holds a code span written with backslash-escaped backticks inside it, around a Go struct tag on `RegeneratedAtHead`. flai/internal/mdlint ends the span at an escaped backtick, and markdownlint does not. The entry is on main. The curator wrote it through flai digest, so flai's own writer accepted text that its own lint then rejects. The race tier reported no race. The close-out skips this test under flai verify, so it fails only the integration and race tiers outside a close-out. Quoting the finding's context in an issue note trips markdownlint's MD038 too, so this note leaves it out.

## Remediation

Story S-0836 remediates this issue, created from it at 2026-11-17T08:17:31Z.
Closed 2026-11-17T11:42:33Z: S-0836. The issue's account was wrong: mdlint gave no false positive. markdownlint-cli2 0.20.0 reports the same MD038 on D-0592's line 17 as the curator wrote it, since both tools close a code span at the backtick after a backslash, as CommonMark says; main's wording, hand-edited since, passes both. The fixture flai/internal/mdlint/testdata/cases/code-spans.md now holds both wordings (T-3529). The cause was that flai digest new and edit never linted an entry's markdown, so one markdownlint rejects reached main: they now refuse it with markdown.MDnnn (T-3558, ADR-0318).
