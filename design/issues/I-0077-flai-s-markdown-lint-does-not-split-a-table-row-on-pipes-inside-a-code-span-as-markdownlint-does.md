---
id: I-0077
title: flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does
class: defect
status: closed
count: 1
cost: 6m
first_reported: 2026-10-05T03:58:53Z
last_reported: 2026-10-05T03:58:53Z
updated: 2026-10-09T18:29:17Z
---

# I-0077 flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does

## Description
flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does

## Instances

### 2026-10-05T03:58:53Z
Story: S-0253.
design/system/flai-cli.md:52, a table row, wrote the merge-base marker as a code span of seven unescaped pipes. markdownlint-cli2 in the close-out's smoke tier split the row on them (MD056, MD038, 23 errors); flai's mdlint, in TestRepositoryLintsClean, passed it, so the close-out stopped one tier later than it could have. Fixed by escaping each pipe.

## Remediation

Story S-0289 remediates this issue, created from it at 2026-10-06T09:56:52Z.
Closed 2026-10-09T18:29:17Z: S-0289: flai's mdlint now has MD056, table column count, so a row whose unescaped pipes, those inside a code span among them, give it more cells than its header is reported as markdownlint-cli2 reports it (TestTableColumnCountOfI0077, testdata/cases/tables.md).
